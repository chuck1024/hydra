.PHONY: run docker_run docker

GOPROXY := https://goproxy.cn,direct
GO111MODULE := on

export GO111MODULE
export GOPROXY

default: run

run:
	go run main.go

docker_run:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hydra main.go
	docker build -t hydra .
	rm -rf hydra
	docker run -p 10240:10240 -d hydra ./server.sh

docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hydra main.go
	docker build -t registry.cn-chengdu.aliyuncs.com/godog/hydra .
	docker push registry.cn-chengdu.aliyuncs.com/godog/hydra
	rm -rf hydra
