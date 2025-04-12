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

	for {

		packet := <-c.Recv
		if ok := handlePacket(packet); !ok {

			break

		}

	}

}

func handlePacket(packet protocol.Packet) {

	switch packet.Type {

	case protocol.PacketTypeChannelState:

	case protocol.PacketTypeBackground:

	case protocol.PacketTypeGame:

	}

}
