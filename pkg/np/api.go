package np

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

type api struct {
	gameNumber string
	key        string
}

func newAPI(gameNumber, apiKey string) api {
	return api{
		gameNumber: gameNumber,
		key:        apiKey,
	}
}

func (a api) GetCurrentSnapshot() ([]byte, error) {
	client := &http.Client{}

	formData := url.Values{}
	formData.Set("game_number", a.gameNumber)
	formData.Set("code", a.key)
	formData.Set("api_version", "0.1")

	req, err := http.NewRequest("POST", apiBaseURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		slog.Error("Error creating request:", "err", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request", "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response:", "err", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, ErrUnexpectedStatusCode
	}

	var errorResponse APIErrorResponse
	if err = json.Unmarshal(body, &errorResponse); err != nil {
		slog.Error("Error decoding error response:", "err", err)
		return nil, err
	} else if errorResponse.ErrorMessage != "" {
		return nil, errorResponse
	}

	return body, nil
}
