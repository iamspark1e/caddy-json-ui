### Build as plugins

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=plugin -ldflags="-w -s" -tags=nomsgpack .
```

> plugin包和普通的Go包没太多区别，只是plugin包有一个约束：其包名必须为main。

### Docs

> https://zh.mojotv.cn/go/golang-plugin-tutorial