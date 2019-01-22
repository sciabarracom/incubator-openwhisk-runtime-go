<!--
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
-->

The runtime implements support for Websocket so you can deploy an Action as a WebSocket server if you have a Kubernetes cluster (or just a Docker server).

Here is a very simple demo https://hellows.sciabarra.net/. It uses a websocket server implemented using the golang runtime and the code of an UNCHANGED OpenWhisk action. All the magic happens at deployment using the descriptor provided.

If you deploy the action now it answers not only to `/run` but also to `/ws` (it is configurable) as a websocket in continuos mode.

You enable the websocket setting the environement variable OW_WEBSOCKET. Also, for the sake of easy deployment, there is also now an autoinit feature. If you set the environment variable to OW_AUTOINIT, it will initialize the runtime from the file you specified in the variable.

You use this feature with a Kubernetes descriptor! Launching the runtime in Kubernetes, provide the main action in it (you can also download it from a git repo or store in a volume), and now you have a websocket server answering to your requets

An example is provided in [../examples/websocket-hello/hellows.yml](this  descriptor):

- it creates a configmap containing the action code
- it launches the image mounting the action code
- the image initialize the action and then listen to the websocket
- the image is also expose the websocket using an ingress


