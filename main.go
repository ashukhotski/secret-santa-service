// main.go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ashukhotski/secret-santa-service/service"

	"github.com/gorilla/mux"
)

func main() {
	dir := filepath.Join(".", "logs")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModeDir)
	}
	fname := time.Now().Format("2006-01-02") + ".txt"
	file, err := os.OpenFile(filepath.Join(dir, fname), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	logger := log.New(mw, "secret-santa-service: ", log.LstdFlags|log.Lshortfile)

	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	connString := fmt.Sprintf("mongodb://%s:%s@%s",
		os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
		os.Getenv("DB_ADDRESS"))
	serviceRepo, err := service.NewServiceRepo(connString, os.Getenv("DB_NAME"))
	if err != nil {
		logger.Println(err)
	}

	h := service.NewHandlers(logger, *serviceRepo)

	router := mux.NewRouter()
	h.SetupRoutes(router)

	go func() {
		errChan <- http.ListenAndServe(":8080", router)
	}()

	log.Fatalln(<-errChan)
}
