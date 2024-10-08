
## Nats

### download

https://docs.nats.io/running-a-nats-service/introduction/installation

### start

```
./nats-server -m 8222 -js
```

API
```
http://localhost:8222/connz
```

## Running the example

### Cloud Example
```
go run main.go --uuid=server --port=:8080
```

### Edge Example
```
go run main.go --uuid=edge --port=:8081
```

Ping all clients connected to the `nats-server`

```
http://0.0.0.0:8080/hosts/remote/ping/all
```

Ping a client by its `uuid` on the `nats-server`

```
http://0.0.0.0:8080/hosts/remote/ping/<UUID>
```



## Module Example

start the module example `modules/module-abc/main.go`


```curl
curl -X POST http://localhost:8080/modules/my-module \
-H "Content-Type: application/json" \
-d '{
    "requestUUID": "abc123",
    "method": "hello",
    "payload": "hello there"
}'

```

using the nats `cli`

wild card `sub`
```
./nats sub ">"
```

publish in 2nd terminal
```
./nats req module.my-module.hello "whats up"
```