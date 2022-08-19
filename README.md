# Go microservice

### Run front-end locally
``` shell
cd ./fron-end
go run cmd/web/main.go # access http://localhost:80
```

### Run broker service locally
``` shell
cd ./broker-service
go run cmd/api # access http://localhost:80
```

### Run project in docker
```
cd project
docker-compose up -d 
# Access http://localhost:8080 # broker service
```
