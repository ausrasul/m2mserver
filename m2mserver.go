package m2mserver
import (
    "net"
	"errors"
	"log"
	"sync"
	"time"
)

type M2mConf struct {
	Ttl	int
	Port		string
}

type M2m struct {
	initialized bool
	conf *M2mConf
	authenticate func(string)bool
	handlers func(*Client)bool
	onDisconnect func(string)
	clients map[string]*Client
}

var std = &M2m{
	initialized: false,
	conf: &M2mConf{
		Ttl: 10,
		Port: "7000",
	},
	clients: make(map[string]*Client),
}

func (m *M2m) Prime (auth func(string)bool, handlers func(*Client)bool, onDisconnect func(string)){
	if auth != nil {
		m.authenticate = auth
	}
	if handlers != nil {
		m.handlers = handlers
	}
	if onDisconnect != nil {
		m.onDisconnect = onDisconnect
	}
}

func (m *M2m) Configure (cnf M2mConf) {
	if cnf.Ttl > 0 {
		m.conf.Ttl = cnf.Ttl
	}
	if cnf.Port != "" {
		m.conf.Port = cnf.Port
	}
	std.initialized = true
}

func Configure(cnf M2mConf) {
	std.Configure(cnf)
}

func Prime(auth func(string)bool, handlers func(*Client)bool, onDisconnect func(string)){
	std.Prime(auth, handlers, onDisconnect)
}

type Cmd struct{
	Name string
	Serial string
	SchType int
	SchTime time.Time
	SchPeriod time.Duration
	Param string
}

type Client struct{
	active bool
	Uid string
	msgQ chan Cmd
	handler map[string]func(*Client, string)
}

func newClient(uid string) (*Client){
	return &Client{
		active: true,
		Uid: uid,
		msgQ: make(chan Cmd, 1),
		handler: make(map[string]func(*Client, string)),
	}
}

func (c *Client) SendCmd(cmd Cmd) bool{
	var mutex = &sync.Mutex{}
	log.Print("active stat: ", c.active)
	mutex.Lock()
	if !c.active {
		close(c.msgQ)
		mutex.Unlock()
		return false
	}
	log.Print("sending msg")
	c.msgQ <- cmd
	log.Print("sent msg")
	mutex.Unlock()
	return true
}

func (c Client) IsActive() bool{
	return c.active
}

func (c *Client) AddHandler(cmdName string, handler func(*Client, string)){
	c.handler[cmdName] = handler
}

func (c *Client) disconnect(){
	c.active = false
}

func (c *Client) HasHandler(cmdName string) bool{
	_, ok := c.handler[cmdName]
	return ok
}

func Listen(){
	
    // Listen for incoming connections.
	for {
		l, err := net.Listen("tcp", "0.0.0.0:" + std.conf.Port)
        if err != nil {
			log.Print("Error listening:", err.Error())
            continue
        }
        // Close the listener when the application closes.
        defer l.Close()
    
        for {
            // Listen for an incoming connection.
            conn, err := l.Accept()
            if err != nil {
                continue
            } else {
                // Handle connections in a new goroutine.
				go connHandler(conn)
            }
        }
    }
}

func GetClient(esn string) (*Client, error){
	client, ok := std.clients[esn]
	if !ok{
		return &Client{}, errors.New("No such client connected")
	}
	return client, nil
}
