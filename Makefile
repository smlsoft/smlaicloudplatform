runswagger:
	swag init
	go run main.go

docker_build_authen:
	docker build -t authenticationservice -f ./cmd/authenticationservice/Dockerfile .

docker_build_authen_and_ship:
	docker build -t smlsoft/smlcloudplatform:authen -f ./cmd/authenticationservice/Dockerfile .
	docker push smlsoft/smlcloudplatform:authen

docker_build_swagger_and_ship:
	swag init
	docker build -t smlsoft/smlcloudplatform:swagger .
	docker push smlsoft/smlcloudplatform:swagger

docker_build_shop_and_ship:
	docker build -t smlsoft/smlcloudplatform:shop -f ./cmd/shopservice/Dockerfile .
	docker push smlsoft/smlcloudplatform:shop

docker_build_inventory:
	docker build -t smlsoft/smlcloudplatform:inventory -f ./cmd/inventoryservice/Dockerfile .

docker_build_inventory_and_ship:
	docker build -t smlsoft/smlcloudplatform:inventory -f ./cmd/inventoryservice/Dockerfile .
	docker push smlsoft/smlcloudplatform:inventory

docker_build_inventoryimport_and_ship:
	docker build -t smlsoft/smlcloudplatform:inventoryimport -f ./cmd/inventoryimportservice/Dockerfile .
	docker push smlsoft/smlcloudplatform:inventoryimport

docker_build_masterservice_and_ship:
	docker build -t smlsoft/smlcloudplatform:masterdata -f ./cmd/masterdataservice/Dockerfile .
	docker push smlsoft/smlcloudplatform:masterdata

docker_build_imageservice_and_ship:
	docker build -t smlsoft/smlcloudplatform:imageuploadservice -f ./cmd/imageuploadservice/Dockerfile .
	docker push smlsoft/smlcloudplatform:imageuploadservice

docker_build_imageservice:
	docker build -t smlsoft/smlcloudplatform:imageuploadservice -f ./cmd/imageuploadservice/Dockerfile .

docker_build_memberservice_and_ship:
	docker build -t smlsoft/smlcloudplatform:member -f ./cmd/memberservice/Dockerfile .
	docker push smlsoft/smlcloudplatform:member

run_docker_cluster:
	docker start zookeeper server-redis-1 server-mongodb-1 kafka

docker_build_app_dev:
	swag init
	docker build -t smlsoft/smlcloudplatform:appdev .
	docker push smlsoft/smlcloudplatform:appdev

docker_build_api_dev:
	swag init
	docker build -t smlsoft/smlcloudplatform:apidev .
	docker push smlsoft/smlcloudplatform:apidev

run_app_dev:
	swag init
	go run main.go