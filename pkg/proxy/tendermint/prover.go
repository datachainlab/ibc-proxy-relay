package tendermint

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"
	"github.com/datachainlab/ibc-proxy-relay/pkg/proxy"
	proxytypes "github.com/datachainlab/ibc-proxy/modules/light-clients/xx-proxy/types"
	ibcproxytypes "github.com/datachainlab/ibc-proxy/modules/proxy/types"
	"github.com/hyperledger-labs/yui-relayer/chains/tendermint"
	"github.com/hyperledger-labs/yui-relayer/core"
)

var _ proxy.ProxyChainProverConfigI = (*ProxyChainProverConfig)(nil)

func (c *ProxyChainProverConfig) Build(proxyChain proxy.ProxyChainI) (proxy.ProxyChainProverI, error) {
	return NewProxyChainProver(c, proxyChain), nil
}

type ProxyChainProver struct {
	proxyChain *TendermintProxyChain
	*tendermint.Prover
}

var _ proxy.ProxyChainProverI = (*ProxyChainProver)(nil)

func NewProxyChainProver(cfg *ProxyChainProverConfig, proxyChain proxy.ProxyChainI) *ProxyChainProver {
	chain := proxyChain.(*TendermintProxyChain)

	return &ProxyChainProver{
		proxyChain: chain,
		Prover:     tendermint.NewProver(chain.Chain, *cfg.ProverConfig),
	}
}

func (p *ProxyChainProver) Codec() codec.ProtoCodecMarshaler {
	return p.proxyChain.Codec()
}

// CreateMsgCreateClient creates a CreateClientMsg to this chain
func (p *ProxyChainProver) CreateMsgCreateClient(clientID string, dstHeader core.HeaderI, signer sdk.AccAddress) (*clienttypes.MsgCreateClient, error) {
	msg, err := p.Prover.CreateMsgCreateClient(clientID, dstHeader, signer)
	if err != nil {
		return nil, err
	}

	ibcPrefix := commitmenttypes.NewMerklePrefix([]byte(host.StoreKey))
	proxyPrefix := commitmenttypes.NewMerklePrefix([]byte(ibcproxytypes.StoreKey))
	clientState := &proxytypes.ClientState{
		UpstreamClientId: p.proxyChain.upstreamPathEnd.UpstreamClientId,
		ProxyClientState: msg.ClientState,
		IbcPrefix:        &ibcPrefix,
		ProxyPrefix:      &proxyPrefix,
		TrustedSetup:     true,
	}
	anyClientState, err := clienttypes.PackClientState(clientState)
	if err != nil {
		return nil, err
	}
	consensusState := proxytypes.NewConsensusState(msg.ConsensusState)
	anyConsensusState, err := clienttypes.PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}
	return &clienttypes.MsgCreateClient{
		ClientState:    anyClientState,
		ConsensusState: anyConsensusState,
		Signer:         msg.Signer,
	}, nil
}

func (p *ProxyChainProver) QueryProxyClientStateWithProof(height int64) (*clienttypes.QueryClientStateResponse, error) {
	return p.proxyChain.ABCIQueryProxyClientState(height, p.proxyChain.upstreamPathEnd.ClientId, true)
}

func (p *ProxyChainProver) QueryProxyClientConsensusStateWithProof(height int64, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
	return p.proxyChain.ABCIQueryProxyConsensusState(height, p.proxyChain.upstreamPathEnd.ClientId, dstClientConsHeight, true)
}

func (p *ProxyChainProver) QueryProxyConnectionStateWithProof(height int64) (*connectiontypes.QueryConnectionResponse, error) {
	return p.proxyChain.ABCIQueryProxyConnection(height, p.proxyChain.upstreamPathEnd.ConnectionId, true)
}

func (p *ProxyChainProver) QueryProxyChannelWithProof(height int64) (chanRes *channeltypes.QueryChannelResponse, err error) {
	return p.proxyChain.ABCIQueryProxyChannel(height, p.proxyChain.upstreamPathEnd.PortId, p.proxyChain.upstreamPathEnd.ChannelId, true)
}

func (p *ProxyChainProver) QueryProxyPacketCommitmentWithProof(height int64, seq uint64) (comRes *channeltypes.QueryPacketCommitmentResponse, err error) {
	return p.proxyChain.ABCIQueryProxyPacketCommitment(height, seq, true)
}

func (p *ProxyChainProver) QueryProxyPacketAcknowledgementCommitmentWithProof(height int64, seq uint64) (ackRes *channeltypes.QueryPacketAcknowledgementResponse, err error) {
	return p.proxyChain.ABCIQueryProxyPacketAcknowledgementCommitment(height, seq, true)
}
