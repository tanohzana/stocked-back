package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron"
	mapFunctions "github.com/tanohzana/stockedBackend/utils"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}
var serverPort = ":8080"
var defaultStocks = "aapl,fb,tsla"
var iexURL = "https://api.iextrading.com"

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)

	go handleStockChangesWithCron()

	log.Println("Magic is happening on port ", serverPort)
	err := http.ListenAndServe(serverPort, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Headers nil here
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true
}

func handleStockChangesWithCron() {
	c := cron.New()
	var prevStocks map[string]interface{}

	c.AddFunc("@every 10s", func() {
		cronTask(&prevStocks)
	})

	c.Start()
}

func cronTask(prevStocks *map[string]interface{}) {
	fmt.Println("Fetching stocks")
	resp, err1 := http.Get(iexURL + "/1.0/stock/market/batch?symbols=" + defaultStocks + "&types=quote,chart&range=1m&last=5")

	if err1 != nil {
		log.Fatal(err1)
	}

	stocks, err2 := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err2 != nil {
		log.Fatal(err2)
	}

	if len(*prevStocks) > 0 {
		stocksJson := getJsonFormatFromString(string(stocks))
		filteredMap := mapFunctions.Filter(*prevStocks, func(key string, value interface{}) bool {
			// implement this
			return true
		})

		fmt.Println(stocksJson)
		fmt.Println(filteredMap)
		// io emit here
	} else {
		*prevStocks = getJsonFormatFromString(string(stocks))
	}
}

func getJsonFormatFromString(jsonString string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonString), &result)

	return result
}

func getDifferencesInStocks(prevStocks, newStocks map[string]interface{}) {

}
