package proxy

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	conntypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	chantypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"
	multivtypes "github.com/datachainlab/ibc-proxy/modules/light-clients/xx-multiv/types"
	proxytypes "github.com/datachainlab/ibc-proxy/modules/proxy/types"

	"github.com/hyperledger-labs/yui-relayer/core"
)

type Prover struct {
	chain           core.ChainI
	prover          core.ProverI
	upstreamProxy   *ProxyProvableChain
	downstreamProxy *ProxyProvableChain

	path *core.PathEnd

	proxySynchronizer *ProxySynchronizer
}

var (
	_ core.ProverI = (*Prover)(nil)
)

func NewProver(chain core.ChainI, prover core.ProverI, upstreamConfig *UpstreamConfig, downstreamConfig *DownstreamConfig) (*Prover, error) {
	if upstreamConfig == nil && downstreamConfig == nil {
		return nil, fmt.Errorf("either upstream or downstream must be not nil")
	} else if downstreamConfig != nil {
		prover = NewMultiVProver(prover)
	}
	return &Prover{
		chain:           chain,
		prover:          prover,
		upstreamProxy:   NewUpstreamProxy(upstreamConfig, chain),
		downstreamProxy: NewDownstreamProxy(downstreamConfig, chain),
	}, nil
}

// SetRelayInfo sets source's path and counterparty's info to the chain
func (pr *Prover) SetRelayInfo(path *core.PathEnd, counterparty *core.ProvableChain, counterpartyPath *core.PathEnd) error {
	downstreamProxyProver, ok := counterparty.ProverI.(*Prover)
	if !ok {
		return fmt.Errorf("counterparty's prover must be %T, but got %T", &Prover{}, counterparty.ProverI)
	}
	pr.path = path
	if pr.upstreamProxy != nil {
		pr.proxySynchronizer = NewProxySynchronizer(path, core.NewProvableChain(pr.chain, pr.prover), downstreamProxyProver.prover, pr.upstreamProxy, pr.downstreamProxy)
		pr.chain.RegisterMsgEventListener(NewProxyUpdater(pr.proxySynchronizer))
	}
	return nil
}

func (pr *Prover) SetupForRelay(ctx context.Context) error {
	// TODO sync with the upstream
	// return pr.proxySynchronizer.SyncALL()
	return nil
}

func (pr *Prover) GetUnderlyingProver() core.ProverI {
	switch prover := pr.prover.(type) {
	case *MultiVProver:
		return prover.ProverI
	default:
		return prover
	}
}

// GetChainID returns the chain ID
func (pr *Prover) GetChainID() string {
	return pr.chain.ChainID()
}

// QueryLatestHeader returns the latest header from the chain
func (pr *Prover) QueryLatestHeader() (out core.HeaderI, err error) {
	if pr.upstreamProxy != nil {
		return pr.upstreamProxy.QueryLatestHeader()
	} else {
		return pr.prover.QueryLatestHeader()
	}
}

// GetLatestLightHeight returns the latest height on the light client
func (pr *Prover) GetLatestLightHeight() (int64, error) {
	return pr.prover.GetLatestLightHeight()
}

// CreateMsgCreateClient creates a CreateClientMsg to this chain
func (pr *Prover) CreateMsgCreateClient(clientID string, dstHeader core.HeaderI, signer sdk.AccAddress) (*clienttypes.MsgCreateClient, error) {
	if pr.upstreamProxy != nil {
		return pr.upstreamProxy.CreateMsgCreateClient(clientID, dstHeader, signer)
	} else {
		return pr.prover.CreateMsgCreateClient(clientID, dstHeader, signer)
	}
}

// SetupHeader creates a new header based on a given header
func (pr *Prover) SetupHeader(dst core.LightClientIBCQueryierI, baseSrcHeader core.HeaderI) (core.HeaderI, error) {
	if pr.upstreamProxy != nil {
		return pr.upstreamProxy.SetupHeader(dst, baseSrcHeader)
	} else {
		return pr.prover.SetupHeader(dst, baseSrcHeader)
	}
}

