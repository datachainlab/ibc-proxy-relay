package tendermint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	chantypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"
	"github.com/datachainlab/ibc-proxy-prover/pkg/proxy"
	"github.com/hyperledger-labs/yui-relayer/chains/tendermint"
	"github.com/hyperledger-labs/yui-relayer/core"
)

var _ proxy.ProxyChainProverConfigI = (*ProxyChainProverConfig)(nil)

func (c *ProxyChainProverConfig) Build(proxyChain proxy.ProxyChainI) (proxy.ProxyChainProverI, error) {
	return NewProxyChainProver(c, proxyChain), nil
}

type ProxyChainProver struct {
	proxyChain proxy.ProxyChainI
	*tendermint.Prover
}

var _ proxy.ProxyChainProverI = (*ProxyChainProver)(nil)

func NewProxyChainProver(cfg *ProxyChainProverConfig, proxyChain proxy.ProxyChainI) *ProxyChainProver {
	chain := proxyChain.(*TendermintProxyChain)

	return &ProxyChainProver{
		proxyChain: proxyChain,
		Prover:     tendermint.NewProver(chain.Chain, *cfg.ProverConfig),
	}
}

// CreateMsgCreateClient creates a CreateClientMsg to this chain
func (p *ProxyChainProver) CreateMsgCreateClient(clientID string, dstHeader core.HeaderI, signer sdk.AccAddress) (*clienttypes.MsgCreateClient, error) {
	// TODO returns a msg that corresponding to ProxyClient https://github.com/datachainlab/ibc-proxy/blob/945de5a19447cf9369a8b6d38fa8b29d710849aa/modules/light-clients/xx-proxy/types/proxy.pb.go#L29
	panic("not implemented error")
}

// TODO other lightclient's methods should be also implemented

func (p *ProxyChainProver) QueryProxyConnectionStateWithProof(height int64, upstreamClientID string) (*connectiontypes.QueryConnectionResponse, error) {
	panic("not implemented error")
}

func (p *ProxyChainProver) QueryProxyClientStateWithProof(height int64, upstreamClientID string) (*clienttypes.QueryClientStateResponse, error) {
	panic("not implemented error")
}

func (p *ProxyChainProver) QueryProxyClientConsensusStateWithProof(height int64, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
	panic("not implemented error")
}

func (p *ProxyChainProver) QueryProxyChannelWithProof(height int64, upstreamClientID string) (chanRes *chantypes.QueryChannelResponse, err error) {
	panic("not implemented error")
}

func (p *ProxyChainProver) QueryProxyPacketCommitmentWithProof(height int64, seq uint64, upstreamClientID string) (comRes *chantypes.QueryPacketCommitmentResponse, err error) {
	panic("not implemented error")
}

func (p *ProxyChainProver) QueryProxyPacketAcknowledgementCommitmentWithProof(height int64, seq uint64, upstreamClientID string) (ackRes *chantypes.QueryPacketAcknowledgementResponse, err error) {
	panic("not implemented error")
}
