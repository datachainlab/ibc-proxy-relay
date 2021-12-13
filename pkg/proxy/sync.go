package proxy

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"
	multivtypes "github.com/datachainlab/ibc-proxy/modules/light-clients/xx-multiv/types"
	proxytypes "github.com/datachainlab/ibc-proxy/modules/proxy/types"
	"github.com/hyperledger-labs/yui-relayer/core"
)

var (
	ConnectionVersion         = connectiontypes.ExportedVersionsToProto(connectiontypes.GetCompatibleVersions())[0]
	DefaultDelayPeriod uint64 = 0
)

type ProxySynchronizer struct {
	path            *core.PathEnd
	upstream        *core.ProvableChain
	upstreamProxy   *UpstreamProxy
	downstreamProxy *DownstreamProxy
}

func NewProxySynchronizer(
	upstream *core.ProvableChain,
	upstreamProxy *UpstreamProxy,
	downstreamProxy *DownstreamProxy,
) *ProxySynchronizer {
	return &ProxySynchronizer{
		upstream:        upstream,
		upstreamProxy:   upstreamProxy,
		downstreamProxy: downstreamProxy,
	}
}

func (ps *ProxySynchronizer) SetPath(path *core.PathEnd) {
	ps.path = path
}

func (ps ProxySynchronizer) SyncALL() error {
	if err := ps.TrySyncClientState(); err != nil {
		return err
	}
	if err := ps.TrySyncConnectionState(); err != nil {
		return err
	}
	if err := ps.TrySyncChannelState(); err != nil {
		return err
	}
	if err := ps.TrySyncPacketState(); err != nil {
		return err
	}
	return nil
}

func (ps ProxySynchronizer) TrySyncClientState() error {
	return retry.Do(
		func() error {
			_, err := ps.upstreamProxy.QueryClientState(0)
			if err == nil {
				return nil
			} else if !strings.Contains(err.Error(), "lightclient not found") {
				return err
			}
			return ps.SyncCreateClient()
		},
		retry.Delay(1*time.Second),
		retry.Attempts(30),
	)
}

func (ps ProxySynchronizer) TrySyncConnectionState() error {
	panic("not implemented error")
	// return retry.Do(
	// 	func() error {
	// 		// check if the connection update exists
	// 		connRes, err := ps.upstream.QueryConnection(0)
	// 		if err != nil { // TODO if not found, returns nil
	// 			return err
	// 		}
	// 		proxyRes, err := ps.upstreamProxy.QueryProxyConnectionState(0)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		if connRes.Connection.State != proxyRes.Connection.State {
	// 			return nil
	// 		}

	// 		// sync with the upstream state
	// 		// TODO fix state checks
	// 		switch connRes.Connection.State {
	// 		case connectiontypes.INIT:
	// 			return ps.SyncConnectionOpenInit(connRes.Connection.Counterparty)
	// 		case connectiontypes.TRYOPEN:
	// 			return ps.SyncConnectionOpenTry(connRes.Connection.Counterparty)
	// 		case connectiontypes.OPEN:
	// 			return ps.SyncConnectionOpenAck(connRes.Connection.Counterparty.ConnectionId)
	// 		default:
	// 			return fmt.Errorf("unexpected state '%v'", connRes.Connection.State)
	// 		}
	// 	},
	// 	retry.Delay(1*time.Second),
	// 	retry.Attempts(30),
	// )
}

func (ps ProxySynchronizer) TrySyncChannelState() error {
	panic("not implemented error")
	// return retry.Do(
	// 	func() error {
	// 		chanRes, err := ps.upstream.QueryChannel(0)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		proxyRes, err := ps.upstreamProxy.QueryProxyChannel(0)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if chanRes.Channel.State != proxyRes.Channel.State {
	// 			return nil
	// 		}
	// 		return nil
	// 	},
	// 	retry.Delay(1*time.Second),
	// 	retry.Attempts(30),
	// )
}

func (ps ProxySynchronizer) TrySyncPacketState() error {
	panic("not implemented error")
	// return retry.Do(
	// 	func() error {
	// 		return nil
	// 	},
	// 	retry.Delay(1*time.Second),
	// 	retry.Attempts(30),
	// )
}

