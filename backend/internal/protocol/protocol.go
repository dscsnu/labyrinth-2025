package protocol

type PacketType int

const (
	PacketTypeBackground PacketType = iota
	PacketTypeGame
	PacketTypeChannelState
)

type Packet struct {
	Type    PacketType
	Message []byte
}

type BackgroundMessage struct{}

func DecodeBackgroundMessage(msg []byte) (BackgroundMessage, error) {
	return BackgroundMessage{}, nil
}

type ChannelStateMessage struct{}

func DecodeChannelStateMessage(msg []byte) (ChannelStateMessage, error) {
	return ChannelStateMessage{}, nil
}

type GameMessage struct{}

func DecodeGameMessage(msg []byte) (GameMessage, error) {
	return GameMessage{}, nil
}
