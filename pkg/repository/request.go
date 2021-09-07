package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"
)

type RequestApiRepository struct {
	config *config.Config
}

func NewRequestApiRepository(config *config.Config) *RequestApiRepository {
	return &RequestApiRepository{config: config}
}

func (rr *RequestApiRepository) SendRequest(q *url.Values, method string, expectedResult string) (interface{}, error) {
	var err error
	u := url.URL{
		Scheme: "https",
		Host:   rr.config.VKApiURL,
		Path:   "method/" + method,
	}
	q.Add("v", "5.131")
	q.Add("access_token", rr.config.Token)

	u.RawQuery, err = url.QueryUnescape(q.Encode())
	if err != nil {
		return nil, fmt.Errorf("error while unescape URL: %v", err)
	}
	cl := http.Client{}
	res, err := cl.Get(u.String())
	if err != nil {
		return "", fmt.Errorf("error while requesting: %v", err)
	}
	result, err := handleResponse(res, expectedResult)
	if err != nil {
		return "", fmt.Errorf("error while handling response: %v", err)
	}
	return result, nil
}

func handleResponse(res *http.Response, param string) (interface{}, error) {
	var result interface{}
	err := json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding response, %v", err)
	}
	if val, ok := result.(map[string]interface{})["error"]; ok {
		return "", fmt.Errorf("error while requesting: %s", val.(map[string]interface{})["error_msg"])
	}
	if val, ok := result.(map[string]interface{})["response"]; !ok {
		return "", fmt.Errorf("error while requesting: server not returned 'response' json")
	} else {
		if param != "" {
			return val.(map[string]interface{})[param], nil
		} else {
			return val, nil
		}
	}
}
