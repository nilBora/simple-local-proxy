# Local proxy

Proxy for use local domains. Data come from ngirok host

```
go run app/main.go --target=https://host.local/
```

### Test

`go test -v -run Test_Main`

`go test -v`

## Make commands
```
make rundev https://local.host/
```
