package proxy

import (
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hyperledger-labs/yui-relayer/core"
)

type ProxiableChain struct {
	core.ChainI
	eventListener ProxyEventListener
}

var _ ProxiableChainI = (*ProxiableChain)(nil)

func NewProxiableChain(chain core.ChainI) *ProxiableChain {
	return &ProxiableChain{
		ChainI: chain,
	}
}

func (c *ProxiableChain) RegisterEventListener(listner ProxyEventListener) {
	c.eventListener = listner
}

func (c *ProxiableChain) SendMsgs(msgs []sdk.Msg) ([]byte, error) {
	res, err := c.ChainI.SendMsgs(msgs)
	if err != nil {
		return nil, err
	}
	return res, c.eventListener.OnSentMsg(c.Path(), msgs)
}

func (c *ProxiableChain) Send(msgs []sdk.Msg) bool {
	_, err := c.SendMsgs(msgs)
	if err != nil {
		log.Println(err)
	}
	return err == nil
}
