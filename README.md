## Install & Run
**System requirements:**

**On a Linux host:** docker 17.06.0-ce+ and docker-compose 1.18.0+ .
## Install
  docker build -t="mongo-job:version" .
## Run
**go run: master**
  docker run -d --name mongo-job-master -v /app/logs/job/master:/app/logs  -p ${port}:8341 -e "MJOBUNIQUEKEY=${uuid}"  mongo-job:${version} ./mongodb-job master --conf=/app/mongodb-job.yaml --serviceName=${mongodb-service-address}
**go run: worker1**
  docker run -d --name mongo-job-worker1 -v /app/logs/job/worker1:/app/logs -p ${port}:8341  -e "MJOBUNIQUEKEY=${uuid}" mongo-job:${version} ./mongodb-job worker --conf=/app/mongodb-job.yaml --serviceName=${mongodb-service-address}
**go run: worker2**
  docker run -d --name mongo-job-worker2 -v /app/logs/job/worker2:/app/logs -p ${port}:8341  -e "MJOBUNIQUEKEY=${uuid}" mongo-job:${version} ./mongodb-job worker --conf=/app/mongodb-job.yaml --serviceName=${mongodb-service-address}

## Build 
**System requirements:** 
**go version: 1.14+**
## Linux 
**linux to windows**
```    
       SET CGO_ENABLED=0  
       SET GOOS=windows
       SET GOARCH=amd64
       cd cmd && go build -o mongodb-job
```

## Windows
**windows to linux**
```    
       SET CGO_ENABLED=0  
       SET GOOS=linux
       SET GOARCH=amd64
       cd cmd && go build -o mongodb-job
       #go build -gcflags=-trimpath=${GOPATH}-asmflags=-trimpath=${GOPATH}
```
## Mac
**Mac to linux**
```    
       CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mongodb-job
```
   
