#!/bin/bash
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

cd "$(dirname $0)"

# trick to avoid rebuilding all the time with vscode when running tests
if [[ -n "$VSCODE_PID" ]] && [[ -e "/tmp/openwisk-runtime-go.$VSCODE_PID" ]]
then exit 0
else touch /tmp/openwisk-runtime-go.$VSCODE_PID
fi

function build {
   test -e exec && rm exec
   cp $1.src $1.go
   GOPATH=$PWD go build -a -o exec $1.go
   rm $1.go
}

function build_main_with_net {
   test -e exec && rm exec
   cat ../../common/gobuild.py.launcher.go \
   | awk '/"bufio"/ {  print "\"net\"" }1' >$1.go
   cat $1.src >>$1.go
   go build -a -o exec $1.go
   rm $1.go
}

function build_main {
   test -e exec && rm exec
   cat ../../common/gobuild.py.launcher.go >$1.go
   cat $1.src >>$1.go
   go build -a -o exec $1.go
   rm $1.go
}


# build zip
rm action.zip 2>/dev/null
zip -r -q action.zip action

# test tcp
build usvr
mv exec usvr
build tcli
mv exec tcli

# test actions
build hi
zip -q hi.zip exec
cp exec hi

build_main hello_message
zip -q hello_message.zip exec
cp exec hello_message

build_main hello_greeting
zip -q hello_greeting.zip exec
cp exec hello_greeting

build_main_with_net hello_debugger
zip -q hello_debugger.zip exec
cp exec hello_debugger

test -e hello.zip && rm hello.zip
cd src
zip -q -r ../hello.zip main.go hello
cd ..

test -e sample.jar && rm sample.jar
cd jar ; zip -q -r ../sample.jar * ; cd ..

build exec
test -e exec.zip && rm exec.zip
zip -q -r exec.zip exec etc dir

echo exec/env >helloack/exec.env
zip -j helloack.zip helloack/*