// SyncCreateClient creates an upstream client on the proxy
func (ps ProxySynchronizer) SyncCreateClient() error {
	header, err := ps.upstream.QueryLatestHeader()
	if err != nil {
		return err
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}
	proxyMsg, err := ps.upstream.CreateMsgCreateClient(ps.upstreamProxy.ProxyPath().UpstreamClientID, header, signer)
	if err != nil {
		return err
	}
	if _, err := ps.upstreamProxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return err
	}
	return nil
}

func (ps ProxySynchronizer) SyncClientState() error {
	provableHeight, err := ps.updateProxyUpstreamClient()
	if err != nil {
		return fmt.Errorf("failed to updateProxyUpstreamClient: %w", err)
	}
	clientRes, err := ps.upstream.QueryClientStateWithProof(provableHeight)
	if err != nil {
		return fmt.Errorf("failed to QueryClientStateWithProof: %w", err)
	}
	clientState := clientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
	consensusHeight := clientState.GetLatestHeight()
	consensusRes, err := ps.upstream.QueryClientConsensusStateWithProof(provableHeight, consensusHeight)
	if err != nil {
		return fmt.Errorf("failed to QueryClientConsensusStateWithProof: %w", err)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}
	var proxyMsg *proxytypes.MsgProxyClientState
	if ps.downstreamProxy == nil {
		proxyMsg = &proxytypes.MsgProxyClientState{
			UpstreamClientId:     ps.upstreamProxy.UpstreamClientID,
			UpstreamPrefix:       commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			CounterpartyClientId: ps.path.ClientID,
			ClientState:          clientRes.ClientState,
			ConsensusState:       consensusRes.ConsensusState,
			ProofClient:          clientRes.Proof,
			ProofConsensus:       consensusRes.Proof,
			ProofHeight:          clientRes.ProofHeight,
			ConsensusHeight:      consensusHeight.(clienttypes.Height),
			Signer:               signer.String(),
		}
	} else {
		head := &multivtypes.BranchProof{
			ClientProof:     clientRes.Proof,
			ClientState:     clientRes.ClientState,
			ConsensusProof:  consensusRes.Proof,
			ConsensusState:  consensusRes.ConsensusState,
			ConsensusHeight: consensusHeight.(clienttypes.Height),
			ProofHeight:     clientRes.ProofHeight,
		}
		proxyClientRes, err := ps.downstreamProxy.QueryClientStateWithProof(int64(clientState.GetLatestHeight().GetRevisionHeight()) - 1)
		if err != nil {
			return fmt.Errorf("failed to downstreamProxy.QueryClientStateWithProof: %w", err)
		}
		leafClient := &multivtypes.LeafClientProof{
			Proof:       proxyClientRes.Proof,
			ProofHeight: proxyClientRes.ProofHeight,
		}
		proofClient := makeClientStateProof(ps.upstream.Codec(), leafClient, head)
		lc := proxyClientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
		proxyConsensusRes, err := ps.downstreamProxy.QueryClientConsensusStateWithProof(int64(clientState.GetLatestHeight().GetRevisionHeight())-1, lc.GetLatestHeight())
		if err != nil {
			return fmt.Errorf("failed to downstreamProxy.QueryClientConsensusStateWithProof: %w", err)
		}
		leafConsensus := &multivtypes.LeafConsensusProof{
			Proof:           proxyConsensusRes.Proof,
			ProofHeight:     proxyConsensusRes.ProofHeight,
			ConsensusHeight: lc.GetLatestHeight().(clienttypes.Height),
		}
		proofConsensus := makeConsensusStateProof(ps.upstream.Codec(), leafConsensus, head)

		proxyMsg = &proxytypes.MsgProxyClientState{
			UpstreamClientId:     ps.upstreamProxy.UpstreamClientID,
			UpstreamPrefix:       commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			CounterpartyClientId: ps.path.ClientID,
			ClientState:          proxyClientRes.ClientState,
			ConsensusState:       proxyConsensusRes.ConsensusState,
			ProofClient:          proofClient,
			ProofConsensus:       proofConsensus,
			ProofHeight:          clientRes.ProofHeight,
			ConsensusHeight:      leafConsensus.ConsensusHeight,
			Signer:               signer.String(),
		}
	}
	if _, err := ps.upstreamProxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return fmt.Errorf("failed to SendMsgs: %w", err)
	}
	return nil
}

