PROTO_SRC=proto/user.proto

gen:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       $(PROTO_SRC)

tidy:
	go mod tidy

run:
	SEED=$${SEED:-true} go run ./cmd/server
