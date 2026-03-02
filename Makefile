PROTO_DIR      := proto
GEN_GO_DIR     := server/gen
GEN_PY_DIR     := sdk/python/gen
PROTOC_INCLUDE := $(HOME)/.local/include

.PHONY: proto-go proto-python proto-all clean-gen

proto-go:
	@mkdir -p $(GEN_GO_DIR)
	protoc \
		--proto_path=$(PROTO_DIR) \
		--proto_path=$(PROTOC_INCLUDE) \
		--go_out=$(GEN_GO_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_GO_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/uos/kernel/v1/kernel.proto

proto-python:
	@mkdir -p $(GEN_PY_DIR)
	python3 -m grpc_tools.protoc \
		--proto_path=$(PROTO_DIR) \
		--proto_path=$(PROTOC_INCLUDE) \
		--python_out=$(GEN_PY_DIR) \
		--grpc_python_out=$(GEN_PY_DIR) \
		--pyi_out=$(GEN_PY_DIR) \
		$(PROTO_DIR)/uos/kernel/v1/kernel.proto

proto-all: proto-go proto-python

clean-gen:
	rm -rf $(GEN_GO_DIR) $(GEN_PY_DIR)
