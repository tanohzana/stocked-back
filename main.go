package main

import (
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	sioServer, err := socketio.NewServer(nil)
	serverPort := ":8080"

	if err != nil {
		log.Fatal(err)
	}

	go sioServer.Serve()
	defer sioServer.Close()

	http.Handle("/socketio/", sioServer)
	http.Handle("/", http.FileServer(http.Dir("./public")))

	log.Println("Magic is happening on port ", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}
