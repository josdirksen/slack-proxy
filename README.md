#Slack proxy

This project provides a simple proxy betwee slack and docker. To develop on this project make sure your 
 GOPATH is set correctly and install the required go-dockerclient package.

```
 go get github.com/fsouza/go-dockerclient
```

Alternatively you can just use docker to run the application:

``` 
$ docker build .
Sending build context to Docker daemon 45.16 MB
Step 1 : FROM golang:latest
 ---> bc422006801e
Step 2 : COPY src/ /go/src/
 ---> 7ba485bdb838
Removing intermediate container 8cc27ecc444f
Step 3 : COPY resources/config.json /
 ---> aec547bee967
Removing intermediate container 0f1a3b96b023
Step 4 : EXPOSE 9000
 ---> Running in 0ad263c5c74e
 ---> 9df08e9a3708
Removing intermediate container 0ad263c5c74e
Step 5 : RUN go get github.com/fsouza/go-dockerclient && go install github.com/fsouza/go-dockerclient
 ---> Running in c9969ed17c30
 ---> 7f7461052097
Removing intermediate container c9969ed17c30
Step 6 : WORKDIR src
 ---> Running in 93f650e9c638
 ---> 622aabe8136e
Removing intermediate container 93f650e9c638
Step 7 : RUN go build -o /app/main github.com/josdirksen/slackproxy/slack-proxy.go
 ---> Running in c30989ca1535
 ---> c49c5c5cfccc
Removing intermediate container c30989ca1535
Step 8 : CMD /app/main --config /config.json
 ---> Running in 459d0a44b4fd
 ---> 82e4ede226e9
Removing intermediate container 459d0a44b4fd
Successfully built 82e4ede226e9
 
$ docker run 82e4ede226e9
...
```