Hello whiskers!

Sorry it is a bit long, so I split it into parts with headlines.

TL;DR

I implemented support for Websocket so you can deploy an Action as a WebSocket server if you have a Kubernetes cluster (or just a Docker server). See at the end of this mail for the example of a Kubernetes deployment descriptor.

Here is a very simple demo https://hellows.sciabarra.net/. It uses a websocket server implemented using the golang runtime and the code of an UNCHANGED OpenWhisk action. All the magic happens at deployment using the descriptor provided.

I believe it is a good foundation for implementing websocket support in OpenWhisk. The next step is to provide it at the API level.

After reading the rest, the question is: does the community approve this feature? If yes, I can submit a PR for include it in the next release of the actionloop.

1. Motivation: why I did this

A few days ago I asked what was the problem in having an action that creates a persistent connection to Kafka. I was answered that a Serverless environment can flood Kafka with requests, because more actions are spawn when the load increase.

The solution is to create a separate server to be deployed somewhere, for example Kubernetes, maybe using websockets to communicate with Kafka. In short I had the need to transform an action in a websocket server.

Hence I had the idea of adding websocket support  to  Action as WebSocket server, adding support for WebSocket to ActionLoop, so I could create a WebSocket server in the same way as you write an action. 

2. What I did

I implemented websocket support in the action runtime. If you deploy the action now it answers not only to `/run` but also to `/ws` (it is configurable) as a websocket in continuos mode.

You enable the websocket setting the environement variable OW_WEBSOCKET. Also, for the sake of easy deployment, there is also now an autoinit feature. If you set the environment variable to OW_AUTOINIT, it will initialize the runtime from the file you specified in the variable.

Ok, fine you can say, but how can I use it?

With a Kubernetes descriptor! You lanunch the runtime in Kubernetes, provide the main action in it (you can also download it from a git repo or store in a volume), and now you have a websocket server answering to your requets


Look to the following descriptor

It is a bit long, this is what it does

- it creates a configmap containing the action code
- it launches the image mounting the action code
- the image initialize the action and then listen to the websocket 
- the image is also expose the websocket using an ingress


