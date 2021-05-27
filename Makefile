BINARY=go-restful
version?=1.0

default::swagger
	@echo Build go-restful API binary
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}
	@echo Build success

mac::swagger
	@echo Build go-restful API binary
	@go build -o ${BINARY}
	@echo Build success

swagger::
	@swag init

image::swagger
	@docker build -t="icyboy/go-restful:${version}" .	
	@docker push icyboy/go-restful:${version}

run::swagger
	air -c .air.toml
