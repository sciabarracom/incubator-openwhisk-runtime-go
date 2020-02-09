# Debugger Support

ActionLoop supports debugging proving the following features:

- Debugging is globally enabled if the runtime receives the enviroment variable __OW_DEBUG_PORT

- Debugging is implemented asking runtimes to launch the action under a runtime and then performing an handshake to allow a client to connect to the action. 

The handshake works with way

- When an action is initialized, if the environment variable `__OW_DEBUG_IP` is set, then it will start a server in `__OW_DEBUG_PORT`+1

- When an action is run, the proxy expects a client listening in `__OW_DEBUG_PORT`+1, it will open a TCP connection to the  `__OW_DEBUG_IP` in port  and write  its own ip address (the first address of the available interfaces that is not localhost) as `IP:<IP>`. Anything else should be considered an error message. Then it waits for the target to close the connection

- A client (ideally embedded in the `wsk` client) will perform intialization of an action under debugger and then run