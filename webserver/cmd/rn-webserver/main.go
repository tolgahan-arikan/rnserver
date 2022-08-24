package main

import (
	"fmt"
	"log"
	"time"
	"webserver"
)

func main() {
	port := 0
	// port := 41214
	fileDir := "/home/peter/Dev/other/rnwserver/cmd/rn-webserver/files"

	// NOTE: here is one way to use the library, where it will start the server
	// and block while running.
	// err := webserver.StartWebServer(port, fileDir)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// NOTE: here is one way to use the library, based on the old webserver.go config.
	// This will start the server in the background and return the server url.
	webserver.SetPort(port)
	webserver.SetConfig(fileDir)
	serverURL, err := webserver.Start()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server listening on url:", serverURL)

	time.Sleep(10 * time.Minute)
}
