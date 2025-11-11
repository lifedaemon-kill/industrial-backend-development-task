API_DIR := api
THIRD_PARTY_DIR := .vendor-proto
GOOGLE_API_DIR := $(THIRD_PARTY_DIR)/google/api
VALIDATE_DIR := $(THIRD_PARTY_DIR)/validate
BIN_DIR := bin
GEN_DIR := pkg/protogen
PROTO_FILES := $(wildcard $(API_DIR)/**/*.proto) $(wildcard $(API_DIR)/*.proto)

PROTOC_GEN_GO := $(BIN_DIR)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(BIN_DIR)/protoc-gen-go-grpc
PROTOC_GEN_GRPC_GATEWAY := $(BIN_DIR)/protoc-gen-grpc-gateway
PROTOC_GEN_OPENAPIV2 := $(BIN_DIR)/protoc-gen-openapiv2


.PHONY: all deps generate clean run up build

all: run

run: deps generate build

build:
	docker-compose up --build
up:
	docker-compose up

deps:
	@echo "Installing protoc plugins into $(BIN_DIR)..."
	mkdir -p $(BIN_DIR)

	GOBIN=$(abspath $(BIN_DIR)) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(abspath $(BIN_DIR)) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(abspath $(BIN_DIR)) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(abspath $(BIN_DIR)) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

	@echo "Downloading Google API proto files..."
	mkdir -p $(GOOGLE_API_DIR)
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto -o $(GOOGLE_API_DIR)/annotations.proto
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto -o $(GOOGLE_API_DIR)/http.proto

	@echo "Downloading validate.proto..."
	mkdir -p $(VALIDATE_DIR)
	curl -sSL https://raw.githubusercontent.com/bufbuild/protoc-gen-validate/main/validate/validate.proto -o $(VALIDATE_DIR)/validate.proto

generate:
	@echo "Generating code..."
	PATH=$(abspath $(BIN_DIR)):$$PATH protoc -I $(API_DIR) \
		-I $(THIRD_PARTY_DIR) \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(GEN_DIR) --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=./api/docs --openapiv2_opt=logtostderr=true \
		$(PROTO_FILES)

clean:
	@echo "Cleaning up generated files..."
	find $(API_DIR) -name "*.pb.go" -delete
	find $(API_DIR) -name "*.gw.go" -delete
	rm -rf $(GOOGLE_API_DIR)
	rm -rf $(VALIDATE_DIR)
	rm -rf $(API_DIR)/docs
	rm -rf $(BIN_DIR)
