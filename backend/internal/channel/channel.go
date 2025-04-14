package channel

import (
	"labyrinth/internal/protocol"
	"sync"
)

type ChannelPool struct {
	mut   sync.Mutex
	cpool map[string]*Channel
}

func NewChannelPool() *ChannelPool {

	return &ChannelPool{cpool: map[string]*Channel{}}

}

func (c *ChannelPool) AddChannel(teamId string, channel *Channel) {

	c.mut.Lock()
	c.cpool[teamId] = channel
	c.mut.Unlock()
}

func (c *ChannelPool) DeleteChannel(teamId string) {

	c.mut.Lock()
	delete(c.cpool, teamId)
	c.mut.Unlock()
}

func (c *ChannelPool) GetChannel(teamId string) *Channel {

	c.mut.Lock()
	defer c.mut.Unlock()

	return c.cpool[teamId]
}

type Channel struct {
	Recv             chan protocol.Packet
	BroadcastClients struct {
		bc  map[chan<- protocol.Packet]struct{}
		mut sync.Mutex
	}
}

func NewChannel() *Channel {

	return &Channel{Recv: make(chan protocol.Packet), BroadcastClients: struct {
		bc  map[chan<- protocol.Packet]struct{}
		mut sync.Mutex
	}{

		bc: map[chan<- protocol.Packet]struct{}{},
	}}

}

func (c *Channel) AddMember(memberChannel chan<- protocol.Packet) {

	c.BroadcastClients.mut.Lock()
	c.BroadcastClients.bc[memberChannel] = struct{}{}
	c.BroadcastClients.mut.Unlock()

}

func (c *Channel) Broadcast(packet protocol.Packet) {

	c.Recv <- packet

}

func (c *Channel) Start() {

	for packet := range c.Recv {

		if ok := c.handlePacket(packet); !ok {

			break

		}

		c.BroadcastClients.mut.Lock()
		for client := range c.BroadcastClients.bc {

			client <- packet
		}
		c.BroadcastClients.mut.Unlock()

	}

}

// change backend state with packets
func (c *Channel) handlePacket(packet protocol.Packet) bool {

	switch packet.Type {

	case "BackgroundMessage":

	case "ChannelStateMessage":

	case "GameMessage":

	}

	return true
}
