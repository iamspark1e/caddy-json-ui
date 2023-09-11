package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	utils "caddy-json-ui/utils"
)

// The URL will be augmented with the following query string parameters:
//
//	    server_name = SNI value,
//		signature_schemes = comma-separated list of hex IDs of signature algorithms,
//		and cipher_suites = comma-separated list of hex IDS of cipher suites.
type TokenizedUri struct {
	Token string `uri:"token" binding:"required"`
}

type pluginConfig struct {
	CertDir string `yaml:"cert_dir"`
	Token   string `yaml:"token"`
}

var (
	Trace *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func Logger(message string, reason error) {
	if reason != nil {
		Error.Printf("%s. (%s)", message, reason.Error())
	} else {
		Info.Printf("%s", message)
	}
}

func InitPlugin(plugin_config_yaml_path string) (func(c *gin.Context), error) {
	// FIXME: below queries has been skipped.
	// signature_schemes := c.Query("signature_schemes")
	// cipher_suites := c.Query("cipher_suites")
	gin.SetMode(gin.ReleaseMode)
	plugin_conf := &pluginConfig{}
	load_err := utils.LoadYAML(plugin_config_yaml_path, plugin_conf)
	if load_err != nil {
		return nil, load_err
	}
	return func(c *gin.Context) {
		server_name := c.Query("server_name")
		// Logger(fmt.Sprintf("server_name: %s, signature_schemes: %s, cipher_suites: %s", server_name, signature_schemes, cipher_suites), nil)
		if server_name == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var tokenizedUri TokenizedUri
		if err := c.ShouldBindUri(&tokenizedUri); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if plugin_conf.Token != "" && tokenizedUri.Token != plugin_conf.Token {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		// if server_name can match wildcard cert
		var topDomain string = getTopDomainFromServerName(server_name)
		Logger(fmt.Sprintf("[Caddy TLS Request: %s] A new request from %s, wildcard should be: %s", server_name, c.ClientIP(), topDomain), nil)
		var requestFilePrefix string = plugin_conf.CertDir + "/" + topDomain + "/" + topDomain
		var matchedFiles []string
		if _, err := os.Stat(plugin_conf.CertDir + "/" + topDomain); err != nil {
			if os.IsNotExist(err) {
				Logger(fmt.Sprintf("[Caddy TLS Request: %s] Wildcard cert not found in path: %s", server_name, plugin_conf.CertDir+"/"+topDomain), err)
			} else {
				Logger(fmt.Sprintf("[Caddy TLS Request: %s] Check wildcard cert exist failed", server_name), err)
			}
			// else, try exact server_name cert
			Logger(fmt.Sprintf("[Caddy TLS Request: %s] Fallback to exact server_name", server_name), nil)
			if _, err := os.Stat(plugin_conf.CertDir + "/" + server_name); err != nil {
				Logger(fmt.Sprintf("[Caddy TLS Request: %s] Check exact server_name cert failed", server_name), err)
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			requestFilePrefix = plugin_conf.CertDir + "/" + server_name + "/" + server_name
			matchedFiles = append(matchedFiles, requestFilePrefix+".pem")
		} else {
			matchedFiles = append(matchedFiles, plugin_conf.CertDir+"/"+topDomain+"/"+"fullchain.pem")
		}
		matchedFiles = append(matchedFiles, requestFilePrefix+".pem")
		var rBody string = ""
		for _, matched := range matchedFiles {
			strByte, _ := os.ReadFile(matched)
			rBody += string(strByte) + "\n"
		}
		c.Data(200, "application/x-pem-file", []byte(rBody))
	}, nil
}

func getTopDomainFromServerName(server_name string) string {
	if server_name == "" {
		return ""
	}
	var topDomain string = ""
	fullDomain := strings.Split(server_name, ".")
	topDomain = fullDomain[len(fullDomain)-2] + "." + fullDomain[len(fullDomain)-1]
	return topDomain
}
