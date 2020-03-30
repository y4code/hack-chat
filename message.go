package main

//{"cmd":"join","channel":"your-channel","nick":"woshimingzi"}
// TODO 需要优化
type Message struct {
	Cmd     string `json:"cmd"`
	Channel string `json:"channel"`
	Nick    string `json:"nick"`
	Text    string `json:"text"`
}

type JoinMessage struct {
	Cmd     string `json:"cmd"`
	Channel string `json:"channel"`
	Nick    string `json:"nick"`
}

type ChatMessage struct {
	Cmd  string `json:"cmd"`
	Nick string `json:"nick"`
	Text string `json:"text"`
}
