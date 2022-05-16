
# SMLCLoudPlatForm

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