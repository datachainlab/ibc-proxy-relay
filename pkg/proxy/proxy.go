package proxy

import (
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

	SetProxyPath(ProxyPath)
	ProxyPath() ProxyPath
}

type ProxyPath struct {
	UpstreamClientID string
	UpstreamChain    core.ChainI
}

type ProxyChainQueryierI interface {
	QueryProxyClientState(height int64) (*clienttypes.QueryClientStateResponse, error)
	QueryProxyClientConsensusState(height int64, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error)
	QueryProxyConnectionState(height int64) (*connectiontypes.QueryConnectionResponse, error)
	QueryProxyChannel(height int64) (chanRes *chantypes.QueryChannelResponse, err error)
	QueryProxyPacketCommitment(height int64, seq uint64) (comRes *chantypes.QueryPacketCommitmentResponse, err error)
	QueryProxyPacketAcknowledgementCommitment(height int64, seq uint64) (ackRes *chantypes.QueryPacketAcknowledgementResponse, err error)
}

type ProxyChainProverI interface {
	core.ProverI
	ProxyChainProverQueryierI
}

type ProxyChainProverQueryierI interface {
	QueryProxyClientStateWithProof(height int64) (*clienttypes.QueryClientStateResponse, error)
	QueryProxyClientConsensusStateWithProof(height int64, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error)
	QueryProxyConnectionStateWithProof(height int64) (*connectiontypes.QueryConnectionResponse, error)
	QueryProxyChannelWithProof(height int64) (chanRes *chantypes.QueryChannelResponse, err error)
	QueryProxyPacketCommitmentWithProof(height int64, seq uint64) (comRes *chantypes.QueryPacketCommitmentResponse, err error)
	QueryProxyPacketAcknowledgementCommitmentWithProof(height int64, seq uint64) (ackRes *chantypes.QueryPacketAcknowledgementResponse, err error)
}
