// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fury/bep3/v1beta1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the pricefeed module's genesis state.
type GenesisState struct {
	// params defines all the paramaters of the module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	// atomic_swaps represents the state of stored atomic swaps
	AtomicSwaps AtomicSwaps `protobuf:"bytes,2,rep,name=atomic_swaps,json=atomicSwaps,proto3,castrepeated=AtomicSwaps" json:"atomic_swaps"`
	// supplies represents the supply information of each atomic swap
	Supplies AssetSupplies `protobuf:"bytes,3,rep,name=supplies,proto3,castrepeated=AssetSupplies" json:"supplies"`
	// previous_block_time represents the time of the previous block
	PreviousBlockTime time.Time `protobuf:"bytes,4,opt,name=previous_block_time,json=previousBlockTime,proto3,stdtime" json:"previous_block_time"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_808f03ab1253b770, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetAtomicSwaps() AtomicSwaps {
	if m != nil {
		return m.AtomicSwaps
	}
	return nil
}

func (m *GenesisState) GetSupplies() AssetSupplies {
	if m != nil {
		return m.Supplies
	}
	return nil
}

func (m *GenesisState) GetPreviousBlockTime() time.Time {
	if m != nil {
		return m.PreviousBlockTime
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "fury.bep3.v1beta1.GenesisState")
}

func init() { proto.RegisterFile("fury/bep3/v1beta1/genesis.proto", fileDescriptor_808f03ab1253b770) }

var fileDescriptor_808f03ab1253b770 = []byte{
	// 364 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0x31, 0x6f, 0xda, 0x40,
	0x14, 0xc7, 0x6d, 0x40, 0x08, 0xd9, 0x74, 0xc0, 0xb4, 0x92, 0x8b, 0x5a, 0x1b, 0x75, 0x62, 0x68,
	0xef, 0x04, 0x0c, 0x9d, 0xeb, 0xa1, 0x5d, 0x1b, 0x43, 0x96, 0x2c, 0xe8, 0x6c, 0x5d, 0x9c, 0x13,
	0xd8, 0x77, 0xf2, 0x3b, 0x43, 0xf8, 0x16, 0x7c, 0x8e, 0x7c, 0x91, 0x30, 0x32, 0x66, 0x0a, 0x11,
	0x7c, 0x91, 0xe8, 0xce, 0x26, 0x0c, 0xb0, 0xf9, 0xee, 0xfd, 0xde, 0xef, 0xf9, 0xfe, 0xcf, 0xf2,
	0x33, 0x9a, 0x72, 0x1c, 0x51, 0x31, 0xc6, 0xcb, 0x61, 0x44, 0x25, 0x19, 0xe2, 0x84, 0x66, 0x14,
	0x18, 0x20, 0x91, 0x73, 0xc9, 0x9d, 0x8e, 0x02, 0x90, 0x02, 0x50, 0x05, 0xf4, 0x3e, 0x27, 0x3c,
	0xe1, 0xba, 0x8a, 0xd5, 0x57, 0x09, 0xf6, 0xfc, 0x84, 0xf3, 0x64, 0x41, 0xb1, 0x3e, 0x45, 0xc5,
	0x3d, 0x96, 0x2c, 0xa5, 0x20, 0x49, 0x2a, 0x2a, 0xe0, 0xdb, 0xe5, 0x28, 0xad, 0xd5, 0xd5, 0x1f,
	0xcf, 0x35, 0xab, 0xfd, 0xaf, 0x9c, 0x3c, 0x91, 0x44, 0x52, 0xe7, 0xb7, 0xd5, 0x14, 0x24, 0x27,
	0x29, 0xb8, 0x66, 0xdf, 0x1c, 0xd8, 0xa3, 0xaf, 0xe8, 0xe2, 0x4f, 0xd0, 0x7f, 0x0d, 0x04, 0x8d,
	0xed, 0xab, 0x6f, 0x84, 0x15, 0xee, 0xdc, 0x5a, 0x6d, 0x22, 0x79, 0xca, 0xe2, 0x19, 0xac, 0x88,
	0x00, 0xb7, 0xd6, 0xaf, 0x0f, 0xec, 0xd1, 0xf7, 0x2b, 0xed, 0x7f, 0x34, 0x36, 0x59, 0x11, 0x11,
	0x74, 0x95, 0xe2, 0x69, 0xef, 0xdb, 0xe7, 0x3b, 0x08, 0x6d, 0x72, 0x3e, 0x38, 0x37, 0x56, 0x0b,
	0x0a, 0x21, 0x16, 0x8c, 0x82, 0x5b, 0xd7, 0x4a, 0xef, 0x9a, 0x12, 0x80, 0xca, 0x89, 0xe2, 0xd6,
	0xc1, 0x97, 0xca, 0xf9, 0xe9, 0x7c, 0xc9, 0x28, 0x84, 0x1f, 0x1a, 0x67, 0x6a, 0x75, 0x45, 0x4e,
	0x97, 0x8c, 0x17, 0x30, 0x8b, 0x16, 0x3c, 0x9e, 0xcf, 0x54, 0x66, 0x6e, 0x43, 0xbf, 0xb7, 0x87,
	0xca, 0x40, 0xd1, 0x29, 0x50, 0x34, 0x3d, 0x05, 0x1a, 0xb4, 0x94, 0x79, 0xb3, 0xf7, 0xcd, 0xb0,
	0x73, 0x12, 0x04, 0xaa, 0x5f, 0x11, 0xc1, 0xdf, 0xed, 0xc1, 0x33, 0x77, 0x07, 0xcf, 0x7c, 0x3b,
	0x78, 0xe6, 0xe6, 0xe8, 0x19, 0xbb, 0xa3, 0x67, 0xbc, 0x1c, 0x3d, 0xe3, 0xee, 0x67, 0xc2, 0xe4,
	0x43, 0x11, 0xa1, 0x98, 0xa7, 0x98, 0x65, 0x71, 0x11, 0x15, 0xf0, 0x2b, 0xa3, 0x72, 0xc5, 0xf3,
	0x39, 0xd6, 0xcb, 0x79, 0x2c, 0xd7, 0x23, 0xd7, 0x82, 0x42, 0xd4, 0xd4, 0x83, 0xc7, 0xef, 0x01,
	0x00, 0x00, 0xff, 0xff, 0xd2, 0x60, 0xb5, 0xc6, 0x23, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.PreviousBlockTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.PreviousBlockTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintGenesis(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	if len(m.Supplies) > 0 {
		for iNdEx := len(m.Supplies) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Supplies[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.AtomicSwaps) > 0 {
		for iNdEx := len(m.AtomicSwaps) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AtomicSwaps[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.AtomicSwaps) > 0 {
		for _, e := range m.AtomicSwaps {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Supplies) > 0 {
		for _, e := range m.Supplies {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.PreviousBlockTime)
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AtomicSwaps", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AtomicSwaps = append(m.AtomicSwaps, AtomicSwap{})
			if err := m.AtomicSwaps[len(m.AtomicSwaps)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Supplies", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Supplies = append(m.Supplies, AssetSupply{})
			if err := m.Supplies[len(m.Supplies)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PreviousBlockTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.PreviousBlockTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
