.PHONY: build-image
build-image: proxy
	cp ../common/* .
	docker image  build . --squash  -t actionloop/golang-v1.12:knative-whisk 
	docker tag actionloop/golang-v1.12:knative-whisk registry.kube-system/actionloop/golang-v1.12:knative-whisk
	docker push registry.kube-system/actionloop/golang-v1.12:knative-whisk

.PHONY: push
push: build-image
	docker tag actionloop/golang-v1.12:knative-whisk docker.io/actionloop/golang-v1.12:knative-whisk
	docker login
	docker push docker.io/actionloop/golang-v1.12:knative-whisk

.PHONY: test-compile
test-compile:
	docker run -i actionloop/golang-v1.12:knative-whisk \
	-compile main \
	 <../examples/golang-main-single/main.go \
	 >exec.zip
	 unzip -l exec.zip

.PHONY: test-build
test-build:
	-mkdir workspace
	cp ../examples/golang-main-single/main.go workspace
	docker run \
	-v $(PWD)/workspace:/workspace \
	-i actionloop/golang-v1.12:knative-whisk -build \
	main /workspace/main.go \
	registry.kube-system/actionloop/golang-v1.12:knative-whisk \
	registry.kube-system/actionloop/hello-world

catalog:
	curl -k -u registry:password https://registry.kube-system/v2/_catalog

proxy:
	env CGO_ENABLED="0" GOOS="linux" GOARCH="amd64" \
	go build -o proxy -a -ldflags "-extldflags -static" ../main/proxy.go
