CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
PACKAGE=route256/test

build: bindir
	go test -c -o ${BINDIR}/app

bindir:
	mkdir -p ${BINDIR}

# Используем bin в текущей директории для установки плагинов protoc
bin: bindir
	GOBIN=$(BINDIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(BINDIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(BINDIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(BINDIR) go install github.com/envoyproxy/protoc-gen-validate@latest

# Добавляем bin в текущей директории в PATH при запуске protoc
PROTOC = PATH="$$PATH:$(BINDIR)" protoc

# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis &&\
 	cd vendor-proto/googleapis &&\
	git sparse-checkout set --no-cone google/api &&\
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

# Устанавливаем proto описания validate
vendor-proto/validate:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/validate-repo &&\
	cd vendor-proto/validate-repo &&\
	git sparse-checkout set --no-cone validate &&\
	git checkout
	mkdir -p  vendor-proto
	mv vendor-proto/validate-repo/validate vendor-proto
	rm -rf vendor-proto/validate-repo

generate_checkout: bin vendor-proto/google/api vendor-proto/google/protobuf vendor-proto/validate
	mkdir -p pkg/proto/checkout_v1
	$(PROTOC) -I api/proto/checkout_v1 -I vendor-proto \
	--go_out pkg/proto/checkout_v1 --go_opt paths=source_relative \
	--go-grpc_out pkg/proto/checkout_v1 --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pkg/proto/checkout_v1 --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:pkg/proto/checkout_v1" \
	api/proto/checkout_v1/checkout.proto

generate_product_service: bin vendor-proto/google/api vendor-proto/google/protobuf vendor-proto/validate
	mkdir -p pkg/proto/product_service
	$(PROTOC) -I api/proto/product_service -I vendor-proto \
	--go_out pkg/proto/product_service --go_opt paths=source_relative \
	--go-grpc_out pkg/proto/product_service --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pkg/proto/product_service --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:pkg/proto/product_service" \
	api/proto/product_service/product_service.proto

generate_loms: bin vendor-proto/google/api vendor-proto/google/protobuf vendor-proto/validate
	mkdir -p pkg/proto/loms_v1
	$(PROTOC) -I api/proto/loms_v1 -I vendor-proto \
	--go_out pkg/proto/loms_v1 --go_opt paths=source_relative \
	--go-grpc_out pkg/proto/loms_v1 --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pkg/proto/loms_v1 --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:pkg/proto/loms_v1" \
	api/proto/loms_v1/loms.proto

generate: generate_checkout generate_product_service generate_loms
