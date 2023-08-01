package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type SNlM0e string

func InitBard(sessionId string) (SNlM0e, error) {
	request, err := http.NewRequest("GET", "https://bard.google.com", nil)
	if err != nil {
		return "", errors.New("failed to create request: " + err.Error())
	}
	request.Header.Set("Cookie", "__Secure-1PSID="+sessionId)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", errors.New("failed to send request: " + err.Error())
	}
	defer response.Body.Close()

	// read body as string
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("failed to read response body: " + err.Error())
	}

	// regex match for SNlM0e /SNlM0e":"(.*?)"/
	matches := regexp.MustCompile(`SNlM0e":"(.*?)"`).FindSubmatch(body)
	if len(matches) < 2 {
		return "", errors.New("failed to find SNlM0e in response body")
	}

	return SNlM0e(matches[1]), nil
}

func ReadSavedSessionId() string {
	bytes, err := os.ReadFile(os.Getenv("HOME") + "/.config/bard/session_id")
	if err != nil {
		return ""
	}
	str := string(bytes)
	return strings.Fields(str)[0]
}

func WriteSessionId(sessionId string) error {
	fmt.Println("Writing session ID to ~/.config/bard/session_id")
	_ = os.MkdirAll(os.Getenv("HOME")+"/.config/bard", 0755)
	err := os.WriteFile(os.Getenv("HOME")+"/.config/bard/session_id", []byte(sessionId), 0644)
	if err != nil {
		return errors.New("failed to write session ID: " + err.Error())
	}
	return nil
}

func DeleteSessionId() error {
	fmt.Println("Deleting session ID from ~/.config/bard/session_id")
	err := os.Remove("~/.config/bard/session_id")
	if err != nil {
		return errors.New("failed to delete session ID: " + err.Error())
	}
	return nil
}

func PromptSessionId() string {
	sessionId := ReadSavedSessionId()
	if sessionId != "" {
		fmt.Println("Using saved session ID")
		return sessionId
	}
	fmt.Println("Go to https://bard.google.com and log in")
	fmt.Println("Open the developer console and look for the cookie named __Secure-1PSID")
	fmt.Println("Enter your session ID: ")
	_ = os.Stdout.Sync()
	_, err := fmt.Scanln(&sessionId)
	if err != nil {
		panic("failed to read session ID: " + err.Error())
	}
	err = WriteSessionId(sessionId)
	if err != nil {
		panic("failed to write session ID: " + err.Error())
	}
	return sessionId
}
