package main

import (
	"log"
	"time"

	"github.com/preslavrachev/req"
)

type ReqBinResponse struct {
	Success string `json:"success,omitempty"`
}

func main() {
	res, err := req.Get[ReqBinResponse]("https://reqbin.com/echo/get/json", req.WithTimeout(1*time.Millisecond))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", res)
}