func (ps ProxySynchronizer) SyncConnectionOpenInit(connCP connectiontypes.Counterparty) error {
	provableHeight, err := ps.updateProxyUpstreamClient()
	if err != nil {
		return fmt.Errorf("failed to updateProxyUpstreamClient: %w", err)
	}
	clientRes, err := ps.upstream.QueryClientStateWithProof(provableHeight)
	if err != nil {
		return fmt.Errorf("failed to QueryClientStateWithProof: %w", err)
	}
	clientState := clientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
	consensusHeight := clientState.GetLatestHeight()
	consensusRes, err := ps.upstream.QueryClientConsensusStateWithProof(provableHeight, consensusHeight)
	if err != nil {
		return fmt.Errorf("failed to QueryClientConsensusStateWithProof: %w", err)
	}
	connRes, err := ps.upstream.QueryConnectionWithProof(provableHeight)
	if err != nil {
		return fmt.Errorf("failed to QueryConnectionWithProof: %w", err)
	}
	log.Println("onConnectionOpenInit:", connRes.Proof, connRes.Connection)
	if len(connRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the connection(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}

	var proxyMsg *proxytypes.MsgProxyConnectionOpenTry
	if ps.downstreamProxy == nil {
		proxyMsg = &proxytypes.MsgProxyConnectionOpenTry{
			ConnectionId:     ps.path.ConnectionID,
			UpstreamClientId: ps.upstreamProxy.UpstreamClientID,
			UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			Connection: connectiontypes.NewConnectionEnd(
				connectiontypes.INIT,
				ps.path.ClientID,
				connCP,
				[]*connectiontypes.Version{ConnectionVersion}, DefaultDelayPeriod,
			),
			ClientState:     clientRes.ClientState,
			ConsensusState:  consensusRes.ConsensusState,
			ProofInit:       connRes.Proof,
			ProofClient:     clientRes.Proof,
			ProofConsensus:  consensusRes.Proof,
			ProofHeight:     clientRes.ProofHeight,
			ConsensusHeight: consensusHeight.(clienttypes.Height),
			Signer:          signer.String(),
		}
	} else {
		head := &multivtypes.BranchProof{
			ClientProof:     clientRes.Proof,
			ClientState:     clientRes.ClientState,
			ConsensusProof:  consensusRes.Proof,
			ConsensusState:  consensusRes.ConsensusState,
			ConsensusHeight: consensusHeight.(clienttypes.Height),
			ProofHeight:     clientRes.ProofHeight,
		}
		proxyClientRes, err := ps.downstreamProxy.QueryClientStateWithProof(int64(clientState.GetLatestHeight().GetRevisionHeight()) - 1)
		if err != nil {
			return err
		}
		leafClient := &multivtypes.LeafClientProof{
			Proof:       proxyClientRes.Proof,
			ProofHeight: proxyClientRes.ProofHeight,
		}
		proofClient := makeClientStateProof(ps.upstream.Codec(), leafClient, head)
		lc := proxyClientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
		proxyConsensusRes, err := ps.downstreamProxy.QueryClientConsensusStateWithProof(int64(clientState.GetLatestHeight().GetRevisionHeight())-1, lc.GetLatestHeight())
		if err != nil {
			return err
		}
		leafConsensus := &multivtypes.LeafConsensusProof{
			Proof:           proxyConsensusRes.Proof,
			ProofHeight:     proxyConsensusRes.ProofHeight,
			ConsensusHeight: lc.GetLatestHeight().(clienttypes.Height),
		}
		proofConsensus := makeConsensusStateProof(ps.upstream.Codec(), leafConsensus, head)

		proxyMsg = &proxytypes.MsgProxyConnectionOpenTry{
			ConnectionId:     ps.path.ConnectionID,
			UpstreamClientId: ps.upstreamProxy.UpstreamClientID,
			UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			Connection: connectiontypes.NewConnectionEnd(
				connectiontypes.INIT,
				ps.path.ClientID,
				connCP,
				[]*connectiontypes.Version{ConnectionVersion}, DefaultDelayPeriod,
			),
			ClientState:     proxyClientRes.ClientState,
			ConsensusState:  proxyConsensusRes.ConsensusState,
			ProofInit:       connRes.Proof,
			ProofClient:     proofClient,
			ProofConsensus:  proofConsensus,
			ProofHeight:     connRes.ProofHeight,
			ConsensusHeight: leafConsensus.ConsensusHeight,
			Signer:          signer.String(),
		}
	}

	if _, err := ps.upstreamProxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return fmt.Errorf("failed to SendMsgs: %w", err)
	}
	return nil
}

