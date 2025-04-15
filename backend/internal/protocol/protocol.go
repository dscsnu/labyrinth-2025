package protocol

type Packet struct {
	Type                string `json:"type"`
	BackgroundMessage   `json:"backgroundMessage"`
	ChannelStateMessage `json:"channelStateMessage"`
	GameMessage         `json:"gameMessage"`
}

type TeamPacket struct {
	Relay      string `json:"message"`
	MsgContext string `json:"msgcontext"`
}

type BackgroundMessage struct {
	Relay      string `json:"message"`
	MsgContext string `json:"msgcontext"`
}
type ChannelStateMessage struct {
	Relay      string `json:"message"`
	MsgContext string `json:"msgcontext"`
}
type GameMessage struct {
	Relay      string `json:"message"`
	MsgContext string `json:"msgcontext"`
}
