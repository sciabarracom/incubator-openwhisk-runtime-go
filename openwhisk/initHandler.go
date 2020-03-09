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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type initBodyRequest struct {
	Code   string                 `json:"code,omitempty"`
	Binary bool                   `json:"binary,omitempty"`
	Main   string                 `json:"main,omitempty"`
	Env    map[string]interface{} `json:"env,omitempty"`
}

type initRequest struct {
	Value initBodyRequest `json:"value,omitempty"`
}

func sendOK(w http.ResponseWriter, debugger bool) {
	// answer OK
	w.Header().Set("Content-Type", "application/json")
	key := "ok"
	if debugger {
		key = "debug"
	}
	buf := []byte(fmt.Sprintf("{\"%s\":true}\n", key))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(buf)))
	w.Write(buf)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func (ap *ActionProxy) initHandler(w http.ResponseWriter, r *http.Request) {

	// you can do muliple initializations when debugging
	if ap.initialized && !Debugging {
		msg := "Cannot initialize the action more than once."
		sendError(w, http.StatusForbidden, msg)
		log.Println(msg)
		return
	}

	// read body of the request
	if ap.compiler != "" {
		Debug("compiler: " + ap.compiler)
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	// decode request parameters
	if len(body) < 1000 {
		Debug("init: decoding %s\n", string(body))
	}

	var request initRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Error unmarshaling request: %v", err))
		return
	}

	// request with empty code - stop any executor but return ok
	if request.Value.Code == "" {
		sendError(w, http.StatusForbidden, "Missing main/no code to execute.")
		return
	}

	// passing the env to the action proxy
	ap.SetEnv(request.Value.Env)

	// setting main
	main := request.Value.Main
	if main == "" {
		main = "main"
	}

	// extract code eventually decoding it
	var buf []byte
	if request.Value.Binary {
		Debug("it is binary code")
		buf, err = base64.StdEncoding.DecodeString(request.Value.Code)
		if err != nil {
			sendError(w, http.StatusBadRequest, "cannot decode the request: "+err.Error())
			return
		}
	} else {
		Debug("it is source code")
		buf = []byte(request.Value.Code)
	}

	// if a compiler is defined try to compile
	_, err = ap.ExtractAndCompile(&buf, main)
	if err != nil {
		if os.Getenv("OW_LOG_INIT_ERROR") == "" {
			sendError(w, http.StatusBadGateway, err.Error())
		} else {
			ap.errFile.Write([]byte(err.Error() + "\n"))
			ap.outFile.Write([]byte(OutputGuard))
			ap.errFile.Write([]byte(OutputGuard))
			sendError(w, http.StatusBadGateway, "The action failed to generate or locate a binary. See logs for details.")
		}
		return
	}

	// start an action
	err = ap.StartLatestAction()
	if err != nil {
		if os.Getenv("OW_LOG_INIT_ERROR") == "" {
			sendError(w, http.StatusBadGateway, "cannot start action: "+err.Error())
		} else {
			ap.errFile.Write([]byte(err.Error() + "\n"))
			ap.outFile.Write([]byte(OutputGuard))
			ap.errFile.Write([]byte(OutputGuard))
			sendError(w, http.StatusBadGateway, "Cannot start action. Check logs for details.")
		}
		return
	}
	ap.initialized = true
	debugger := StartDebuggerIfRequired(request.Value.Env)
	sendOK(w, debugger)
}

func getValueFromEnvOrMap(key string, env map[string]interface{}) string {
	value, ok := env[key].(string)
	if ok {
		return value
	}
	return os.Getenv(key)
}

// StartDebuggerIfRequired checks environment variables
// for __OW_DEBUG_PORT __OW_DEBUG_PROXY __OW_DEBUG_AUTH
// either in the map and in the environment and starts debugger if required
func StartDebuggerIfRequired(env map[string]interface{}) bool {
	debugPort := getValueFromEnvOrMap("__OW_DEBUG_PORT", env)
	debugProxy := getValueFromEnvOrMap("__OW_DEBUG_PROXY", env)
	debugAuth := getValueFromEnvOrMap("__OW_DEBUG_AUTH", env)
	if debugPort == "" || debugProxy == "" || debugAuth == "" {
		Debug("Some paramter is empty: port:%s proxy:%s auth:%s", debugPort, debugProxy, debugAuth)
		return false
	}

	Debug("ConnectDebugger: port:%s proxy:%s auth:%s", debugPort, debugProxy, debugAuth)
	fwd, err := ConnectDebugger(debugPort, debugProxy, debugAuth)
	if err != nil {
		Debug("connect error %s", err)
		return false
	}
	debugTarget := GetCurrentIP()
	Debug("RequestReverseProxy: %s %s %s", debugProxy, debugAuth, debugTarget)
	res, err := RequestReverseProxy(debugProxy, debugAuth, debugTarget)
	if err != nil {
		Debug("reverse proxy error %s", err)
		fwd.Close()
		return false
	}
	Debug("ok: %s", res)
	return true
}

// ExtractAndCompile decode the buffer and if a compiler is defined, compile it also
func (ap *ActionProxy) ExtractAndCompile(buf *[]byte, main string) (string, error) {

	// extract action in src folder
	file, err := ap.ExtractAction(buf, "src")
	if err != nil {
		return "", err
	}
	if file == "" {
		return "", fmt.Errorf("empty filename")
	}

	// some path surgery
	dir := filepath.Dir(file)
	parent := filepath.Dir(dir)
	srcDir := filepath.Join(parent, "src")
	binDir := filepath.Join(parent, "bin")
	binFile := filepath.Join(binDir, "exec")

	// if the file is already compiled or there is no compiler just move it from src to bin
	if ap.compiler == "" || isCompiled(file) {
		os.Rename(srcDir, binDir)
		return binFile, nil
	}

	// ok let's try to compile
	Debug("compiling: %s main: %s", file, main)
	os.Mkdir(binDir, 0755)
	err = ap.CompileAction(main, srcDir, binDir)
	if err != nil {
		return "", err
	}

	// check only if the file exist
	if _, err := os.Stat(binFile); os.IsNotExist(err) {
		return "", fmt.Errorf("cannot compile")
	}
	return binFile, nil
}
