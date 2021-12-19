package tendermint

import (
	"context"
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	chantypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/rpc/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/datachainlab/ibc-proxy-relay/pkg/proxy"
	ibcproxytypes "github.com/datachainlab/ibc-proxy/modules/proxy/types"
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
	upstreamPathEnd proxy.ProxyPathEnd
}

var _ proxy.ProxyChainI = (*TendermintProxyChain)(nil)

func NewTendermintProxyChain(chain *tendermint.Chain) *TendermintProxyChain {
	return &TendermintProxyChain{Chain: chain}
}

func (c *TendermintProxyChain) SetUpstreamPathEnd(path proxy.ProxyPathEnd) {
	c.upstreamPathEnd = path
}

func (p *TendermintProxyChain) upstreamPrefix() *commitmenttypes.MerklePrefix {
	prefix := commitmenttypes.NewMerklePrefix([]byte(host.StoreKey))
	return &prefix
}

func (c *TendermintProxyChain) QueryProxyClientState(height int64) (*clienttypes.QueryClientStateResponse, error) {
	return c.ABCIQueryProxyClientState(height, c.upstreamPathEnd.ClientID(), false)
}

func (c *TendermintProxyChain) QueryProxyClientConsensusState(height int64, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
	return c.ABCIQueryProxyConsensusState(height, c.upstreamPathEnd.ClientID(), dstClientConsHeight, false)
}

func (c *TendermintProxyChain) QueryProxyConnectionState(height int64) (*connectiontypes.QueryConnectionResponse, error) {
	return c.ABCIQueryProxyConnection(height, c.upstreamPathEnd.ConnectionID(), false)
}

func (c *TendermintProxyChain) QueryProxyChannel(height int64) (chanRes *chantypes.QueryChannelResponse, err error) {
	return c.ABCIQueryProxyChannel(height, c.upstreamPathEnd.PortID(), c.upstreamPathEnd.ChannelID(), false)
}

func (c *TendermintProxyChain) QueryProxyPacketCommitment(height int64, seq uint64) (comRes *chantypes.QueryPacketCommitmentResponse, err error) {
	return c.ABCIQueryProxyPacketCommitment(height, seq, false)
}

func (c *TendermintProxyChain) QueryProxyPacketAcknowledgementCommitment(height int64, seq uint64) (ackRes *chantypes.QueryPacketAcknowledgementResponse, err error) {
	return c.ABCIQueryProxyPacketAcknowledgementCommitment(height, seq, false)
}

