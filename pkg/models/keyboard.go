package models

type Keyboard struct {
	OneTime bool       `json:"one_time"`
	Inline  bool       `json:"inline"`
	Buttons [][]Button `json:"buttons"`
}

type Button struct {
	Action `json:"action"`
	color  string `json:"color"`
}

type Action struct {
	Type    string   `json:"type"`
	Label   string   `json:"label"`
	Link    string   `json:"link,omitempty"`
	Payload []string `json:"payload,omitempty"`
}