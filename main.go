package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	bootstrap "caddy-json-ui/bootstrap"
	internal "caddy-json-ui/internal"
	middleware "caddy-json-ui/middleware"
	utils "caddy-json-ui/utils"

	"github.com/gin-gonic/gin"

	_ "net/http/pprof"
)

func stashLatestCaddyConfig(jsonContent []byte) error {
	err := os.WriteFile("./config/config.json", jsonContent, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	endpoint := os.Getenv("CADDY_API_ENDPOINT")
	if endpoint == "" {
		endpoint = bootstrap.CaddyAPIEndpoint
	}
	caddySrv := internal.NewCaddyServer(endpoint)
	exist, existErr := utils.PathExists("./config/config.json")
	if existErr != nil {
		log.Print("Stat config dir failed")
		log.Panic(existErr.Error())
		return
	}

	if !exist {
		err := os.Mkdir("./config", os.ModePerm)
		if err != nil {
			log.Print("Mkdir config dir failed")
			log.Panic(err.Error())
			return
		} else {
			err = os.WriteFile("./config/config.json", []byte("{\"admin\": {\"listen\": \"127.0.0.1:2019\"}}"), os.ModePerm)
			if err != nil {
				log.Print("Write config file failed")
				log.Panic(err.Error())
				return
			}
		}
	}

	gin.SetMode(gin.ReleaseMode)
	apiEng := gin.New()
	api := apiEng.Group("/api")
	{
		api.POST("/save", func(c *gin.Context) {
			jsonByte, err := io.ReadAll(c.Request.Body)
			if err != nil {
				// Handle error
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"ok":  0,
					"msg": err.Error(),
				})
			}
			err = caddySrv.SetCaddyConfig(jsonByte)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"ok":  0,
					"msg": err.Error(),
				})
				return
			}
			stashLatestCaddyConfig(jsonByte)
			c.JSON(http.StatusOK, gin.H{
				"ok":  1,
				"msg": "Success",
			})
		})
		api.GET("/load", func(c *gin.Context) {
			data, err := caddySrv.GetCaddyConfig()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"ok":  0,
					"msg": err.Error(),
				})
				return
			}
			c.JSON(200, gin.H{
				"ok":   1,
				"msg":  "Success",
				"data": string(data),
			})
		})
		api.GET("/healthy", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"ok":  1,
				"msg": "healthy",
			})
		})
	}
	pubEng := gin.New()
	pubEng.Static("/", "./public")
	r := gin.Default()
	r.Use(middleware.MidWhitelistFilter(bootstrap.WhitelistPath))
	r.Any("/*any", func(c *gin.Context) {
		path := c.Param("any")
		if strings.HasPrefix(path, "/api") {
			apiEng.HandleContext(c)
		} else {
			pubEng.HandleContext(c)
		}
	})

	e := r.Run(fmt.Sprintf(":%d", bootstrap.Port))
	if e != nil {
		log.Fatal(e)
	}
}
