package m2mserver
import (
    "log"
	"errors"
    "net"
)

func handShake(conn net.Conn)(string, error){
	if !sendAccept(conn) {
		return "", errors.New("SendAccept failed")
    }
    // Validate ESN
    uid, stat := receiveUid(conn)
    if !stat {
        return "", errors.New("Invalid ESN")
    }
    // send ESN_ACK
    if !sendUidAck(conn) {
        return "", errors.New("ESN Ack not sent")
    }
	return uid, nil
}

func sendAccept(conn net.Conn)(stat bool){
    bytesSent, err := conn.Write([]byte("ACCEPTED"))
    if err != nil || bytesSent != 8 {
        log.Print("Error sending connection confirmation")
        return false
    }
    return true
}

func receiveUid(conn net.Conn)(string, bool){
    // Make a buffer to hold incoming data.
    buf := make([]byte, 1024)
    
    // Read the incoming connection into the buffer.
    reqLen, err := conn.Read(buf)
    if err != nil {
        log.Print("Error reading:", err.Error())
        return "", false
    }
    if reqLen != 12 && reqLen != 14 {
        log.Print("Invalid Uid")
        return "", false
    }
	return string(buf[:12]), true
}

func sendUidAck(conn net.Conn)(stat bool){
    bytesSent, err := conn.Write([]byte("UID_ACK"))
    if err != nil || bytesSent != 7 {
        log.Print("Error sending UID ACK")
        return false
    }
    return true
}
