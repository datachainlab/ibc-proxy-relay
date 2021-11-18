package proxy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"
	proxytypes "github.com/datachainlab/ibc-proxy/modules/proxy/types"
	"github.com/hyperledger-labs/yui-relayer/core"
)

func UpdateUpstreamClient(upstream *core.ProvableChain) error {
	pr := upstream.ProverI.(*Prover)
	pr.xxxInitChains()

	// 1. update the lightDB corresponding to the proxy
	if _, _, _, err := pr.upstream.Proxy.UpdateLightWithHeader(); err != nil {
		return err
	}

	// 2. update the upstream client state on the proxy
	provableHeight, err := pr.updateProxyUpstreamClient()
	if err != nil {
		return err
	}

	// 3. make a msg and submit it to the proxy

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
	signer, err := pr.upstream.Proxy.GetAddress()
	if err != nil {
		return err
	}
	msg := &proxytypes.MsgProxyClientState{
		UpstreamClientId:     pr.upstream.UpstreamClientID,
		UpstreamPrefix:       commitmenttypes.NewMerklePrefix([]byte(host.StoreKey)),
		CounterpartyClientId: upstream.Path().ClientID,
		ClientState:          clientRes.GetClientState(),
		ConsensusState:       consensusRes.GetConsensusState(),
		ProofClient:          clientRes.Proof,
		ProofConsensus:       consensusRes.Proof,
		ProofHeight:          clientRes.ProofHeight,
		ConsensusHeight:      consensusHeight.(clienttypes.Height),
		Signer:               signer.String(),
	}
	if _, err := pr.upstream.Proxy.SendMsgs([]sdk.Msg{msg}); err != nil {
		return err
	}
	return nil
}
