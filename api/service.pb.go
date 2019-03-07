// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import duration "github.com/golang/protobuf/ptypes/duration"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Indicates the type of speech event.
type StreamingRecognizeResponse_SpeechEventType int32

const (
	// No speech event specified.
	StreamingRecognizeResponse_SPEECH_EVENT_UNSPECIFIED StreamingRecognizeResponse_SpeechEventType = 0
	// This event indicates that the server has detected the end of the user's
	// speech utterance and expects no additional speech. Therefore, the server
	// will not process additional audio (although it may subsequently return
	// additional results). The client should stop sending additional audio
	// data, half-close the gRPC connection, and wait for any additional results
	// until the server closes the gRPC connection. This event is only sent if
	// `single_utterance` was set to `true`, and is not used otherwise.
	StreamingRecognizeResponse_END_OF_SINGLE_UTTERANCE StreamingRecognizeResponse_SpeechEventType = 1
)

var StreamingRecognizeResponse_SpeechEventType_name = map[int32]string{
	0: "SPEECH_EVENT_UNSPECIFIED",
	1: "END_OF_SINGLE_UTTERANCE",
}
var StreamingRecognizeResponse_SpeechEventType_value = map[string]int32{
	"SPEECH_EVENT_UNSPECIFIED": 0,
	"END_OF_SINGLE_UTTERANCE":  1,
}

func (x StreamingRecognizeResponse_SpeechEventType) String() string {
	return proto.EnumName(StreamingRecognizeResponse_SpeechEventType_name, int32(x))
}
func (StreamingRecognizeResponse_SpeechEventType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_service_339644b7ccc8fec8, []int{3, 0}
}

