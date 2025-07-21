PROJECT_NAME = ms_exchange
OS = linux
ARCH = amd64
BUILD_FROM = ./cmd/${PROJECT_NAME}
BUILD_TO = /app/${PROJECT_NAME}

init:
	go mod init ${PROJECT_NAME} && go mod tidy

build:
	GOOS=${OS} GOARCH=${ARCH} CGO_ENABLED=1 go build -a -installsuffix cgo -ldflags="-w -s" -o ${BUILD_TO} ${BUILD_FROM}
