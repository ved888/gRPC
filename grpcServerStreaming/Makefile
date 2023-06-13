dependency:
          go get google.golanmg.org/grpc
          go get -u google.golang.org/protobuf/proto

gen:
   protoc --go_out=pb --proto_path=proto proto/*.proto
   protoc --go_out=pb --proto_path=proto proto/processer_message.proto
   protoc --go_out=pb --go-grpc_out=pb --proto_path=proto proto/*.proto

clean:
    rm pb/*.go

test:
     go test ./...
     go test ./serializer

# Run the project
server:
       go run cmd/server/main.go -port :8080
client:
       go run cmd/client/main.go -address 0.0.0.0:8080

