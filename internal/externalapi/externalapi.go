package externalapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type APIResponse struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func GetDataFromExternalAPI(apiBaseURL, group, song string) (APIResponse, error) {
	apiURL := fmt.Sprintf("%s/info?group=%s&song=%s", apiBaseURL, url.QueryEscape(group), url.QueryEscape(song))

	resp, err := http.Get(apiURL)
	if err != nil {
		return APIResponse{}, fmt.Errorf("failed to call external API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return APIResponse{}, fmt.Errorf("external API returned non-OK status: %s", resp.Status)
	}

	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return APIResponse{}, fmt.Errorf("failed to decode API response: %w", err)
	}

	return apiResponse, nil
}
