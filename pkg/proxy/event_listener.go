package proxy

import (
	"fmt"
	"strings"
	"time"

	"github.com/avast/retry-go"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	"github.com/hyperledger-labs/yui-relayer/core"
)

type ProxyUpdater struct {
	synchronizer *ProxySynchronizer
}

var _ core.MsgEventListener = (*ProxyUpdater)(nil)

func NewProxyUpdater(synchronizer *ProxySynchronizer) ProxyUpdater {
	return ProxyUpdater{synchronizer: synchronizer}
}

// path: a path kept on upstream chain
// WARNING: This callback function calling may be skipped if the relayer is stopped for some reason.
func (up ProxyUpdater) OnSentMsg(msgs []sdk.Msg) error {
	for _, msg := range msgs {
		err := retry.Do(func() error {
			switch msg := msg.(type) {
			case *clienttypes.MsgCreateClient:
				return up.synchronizer.SyncCreateClient()
			case *clienttypes.MsgUpdateClient:
				// nop
			case *connectiontypes.MsgConnectionOpenInit:
				err := up.synchronizer.SyncConnectionOpenInit(msg.Counterparty)
				if !checkSkippableConnectionError(err) {
					return err
				}
			case *connectiontypes.MsgConnectionOpenTry:
				err := up.synchronizer.SyncConnectionOpenTry(msg.Counterparty)
				if !checkSkippableConnectionError(err) {
					return err
				}
			case *connectiontypes.MsgConnectionOpenAck:
				err := up.synchronizer.SyncConnectionOpenAck(msg.CounterpartyConnectionId)
				if !checkSkippableConnectionError(err) {
					return err
				}
			case *connectiontypes.MsgConnectionOpenConfirm:
				err := up.synchronizer.SyncConnectionOpenConfirm()
				if !checkSkippableConnectionError(err) {
					return err
				}
			case *channeltypes.MsgChannelOpenInit:
				return up.synchronizer.SyncChannelOpenInit()
			case *channeltypes.MsgChannelOpenTry:
				return up.synchronizer.SyncChannelOpenTry()
			case *channeltypes.MsgChannelOpenAck:
				return up.synchronizer.SyncChannelOpenAck()
			case *channeltypes.MsgChannelOpenConfirm:
				return up.synchronizer.SyncChannelOpenConfirm()
			case *channeltypes.MsgRecvPacket:
				return up.synchronizer.SyncRecvPacket(msg.Packet)
			case *channeltypes.MsgAcknowledgement:
				// nop
			}
			return nil
		},
			retry.Delay(1*time.Second),
			retry.Attempts(30),
		)
		if err != nil {
			panic(fmt.Errorf("failed to OnSentMsg: %T %v", msg, err))
		}
	}
	return nil
}

func checkSkippableConnectionError(err error) bool {
	if err == nil {
		return true
	} else if strings.Contains(err.Error(), "already exists") {
		return true
	} else {
		return false
	}
}
