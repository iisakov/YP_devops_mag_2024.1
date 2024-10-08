package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	r, err := http.Get("srv.msk01.gigacorp.local/_stats")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r)
}
