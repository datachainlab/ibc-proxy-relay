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
	path             *core.PathEnd
	upstream         *core.ProvableChain
	downstreamProver core.ProverI
	upstreamProxy    *ProxyProvableChain
	downstreamProxy  *ProxyProvableChain
}

func NewProxySynchronizer(
	path *core.PathEnd,
	upstream *core.ProvableChain,
	downstreamProver core.ProverI,
	upstreamProxy *ProxyProvableChain,
	downstreamProxy *ProxyProvableChain,
) *ProxySynchronizer {
	return &ProxySynchronizer{
		path:             path,
		upstream:         upstream,
		downstreamProver: downstreamProver,
		upstreamProxy:    upstreamProxy,
		downstreamProxy:  downstreamProxy,
	}
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
			UpstreamClientId:     ps.upstreamProxy.Path().ClientID,
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
		head := &multivtypes.Proof{
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
		leafClient := &multivtypes.LeafProof{
			Proof:       proxyClientRes.Proof,
			ProofHeight: proxyClientRes.ProofHeight,
		}
		proofClient := makeMultiProof(ps.upstream.Codec(), head, nil, leafClient)
		lc := proxyClientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
		proxyConsensusRes, err := ps.downstreamProxy.QueryClientConsensusStateWithProof(int64(clientState.GetLatestHeight().GetRevisionHeight())-1, lc.GetLatestHeight())
		if err != nil {
			return fmt.Errorf("failed to downstreamProxy.QueryClientConsensusStateWithProof: %w", err)
		}
		leafConsensus := &multivtypes.LeafProof{
			Proof:       proxyConsensusRes.Proof,
			ProofHeight: proxyConsensusRes.ProofHeight,
		}
		proofConsensus := makeMultiProof(ps.upstream.Codec(), head, nil, leafConsensus)

		proxyMsg = &proxytypes.MsgProxyClientState{
			UpstreamClientId:     ps.upstreamProxy.Path().ClientID,
			UpstreamPrefix:       commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			CounterpartyClientId: ps.path.ClientID,
			ClientState:          proxyClientRes.ClientState,
			ConsensusState:       proxyConsensusRes.ConsensusState,
			ProofClient:          proofClient,
			ProofConsensus:       proofConsensus,
			ProofHeight:          clientRes.ProofHeight,
			ConsensusHeight:      lc.GetLatestHeight().(clienttypes.Height),
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
	if len(connRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the connection(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}

	log.Printf("SyncConnectionOpenInit: connection=%v proof=%v", connRes.Connection, connRes.Proof)

	var proxyMsg *proxytypes.MsgProxyConnectionOpenTry
	if ps.downstreamProxy == nil {
		downstreamClientRes, err := ps.downstreamProver.QueryClientStateWithProof(int64(clientState.GetLatestHeight().GetRevisionHeight()) - 1)
		if err != nil {
			return err
		}
		lc := downstreamClientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
		downstreamConsensusRes, err := ps.downstreamProver.QueryClientConsensusStateWithProof(
			int64(clientState.GetLatestHeight().GetRevisionHeight())-1,
			lc.GetLatestHeight(),
		)
		if err != nil {
			return err
		}
		proxyMsg = &proxytypes.MsgProxyConnectionOpenTry{
			ConnectionId:   ps.path.ConnectionID,
			UpstreamPrefix: commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			Connection: connectiontypes.NewConnectionEnd(
				connectiontypes.INIT,
				ps.path.ClientID,
				connCP,
				[]*connectiontypes.Version{ConnectionVersion}, DefaultDelayPeriod,
			),
			DownstreamClientState:    clientRes.ClientState,
			DownstreamConsensusState: consensusRes.ConsensusState,
			ProxyClientState:         downstreamClientRes.ClientState,
			ProofInit:                connRes.Proof,
			ProofClient:              clientRes.Proof,
			ProofConsensus:           consensusRes.Proof,
			ProofHeight:              clientRes.ProofHeight,
			ConsensusHeight:          consensusHeight.(clienttypes.Height),
			ProofProxyClient:         downstreamClientRes.Proof,
			ProofProxyConsensus:      downstreamConsensusRes.Proof,
			ProofProxyHeight:         downstreamClientRes.ProofHeight,
			ProxyConsensusHeight:     lc.GetLatestHeight().(clienttypes.Height),
			Signer:                   signer.String(),
		}
	} else {
		head := &multivtypes.Proof{
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
		leafClient := &multivtypes.LeafProof{
			Proof:       proxyClientRes.Proof,
			ProofHeight: proxyClientRes.ProofHeight,
		}
		proofClient := makeMultiProof(ps.upstream.Codec(), head, nil, leafClient)
		lc := proxyClientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
		proxyConsensusRes, err := ps.downstreamProxy.QueryClientConsensusStateWithProof(int64(clientState.GetLatestHeight().GetRevisionHeight())-1, lc.GetLatestHeight())
		if err != nil {
			return err
		}
		leafConsensus := &multivtypes.LeafProof{
			Proof:       proxyConsensusRes.Proof,
			ProofHeight: proxyConsensusRes.ProofHeight,
		}
		proofConsensus := makeMultiProof(ps.upstream.Codec(), head, nil, leafConsensus)

		downstreamClientRes, err := ps.downstreamProver.QueryClientStateWithProof(
			int64(lc.GetLatestHeight().GetRevisionHeight()) - 1,
		)
		if err != nil {
			return err
		}
		lc2 := downstreamClientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
		downstreamConsensusRes, err := ps.downstreamProver.QueryClientConsensusStateWithProof(
			int64(lc.GetLatestHeight().GetRevisionHeight())-1,
			lc2.GetLatestHeight(),
		)

		proxyMsg = &proxytypes.MsgProxyConnectionOpenTry{
			ConnectionId:   ps.path.ConnectionID,
			UpstreamPrefix: commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			Connection: connectiontypes.NewConnectionEnd(
				connectiontypes.INIT,
				ps.path.ClientID,
				connCP,
				[]*connectiontypes.Version{ConnectionVersion}, DefaultDelayPeriod,
			),
			DownstreamClientState:    proxyClientRes.ClientState,
			DownstreamConsensusState: proxyConsensusRes.ConsensusState,
			ProxyClientState:         downstreamClientRes.ClientState,
			ProofInit:                connRes.Proof,
			ProofClient:              proofClient,
			ProofConsensus:           proofConsensus,
			ProofHeight:              connRes.ProofHeight,
			ConsensusHeight:          lc.GetLatestHeight().(clienttypes.Height),
			ProofProxyClient:         downstreamClientRes.Proof,
			ProofProxyConsensus:      downstreamConsensusRes.Proof,
			ProofProxyHeight:         lc.GetLatestHeight().(clienttypes.Height),
			ProxyConsensusHeight:     lc2.GetLatestHeight().(clienttypes.Height),
			Signer:                   signer.String(),
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
	if len(connRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the connection(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}

	log.Printf("SyncConnectionOpenTry: connection=%v proof=%v", connRes.Connection, connRes.Proof)

	var proxyMsg *proxytypes.MsgProxyConnectionOpenAck
	if ps.downstreamProxy == nil {
		downstreamClientRes, err := ps.downstreamProver.QueryClientStateWithProof(int64(clientState.GetLatestHeight().GetRevisionHeight()) - 1)
		if err != nil {
			return err
		}
		lc := downstreamClientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
		downstreamConsensusRes, err := ps.downstreamProver.QueryClientConsensusStateWithProof(
			int64(clientState.GetLatestHeight().GetRevisionHeight())-1,
			lc.GetLatestHeight(),
		)
		if err != nil {
			return err
		}

		proxyMsg = &proxytypes.MsgProxyConnectionOpenAck{
			ConnectionId:   ps.path.ConnectionID,
			UpstreamPrefix: commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			Connection: connectiontypes.NewConnectionEnd(
				connectiontypes.TRYOPEN,
				ps.path.ClientID,
				connCP,
				[]*connectiontypes.Version{ConnectionVersion}, DefaultDelayPeriod,
			),
			DownstreamClientState:    clientRes.ClientState,
			DownstreamConsensusState: consensusRes.ConsensusState,
			ProxyClientState:         downstreamClientRes.ClientState,
			ProofTry:                 connRes.Proof,
			ProofClient:              clientRes.Proof,
			ProofConsensus:           consensusRes.Proof,
			ProofHeight:              clientRes.ProofHeight,
			ConsensusHeight:          consensusHeight.(clienttypes.Height),
			ProofProxyClient:         downstreamClientRes.Proof,
			ProofProxyConsensus:      downstreamConsensusRes.Proof,
			ProofProxyHeight:         downstreamClientRes.ProofHeight,
			ProxyConsensusHeight:     lc.GetLatestHeight().(clienttypes.Height),
			Signer:                   signer.String(),
		}
	} else {
		head := &multivtypes.Proof{
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
		leafClient := &multivtypes.LeafProof{
			Proof:       proxyClientRes.Proof,
			ProofHeight: proxyClientRes.ProofHeight,
		}
		proofClient := makeMultiProof(ps.upstream.Codec(), head, nil, leafClient)

		lc := proxyClientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
		proxyConsensusRes, err := ps.downstreamProxy.QueryClientConsensusStateWithProof(int64(clientState.GetLatestHeight().GetRevisionHeight())-1, lc.GetLatestHeight())
		if err != nil {
			return err
		}
		leafConsensus := &multivtypes.LeafProof{
			Proof:       proxyConsensusRes.Proof,
			ProofHeight: proxyConsensusRes.ProofHeight,
		}
		proofConsensus := makeMultiProof(ps.upstream.Codec(), head, nil, leafConsensus)

		downstreamClientRes, err := ps.downstreamProver.QueryClientStateWithProof(
			int64(lc.GetLatestHeight().GetRevisionHeight()) - 1,
		)
		if err != nil {
			return err
		}
		lc2 := downstreamClientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
		downstreamConsensusRes, err := ps.downstreamProver.QueryClientConsensusStateWithProof(
			int64(lc.GetLatestHeight().GetRevisionHeight())-1,
			lc2.GetLatestHeight(),
		)

		proxyMsg = &proxytypes.MsgProxyConnectionOpenAck{
			ConnectionId:   ps.path.ConnectionID,
			UpstreamPrefix: commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			Connection: connectiontypes.NewConnectionEnd(
				connectiontypes.TRYOPEN,
				ps.path.ClientID,
				connCP,
				[]*connectiontypes.Version{ConnectionVersion}, DefaultDelayPeriod,
			),
			DownstreamClientState:    proxyClientRes.ClientState,
			DownstreamConsensusState: proxyConsensusRes.ConsensusState,
			ProxyClientState:         downstreamClientRes.ClientState,
			ProofTry:                 connRes.Proof,
			ProofClient:              proofClient,
			ProofConsensus:           proofConsensus,
			ProofHeight:              connRes.ProofHeight,
			ConsensusHeight:          lc.GetLatestHeight().(clienttypes.Height),
			ProofProxyClient:         downstreamClientRes.Proof,
			ProofProxyConsensus:      downstreamConsensusRes.Proof,
			ProofProxyHeight:         lc.GetLatestHeight().(clienttypes.Height),
			ProxyConsensusHeight:     lc2.GetLatestHeight().(clienttypes.Height),
			Signer:                   signer.String(),
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
	if len(connRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the connection(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}

	log.Printf("SyncConnectionOpenAck: connection=%v proof=%v", connRes.Connection, connRes.Proof)

	proxyMsg := &proxytypes.MsgProxyConnectionOpenConfirm{
		ConnectionId:             ps.path.ConnectionID,
		UpstreamClientId:         ps.upstreamProxy.Path().ClientID,
		UpstreamPrefix:           commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		CounterpartyConnectionId: counterpartyConnectionID,
		ProofAck:                 connRes.Proof,
		ProofHeight:              connRes.ProofHeight,
		Signer:                   signer.String(),
	}
	if _, err := ps.upstreamProxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return fmt.Errorf("failed to SendMsgs: %w", err)
	}
	return nil
}

func (ps ProxySynchronizer) SyncConnectionOpenConfirm() error {
	provableHeight, err := ps.updateProxyUpstreamClient()
	if err != nil {
		return err
	}
	connRes, err := ps.upstream.QueryConnectionWithProof(provableHeight)
	if err != nil {
		return err
	}
	if len(connRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the connection(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}

	log.Printf("SyncConnectionOpenConfirm: connection=%v proof=%v", connRes.Connection, connRes.Proof)

	proxyMsg := &proxytypes.MsgProxyConnectionOpenFinalize{
		ConnectionId:     ps.path.ConnectionID,
		UpstreamClientId: ps.upstreamProxy.Path().ClientID,
		UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		ProofConfirm:     connRes.Proof,
		ProofHeight:      connRes.ProofHeight,
		Signer:           signer.String(),
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
	if len(chanRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the channel(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}

	log.Printf("SyncChannelOpenInit: channel=%v proof=%v", chanRes.Channel, chanRes.Proof)

	proxyMsg := &proxytypes.MsgProxyChannelOpenTry{
		UpstreamClientId: ps.upstreamProxy.Path().ClientID,
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
	if len(chanRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the channel(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}

	log.Printf("SyncChannelOpenTry: channel=%v proof=%v", chanRes.Channel, chanRes.Proof)

	proxyMsg := &proxytypes.MsgProxyChannelOpenAck{
		UpstreamClientId:    ps.upstreamProxy.Path().ClientID,
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
	if len(chanRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the channel(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}

	log.Printf("SyncChannelOpenAck: channel=%v proof=%v", chanRes.Channel, chanRes.Proof)

	proxyMsg := &proxytypes.MsgProxyChannelOpenConfirm{
		UpstreamClientId:    ps.upstreamProxy.Path().ClientID,
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

func (ps ProxySynchronizer) SyncChannelOpenConfirm() error {
	provableHeight, err := ps.updateProxyUpstreamClient()
	if err != nil {
		return err
	}
	chanRes, err := ps.upstream.QueryChannelWithProof(provableHeight)
	if err != nil {
		return err
	}
	if len(chanRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the channel(height=%v)", provableHeight)
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}

	log.Printf("SyncChannelOpenConfirm: channel=%v proof=%v", chanRes.Channel, chanRes.Proof)

	proxyMsg := &proxytypes.MsgProxyChannelOpenFinalize{
		UpstreamClientId: ps.upstreamProxy.Path().ClientID,
		UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		PortId:           ps.path.PortID,
		ChannelId:        ps.path.ChannelID,
		ProofConfirm:     chanRes.Proof,
		ProofHeight:      chanRes.ProofHeight,
		Signer:           signer.String(),
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
	ack, err := ps.upstream.QueryPacketAcknowledgement(provableHeight, packet.Sequence)
	if err != nil {
		return err
	}
	signer, err := ps.upstreamProxy.GetAddress()
	if err != nil {
		return err
	}

	log.Printf("SyncRecvPacket: packet=%v ack=%v proof=%v", packet, res.Acknowledgement, res.Proof)

	proxyMsg := &proxytypes.MsgProxyAcknowledgePacket{
		UpstreamClientId: ps.upstreamProxy.Path().ClientID,
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

func makeMultiProof(
	cdc codec.Codec,
	head *multivtypes.Proof,
	branches []*multivtypes.Proof,
	leafClient *multivtypes.LeafProof,
) []byte {
	var mp multivtypes.MultiProof
	mp.Head = *head
	for _, branch := range branches {
		mp.Branches = append(mp.Branches, *branch)
	}
	mp.Leaf = *leafClient
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
