package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type BardSession struct {
	SessionId      string  `json:"session_id"`
	SNlM0e         string  `json:"SNlM0e"`
	ReqID          string  `json:"req_id"`
	ConversationID *string `json:"conversation_id"`
	ResponseID     *string `json:"response_id"`
	ChoiceID       *string `json:"choice_id"`
}

func InitBard(sessionId string) (*BardSession, error) {
	request, err := http.NewRequest("GET", "https://bard.google.com", nil)
	if err != nil {
		return nil, errors.New("failed to create request: " + err.Error())
	}
	request.Header.Set("Cookie", "__Secure-1PSID="+sessionId)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.New("failed to send request: " + err.Error())
	}
	defer response.Body.Close()

	// read body as string
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	// regex match for SNlM0e /SNlM0e":"(.*?)"/
	matches := regexp.MustCompile(`SNlM0e":"(.*?)"`).FindSubmatch(body)
	if len(matches) < 2 {
		return nil, errors.New("failed to find SNlM0e in response body")
	}

	return &BardSession{
		SNlM0e:    string(matches[1]),
		SessionId: sessionId,
		ReqID:     "0",
	}, nil
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

func (s BardSession) Ask(message string) (string, error) {
	request, err := http.NewRequest("POST", "https://bard.google.com/_/BardChatUi/data/assistant.lamda.BardFrontendService/StreamGenerate", nil)
	if err != nil {
		return "", errors.New("failed to create request: " + err.Error())
	}
	urlQuery := request.URL.Query()
	urlQuery.Add("bl", "boq_assistant-bard-web-server_20230730.19_p0")
	urlQuery.Add("_reqID", s.ReqID)
	urlQuery.Add("rt", "c")
	request.URL.RawQuery = urlQuery.Encode()
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Cookie", "__Secure-1PSID="+s.SessionId)

	freqRight := []interface{}{[]string{message}, nil, []*string{s.ConversationID, s.ResponseID, s.ChoiceID}}
	freqRightBytes, err := json.Marshal(freqRight)
	if err != nil {
		return "", errors.New("failed to marshal freq right: " + err.Error())
	}
	freqRightString := string(freqRightBytes)
	freq := []interface{}{nil, freqRightString}
	freqBytes, err := json.Marshal(freq)
	if err != nil {
		return "", errors.New("failed to marshal freq: " + err.Error())
	}
	freqString := string(freqBytes)
	formData := url.QueryEscape("f.req") + "=" + url.QueryEscape(freqString) + "&" + url.QueryEscape("at") + "=" + url.QueryEscape(s.SNlM0e)
	request.Body = io.NopCloser(strings.NewReader(formData))

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
	bodyStringParts := strings.Split(string(body), "\n")
	if len(bodyStringParts) < 3 {
		return "", errors.New("failed to parse response body")
	}

	return bodyStringParts[3], nil
}
