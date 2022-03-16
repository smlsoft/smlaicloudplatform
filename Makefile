runswagger:
	swag init
	go run main.go

docker_build_authen:
	docker build -t authenticationservice -f ./cmd/authenticationservice/Dockerfile .

docker_build_authen_and_ship:
	docker build -t smlsoft/smlcloudplatform:authen -f ./cmd/authenticationservice/Dockerfile .
	docker push smlsoft/smlcloudplatform:authen

docker_build_swagger_and_ship:
	docker build -t smlsoft/smlcloudplatform:swagger .
	docker push smlsoft/smlcloudplatform:swagger

docker_build_shop_and_ship:
	docker build -t smlsoft/smlcloudplatform:shop -f ./cmd/shopservice/Dockerfile .
	docker push smlsoft/smlcloudplatform:shop
