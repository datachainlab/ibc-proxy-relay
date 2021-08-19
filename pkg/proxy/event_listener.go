package proxy

import (
	"fmt"
	"log"
	"time"

	"github.com/avast/retry-go"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	conntypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	chantypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"
	proxytypes "github.com/datachainlab/ibc-proxy/modules/proxy/types"
	"github.com/hyperledger-labs/yui-relayer/core"
)

var (
	ConnectionVersion         = conntypes.ExportedVersionsToProto(conntypes.GetCompatibleVersions())[0]
	DefaultDelayPeriod uint64 = 0
)

// path: a path kept on upstream chain
// WARNING: This callback function calling may be skipped if the relayer is stopped for some reason.
func (pr *Prover) OnSentMsg(path *core.PathEnd, msgs []sdk.Msg) error {
	for _, msg := range msgs {
		log.Printf("Called OnSentMsg: %T", msg)
		var err error
		switch msg := msg.(type) {
		case *clienttypes.MsgCreateClient:
			err = pr.onCreateClient(path, msg)
		case *clienttypes.MsgUpdateClient:
			// nop
		case *conntypes.MsgConnectionOpenInit:
			err = pr.onConnectionOpenInit(path, msg)
		case *conntypes.MsgConnectionOpenTry:
			err = pr.onConnectionOpenTry(path, msg)
		case *conntypes.MsgConnectionOpenAck:
			err = pr.onConnectionOpenAck(path, msg)
		case *conntypes.MsgConnectionOpenConfirm:
			// nop
		case *chantypes.MsgChannelOpenInit:
			panic("not implemented error")
		case *chantypes.MsgChannelOpenTry:
			err = pr.onChannelOpenTry(path, msg)
		case *chantypes.MsgChannelOpenAck:
			panic("not implemented error")
		case *chantypes.MsgChannelOpenConfirm:
			// nop
		case *chantypes.MsgRecvPacket:
			err = pr.onRecvPacket(path, msg)
		case *chantypes.MsgAcknowledgement:
			panic("not implemented error")
		}
		if err != nil {
			log.Println("OnSentMsg:", err)
			return err
		}
	}
	return nil
}

func (pr *Prover) onCreateClient(path *core.PathEnd, msg *clienttypes.MsgCreateClient) error {

	// 1. creates an upstream client on the proxy

	header, err := pr.prover.QueryLatestHeader()
	if err != nil {
		return err
	}
	signer, err := pr.upstream.Proxy.GetAddress()
	if err != nil {
		return err
	}
	proxyMsg, err := pr.prover.CreateMsgCreateClient(pr.upstream.Proxy.ProxyPath().UpstreamClientID, header, signer)
	if err != nil {
		return err
	}
	if _, err := pr.upstream.Proxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return err
	}

	// 2. proxies the downstream client in the upstream

	// TODO Is it necessary to implement proxyClientState function in the proxy module?

	return nil
}

