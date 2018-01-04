package portfwd

// HostList is a collection of hosts we may be dialing
type HostList interface {
	// Host returns a host to connect to based on its own internal rules
	Host() string
}
