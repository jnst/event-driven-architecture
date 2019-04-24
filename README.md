# event-sourcing

## Concept

### Publisher

Publisher are responsible for messaging using Amazon SNS.

### Subscriber

Subscribers are responsible for polling for queuing.

### Butler

Butler does chores such as creation, destruction of SNS/SQS.  
FYI: This has nothing to do with the concept of the Event sourcing architecture. 

## Setup

Requires version 1.12 or higher to use go modules.

```bash
$ go version
go version go1.12.4 darwin/amd64
```

Since aws-sdk-go needs credentials, please set "dummy" if you use localstack.

```bash
$ cat ~/.aws/credentials

[dummy]
aws_access_key_id = dummy
aws_secret_access_key = dummy 
```

## Run

Start docker for localstack.

```bash
$ docker-compose up -d
$ env GO111MODULE=on AWS_PROFILE=dummy go run main.go
```

Stop docker.

```bash
$ docker-compose down
```
