// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fury/committee/v1beta1/committee.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
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

// TallyOption enumerates the valid types of a tally.
type TallyOption int32

const (
	// TALLY_OPTION_UNSPECIFIED defines a null tally option.
	TALLY_OPTION_UNSPECIFIED TallyOption = 0
	// Votes are tallied each block and the proposal passes as soon as the vote threshold is reached
	TALLY_OPTION_FIRST_PAST_THE_POST TallyOption = 1
	// Votes are tallied exactly once, when the deadline time is reached
	TALLY_OPTION_DEADLINE TallyOption = 2
)

var TallyOption_name = map[int32]string{
	0: "TALLY_OPTION_UNSPECIFIED",
	1: "TALLY_OPTION_FIRST_PAST_THE_POST",
	2: "TALLY_OPTION_DEADLINE",
}

var TallyOption_value = map[string]int32{
	"TALLY_OPTION_UNSPECIFIED":         0,
	"TALLY_OPTION_FIRST_PAST_THE_POST": 1,
	"TALLY_OPTION_DEADLINE":            2,
}

func (x TallyOption) String() string {
	return proto.EnumName(TallyOption_name, int32(x))
}

func (TallyOption) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f8b8f3da2af937bb, []int{0}
}

// BaseCommittee is a common type shared by all Committees
type BaseCommittee struct {
	ID          uint64                                          `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Description string                                          `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Members     []github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,3,rep,name=members,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"members,omitempty"`
	Permissions []*types.Any                                    `protobuf:"bytes,4,rep,name=permissions,proto3" json:"permissions,omitempty"`
	// Smallest percentage that must vote for a proposal to pass
	VoteThreshold github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=vote_threshold,json=voteThreshold,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"vote_threshold"`
	// The length of time a proposal remains active for. Proposals will close earlier if they get enough votes.
	ProposalDuration time.Duration `protobuf:"bytes,6,opt,name=proposal_duration,json=proposalDuration,proto3,stdduration" json:"proposal_duration"`
	TallyOption      TallyOption   `protobuf:"varint,7,opt,name=tally_option,json=tallyOption,proto3,enum=fury.committee.v1beta1.TallyOption" json:"tally_option,omitempty"`
}

func (m *BaseCommittee) Reset()      { *m = BaseCommittee{} }
func (*BaseCommittee) ProtoMessage() {}
func (*BaseCommittee) Descriptor() ([]byte, []int) {
	return fileDescriptor_f8b8f3da2af937bb, []int{0}
}
func (m *BaseCommittee) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BaseCommittee) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BaseCommittee.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BaseCommittee) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BaseCommittee.Merge(m, src)
}
func (m *BaseCommittee) XXX_Size() int {
	return m.Size()
}
func (m *BaseCommittee) XXX_DiscardUnknown() {
	xxx_messageInfo_BaseCommittee.DiscardUnknown(m)
}

var xxx_messageInfo_BaseCommittee proto.InternalMessageInfo

// MemberCommittee is an alias of BaseCommittee
type MemberCommittee struct {
	*BaseCommittee `protobuf:"bytes,1,opt,name=base_committee,json=baseCommittee,proto3,embedded=base_committee" json:"base_committee,omitempty"`
}

func (m *MemberCommittee) Reset()      { *m = MemberCommittee{} }
func (*MemberCommittee) ProtoMessage() {}
func (*MemberCommittee) Descriptor() ([]byte, []int) {
	return fileDescriptor_f8b8f3da2af937bb, []int{1}
}
func (m *MemberCommittee) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MemberCommittee) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MemberCommittee.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MemberCommittee) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MemberCommittee.Merge(m, src)
}
func (m *MemberCommittee) XXX_Size() int {
	return m.Size()
}
func (m *MemberCommittee) XXX_DiscardUnknown() {
	xxx_messageInfo_MemberCommittee.DiscardUnknown(m)
}

var xxx_messageInfo_MemberCommittee proto.InternalMessageInfo

