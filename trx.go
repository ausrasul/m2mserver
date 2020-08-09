package m2mserver
import (
	"log"
	"net"
	"time"
	"bufio"
	"encoding/json"
	"strconv"
)

/*
func sender(conn net.Conn, outbox <-chan Cmd, stopSend <-chan bool, sendStopped chan<- bool){
	for {
		select{
		case <- stopSend:
			sendStopped <- true
			log.Print("sender stopping.")
			return
		case command:= <- outbox :
			msg, err := json.Marshal(command)
			if err != nil {
				log.Print("Error marshaling msg")
				sendStopped <- true
				return
			}
			//w   := bufio.NewWriter(conn)
			w := bufio.NewWriterSize(conn, len(msg) + 1)

			bytesSent, err := w.Write(append(msg, byte('\n')))
			if err != nil || bytesSent < len(msg) + 1{
				log.Print("Error buffering msg")
				sendStopped <- true
				return
			}
			err = w.Flush()
			if err != nil {
				log.Print("Error sending msg")
				sendStopped <- true
				return
			}
			log.Print("bytes sent: ", bytesSent)
		}
	}
}
*/
func sender(conn net.Conn, outbox <-chan Cmd, stopSend <-chan bool, sendStopped chan<- bool){
	for {
		select{
		case <- stopSend:
			sendStopped <- true
			log.Print("sender stopping.")
			return
		case command:= <- outbox :
			msg, err := json.Marshal(command)
			if err != nil {
				log.Print("Error marshaling msg")
				sendStopped <- true
				return
			}
			// sending message header, the data size.
			dataSize := []byte(strconv.Itoa(len(msg)))
			dataSize = append(dataSize, byte('\n'))
			w := bufio.NewWriterSize(conn, len(dataSize))

			bytesSent, err := w.Write(dataSize)
			if err != nil || bytesSent < len(dataSize) {
				log.Print("Error buffering msg len")
				sendStopped <- true
				return
			}
			err = w.Flush()
			if err != nil {
				log.Print("Error sending msg")
				sendStopped <- true
				return
			}
			log.Print("bytes sent (header): ", bytesSent)

			// sending the data
			w = bufio.NewWriterSize(conn, len(msg))

			bytesSent, err = w.Write(msg)
			if err != nil || bytesSent < len(msg){
				log.Print("Error buffering msg")
				sendStopped <- true
				return
			}
			err = w.Flush()
			if err != nil {
				log.Print("Error sending msg")
				sendStopped <- true
				return
			}
			log.Print("bytes sent: ", bytesSent)
		}
	}
}

func receiver(conn net.Conn, inbox chan<- Cmd, stopRcv <-chan bool, rcvStopped chan<- bool){
	for {
		select {
		case <- stopRcv:
			rcvStopped <- true
			log.Print("receiver stopping.")
			return
		default:
		}
		conn.SetDeadline(time.Now().Add(1e9))
		buf, err := bufio.NewReaderSize(conn, 65537).ReadString('\n')
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Print("Error reading:", err.Error())
			rcvStopped <- true
			return
		}
		log.Print(":", buf)
		var cmd Cmd
		err = json.Unmarshal ([]byte(buf) , &cmd)
		if err != nil {
			log.Print("Error unmarshaling msg:", err)
			rcvStopped <- true
			return
		}
		inbox <- cmd
	}
}
