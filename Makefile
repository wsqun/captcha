PROJECT_NAME=captcha

build:
	go build -tags=jsoniter -o ./${PROJECT_NAME} ./examples

buildCompare:
	go build -tags=jsoniter -o ./${PROJECT_NAME} ./example_compare

buildImg:
	docker build -t wsqun/svc-captcha .
	docker tag wsqun/svc-captcha wsqun/svc-captcha:1.0.1

buildImgCompare:
	docker build -t wsqun/svc-captcha-compare .
	docker tag wsqun/svc-captcha-compare wsqun/svc-captcha-compare:1.0.1