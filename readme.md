#CMPE 273 - Lab 3 Part I - Server Side Cache Data Store
## Building CMPE 273 - Lab 3 (Part I)
```
cd CMPE273-lab3/CMPE273-lab3-server
export GOPATH=$PWD # Set GO PATH
go get github.com/julienschmidt/httprouter #Get dependencies
go build src/server/Server.go # Build Server
```

## Run the server using 

```
go run src/server/Server.go
```

or

```
./Server
```

## API Calls

### PUT key request

Put key/value pair
```
curl -v -X PUT  http://localhost:3000/keys/{key-id}/{value}
```

Sample response is : HTTP Response CODE : 200

### GET key request

Get key based on id
```
curl -X GET -v http://localhost:3000/keys/{key-id}
```
e.g.
```
curl -X GET -v http://localhost:3000/keys/2
```

Sample response is : HTTP Response CODE : 200

```
{ "key" : 2, “value” : "b" }
```

### GET All keys request

Get all keys
```
curl -X GET -v http://localhost:3000/keys
``` 

Sample response is : HTTP Response CODE : 200
```
[{ 
   "key" : 1,
   "value" : "a"
 },
 {
   "key" : 2,
   "value" : "b"
 }
]
```

#CMPE 273 - Lab 3 Part II - Consistent Hashing on Client Side
## Building CMPE 273 - Lab 3 (Part II)
Before running Client, please ensure that the above Server instance (from Part I) is running on three different ports (e.g. 3000, 3001, 3002).
For doing the same, make two more copies of 'CMPE273-lab3-server' folder and change below code (with appropriate port added) in src\server\Server.go
```
    server := http.Server{
            Addr:        "0.0.0.0:3000",
            Handler: mux,
    }
```

This Client uses the Consistent Hashing Algorithm to hash above Server instances and keys (using MD5 as hashing function). It shards the data set among these Servers.

Please note that this Client could be run as REST client or as a STANDALONE GO Application

```
cd CMPE273-lab3/CMPE273-lab3-client
export GOPATH=$PWD # Set GO PATH
go get github.com/julienschmidt/httprouter #Get dependencies
go build -v ./...  # Build the project
go build src/client/Client.go # Build Server
```

## Run the Client using 

```
go run src/client/Client.go
```

or

```
./Client
```
## API Calls with REST Client

### PUT key request

Put key/value pair
```
curl -v -X PUT  http://localhost:8080/keys/{key-id}/{value}
```

Sample response is : HTTP Response CODE : 200

### GET key request

Get key based on id
```
curl -X GET -v http://localhost:8080/keys/{key-id}
```
e.g.
```
curl -X GET -v http://localhost:8080/keys/2
```

Sample response is : HTTP Response CODE : 200

```
{ "key" : 2, “value” : "b" }
```

## The Outputs with Running Client as STANDALONE Go Application

Please follow the loggers on console (For both - Client as well as Server instances) to verify how data is getting stored and retrieved from three Server instances.
Alternatively, one can also make GET calls to respective Server instance to ensure that the keys are retrieved from there correctly.