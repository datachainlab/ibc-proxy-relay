
.PHONY: build
build:
	go build -o ./build/uly ./main.go

.PHONY: test
test:
	go test -v ./modules/...

.PHONY: proto-gen
proto-gen:
	@echo "Generating Protobuf files"
	docker run -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen.sh
