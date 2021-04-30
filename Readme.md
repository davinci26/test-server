## Simple HTTP(s) server used for testing

### Features
1. `HTTPS` server listening on port 5055
2. `HTTP/2` cleartext listening on port 5050.


### Endpoints
1. `\ping`: simple endpoint to check if the server is up and running
2. `\request?size=<int>`: returns a response with fixed size of bytes in the body. Note: A byte array is allocated per request
which means that you can get the server OOM with a big enough size.


### Build

```
go build
```