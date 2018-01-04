package portfwd

import (
	"errors"
	"math/rand"
	"sync"
)

// HostList is a collection of hosts we may be dialing
type HostList interface {
	// Host returns a host to connect to based on its own internal rules
	Host() (string, error)

	// AddHost adds a new host to the list
	AddHost(host string) error
}

// The base host list
type hostList struct {
	mu    sync.Mutex
	hosts []string
}

func newHostList() *hostList {
	return &hostList{
		hosts: make([]string, 0),
	}
}

func (l *hostList) addHost(host string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.containsHost(host) {
		return errors.New("Duplicate host: " + host)
	}
	l.hosts = append(l.hosts, host)
	return nil
}

func (l *hostList) containsHost(host string) bool {
	for _, lhost := range l.hosts {
		if lhost == host {
			return true
		}
	}
	return false
}

func (l *hostList) getHost(i int) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if i >= len(l.hosts) {
		return "", errors.New("Invalid host index")
	}
	return l.hosts[i], nil
}

// RandomHostList returns a host from its list randomly
type RandomHostList struct {
	*hostList
}

// NewRandomHostList initializes a new HostList returning hosts in random
// order
func NewRandomHostList() *RandomHostList {
	return &RandomHostList{
		hostList: newHostList(),
	}
}

func (l *RandomHostList) Host() (string, error) {
	if len(l.hosts) == 0 {
		return "", errors.New("No hosts available")
	} else {
		return l.getHost(rand.Intn(len(l.hosts)))
	}
}

func (l *RandomHostList) AddHost(host string) error {
	return l.addHost(host)
}
