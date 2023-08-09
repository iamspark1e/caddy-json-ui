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

	"github.com/gin-gonic/gin"
)

func main() {
	endpoint := os.Getenv("CADDY_API_ENDPOINT")
	if endpoint == "" {
		endpoint = bootstrap.CaddyAPIEndpoint
	}
	caddySrv := internal.NewCaddyServer(endpoint)
	// check if not config.json exists
	file, err := os.Open("./config/config.json")
	if err != nil {
		if _, err := os.Stat("./config"); os.IsNotExist(err) {
			err = os.Mkdir("./config", 0777)
			if err != nil {
				log.Print("Make config dir failed")
				return
			}
		}
		err = os.WriteFile("./config/config.json", []byte("{\"admin\": {\"listen\": \"127.0.0.1:2019\"}}"), 0777)
		if err != nil {
			log.Print("Write default config failed")
			log.Print(err.Error())
		}
	}
	defer file.Close()

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
			}
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

	err = r.Run(fmt.Sprintf(":%d", bootstrap.Port))
	if err != nil {
		log.Fatal(err)
	}
}
