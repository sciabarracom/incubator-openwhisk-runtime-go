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

export GOPATH=$PWD

function build {
   cp $1.src $1.go
   go build -a -o $1 $1.go
   rm $1.go
}

function build_main {
   cp ../../common/gobuild.py.launcher.go $1.go
   cat $1.src >>$1.go
   go build -a -o $1 $1.go
   rm $1.go
}

build hi
cp -f hi exec 
zip hi.zip exec

build_main hello_message
cp -f hello_message exec 
zip hello_message.zip exec

build_main hello_greeting
cp -f hello_greeting exec
zip hello_greeting.zip exec

test -e hello.zip && rm hello.zip
cd src
zip -q -r ../hello.zip main.go hello
cd ..

build empty
cp -f empty exec
test -e exec.zip && rm exec.zip
zip -q -r exec.zip exec etc dir
