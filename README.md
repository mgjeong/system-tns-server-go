System TNS Server
=======================================

This provides functionalities to register, retreive, unregister topic information that could be used for ezMQ-plus publishers/subscribers.

## Prerequisites ##
- docker-ce
    - Version: 17.09
    - [How to install](https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/)
- go compiler
    - Version: 1.8 or above
    - [How to install](https://golang.org/dl/)

## How to build ##
This provides how to build sources codes to an excutable binary and dockerize it to create a Docker image.

#### 1. Executable binary ####
```shell
$ ./build.sh
```
If source codes are successfully built, you can find an output binary file, **tns-server**, on a root of project folder.

#### 2. Docker Image  ####
Next, you can create it to a Docker image.
```shell
$ docker build -t system-tns-server-go -f Dockerfile .
```
If it succeeds, you can see the built image as follows:
```shell
$ sudo docker images
REPOSITORY                         TAG        IMAGE ID        CREATED           SIZE
system-tns-server-go               latest     f8123bce4802    29 minutes ago    159MB
```

## How to run with Docker image ##
Required options to run Docker image
- port
    - 48323:48323
- volume
    - "host folder"/data/db:/data/db (Note that you should replace "host folder" to a desired folder on your host machine)

You can execute it with a Docker image as follows:
```shell
$ docker run -it -p 48323:48323 -v /tns-server/data/db:/data/db system-tns-server-go
```
If it succeeds, you can see log messages on your screen as follows:
```shell
$ docker run -it -p 48323:48323 -v /tns-server/data/db:/data/db system-tns-server-go
2018-05-08T06:46:02.558+0000 I CONTROL  [initandlisten] MongoDB starting : pid=6 port=27017 dbpath=/data/db 64-bit host=50db434b5682
2018-05-08T06:46:02.558+0000 I CONTROL  [initandlisten] db version v3.4.4
2018-05-08T06:46:02.558+0000 I CONTROL  [initandlisten] git version: 888390515874a9debd1b6c5d36559ca86b44babd
2018-05-08T06:46:02.558+0000 I CONTROL  [initandlisten] OpenSSL version: LibreSSL 2.5.5
2018-05-08T06:46:02.558+0000 I CONTROL  [initandlisten] allocator: system
2018-05-08T06:46:02.558+0000 I CONTROL  [initandlisten] modules: none
2018-05-08T06:46:02.558+0000 I CONTROL  [initandlisten] build environment:
2018-05-08T06:46:02.558+0000 I CONTROL  [initandlisten]     distarch: x86_64
2018-05-08T06:46:02.559+0000 I CONTROL  [initandlisten]     target_arch: x86_64
2018-05-08T06:46:02.559+0000 I CONTROL  [initandlisten] options: { repair: true }
2018-05-08T06:46:02.559+0000 W -        [initandlisten] Detected unclean shutdown - /data/db/mongod.lock is not empty.
2018-05-08T06:46:02.570+0000 I -        [initandlisten] Detected data files in /data/db created by the 'wiredTiger' storage engine, so setting the active storage engine to 'wiredTiger'.
2018-05-08T06:46:02.570+0000 W STORAGE  [initandlisten] Recovering data from the last clean checkpoint.

...

2018-05-08T06:46:03.396+0000 I FTDC     [initandlisten] Initializing full-time diagnostic data capture with directory '/data/db/diagnostic.data'
2018-05-08T06:46:03.396+0000 I NETWORK  [thread1] waiting for connections on port 27017
2018-05-08T06:46:03.829+0000 I NETWORK  [thread1] connection accepted from 127.0.0.1:34936 #1 (1 connection now open)
[DEBUG][TNSSVR]2018/05/08 06:46:03 tns/controller/keepalive.Executor keepalive.go InitKeepAlive : 55 [IN]
[DEBUG][TNSSVR]2018/05/08 06:46:03 tns/controller/keepalive.Executor keepalive.go InitKeepAlive : 77 [OUT]
[DEBUG][TNSSVR]2018/05/08 06:46:03 tns/controller/keepalive keepalive.go keepAliveTimerLoop : 150 [Start KeepAlive Timer loop]
2018-05-08T06:46:04.015+0000 I FTDC     [ftdc] Unclean full-time diagnostic data capture shutdown detected, found interim file, some metrics may have been lost. OK

```
## API Document ##
TNS Server provides a set of REST APIs for its operations. Descriptions for the APIs are stored in <root>/doc folder.
- **[tns.yaml](https://github.sec.samsung.net/RS7-EdgeComputing/system-tns-server-go/blob/master/doc/tns.yaml)**

Note that you can visit [Swagger Editor](https://editor.swagger.io/) to graphically investigate the REST APIs in YAML.
