package httpmode

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"claude-proxy/bedrock"
)

type InputPayload struct {
	InputText string `json:"input"`
}

func Start() {
	http.HandleFunc("/invoke", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var payload InputPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	text, err := callClaudeAndParse(payload.InputText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"response": text,
	})
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
		return "", err
	}
	return parsed.Content[0].Text, nil
}
