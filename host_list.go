package portfwd

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	// Names of strategies to hostlists
	strategies = make(map[string]func() HostList)
)

func init() {
	// Register available strategies
	strategies["random"] = func() HostList { return NewRandomHostList() }
	strategies["roundrobin"] = func() HostList { return NewRoundRobinHostList() }
}

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

func (l *hostList) String() string {
	return fmt.Sprint(l.hosts)
}

// HostListForStrategy returns a new HostList for the given strategy
func HostListForStrategy(strategyName string) (HostList, error) {
	f := strategies[strings.ToLower(strategyName)]
	if f == nil {
		return nil, errors.New("Invalid strategy: " + strategyName)
	}
	return f(), nil
}
