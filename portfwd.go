package portfwd

var (
	// Protocol to listen on
	proto = "tcp"

	// Map of listen addresses to host list for requests to this address
	listeners = make(map[string]HostList)
)