// UpdateLightWithHeader updates a header on the light client and returns the header and height corresponding to the chain
func (pr *Prover) UpdateLightWithHeader() (header core.HeaderI, provableHeight int64, queryableHeight int64, err error) {
	if pr.upstreamProxy != nil {
		_, _, qh, err := pr.prover.UpdateLightWithHeader()
		if err != nil {
			return nil, 0, 0, err
		}
		h, ph, _, err := pr.upstreamProxy.UpdateLightWithHeader()
		if err != nil {
			return nil, 0, 0, err
		}
		return h, ph, qh, nil
	} else {
		return pr.prover.UpdateLightWithHeader()
	}
}

// QueryClientConsensusState returns the ClientConsensusState and its proof
func (pr *Prover) QueryClientConsensusStateWithProof(height int64, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
	if pr.upstreamProxy != nil {
		return pr.upstreamProxy.QueryProxyClientConsensusStateWithProof(height, dstClientConsHeight)
	} else {
		// in case we are downstream:

		clientRes, err := pr.prover.QueryClientStateWithProof(height)
		if err != nil {
			return nil, err
		}
		proxyClientState, err := clienttypes.UnpackClientState(clientRes.ClientState)
		if err != nil {
			return nil, err
		}
		consRes, err := pr.prover.QueryClientConsensusStateWithProof(height, proxyClientState.GetLatestHeight())
		if err != nil {
			return nil, err
		}
		head := &multivtypes.Proof{
			ClientProof:     clientRes.Proof,
			ClientState:     clientRes.ClientState,
			ConsensusProof:  consRes.Proof,
			ConsensusState:  consRes.ConsensusState,
			ProofHeight:     clientRes.ProofHeight,
			ConsensusHeight: proxyClientState.GetLatestHeight().(clienttypes.Height),
		}
		proxyConsRes, err := pr.downstreamProxy.QueryClientConsensusStateWithProof(
			int64(proxyClientState.GetLatestHeight().GetRevisionHeight()-1),
			dstClientConsHeight,
		)
		if err != nil {
			return nil, err
		}
		leafClient := &multivtypes.LeafProof{
			Proof:       proxyConsRes.Proof,
			ProofHeight: proxyConsRes.ProofHeight,
		}
		proof := makeMultiProof(pr.chain.Codec(), head, nil, leafClient)
		return clienttypes.NewQueryConsensusStateResponse(proxyConsRes.ConsensusState, proof, dstClientConsHeight.(clienttypes.Height)), nil
	}
}

// QueryClientStateWithProof returns the ClientState and its proof
func (pr *Prover) QueryClientStateWithProof(height int64) (*clienttypes.QueryClientStateResponse, error) {
	if pr.upstreamProxy != nil {
		return pr.upstreamProxy.QueryProxyClientStateWithProof(height)
	} else {
		// in case we are downstream:

		// down-(head)>proxy-(leaf)>upstream

		clientRes, err := pr.prover.QueryClientStateWithProof(height)
		if err != nil {
			return nil, err
		}
		proxyClientState, err := clienttypes.UnpackClientState(clientRes.ClientState)
		if err != nil {
			return nil, err
		}
		consRes, err := pr.prover.QueryClientConsensusStateWithProof(height, proxyClientState.GetLatestHeight())
		if err != nil {
			return nil, err
		}
		head := &multivtypes.Proof{
			ClientProof:     clientRes.Proof,
			ClientState:     clientRes.ClientState,
			ConsensusProof:  consRes.Proof,
			ConsensusState:  consRes.ConsensusState,
			ProofHeight:     clientRes.ProofHeight,
			ConsensusHeight: proxyClientState.GetLatestHeight().(clienttypes.Height),
		}
		// TODO (height-1) is for tendermint specific
		proxyClientRes, err := pr.downstreamProxy.QueryClientStateWithProof(int64(proxyClientState.GetLatestHeight().GetRevisionHeight() - 1))
		if err != nil {
			return nil, err
		}
		leafClient := &multivtypes.LeafProof{
			Proof:       proxyClientRes.Proof,
			ProofHeight: proxyClientRes.ProofHeight,
		}
		proof := makeMultiProof(pr.chain.Codec(), head, nil, leafClient)
		return clienttypes.NewQueryClientStateResponse(proxyClientRes.ClientState, proof, clientRes.ProofHeight), nil
	}
}

// QueryConnectionWithProof returns the Connection and its proof
func (pr *Prover) QueryConnectionWithProof(height int64) (*conntypes.QueryConnectionResponse, error) {
	if pr.upstreamProxy != nil {
		return pr.upstreamProxy.QueryProxyConnectionStateWithProof(height)
	} else {
		return pr.prover.QueryConnectionWithProof(height)
	}
}

