package zedmode

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"claude-proxy/bedrock"
)

type ZedMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ZedInput struct {
	Messages []ZedMessage `json:"messages"`
}

type ZedResponse struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func Run(stdin io.Reader, stdout, stderr io.Writer) {
	var input ZedInput
	if err := json.NewDecoder(bufio.NewReader(stdin)).Decode(&input); err != nil {
		fmt.Fprintf(stderr, "Failed to parse Zed input: %v\n", err)
		os.Exit(1)
	}

	var userText string
	for i := len(input.Messages) - 1; i >= 0; i-- {
		if input.Messages[i].Role == "user" {
			userText = input.Messages[i].Content
			break
		}
	}

	responseText, err := callClaudeAndParse(userText)
	if err != nil {
		fmt.Fprintf(stderr, "Failed to call Claude: %v\n", err)
		os.Exit(1)
	}

	resp := ZedResponse{
		Role:    "assistant",
		Content: responseText,
	}
	json.NewEncoder(stdout).Encode(resp)
}

func callClaudeAndParse(input string) (string, error) {
	resp, err := bedrock.CallClaude(context.Background(), input)
	if err != nil {
		return "", err
	}

	var parsed struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal([]byte(resp), &parsed); err != nil || len(parsed.Content) == 0 {
		return "", fmt.Errorf("unable to parse Claude response")
	}
	return parsed.Content[0].Text, nil
}