type TTSRequest struct {
	Text                 string   `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TTSRequest) Reset()         { *m = TTSRequest{} }
func (m *TTSRequest) String() string { return proto.CompactTextString(m) }
func (*TTSRequest) ProtoMessage()    {}
func (*TTSRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_339644b7ccc8fec8, []int{0}
}
func (m *TTSRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TTSRequest.Unmarshal(m, b)
}
func (m *TTSRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TTSRequest.Marshal(b, m, deterministic)
}
func (dst *TTSRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TTSRequest.Merge(dst, src)
}
func (m *TTSRequest) XXX_Size() int {
	return xxx_messageInfo_TTSRequest.Size(m)
}
func (m *TTSRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TTSRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TTSRequest proto.InternalMessageInfo

func (m *TTSRequest) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type TTSReply struct {
	MP3Response          []byte   `protobuf:"bytes,1,opt,name=MP3Response,proto3" json:"MP3Response,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TTSReply) Reset()         { *m = TTSReply{} }
func (m *TTSReply) String() string { return proto.CompactTextString(m) }
func (*TTSReply) ProtoMessage()    {}
func (*TTSReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_339644b7ccc8fec8, []int{1}
}
func (m *TTSReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TTSReply.Unmarshal(m, b)
}
func (m *TTSReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TTSReply.Marshal(b, m, deterministic)
}
func (dst *TTSReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TTSReply.Merge(dst, src)
}
func (m *TTSReply) XXX_Size() int {
	return xxx_messageInfo_TTSReply.Size(m)
}
func (m *TTSReply) XXX_DiscardUnknown() {
	xxx_messageInfo_TTSReply.DiscardUnknown(m)
}

var xxx_messageInfo_TTSReply proto.InternalMessageInfo

func (m *TTSReply) GetMP3Response() []byte {
	if m != nil {
		return m.MP3Response
	}
	return nil
}

type STTRequest struct {
	STTBuffer            []byte   `protobuf:"bytes,1,opt,name=STTBuffer,proto3" json:"STTBuffer,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *STTRequest) Reset()         { *m = STTRequest{} }
func (m *STTRequest) String() string { return proto.CompactTextString(m) }
func (*STTRequest) ProtoMessage()    {}
func (*STTRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_339644b7ccc8fec8, []int{2}
}
func (m *STTRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_STTRequest.Unmarshal(m, b)
}
func (m *STTRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_STTRequest.Marshal(b, m, deterministic)
}
func (dst *STTRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_STTRequest.Merge(dst, src)
}
func (m *STTRequest) XXX_Size() int {
	return xxx_messageInfo_STTRequest.Size(m)
}
func (m *STTRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_STTRequest.DiscardUnknown(m)
}

var xxx_messageInfo_STTRequest proto.InternalMessageInfo

func (m *STTRequest) GetSTTBuffer() []byte {
	if m != nil {
		return m.STTBuffer
	}
	return nil
}

type StreamingRecognizeResponse struct {
	// Output only. This repeated list contains zero or more results that
	// correspond to consecutive portions of the audio currently being processed.
	// It contains zero or one `is_final=true` result (the newly settled portion),
	// followed by zero or more `is_final=false` results (the interim results).
	Results []*StreamingRecognitionResult `protobuf:"bytes,2,rep,name=results,proto3" json:"results,omitempty"`
	// Output only. Indicates the type of speech event.
	SpeechEventType      StreamingRecognizeResponse_SpeechEventType `protobuf:"varint,4,opt,name=speech_event_type,json=speechEventType,proto3,enum=api.StreamingRecognizeResponse_SpeechEventType" json:"speech_event_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                   `json:"-"`
	XXX_unrecognized     []byte                                     `json:"-"`
	XXX_sizecache        int32                                      `json:"-"`
}

func (m *StreamingRecognizeResponse) Reset()         { *m = StreamingRecognizeResponse{} }
func (m *StreamingRecognizeResponse) String() string { return proto.CompactTextString(m) }
func (*StreamingRecognizeResponse) ProtoMessage()    {}
func (*StreamingRecognizeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_339644b7ccc8fec8, []int{3}
}
func (m *StreamingRecognizeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamingRecognizeResponse.Unmarshal(m, b)
}
func (m *StreamingRecognizeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamingRecognizeResponse.Marshal(b, m, deterministic)
}
func (dst *StreamingRecognizeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamingRecognizeResponse.Merge(dst, src)
}
func (m *StreamingRecognizeResponse) XXX_Size() int {
	return xxx_messageInfo_StreamingRecognizeResponse.Size(m)
}
func (m *StreamingRecognizeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamingRecognizeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StreamingRecognizeResponse proto.InternalMessageInfo

func (m *StreamingRecognizeResponse) GetResults() []*StreamingRecognitionResult {
	if m != nil {
		return m.Results
	}
	return nil
}

func (m *StreamingRecognizeResponse) GetSpeechEventType() StreamingRecognizeResponse_SpeechEventType {
	if m != nil {
		return m.SpeechEventType
	}
	return StreamingRecognizeResponse_SPEECH_EVENT_UNSPECIFIED
}

type StreamingRecognitionResult struct {
	Alternatives         []*SpeechRecognitionAlternative `protobuf:"bytes,1,rep,name=alternatives,proto3" json:"alternatives,omitempty"`
	IsFinal              bool                            `protobuf:"varint,2,opt,name=is_final,json=isFinal,proto3" json:"is_final,omitempty"`
	Stability            float32                         `protobuf:"fixed32,3,opt,name=stability,proto3" json:"stability,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *StreamingRecognitionResult) Reset()         { *m = StreamingRecognitionResult{} }
func (m *StreamingRecognitionResult) String() string { return proto.CompactTextString(m) }
func (*StreamingRecognitionResult) ProtoMessage()    {}
func (*StreamingRecognitionResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_339644b7ccc8fec8, []int{4}
}
func (m *StreamingRecognitionResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamingRecognitionResult.Unmarshal(m, b)
}
func (m *StreamingRecognitionResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamingRecognitionResult.Marshal(b, m, deterministic)
}
func (dst *StreamingRecognitionResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamingRecognitionResult.Merge(dst, src)
}
func (m *StreamingRecognitionResult) XXX_Size() int {
	return xxx_messageInfo_StreamingRecognitionResult.Size(m)
}
func (m *StreamingRecognitionResult) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamingRecognitionResult.DiscardUnknown(m)
}

var xxx_messageInfo_StreamingRecognitionResult proto.InternalMessageInfo

func (m *StreamingRecognitionResult) GetAlternatives() []*SpeechRecognitionAlternative {
	if m != nil {
		return m.Alternatives
	}
	return nil
}

func (m *StreamingRecognitionResult) GetIsFinal() bool {
	if m != nil {
		return m.IsFinal
	}
	return false
}

func (m *StreamingRecognitionResult) GetStability() float32 {
	if m != nil {
		return m.Stability
	}
	return 0
}

type SpeechRecognitionAlternative struct {
	Transcript           string      `protobuf:"bytes,1,opt,name=transcript,proto3" json:"transcript,omitempty"`
	Confidence           float32     `protobuf:"fixed32,2,opt,name=confidence,proto3" json:"confidence,omitempty"`
	Words                []*WordInfo `protobuf:"bytes,3,rep,name=words,proto3" json:"words,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *SpeechRecognitionAlternative) Reset()         { *m = SpeechRecognitionAlternative{} }
func (m *SpeechRecognitionAlternative) String() string { return proto.CompactTextString(m) }
func (*SpeechRecognitionAlternative) ProtoMessage()    {}
func (*SpeechRecognitionAlternative) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_339644b7ccc8fec8, []int{5}
}
func (m *SpeechRecognitionAlternative) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpeechRecognitionAlternative.Unmarshal(m, b)
}
func (m *SpeechRecognitionAlternative) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpeechRecognitionAlternative.Marshal(b, m, deterministic)
}
func (dst *SpeechRecognitionAlternative) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpeechRecognitionAlternative.Merge(dst, src)
}
func (m *SpeechRecognitionAlternative) XXX_Size() int {
	return xxx_messageInfo_SpeechRecognitionAlternative.Size(m)
}
func (m *SpeechRecognitionAlternative) XXX_DiscardUnknown() {
	xxx_messageInfo_SpeechRecognitionAlternative.DiscardUnknown(m)
}

