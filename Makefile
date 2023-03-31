
gen.proto:
	protoc --go_out=. --go-grpc_out=. grpc/gpt.proto