package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func invokeParse(cmd string) bool {
	switch cmd {
	case debugRunCmd.FullCommand():
		ch := make(chan string)
		go doRun(*debugRunIPArg, *debugRunJSONArg, ch)
		fmt.Println(<-ch)
		return true
	case debugInitCmd.FullCommand():
		doInit(*debugInitIPArg, *debugInitFileArg, *debugInitMainArg, *debugInitEnvArg)
		return true
	case debugStartCmd.FullCommand():
		doStart(*debugStartIPArg, *debugStartJSONArg)
		return true
	case debugMyIPCmd.FullCommand():
		fmt.Println(myIPAddress())
		return true
	case debugFwdCmd.FullCommand():
		fmt.Printf("*** forwarding to %s, please connect to the debugger now ***\n", *debugFwdIPCmd)
		conn, err := acceptConnection("127.0.0.1:8081")
		if err == nil {
			portForwarding(conn, *debugFwdIPCmd+":8081")
			fmt.Println("Connection established, press enter to exit:")
			fmt.Scanln() // wait for Enter Key
		}
		return true
	}
	return false
}

func acceptConnection(hostport string) (net.Conn, error) {
	log.Info("Listening to", hostport)
	sock, err := net.Listen("tcp", hostport)
	if err != nil {
		return nil, err
	}
	log.Info("Connected to ", hostport)
	return sock.Accept()
}

func receiveDebugInfo(ch chan string) {
	// accept connection from action for debug info
	log.Info("receiveDebugInfo")
	conn, err := acceptConnection(":7737")
	if err != nil {
		log.Info(err)
		ch <- ""
		return
	}
	// read debuginfo
	log.Info("reading from socket")
	scan := bufio.NewScanner(conn)
	if scan.Scan() {
		msg := scan.Text()
		log.Info("received from socket ", msg)
		ch <- msg
	}
	// wait for an ack
	<-ch
	log.Info("got ack")
	// communicate to the destination
	conn.Write([]byte("\n"))
	conn.Close()
	log.Info("sent ack and closed")
}

func doStart(ip string, json string) {
	ch := make(chan string)
	// receive debug info
	go receiveDebugInfo(ch)
	go doRun(ip, json, ch)
	// invoke the action
	// receive the debug info
	msg := <-ch
	log.Info("got message", msg)
	if msg == "" {
		log.Info("cannot connect to debugger")
		return
	}
	a := strings.Split(msg, ":")
	fmt.Println("*** please connect to the debugger ***")
	conn, err := acceptConnection("127.0.0.1:8081")
	if err != nil {
		log.Info(err)
		ch <- "\n"
		return
	}
	// debugger ready, go on
	portForwarding(conn, a[1]+":8081")
	// acknowlede to run action
	ch <- "\n"
	// wait for run result
	res := <-ch
	fmt.Println(res)
}

func portForwarding(local net.Conn, target string) {
	remote, err := net.Dial("tcp", target)
	if err != nil {
		log.Info(err)
		return
	}
	copyIO := func(src, dst net.Conn) {
		defer src.Close()
		defer dst.Close()
		io.Copy(src, dst)
	}
	log.Info("starting forwarding to ", target)
	go copyIO(local, remote)
	go copyIO(remote, local)
}

func doInit(ip string, file string, main string, env string) {
	url := fmt.Sprintf("http://%s:8080/init", ip)
	data, err := encodeInit(file, main, env)
	FatalIf(err)
	doPost(url, data)
}

func doRun(ip string, json string, ch chan string) {
	myip := myIPAddress()
	url := fmt.Sprintf("http://%s:8080/run", ip)
	encoded := fmt.Sprintf(`{ "debugger": "%s", "value": %s }`, myip, json)
	ch <- doPost(url, []byte(encoded))
}

func doPost(url string, data []byte) string {

	log.Printf(">>> %s\n=== %s\n", url, data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	transport := http.Transport{}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("!!! Connection Error: %s -> %s\n", url, err)
		return ""
	}
	if resp.StatusCode != 200 {
		log.Printf("!!! Http Error: %s -> %s\n", url, resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	resp.Body.Close()
	transport.CloseIdleConnections()
	return string(body)
}

func doRunBegin(ip string, json string) {
	myip := myIPAddress()
	url := fmt.Sprintf("http://%s:8080/run", ip)
	encoded := fmt.Sprintf(`{ "debugger": "%s", "value": %s }`, myip, json)
	doPost(url, []byte(encoded))
}

func encodeInit(filename string, main string, envJSON string) ([]byte, error) {
	buf, _ := ioutil.ReadFile(filename)
	toEncode := make(map[string]interface{})
	toEncode["main"] = main
	var env map[string]interface{}
	json.Unmarshal([]byte(envJSON), &env)
	toEncode["env"] = env
	log.Printf("%s\n", env)
	if strings.HasSuffix(filename, ".zip") || strings.HasSuffix(filename, ".exe") || strings.HasSuffix(filename, ".jar") {
		toEncode["binary"] = true
		toEncode["code"] = base64.StdEncoding.EncodeToString(buf)
	} else {
		toEncode["code"] = string(buf)
	}
	res, err := json.Marshal(map[string]interface{}{"value": toEncode})
	return res, err
}
