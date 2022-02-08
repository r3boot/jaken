package common

type RawMessage struct {
	Channel  string `json:"channel"`
	Hostmask string `json:"hostmask"`
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
}

type CommandMessage struct {
	Channel   string `json:"channel"`
	Hostmask  string `json:"hostmask"`
	Nickname  string `json:"nickname"`
	Command   string `json:"command"`
	Arguments string `json:"arguments"`
}

type FromMessage struct {
	Recipient string `json:"channel"`
	Message   string `json:"message"`
}

type TopicMessage struct {
	Channel string `json:"channel"`
	Topic   string `json:"message"`
}
