package openwhisk

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func ExampleTestHTTPServer() {
	ts := testHTTPServer("Hello")
	defer ts.Close()
	body, code, _ := doGet(ts.URL + "/hello")
	fmt.Print(1, code, " ", body)
	fmt.Print(2, " ", testHTTPServerLastRequest, "\n", testHTTPServerLastBody)
	body, code, _ = doPost(ts.URL+"/hellopost", `{"hello":"world"}`)
	fmt.Print(3, " ", code, " ", body)
	fmt.Print(4, " ", testHTTPServerLastRequest, "\n", testHTTPServerLastBody)
	// Output:
	// 1 200 Hello
	// 2 GET /hello
	// 3 200 Hello
	// 4 POST /hellopost
	// {"hello":"world"}
}

func ExampleTestTCPClient() {
	// start
	svr := exec.Command("go", "run", "_test/usrv.go", "127.0.0.1:9999")
	svr.Start()
	// connnect
	cli := exec.Command("go", "run", "_test/tcli.go", "127.0.0.1:9999", "hello")
	out, _ := cli.Output()
	fmt.Println(string(out))
	// stop
	exec.Command("go", "run", "_test/tcli.go", "127.0.0.1:9999", "*").Run()
	svr.Wait()
	// Output:
	// HELLO
}

func ExampleRequestReverseProxy() {
	fmt.Println("reverse proxy request")
	ts := testHTTPServer("test server")
	defer ts.Close()
	RequestReverseProxy(ts.URL, "/123", "http://1.2.3.4/")
	grep(`rule|url`, testHTTPServerLastBody)
	// Output:
	// reverse proxy request
	// "rule": "PathPrefix:/123"
	// "url": "http://1.2.3.4/"
}

func ExampleServerClient() {
	// server
	svr := exec.Command("bash", "-c", "echo hello | nc -l 9999")
	svr.Start()
	defer svr.Wait()
	// client
	cli := exec.Command("bash", "-c", "nc localhost 9999")
	out, _ := cli.CombinedOutput()
	fmt.Println(string(out))
	// Output:
	// hello
}

func ExampleForwarder() {
	// create a new forwarder in port 8079 that fowards to port 8888
	auth := "debug:pass"
	serverHost := "127.0.0.1"
	serverPort := 8079

	// server
	svr := exec.Command("go", "run", "_test/usrv.go", "127.0.0.1:9999")
	svr.Start()

	// create a forwarder server
	fwd, err := NewForwarder(auth, serverHost, serverPort)
	if err != nil {
		return
	}
	fwd.Server.Logger.Info = false
	fwd.Server.Logger.Debug = false

	fmt.Println("started server", serverHost, serverPort, fwd.Start() == nil)
	defer fwd.Server.Close()

	// create a forwarder client
	server := fmt.Sprintf("%s:%d", serverHost, serverPort)
	rev := "7777:127.0.0.1:9999"
	cli, err := ChiselClient(server, rev, auth)
	cli.Logger.Info = false
	cli.Logger.Debug = false

	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println("started client", server, rev, cli.Start(ctx) == nil)
	defer cancel()

	// result
	cmd := exec.Command("go", "run", "_test/tcli.go", "127.0.0.1:7777", "i can", "reach you")
	out, _ := cmd.Output()
	fmt.Println(string(out))

	// close
	exec.Command("go", "run", "_test/tcli.go", "127.0.0.1:9999", "*").Run()
	svr.Wait()

	// Output:
	// started server 127.0.0.1 8079 true
	// started client 127.0.0.1:8079 7777:127.0.0.1:9999 true
	// I CAN REACH YOU
}

func ExampleForwarderCmd() {
	// create a new forwarder in port 8079 that fowards to port 8888
	auth := "debug:pass"
	serverHost := "127.0.0.1"
	serverPort := 8079

	// server
	svr := exec.Command("go", "run", "_test/usrv.go", "127.0.0.1:9999")
	svr.Start()

	// create a forwarder server
	err := NewForwarderCmd(auth, serverHost, serverPort)
	if err != nil {
		return
	}
	fmt.Println("started server", serverHost, serverPort, err == nil)
	time.Sleep(1 * time.Second)

	// create a forwarder client
	server := fmt.Sprintf("http://%s:%d", serverHost, serverPort)
	rev := "7777:127.0.0.1:9999"
	err = ChiselClientCmd(server, rev, auth)
	fmt.Println("started client", server, rev, err == nil)
	time.Sleep(1 * time.Second)

	// result
	cmd := exec.Command("go", "run", "_test/tcli.go", "127.0.0.1:7777", "hello", "world")
	out, _ := cmd.Output()
	fmt.Print(string(out))

	// close
	exec.Command("go", "run", "_test/tcli.go", "127.0.0.1:9999", "*").Run()
	exec.Command("killall", "-9", "chisel").Run()
	svr.Wait()

	// Output:
	// started server 127.0.0.1 8079 true
	// started client http://127.0.0.1:8079 7777:127.0.0.1:9999 true
	// HELLO WORLD

}
