module github.com/apache/openwhisk-runtime-go

go 1.13

replace github.com/apache/openwhisk-runtime-go/main => ./main

replace github.com/apache/openwhisk-runtime-go/openwhisk => ./openwhisk

require github.com/rs/zerolog v1.18.0
