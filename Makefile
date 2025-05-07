PACKAGE_ROOT=api
OUT_DIR=internal/pb
GATEWAY_INCLUDE=$(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway/v2)

PROTOC_GEN_GO=$(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC=$(shell which protoc-gen-go-grpc)
PROTOC_GEN_GRPC_GATEWAY=$(shell which protoc-gen-grpc-gateway)
PROTOC_GEN_OPENAPIV2=$(shell which protoc-gen-openapiv2)

.PHONY: all proto check-tools install-tools

all: proto

check-tools:
ifndef PROTOC_GEN_GO
	$(error protoc-gen-go is not installed)
endif
ifndef PROTOC_GEN_GO_GRPC
	$(error protoc-gen-go-grpc is not installed)
endif
ifndef PROTOC_GEN_GRPC_GATEWAY
	$(error protoc-gen-grpc-gateway is not installed)
endif
ifndef PROTOC_GEN_OPENAPIV2
	$(error protoc-gen-openapiv2 is not installed)
endif

proto: check-tools
	find $(PACKAGE_ROOT) -name "*.proto" | while read -r file; do \
		echo "Generating for $$file"; \
		REL_PATH=$$(dirname $$file); \
		OUT_PATH=$(OUT_DIR)/$$REL_PATH; \
		mkdir -p "$$OUT_PATH"; \
		protoc \
			-I . \
			-I ./vendor.protogen \
			-I $(GATEWAY_INCLUDE) \
			--go_out=$(OUT_DIR) \
			--go_opt=paths=source_relative \
			--go-grpc_out=$(OUT_DIR) \
			--go-grpc_opt=paths=source_relative \
			--grpc-gateway_out=$(OUT_DIR) \
			--grpc-gateway_opt=paths=source_relative \
			--openapiv2_out="$$OUT_PATH" \
			--openapiv2_opt="logtostderr=true" \
			$$file; \
	done

install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest


# ---------- таргет: один сервис ----------
# Пример: make generate-client SERVICE=fields-api
generate-client: check-tools
ifndef SERVICE
	$(error Specify SERVICE=<dir inside vendor.protogen>)
endif
	@echo "▶︎  Generating client for service: $(SERVICE)"

	# все .proto этого сервиса
	$(eval PROTOS := $(shell find vendor.protogen/$(SERVICE) -name "*.proto"))

	# куда кладём
	$(eval OUT_ROOT := internal/pb/$(SERVICE))
	@mkdir -p $(OUT_ROOT)

	# единый protoc
	protoc \
		-I vendor.protogen \
		-I vendor.protogen/$(SERVICE) \
		--go_out=$(OUT_ROOT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_ROOT) \
		--go-grpc_opt=paths=source_relative \
		$(PROTOS)

	@echo "✅  OK — client in $(OUT_ROOT)"


# ---------- таргет: все сервисы ----------
# make generate-all   (обойдёт все подпапки vendor.protogen/*)
generate-all: check-tools
	@for svc in $$(ls vendor.protogen); do \
		$(MAKE) generate-client SERVICE=$$svc; \
	done