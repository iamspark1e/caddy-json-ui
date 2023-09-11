package mids

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func in(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}

func MidWhitelistFilter(whitelistPath string) gin.HandlerFunc {
	// Parse whitelist
	var whitelist []string
	whitelist = append(whitelist, "127.0.0.1")
	whitelist = append(whitelist, "::1") // allow localhost by default
	file, err := os.Open(whitelistPath)
	if err != nil {
		log.Print("open whitelist failed")
	} else {
		reader := bufio.NewReader(file)
		for {
			bt, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			str := strings.Replace(string(bt), "\n", "", -1)
			whitelist = append(whitelist, str)
		}
	}
	defer file.Close()

	return func(c *gin.Context) {
		client := c.ClientIP()
		result := in(client, whitelist)
		if !result {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Next()
	}
}
