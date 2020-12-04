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
       SET CGO_ENABLED=0  //不设置也可以，原因不明
       SET GOOS=windows
       SET GOARCH=amd64
       cd cmd && go build -o mongodb-job
```

## Windows
**windows to linux**
```    
       SET CGO_ENABLED=0  //不设置也可以，原因不明
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

## 部署信息
```
kk8s环境详情
kk8s预生产环境：
1、集群内部调用地址： master: mongodb-job-service 、worker1: mongodb-job-worker1 、 worker2：mongodb-job-worker2
   集群外部调试地址：  
   
2、 日志查看地址：   

```




```

k8s预生产环境：
1、集群内部调用地址：  
   集群外部调试地址： 
   
2、 日志查看地址：   

```

   