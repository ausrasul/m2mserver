package m2mserver

import (
	"log"
	"net"
	"time"
)

// Handles incoming requests.
func connHandler(conn net.Conn) {
	connTimeout := time.Duration(std.conf.Ttl)
	defer conn.Close()

	// Authenticate a client
	uid, err := handShake(conn)
	if err != nil {
		return
	}
	if !std.authenticate(uid){
		return
	}
	std.clients[uid] = newClient(uid)
	client, _ := GetClient(uid)
	std.handlers(client)
	clientStat := &std.clients[uid].active

	defer func(clientStat *bool, uid string) {
		*clientStat = false
		std.onDisconnect(uid)
		log.Print("closing client")
	}(clientStat, uid)

	msgQ := std.clients[uid].msgQ
	msgQ = msgQ

	outbox := make(chan Cmd, 1)
	stopSend := make(chan bool, 1)
	sendStopped := make(chan bool, 1)

	inbox := make(chan Cmd, 1)
	stopRcv := make(chan bool, 1)
	rcvStopped := make(chan bool, 1)

	defer close(outbox)
	defer close(stopSend)
	defer close(sendStopped)
	defer close(inbox)
	defer close(stopRcv)
	defer close(rcvStopped)

	go sender(conn, outbox, stopSend, sendStopped)
	go receiver(conn, inbox, stopRcv, rcvStopped)

	timer := time.NewTimer(time.Second * connTimeout)
	ok := true
	for ok {
		select {
		case <-timer.C:
			log.Print("No heartbeat from client.")
			*clientStat = false
			stopRcv <- true
			stopSend <- true
			<-rcvStopped
			<-sendStopped
			select {
			case <-msgQ:
			default:
			}
			log.Print("all stopped.")
			ok = false
			continue
		case cmd := <-inbox:
			log.Print("handle response")
			if cmd.Name == "hb" {
				log.Print("HB <--")
				timer = time.NewTimer(time.Second * connTimeout)
				cmd.Name = "hb_ack"
				outbox <- cmd
			} else {
				callback, ok := client.handler[cmd.Name]
				log.Print(ok)
				if !ok {
					log.Print("No handler for command ", cmd.Name)
				} else {
					callback(client, cmd.Param)
				}
			}
		case msg := <-msgQ:
			log.Print("Push notification: ", msg.Name)
			outbox <- msg
		case <-sendStopped:
			*clientStat = false
			log.Print("Sender stopped.")
			stopRcv <- true
			<-rcvStopped
			select {
			case <-msgQ:
			default:
			}
			ok = false
			continue
		case <-rcvStopped:
			*clientStat = false
			log.Print("Receiver stopped.")
			stopSend <- true
			<-sendStopped
			select {
			case <-msgQ:
			default:
			}
			ok = false
			continue
		}
	}

	log.Print("handler closing ")
}