var xxx_messageInfo_SpeechRecognitionAlternative proto.InternalMessageInfo

func (m *SpeechRecognitionAlternative) GetTranscript() string {
	if m != nil {
		return m.Transcript
	}
	return ""
}

func (m *SpeechRecognitionAlternative) GetConfidence() float32 {
	if m != nil {
		return m.Confidence
	}
	return 0
}

func (m *SpeechRecognitionAlternative) GetWords() []*WordInfo {
	if m != nil {
		return m.Words
	}
	return nil
}

type WordInfo struct {
	StartTime            *duration.Duration `protobuf:"bytes,1,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime              *duration.Duration `protobuf:"bytes,2,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	Word                 string             `protobuf:"bytes,3,opt,name=word,proto3" json:"word,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *WordInfo) Reset()         { *m = WordInfo{} }
func (m *WordInfo) String() string { return proto.CompactTextString(m) }
func (*WordInfo) ProtoMessage()    {}
func (*WordInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_339644b7ccc8fec8, []int{6}
}
func (m *WordInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WordInfo.Unmarshal(m, b)
}
func (m *WordInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WordInfo.Marshal(b, m, deterministic)
}
func (dst *WordInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WordInfo.Merge(dst, src)
}
func (m *WordInfo) XXX_Size() int {
	return xxx_messageInfo_WordInfo.Size(m)
}
func (m *WordInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_WordInfo.DiscardUnknown(m)
}

var xxx_messageInfo_WordInfo proto.InternalMessageInfo

func (m *WordInfo) GetStartTime() *duration.Duration {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *WordInfo) GetEndTime() *duration.Duration {
	if m != nil {
		return m.EndTime
	}
	return nil
}