func (pr *Prover) onConnectionOpenInit(path *core.PathEnd, msg *conntypes.MsgConnectionOpenInit) error {
	provableHeight, err := pr.updateProxyUpstreamClient()
	if err != nil {
		return err
	}
	clientRes, err := pr.prover.QueryClientStateWithProof(provableHeight)
	if err != nil {
		return err
	}

	clientState := clientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
	consensusHeight := clientState.GetLatestHeight()
	consensusRes, err := pr.prover.QueryClientConsensusStateWithProof(provableHeight, consensusHeight)
	if err != nil {
		return err
	}
	connRes, err := pr.prover.QueryConnectionWithProof(provableHeight)
	if err != nil {
		return err
	}
	log.Println("onConnectionOpenInit:", connRes.Proof, connRes.Connection)
	if len(connRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the connection(height=%v)", provableHeight)
	}

	signer, err := pr.upstream.Proxy.GetAddress()
	if err != nil {
		return err
	}
	proxyMsg := &proxytypes.MsgProxyConnectionOpenTry{
		ConnectionId:     path.ConnectionID,
		UpstreamClientId: pr.upstream.UpstreamClientID,
		UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		Connection: conntypes.NewConnectionEnd(
			conntypes.INIT,
			path.ClientID,
			msg.Counterparty,
			[]*conntypes.Version{ConnectionVersion}, DefaultDelayPeriod,
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
	if _, err := pr.upstream.Proxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return err
	}
	return nil
}

func (pr *Prover) onConnectionOpenTry(path *core.PathEnd, msg *conntypes.MsgConnectionOpenTry) error {
	provableHeight, err := pr.updateProxyUpstreamClient()
	if err != nil {
		return err
	}
	clientRes, err := pr.prover.QueryClientStateWithProof(provableHeight)
	if err != nil {
		return err
	}

	clientState := clientRes.ClientState.GetCachedValue().(ibcexported.ClientState)
	consensusHeight := clientState.GetLatestHeight()
	consensusRes, err := pr.prover.QueryClientConsensusStateWithProof(provableHeight, consensusHeight)
	if err != nil {
		return err
	}
	// TODO it should retry to query until a new block is created in the upstream
	connRes, err := pr.prover.QueryConnectionWithProof(provableHeight)
	if err != nil {
		return err
	}
	log.Println("onConnectionOpenTry:", connRes.Proof, connRes.Connection)
	if len(connRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the connection(height=%v)", provableHeight)
	}

	signer, err := pr.upstream.Proxy.GetAddress()
	if err != nil {
		return err
	}
	proxyMsg := &proxytypes.MsgProxyConnectionOpenAck{
		ConnectionId:     path.ConnectionID,
		UpstreamClientId: pr.upstream.UpstreamClientID,
		UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		Connection: conntypes.NewConnectionEnd(
			conntypes.TRYOPEN,
			path.ClientID,
			msg.Counterparty,
			[]*conntypes.Version{ConnectionVersion}, DefaultDelayPeriod,
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
	if _, err := pr.upstream.Proxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return err
	}
	return nil
}

func (pr *Prover) onConnectionOpenAck(path *core.PathEnd, msg *conntypes.MsgConnectionOpenAck) error {
	provableHeight, err := pr.updateProxyUpstreamClient()
	if err != nil {
		return err
	}
	connRes, err := pr.prover.QueryConnectionWithProof(provableHeight)
	if err != nil {
		return err
	}
	log.Println("onConnectionOpenAck:", connRes.Proof, connRes.Connection)
	if len(connRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the connection(height=%v)", provableHeight)
	}
	signer, err := pr.upstream.Proxy.GetAddress()
	if err != nil {
		return err
	}
	proxyMsg := &proxytypes.MsgProxyConnectionOpenConfirm{
		ConnectionId:     path.ConnectionID,
		UpstreamClientId: pr.upstream.UpstreamClientID,
		UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		Connection: conntypes.NewConnectionEnd(
			conntypes.OPEN,
			path.ClientID,
			conntypes.Counterparty{
				ClientId:     "", // TODO query the previous connection state and set this field to the value
				ConnectionId: msg.CounterpartyConnectionId,
				Prefix:       commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			},
			[]*conntypes.Version{ConnectionVersion}, DefaultDelayPeriod,
		),
		ProofAck:    connRes.Proof,
		ProofHeight: connRes.ProofHeight,
		Signer:      signer.String(),
	}
	if _, err := pr.upstream.Proxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return err
	}
	return nil
}

func (pr *Prover) onChannelOpenTry(path *core.PathEnd, msg *chantypes.MsgChannelOpenTry) error {
	provableHeight, err := pr.updateProxyUpstreamClient()
	if err != nil {
		return err
	}
	chanRes, err := pr.prover.QueryChannelWithProof(provableHeight)
	if err != nil {
		return err
	}
	log.Println("onChannelOpenTry:", chanRes.Proof, chanRes.Channel)
	if len(chanRes.Proof) == 0 {
		return fmt.Errorf("failed to query a proof of the channel(height=%v)", provableHeight)
	}
	signer, err := pr.upstream.Proxy.GetAddress()
	if err != nil {
		return err
	}
	proxyMsg := &proxytypes.MsgProxyChannelOpenAck{
		UpstreamClientId:    pr.upstream.UpstreamClientID,
		UpstreamPrefix:      commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		Order:               chanRes.Channel.Ordering,
		ConnectionHops:      chanRes.Channel.ConnectionHops,
		PortId:              path.PortID,
		ChannelId:           path.ChannelID,
		Counterparty:        chanRes.Channel.Counterparty,
		Version:             chanRes.Channel.Version,
		CounterpartyVersion: chanRes.Channel.Version, // TODO this field must be the version that is provided by counterparty chain
		ProofTry:            chanRes.Proof,
		ProofHeight:         chanRes.ProofHeight,
		Signer:              signer.String(),
	}
	if _, err := pr.upstream.Proxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return err
	}
	return nil
}

func (pr *Prover) onRecvPacket(path *core.PathEnd, msg *chantypes.MsgRecvPacket) error {
	var (
		provableHeight int64
		res            *chantypes.QueryPacketAcknowledgementResponse
	)

	err := retry.Do(
		func() error {
			var err error
			provableHeight, err = pr.updateProxyUpstreamClient()
			if err != nil {
				return err
			}
			res, err = pr.prover.QueryPacketAcknowledgementCommitmentWithProof(provableHeight, msg.Packet.Sequence)
			return err
		},
		retry.Delay(1*time.Second),
		retry.Attempts(30),
	)
	if err != nil {
		return err
	}
	log.Println("onRecvPacket:", res.Proof, res.Acknowledgement)

	ack, err := pr.chain.QueryPacketAcknowledgement(provableHeight, msg.Packet.Sequence)
	if err != nil {
		return err
	}
	signer, err := pr.upstream.Proxy.GetAddress()
	if err != nil {
		return err
	}
	proxyMsg := &proxytypes.MsgProxyAcknowledgePacket{
		UpstreamClientId: pr.upstream.UpstreamClientID,
		UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		Packet:           msg.Packet,
		Acknowledgement:  ack,
		Proof:            res.Proof,
		ProofHeight:      res.ProofHeight,
		Signer:           signer.String(),
	}
	if _, err := pr.upstream.Proxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
		return err
	}
	return nil
}

// updateProxyUpstreamClient updates the upstream client on the proxy
func (pr *Prover) updateProxyUpstreamClient() (int64, error) {
	h, provableHeight, _, err := pr.prover.UpdateLightWithHeader()
	if err != nil {
		return 0, err
	}
	header, err := pr.prover.SetupHeader(pr.upstream.Proxy, h)
	if err != nil {
		return 0, err
	}
	addr, err := pr.upstream.Proxy.GetAddress()
	if err != nil {
		return 0, err
	}
	_, err = pr.upstream.Proxy.SendMsgs(
		[]sdk.Msg{pr.upstream.Proxy.Path().UpdateClient(header, addr)},
	)
	if err != nil {
		return 0, err
	}
	return provableHeight, nil
}
