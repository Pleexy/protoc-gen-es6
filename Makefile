genextensions:
	protoc --go_out=paths=source_relative:. ./proto/*.proto

genjs:
	protoc -I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I./ \
	--js_out=import_style=commonjs,binary:. \
	--flow_out=.  \
	protos/scalars/*.proto
	protoc -I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I./ \
	--js_out=import_style=commonjs,binary:. \
	--flow_out=.  \
	protos/complex/*.proto


gendebug:
	protoc -I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I./ \
	-I${GOPATH}/src/github.com/Pleexy/protoc-gen-es6/proto -I./ \
	--debug_out=".:." ./protos/complex/*.proto

genes:
	go install
	protoc -I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I./ \
	-I${GOPATH}/src/github.com/Pleexy/protoc-gen-es6/proto -I./ \
	--es6_out=".:." ./protos/complex/*.proto

gentests:
	go install
	protoc -I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I./ \
	--es6_out=".:." \
	--js_out=import_style=commonjs,binary:tests/proto/google \
	./tests/proto/*.proto


