// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proxy/prover.proto

package proxy

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ProverConfig struct {
	Prover     *types.Any        `protobuf:"bytes,1,opt,name=prover,proto3" json:"prover,omitempty"`
	Upstream   *UpstreamConfig   `protobuf:"bytes,2,opt,name=upstream,proto3" json:"upstream,omitempty"`
	Downstream *DownstreamConfig `protobuf:"bytes,3,opt,name=downstream,proto3" json:"downstream,omitempty"`
}

func (m *ProverConfig) Reset()         { *m = ProverConfig{} }
func (m *ProverConfig) String() string { return proto.CompactTextString(m) }
func (*ProverConfig) ProtoMessage()    {}
func (*ProverConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_362c6678412ad31f, []int{0}
}
func (m *ProverConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProverConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProverConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProverConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProverConfig.Merge(m, src)
}
func (m *ProverConfig) XXX_Size() int {
	return m.Size()
}
func (m *ProverConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ProverConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ProverConfig proto.InternalMessageInfo

func (m *ProverConfig) GetProver() *types.Any {
	if m != nil {
		return m.Prover
	}
	return nil
}

func (m *ProverConfig) GetUpstream() *UpstreamConfig {
	if m != nil {
		return m.Upstream
	}
	return nil
}

func (m *ProverConfig) GetDownstream() *DownstreamConfig {
	if m != nil {
		return m.Downstream
	}
	return nil
}

type UpstreamConfig struct {
	ProxyChain       *types.Any `protobuf:"bytes,1,opt,name=proxy_chain,json=proxyChain,proto3" json:"proxy_chain,omitempty"`
	ProxyChainProver *types.Any `protobuf:"bytes,2,opt,name=proxy_chain_prover,json=proxyChainProver,proto3" json:"proxy_chain_prover,omitempty"`
	// the client id corresponding to the upstream on the proxy machine
	// TODO this parameter should be moved into a path configuration
	UpstreamClientId string `protobuf:"bytes,3,opt,name=upstream_client_id,json=upstreamClientId,proto3" json:"upstream_client_id,omitempty"`
}

func (m *UpstreamConfig) Reset()         { *m = UpstreamConfig{} }
func (m *UpstreamConfig) String() string { return proto.CompactTextString(m) }
func (*UpstreamConfig) ProtoMessage()    {}
func (*UpstreamConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_362c6678412ad31f, []int{1}
}
func (m *UpstreamConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UpstreamConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UpstreamConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UpstreamConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpstreamConfig.Merge(m, src)
}
func (m *UpstreamConfig) XXX_Size() int {
	return m.Size()
}
func (m *UpstreamConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_UpstreamConfig.DiscardUnknown(m)
}

var xxx_messageInfo_UpstreamConfig proto.InternalMessageInfo

func (m *UpstreamConfig) GetProxyChain() *types.Any {
	if m != nil {
		return m.ProxyChain
	}
	return nil
}

func (m *UpstreamConfig) GetProxyChainProver() *types.Any {
	if m != nil {
		return m.ProxyChainProver
	}
	return nil
}

func (m *UpstreamConfig) GetUpstreamClientId() string {
	if m != nil {
		return m.UpstreamClientId
	}
	return ""
}

type DownstreamConfig struct {
	ProxyChain       *types.Any `protobuf:"bytes,1,opt,name=proxy_chain,json=proxyChain,proto3" json:"proxy_chain,omitempty"`
	ProxyChainProver *types.Any `protobuf:"bytes,2,opt,name=proxy_chain_prover,json=proxyChainProver,proto3" json:"proxy_chain_prover,omitempty"`
	// the client id corresponding to the upstream on the proxy machine
	// TODO this parameter should be moved into a path configuration
	UpstreamClientId string `protobuf:"bytes,3,opt,name=upstream_client_id,json=upstreamClientId,proto3" json:"upstream_client_id,omitempty"`
}

func (m *DownstreamConfig) Reset()         { *m = DownstreamConfig{} }
func (m *DownstreamConfig) String() string { return proto.CompactTextString(m) }
func (*DownstreamConfig) ProtoMessage()    {}
func (*DownstreamConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_362c6678412ad31f, []int{2}
}
func (m *DownstreamConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DownstreamConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DownstreamConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DownstreamConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DownstreamConfig.Merge(m, src)
}
func (m *DownstreamConfig) XXX_Size() int {
	return m.Size()
}
func (m *DownstreamConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_DownstreamConfig.DiscardUnknown(m)
}

var xxx_messageInfo_DownstreamConfig proto.InternalMessageInfo

func (m *DownstreamConfig) GetProxyChain() *types.Any {
	if m != nil {
		return m.ProxyChain
	}
	return nil
}

func (m *DownstreamConfig) GetProxyChainProver() *types.Any {
	if m != nil {
		return m.ProxyChainProver
	}
	return nil
}

func (m *DownstreamConfig) GetUpstreamClientId() string {
	if m != nil {
		return m.UpstreamClientId
	}
	return ""
}

func init() {
	proto.RegisterType((*ProverConfig)(nil), "ibc.proxy.prover.v1.ProverConfig")
	proto.RegisterType((*UpstreamConfig)(nil), "ibc.proxy.prover.v1.UpstreamConfig")
	proto.RegisterType((*DownstreamConfig)(nil), "ibc.proxy.prover.v1.DownstreamConfig")
}

func init() { proto.RegisterFile("proxy/prover.proto", fileDescriptor_362c6678412ad31f) }

var fileDescriptor_362c6678412ad31f = []byte{
	// 341 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x90, 0xbf, 0x4e, 0x42, 0x31,
	0x14, 0xc6, 0x29, 0x26, 0x44, 0x8b, 0x31, 0xa4, 0x32, 0x20, 0xc3, 0x8d, 0xc1, 0x98, 0x38, 0x40,
	0x1b, 0x34, 0xce, 0x46, 0xd0, 0xc1, 0xb8, 0x18, 0x12, 0x17, 0x97, 0x9b, 0xde, 0x3f, 0x94, 0xc6,
	0x4b, 0x7b, 0x53, 0x2f, 0xe8, 0x7d, 0x0b, 0x5f, 0xc7, 0xb8, 0x3a, 0x38, 0x32, 0x3a, 0x1a, 0x78,
	0x11, 0xc3, 0x69, 0xc1, 0x3f, 0x21, 0xec, 0x6e, 0xed, 0x39, 0xdf, 0xef, 0x3b, 0xe7, 0x3b, 0x98,
	0xa4, 0x46, 0x3f, 0xe5, 0x2c, 0x35, 0x7a, 0x1c, 0x1b, 0x9a, 0x1a, 0x9d, 0x69, 0xb2, 0x2b, 0x83,
	0x90, 0x42, 0x9d, 0xba, 0xfa, 0xb8, 0x5d, 0xaf, 0x0a, 0x2d, 0x34, 0xf4, 0xd9, 0xfc, 0x65, 0xa5,
	0xf5, 0x3d, 0xa1, 0xb5, 0x48, 0x62, 0x06, 0xbf, 0x60, 0xd4, 0x67, 0x5c, 0xe5, 0xb6, 0xd5, 0x78,
	0x43, 0x78, 0xfb, 0x06, 0xf0, 0xae, 0x56, 0x7d, 0x29, 0x48, 0x13, 0x97, 0xac, 0x5d, 0x0d, 0xed,
	0xa3, 0xa3, 0xf2, 0x71, 0x95, 0x5a, 0x98, 0x2e, 0x60, 0x7a, 0xae, 0xf2, 0x9e, 0xd3, 0x90, 0x33,
	0xbc, 0x39, 0x4a, 0x1f, 0x32, 0x13, 0xf3, 0x61, 0xad, 0x08, 0xfa, 0x03, 0xba, 0x62, 0x2f, 0x7a,
	0xeb, 0x44, 0x76, 0x48, 0x6f, 0x09, 0x91, 0x4b, 0x8c, 0x23, 0xfd, 0xa8, 0x9c, 0xc5, 0x06, 0x58,
	0x1c, 0xae, 0xb4, 0xb8, 0x58, 0xca, 0x9c, 0xc9, 0x0f, 0xb0, 0xf1, 0x82, 0xf0, 0xce, 0xef, 0x19,
	0xe4, 0x14, 0x97, 0xc1, 0xc2, 0x0f, 0x07, 0x5c, 0xaa, 0xb5, 0x69, 0x30, 0x08, 0xbb, 0x73, 0x1d,
	0xe9, 0xb8, 0x63, 0x5b, 0xcc, 0x77, 0xb7, 0x28, 0xae, 0xa1, 0x2b, 0xdf, 0xb4, 0xbd, 0x24, 0x69,
	0x62, 0xb2, 0x08, 0xe8, 0x87, 0x89, 0x8c, 0x55, 0xe6, 0xcb, 0x08, 0xc2, 0x6d, 0xf5, 0x2a, 0x8b,
	0x4e, 0x17, 0x1a, 0x57, 0x51, 0xe3, 0x15, 0xe1, 0xca, 0xdf, 0x70, 0xff, 0x66, 0xfb, 0xce, 0xf5,
	0xfb, 0xd4, 0x43, 0x93, 0xa9, 0x87, 0x3e, 0xa7, 0x1e, 0x7a, 0x9e, 0x79, 0x85, 0xc9, 0xcc, 0x2b,
	0x7c, 0xcc, 0xbc, 0xc2, 0x5d, 0x5b, 0xc8, 0x6c, 0x30, 0x0a, 0x68, 0xa8, 0x87, 0x2c, 0xe2, 0x19,
	0x87, 0x95, 0x12, 0x1e, 0x30, 0x19, 0x84, 0x2d, 0x98, 0xda, 0x32, 0x71, 0xc2, 0x73, 0x96, 0xde,
	0x0b, 0x06, 0xff, 0xa0, 0x04, 0xab, 0x9d, 0x7c, 0x05, 0x00, 0x00, 0xff, 0xff, 0xd3, 0x63, 0xb0,
	0x10, 0xf0, 0x02, 0x00, 0x00,
}

func (m *ProverConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProverConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProverConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Downstream != nil {
		{
			size, err := m.Downstream.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProver(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Upstream != nil {
		{
			size, err := m.Upstream.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProver(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Prover != nil {
		{
			size, err := m.Prover.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProver(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *UpstreamConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpstreamConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UpstreamConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.UpstreamClientId) > 0 {
		i -= len(m.UpstreamClientId)
		copy(dAtA[i:], m.UpstreamClientId)
		i = encodeVarintProver(dAtA, i, uint64(len(m.UpstreamClientId)))
		i--
		dAtA[i] = 0x1a
	}
	if m.ProxyChainProver != nil {
		{
			size, err := m.ProxyChainProver.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProver(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.ProxyChain != nil {
		{
			size, err := m.ProxyChain.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProver(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DownstreamConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DownstreamConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DownstreamConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.UpstreamClientId) > 0 {
		i -= len(m.UpstreamClientId)
		copy(dAtA[i:], m.UpstreamClientId)
		i = encodeVarintProver(dAtA, i, uint64(len(m.UpstreamClientId)))
		i--
		dAtA[i] = 0x1a
	}
	if m.ProxyChainProver != nil {
		{
			size, err := m.ProxyChainProver.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProver(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.ProxyChain != nil {
		{
			size, err := m.ProxyChain.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProver(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintProver(dAtA []byte, offset int, v uint64) int {
	offset -= sovProver(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ProverConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Prover != nil {
		l = m.Prover.Size()
		n += 1 + l + sovProver(uint64(l))
	}
	if m.Upstream != nil {
		l = m.Upstream.Size()
		n += 1 + l + sovProver(uint64(l))
	}
	if m.Downstream != nil {
		l = m.Downstream.Size()
		n += 1 + l + sovProver(uint64(l))
	}
	return n
}

func (m *UpstreamConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ProxyChain != nil {
		l = m.ProxyChain.Size()
		n += 1 + l + sovProver(uint64(l))
	}
	if m.ProxyChainProver != nil {
		l = m.ProxyChainProver.Size()
		n += 1 + l + sovProver(uint64(l))
	}
	l = len(m.UpstreamClientId)
	if l > 0 {
		n += 1 + l + sovProver(uint64(l))
	}
	return n
}

func (m *DownstreamConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ProxyChain != nil {
		l = m.ProxyChain.Size()
		n += 1 + l + sovProver(uint64(l))
	}
	if m.ProxyChainProver != nil {
		l = m.ProxyChainProver.Size()
		n += 1 + l + sovProver(uint64(l))
	}
	l = len(m.UpstreamClientId)
	if l > 0 {
		n += 1 + l + sovProver(uint64(l))
	}
	return n
}

func sovProver(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProver(x uint64) (n int) {
	return sovProver(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ProverConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProver
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProverConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProverConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prover", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Prover == nil {
				m.Prover = &types.Any{}
			}
			if err := m.Prover.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Upstream", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Upstream == nil {
				m.Upstream = &UpstreamConfig{}
			}
			if err := m.Upstream.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Downstream", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Downstream == nil {
				m.Downstream = &DownstreamConfig{}
			}
			if err := m.Downstream.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProver
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UpstreamConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProver
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: UpstreamConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UpstreamConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProxyChain", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ProxyChain == nil {
				m.ProxyChain = &types.Any{}
			}
			if err := m.ProxyChain.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProxyChainProver", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ProxyChainProver == nil {
				m.ProxyChainProver = &types.Any{}
			}
			if err := m.ProxyChainProver.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpstreamClientId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProver
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UpstreamClientId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProver
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DownstreamConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProver
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DownstreamConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DownstreamConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProxyChain", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ProxyChain == nil {
				m.ProxyChain = &types.Any{}
			}
			if err := m.ProxyChain.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProxyChainProver", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ProxyChainProver == nil {
				m.ProxyChainProver = &types.Any{}
			}
			if err := m.ProxyChainProver.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpstreamClientId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProver
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UpstreamClientId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProver
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipProver(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProver
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProver
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProver
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthProver
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProver
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProver
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProver        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProver          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProver = fmt.Errorf("proto: unexpected end of group")
)
