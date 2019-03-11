# nuclio-php-runtime
Runtime for php applications in [nuclio](https://nuclio.io/)

This repository provides a small shim written in golang, which forwards a [nuclio event](https://github.com/nuclio/nuclio-sdk-go/blob/master/event.go) via fcgi to `php-fpm` through a unix socket.

### How does it work

On the first request (coldstart) received by the function a php-fpm server will be started.
  The php-fpm server is configured to listen on a linux socket ( `/var/task/fpm.sock`).
  The request will be forwarded to fpm via the socket and the response will be returned to the clients
  
Subsequent request to the function will profit from the already running `php-fpm` server and handle requests faster.

### Configuration & Caveats

The current version is opinionated:

- Per default we assume that a central php script is provided which all requests will be handled (it is responsible for application internal routing).
  The script can be defined in the `Dockerfile` (`PHP_SCRIPT`) or later via the environment vars in the nuclio interface
- It assumes that the working directory of the function is `/var/task` - thus the fpm configuration file as well as the source code is put there
- The `Dockerfile` utilizes the standard `php:fpm` container as the basis and adds the nuclio specific options. If you use a different container, please ensure you set `PHP_FPM_BIN` to the correct path.



### Example

The repository contains a example `Dockerfile` on how to utilize the handler.
It can easily be deployed by building the docker image via:

```
cd example && docker build -t nuclio-php-example:latest
```

Once it is build, it can be deployed with the following nuclio configuration

```
apiVersion: "nuclio.io/v1"
kind: NuclioFunction
metadata:
  name: php-example
spec:
  image: nuclio-php-example:latest
  handler: main:Handler
  runtime: golang
```

For more information on deploying functions via Dockerfiles please refer to the nuclio documentation:
- https://nuclio.io/docs/latest/tasks/deploy-functions-from-dockerfile/
- https://nuclio.io/docs/latest/tasks/deploying-pre-built-functions/

There is also a [demo project](https://github.com/patrickjahns/nuclio-symfony) utilizing the symfony framework


### Other Projects
- [bref.sh](https://bref.sh/)
  PHP Runtime for AWS Lambda


### Authors

* [Patrick Jahns](https://github.com/patrickjahns)


### License

Apache


### Copyright

```
Copyright (c) 2019 Patrick Jahns 
```