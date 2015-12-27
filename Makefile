
install-protoc-gen-goevpb:
	go install github.com/crackcomm/evpb/protoc-gen-goevpb

protoc: install-protoc-gen-goevpb
	protoc -I${GOPATH}/src --goevpb_out=plugins=evpb:${GOPATH}/src \
		${GOPATH}/src/github.com/crackcomm/evpb/example/pb/message.proto
