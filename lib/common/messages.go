package common

type ToMessage struct {
	Channel  string `json:"channel"`
	Hostmask string `json:"hostmask"`
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
}

type FromMessage struct {
	Channel  string `json:"channel"`
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
}
