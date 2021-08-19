package tendermint

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"
	"github.com/datachainlab/ibc-proxy-prover/pkg/proxy"
	proxytypes "github.com/datachainlab/ibc-proxy/modules/light-clients/xx-proxy/types"
	ibcproxytypes "github.com/datachainlab/ibc-proxy/modules/proxy/types"
	"github.com/hyperledger-labs/yui-relayer/chains/tendermint"
	"github.com/hyperledger-labs/yui-relayer/core"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/rpc/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (p *ProxyChainProver) upstreamPathEnd() *core.PathEnd {
	return p.proxyChain.ProxyPath().UpstreamChain.Path()
}

func (p *ProxyChainProver) upstreamClientID() string {
	return p.proxyChain.ProxyPath().UpstreamClientID
}

func (p *ProxyChainProver) upstreamPrefix() *commitmenttypes.MerklePrefix {
	prefix := commitmenttypes.NewMerklePrefix([]byte(host.StoreKey))
	return &prefix
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
		UpstreamClientId: p.upstreamClientID(),
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

// TODO other lightclient's methods should be also implemented

func (p *ProxyChainProver) QueryProxyClientStateWithProof(height int64) (*clienttypes.QueryClientStateResponse, error) {
	return p.queryProxyClientState(height, p.upstreamPathEnd().ClientID)
}

func (p *ProxyChainProver) QueryProxyClientConsensusStateWithProof(height int64, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
	return p.queryProxyClientConsensusState(height, p.upstreamPathEnd().ClientID, dstClientConsHeight)
}

func (p *ProxyChainProver) QueryProxyConnectionStateWithProof(height int64) (*connectiontypes.QueryConnectionResponse, error) {
	res, err := p.queryProxyConnection(height, p.upstreamPathEnd().ConnectionID)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return emptyConnRes, nil
	} else if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *ProxyChainProver) QueryProxyChannelWithProof(height int64) (chanRes *channeltypes.QueryChannelResponse, err error) {
	res, err := p.queryProxyChannel(height, p.upstreamPathEnd().PortID, p.upstreamPathEnd().ChannelID)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return emptyChannelRes, nil
	} else if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *ProxyChainProver) QueryProxyPacketCommitmentWithProof(height int64, seq uint64) (comRes *channeltypes.QueryPacketCommitmentResponse, err error) {
	panic("not implemented error")
}

func (p *ProxyChainProver) QueryProxyPacketAcknowledgementCommitmentWithProof(height int64, seq uint64) (ackRes *channeltypes.QueryPacketAcknowledgementResponse, err error) {
	return p.queryProxyAcknowledgementCommitment(height, p.upstreamPathEnd().PortID, p.upstreamPathEnd().ChannelID, seq)
}

func (p *ProxyChainProver) queryProxyClientConsensusState(height int64, clientID string, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
	value, proof, proofHeight, err := p.queryProxy(height, ibcproxytypes.ProxyConsensusStateKey(p.upstreamPrefix(), p.upstreamClientID(), clientID, dstClientConsHeight))
	if err != nil {
		return nil, err
	}
	// check if client consensus exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	consensusState, err := clienttypes.UnmarshalConsensusState(p.Codec(), value)
	if err != nil {
		return nil, err
	}
	anyConsensusState, err := clienttypes.PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}
	return clienttypes.NewQueryConsensusStateResponse(anyConsensusState, proof, proofHeight), nil
}

func (p *ProxyChainProver) queryProxyClientState(height int64, clientID string) (*clienttypes.QueryClientStateResponse, error) {
	value, proof, proofHeight, err := p.queryProxy(height, ibcproxytypes.ProxyClientStateKey(p.upstreamPrefix(), p.upstreamClientID(), clientID))
	if err != nil {
		return nil, err
	}
	// check if client exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	clientState, err := clienttypes.UnmarshalClientState(p.Codec(), value)
	if err != nil {
		return nil, err
	}
	anyClientState, err := clienttypes.PackClientState(clientState)
	if err != nil {
		return nil, err
	}
	return clienttypes.NewQueryClientStateResponse(anyClientState, proof, proofHeight), nil
}

