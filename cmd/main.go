package main

import (
	"log"
	"net/http"
	"os"
    
	"github.com/joho/godotenv"
	mw "github.com/saeidalz13/verticalvibes/middlewares"
	"github.com/saeidalz13/verticalvibes/token"
	"github.com/saeidalz13/verticalvibes/routes"
)

func main() {
	logger := log.New(os.Stdout, "MAIN: ", log.Ldate|log.Ltime|log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	tokenManager, err := token.BuildTokenManager()
	if err != nil {
		logger.Fatalln(err)
	}

	mux := http.NewServeMux()
    routes.Setup(mux)
	middleWare := mw.NewMiddlewareHanlder(logger, tokenManager)

	wrappedMux := middleWare.Chain(mux)

	hostname := os.Getenv("HOSTNAME")
	port := os.Getenv("PORT")
	addr := hostname + ":" + port

	server := &http.Server{
		Addr:    addr,
		Handler: wrappedMux,
	}

	logger.Printf("Listening to %s...\n", addr)
	logger.Fatalln(server.ListenAndServe())
}

