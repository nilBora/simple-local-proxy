# Local proxy

Proxy for use local domains. Data come from ngirok host

```
go run app/main.go --target=https://host.local/
```

### Test

`go test -v -run Test_Main`

## Make commands
```
make rundev https://local.host/
```
