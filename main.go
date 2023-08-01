package main

import "log"

func main() {
	sessionId := PromptSessionId()
	bardCli, err := InitBard(sessionId)
	if err != nil {
		panic("failed to initialize bard: " + err.Error())
	}
	log.Println("bardCli: ", bardCli)

	str, err := bardCli.Ask("What time is it?")
	if err != nil {
		panic("failed to ask question: " + err.Error())
	}
	log.Println("str: ", str)
}