func (ps ProxySynchronizer) SyncConnectionOpenTry(connCP connectiontypes.Counterparty) error {
	provableHeight, err := ps.updateProxyUpstreamClient()
	if err != nil {
		return fmt.Errorf("failed to updateProxyUpstreamClient: %w", err)
	}
	clientRes, err := ps.upstream.QueryClientStateWithProof(provableHeight)
	if err != nil {
		return fmt.Errorf("failed to QueryClientStateWithProof: %w", err)
	}
	clientState := clientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
	consensusHeight := clientState.GetLatestHeight()
	consensusRes, err := ps.upstream.QueryClientConsensusStateWithProof(provableHeight, consensusHeight)
	if err != nil {
		return fmt.Errorf("failed to QueryClientConsensusStateWithProof: %w", err)
	}
	connRes, err := ps.upstream.QueryConnectionWithProof(provableHeight)
	if err != nil {
		return fmt.Errorf("failed to QueryConnectionWithProof: %w", err)
	}
	log.Println("onConnectionOpenTry:", connRes.Proof, connRes.Connection)
	if len(connRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the connection(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}

	var proxyMsg *proxytypes.MsgProxyConnectionOpenAck
	if ps.downstreamProxy == nil {
		proxyMsg = &proxytypes.MsgProxyConnectionOpenAck{
			ConnectionId:     ps.path.ConnectionID,
			UpstreamClientId: ps.upstreamProxy.UpstreamClientID,
			UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			Connection: connectiontypes.NewConnectionEnd(
				connectiontypes.TRYOPEN,
				ps.path.ClientID,
				connCP,
				[]*connectiontypes.Version{ConnectionVersion}, DefaultDelayPeriod,
			),
			ClientState:     clientRes.ClientState,
			ConsensusState:  consensusRes.ConsensusState,
			ProofTry:        connRes.Proof,
			ProofClient:     clientRes.Proof,
			ProofConsensus:  consensusRes.Proof,
			ProofHeight:     clientRes.ProofHeight,
			ConsensusHeight: consensusHeight.(clienttypes.Height),
			Signer:          signer.String(),
		}
	} else {
		head := &multivtypes.BranchProof{
			ClientProof:     clientRes.Proof,
			ClientState:     clientRes.ClientState,
			ConsensusProof:  consensusRes.Proof,
			ConsensusState:  consensusRes.ConsensusState,
			ConsensusHeight: consensusHeight.(clienttypes.Height),
			ProofHeight:     clientRes.ProofHeight,
		}
		proxyClientRes, err := ps.downstreamProxy.QueryClientStateWithProof(int64(clientState.GetLatestHeight().GetRevisionHeight()) - 1)
		if err != nil {
			return err
		}
		leafClient := &multivtypes.LeafClientProof{
			Proof:       proxyClientRes.Proof,
			ProofHeight: proxyClientRes.ProofHeight,
		}
		proofClient := makeClientStateProof(ps.upstream.Codec(), leafClient, head)

		lc := proxyClientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
		proxyConsensusRes, err := ps.downstreamProxy.QueryClientConsensusStateWithProof(int64(clientState.GetLatestHeight().GetRevisionHeight())-1, lc.GetLatestHeight())
		if err != nil {
			return err
		}
		leafConsensus := &multivtypes.LeafConsensusProof{
			Proof:           proxyConsensusRes.Proof,
			ProofHeight:     proxyConsensusRes.ProofHeight,
			ConsensusHeight: lc.GetLatestHeight().(clienttypes.Height),
		}
		proofConsensus := makeConsensusStateProof(ps.upstream.Codec(), leafConsensus, head)

		proxyMsg = &proxytypes.MsgProxyConnectionOpenAck{
			ConnectionId:     ps.path.ConnectionID,
			UpstreamClientId: ps.upstreamProxy.UpstreamClientID,
			UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			Connection: connectiontypes.NewConnectionEnd(
				connectiontypes.TRYOPEN,
				ps.path.ClientID,
				connCP,
				[]*connectiontypes.Version{ConnectionVersion}, DefaultDelayPeriod,
			),
			ClientState:     proxyClientRes.ClientState,
			ConsensusState:  proxyConsensusRes.ConsensusState,
			ProofTry:        connRes.Proof,
			ProofClient:     proofClient,
			ProofConsensus:  proofConsensus,
			ProofHeight:     connRes.ProofHeight,
			ConsensusHeight: leafConsensus.ConsensusHeight,
			Signer:          signer.String(),
		}
	}
	if _, err := ps.upstreamProxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return fmt.Errorf("failed to SendMsgs: %w", err)
	}
	return nil
}

