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

}
