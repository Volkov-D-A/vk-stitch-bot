package repository

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"

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

	u.RawQuery = q.Encode()
	fmt.Println(u.String())
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

func (rr *RequestApiRepository) SendMessage(text string, keyboard interface{}, mr *models.MessageRecipient) error {
	val := url.Values{}
	val.Add("message", text)
	if keyboard != nil {
		js, err := json.Marshal(keyboard.(*models.Keyboard))
		if err != nil {
			return fmt.Errorf("error marshalling keyboard: %v", err)
		}
		val.Add("keyboard", string(js))
	}
	val.Add("peer_id", strconv.Itoa(mr.Id))
	val.Add("random_id", getRandomId())
	val.Add("dont_parse_links", "1")
	_, err := rr.SendRequest(&val, "messages.send", "")
	if err != nil {
		return fmt.Errorf("error while sending message: %v", err)
	}
	return nil
}

func (rr *RequestApiRepository) GetCallbackServerInfo() ([]models.CallbackServerItem, error) {
	var srv []models.CallbackServerItem
	val := url.Values{}
	val.Add("group_id", rr.config.VKGroupID)
	result, err := rr.SendRequest(&val, "groups.getCallbackServers", "items")
	if err != nil {
		return nil, fmt.Errorf("error while sending 'getCallbackServers' request: %v", err)
	}
	js, _ := json.Marshal(result)
	if err = json.Unmarshal(js, &srv); err != nil {
		return nil, fmt.Errorf("error while unmarshaling 'getCallbackServers' request: %v", err)
	}
	return srv, nil
}

func (rr *RequestApiRepository) SetCallbackUrl() (string, error) {
	val := url.Values{}
	val.Add("url", rr.config.CallbackUrl)
	val.Add("title", "VKBot")
	val.Add("group_id", rr.config.VKGroupID)
	val.Add("secret_key", rr.config.VKCallbackSecret)
	result, err := rr.SendRequest(&val, "groups.addCallbackServer", "server_id")
	if err != nil {
		return "", fmt.Errorf("error on sending 'addCallbackServer' request: %v", err)
	}
	return fmt.Sprintf("%.0f", result.(float64)), nil
}

func (rr *RequestApiRepository) SetupCallbackService(srvId string) error {
	val := url.Values{}
	val.Add("group_id", rr.config.VKGroupID)
	val.Add("server_id", srvId)
	val.Add("api_version", "5.131")
	val.Add("message_allow", "1")
	val.Add("message_deny", "1")
	val.Add("message_new", "1")
	_, err := rr.SendRequest(&val, "groups.setCallbackSettings", "")
	if err != nil {
		return fmt.Errorf("error on sending 'setCallbackSettings' request: %v", err)
	}
	return nil
}

func (rr *RequestApiRepository) GetConfirmationCode() (string, error) {
	val := url.Values{}
	val.Add("group_id", rr.config.VKGroupID)
	result, err := rr.SendRequest(&val, "groups.getCallbackConfirmationCode", "code")
	if err != nil {
		return "", fmt.Errorf("error on sending 'getCallbackConfirmationCode' request: %v", err)
	}
	return result.(string), nil
}

func (rr *RequestApiRepository) RemoveCallbackServer(id string) error {
	val := url.Values{}
	val.Add("group_id", rr.config.VKGroupID)
	val.Add("server_id", id)
	_, err := rr.SendRequest(&val, "groups.deleteCallbackServer", "")
	if err != nil {
		return fmt.Errorf("error while sending 'deleteCallbackServer' request: %v", err)
	}
	return nil
}

func handleResponse(res *http.Response, param string) (interface{}, error) {
	var result interface{}
	fmt.Println(res.Header)
	err := json.NewDecoder(res.Body).Decode(&result)
	fmt.Println(result)
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

func getRandomId() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Intn(1000000000))

}
