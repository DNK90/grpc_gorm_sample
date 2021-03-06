* Install

```
go install \
    github.com/golang/protobuf/protoc-gen-go \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
    github.com/infobloxopen/protoc-gen-gorm
```

After install, `protoc-gen-go` `protoc-gen-grpc-gateway` `protoc-gen-swagger` and `protoc-gen-gorm` will be added into `$GOPATH/bin`

* In order to generate `sample.pb.go` `sample.pb.gorm.go` and `sample.pb.gw.go` run the following script

```shell script
protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. --gorm_out=. --grpc-gateway_out=logtostderr=true:. -grpc-gateway_out=logtostderr=true:. sample/sample.proto sample/sample.proto
```

* For more detail on configuring gorm in protobuf, refer the following link:
https://github.com/infobloxopen/protoc-gen-gorm

* Run `dep ensure` to install all dependencies then run `main.go`

* Add Item
```
POST /v1/addItem
BODY
{
    "id": "123"
    "name": "item1"
    "description": "itemDescription"
}
```

* Get Item By ID

```
POST /v1/getItem
BODY
{
    "id": "123"
}
```

* List all items

```
POST /v1/listItems
```
