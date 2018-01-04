package portfwd

import (
	"errors"
	"io"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// Protocol to listen on
	proto = "tcp"

	// Default Dial Timeout
	Timeout = 30 * time.Second

	forwardersMu sync.Mutex
	// Map of listen addresses to host list for requests to this address
	forwarders = make(map[string]HostList)

	// Listening connections by local host
	listeners = make(map[string]net.Listener)

	// Channel to signal shutdown
	chSig = make(chan bool, 1)
)

// Goroutine that connects and forwards traffic
func forwarder(src, dst net.Conn) {
	if src == nil {
		err := errors.New("No local connection")
		logErr("", err)
		return
	}
	if dst == nil {
		err := errors.New("No remote connection")
		logErr("", err)
		return
	}

	go func() {
		_, err := io.Copy(src, dst)
		if err != nil {
			logErr(dst.RemoteAddr().String(), err)
		}
	}()
	go func() {
		_, err := io.Copy(dst, src)
		if err != nil {
			logErr(dst.RemoteAddr().String(), err)
		}
	}()
}

// Add a listener
// TODO check if it already exists
func addListener(lhost string, hosts HostList) {
	var err error
	listeners[lhost], err = net.Listen(proto, lhost)
	if err != nil {
		logErr(lhost, err)
		return
	}
	go startListener(lhost)
}

// Start listener on specified host and forward connections
func startListener(lhost string) {
	listener := listeners[lhost]
	if listener == nil {
		logErr(lhost, errors.New("No listener at this address"))
		return
	}
	for {
		select {
		case <-chSig:
			listener.Close()
			return
		default:
			src, err := listener.Accept()
			if err != nil {
				logErr(lhost, err)
				continue
			}
			dst, err := dialForwardHost(lhost)
			if err != nil {
				logErr(lhost, err)
				continue
			}
			go forwarder(src, dst)
		}
	}
}

// find a host to dial out to and connect
func dialForwardHost(lhost string) (net.Conn, error) {
	rhost, err := getHostForListener(lhost)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTimeout(proto, rhost, Timeout)
	if err != nil {
		return nil, err
	}
	log.WithField("from", lhost).WithField("to", rhost).Info("Established new connection")
	return conn, nil
}

func getHostForListener(lhost string) (string, error) {
	hostList := forwarders[lhost]
	if hostList == nil {
		return "", errors.New("No registered listeners for " + lhost)
	}
	return hostList.Host()
}

// Log error
func logErr(host string, err error) {
	if err != nil {
		log.WithField("host", host).Error(err)
	}
}

// Start starts all forwarders
// No new forwarders may be added at this point
func Start() {
	log.Infof("Starting %d forwarders", len(forwarders))
	for lhost, hostList := range forwarders {
		addListener(lhost, hostList)
	}
}

// Stop shuts down all forwarders
func Stop() {
	log.Info("Shutting down...")
	close(chSig)
	// TODO proper waitgroup for shutdown
	time.Sleep(100 * time.Millisecond)
}

// AddForwarder creates a new forwarding listener
func AddForwarder(lhost string, hosts HostList) {
	forwardersMu.Lock()
	defer forwardersMu.Unlock()
	forwarders[lhost] = hosts
}