func (ps ProxySynchronizer) SyncConnectionOpenAck(counterpartyConnectionID string) error {
	provableHeight, err := ps.updateProxyUpstreamClient()
	if err != nil {
		return err
	}
	connRes, err := ps.upstream.QueryConnectionWithProof(provableHeight)
	if err != nil {
		return err
	}
	log.Println("onConnectionOpenAck:", connRes.Proof, connRes.Connection)
	if len(connRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the connection(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}
	proxyMsg := &proxytypes.MsgProxyConnectionOpenConfirm{
		ConnectionId:     ps.path.ConnectionID,
		UpstreamClientId: ps.upstreamProxy.UpstreamClientID,
		UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		Connection: connectiontypes.NewConnectionEnd(
			connectiontypes.OPEN,
			ps.path.ClientID,
			connectiontypes.Counterparty{
				ClientId:     connRes.Connection.Counterparty.ClientId,
				ConnectionId: counterpartyConnectionID,
				Prefix:       commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			},
			[]*connectiontypes.Version{ConnectionVersion}, DefaultDelayPeriod,
		),
		ProofAck:    connRes.Proof,
		ProofHeight: connRes.ProofHeight,
		Signer:      signer.String(),
	}
	if _, err := ps.upstreamProxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return fmt.Errorf("failed to SendMsgs: %w", err)
	}
	return nil
}

func (ps ProxySynchronizer) SyncChannelOpenInit() error {
	provableHeight, err := ps.updateProxyUpstreamClient()
	if err != nil {
		return err
	}
	chanRes, err := ps.upstream.QueryChannelWithProof(provableHeight)
	if err != nil {
		return err
	}
	log.Println("onChannelOpenInit:", chanRes.Proof, chanRes.Channel)
	if len(chanRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the channel(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}
	proxyMsg := &proxytypes.MsgProxyChannelOpenTry{
		UpstreamClientId: ps.upstreamProxy.UpstreamClientID,
		UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		Order:            chanRes.Channel.Ordering,
		ConnectionHops:   chanRes.Channel.ConnectionHops,
		PortId:           ps.path.PortID,
		ChannelId:        ps.path.ChannelID,
		DownstreamPortId: chanRes.Channel.Counterparty.PortId,
		Version:          chanRes.Channel.Version,
		ProofInit:        chanRes.Proof,
		ProofHeight:      chanRes.ProofHeight,
		Signer:           signer.String(),
	}
	if _, err := ps.upstreamProxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return err
	}
	return nil
}

func (ps ProxySynchronizer) SyncChannelOpenTry() error {
	provableHeight, err := ps.updateProxyUpstreamClient()
	if err != nil {
		return err
	}
	chanRes, err := ps.upstream.QueryChannelWithProof(provableHeight)
	if err != nil {
		return err
	}
	log.Println("onChannelOpenTry:", chanRes.Proof, chanRes.Channel)
	if len(chanRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the channel(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}
	proxyMsg := &proxytypes.MsgProxyChannelOpenAck{
		UpstreamClientId:    ps.upstreamProxy.UpstreamClientID,
		UpstreamPrefix:      commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		Order:               chanRes.Channel.Ordering,
		ConnectionHops:      chanRes.Channel.ConnectionHops,
		PortId:              ps.path.PortID,
		ChannelId:           ps.path.ChannelID,
		DownstreamPortId:    chanRes.Channel.Counterparty.PortId,
		DownstreamChannelId: chanRes.Channel.Counterparty.ChannelId,
		Version:             chanRes.Channel.Version,
		ProofTry:            chanRes.Proof,
		ProofHeight:         chanRes.ProofHeight,
		Signer:              signer.String(),
	}
	if _, err := ps.upstreamProxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return err
	}
	return nil
}

