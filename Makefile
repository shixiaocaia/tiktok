# Go编译器
GO := go

.PHONY: all clean gatewaysvr usersvr videosvr favoritesvr

all: gatewaysvr usersvr videosvr favoritesvr

gatewaysvr:
	cd cmd/gatewaysvr && $(GO) run main.go

usersvr:
	cd cmd/usersvr && $(GO) run main.go

videosvr:
	cd cmd/videosvr && $(GO) run main.go

favoritesvr:
	cd cmd/favoritesvr && $(GO) run main.go

clean:
	# 清理操作...
