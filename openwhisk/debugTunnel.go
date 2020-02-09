package openwhisk

// StartForwarderIfDebugging will start a port fowarder, if required
// it will check if __OW_DEBUG_IP is defined,
// and then it will start a forwarder listening in port __OW_DEBUG_PORT+1
func StartForwarderIfDebugging() {

}

// WaitAndConnect will open a socket in the givent port,
// then wait to receive an IP
// and finally it will establish a connection back
func WaitAndConnect(port int) error {
	return nil
}

// EstablishTunnel will create a tunnel for port forwarding
// First, it will open back a connection to the given ip and port
// then it will write its own ip, then finally wait for the remote port to be closed before continuing
func EstablishTunnel(string, int) {

}
