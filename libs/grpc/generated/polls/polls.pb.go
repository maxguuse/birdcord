// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.4
// source: polls.proto

package polls

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreatePollRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title           string   `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Options         []string `protobuf:"bytes,2,rep,name=options,proto3" json:"options,omitempty"`
	DiscordId       string   `protobuf:"bytes,3,opt,name=discord_id,json=discordId,proto3" json:"discord_id,omitempty"`
	DiscordAuthorId string   `protobuf:"bytes,4,opt,name=discord_author_id,json=discordAuthorId,proto3" json:"discord_author_id,omitempty"`
	DiscordGuildId  string   `protobuf:"bytes,5,opt,name=discord_guild_id,json=discordGuildId,proto3" json:"discord_guild_id,omitempty"`
	ChannelId       string   `protobuf:"bytes,6,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
}

func (x *CreatePollRequest) Reset() {
	*x = CreatePollRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polls_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePollRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePollRequest) ProtoMessage() {}

func (x *CreatePollRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polls_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePollRequest.ProtoReflect.Descriptor instead.
func (*CreatePollRequest) Descriptor() ([]byte, []int) {
	return file_polls_proto_rawDescGZIP(), []int{0}
}

func (x *CreatePollRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *CreatePollRequest) GetOptions() []string {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *CreatePollRequest) GetDiscordId() string {
	if x != nil {
		return x.DiscordId
	}
	return ""
}

func (x *CreatePollRequest) GetDiscordAuthorId() string {
	if x != nil {
		return x.DiscordAuthorId
	}
	return ""
}

func (x *CreatePollRequest) GetDiscordGuildId() string {
	if x != nil {
		return x.DiscordGuildId
	}
	return ""
}

func (x *CreatePollRequest) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

type CreatePollResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PollId  int32     `protobuf:"varint,1,opt,name=poll_id,json=pollId,proto3" json:"poll_id,omitempty"`
	Options []*Option `protobuf:"bytes,2,rep,name=options,proto3" json:"options,omitempty"`
}

func (x *CreatePollResponse) Reset() {
	*x = CreatePollResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polls_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePollResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePollResponse) ProtoMessage() {}

func (x *CreatePollResponse) ProtoReflect() protoreflect.Message {
	mi := &file_polls_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePollResponse.ProtoReflect.Descriptor instead.
func (*CreatePollResponse) Descriptor() ([]byte, []int) {
	return file_polls_proto_rawDescGZIP(), []int{1}
}

func (x *CreatePollResponse) GetPollId() int32 {
	if x != nil {
		return x.PollId
	}
	return 0
}

func (x *CreatePollResponse) GetOptions() []*Option {
	if x != nil {
		return x.Options
	}
	return nil
}

type VoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PollId        int32  `protobuf:"varint,1,opt,name=poll_id,json=pollId,proto3" json:"poll_id,omitempty"`
	OptionId      int32  `protobuf:"varint,2,opt,name=option_id,json=optionId,proto3" json:"option_id,omitempty"`
	DiscordUserId string `protobuf:"bytes,3,opt,name=discord_user_id,json=discordUserId,proto3" json:"discord_user_id,omitempty"`
}

func (x *VoteRequest) Reset() {
	*x = VoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polls_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequest) ProtoMessage() {}

func (x *VoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polls_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequest.ProtoReflect.Descriptor instead.
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return file_polls_proto_rawDescGZIP(), []int{2}
}

func (x *VoteRequest) GetPollId() int32 {
	if x != nil {
		return x.PollId
	}
	return 0
}

func (x *VoteRequest) GetOptionId() int32 {
	if x != nil {
		return x.OptionId
	}
	return 0
}

func (x *VoteRequest) GetDiscordUserId() string {
	if x != nil {
		return x.DiscordUserId
	}
	return ""
}

type VoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DiscordToken string    `protobuf:"bytes,1,opt,name=discord_token,json=discordToken,proto3" json:"discord_token,omitempty"`
	Options      []*Option `protobuf:"bytes,2,rep,name=options,proto3" json:"options,omitempty"`
	TotalVotes   int32     `protobuf:"varint,3,opt,name=total_votes,json=totalVotes,proto3" json:"total_votes,omitempty"`
	Title        string    `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	Success      bool      `protobuf:"varint,5,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *VoteResponse) Reset() {
	*x = VoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polls_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteResponse) ProtoMessage() {}

func (x *VoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_polls_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteResponse.ProtoReflect.Descriptor instead.
func (*VoteResponse) Descriptor() ([]byte, []int) {
	return file_polls_proto_rawDescGZIP(), []int{3}
}

func (x *VoteResponse) GetDiscordToken() string {
	if x != nil {
		return x.DiscordToken
	}
	return ""
}

func (x *VoteResponse) GetOptions() []*Option {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *VoteResponse) GetTotalVotes() int32 {
	if x != nil {
		return x.TotalVotes
	}
	return 0
}

func (x *VoteResponse) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *VoteResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type Option struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title      string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Id         int32  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	TotalVotes int32  `protobuf:"varint,3,opt,name=total_votes,json=totalVotes,proto3" json:"total_votes,omitempty"`
}

func (x *Option) Reset() {
	*x = Option{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polls_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Option) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Option) ProtoMessage() {}

func (x *Option) ProtoReflect() protoreflect.Message {
	mi := &file_polls_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Option.ProtoReflect.Descriptor instead.
func (*Option) Descriptor() ([]byte, []int) {
	return file_polls_proto_rawDescGZIP(), []int{4}
}

func (x *Option) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Option) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Option) GetTotalVotes() int32 {
	if x != nil {
		return x.TotalVotes
	}
	return 0
}

type GetActivePollsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DiscordGuildId string `protobuf:"bytes,1,opt,name=discord_guild_id,json=discordGuildId,proto3" json:"discord_guild_id,omitempty"`
}

func (x *GetActivePollsRequest) Reset() {
	*x = GetActivePollsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polls_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetActivePollsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetActivePollsRequest) ProtoMessage() {}

func (x *GetActivePollsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polls_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetActivePollsRequest.ProtoReflect.Descriptor instead.
func (*GetActivePollsRequest) Descriptor() ([]byte, []int) {
	return file_polls_proto_rawDescGZIP(), []int{5}
}

func (x *GetActivePollsRequest) GetDiscordGuildId() string {
	if x != nil {
		return x.DiscordGuildId
	}
	return ""
}

type GetActivePollsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Polls []*Poll `protobuf:"bytes,1,rep,name=polls,proto3" json:"polls,omitempty"`
}

func (x *GetActivePollsResponse) Reset() {
	*x = GetActivePollsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polls_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetActivePollsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetActivePollsResponse) ProtoMessage() {}

func (x *GetActivePollsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_polls_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetActivePollsResponse.ProtoReflect.Descriptor instead.
func (*GetActivePollsResponse) Descriptor() ([]byte, []int) {
	return file_polls_proto_rawDescGZIP(), []int{6}
}

func (x *GetActivePollsResponse) GetPolls() []*Poll {
	if x != nil {
		return x.Polls
	}
	return nil
}

type Poll struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Id    int32  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Poll) Reset() {
	*x = Poll{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polls_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Poll) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Poll) ProtoMessage() {}

func (x *Poll) ProtoReflect() protoreflect.Message {
	mi := &file_polls_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Poll.ProtoReflect.Descriptor instead.
func (*Poll) Descriptor() ([]byte, []int) {
	return file_polls_proto_rawDescGZIP(), []int{7}
}

func (x *Poll) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Poll) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type StopPollRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PollId int32 `protobuf:"varint,1,opt,name=poll_id,json=pollId,proto3" json:"poll_id,omitempty"`
}

func (x *StopPollRequest) Reset() {
	*x = StopPollRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polls_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopPollRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopPollRequest) ProtoMessage() {}

