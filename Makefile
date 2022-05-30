swagger:
	GO111MODULE=on swagger generate spec -o ./swagger.yaml --scan-models

proto:
	protoc -I protos/ protos/product/product.proto --go_out=plugins=grpc:protos/product