package main

import (
	"log"
	"net/http"
	"yoon/salatmv/lib"
)

func main() {
	log.Fatal(http.ListenAndServe(":3407", lib.NewServer()))
}
