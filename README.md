# M2M Server

A Go Object Request Broker

This is a server client library that lets you send commands in json format to execute remote code
 to multiple clients/workers in the form of broadcast, multicast, individual.

## When to use Go ORB
- To run functions or OS commands on multiple clients and receive responses.
- To run the above as a scheduled task (like cron job)

## Features:
- Add clients at any time.
- Manage authentication yourself.
- Commands are in json format.
- Immediate or scheduled one time execution, or scheduled repeated task.
- Allowed commands by client and server are user defined.
- Command execution is user implemented.

## Work mechanism:

![alt text](https://raw.githubusercontent.com/ausrasul/m2mserver/master/diagram.png)

The server starts at a given port, and listen to connections.
When a client tries to connect, it is authenticated.
The client or the server can initiate communication (defined by the user)
When the server sends a command to client, the client lookup the command name in a map of handlers.
The corresponding handler is then called and the result is sent back together with the command tracking id (serial number)
The response sent back by the client is considered a "command" by the server, there is no difference.

The server can process or leave that command depends if the user have defined a handler for it.

Server have also a list of handlers mapped per command name.

The server and client have heartbeat, if that is broken, the client is considered disconnected.
The client will automatically try to reconnect after a given time.
The user can see how many clients are defined and how many clients are actually connected to the server.

## Dependencies
None

## Usage

```

package main

import (
	m "github.com/ausrasul/m2mserver"
	"log"
)

func main(){
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Configure the server

	m.Configure(
		m.M2mConf{
			Ttl:      10,  // Heartbeat timeout in seconds
			Port: "7000",  // Listening port number (IP is always localhost)
		}
	)

	// Load hooks
	m.Prime(
		onAuthenticate,  // Called when any client tries to connect.
		onConnect,       // Called right after onAuthenticate, here you can add commands handlers.
		onDisconnect,    // Called after a client is disconnected.
	)
	m.Listen()           // start the server.
}

func onConnect(c *m.Client)bool{
	log.Print("adding command handlers for ", c.Uid)  // c.Uid is the identification of the client.
	// A handler is a function that runs when a client sends a command to server.
	// Each client should be given its own handelr(s)
	// one handler per command.
	if c.Uid == "000c29620510" {
		log.Print("Adding a handler for command test_ack")
		c.AddHandler("test_ack", t_ack)
	}

	// You don't have to send any command here or now, this is just for demo.
	command := m.Cmd{
		Name: "test",          // command we're sending to client, hopefully the client has
		                       // a handler for it.
		Param: "test2341111",  // you can send any parameters with the command,
	}
	c.SendCmd(command)

	return true
}

func onAuthenticate(Uid string) bool {
	fmt.Println("Looking up auth database")
	fmt.Println("Client authenticated")
	// Note that we use Uid for authentication, TODO: we'll add a more secure method.
	return true
}

func onDisconnect(Uid string){
	fmt.Println("Client disconnected, cleaning up...")
	fmt.Println("Should print out Uid ", Uid)
}

func t_ack(c *m.Client, param string){
	fmt.Println("we got test_ack command, doing stuff with param: ", param)
}

```
