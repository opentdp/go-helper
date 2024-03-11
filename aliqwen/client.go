package aliqwen

import (
	"encoding/json"

	"github.com/opentdp/go-helper/request"
)

type Client struct {
	ApiBaseUrl string
	ApiVersion string
	ApiKey     string
	Model      string
	Params     *Parameters
}

func NewClient(key string) *Client {

	return &Client{
		ApiBaseUrl: ApiBaseUrl,
		ApiVersion: ApiVersion,
		ApiKey:     key,
		Model:      "qwen-max",
		Params:     &Parameters{EnableSearch: true},
	}

}

func (c *Client) CreateChatCompletion(messages []*Messages) (*ResponseBody, error) {

	req := RequestBody{
		Model:      c.Model,
		Input:      Input{messages},
		Parameters: c.Params,
	}
	body, _ := json.Marshal(req)

	heaner := request.H{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.ApiKey,
	}

	url := c.ApiBaseUrl + "/api/" + c.ApiVersion + "/services/aigc/text-generation/generation"
	response, err := request.Post(url, string(body), heaner)
	if err != nil {
		return nil, err
	}

	var resp ResponseBody
	err = json.Unmarshal(response, &resp)

	return &resp, err

}
