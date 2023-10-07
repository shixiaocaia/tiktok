package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"io/ioutil"
	"math/rand"
	"time"
)

func GetCurrentTime() int64 {
	return time.Now().UnixNano() / 1e6
}

func RandomString() string {
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result := make([]byte, 16)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func GetRequestInfo(c *gin.Context) {
	headers := c.Request.Header
	for key, values := range headers {
		for _, value := range values {
			log.Debugf("Request Headers: %s: %s\n", key, value)
		}
	}

	// 获取请求体
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Debugf("Request Body: %s\n", string(body))

	// 获取请求参数
	params := c.Request.Form
	for key, values := range params {
		for _, value := range values {
			log.Debugf("Request Parameters: %s: %s\n", key, value)
		}
	}
}
