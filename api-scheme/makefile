api:
	protoc -I ./proto --go_out=proto/src/golang --go_opt=paths=source_relative \
    --go-grpc_out=proto/src/golang --go-grpc_opt=paths=source_relative \
	proto/*.proto
	

json:
	protoc -I . --openapiv2_out proto/json \
	--openapiv2_opt logtostderr=true proto/create_link.proto