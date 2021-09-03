package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"
)

type CallbackSetupService struct {
	config *config.Config
}

func NewCallbackSetupService(config *config.Config) *CallbackSetupService {
	return &CallbackSetupService{
		config: config,
	}
}

func (cs *CallbackSetupService) SendConfirmationResponse(w http.ResponseWriter) error {
	code, err := cs.GetConfirmationCode()
	if err != nil {
		return fmt.Errorf("error while getting confirmation code: %v", err)
	}
	//Send reply to VK api
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(code))
	if err != nil {
		return fmt.Errorf("error while sending response: %v", err)
	}
	return nil
}

func (cs *CallbackSetupService) SetCallbackUrl() error {
	//https://api.vk.com/method/groups.addCallbackServer?group_id=...&url=...&title=...&secret_key=...&access_token=TOKEN&v=V
	var err error
	u := url.URL{
		Scheme: "https",
		Host:   cs.config.VKApiURL,
		Path:   "method/groups.addCallbackServer",
	}
	q := u.Query()
	q.Add("group_id", cs.config.VKGroupID)
	q.Add("url", "https://ae31-94-24-252-226.ngrok.io/callback")
	q.Add("secret_key", cs.config.VKCallbackSecret)
	q.Add("access_token", cs.config.Token)
	q.Add("v", "5.131")
	q.Add("title", "VKBot")
	u.RawQuery, err = url.QueryUnescape(q.Encode())

	cl := http.Client{}
	res, err := cl.Get(u.String())
	if err != nil {
		return fmt.Errorf("error while sending callback url configuration: %v", err)
	}
	result, err := cs.handleResponse(res, "server_id")
	if err != nil {
		return fmt.Errorf("error while handle response from from api on 'addCallbackServer': %v", err)
	}
	if err = cs.SetupCallbackService(fmt.Sprintf("%.0f", result.(float64))); err != nil {
		return fmt.Errorf("error while setting up callback service: %v", err)
	}
	return nil
}

func (cs *CallbackSetupService) SetupCallbackService(srvId string) error {

	//https://api.vk.com/method/groups.setCallbackSettings?group_id=...&url=...&title=...&secret_key=...&access_token=TOKEN&v=V
	u := url.URL{
		Scheme: "https",
		Host:   cs.config.VKApiURL,
		Path:   "method/groups.setCallbackSettings",
	}
	q := u.Query()
	q.Add("group_id", cs.config.VKGroupID)
	q.Add("access_token", cs.config.Token)
	q.Add("v", "5.131")
	q.Add("server_id", srvId)
	q.Add("api_version", "5.131")
	q.Add("message_allow", "1")
	q.Add("message_deny", "1")
	u.RawQuery = q.Encode()

	cl := http.Client{}
	res, err := cl.Get(u.String())
	if err != nil {
		return fmt.Errorf("error while sending configuration api callback: %v", err)
	}
	_, err = cs.handleResponse(res, "")
	if err != nil {
		return fmt.Errorf("error while handling response from api on 'groups.setCallbackSettings': %v", err)
	}
	return nil
}

func (cs *CallbackSetupService) GetConfirmationCode() (string, error) {
	//https://api.vk.com/method/groups.getCallbackConfirmationCode?group_id=...&v=V
	u := url.URL{
		Scheme: "https",
		Host:   cs.config.VKApiURL,
		Path:   "method/groups.getCallbackConfirmationCode",
	}
	q := u.Query()
	q.Add("group_id", cs.config.VKGroupID)
	q.Add("access_token", cs.config.Token)
	q.Add("v", "5.131")
	u.RawQuery = q.Encode()

	cl := http.Client{}
	res, err := cl.Get(u.String())
	if err != nil {
		return "", fmt.Errorf("error while requesting callback confirmation code: %v", err)
	}
	result, err := cs.handleResponse(res, "code")
	if err != nil {
		return "", fmt.Errorf("error while handle response from from api on 'getCallbackConfirmationCode': %v", err)
	}

	return result.(string), nil
}

func (cs *CallbackSetupService) handleResponse(res *http.Response, param string) (interface{}, error) {
	var result interface{}
	err := json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding response, %v", err)
	}
	if val, ok := result.(map[string]interface{})["error"]; ok {
		return "", fmt.Errorf("error while requesting confirmation code: %s", val.(map[string]interface{})["error_msg"])
	}
	if val, ok := result.(map[string]interface{})["response"]; !ok {
		return "", fmt.Errorf("error while requesting confirmation code: server not returned 'response' json")
	} else {
		if param != "" {
			return val.(map[string]interface{})[param], nil
		} else {
			return val, nil
		}
	}
}
