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

}

type ChannelStateMessage struct{}

func DecodeChannelStateMessage(msg []byte) (ChannelStateMessage, error) {

}

type GameMessage struct{}

func DecodeGameMessage(msg []byte) (GameMessage, error) {

}
