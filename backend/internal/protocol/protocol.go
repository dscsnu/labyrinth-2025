package protocol

import (
	"encoding/binary"
)

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

type BackgroundMessageContext int

const (
	JoinBackgroundMessageContext BackgroundMessageContext = iota
	LeaveBackgroundMessageContext
)

type BackgroundMessage struct {
	Message    [128]byte
	MsgContext BackgroundMessageContext
}

func DecodeBackgroundMessage(msg []byte) (BackgroundMessage, error) {

	backgroundMessage := BackgroundMessage{}
	_, err := binary.Decode(msg, binary.LittleEndian, &backgroundMessage)

	return backgroundMessage, err
}

type ChannelStateMessage struct {
	Open bool
}

func DecodeChannelStateMessage(msg []byte) (ChannelStateMessage, error) {

	channelStateMessage := ChannelStateMessage{}
	_, err := binary.Decode(msg, binary.LittleEndian, &channelStateMessage)

	return channelStateMessage, err
}

type GameMessage struct {

	// game message fields

}

func DecodeGameMessage(msg []byte) (GameMessage, error) {
	gameMessage := GameMessage{}
	_, err := binary.Decode(msg, binary.LittleEndian, &gameMessage)
	return gameMessage, err
}
