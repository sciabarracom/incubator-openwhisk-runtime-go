IMG?=openwhisk/actionloop-golang-v1.11:410d006
EXT=go
NAME=under-test-
START=10000
N?=100
END=$(shell expr $(START) + $(N))

report.txt: main.test main-src.test main-bin.test exec.test
	tail -n +1 *.test | tee report.txt

%.test: %.json multipost
	make start
	seq  $(START) $(END) | time xargs ./multipost -init $< -run run.dat 2>$@
	make stop

clean:
	-rm *.json *.zip *.b64 *.test exec report.txt multipost

main-src.zip: main.$(EXT)
	zip main-src.zip main.$(EXT)

main-bin.zip: main-src.zip
	docker run -i $(IMG) -compile main <main-src.zip >main-bin.zip

multipost: multipost.go
	go build -o multipost multipost.go

exec: main-bin.zip
	unzip -o main-bin.zip

exec.json: exec
	base64 exec >exec.b64
	jq -n --rawfile f exec.b64  '{"value":{"binary":true,"code": $$f}}' >exec.json

%.json: %.zip
	base64 $< >$*.b64
	jq -n --rawfile f $*.b64  '{"value":{"binary":true,"code": $$f}}' >$*.json

%.json: %.$(EXT)
	jq -n --rawfile f $<  '{"value":{"code": $$f}}' >$*.json

start:
	seq $(START) $(END) | while read i; do \
	docker run -d -p $$i:8080 --name $(NAME)$$i --rm $(IMG)  ;\
	done

stop:
	seq $(START) $(END) | while read i; do \
	docker kill $(NAME)$$i ;\
	done

.PHONY: start stop all