// TokenCommittee supports voting on proposals by token holders
type TokenCommittee struct {
	*BaseCommittee `protobuf:"bytes,1,opt,name=base_committee,json=baseCommittee,proto3,embedded=base_committee" json:"base_committee,omitempty"`
	Quorum         github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=quorum,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"quorum"`
	TallyDenom     string                                 `protobuf:"bytes,3,opt,name=tally_denom,json=tallyDenom,proto3" json:"tally_denom,omitempty"`
}

func (m *TokenCommittee) Reset()      { *m = TokenCommittee{} }
func (*TokenCommittee) ProtoMessage() {}
func (*TokenCommittee) Descriptor() ([]byte, []int) {
	return fileDescriptor_f8b8f3da2af937bb, []int{2}
}
func (m *TokenCommittee) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TokenCommittee) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TokenCommittee.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TokenCommittee) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenCommittee.Merge(m, src)
}
func (m *TokenCommittee) XXX_Size() int {
	return m.Size()
}
func (m *TokenCommittee) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenCommittee.DiscardUnknown(m)
}

var xxx_messageInfo_TokenCommittee proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("fury.committee.v1beta1.TallyOption", TallyOption_name, TallyOption_value)
	proto.RegisterType((*BaseCommittee)(nil), "fury.committee.v1beta1.BaseCommittee")
	proto.RegisterType((*MemberCommittee)(nil), "fury.committee.v1beta1.MemberCommittee")
	proto.RegisterType((*TokenCommittee)(nil), "fury.committee.v1beta1.TokenCommittee")
}

func init() {
	proto.RegisterFile("fury/committee/v1beta1/committee.proto", fileDescriptor_f8b8f3da2af937bb)
}

var fileDescriptor_f8b8f3da2af937bb = []byte{
	// 655 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xb6, 0x93, 0x90, 0xd2, 0x4d, 0x1b, 0xd2, 0xa5, 0x54, 0x4e, 0x85, 0x6c, 0xab, 0x40, 0x15,
	0x21, 0xc5, 0x56, 0xc3, 0x8d, 0x5b, 0x5c, 0x27, 0x6a, 0xa4, 0xd2, 0x44, 0x8e, 0x7b, 0x80, 0x8b,
	0xe5, 0x9f, 0x25, 0xb5, 0x1a, 0x7b, 0x83, 0x77, 0x5d, 0x9a, 0x37, 0xe0, 0xc8, 0xb1, 0x47, 0x24,
	0x5e, 0xa1, 0x0f, 0x51, 0xf5, 0x54, 0x71, 0x42, 0x1c, 0x42, 0x49, 0x9f, 0x02, 0x4e, 0xc8, 0x7f,
	0x4d, 0x0a, 0x45, 0x82, 0x03, 0x27, 0x7b, 0xbf, 0xf9, 0x66, 0x66, 0xbf, 0x99, 0x4f, 0x0b, 0x36,
	0x7d, 0xe4, 0x61, 0xd9, 0xc6, 0x9e, 0xe7, 0x52, 0x8a, 0x90, 0x7c, 0xb4, 0x65, 0x21, 0x6a, 0x6e,
	0xcd, 0x10, 0x69, 0x14, 0x60, 0x8a, 0xe1, 0x5a, 0xc4, 0x93, 0x66, 0x68, 0xca, 0x5b, 0xaf, 0xda,
	0x98, 0x78, 0x98, 0x18, 0x31, 0x4b, 0x4e, 0x0e, 0x49, 0xca, 0xfa, 0xea, 0x00, 0x0f, 0x70, 0x82,
	0x47, 0x7f, 0x29, 0x5a, 0x1d, 0x60, 0x3c, 0x18, 0x22, 0x39, 0x3e, 0x59, 0xe1, 0x6b, 0xd9, 0xf4,
	0xc7, 0x69, 0x88, 0xff, 0x35, 0xe4, 0x84, 0x81, 0x49, 0x5d, 0xec, 0x27, 0xf1, 0x8d, 0xef, 0x79,
	0xb0, 0xac, 0x98, 0x04, 0x6d, 0x67, 0xb7, 0x80, 0x6b, 0x20, 0xe7, 0x3a, 0x1c, 0x2b, 0xb2, 0xb5,
	0x82, 0x52, 0x9c, 0x4e, 0x84, 0x5c, 0x47, 0xd5, 0x72, 0xae, 0x03, 0x45, 0x50, 0x72, 0x10, 0xb1,
	0x03, 0x77, 0x14, 0xa5, 0x73, 0x39, 0x91, 0xad, 0x2d, 0x6a, 0xf3, 0x10, 0xb4, 0xc0, 0x82, 0x87,
	0x3c, 0x0b, 0x05, 0x84, 0xcb, 0x8b, 0xf9, 0xda, 0x92, 0xb2, 0xf3, 0x63, 0x22, 0xd4, 0x07, 0x2e,
	0x3d, 0x08, 0xad, 0x48, 0x66, 0x2a, 0x25, 0xfd, 0xd4, 0x89, 0x73, 0x28, 0xd3, 0xf1, 0x08, 0x11,
	0xa9, 0x69, 0xdb, 0x4d, 0xc7, 0x09, 0x10, 0x21, 0x9f, 0x4e, 0xeb, 0xf7, 0x53, 0xc1, 0x29, 0xa2,
	0x8c, 0x29, 0x22, 0x5a, 0x56, 0x18, 0xb6, 0x41, 0x69, 0x84, 0x02, 0xcf, 0x25, 0xc4, 0xc5, 0x3e,
	0xe1, 0x0a, 0x62, 0xbe, 0x56, 0x6a, 0xac, 0x4a, 0x89, 0x4a, 0x29, 0x53, 0x29, 0x35, 0xfd, 0xb1,
	0x52, 0x3e, 0x3f, 0xad, 0x83, 0xde, 0x35, 0x59, 0x9b, 0x4f, 0x84, 0xfb, 0xa0, 0x7c, 0x84, 0x29,
	0x32, 0xe8, 0x41, 0x80, 0xc8, 0x01, 0x1e, 0x3a, 0xdc, 0x9d, 0x48, 0x90, 0x22, 0x9d, 0x4d, 0x04,
	0xe6, 0xcb, 0x44, 0xd8, 0xfc, 0x8b, 0x6b, 0xab, 0xc8, 0xd6, 0x96, 0xa3, 0x2a, 0x7a, 0x56, 0x04,
	0xf6, 0xc0, 0xca, 0x28, 0xc0, 0x23, 0x4c, 0xcc, 0xa1, 0x91, 0x4d, 0x9a, 0x2b, 0x8a, 0x6c, 0xad,
	0xd4, 0xa8, 0xfe, 0x76, 0x49, 0x35, 0x25, 0x28, 0x77, 0xa3, 0xa6, 0x27, 0x5f, 0x05, 0x56, 0xab,
	0x64, 0xd9, 0x59, 0x0c, 0xb6, 0xc1, 0x12, 0x35, 0x87, 0xc3, 0xb1, 0x81, 0x93, 0xb9, 0x2f, 0x88,
	0x6c, 0xad, 0xdc, 0x78, 0x24, 0xdd, 0xee, 0x1d, 0x49, 0x8f, 0xb8, 0xdd, 0x98, 0xaa, 0x95, 0xe8,
	0xec, 0xf0, 0x7c, 0xe5, 0xe4, 0x83, 0xc0, 0x9c, 0x9f, 0xd6, 0x17, 0xaf, 0x37, 0xbd, 0x71, 0x0c,
	0xee, 0xbd, 0x88, 0xc7, 0x3a, 0x5b, 0xbe, 0x06, 0xca, 0x96, 0x49, 0x90, 0x71, 0x5d, 0x38, 0x36,
	0x42, 0xa9, 0xf1, 0xe4, 0x4f, 0xfd, 0x6e, 0x78, 0x47, 0x29, 0x5c, 0x4c, 0x04, 0x56, 0x5b, 0xb6,
	0xe6, 0xc1, 0xdb, 0x3a, 0x5f, 0xb2, 0xa0, 0xac, 0xe3, 0x43, 0xe4, 0xff, 0xd7, 0xce, 0xb0, 0x0d,
	0x8a, 0x6f, 0x42, 0x1c, 0x84, 0x5e, 0xe2, 0xd6, 0x7f, 0x5e, 0x6e, 0x9a, 0x0d, 0x05, 0x90, 0x8c,
	0xd2, 0x70, 0x90, 0x8f, 0x3d, 0x2e, 0x1f, 0x5b, 0x1f, 0xc4, 0x90, 0x1a, 0x21, 0xb7, 0x48, 0x7c,
	0x1a, 0x80, 0xd2, 0xdc, 0x2e, 0xe0, 0x43, 0xc0, 0xe9, 0xcd, 0xdd, 0xdd, 0x97, 0x46, 0xb7, 0xa7,
	0x77, 0xba, 0x7b, 0xc6, 0xfe, 0x5e, 0xbf, 0xd7, 0xda, 0xee, 0xb4, 0x3b, 0x2d, 0xb5, 0xc2, 0xc0,
	0xc7, 0x40, 0xbc, 0x11, 0x6d, 0x77, 0xb4, 0xbe, 0x6e, 0xf4, 0x9a, 0x7d, 0xdd, 0xd0, 0x77, 0x5a,
	0x46, 0xaf, 0xdb, 0xd7, 0x2b, 0x2c, 0xac, 0x82, 0x07, 0x37, 0x58, 0x6a, 0xab, 0xa9, 0xee, 0x76,
	0xf6, 0x5a, 0x95, 0xdc, 0x7a, 0xe1, 0xdd, 0x47, 0x9e, 0x51, 0xba, 0x67, 0xdf, 0x78, 0xe6, 0x6c,
	0xca, 0xb3, 0x17, 0x53, 0x9e, 0xbd, 0x9c, 0xf2, 0xec, 0xfb, 0x2b, 0x9e, 0xb9, 0xb8, 0xe2, 0x99,
	0xcf, 0x57, 0x3c, 0xf3, 0x6a, 0x6b, 0x4e, 0xb5, 0xeb, 0xdb, 0xa1, 0x15, 0x92, 0xba, 0x8f, 0xe8,
	0x5b, 0x1c, 0x1c, 0xca, 0xf1, 0x8b, 0x75, 0x3c, 0xf7, 0x66, 0xc5, 0x43, 0xb0, 0x8a, 0xb1, 0x57,
	0x9f, 0xfd, 0x0c, 0x00, 0x00, 0xff, 0xff, 0x0f, 0x2d, 0x28, 0xb5, 0xd2, 0x04, 0x00, 0x00,
}

func (m *BaseCommittee) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BaseCommittee) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BaseCommittee) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.TallyOption != 0 {
		i = encodeVarintCommittee(dAtA, i, uint64(m.TallyOption))
		i--
		dAtA[i] = 0x38
	}
	n1, err1 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.ProposalDuration, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.ProposalDuration):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintCommittee(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x32
	{
		size := m.VoteThreshold.Size()
		i -= size
		if _, err := m.VoteThreshold.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintCommittee(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if len(m.Permissions) > 0 {
		for iNdEx := len(m.Permissions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Permissions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintCommittee(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Members) > 0 {
		for iNdEx := len(m.Members) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Members[iNdEx])
			copy(dAtA[i:], m.Members[iNdEx])
			i = encodeVarintCommittee(dAtA, i, uint64(len(m.Members[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintCommittee(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if m.ID != 0 {
		i = encodeVarintCommittee(dAtA, i, uint64(m.ID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MemberCommittee) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MemberCommittee) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MemberCommittee) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.BaseCommittee != nil {
		{
			size, err := m.BaseCommittee.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintCommittee(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TokenCommittee) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TokenCommittee) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TokenCommittee) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TallyDenom) > 0 {
		i -= len(m.TallyDenom)
		copy(dAtA[i:], m.TallyDenom)
		i = encodeVarintCommittee(dAtA, i, uint64(len(m.TallyDenom)))
		i--
		dAtA[i] = 0x1a
	}
	{
		size := m.Quorum.Size()
		i -= size
		if _, err := m.Quorum.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintCommittee(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.BaseCommittee != nil {
		{
			size, err := m.BaseCommittee.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintCommittee(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintCommittee(dAtA []byte, offset int, v uint64) int {
	offset -= sovCommittee(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BaseCommittee) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovCommittee(uint64(m.ID))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovCommittee(uint64(l))
	}
	if len(m.Members) > 0 {
		for _, b := range m.Members {
			l = len(b)
			n += 1 + l + sovCommittee(uint64(l))
		}
	}
	if len(m.Permissions) > 0 {
		for _, e := range m.Permissions {
			l = e.Size()
			n += 1 + l + sovCommittee(uint64(l))
		}
	}
	l = m.VoteThreshold.Size()
	n += 1 + l + sovCommittee(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.ProposalDuration)
	n += 1 + l + sovCommittee(uint64(l))
	if m.TallyOption != 0 {
		n += 1 + sovCommittee(uint64(m.TallyOption))
	}
	return n
}

func (m *MemberCommittee) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BaseCommittee != nil {
		l = m.BaseCommittee.Size()
		n += 1 + l + sovCommittee(uint64(l))
	}
	return n
}

func (m *TokenCommittee) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BaseCommittee != nil {
		l = m.BaseCommittee.Size()
		n += 1 + l + sovCommittee(uint64(l))
	}
	l = m.Quorum.Size()
	n += 1 + l + sovCommittee(uint64(l))
	l = len(m.TallyDenom)
	if l > 0 {
		n += 1 + l + sovCommittee(uint64(l))
	}
	return n
}

func sovCommittee(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCommittee(x uint64) (n int) {
	return sovCommittee(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BaseCommittee) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommittee
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
			return fmt.Errorf("proto: BaseCommittee: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BaseCommittee: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommittee
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommittee
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
				return ErrInvalidLengthCommittee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommittee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Members", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommittee
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthCommittee
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthCommittee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Members = append(m.Members, make([]byte, postIndex-iNdEx))
			copy(m.Members[len(m.Members)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Permissions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommittee
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
				return ErrInvalidLengthCommittee
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommittee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Permissions = append(m.Permissions, &types.Any{})
			if err := m.Permissions[len(m.Permissions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VoteThreshold", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommittee
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
				return ErrInvalidLengthCommittee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommittee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.VoteThreshold.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProposalDuration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommittee
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
				return ErrInvalidLengthCommittee
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommittee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.ProposalDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TallyOption", wireType)
			}
			m.TallyOption = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommittee
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TallyOption |= TallyOption(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCommittee(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommittee
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
func (m *MemberCommittee) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommittee
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
			return fmt.Errorf("proto: MemberCommittee: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MemberCommittee: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseCommittee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommittee
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
				return ErrInvalidLengthCommittee
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommittee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.BaseCommittee == nil {
				m.BaseCommittee = &BaseCommittee{}
			}
			if err := m.BaseCommittee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommittee(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommittee
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
func (m *TokenCommittee) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommittee
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
			return fmt.Errorf("proto: TokenCommittee: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TokenCommittee: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseCommittee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommittee
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
				return ErrInvalidLengthCommittee
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommittee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.BaseCommittee == nil {
				m.BaseCommittee = &BaseCommittee{}
			}
			if err := m.BaseCommittee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Quorum", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommittee
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
				return ErrInvalidLengthCommittee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommittee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Quorum.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TallyDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommittee
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
				return ErrInvalidLengthCommittee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommittee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TallyDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommittee(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommittee
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
func skipCommittee(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCommittee
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
					return 0, ErrIntOverflowCommittee
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
					return 0, ErrIntOverflowCommittee
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
				return 0, ErrInvalidLengthCommittee
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCommittee
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCommittee
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCommittee        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCommittee          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCommittee = fmt.Errorf("proto: unexpected end of group")
)
