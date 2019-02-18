GO_PATH := $(shell go env GOPATH)
GO_BIN := $(GO_PATH)/bin
GO_PKGS := $(shell go list ./...)
PROTOC_GEN_GO := $(GO_BIN)/protoc-gen-go
PROTOC_GEN_GO_SRC := vendor/github.com/golang/protobuf/protoc-gen-go

PROTO_FILES := $(shell ls proto/*.proto)
GO_PROTO := $(patsubst proto%,service%,$(PROTO_FILES:proto=pb.go))

.PHONY: gotest
gotest:
	go test $(GO_PKGS)

.PHONY: go_proto
go_proto: $(GO_PROTO)
	echo $(GO_PROTO)

service/%.pb.go: $(PROTOC_GEN_GO) proto/%.proto
	protoc -I proto --go_out=plugins=grpc:service proto/*.proto
