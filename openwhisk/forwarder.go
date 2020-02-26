/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package openwhisk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	chclient "github.com/jpillora/chisel/client"
	chserver "github.com/jpillora/chisel/server"
)

// Forwarder is a forwarder
type Forwarder struct {
	*chserver.Server
	host string
	port string
}

// NewForwarder intializes ans start the forwarder server
// it expects the authentication key, and the port it should listen
func NewForwarder(auth string, host string, port int) (*Forwarder, error) {
	server, err := chserver.NewServer(&chserver.Config{
		Auth:     auth,
		KeySeed:  "",
		AuthFile: "",
		Proxy:    "",
		Socks5:   false,
		Reverse:  false,
	})
	if err != nil {
		return nil, err
	}
	fwd := Forwarder{
		Server: server,
		host:   host,
		port:   strconv.Itoa(port),
	}
	Debug("started chisel server %s %d", host, port)
	return &fwd, nil
}

// NewForwarderCmd starts the chisel server as a command
func NewForwarderCmd(auth string, host string, port int) error {
	cmd := exec.Command("chisel", "server",
		"--auth", auth,
		"--host", host,
		"--port", strconv.Itoa(port))
	return cmd.Start()
}

// Start forwarder
func (fwd *Forwarder) Start() error {
	return fwd.Server.Start(fwd.host, fwd.port)
}

// RequestReverseProxy requests a reverse proxy
func RequestReverseProxy(proxy string, prefix string, target string) ([]byte, error) {
	data := strings.ReplaceAll(fmt.Sprintf(`{
  "frontends": {
    "frontend$N$": {
      "backend": "backend$N$",
      "routes": {
        "route$N$": {
          "rule": "PathPrefix:%s"
        }
      }
    }
  },
  "backends": {
    "backend$N$": {
      "servers": {
        "server$N$": {
          "url": "%s"
        }
      }
    }
  }
}`, prefix, target), "$N$", "1")
	client := &http.Client{}
	restReq := proxy + "/api/providers/rest"
	req, err := http.NewRequest(http.MethodPut, restReq, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

// ChiselClient starts a chisel client
func ChiselClient(server string, remote string, auth string) (*chclient.Client, error) {
	return chclient.NewClient(&chclient.Config{
		Server:           server,
		Remotes:          []string{remote},
		Auth:             auth,
		Fingerprint:      "",
		KeepAlive:        0,
		MaxRetryCount:    -1,
		MaxRetryInterval: 0,
		HTTPProxy:        "",
		HostHeader:       "",
	})
}

// ChiselClientCmd is it
func ChiselClientCmd(server string, remote string, auth string) error {
	cmd := exec.Command("chisel", "client",
		"--auth", auth, server, remote)
	return cmd.Start()
}
