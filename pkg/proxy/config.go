package proxy

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger-labs/yui-relayer/core"
)

type ProxyChainConfigI interface {
	proto.Message
	Build() (ProxyChainI, error)
}

type ProxyChainProverConfigI interface {
	proto.Message
	Build(ProxyChainI) (ProxyChainProverI, error)
}

var _ core.ProverConfigI = (*ProverConfig)(nil)

func (pc ProverConfig) Build(chain core.ChainI) (core.ProverI, error) {
	prover, err := pc.Prover.GetCachedValue().(core.ProverConfigI).Build(chain)
	if err != nil {
		return nil, err
	}
	return NewProver(chain, prover, pc.Upstream, pc.Downstream)
}

type UpstreamProxy struct {
	ProxyProvableChain
	UpstreamClientID string
}

func NewUpstreamProxy(config *UpstreamConfig, chain core.ChainI) *UpstreamProxy {
	if config == nil {
		return nil
	}
	proxyChain, err := config.ProxyChain.GetCachedValue().(ProxyChainConfigI).Build()
	if err != nil {
		panic(err)
	}
	proxyChain.SetProxyPath(ProxyPath{
		UpstreamClientID: config.UpstreamClientId,
		UpstreamChain:    chain,
	})
	proxyChainProver, err := config.ProxyChainProver.GetCachedValue().(ProxyChainProverConfigI).Build(proxyChain)
	if err != nil {
		panic(err)
	}
	return &UpstreamProxy{
		ProxyProvableChain: ProxyProvableChain{ProxyChainI: proxyChain, ProxyChainProverI: proxyChainProver},
		UpstreamClientID:   config.UpstreamClientId,
	}
}

type DownstreamProxy struct {
	ProxyProvableChain
	UpstreamClientID string
}

func NewDownstreamProxy(config *DownstreamConfig, chain core.ChainI) *DownstreamProxy {
	if config == nil {
		return nil
	}
	proxyChain, err := config.ProxyChain.GetCachedValue().(ProxyChainConfigI).Build()
	if err != nil {
		panic(err)
	}
	proxyChainProver, err := config.ProxyChainProver.GetCachedValue().(ProxyChainProverConfigI).Build(proxyChain)
	if err != nil {
		panic(err)
	}
	return &DownstreamProxy{
		ProxyProvableChain: ProxyProvableChain{ProxyChainI: proxyChain, ProxyChainProverI: proxyChainProver},
		UpstreamClientID:   config.UpstreamClientId,
	}
}

var _, _, _ codectypes.UnpackInterfacesMessage = (*ProverConfig)(nil), (*UpstreamConfig)(nil), (*DownstreamConfig)(nil)

func (cfg *ProverConfig) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if cfg == nil {
		return nil
	}
	if err := unpacker.UnpackAny(cfg.Prover, new(core.ProverConfigI)); err != nil {
		return err
	}
	if err := cfg.Upstream.UnpackInterfaces(unpacker); err != nil {
		return err
	}
	if err := cfg.Downstream.UnpackInterfaces(unpacker); err != nil {
		return err
	}
	return nil
}

func (cfg *UpstreamConfig) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if cfg == nil {
		return nil
	}
	if err := unpacker.UnpackAny(cfg.ProxyChain, new(ProxyChainConfigI)); err != nil {
		return err
	}
	if err := unpacker.UnpackAny(cfg.ProxyChainProver, new(ProxyChainProverConfigI)); err != nil {
		return err
	}
	return nil
}

func (cfg *DownstreamConfig) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if cfg == nil {
		return nil
	}
	if err := unpacker.UnpackAny(cfg.ProxyChain, new(ProxyChainConfigI)); err != nil {
		return err
	}
	if err := unpacker.UnpackAny(cfg.ProxyChainProver, new(ProxyChainProverConfigI)); err != nil {
		return err
	}
	return nil
}
