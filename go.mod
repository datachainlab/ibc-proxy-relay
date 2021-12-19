module github.com/datachainlab/ibc-proxy-relay

go 1.16

replace (
	github.com/cosmos/ibc-go => github.com/datachainlab/ibc-go v0.0.0-20210623043207-6582d8c965f8
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)

require (
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/cosmos/cosmos-sdk v0.43.0-beta1
	github.com/cosmos/ibc-go v1.0.0-beta1
	github.com/datachainlab/ibc-proxy v0.0.0-20211216053654-43b54f292c56
	github.com/gogo/protobuf v1.3.3
	github.com/hyperledger-labs/yui-relayer v0.1.1-0.20211206093912-53858db80d32
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/tendermint/tendermint v0.34.10
	google.golang.org/grpc v1.37.0
)
