package proxy

import (
	"fmt"
	"strings"

	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	"github.com/hyperledger-labs/yui-relayer/core"
)

var _ core.PathEndI = (*ProxyPathEnd)(nil)

func (p *ProxyPathEnd) Type() string {
	return "proxy"
}

func (p *ProxyPathEnd) ChainID() string {
	return p.ChainId
}

func (p *ProxyPathEnd) ClientID() string {
	return p.ClientId
}

func (p *ProxyPathEnd) ConnectionID() string {
	return p.ConnectionId
}

func (p *ProxyPathEnd) ChannelID() string {
	return p.ChannelId
}

func (p *ProxyPathEnd) PortID() string {
	return p.PortId
}

func (p *ProxyPathEnd) ChannelOrder() channeltypes.Order {
	return core.OrderFromString(strings.ToUpper(p.Order))
}

func (p *ProxyPathEnd) ChannelVersion() string {
	return p.Version
}

func (p *ProxyPathEnd) Validate() error {
	pe := &core.PathEnd{
		ChainId:      p.ChainId,
		ClientId:     p.ClientId,
		ConnectionId: p.ConnectionId,
		ChannelId:    p.ChannelId,
		PortId:       p.PortId,
		Order:        p.Order,
		Version:      p.Version,
	}
	return pe.Validate()
}

type clientPathEnd struct {
	ChainId  string
	ClientId string
}

var _ core.PathEndI = (*clientPathEnd)(nil)

func (p *clientPathEnd) Type() string {
	return "client"
}

func (p *clientPathEnd) ChainID() string {
	return p.ChainId
}

func (p *clientPathEnd) ClientID() string {
	return p.ClientId
}

func (p *clientPathEnd) ConnectionID() string {
	panic("not supported")
}

func (p *clientPathEnd) ChannelID() string {
	panic("not supported")
}

func (p *clientPathEnd) PortID() string {
	panic("not supported")
}

func (p *clientPathEnd) ChannelOrder() channeltypes.Order {
	panic("not supported")
}

func (p *clientPathEnd) ChannelVersion() string {
	panic("not supported")
}

func (p *clientPathEnd) Validate() error {
	if p.ChainId == "" || p.ClientId == "" {
		return fmt.Errorf("each field must be given non-empty string")
	}
	return nil
}

func (p *clientPathEnd) String() string {
	return fmt.Sprintf("chainID=%v clientID=%v", p.ChainId, p.ClientID())
}