func (x *StopPollRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polls_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopPollRequest.ProtoReflect.Descriptor instead.
func (*StopPollRequest) Descriptor() ([]byte, []int) {
	return file_polls_proto_rawDescGZIP(), []int{8}
}

func (x *StopPollRequest) GetPollId() int32 {
	if x != nil {
		return x.PollId
	}
	return 0
}

type StopPollResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DiscordToken string    `protobuf:"bytes,1,opt,name=discord_token,json=discordToken,proto3" json:"discord_token,omitempty"`
	Winners      []*Option `protobuf:"bytes,2,rep,name=winners,proto3" json:"winners,omitempty"`
	TotalVotes   int32     `protobuf:"varint,3,opt,name=total_votes,json=totalVotes,proto3" json:"total_votes,omitempty"`
	Title        string    `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
}

func (x *StopPollResponse) Reset() {
	*x = StopPollResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polls_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopPollResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopPollResponse) ProtoMessage() {}

func (x *StopPollResponse) ProtoReflect() protoreflect.Message {
	mi := &file_polls_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopPollResponse.ProtoReflect.Descriptor instead.
func (*StopPollResponse) Descriptor() ([]byte, []int) {
	return file_polls_proto_rawDescGZIP(), []int{9}
}

func (x *StopPollResponse) GetDiscordToken() string {
	if x != nil {
		return x.DiscordToken
	}
	return ""
}

func (x *StopPollResponse) GetWinners() []*Option {
	if x != nil {
		return x.Winners
	}
	return nil
}

func (x *StopPollResponse) GetTotalVotes() int32 {
	if x != nil {
		return x.TotalVotes
	}
	return 0
}

func (x *StopPollResponse) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

var File_polls_proto protoreflect.FileDescriptor

var file_polls_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x6f, 0x6c, 0x6c, 0x73, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xd7, 0x01, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x6c, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x73,
	0x63, 0x6f, 0x72, 0x64, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72,
	0x64, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x49, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x67, 0x75,
	0x69, 0x6c, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x72, 0x64, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a,
	0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64, 0x22, 0x56, 0x0a, 0x12, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x6f, 0x6c, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x70, 0x6f, 0x6c, 0x6c, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x07, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x6f,
	0x6c, 0x6c, 0x73, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x22, 0x6b, 0x0a, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x6f, 0x6c, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x06, 0x70, 0x6f, 0x6c, 0x6c, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x72, 0x64, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x22, 0xad, 0x01, 0x0a, 0x0c, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x23, 0x0a, 0x0d, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72,
	0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x27, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x2e,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x56, 0x6f, 0x74, 0x65, 0x73,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x22, 0x4f, 0x0a, 0x06, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x56, 0x6f, 0x74, 0x65,
	0x73, 0x22, 0x41, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x50, 0x6f,
	0x6c, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x47, 0x75, 0x69,
	0x6c, 0x64, 0x49, 0x64, 0x22, 0x3b, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x50, 0x6f, 0x6c, 0x6c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21,
	0x0a, 0x05, 0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e,
	0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x2e, 0x50, 0x6f, 0x6c, 0x6c, 0x52, 0x05, 0x70, 0x6f, 0x6c, 0x6c,
	0x73, 0x22, 0x2c, 0x0a, 0x04, 0x50, 0x6f, 0x6c, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x2a, 0x0a, 0x0f, 0x53, 0x74, 0x6f, 0x70, 0x50, 0x6f, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x6f, 0x6c, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x06, 0x70, 0x6f, 0x6c, 0x6c, 0x49, 0x64, 0x22, 0x97, 0x01, 0x0a, 0x10,
	0x53, 0x74, 0x6f, 0x70, 0x50, 0x6f, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x23, 0x0a, 0x0d, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x27, 0x0a, 0x07, 0x77, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x2e, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x77, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x73, 0x12, 0x1f,
	0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x56, 0x6f, 0x74, 0x65, 0x73, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x32, 0x87, 0x02, 0x0a, 0x05, 0x50, 0x6f, 0x6c, 0x6c, 0x73, 0x12,
	0x41, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x6c, 0x6c, 0x12, 0x18, 0x2e,
	0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x6c, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x50,
	0x6f, 0x6c, 0x6c, 0x73, 0x12, 0x1c, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x2e, 0x47, 0x65, 0x74,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x50, 0x6f, 0x6c, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x50, 0x6f, 0x6c, 0x6c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3b, 0x0a, 0x08, 0x53, 0x74, 0x6f, 0x70, 0x50, 0x6f, 0x6c, 0x6c, 0x12, 0x16, 0x2e,
	0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x50, 0x6f, 0x6c, 0x6c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x2e, 0x53, 0x74,
	0x6f, 0x70, 0x50, 0x6f, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f,
	0x0a, 0x04, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x12, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x2e, 0x56,
	0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x70, 0x6f, 0x6c,
	0x6c, 0x73, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x61,
	0x78, 0x67, 0x75, 0x75, 0x73, 0x65, 0x2f, 0x62, 0x69, 0x72, 0x64, 0x63, 0x6f, 0x72, 0x64, 0x2f,
	0x6c, 0x69, 0x62, 0x73, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x74, 0x65, 0x64, 0x2f, 0x70, 0x6f, 0x6c, 0x6c, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_polls_proto_rawDescOnce sync.Once
	file_polls_proto_rawDescData = file_polls_proto_rawDesc
)

func file_polls_proto_rawDescGZIP() []byte {
	file_polls_proto_rawDescOnce.Do(func() {
		file_polls_proto_rawDescData = protoimpl.X.CompressGZIP(file_polls_proto_rawDescData)
	})
	return file_polls_proto_rawDescData
}

var file_polls_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_polls_proto_goTypes = []interface{}{
	(*CreatePollRequest)(nil),      // 0: polls.CreatePollRequest
	(*CreatePollResponse)(nil),     // 1: polls.CreatePollResponse
	(*VoteRequest)(nil),            // 2: polls.VoteRequest
	(*VoteResponse)(nil),           // 3: polls.VoteResponse
	(*Option)(nil),                 // 4: polls.Option
	(*GetActivePollsRequest)(nil),  // 5: polls.GetActivePollsRequest
	(*GetActivePollsResponse)(nil), // 6: polls.GetActivePollsResponse
	(*Poll)(nil),                   // 7: polls.Poll
	(*StopPollRequest)(nil),        // 8: polls.StopPollRequest
	(*StopPollResponse)(nil),       // 9: polls.StopPollResponse
}
var file_polls_proto_depIdxs = []int32{
	4, // 0: polls.CreatePollResponse.options:type_name -> polls.Option
	4, // 1: polls.VoteResponse.options:type_name -> polls.Option
	7, // 2: polls.GetActivePollsResponse.polls:type_name -> polls.Poll
	4, // 3: polls.StopPollResponse.winners:type_name -> polls.Option
	0, // 4: polls.Polls.CreatePoll:input_type -> polls.CreatePollRequest
	5, // 5: polls.Polls.GetActivePolls:input_type -> polls.GetActivePollsRequest
	8, // 6: polls.Polls.StopPoll:input_type -> polls.StopPollRequest
	2, // 7: polls.Polls.Vote:input_type -> polls.VoteRequest
	1, // 8: polls.Polls.CreatePoll:output_type -> polls.CreatePollResponse
	6, // 9: polls.Polls.GetActivePolls:output_type -> polls.GetActivePollsResponse
	9, // 10: polls.Polls.StopPoll:output_type -> polls.StopPollResponse
	3, // 11: polls.Polls.Vote:output_type -> polls.VoteResponse
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_polls_proto_init() }
func file_polls_proto_init() {
	if File_polls_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_polls_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePollRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_polls_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePollResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_polls_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_polls_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_polls_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Option); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_polls_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetActivePollsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_polls_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetActivePollsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_polls_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Poll); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_polls_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopPollRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_polls_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopPollResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_polls_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_polls_proto_goTypes,
		DependencyIndexes: file_polls_proto_depIdxs,
		MessageInfos:      file_polls_proto_msgTypes,
	}.Build()
	File_polls_proto = out.File
	file_polls_proto_rawDesc = nil
	file_polls_proto_goTypes = nil
	file_polls_proto_depIdxs = nil
}
