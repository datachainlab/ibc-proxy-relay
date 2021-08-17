package tendermint

import (
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	chantypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"

	"github.com/datachainlab/ibc-proxy-prover/pkg/proxy"
	"github.com/hyperledger-labs/yui-relayer/chains/tendermint"
)

var _ proxy.ProxyChainConfigI = (*ProxyChainConfig)(nil)

func (c *ProxyChainConfig) Build() (proxy.ProxyChainI, error) {
	chain, err := c.ChainConfig.Build()
	if err != nil {
		return nil, err
	}
	return NewTendermintProxyChain(chain.(*tendermint.Chain)), nil
}

type TendermintProxyChain struct {
	*tendermint.Chain
	proxyPath proxy.ProxyPath
}

var _ proxy.ProxyChainI = (*TendermintProxyChain)(nil)

func NewTendermintProxyChain(chain *tendermint.Chain) *TendermintProxyChain {
	return &TendermintProxyChain{Chain: chain}
}

func (c *TendermintProxyChain) SetProxyPath(path proxy.ProxyPath) {
	c.proxyPath = path
}

func (c *TendermintProxyChain) ProxyPath() proxy.ProxyPath {
	return c.proxyPath
}

func (c *TendermintProxyChain) QueryProxyClientState(height int64, upstreamClientID string) (*clienttypes.QueryClientStateResponse, error) {
	// c.Client.ABCIQuery()
	panic("not implemented error")
}

func (c *TendermintProxyChain) QueryProxyClientConsensusState(height int64, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
	panic("not implemented error")
}

func (c *TendermintProxyChain) QueryProxyConnectionState(height int64, upstreamClientID string) (*connectiontypes.QueryConnectionResponse, error) {
	panic("not implemented error")
}

func (c *TendermintProxyChain) QueryProxyChannel(height int64, upstreamClientID string) (chanRes *chantypes.QueryChannelResponse, err error) {
	panic("not implemented error")
}

func (c *TendermintProxyChain) QueryProxyPacketCommitment(height int64, seq uint64, upstreamClientID string) (comRes *chantypes.QueryPacketCommitmentResponse, err error) {
	panic("not implemented error")
}

func (c *TendermintProxyChain) QueryProxyPacketAcknowledgementCommitment(height int64, seq uint64, upstreamClientID string) (ackRes *chantypes.QueryPacketAcknowledgementResponse, err error) {
	panic("not implemented error")
}
