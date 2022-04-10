package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	store "the-api/store"
	utils "the-api/utils"

	"github.com/robfig/cron"

	"github.com/gin-gonic/gin"
)

var simulateNetworkPartition bool = false // if set to true: simulate a random network partition event between the-api and the-upstream apps
var networkPartitioned bool = false
var router *gin.Engine
var port int
var upstreamHost string
var origUpstreamPort int
var upstreamPort int

func init() {
	log.Printf("app initializing")
	gin.SetMode(gin.DebugMode)
	router = gin.Default()

	flag.IntVar(&port, "port", 8080, "app port (default: 8080)")
	flag.StringVar(&upstreamHost, "upstreamHost", "127.0.0.1", "the-upstream app host/ip (default: 127.0.0.1)")
	flag.IntVar(&upstreamPort, "upstreamPort", 8081, "the-upstream app port (default: 8081)")
	flag.BoolVar(&simulateNetworkPartition, "simulateNetworkPartition", false, "if set to true: simulate a random network partition event between the-api and the-upstream apps (default: false)")
	flag.Parse()
	log.Printf("the-upstream located at %s:%d", upstreamHost, upstreamPort)
	log.Printf("simulate a network partition between the-api and the-upstream apps: %v", simulateNetworkPartition)
	origUpstreamPort = upstreamPort

	if simulateNetworkPartition {
		networkPartition := cron.New()
		networkPartition.AddFunc("@every 60s", func() {
			networkPartitioned = utils.RandomBool()
			if networkPartitioned {
				log.Printf(">>> [NETWORK PARTITION] the-upstream app will be unreachable for 60s <<<")
				upstreamPort = origUpstreamPort + 1
			} else {
				log.Printf("network is normal")
				upstreamPort = origUpstreamPort
			}
		})
		networkPartition.Start()
	}
}

func main() {

	root := router.Group("/")
	{
		root.GET("/", getHomepage)
		root.GET("/ready", getReady)
	}

	network := router.Group("/network")
	{
		network.GET("/status", getNetworkStatus)
	}

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/operations", getOperations)
		apiV1.GET("/operation/:id", getOperationById)
		apiV1.POST("/operation", postOperation)
	}

	router.Run(fmt.Sprintf(":%d", port))
}

func getHomepage(ctx *gin.Context) {
	ctx.String(http.StatusOK,
		"the-api\n\n\n[GET] /ready\n[GET] /network/status\n[GET] /api/v1/operations\n[GET] /api/v1/operation/:id\n[POST] /api/v1/operation\n%s",
		"{\n \"name\": \"<value>\"\n}")
}

func getReady(ctx *gin.Context) {
	log.Printf("request hit for %s", ctx.FullPath())

	if !isReady() {
		ctx.String(http.StatusServiceUnavailable, "service unavailable")
		return
	}

	ctx.String(http.StatusOK, "")
}

func getNetworkStatus(ctx *gin.Context) {
	log.Printf("request hit for %s", ctx.FullPath())
	status := "normal"
	if networkPartitioned {
		status = "partitioned"
	}
	ctx.JSON(http.StatusOK, gin.H{"status": status})
}

func getOperations(ctx *gin.Context) {
	log.Printf("request hit for %s", ctx.FullPath())

	if !isReady() {
		ctx.String(http.StatusServiceUnavailable, "service unavailable")
		return
	}

	ctx.JSON(http.StatusOK, store.ReadAll())
}

func getOperationById(ctx *gin.Context) {
	log.Printf("request hit for %s", ctx.FullPath())

	if !isReady() {
		ctx.String(http.StatusServiceUnavailable, "service unavailable")
		return
	}

	exists, op := store.ReadById(ctx.Param("id"))
	if !exists {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, op)
}

func postOperation(ctx *gin.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	log.Printf("request hit for %s with body %v", ctx.FullPath(), string(body))

	if !isReady() {
		ctx.String(http.StatusServiceUnavailable, "service unavailable")
		return
	}

	var op store.Operation
	json.Unmarshal(body, &op)

	if utils.IsEmpty(op.Name) {
		ctx.String(http.StatusBadRequest, "missing required param: name")
		return
	}

	ok, upstreamOp := postUpstream(op)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "unable to create operation"})
		return
	}

	ok, op = store.Create(upstreamOp.Name, upstreamOp.UpstreamId)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "unable to create operation"})
		return
	}
	_, op = store.ReadById(op.Id)
	ctx.JSON(http.StatusOK, op)
}

func postUpstream(op store.Operation) (bool, store.Operation) {
	url := fmt.Sprintf("http://%s:%d/api/v1/upstream-operation", upstreamHost, upstreamPort)
	log.Printf("POSTing to %s with data %v", url, op)

	data, _ := json.Marshal(map[string]string{"id": op.Id, "name": op.Name})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	if err != nil {
		log.Printf("%v", err.Error())
		return false /*success*/, store.Operation{}
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var upstreamOp store.Operation
	json.Unmarshal(body, &upstreamOp)

	log.Printf("%s returned HTTP status %d and body %v", url, resp.StatusCode, upstreamOp)
	return (resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMovedPermanently), upstreamOp
}

func isReady() bool {
	url := fmt.Sprintf("http://%s:%d/ready", upstreamHost, upstreamPort)
	log.Printf("GET %s", url)
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("%v", err.Error())
		return false /*ready*/
	}

	defer resp.Body.Close()
	log.Printf("%s returned HTTP status %d", url, resp.StatusCode)
	return resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMovedPermanently
}