func (c *TendermintProxyChain) ABCIQueryProxy(height int64, key []byte, prove bool) (value []byte, proof []byte, proofHeight clienttypes.Height, err error) {
	r, err := c.Client.ABCIQueryWithOptions(
		context.TODO(),
		fmt.Sprintf("store/%s/key", ibcproxytypes.StoreKey),
		key,
		client.ABCIQueryOptions{
			Height: height,
			Prove:  prove,
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

	if prove {
		var merkleProof commitmenttypes.MerkleProof
		merkleProof, err = commitmenttypes.ConvertProofs(res.ProofOps)
		if err != nil {
			return
		}

		proof, err = c.Codec().Marshal(&merkleProof)
		if err != nil {
			return
		}

		revision := clienttypes.ParseChainID(c.ChainID())
		proofHeight = clienttypes.NewHeight(revision, uint64(res.Height)+1)
	}

	return res.Value, proof, proofHeight, nil
}

func (c *TendermintProxyChain) ABCIQueryProxyClientState(height int64, clientID string, prove bool) (*clienttypes.QueryClientStateResponse, error) {
	value, proof, proofHeight, err := c.ABCIQueryProxy(height, ibcproxytypes.ProxyClientStateKey(c.upstreamPrefix(), c.upstreamPathEnd.UpstreamClientId, clientID), prove)
	if err != nil {
		return nil, err
	}
	// check if client exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	clientState, err := clienttypes.UnmarshalClientState(c.Codec(), value)
	if err != nil {
		return nil, err
	}
	anyClientState, err := clienttypes.PackClientState(clientState)
	if err != nil {
		return nil, err
	}
	return clienttypes.NewQueryClientStateResponse(anyClientState, proof, proofHeight), nil
}

func (c *TendermintProxyChain) ABCIQueryProxyConsensusState(height int64, clientID string, dstClientConsHeight ibcexported.Height, prove bool) (*clienttypes.QueryConsensusStateResponse, error) {
	value, proof, proofHeight, err := c.ABCIQueryProxy(height, ibcproxytypes.ProxyConsensusStateKey(c.upstreamPrefix(), c.upstreamPathEnd.UpstreamClientId, clientID, dstClientConsHeight), prove)
	if err != nil {
		return nil, err
	}
	// check if client consensus exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	consensusState, err := clienttypes.UnmarshalConsensusState(c.Codec(), value)
	if err != nil {
		return nil, err
	}
	anyConsensusState, err := clienttypes.PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}
	return clienttypes.NewQueryConsensusStateResponse(anyConsensusState, proof, proofHeight), nil
}

func (c *TendermintProxyChain) ABCIQueryProxyConnection(height int64, connectionID string, prove bool) (*connectiontypes.QueryConnectionResponse, error) {
	value, proof, proofHeight, err := c.ABCIQueryProxy(height, ibcproxytypes.ProxyConnectionKey(c.upstreamPrefix(), c.upstreamPathEnd.UpstreamClientId, connectionID), prove)
	if err != nil {
		return nil, err
	}
	// check if connection exists
	if len(value) == 0 {
		return emptyConnRes, nil
	}

	var connection connectiontypes.ConnectionEnd
	if err := c.Codec().Unmarshal(value, &connection); err != nil {
		return nil, err
	}
	return connectiontypes.NewQueryConnectionResponse(connection, proof, proofHeight), nil
}

func (c *TendermintProxyChain) ABCIQueryProxyChannel(height int64, portID string, channelID string, prove bool) (*channeltypes.QueryChannelResponse, error) {
	value, proof, proofHeight, err := c.ABCIQueryProxy(height, ibcproxytypes.ProxyChannelKey(c.upstreamPrefix(), c.upstreamPathEnd.UpstreamClientId, portID, channelID), prove)
	if err != nil {
		return nil, err
	}
	// check if channel exists
	if len(value) == 0 {
		return emptyChannelRes, nil
	}

	var channel channeltypes.Channel
	if err := c.Codec().Unmarshal(value, &channel); err != nil {
		return nil, err
	}
	return channeltypes.NewQueryChannelResponse(channel, proof, proofHeight), nil
}

func (c *TendermintProxyChain) ABCIQueryProxyPacketCommitment(height int64, seq uint64, prove bool) (comRes *channeltypes.QueryPacketCommitmentResponse, err error) {
	portID, channelID := c.upstreamPathEnd.PortId, c.upstreamPathEnd.ChannelId
	value, proof, proofHeight, err := c.ABCIQueryProxy(height, ibcproxytypes.ProxyPacketCommitmentKey(c.upstreamPrefix(), c.upstreamPathEnd.UpstreamClientId, portID, channelID, seq), prove)
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return nil, sdkerrors.Wrapf(channeltypes.ErrPacketCommitmentNotFound, "portID=%v channelID=%v sequence=%v", portID, channelID, seq)
	}
	return channeltypes.NewQueryPacketCommitmentResponse(value, proof, proofHeight), nil
}

func (c *TendermintProxyChain) ABCIQueryProxyPacketAcknowledgementCommitment(height int64, seq uint64, prove bool) (*channeltypes.QueryPacketAcknowledgementResponse, error) {
	portID, channelID := c.upstreamPathEnd.PortId, c.upstreamPathEnd.ChannelId
	value, proof, proofHeight, err := c.ABCIQueryProxy(height, ibcproxytypes.ProxyAcknowledgementKey(c.upstreamPrefix(), c.upstreamPathEnd.UpstreamClientId, portID, channelID, seq), prove)
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return nil, sdkerrors.Wrapf(channeltypes.ErrAcknowledgementExists, "portID=%v channelID=%v sequence=%v", portID, channelID, seq)
	}
	return channeltypes.NewQueryPacketAcknowledgementResponse(value, proof, proofHeight), nil
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
