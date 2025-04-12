package channel

import "labyrinth/internal/protocol"

type Channel struct {
	Recv             <-chan protocol.Packet
	BroadcastClients map[chan<- protocol.Packet]struct{}
}

func NewChannel() *Channel {

	return &Channel{Recv: make(<-chan protocol.Packet), BroadcastClients: map[chan<- protocol.Packet]struct{}{}}

}

func (c *Channel) Start() {

	for packet := range c.Recv {

		if ok := c.handlePacket(packet); !ok {

			break

		}

	}

}

func (c *Channel) handlePacket(packet protocol.Packet) bool {

	switch packet.Type {

	case protocol.PacketTypeChannelState:

		//channelStateMessage, err := protocol.DecodeChannelStateMessage(packet.Message)

	case protocol.PacketTypeBackground:
		//backgroundMessage, err := protocol.DecodeBackgroundMessage(packet.Message)

	case protocol.PacketTypeGame:
		//gameMessage, err := protocol.DecodeGameMessage(packet.Message)

	}

	return true
}
