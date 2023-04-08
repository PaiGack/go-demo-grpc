cd proto

set p1=.
set p2=order_managemant.proto
set p3=../..
set p4=%GOPATH%/src

protoc -I %p4% -I %p1% %p2% --go-grpc_out=%p3% --go_out=%p3%

cd ..


go build ./service
go build ./client