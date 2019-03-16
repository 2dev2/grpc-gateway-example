# gRPC + REST Gateway Play

_Experimentations around grpc-gateway_

**[master branch](https://github.com/stoakes/grpc-gateway-example)** A modified version of the grpc-gateway-example to have both HTTP and gRPC on the same port, without TLS. Server implementation is largely inspired from [Zenithar/go-password](https://github.com/Zenithar/go-password).

**[gorilla-mux](https://github.com/Stoakes/grpc-gateway-example/tree/gorilla-mux)** grpc-gateway-example modified to use `gorilla/mux` as HTTP handler for routes standing aside of grpc-gateway.

## Installation & usage

To try gRPC gateway with gorilla mux router, do this:

```bash
$ go get -u github.com/gorilla/mux
$ go get -u github.com/Stoakes/grpc-gateway-example
$ grpc-gateway-example serve
# In another terminal
$ grpc-gateway-example echo "my first rpc echo"
$ curl -X POST -k http://localhost:10000/v1/echo -H "Content-Type: text/plain" -d '{"value": "foo"}'
$ curl http://localhost:10000/bonjour/foo # foo is matched by gorilla mux URL variable system
```

## Features

- A basic echo service returning value send as parameter. Available on port 10000, in both HTTP & GRPC, without TLS.
- A swagger description of the service at `localhost:10000/swagger.json`
- A swagger UI to interact with the REST API at `localhost:10000/swagger-ui/`
- Advanced HTTP routing with Gorilla Mux. Example at `localhost:10000/bonjour/{your name}`
