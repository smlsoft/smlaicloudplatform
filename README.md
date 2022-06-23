
# SMLCLoudPlatForm

## Environment Variable

### MongoDB
| Name        | Description            | Value |
|-------------|------------------------|-------|
| MONGODB_URI | Mongodb Connection URI | ''    |
| MONGODB_DB  | Mongodb Database Name  | ''    |


### Redis

| Name                 | Description          | Value |
|----------------------|----------------------|-------|
| REDIS_CACHE_URI      | Redis connection uri | ''    |
| REDIS_CACHE_PASSWORD | Redis Password       | ''    |


## For wsl(ubuntu) Please Read

Install Kafkalib , Gcc
```

sudo apt-get install build-essential

```


## gen swagger

```
swag init -g swaggergen.go -o ./api/swagger/

```

### generate authorization pem key

```
openssl genrsa -out private.key 4096
```

```
openssl rsa -in private.key -pubout -out public.key
```

### Run Swagger 
```
swag init
go run main.go
```

```
http://localhost:1323/swagger/index.html
http://localhost:1323/swagger/doc.json

```

### Build Docker Command
```
docker build -t inventoryservice -f cmd/inventoryservice/Dockerfile .
docker build -t smlsoft/cloudauthentication -f ./cmd/authenticationservice/Dockerfile .
```

Get Mock Package
```
go install github.com/vektra/mockery/v2@latest
```

### Run Swagger With Make
```
make runswagger
```



## Github Registry 
## https://ghcr.io

```

```


# FOR M1 Run Please Read

```
brew install openssl
brew install librdkafka
brew install pkg-config
export PKG_CONFIG_PATH="/opt/homebrew/opt/openssl@3/lib/pkgconfig"
go build --tags dynamic main.go

```


## FOR M1 Build Docker and push
```
docker buildx create --use
docker buildx build --platform linux/amd64 --push -t <tag_to_push> .
```

