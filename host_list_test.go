package portfwd

import (
	"fmt"
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomHostList(t *testing.T) {
	testHosts := []string{
		"localhost:1234",
		"localhost:5678",
		"localhost:1234",
	}
	assert := assert.New(t)
	l := NewRandomHostList()
	err := l.AddHost(testHosts[0])
	assert.NoError(err)
	err = l.AddHost(testHosts[1])
	assert.NoError(err)
	err = l.AddHost(testHosts[2])
	assert.Error(err)

	host, err := l.Host()
	assert.NoError(err)
	assert.NotEmpty(host)
}

// Test creating a single forwarder from :15000 to :25000
func TestForwarder(t *testing.T) {
	assert := assert.New(t)
	msg := "test\n"
	lHost := "localhost:15000"
	rHost := "localhost:25000"

	// listen for a message on rHost
	go runEchoListener(t, rHost, msg)

	hostList := NewRandomHostList()
	err := hostList.AddHost(rHost)

	AddForwarder(lHost, hostList)
	Start()

	conn, err := net.Dial("tcp", lHost)
	assert.NoError(err)
	_, err = fmt.Fprintf(conn, msg)
	assert.NoError(err)
	conn.Close()

	Stop()

}

func runEchoListener(t *testing.T, host, testMsg string) {
	assert := assert.New(t)
	l, err := net.Listen("tcp", host)
	assert.NoError(err)
	defer l.Close()
	lConn, err := l.Accept()
	// Accept single connection and read message
	assert.NoError(err)
	buf, err := ioutil.ReadAll(lConn)
	assert.NoError(err)
	bufStr := string(buf)
	t.Log(bufStr)
	assert.Equal(testMsg, bufStr)
}
