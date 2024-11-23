package main

import (
	"code_challenge1/log"
	"code_challenge1/server"
)

func main() {

	s, err := server.NewServer()
	if err != nil {
		panic(err)
	}
	if err := s.Serve(":8080"); err != nil {
		log.Errorf("serve server error: %v", err)
	}
}
