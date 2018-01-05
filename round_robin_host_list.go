package portfwd

import (
	"errors"
	"sync"
)

// RoundRobinHostList returns hosts from its list in order
type RoundRobinHostList struct {
	*hostList

	idxMu sync.Mutex
	idx   int
}

// NewRandomHostList initializes a new HostList returning hosts in random
// order
func NewRoundRobinHostList() *RoundRobinHostList {
	return &RoundRobinHostList{
		hostList: newHostList(),
	}
}

func (l *RoundRobinHostList) Host() (string, error) {
	if len(l.hosts) == 0 {
		return "", errors.New("No hosts available")
	}
	l.idxMu.Lock()
	defer l.idxMu.Unlock()
	host, err := l.getHost(l.idx)
	if err != nil {
		return "", err
	}
	if l.idx >= (len(l.hosts) - 1) {
		l.idx = 0
	} else {
		l.idx++
	}
	return host, nil
}

func (l *RoundRobinHostList) AddHost(host string) error {
	return l.addHost(host)
}
