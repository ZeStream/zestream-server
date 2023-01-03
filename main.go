package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"zestream/constants"
	"zestream/routes"

	"github.com/joho/godotenv"
)

func main() {
	e := godotenv.Load()

	if e != nil {
		fmt.Print(e)
	}

	r := routes.Init()

	port := os.Getenv(constants.PORT)

	if port == "" {
		port = constants.DEFAULT_PORT
	}
	router := http.NewServeMux()
	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println(err)
	}

	r.Run(":" + port)
}