// QueryChannelWithProof returns the Channel and its proof
func (pr *Prover) QueryChannelWithProof(height int64) (chanRes *chantypes.QueryChannelResponse, err error) {
	if pr.upstreamProxy != nil {
		return pr.upstreamProxy.QueryProxyChannelWithProof(height)
	} else {
		return pr.prover.QueryChannelWithProof(height)
	}
}

// QueryPacketCommitmentWithProof returns the packet commitment and its proof
func (pr *Prover) QueryPacketCommitmentWithProof(height int64, seq uint64) (comRes *chantypes.QueryPacketCommitmentResponse, err error) {
	if pr.upstreamProxy != nil {
		// first, query a packet commitment to the proxy
		// if the commitment exists, just returns it
		// otherwise, the proxy proxies the commitment
		res, err := pr.upstreamProxy.QueryProxyPacketCommitmentWithProof(height, seq)
		if err == nil {
			return res, nil
		} else if !strings.Contains(err.Error(), "packet commitment not found") {
			return nil, err
		}
		log.Println("try to perform a relay `upstream->proxy`")
		provableHeight, err := pr.proxySynchronizer.updateProxyUpstreamClient()
		if err != nil {
			return nil, err
		}
		pcRes, err := pr.prover.QueryPacketCommitmentWithProof(provableHeight, seq)
		if err != nil {
			return nil, err
		}
		packet, err := pr.chain.QueryPacket(provableHeight, seq)
		if err != nil {
			return nil, err
		}
		signer, err := pr.upstreamProxy.GetAddress()
		if err != nil {
			return nil, err
		}
		proxyMsg := &proxytypes.MsgProxyRecvPacket{
			UpstreamClientId: pr.upstreamProxy.Path().ClientID,
			UpstreamPrefix:   commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
			Packet:           *packet,
			Proof:            pcRes.Proof,
			ProofHeight:      pcRes.ProofHeight,
			Signer:           signer.String(),
		}
		if _, err := pr.upstreamProxy.SendMsgs([]sdk.Msg{proxyMsg}); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("the relay `upstream->proxy` succeeded, but you need to update the proxy client in the downstream and call `QueryPacketCommitmentWithProof` again")
	} else {
		return pr.prover.QueryPacketCommitmentWithProof(height, seq)
	}
}

// QueryPacketAcknowledgementCommitmentWithProof returns the packet acknowledgement commitment and its proof
func (pr *Prover) QueryPacketAcknowledgementCommitmentWithProof(height int64, seq uint64) (ackRes *chantypes.QueryPacketAcknowledgementResponse, err error) {
	if pr.upstreamProxy != nil {
		return pr.upstreamProxy.QueryProxyPacketAcknowledgementCommitmentWithProof(height, seq)
	} else {
		return pr.prover.QueryPacketAcknowledgementCommitmentWithProof(height, seq)
	}
}

// Init ...
func (pr *Prover) Init(homePath string, timeout time.Duration, codec codec.ProtoCodecMarshaler, debug bool) error {
	if pr.upstreamProxy != nil {
		if err := pr.upstreamProxy.Init(homePath, timeout, codec, debug); err != nil {
			return err
		}
	}
	if pr.downstreamProxy != nil {
		if err := pr.downstreamProxy.Init(homePath, timeout, codec, debug); err != nil {
			return err
		}
	}
	return nil
}

type MultiVProver struct {
	core.ProverI
}

var _ core.ProverI = (*MultiVProver)(nil)

func NewMultiVProver(prover core.ProverI) *MultiVProver {
	return &MultiVProver{ProverI: prover}
}

func (p *MultiVProver) CreateMsgCreateClient(clientID string, dstHeader core.HeaderI, signer sdk.AccAddress) (*clienttypes.MsgCreateClient, error) {
	msg, err := p.ProverI.CreateMsgCreateClient(clientID, dstHeader, signer)
	if err != nil {
		return nil, err
	}
	clientState := &multivtypes.ClientState{UnderlyingClientState: msg.ClientState, Depth: 0}
	anyClientState, err := clienttypes.PackClientState(clientState)
	if err != nil {
		return nil, err
	}
	msg.ClientState = anyClientState
	return msg, nil
}
