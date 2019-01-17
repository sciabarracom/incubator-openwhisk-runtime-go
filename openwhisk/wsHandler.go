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
	"bytes"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var origin = os.Getenv("OW_WEBSOCKET_ORIGIN")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		if origin == "" {
			return true
		}
		return r.Host == origin
	},
}

func (ap *ActionProxy) wsHandler(w http.ResponseWriter, r *http.Request) {

	// enable a websocket
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		Debug("cannot upgrade websocket: %s", err.Error())
		return
	}
	defer c.Close()

	// loop
	for {
		// read a message in the websocket
		mt, body, err := c.ReadMessage()
		if err != nil {
			Debug("read error: %s", err.Error())
			break
		}
		DebugLimit("recv: %s", body, 120)

		// check if you have an action
		if ap.theExecutor == nil {
			c.WriteMessage(mt, []byte(`{"error":"no action defined yet"}`))
			return
		}

		// remove newlines and wrap in value
		body = bytes.Replace(body, []byte("\n"), []byte(""), -1)
		body = append([]byte(`{"value":`), body...)
		body = append(body, byte('}'))

		// execute the action
		response, err := ap.theExecutor.Interact(body)

		// check for early termination
		if err != nil {
			Debug("WARNING! Command exited: %s", err.Error())
			ap.theExecutor = nil
			c.WriteMessage(mt, []byte(`{"error": "command exited"}`))
			return
		}
		DebugLimit("received:", response, 120)

		// check if the underlying process exited
		if ap.theExecutor.Exited() {
			c.WriteMessage(mt, []byte(`{"error":"command exited"}`))
			return
		}

		// write the response
		err = c.WriteMessage(mt, response)
		if err != nil {
			Debug("write error: %s", err.Error())
			break
		}
	}
}
