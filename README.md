# Caddy JSON UI

## Quick Start

### Docker

```bash
docker run -d -it -e CADDY_API_ENDPOINT="http://127.0.0.1:2019" -p 8045:8045 --name my-json-ui caddy-json-ui
```

> It's suggested to use `docker compose` with a "Caddy Container".

### Binary

```bash
./caddy-json-ui --port 30081 --caddyadminapi "http://127.0.0.1:2019"
```

## Advanced usage

### Generate a schema file for Editor IntelliSense?

> How to use [xcaddy](https://github.com/caddyserver/xcaddy) ?

```bash
xcaddy build \
    --with github.com/abiosoft/caddy-json-schema \
    # any other module you want to include in the generated schema

caddy json-schema --vscode
```

And put the `.vscode/caddy_schema.json` into `public` folder.

Declare the `"$scheme"` key and use full url to refer it.

> Container image "ghcr.io/iamspark1e/caddy-json-ui" already has a schema file under "./public" by using Github Action, you can directly use it with `{"$schema": "http://127.0.0.1:8045/caddy_schema.json"}` (replace the url in your condition).

### Using plugin econosystem

An example `plugins.yml`,

```yml
enable: true
plugins:
  - name: http-tls-provider
    route: tls-provider # So the real endpoint url will be `localhost:8045/plugins/tls-provider`
    config_path: ./plugins/http-tls-provider/config.yml
```

A plugin's `*.so` file will expose an `func InitPlugin(plugin_config_yaml_path string) (func(c *gin.Context), error)`

## Development

### Run

> Remember to allow UI server access caddy server if you have firewall running

```bash
go run main.go --port 30081 --caddyadminapi "http://127.0.0.1:2019"
```

### Build

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-w -s" -tags=nomsgpack .
```