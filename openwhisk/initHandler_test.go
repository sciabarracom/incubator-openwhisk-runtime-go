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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

/**
The _test/build.sh script builds some binaries that are used to actually run tests.
Tests basically submit various binaries to the handler, simulating the init of the runtime.
**/

func Example_json_init() {
	fmt.Println(initCode("", ""))
	fmt.Println(initCode("_test/etc", ""))
	fmt.Println(initCode("_test/etc", "world"))
	fmt.Println(initBinary("_test/etc", ""))
	fmt.Println(initBinary("_test/etc", "hello"))
	// Output:
	// {"value":{}}
	// {"value":{"code":"1\n"}}
	// {"value":{"code":"1\n","main":"world"}}
	// {"value":{"code":"MQo=","binary":true}}
	// {"value":{"code":"MQo=","binary":true,"main":"hello"}}
}

func Example_bininit_nocompiler() {
	ts, cur, log := startTestServer("")
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello_message", ""))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	ts, cur, log = startTestServer("")
	doInit(ts, initBinary("_test/hello_greeting", ""))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 500 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"message":"Hello, Mike!"}
	// name=Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// 200 {"ok":true}
	// 200 {"greetings":"Hello, Mike"}
	// Hello, Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_zipinit_nocompiler() {
	ts, cur, log := startTestServer("")
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello_greeting.zip", ""))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	ts, cur, log = startTestServer("")
	doInit(ts, initBinary("_test/hello_message.zip", ""))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 500 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"greetings":"Hello, Mike"}
	// Hello, Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// 200 {"ok":true}
	// 200 {"message":"Hello, Mike!"}
	// name=Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_shell_nocompiler() {
	ts, cur, log := startTestServer("")
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello.sh", ""))
	doRun(ts, "")
	doRun(ts, `{"name":"world"}`)
	stopTestServer(ts, cur, log)
	// Output:
	// 500 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"hello": "Mike"}
	// 200 {"hello": "world"}
	// msg=hello Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// msg=hello world
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_main_nocompiler() {
	ts, cur, log := startTestServer("")
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello_message", "message"))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	ts, cur, log = startTestServer("")
	doInit(ts, initBinary("_test/hello_greeting", "greeting"))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 500 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"message":"Hello, Mike!"}
	// name=Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// 200 {"ok":true}
	// 200 {"greetings":"Hello, Mike"}
	// Hello, Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_main_zipinit_nocompiler() {
	ts, cur, log := startTestServer("")
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello_greeting.zip", "greeting"))
	doRun(ts, "")
	stopTestServer(ts, cur, log)

	ts, cur, log = startTestServer("")
	doInit(ts, initBinary("_test/hello_message.zip", "message"))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 500 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"greetings":"Hello, Mike"}
	// Hello, Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// 200 {"ok":true}
	// 200 {"message":"Hello, Mike!"}
	// name=Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_compile_simple() {
	comp, _ := filepath.Abs("../common/gobuild.py")
	ts, cur, log := startTestServer(comp)
	doRun(ts, "")
	doInit(ts, initCode("_test/hello.src", ""))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 500 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"message":"Hello, Mike!"}
	// name=Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_compile_withMain() {
	comp, _ := filepath.Abs("../common/gobuild.py")
	ts, cur, log := startTestServer(comp)
	doRun(ts, "")
	doInit(ts, initCode("_test/hello1.src", "hello"))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 500 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"hello":"Hello, Mike!"}
	// name=Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_compile_withZipSrc() {
	comp, _ := filepath.Abs("../common/gobuild.py")
	ts, cur, log := startTestServer(comp)
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello.zip", ""))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 500 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"greetings":"Hello, Mike"}
	// Main
	// Hello, Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_badinit_nocompiler() {
	ts, cur, log := startTestServer("")
	doRun(ts, "")
	doInit(ts, "{}")
	//sys("ls", "_test/exec")
	doInit(ts, initBinary("_test/exec", ""))      // empty
	doInit(ts, initBinary("_test/hi", ""))        // say hi
	doInit(ts, initBinary("_test/hello.src", "")) // source not excutable
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 500 {"error":"no action defined yet"}
	// 403 {"error":"Missing main/no code to execute."}
	// 502 {"error":"cannot start action: command exited"}
	// 502 {"error":"cannot start action: command exited"}
	// 502 {"error":"cannot start action: command exited"}
	// 500 {"error":"no action defined yet"}
	// hi
}

func Example_zip_init() {
	ts, cur, log := startTestServer("")
	buf, _ := Zip("_test/pysample")
	doInit(ts, initBytes(buf, ""))
	doRun(ts, ``)
	doRun(ts, `{"name":"World"}`)
	stopTestServer(ts, cur, log)
	// Output:
	// 200 {"ok":true}
	// 200 {"python": "Hello, Mike"}
	// 200 {"python": "Hello, World"}
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_parse_env() {
	var request initBodyRequest
	body := []byte(`{"code":"code"}`)
	json.Unmarshal(body, &request)
	fmt.Println(request.Env)

	var request1 initBodyRequest
	body = []byte(`{"code":"code", "env":{"hello":"world"}}`)
	json.Unmarshal(body, &request1)
	fmt.Println(request1.Env["hello"])

	var request2 initBodyRequest
	body = []byte(`{"code":"code", "env": { "hello": "world", "hi": "all"}}`)
	json.Unmarshal(body, &request2)
	fmt.Println(request2.Env["hello"], request2.Env["hi"])
	// Output:
	// map[]
	// world
	// world all
}

func Example_debugger() {

	// proxy requests receiver
	dts := testHTTPServer("test server")
	os.Setenv("__OW_DEBUG_PROXY", dts.URL)

	// test server
	ts, cur, log := startTestServer("")

	// intialized expecting the debugger
	env := map[string]interface{}{
		"__OW_DEBUG_PORT": "8081",
		"__OW_DEBUG_AUTH": "123456",
	}
	doInit(ts, initBinaryEnv("_test/hello_debugger.zip", "main", env))

	// check the debugger
	out, err := exec.Command("_test/tcli", "127.0.0.1:8081",
		"hello").Output()
	fmt.Print(1, err, string(out))

	// check the forwarder
	server := "127.0.0.1:8079"
	rev := "8079:127.0.0.1:8081"
	cli, err := ChiselClient(server, rev, "123456")
	cli.Logger.Info = false
	cli.Logger.Debug = false
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println(2, "started client", server, rev, cli.Start(ctx) == nil)

	// check the debugger
	out, err = exec.Command("_test/tcli", "127.0.0.1:8081",
		"forward me").Output()
	fmt.Print(3, err, string(out))

	// check the proxy server rules
	fmt.Println(4, "checking rules and run")
	grep(`rule|url`, replace("\\d+", "XXX", testHTTPServerLastBody))

	// try a run and finish
	doRun(ts, `{"name":"World"}`)
	stopTestServer(ts, cur, log)
	dts.Close()
	cancel()
	// Output:
	// 200 {"debug":true}
	// 1 <nil>HELLO
	// 2 started client 127.0.0.1:8079 8079:127.0.0.1:8081 true
	// 3 <nil>FORWARD ME
	// 4 checking rules and run
	// "rule": "PathPrefix:/XXX"
	// "url": "http://XXX.XXX.XXX.XXX:XXX"
	// 200 {"message":"Hello, World!"}
	// name=World
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}
