package portfwd

import (
	"errors"
	"math/rand"
)

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
