package main

import (
	"log"
	"net/http"

	"github.com/throwpedro/ground_team_training_scanner/routes"
)

func main() {
	http.HandleFunc("/", routes.GetGroundTimes)

	println("Listening on port 8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
