
install-protoc-gen-go:
	go install github.com/crackcomm/evpb/protoc-gen-go

protoc:
	protoc -I${GOPATH}/src --go_out=plugins=evpb:${GOPATH}/src \
		${GOPATH}/src/github.com/crackcomm/evpb/example/pb/message.proto
