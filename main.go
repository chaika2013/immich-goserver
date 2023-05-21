package main

import (
	"flag"
	"log"

	"github.com/chaika2013/immich-goserver/model"
	"github.com/chaika2013/immich-goserver/pipeline"
	"github.com/chaika2013/immich-goserver/router"
	"github.com/chaika2013/immich-goserver/session"
	"github.com/gin-gonic/gin"
)

func main() {
	flag.Parse()

	// initialize database
	err := model.Setup()
	if err != nil {
		log.Fatal(err.Error())
	}

	// run processing pipeline
	pipeline.Setup()

	// run gin
	gin := gin.Default()

	// working with session
	gin.Use(session.Setup())

	// setup routers
	router.Setup(gin)
	gin.Run(":8080")
}
