cd proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative product_info.proto
cd ..

go build ./client
go build ./server
