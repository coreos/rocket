// Code generated by protoc-gen-go.
// source: test.proto
// DO NOT EDIT!

/*
Package grpc_testing is a generated protocol buffer package.

It is generated from these files:
	test.proto

It has these top-level messages:
	StatsRequest
	ServerStats
	Payload
	HistogramData
	ClientConfig
	Mark
	ClientArgs
	ClientStats
	ClientStatus
	ServerConfig
	ServerArgs
	ServerStatus
	SimpleRequest
	SimpleResponse
*/
package grpc_testing

import proto "github.com/coreos/rkt/Godeps/_workspace/src/github.com/golang/protobuf/proto"

import (
	context "github.com/coreos/rkt/Godeps/_workspace/src/golang.org/x/net/context"
	grpc "github.com/coreos/rkt/Godeps/_workspace/src/google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type PayloadType int32

const (
	// Compressable text format.
	PayloadType_COMPRESSABLE PayloadType = 0
	// Uncompressable binary format.
	PayloadType_UNCOMPRESSABLE PayloadType = 1
	// Randomly chosen from all other formats defined in this enum.
	PayloadType_RANDOM PayloadType = 2
)

var PayloadType_name = map[int32]string{
	0: "COMPRESSABLE",
	1: "UNCOMPRESSABLE",
	2: "RANDOM",
}
var PayloadType_value = map[string]int32{
	"COMPRESSABLE":   0,
	"UNCOMPRESSABLE": 1,
	"RANDOM":         2,
}

func (x PayloadType) String() string {
	return proto.EnumName(PayloadType_name, int32(x))
}

type ClientType int32

const (
	ClientType_SYNCHRONOUS_CLIENT ClientType = 0
	ClientType_ASYNC_CLIENT       ClientType = 1
)

var ClientType_name = map[int32]string{
	0: "SYNCHRONOUS_CLIENT",
	1: "ASYNC_CLIENT",
}
var ClientType_value = map[string]int32{
	"SYNCHRONOUS_CLIENT": 0,
	"ASYNC_CLIENT":       1,
}

func (x ClientType) String() string {
	return proto.EnumName(ClientType_name, int32(x))
}

type ServerType int32

const (
	ServerType_SYNCHRONOUS_SERVER ServerType = 0
	ServerType_ASYNC_SERVER       ServerType = 1
)

var ServerType_name = map[int32]string{
	0: "SYNCHRONOUS_SERVER",
	1: "ASYNC_SERVER",
}
var ServerType_value = map[string]int32{
	"SYNCHRONOUS_SERVER": 0,
	"ASYNC_SERVER":       1,
}

func (x ServerType) String() string {
	return proto.EnumName(ServerType_name, int32(x))
}

type RpcType int32

const (
	RpcType_UNARY     RpcType = 0
	RpcType_STREAMING RpcType = 1
)

var RpcType_name = map[int32]string{
	0: "UNARY",
	1: "STREAMING",
}
var RpcType_value = map[string]int32{
	"UNARY":     0,
	"STREAMING": 1,
}

func (x RpcType) String() string {
	return proto.EnumName(RpcType_name, int32(x))
}

type StatsRequest struct {
	// run number
	TestNum int32 `protobuf:"varint,1,opt,name=test_num" json:"test_num,omitempty"`
}

func (m *StatsRequest) Reset()         { *m = StatsRequest{} }
func (m *StatsRequest) String() string { return proto.CompactTextString(m) }
func (*StatsRequest) ProtoMessage()    {}

type ServerStats struct {
	// wall clock time
	TimeElapsed float64 `protobuf:"fixed64,1,opt,name=time_elapsed" json:"time_elapsed,omitempty"`
	// user time used by the server process and threads
	TimeUser float64 `protobuf:"fixed64,2,opt,name=time_user" json:"time_user,omitempty"`
	// server time used by the server process and all threads
	TimeSystem float64 `protobuf:"fixed64,3,opt,name=time_system" json:"time_system,omitempty"`
}

func (m *ServerStats) Reset()         { *m = ServerStats{} }
func (m *ServerStats) String() string { return proto.CompactTextString(m) }
func (*ServerStats) ProtoMessage()    {}

type Payload struct {
	// The type of data in body.
	Type PayloadType `protobuf:"varint,1,opt,name=type,enum=grpc.testing.PayloadType" json:"type,omitempty"`
	// Primary contents of payload.
	Body []byte `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (m *Payload) Reset()         { *m = Payload{} }
func (m *Payload) String() string { return proto.CompactTextString(m) }
func (*Payload) ProtoMessage()    {}

type HistogramData struct {
	Bucket       []uint32 `protobuf:"varint,1,rep,name=bucket" json:"bucket,omitempty"`
	MinSeen      float64  `protobuf:"fixed64,2,opt,name=min_seen" json:"min_seen,omitempty"`
	MaxSeen      float64  `protobuf:"fixed64,3,opt,name=max_seen" json:"max_seen,omitempty"`
	Sum          float64  `protobuf:"fixed64,4,opt,name=sum" json:"sum,omitempty"`
	SumOfSquares float64  `protobuf:"fixed64,5,opt,name=sum_of_squares" json:"sum_of_squares,omitempty"`
	Count        float64  `protobuf:"fixed64,6,opt,name=count" json:"count,omitempty"`
}

func (m *HistogramData) Reset()         { *m = HistogramData{} }
func (m *HistogramData) String() string { return proto.CompactTextString(m) }
func (*HistogramData) ProtoMessage()    {}

type ClientConfig struct {
	ServerTargets             []string   `protobuf:"bytes,1,rep,name=server_targets" json:"server_targets,omitempty"`
	ClientType                ClientType `protobuf:"varint,2,opt,name=client_type,enum=grpc.testing.ClientType" json:"client_type,omitempty"`
	EnableSsl                 bool       `protobuf:"varint,3,opt,name=enable_ssl" json:"enable_ssl,omitempty"`
	OutstandingRpcsPerChannel int32      `protobuf:"varint,4,opt,name=outstanding_rpcs_per_channel" json:"outstanding_rpcs_per_channel,omitempty"`
	ClientChannels            int32      `protobuf:"varint,5,opt,name=client_channels" json:"client_channels,omitempty"`
	PayloadSize               int32      `protobuf:"varint,6,opt,name=payload_size" json:"payload_size,omitempty"`
	// only for async client:
	AsyncClientThreads int32   `protobuf:"varint,7,opt,name=async_client_threads" json:"async_client_threads,omitempty"`
	RpcType            RpcType `protobuf:"varint,8,opt,name=rpc_type,enum=grpc.testing.RpcType" json:"rpc_type,omitempty"`
}

func (m *ClientConfig) Reset()         { *m = ClientConfig{} }
func (m *ClientConfig) String() string { return proto.CompactTextString(m) }
func (*ClientConfig) ProtoMessage()    {}

// Request current stats
type Mark struct {
}

func (m *Mark) Reset()         { *m = Mark{} }
func (m *Mark) String() string { return proto.CompactTextString(m) }
func (*Mark) ProtoMessage()    {}

type ClientArgs struct {
	Setup *ClientConfig `protobuf:"bytes,1,opt,name=setup" json:"setup,omitempty"`
	Mark  *Mark         `protobuf:"bytes,2,opt,name=mark" json:"mark,omitempty"`
}

func (m *ClientArgs) Reset()         { *m = ClientArgs{} }
func (m *ClientArgs) String() string { return proto.CompactTextString(m) }
func (*ClientArgs) ProtoMessage()    {}

func (m *ClientArgs) GetSetup() *ClientConfig {
	if m != nil {
		return m.Setup
	}
	return nil
}

func (m *ClientArgs) GetMark() *Mark {
	if m != nil {
		return m.Mark
	}
	return nil
}

type ClientStats struct {
	Latencies   *HistogramData `protobuf:"bytes,1,opt,name=latencies" json:"latencies,omitempty"`
	TimeElapsed float64        `protobuf:"fixed64,3,opt,name=time_elapsed" json:"time_elapsed,omitempty"`
	TimeUser    float64        `protobuf:"fixed64,4,opt,name=time_user" json:"time_user,omitempty"`
	TimeSystem  float64        `protobuf:"fixed64,5,opt,name=time_system" json:"time_system,omitempty"`
}

func (m *ClientStats) Reset()         { *m = ClientStats{} }
func (m *ClientStats) String() string { return proto.CompactTextString(m) }
func (*ClientStats) ProtoMessage()    {}

func (m *ClientStats) GetLatencies() *HistogramData {
	if m != nil {
		return m.Latencies
	}
	return nil
}

type ClientStatus struct {
	Stats *ClientStats `protobuf:"bytes,1,opt,name=stats" json:"stats,omitempty"`
}

func (m *ClientStatus) Reset()         { *m = ClientStatus{} }
func (m *ClientStatus) String() string { return proto.CompactTextString(m) }
func (*ClientStatus) ProtoMessage()    {}

func (m *ClientStatus) GetStats() *ClientStats {
	if m != nil {
		return m.Stats
	}
	return nil
}

type ServerConfig struct {
	ServerType ServerType `protobuf:"varint,1,opt,name=server_type,enum=grpc.testing.ServerType" json:"server_type,omitempty"`
	Threads    int32      `protobuf:"varint,2,opt,name=threads" json:"threads,omitempty"`
	EnableSsl  bool       `protobuf:"varint,3,opt,name=enable_ssl" json:"enable_ssl,omitempty"`
}

func (m *ServerConfig) Reset()         { *m = ServerConfig{} }
func (m *ServerConfig) String() string { return proto.CompactTextString(m) }
func (*ServerConfig) ProtoMessage()    {}

type ServerArgs struct {
	Setup *ServerConfig `protobuf:"bytes,1,opt,name=setup" json:"setup,omitempty"`
	Mark  *Mark         `protobuf:"bytes,2,opt,name=mark" json:"mark,omitempty"`
}

func (m *ServerArgs) Reset()         { *m = ServerArgs{} }
func (m *ServerArgs) String() string { return proto.CompactTextString(m) }
func (*ServerArgs) ProtoMessage()    {}

func (m *ServerArgs) GetSetup() *ServerConfig {
	if m != nil {
		return m.Setup
	}
	return nil
}

func (m *ServerArgs) GetMark() *Mark {
	if m != nil {
		return m.Mark
	}
	return nil
}

type ServerStatus struct {
	Stats *ServerStats `protobuf:"bytes,1,opt,name=stats" json:"stats,omitempty"`
	Port  int32        `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
}

func (m *ServerStatus) Reset()         { *m = ServerStatus{} }
func (m *ServerStatus) String() string { return proto.CompactTextString(m) }
func (*ServerStatus) ProtoMessage()    {}

func (m *ServerStatus) GetStats() *ServerStats {
	if m != nil {
		return m.Stats
	}
	return nil
}

type SimpleRequest struct {
	// Desired payload type in the response from the server.
	// If response_type is RANDOM, server randomly chooses one from other formats.
	ResponseType PayloadType `protobuf:"varint,1,opt,name=response_type,enum=grpc.testing.PayloadType" json:"response_type,omitempty"`
	// Desired payload size in the response from the server.
	// If response_type is COMPRESSABLE, this denotes the size before compression.
	ResponseSize int32 `protobuf:"varint,2,opt,name=response_size" json:"response_size,omitempty"`
	// Optional input payload sent along with the request.
	Payload *Payload `protobuf:"bytes,3,opt,name=payload" json:"payload,omitempty"`
}

func (m *SimpleRequest) Reset()         { *m = SimpleRequest{} }
func (m *SimpleRequest) String() string { return proto.CompactTextString(m) }
func (*SimpleRequest) ProtoMessage()    {}

func (m *SimpleRequest) GetPayload() *Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

type SimpleResponse struct {
	Payload *Payload `protobuf:"bytes,1,opt,name=payload" json:"payload,omitempty"`
}

func (m *SimpleResponse) Reset()         { *m = SimpleResponse{} }
func (m *SimpleResponse) String() string { return proto.CompactTextString(m) }
func (*SimpleResponse) ProtoMessage()    {}

func (m *SimpleResponse) GetPayload() *Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func init() {
	proto.RegisterEnum("grpc.testing.PayloadType", PayloadType_name, PayloadType_value)
	proto.RegisterEnum("grpc.testing.ClientType", ClientType_name, ClientType_value)
	proto.RegisterEnum("grpc.testing.ServerType", ServerType_name, ServerType_value)
	proto.RegisterEnum("grpc.testing.RpcType", RpcType_name, RpcType_value)
}

// Client API for TestService service

type TestServiceClient interface {
	// One request followed by one response.
	// The server returns the client payload as-is.
	UnaryCall(ctx context.Context, in *SimpleRequest, opts ...grpc.CallOption) (*SimpleResponse, error)
	// One request followed by one response.
	// The server returns the client payload as-is.
	StreamingCall(ctx context.Context, opts ...grpc.CallOption) (TestService_StreamingCallClient, error)
}

type testServiceClient struct {
	cc *grpc.ClientConn
}

func NewTestServiceClient(cc *grpc.ClientConn) TestServiceClient {
	return &testServiceClient{cc}
}

func (c *testServiceClient) UnaryCall(ctx context.Context, in *SimpleRequest, opts ...grpc.CallOption) (*SimpleResponse, error) {
	out := new(SimpleResponse)
	err := grpc.Invoke(ctx, "/grpc.testing.TestService/UnaryCall", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) StreamingCall(ctx context.Context, opts ...grpc.CallOption) (TestService_StreamingCallClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_TestService_serviceDesc.Streams[0], c.cc, "/grpc.testing.TestService/StreamingCall", opts...)
	if err != nil {
		return nil, err
	}
	x := &testServiceStreamingCallClient{stream}
	return x, nil
}

type TestService_StreamingCallClient interface {
	Send(*SimpleRequest) error
	Recv() (*SimpleResponse, error)
	grpc.ClientStream
}

type testServiceStreamingCallClient struct {
	grpc.ClientStream
}

func (x *testServiceStreamingCallClient) Send(m *SimpleRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *testServiceStreamingCallClient) Recv() (*SimpleResponse, error) {
	m := new(SimpleResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for TestService service

type TestServiceServer interface {
	// One request followed by one response.
	// The server returns the client payload as-is.
	UnaryCall(context.Context, *SimpleRequest) (*SimpleResponse, error)
	// One request followed by one response.
	// The server returns the client payload as-is.
	StreamingCall(TestService_StreamingCallServer) error
}

func RegisterTestServiceServer(s *grpc.Server, srv TestServiceServer) {
	s.RegisterService(&_TestService_serviceDesc, srv)
}

func _TestService_UnaryCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(SimpleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(TestServiceServer).UnaryCall(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _TestService_StreamingCall_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TestServiceServer).StreamingCall(&testServiceStreamingCallServer{stream})
}

type TestService_StreamingCallServer interface {
	Send(*SimpleResponse) error
	Recv() (*SimpleRequest, error)
	grpc.ServerStream
}

type testServiceStreamingCallServer struct {
	grpc.ServerStream
}

func (x *testServiceStreamingCallServer) Send(m *SimpleResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *testServiceStreamingCallServer) Recv() (*SimpleRequest, error) {
	m := new(SimpleRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _TestService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.testing.TestService",
	HandlerType: (*TestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UnaryCall",
			Handler:    _TestService_UnaryCall_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamingCall",
			Handler:       _TestService_StreamingCall_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
}

// Client API for Worker service

type WorkerClient interface {
	// Start test with specified workload
	RunTest(ctx context.Context, opts ...grpc.CallOption) (Worker_RunTestClient, error)
	// Start test with specified workload
	RunServer(ctx context.Context, opts ...grpc.CallOption) (Worker_RunServerClient, error)
}

type workerClient struct {
	cc *grpc.ClientConn
}

func NewWorkerClient(cc *grpc.ClientConn) WorkerClient {
	return &workerClient{cc}
}

func (c *workerClient) RunTest(ctx context.Context, opts ...grpc.CallOption) (Worker_RunTestClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Worker_serviceDesc.Streams[0], c.cc, "/grpc.testing.Worker/RunTest", opts...)
	if err != nil {
		return nil, err
	}
	x := &workerRunTestClient{stream}
	return x, nil
}

type Worker_RunTestClient interface {
	Send(*ClientArgs) error
	Recv() (*ClientStatus, error)
	grpc.ClientStream
}

type workerRunTestClient struct {
	grpc.ClientStream
}

func (x *workerRunTestClient) Send(m *ClientArgs) error {
	return x.ClientStream.SendMsg(m)
}

func (x *workerRunTestClient) Recv() (*ClientStatus, error) {
	m := new(ClientStatus)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *workerClient) RunServer(ctx context.Context, opts ...grpc.CallOption) (Worker_RunServerClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Worker_serviceDesc.Streams[1], c.cc, "/grpc.testing.Worker/RunServer", opts...)
	if err != nil {
		return nil, err
	}
	x := &workerRunServerClient{stream}
	return x, nil
}

type Worker_RunServerClient interface {
	Send(*ServerArgs) error
	Recv() (*ServerStatus, error)
	grpc.ClientStream
}

type workerRunServerClient struct {
	grpc.ClientStream
}

func (x *workerRunServerClient) Send(m *ServerArgs) error {
	return x.ClientStream.SendMsg(m)
}

func (x *workerRunServerClient) Recv() (*ServerStatus, error) {
	m := new(ServerStatus)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Worker service

type WorkerServer interface {
	// Start test with specified workload
	RunTest(Worker_RunTestServer) error
	// Start test with specified workload
	RunServer(Worker_RunServerServer) error
}

func RegisterWorkerServer(s *grpc.Server, srv WorkerServer) {
	s.RegisterService(&_Worker_serviceDesc, srv)
}

func _Worker_RunTest_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(WorkerServer).RunTest(&workerRunTestServer{stream})
}

type Worker_RunTestServer interface {
	Send(*ClientStatus) error
	Recv() (*ClientArgs, error)
	grpc.ServerStream
}

type workerRunTestServer struct {
	grpc.ServerStream
}

func (x *workerRunTestServer) Send(m *ClientStatus) error {
	return x.ServerStream.SendMsg(m)
}

func (x *workerRunTestServer) Recv() (*ClientArgs, error) {
	m := new(ClientArgs)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Worker_RunServer_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(WorkerServer).RunServer(&workerRunServerServer{stream})
}

type Worker_RunServerServer interface {
	Send(*ServerStatus) error
	Recv() (*ServerArgs, error)
	grpc.ServerStream
}

type workerRunServerServer struct {
	grpc.ServerStream
}

func (x *workerRunServerServer) Send(m *ServerStatus) error {
	return x.ServerStream.SendMsg(m)
}

func (x *workerRunServerServer) Recv() (*ServerArgs, error) {
	m := new(ServerArgs)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Worker_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.testing.Worker",
	HandlerType: (*WorkerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RunTest",
			Handler:       _Worker_RunTest_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "RunServer",
			Handler:       _Worker_RunServer_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
}
