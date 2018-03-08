package main

import (
	"flag"
	"os"
	"path"

	"fmt"

	"github.com/schollz/find3/server/main/src/api"
	"github.com/schollz/find3/server/main/src/database"
	"github.com/schollz/find3/server/main/src/mqtt"
	"github.com/schollz/find3/server/main/src/server"
)

func main() {
	externalAddress := flag.String("external", "http://127.0.0.1:8003", "external address")
	aiPort := flag.String("ai", "8002", "port for the AI server")
	port := flag.String("port", "8003", "port for the data (this) server")
	debug := flag.Bool("debug", false, "turn on debug mode")
	mqttFlag := flag.Bool("mqtt", false, "turn on mqtt")

	var dataFolder string
	flag.StringVar(&dataFolder, "data", "", "location to store data")

	var minimumPassive int
	flag.IntVar(&minimumPassive, "min-passive", -1, "minimum number of passive points (passive only)")

	flag.Parse()

	if dataFolder == "" {
		dataFolder, _ = os.Getwd()
		dataFolder = path.Join(dataFolder, "data")
	}
	os.MkdirAll(dataFolder, 0775)

	// setup folders
	database.DataFolder = dataFolder
	api.DataFolder = dataFolder

	// setup debugging
	database.Debug(*debug)
	api.Debug(*debug)
	server.Debug(*debug)
	mqtt.Debug = *debug

	api.AIPort = *aiPort
	server.Port = *port
	if os.Getenv("EXTERNAL_ADDRESS") != "" {
		server.ExternalServerAddress = os.Getenv("EXTERNAL_ADDRESS")
	} else {
		server.ExternalServerAddress = *externalAddress
	}
	server.MinimumPassive = minimumPassive
	server.UseMQTT = *mqttFlag
	err := server.Run()
	if err != nil {
		fmt.Print("error: ")
		fmt.Println(err)
	}
}