func (m *WordInfo) GetWord() string {
	if m != nil {
		return m.Word
	}
	return ""
}

func init() {
	proto.RegisterType((*TTSRequest)(nil), "api.TTSRequest")
	proto.RegisterType((*TTSReply)(nil), "api.TTSReply")
	proto.RegisterType((*STTRequest)(nil), "api.STTRequest")
	proto.RegisterType((*StreamingRecognizeResponse)(nil), "api.StreamingRecognizeResponse")
	proto.RegisterType((*StreamingRecognitionResult)(nil), "api.StreamingRecognitionResult")
	proto.RegisterType((*SpeechRecognitionAlternative)(nil), "api.SpeechRecognitionAlternative")
	proto.RegisterType((*WordInfo)(nil), "api.WordInfo")
	proto.RegisterEnum("api.StreamingRecognizeResponse_SpeechEventType", StreamingRecognizeResponse_SpeechEventType_name, StreamingRecognizeResponse_SpeechEventType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// OpenVAServiceClient is the client API for OpenVAService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OpenVAServiceClient interface {
	TTSStringToMP3(ctx context.Context, in *TTSRequest, opts ...grpc.CallOption) (*TTSReply, error)
	STT(ctx context.Context, in *STTRequest, opts ...grpc.CallOption) (OpenVAService_STTClient, error)
}

type openVAServiceClient struct {
	cc *grpc.ClientConn
}

func NewOpenVAServiceClient(cc *grpc.ClientConn) OpenVAServiceClient {
	return &openVAServiceClient{cc}
}

func (c *openVAServiceClient) TTSStringToMP3(ctx context.Context, in *TTSRequest, opts ...grpc.CallOption) (*TTSReply, error) {
	out := new(TTSReply)
	err := c.cc.Invoke(ctx, "/api.OpenVAService/TTSStringToMP3", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *openVAServiceClient) STT(ctx context.Context, in *STTRequest, opts ...grpc.CallOption) (OpenVAService_STTClient, error) {
	stream, err := c.cc.NewStream(ctx, &_OpenVAService_serviceDesc.Streams[0], "/api.OpenVAService/STT", opts...)
	if err != nil {
		return nil, err
	}
	x := &openVAServiceSTTClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type OpenVAService_STTClient interface {
	Recv() (*StreamingRecognizeResponse, error)
	grpc.ClientStream
}

type openVAServiceSTTClient struct {
	grpc.ClientStream
}

func (x *openVAServiceSTTClient) Recv() (*StreamingRecognizeResponse, error) {
	m := new(StreamingRecognizeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OpenVAServiceServer is the server API for OpenVAService service.
type OpenVAServiceServer interface {
	TTSStringToMP3(context.Context, *TTSRequest) (*TTSReply, error)
	STT(*STTRequest, OpenVAService_STTServer) error
}

func RegisterOpenVAServiceServer(s *grpc.Server, srv OpenVAServiceServer) {
	s.RegisterService(&_OpenVAService_serviceDesc, srv)
}

func _OpenVAService_TTSStringToMP3_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TTSRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OpenVAServiceServer).TTSStringToMP3(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.OpenVAService/TTSStringToMP3",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpenVAServiceServer).TTSStringToMP3(ctx, req.(*TTSRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OpenVAService_STT_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(STTRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OpenVAServiceServer).STT(m, &openVAServiceSTTServer{stream})
}

type OpenVAService_STTServer interface {
	Send(*StreamingRecognizeResponse) error
	grpc.ServerStream
}

type openVAServiceSTTServer struct {
	grpc.ServerStream
}

func (x *openVAServiceSTTServer) Send(m *StreamingRecognizeResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _OpenVAService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.OpenVAService",
	HandlerType: (*OpenVAServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TTSStringToMP3",
			Handler:    _OpenVAService_TTSStringToMP3_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "STT",
			Handler:       _OpenVAService_STT_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "service.proto",
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_service_339644b7ccc8fec8) }

var fileDescriptor_service_339644b7ccc8fec8 = []byte{
	// 574 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0xc1, 0x4f, 0x13, 0x4f,
	0x14, 0x66, 0x5b, 0x7e, 0x3f, 0xca, 0x83, 0x02, 0xce, 0xc5, 0xa5, 0x12, 0x58, 0xd7, 0x4b, 0x63,
	0xcc, 0x62, 0x8a, 0x07, 0x8d, 0x27, 0x84, 0x45, 0x9b, 0x40, 0x69, 0x66, 0x06, 0x3c, 0x78, 0xd8,
	0x2c, 0xed, 0x6b, 0x9d, 0x64, 0x99, 0x59, 0x67, 0xa6, 0x68, 0x3d, 0x19, 0xbd, 0xfb, 0x17, 0xf8,
	0xc7, 0x9a, 0xce, 0x76, 0xb3, 0x05, 0x8d, 0xdc, 0x66, 0xbf, 0xf9, 0xbe, 0x79, 0xef, 0x7b, 0xfb,
	0x3d, 0x68, 0x1a, 0xd4, 0x37, 0x62, 0x80, 0x51, 0xae, 0x95, 0x55, 0xa4, 0x9e, 0xe6, 0xa2, 0xb5,
	0x3b, 0x56, 0x6a, 0x9c, 0xe1, 0xbe, 0x83, 0xae, 0x26, 0xa3, 0xfd, 0xe1, 0x44, 0xa7, 0x56, 0x28,
	0x59, 0x90, 0xc2, 0x00, 0x80, 0x73, 0x46, 0xf1, 0xd3, 0x04, 0x8d, 0x25, 0x04, 0x96, 0x2d, 0x7e,
	0xb1, 0xbe, 0x17, 0x78, 0xed, 0x55, 0xea, 0xce, 0xe1, 0x33, 0x68, 0x38, 0x46, 0x9e, 0x4d, 0x49,
	0x00, 0x6b, 0x67, 0xfd, 0x03, 0x8a, 0x26, 0x57, 0xd2, 0xa0, 0xa3, 0xad, 0xd3, 0x45, 0x28, 0x7c,
	0x0a, 0xc0, 0x38, 0x2f, 0xdf, 0xdb, 0x81, 0x55, 0xc6, 0xf9, 0x9b, 0xc9, 0x68, 0x84, 0x7a, 0xce,
	0xae, 0x80, 0xf0, 0x7b, 0x0d, 0x5a, 0xcc, 0x6a, 0x4c, 0xaf, 0x85, 0x1c, 0x53, 0x1c, 0xa8, 0xb1,
	0x14, 0x5f, 0xb1, 0x7c, 0x8a, 0xbc, 0x82, 0x15, 0x8d, 0x66, 0x92, 0x59, 0xe3, 0xd7, 0x82, 0x7a,
	0x7b, 0xad, 0xb3, 0x17, 0xa5, 0xb9, 0x88, 0xee, 0x2a, 0x66, 0x66, 0xa8, 0xe3, 0xd1, 0x92, 0x4f,
	0x3e, 0xc0, 0x03, 0x93, 0x23, 0x0e, 0x3e, 0x26, 0x78, 0x83, 0xd2, 0x26, 0x76, 0x9a, 0xa3, 0xbf,
	0x1c, 0x78, 0xed, 0x8d, 0xce, 0xfe, 0x5f, 0x1f, 0xa9, 0xca, 0x46, 0xcc, 0x09, 0xe3, 0x99, 0x8e,
	0x4f, 0x73, 0xa4, 0x9b, 0xe6, 0x36, 0x10, 0x9e, 0xc2, 0xe6, 0x1d, 0x0e, 0xd9, 0x01, 0x9f, 0xf5,
	0xe3, 0xf8, 0xe8, 0x5d, 0x12, 0x5f, 0xc6, 0x3d, 0x9e, 0x5c, 0xf4, 0x58, 0x3f, 0x3e, 0xea, 0x9e,
	0x74, 0xe3, 0xe3, 0xad, 0x25, 0xf2, 0x08, 0x1e, 0xc6, 0xbd, 0xe3, 0xe4, 0xfc, 0x24, 0x61, 0xdd,
	0xde, 0xdb, 0xd3, 0x38, 0xb9, 0xe0, 0x3c, 0xa6, 0x87, 0xbd, 0xa3, 0x78, 0xcb, 0x0b, 0x7f, 0x79,
	0x7f, 0x0e, 0xa1, 0xb2, 0x44, 0x62, 0x58, 0x4f, 0x33, 0x8b, 0x5a, 0xa6, 0x56, 0xdc, 0xa0, 0xf1,
	0x3d, 0x37, 0x89, 0xc7, 0x85, 0x09, 0xd7, 0xc5, 0x82, 0xe6, 0xb0, 0x62, 0xd2, 0x5b, 0x32, 0xb2,
	0x0d, 0x0d, 0x61, 0x92, 0x91, 0x90, 0x69, 0xe6, 0xd7, 0x02, 0xaf, 0xdd, 0xa0, 0x2b, 0xc2, 0x9c,
	0xcc, 0x3e, 0x67, 0xff, 0xc8, 0xd8, 0xf4, 0x4a, 0x64, 0xc2, 0x4e, 0xfd, 0x7a, 0xe0, 0xb5, 0x6b,
	0xb4, 0x02, 0xc2, 0x1f, 0x1e, 0xec, 0xfc, 0xab, 0x0e, 0xd9, 0x05, 0xb0, 0x3a, 0x95, 0x66, 0xa0,
	0x45, 0x5e, 0x06, 0x67, 0x01, 0x99, 0xdd, 0x0f, 0x94, 0x1c, 0x89, 0x21, 0xca, 0x01, 0xba, 0xda,
	0x35, 0xba, 0x80, 0x90, 0x27, 0xf0, 0xdf, 0x67, 0xa5, 0x87, 0xc6, 0xaf, 0x3b, 0x67, 0x4d, 0xe7,
	0xec, 0xbd, 0xd2, 0xc3, 0xae, 0x1c, 0x29, 0x5a, 0xdc, 0x85, 0x3f, 0x3d, 0x68, 0x94, 0x18, 0x79,
	0x09, 0x60, 0x6c, 0xaa, 0x6d, 0x62, 0xc5, 0x75, 0x91, 0xc1, 0xb5, 0xce, 0x76, 0x54, 0xe4, 0x3c,
	0x2a, 0x73, 0x1e, 0x1d, 0xcf, 0x73, 0xee, 0xcc, 0x68, 0xcb, 0xc5, 0x35, 0x92, 0x17, 0xd0, 0x40,
	0x39, 0x2c, 0x74, 0xb5, 0xfb, 0x74, 0x2b, 0x28, 0x87, 0x4e, 0x45, 0x60, 0x79, 0xd6, 0x85, 0x9b,
	0xcd, 0x2a, 0x75, 0xe7, 0xce, 0x37, 0x0f, 0x9a, 0xe7, 0x39, 0xca, 0xcb, 0x43, 0x56, 0xec, 0x1c,
	0xe9, 0xc0, 0x06, 0xe7, 0x8c, 0x59, 0x2d, 0xe4, 0x98, 0xab, 0xb3, 0xfe, 0x01, 0xd9, 0x74, 0x56,
	0xaa, 0xed, 0x6a, 0x35, 0x2b, 0x20, 0xcf, 0xa6, 0xe1, 0x12, 0x79, 0x0d, 0x75, 0xc6, 0xf9, 0x9c,
	0x58, 0xad, 0x4d, 0x6b, 0xef, 0x9e, 0x8c, 0x86, 0x4b, 0xcf, 0xbd, 0xab, 0xff, 0x5d, 0xcb, 0x07,
	0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xbc, 0xfb, 0x46, 0x0f, 0xf6, 0x03, 0x00, 0x00,
}