func (ps ProxySynchronizer) SyncChannelOpenAck() error {
	provableHeight, err := ps.updateProxyUpstreamClient()
	if err != nil {
		return err
	}
	chanRes, err := ps.upstream.QueryChannelWithProof(provableHeight)
	if err != nil {
		return err
	}
	log.Println("onChannelOpenAck:", chanRes.Proof, chanRes.Channel)
	if len(chanRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the channel(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}
	proxyMsg := &proxytypes.MsgProxyChannelOpenConfirm{
		UpstreamClientId:    ps.upstreamProxy.UpstreamClientID,
		UpstreamPrefix:      commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		PortId:              ps.path.PortID,
		ChannelId:           ps.path.ChannelID,
		DownstreamChannelId: chanRes.Channel.Counterparty.ChannelId,
		ProofAck:            chanRes.Proof,
		ProofHeight:         chanRes.ProofHeight,
		Signer:              signer.String(),
	}
	if _, err := ps.upstreamProxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return err
	}
	return nil
}

func (ps ProxySynchronizer) SyncRecvPacket(packet channeltypes.Packet) error {
	provableHeight, err := ps.updateProxyUpstreamClient()
	if err != nil {
		return err
	}
	res, err := ps.upstream.QueryPacketAcknowledgementCommitmentWithProof(provableHeight, packet.Sequence)
	if err != nil {
		return err
	}
	log.Println("onRecvPacket:", res.Proof, res.Acknowledgement)

	ack, err := ps.upstream.QueryPacketAcknowledgement(provableHeight, packet.Sequence)
	if err != nil {
		return err
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}
	proxyMsg := &proxytypes.MsgProxyAcknowledgePacket{
		UpstreamClientId: ps.upstreamProxy.UpstreamClientID,
		UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		Packet:           packet,
		Acknowledgement:  ack,
		Proof:            res.Proof,
		ProofHeight:      res.ProofHeight,
		Signer:           signer.String(),
	}
	if _, err := ps.upstreamProxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return err
	}
	return nil
}

// updateProxyUpstreamClient updates the upstream client on the proxy
func (ps ProxySynchronizer) updateProxyUpstreamClient() (int64, error) {
	h, provableHeight, _, err := ps.upstream.UpdateLightWithHeader()
	if err != nil {
		return 0, err
	}
	header, err := ps.upstream.SetupHeader(ps.upstreamProxy, h)
	if err != nil {
		return 0, err
	}
	if header != nil {
		if err := header.ValidateBasic(); err != nil {
			return 0, err
		}
		addr, err := ps.upstreamProxy.GetAddress()
		if err != nil {
			return 0, err
		}
		_, err = ps.upstreamProxy.SendMsgs(
			[]sdk.Msg{ps.upstreamProxy.Path().UpdateClient(header, addr)},
		)
		if err != nil {
			return 0, err
		}
	}
	return provableHeight, nil
}

func makeClientStateProof(
	cdc codec.Codec,
	leafClient *multivtypes.LeafClientProof,
	branches ...*multivtypes.BranchProof,
) []byte {
	var mp multivtypes.MultiProof

	for _, branch := range branches {
		mp.Proofs = append(mp.Proofs, &multivtypes.Proof{
			Proof: &multivtypes.Proof_Branch{Branch: branch},
		})
	}
	mp.Proofs = append(mp.Proofs, &multivtypes.Proof{
		Proof: &multivtypes.Proof_LeafClient{LeafClient: leafClient},
	})

	any, err := codectypes.NewAnyWithValue(&mp)
	if err != nil {
		panic(err)
	}
	bz, err := cdc.Marshal(any)
	if err != nil {
		panic(err)
	}
	return bz
}

func makeConsensusStateProof(
	cdc codec.Codec,
	leafConsensus *multivtypes.LeafConsensusProof,
	branches ...*multivtypes.BranchProof,
) []byte {
	var mp multivtypes.MultiProof

	for _, branch := range branches {
		mp.Proofs = append(mp.Proofs, &multivtypes.Proof{
			Proof: &multivtypes.Proof_Branch{Branch: branch},
		})
	}
	mp.Proofs = append(mp.Proofs, &multivtypes.Proof{
		Proof: &multivtypes.Proof_LeafConsensus{LeafConsensus: leafConsensus},
	})

	any, err := codectypes.NewAnyWithValue(&mp)
	if err != nil {
		panic(err)
	}
	bz, err := cdc.Marshal(any)
	if err != nil {
		panic(err)
	}
	return bz
}
