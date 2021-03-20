package m2mserver

import "testing"

func TestConfigure(t *testing.T) {
	if std.initialized {
		t.Error("Expected initialized false, got true")
	}
	Configure(
		M2mConf{
			Ttl:  20,
			Port: "8000",
		},
	)
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
	if std.conf.Ttl != 20 {
		t.Error("expected Ttl = 20, got ", std.conf.Ttl)
	}
	if std.conf.Port != "8000" {
		t.Error("expected Port = 8000, got ", std.conf.Port)
	}
	if std.authenticate == nil {
		t.Error("expected authentication function, got none.")
	}
	if std.handlers == nil {
		t.Error("expected handlers function, got none.")
	}
}

// function used for testing purpose, see config test above.
func t_ack(c *Client, param string) {}

func TestNewClient(t *testing.T) {
	c := newClient("111111111111")
	if c.Uid != "111111111111" {
		t.Error("failed to initialized a client, expected uid 111111111111, got ", c.Uid)
	}
}

func TestSendCmd(t *testing.T) {
	c := newClient("111111111111")
	c.SendCmd(Cmd{Name: "test", Param: "123"})
	select {
	case m := <-c.msgQ:
		if m.Name != "test" || m.Param != "123" {
			t.Error("expected msg name test, got ", m.Name, " and msg param 123, got ", m.Param)
		}
	default:
		t.Error("message failed to send, msgQ is empty")
	}
}

func TestIsActive(t *testing.T) {
	c := newClient("111111111111")
	if !c.IsActive() {
		t.Error("expecting client active true, got false")
	}
}

func TestAddHandler(t *testing.T) {
	c := newClient("111111111111")
	std.clients["111111111111"] = c
	std.handlers(c)
	_, ok := c.handler["test_ack"]
	if !ok {
		t.Error("expected handler test_ack, got none.")
	}
}

func TestHasHandler(t *testing.T) {
	c := newClient("111111111111")
	std.clients["111111111111"] = c
	std.handlers(c)
	if !c.HasHandler("test_ack") {
		t.Error("expected hasHandler to return true, got false.")
	}
}

func TestGetClient(t *testing.T) {
	c := newClient("111111111111")
	std.clients["111111111111"] = c
	c, err := GetClient("111111111111")
	if err != nil || c.Uid != "111111111111" {
		t.Error("expected to get client 1111.., got ", c.Uid, err)
	}
}
