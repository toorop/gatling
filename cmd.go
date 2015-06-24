package main

import (
	"log"
)

func main() {

	concurrency := 200

	g := gatling{
		Target:   "127.0.0.01:2525",
		MailFrom: "gatling@tmail.io",
		RcptTo:   "gatling@pepper.pm",
		BodySize: 1000,
		Bullets:  100,
	}

	cDone := make(chan bool)
	for i := 0; i < concurrency; i++ {
		go g.Fire(&cDone)
	}

	for i := 0; i < concurrency; i++ {
		select {
		case <-cDone:
			log.Printf("gatling %d/%d done", i, concurrency)
		}
	}

}
