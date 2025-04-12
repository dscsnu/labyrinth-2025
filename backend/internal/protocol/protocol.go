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
