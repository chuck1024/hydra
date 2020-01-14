module hydra

go 1.13

require (
	github.com/chuck1024/doglog v0.0.0-20200114052321-1297eb7c152e
	github.com/chuck1024/gl v0.0.0-20200114031106-5a09ab9144f9
	github.com/chuck1024/godog v0.0.0-20200114053715-b4dab07db07d
	github.com/chuck1024/hydra v0.0.0-20190623085650-e100634444e2
	github.com/chuck1024/redisdb v0.0.0-20190617091652-c849607cda9f
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/gin-gonic/gin v1.5.0
	github.com/go-redis/redis v6.15.6+incompatible // indirect
	github.com/gorilla/websocket v1.2.1-0.20180605202552-5ed622c449da
	google.golang.org/genproto v0.0.0-20191223191004-3caeed10a8bf // indirect
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
