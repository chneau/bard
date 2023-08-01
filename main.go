package main

import "log"

func main() {
	sessionId := PromptSessionId()
	bardCli, err := InitBard(sessionId)
	if err != nil {
		panic("failed to initialize bard: " + err.Error())
	}
	log.Println("bardCli: ", bardCli)
}
