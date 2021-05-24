package main

import(
	//config "github.com/Retler/ART/config"
	//tweets "github.com/Retler/ART/tweets"
	handlers "github.com/Retler/ART/statshandlers"
	"fmt"
	"log"
	"net/http"
	//"encoding/json"
	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", handlers.Home)
	myRouter.HandleFunc("/tweet", handlers.GetTweet)
	myRouter.HandleFunc("/tweets", handlers.GetTweetsSince)
	
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
    fmt.Println("Starting StatsAPI")
    handleRequests()
}
