package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var router *gin.Engine
var port int

type upstreamOperation struct {
	Id         string `json: "id"`
	Name       string `json: "name"`
	UpstreamId string `json: "upstreamId"`
}

func init() {
	log.Printf("app initializing")
	gin.SetMode(gin.DebugMode)
	router = gin.Default()

	flag.IntVar(&port, "port", 8081, "app port (default: 8081)")
	flag.Parse()
}

func main() {

	root := router.Group("/")
	{
		root.GET("/", getHomepage)
		root.GET("/ready", getReady)
	}

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/upstream-operation", postUpstreamOperation)
	}

	router.Run(fmt.Sprintf(":%d", port))
}

func getHomepage(ctx *gin.Context) {
	log.Printf("request hit for %s", ctx.FullPath())
	ctx.String(http.StatusOK, "the-upstream\n\n\n[GET] /ready\n[POST] /api/v1/upstream-operation\n%s",
		"{\n \"id\": \"<value>\"\n \"name\": \"<value>\"\n}")
}

func getReady(ctx *gin.Context) {
	log.Printf("request hit for %s", ctx.FullPath())
	ctx.String(http.StatusOK, "")
}

func postUpstreamOperation(ctx *gin.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	log.Printf("request hit for %s with body %v", ctx.FullPath(), string(body))
	var op upstreamOperation
	json.Unmarshal(body, &op)
	op.UpstreamId = fmt.Sprintf("%v", uuid.New())
	ctx.JSON(http.StatusOK, op)
}
