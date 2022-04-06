PROJECT_NAME=captcha

build:
	go build -tags=jsoniter -o ./${PROJECT_NAME} ./examples