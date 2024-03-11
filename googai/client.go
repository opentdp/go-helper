package googai

import (
	"encoding/json"

	"github.com/opentdp/go-helper/request"
)

type Client struct {
	ApiBaseUrl     string
	ApiVersion     string
	ApiKey         string
	Model          string
	SafetySettings []SafetySetting
}

func NewClient(key string) *Client {

	return &Client{
		ApiBaseUrl: ApiBaseUrl,
		ApiVersion: ApiVersion,
		ApiKey:     key,
		Model:      "gemini-pro",
	}

}

func (c *Client) CreateChatCompletion(contents []Content) (*ResponseBody, error) {

	rq := &RequestBody{
		Contents:       contents,
		SafetySettings: c.SafetySettings,
	}
	body, _ := json.Marshal(rq)

	heaner := request.H{
		"Content-Type":   "application/json",
		"x-goog-api-key": c.ApiKey,
	}

	url := c.ApiBaseUrl + "/" + c.ApiVersion + "/models/" + c.Model + ":generateContent"
	response, err := request.Post(url, string(body), heaner)
	if err != nil {
		return nil, err
	}

	var resp ResponseBody
	err = json.Unmarshal(response, &resp)

	return &resp, err

}

func (c *Client) CreateImageCompletion(contents []Content) (*ResponseBody, error) {

	rq := &RequestBody{
		Contents:       contents,
		SafetySettings: c.SafetySettings,
	}
	body, _ := json.Marshal(rq)

	heaner := request.H{
		"Content-Type":   "application/json",
		"x-goog-api-key": c.ApiKey,
	}

	url := c.ApiBaseUrl + "/" + c.ApiVersion + "/models/gemini-pro-vision:generateContent"
	response, err := request.Post(url, string(body), heaner)
	if err != nil {
		return nil, err
	}

	var resp ResponseBody
	err = json.Unmarshal(response, &resp)

	return &resp, err

}
