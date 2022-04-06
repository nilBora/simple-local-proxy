# Local proxy
[![Build Status](https://github.com/nilBora/simple-local-proxy/workflows/Go/badge.svg)](https://github.com/nilBora/simple-local-proxy/actions)

Proxy for use local domains. Data come from ngirok host

```
go run app/main.go --target=https://host.local/
```

### Test

`go test -v -run Test_Main`

`go test -v ./...`

## Make commands
```
make rundev https://local.host/
```
