package proxy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	chantypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"
	"github.com/hyperledger-labs/yui-relayer/core"
)

type ProxyProvableChain struct {
	ProxyChainI
	ProxyChainProverI
}

func NewProxyProvableChain(chain ProxyChainI, prover ProxyChainProverI) *ProxyProvableChain {
	return &ProxyProvableChain{ProxyChainI: chain, ProxyChainProverI: prover}
}

type ProxyChainI interface {
	core.ChainI
	ProxyChainQueryierI
}

type ProxyChainQueryierI interface {
	QueryProxyClientState(height int64, upstreamClientID string) (*clienttypes.QueryClientStateResponse, error)
	QueryProxyClientConsensusState(height int64, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error)
	QueryProxyConnectionState(height int64, upstreamClientID string) (*connectiontypes.QueryConnectionResponse, error)
	QueryProxyChannel(height int64, upstreamClientID string) (chanRes *chantypes.QueryChannelResponse, err error)
	QueryProxyPacketCommitment(height int64, seq uint64, upstreamClientID string) (comRes *chantypes.QueryPacketCommitmentResponse, err error)
	QueryProxyPacketAcknowledgementCommitment(height int64, seq uint64, upstreamClientID string) (ackRes *chantypes.QueryPacketAcknowledgementResponse, err error)
}

type ProxyChainProverI interface {
	core.ProverI
	ProxyChainProverQueryierI
}

type ProxyChainProverQueryierI interface {
	QueryProxyClientStateWithProof(height int64, upstreamClientID string) (*clienttypes.QueryClientStateResponse, error)
	QueryProxyClientConsensusStateWithProof(height int64, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error)
	QueryProxyConnectionStateWithProof(height int64, upstreamClientID string) (*connectiontypes.QueryConnectionResponse, error)
	QueryProxyChannelWithProof(height int64, upstreamClientID string) (chanRes *chantypes.QueryChannelResponse, err error)
	QueryProxyPacketCommitmentWithProof(height int64, seq uint64, upstreamClientID string) (comRes *chantypes.QueryPacketCommitmentResponse, err error)
	QueryProxyPacketAcknowledgementCommitmentWithProof(height int64, seq uint64, upstreamClientID string) (ackRes *chantypes.QueryPacketAcknowledgementResponse, err error)
}

type ProxyEventListener interface {
	OnSentMsg(path *core.PathEnd, msgs []sdk.Msg) error
}

// TODO this interface should be merged into core.ChainI
type ProxiableChainI interface {
	core.ChainI
	RegisterEventListener(ProxyEventListener)
}
