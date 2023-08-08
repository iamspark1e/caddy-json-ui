package bootstrap

import (
	"flag"
)

var (
	LogLevel         string
	WhitelistPath    string
	Port             int
	CaddyAPIEndpoint string
)

func init() {
	flag.IntVar(&Port, "port", 8045, "Running port")
	flag.StringVar(&CaddyAPIEndpoint, "caddyadminapi", "http://127.0.0.1:2019", "The Caddy Server for Manager communication")
	flag.StringVar(&WhitelistPath, "whitelist", "./whitelist", "Use white list of IPs to filter trusted client, localhost(127.0.0.1 & ::1) is allowed by default")
	flag.StringVar(&LogLevel, "loglevel", "info", "Log level, trace|info|warn|error")
	flag.Parse()
}
