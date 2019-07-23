# sewer
[![Go Report Card](https://goreportcard.com/badge/github.com/0x4139/sewer?style=flat-square)](https://goreportcard.com/report/github.com/0x4139/sewer)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/0x4139/sewer)
[![Release](https://img.shields.io/github/release/golang-standards/project-layout.svg?style=flat-square)](https://github.com/0x4139/sewer/releases/latest)

sewer allows you to use go  channels horizontally with the help of RabbitMQ

## Getting Started

`go get github.com/0x4139/sewer
`
### Prerequisites

sewer only uses a dependency to `	"github.com/streadway/amqp"`

```
module github.com/0x4139/sewer

go 1.12

require github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94

```

### Installing

You can either use go get

```
go get github.com/0x4139/sewer
```

or go mod

```
go mod download  github.com/0x4139/sewer
```


## Running the tests

#Testing using docker

If RabbitMQ is not installed on your machine, the easiest way to install it is using docker container.

## Installing and running RabbitMQ without management tools plugin for the first time

`docker run -d --hostname rabbitmq-test-host --name rabbitmq-test -p 5672:5672 -p 15672:15672 rabbitmq:3`

## Installing and running RabbitMQ with management tools plugin for the first time

`docker run -d --hostname rabbitmq-test-host --name rabbitmq-test -p 5672:5672 -p 15672:15672 rabbitmq:3-management`

## Starting container

`docker start rabbitmq-test`
### Running

From the working directory

```
go test
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **@Vali Malinoiu** - *Initial work* - [0x4139](https://github.com/0x4139)


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
