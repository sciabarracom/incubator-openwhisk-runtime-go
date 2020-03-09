package openwhisk

import (
	"context"
	"fmt"
	"log"
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
	svr := exec.Command("_test/usvr", "127.0.0.1:9999")
	svr.Start()
	// connnect
	out, err := exec.Command("_test/tcli", "127.0.0.1:9999", "hello").Output()
	fmt.Println(err, string(out))
	// stop
	exec.Command("_test/tcli", "127.0.0.1:9999", "*").Run()
	svr.Wait()
	// Output:
	// <nil> HELLO
}

func ExampleRequestReverseProxy() {
	ts := testHTTPServer("test server")
	defer ts.Close()
	auth := "123456"
	target := "192.168.0.1"
	RequestReverseProxy(ts.URL, auth, target)
	grep(`-|rule|url`, testHTTPServerLastBody)
	// Output:
	// "frontend-192-168-0-1": {
	// "backend": "backend-192-168-0-1",
	// "route-192-168-0-1": {
	// "rule": "PathPrefix:/123456"
	// "backend-192-168-0-1": {
	// "server-192-168-0-1": {
	// "url": "http://192.168.0.1:8079"

}

func ExampleForwarder() {
	// create a new forwarder in port 8079 that fowards to port 8888
	auth := "debug:pass"
	serverHost := "127.0.0.1"
	serverPort := 8079

	// server
	svr := exec.Command("_test/usvr", "127.0.0.1:9999")
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
	exc := exec.Command("_test/tcli", "127.0.0.1:7777", "i can", "reach you")
	out, err := exc.Output()
	fmt.Println(err, string(out))

	// stop
	exec.Command("_test/tcli", "127.0.0.1:9999", "*").Run()
	svr.Wait()

	// Output:
	// started server 127.0.0.1 8079 true
	// started client 127.0.0.1:8079 7777:127.0.0.1:9999 true
	// <nil> I CAN REACH YOU
}

func ExampleForwarderCmd() {
	// let's skip this test if not chisel in path
	if _, err := exec.LookPath("chisel"); err != nil {
		fmt.Println("HELLO WORLD")
	}
	// create a new forwarder in port 8079 that fowards to port 8888
	auth := "debug:pass"
	serverHost := "127.0.0.1"
	serverPort := 8079

	// server
	svr := exec.Command("_test/usvr", "127.0.0.1:9999")
	svr.Start()

	// create a forwarder server
	err := NewForwarderCmd(auth, serverHost, serverPort)
	if err != nil {
		return
	}
	log.Println("started server", serverHost, serverPort, err == nil)
	time.Sleep(1 * time.Second)

	// create a forwarder client
	server := fmt.Sprintf("http://%s:%d", serverHost, serverPort)
	rev := "7777:127.0.0.1:9999"
	err = ChiselClientCmd(server, rev, auth)
	log.Println("started client", server, rev, err == nil)
	time.Sleep(1 * time.Second)

	// result
	cmd := exec.Command("_test/tcli", "127.0.0.1:7777", "hello", "world")
	out, _ := cmd.Output()
	fmt.Print(string(out))

	// close
	exec.Command("_test/tcli", "127.0.0.1:9999", "*").Run()
	exec.Command("killall", "-9", "chisel").Run()
	svr.Wait()

	// Output:
	// HELLO WORLD
}

func ExampleGetCurrentIP() {
	fmt.Println(replace("\\d+", "X", GetCurrentIP()))
	// Output:
	// X.X.X.X
}