func (p *ProxyChainProver) queryProxyConnection(height int64, connectionID string) (*connectiontypes.QueryConnectionResponse, error) {
	value, proof, proofHeight, err := p.queryProxy(height, ibcproxytypes.ProxyConnectionKey(p.upstreamPrefix(), p.upstreamClientID(), connectionID))
	if err != nil {
		return nil, err
	}
	// check if connection exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(connectiontypes.ErrConnectionNotFound, connectionID)
	}

	var connection connectiontypes.ConnectionEnd
	if err := p.Codec().Unmarshal(value, &connection); err != nil {
		return nil, err
	}
	return connectiontypes.NewQueryConnectionResponse(connection, proof, proofHeight), nil
}

func (p *ProxyChainProver) queryProxyChannel(height int64, portID string, channelID string) (*channeltypes.QueryChannelResponse, error) {
	value, proof, proofHeight, err := p.queryProxy(height, ibcproxytypes.ProxyChannelKey(p.upstreamPrefix(), p.upstreamClientID(), portID, channelID))
	if err != nil {
		return nil, err
	}
	// check if channel exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrapf(channeltypes.ErrChannelNotFound, "portID=%v channelID=%v", portID, channelID)
	}

	var channel channeltypes.Channel
	if err := p.Codec().Unmarshal(value, &channel); err != nil {
		return nil, err
	}
	return channeltypes.NewQueryChannelResponse(channel, proof, proofHeight), nil
}

func (p *ProxyChainProver) queryProxyAcknowledgementCommitment(height int64, portID string, channelID string, sequence uint64) (*channeltypes.QueryPacketAcknowledgementResponse, error) {
	value, proof, proofHeight, err := p.queryProxy(height, ibcproxytypes.ProxyAcknowledgementKey(p.upstreamPrefix(), p.upstreamClientID(), portID, channelID, sequence))
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return nil, sdkerrors.Wrapf(channeltypes.ErrAcknowledgementExists, "portID=%v channelID=%v sequence=%v", portID, channelID, sequence)
	}
	return channeltypes.NewQueryPacketAcknowledgementResponse(value, proof, proofHeight), nil
}

func (p *ProxyChainProver) queryProxy(height int64, key []byte) (value []byte, proof []byte, proofHeight clienttypes.Height, err error) {
	r, err := p.proxyChain.Client.ABCIQueryWithOptions(
		context.TODO(),
		fmt.Sprintf("store/%s/key", ibcproxytypes.StoreKey),
		key,
		client.ABCIQueryOptions{
			Height: height,
			Prove:  true,
		},
	)
	if err != nil {
		return
	}

	res := r.Response
	if !res.IsOK() {
		err = sdkErrorToGRPCError(res)
		return
	}

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	if err != nil {
		return
	}

	proof, err = p.Codec().Marshal(&merkleProof)
	if err != nil {
		return
	}

	revision := clienttypes.ParseChainID(p.proxyChain.ChainID())

	// proof height + 1 is returned as the proof created corresponds to the height the proof
	// was created in the IAVL tree. Tendermint and subsequently the clients that rely on it
	// have heights 1 above the IAVL tree. Thus we return proof height + 1
	return res.Value, proof, clienttypes.NewHeight(revision, uint64(res.Height)+1), nil
}

func sdkErrorToGRPCError(resp abci.ResponseQuery) error {
	switch resp.Code {
	case sdkerrors.ErrInvalidRequest.ABCICode():
		return status.Error(codes.InvalidArgument, resp.Log)
	case sdkerrors.ErrUnauthorized.ABCICode():
		return status.Error(codes.Unauthenticated, resp.Log)
	case sdkerrors.ErrKeyNotFound.ABCICode():
		return status.Error(codes.NotFound, resp.Log)
	default:
		return status.Error(codes.Unknown, resp.Log)
	}
}

var emptyConnRes = connectiontypes.NewQueryConnectionResponse(
	connectiontypes.NewConnectionEnd(
		connectiontypes.UNINITIALIZED,
		"client",
		connectiontypes.NewCounterparty(
			"client",
			"connection",
			commitmenttypes.NewMerklePrefix([]byte{}),
		),
		[]*connectiontypes.Version{},
		0,
	),
	[]byte{},
	clienttypes.NewHeight(0, 0),
)

var emptyChannelRes = channeltypes.NewQueryChannelResponse(
	channeltypes.NewChannel(
		channeltypes.UNINITIALIZED,
		channeltypes.UNORDERED,
		channeltypes.NewCounterparty(
			"port",
			"channel",
		),
		[]string{},
		"version",
	),
	[]byte{},
	clienttypes.NewHeight(0, 0),
)
