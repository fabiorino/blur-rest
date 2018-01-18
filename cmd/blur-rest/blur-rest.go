package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"blur-rest/config"
	"blur-rest/handlers"
	"blur-rest/store"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
)

var options struct {
	Fqdn   string `short:"f" long:"fqdn" description:"The Fully Qualified Domain Name of the server" required:"true"`
	Port   string `short:"p" long:"port" description:"The listening port for incoming requests" required:"true"`
	BoltDB string `short:"b" long:"boltdb" description:"Absolute path of the BoltDB database file" default:"images.db"`
}

func init() {
	// Parse arguments
	_, err := flags.Parse(&options)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Create logs folder
	logsFolderName := "logs"
	if err := os.MkdirAll(logsFolderName, os.ModeDir); err != nil {
		log.Fatalln("Could not create log folder")
	}

	// Create log file. The name of the file is the current timestamp in the form DDMMYYYYHHMMSS
	t := time.Now()
	timestamp := t.Format("02012006150405")
	logFileName := timestamp + ".txt"
	logFile, err := os.Create(filepath.Join(logsFolderName, logFileName))
	if err != nil {
		errMsg := fmt.Sprintf("Could not create log file: %s", err.Error())
		log.Fatalln(errMsg)
	}

	// Write log messages on the file we've just created and the stdout
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	// Init store
	store, err := store.NewBoltDB(options.BoltDB)
	if err != nil {
		errMsg := fmt.Sprintf("Could not init BoltDB: %s", err.Error())
		log.Fatalln(errMsg)
	}

	// Setup global config
	config.GlobalConfig.Fqdn = options.Fqdn
	config.GlobalConfig.Port = options.Port
	config.GlobalConfig.Store = store

	router := gin.Default()

	// Setup endpoints
	router.POST("/blur", handlers.PostImageMetaHandler)
	router.PUT("/blur/{id}", handlers.UploadImageHandler)

	// Setup server parameters
	srv := &http.Server{
		Addr:    ":" + options.Port,
		Handler: router,
	}

	// Start the server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errMsg := fmt.Sprintf("Error in ListenAndServe: %s", err.Error())
			log.Fatalln(errMsg)
		}
	}()

	handleShutdown(srv)
}

// handleShutdown waits for interrupt signal to gracefully shutdown the server with a timeout of 30 seconds
func handleShutdown(srv *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Error during server shutdown: ", err)
	}
}
