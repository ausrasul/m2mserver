package m2mserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	assert.False(t, std.initialized, "Initialized is false in the beginning")
	assert.NotEqual(t, std.conf.Ttl, 20, "TTL is set")
	assert.NotEqual(t, std.conf.Port, "8000", "Port is set")
	Configure(
		M2mConf{
			Ttl:  20,
			Port: "8000",
		},
	)
	assert.True(t, std.initialized, "Initialized true after configure")
	assert.Equal(t, 20, std.conf.Ttl, "TTL is set")
	assert.Equal(t, std.conf.Port, "8000", "Port is set")
}

func TestPrime(t *testing.T) {
	assert.Nil(t, std.authenticate, "Doesn't have authentication function")
	assert.Nil(t, std.handlers, "Has handler function")
	Prime(
		func(uid string) bool { return true },
		func(c *Client) bool {
			if c.Uid == "111111111111" {
				c.AddHandler("test_ack", t_ack)
			}
			return true
		},
		func(uid string) { return },
	)
	assert.NotNil(t, std.authenticate, "Has authentication function")
	assert.NotNil(t, std.handlers, "Has handler function")
}

// function used for testing, see config test above.
func t_ack(c *Client, param string) {}

func TestNewClient(t *testing.T) {
	c := newClient("111111111111")
	assert.Equal(t, c.Uid, "111111111111", "failed to initialized a client")
}

func TestSendCmd(t *testing.T) {
	c := newClient("111111111111")
	cmd := Cmd{Name: "test", Param: "123"}
	c.SendCmd(cmd)
	select {
	case m := <-c.msgQ:
		assert.True(t, m.Name == cmd.Name && m.Param == cmd.Param, "Unexpected command sent")
	default:
		assert.True(t, false, "message failed to send, msgQ is empty")
	}
}

func TestIsActive(t *testing.T) {
	c := newClient("111111111111")
	assert.True(t, c.IsActive(), "Client is active")
}

func TestAddHandler(t *testing.T) {
	c := newClient("111111111111")
	std.clients["111111111111"] = c
	std.handlers(c)
	_, ok := c.handler["test_ack"]
	assert.True(t, ok, "handler equals test_ack")
}

func TestHasHandler(t *testing.T) {
	c := newClient("111111111111")
	std.clients["111111111111"] = c
	std.handlers(c)
	assert.True(t, c.HasHandler("test_ack"), "handler returns true")
}

func TestGetClient(t *testing.T) {
	c := newClient("111111111111")
	std.clients["111111111111"] = c
	c, err := GetClient("111111111111")
	assert.Nil(t, err, "No errors")
	assert.Equal(t, "111111111111", c.Uid, "Uid is correct")
}

func TestListens(t *testing.T) {
	// mock listner
}

/*what we should be testing is:
it configures
it primes
it accepts connections
it authenticate connections
it keeps heart beat
it sends commands of any size to buffer
it handles received commands of any size
*/
/*
how the flow of information should be?
configure set vars
set listeners, set vars
accept connections
authenticate, this is too complex. should be simpler, just call authenticate.
heartbeat, send and receive
receive commands,
send commands.
*/
