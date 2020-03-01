#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
IMG=whisk/action-golang-v1.13:latest

build:
	GOOS=linux GOARCH=amd64 go build -o proxy -a  -ldflags '-extldflags "-static"' ../main/proxy.go
	docker build -t $(IMG) .

clean:
	docker rmi -f $(IMG)

debug:
	docker run -p 8080:8080 -p 8079:8079 \
	--name go-action --rm -ti --entrypoint=/bin/bash \
	-e __OW_DEBUG_PORT=8081 \
	-v $(PWD):/mnt -e OW_COMPILER=/mnt/bin/compile \
	$(IMG)

server:
	docker exec -ti go-action chisel server -p 8079

client:
	chisel client http://localhost:8079 8081:127.0.0.1:8081


.PHONY: build clean  debug
