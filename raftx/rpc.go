package raftx

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/bytedance/sonic"
)

const ServerStopErr = "Server Stopped"

type RPC struct {
	Command      interface{}
	ResponseChan chan RPCResponse
}

type RPCResponse struct {
	Response interface{}
	Err      error
}

type RequestInterface interface {
	*VoteRequest | *AppendEntriesRequest
}

type ResponseInterface interface {
	*VoteResponse | *AppendEntriesResponse
}

func (r *RPC) Respond(resp interface{}, err error) {
	r.ResponseChan <- RPCResponse{resp, err}
}

// Transporter 负责通过网络发送 Request 和 Response, 可以基于 TCP、UDP、TLS、HTTP、GRPC 来实现
type Transporter interface {
	Start(port int, server *Server)
	AppendEntries(peer *Peer, req *AppendEntriesRequest) (*AppendEntriesResponse, error)
	RequestVote(peer *Peer, req *VoteRequest) (*VoteResponse, error)
}

type HttpTransporter struct {
	appendEntriesPath string
	votePath          string
	httpClient        http.Client
}

func NewHttpTransporter(prefix string, timeout time.Duration) *HttpTransporter {
	t := &HttpTransporter{
		appendEntriesPath: joinUrlPath(prefix, "/ae"),
		votePath:          joinUrlPath(prefix, "/vote"),
	}
	t.httpClient.Timeout = timeout
	return t
}

// Start 启动 HTTP Server
func (transporter *HttpTransporter) Start(port int, server *Server) {
	mux := http.NewServeMux()
	mux.Handle(transporter.appendEntriesPath, AppendEntriesHandler(server))
	mux.Handle(transporter.votePath, VoteHandler(server))
	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), mux); err != nil {
		panic(err)
	}
}

func (transporter *HttpTransporter) AppendEntries(peer *Peer, req *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	var resp AppendEntriesResponse
	err := post(transporter.httpClient, peer, transporter.appendEntriesPath, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (transporter *HttpTransporter) RequestVote(peer *Peer, req *VoteRequest) (*VoteResponse, error) {
	var resp VoteResponse
	err := post(transporter.httpClient, peer, transporter.votePath, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func post[E RequestInterface, S ResponseInterface](httpClient http.Client, peer *Peer, path string, req E, resp S) error {
	// 序列化
	bs, err := sonic.Marshal(req)
	if err != nil {
		slog.Error("Sonic Marshal Request Failed", "request", reflect.TypeOf(req).Elem().Name(), "error", err)
		return err
	}

	url := joinUrlPath(peer.ConnectionString, path)
	httpResp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		slog.Error("HttpClient Post Failed", "request", reflect.TypeOf(req).Elem().Name(), "error", err)
		return err
	}
	defer httpResp.Body.Close()

	// 读取 Response
	bs, err = io.ReadAll(httpResp.Body)
	if err != nil {
		slog.Error("Read Response Failed", "request", reflect.TypeOf(req).Elem().Name(), "error", err)
		return err
	}

	// 判断响应码
	if httpResp.StatusCode != http.StatusOK {
		respBody := string(bs)

		// 获取上一层调用者的方法名
		var methodName string
		funcName, _, _, ok := runtime.Caller(1)
		if ok {
			methodName = runtime.FuncForPC(funcName).Name()
		}

		// 判断是否为 ServerStopErr
		if respBody != ServerStopErr+"\n" {
			slog.Error(methodName+" Got Abnormal Code", "status", httpResp.Status, "msg", respBody)
		}

		return fmt.Errorf("%s Got Abnormal Code, status %s msg %s", methodName, httpResp.Status, respBody)
	}

	err = sonic.Unmarshal(bs, resp)
	if err != nil {
		slog.Error("Sonic Marshal Response Failed", "request", reflect.TypeOf(req).Elem().Name(), "error", err)
		return err
	}
	return nil
}

// AppendEntriesHandler 返回路由所需的 AppendEntriesHandlerFunc
func AppendEntriesHandler(server *Server) http.HandlerFunc {
	return handler(server, &AppendEntriesRequest{})
}

// VoteHandler 返回路由所需的 VoteHandlerFunc
func VoteHandler(server *Server) http.HandlerFunc {
	return handler(server, &VoteRequest{})
}

// 处理请求
func handler[E RequestInterface](server *Server, request E) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查 Raft 节点状态
		if server.GetState() == Stopped {
			http.Error(w, ServerStopErr, http.StatusInternalServerError)
			return
		}

		// 读取 Body
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid Json Request", http.StatusBadRequest)
			return
		}

		// 反序列化
		err = sonic.Unmarshal(bs, request)
		if err != nil {
			http.Error(w, "Sonic Unmarshal Failed", http.StatusInternalServerError)
			return
		}

		// 构造 RPC
		rpc := RPC{
			Command:      request,
			ResponseChan: make(chan RPCResponse),
		}

		// 再检查一次节点状态
		if server.GetState() == Stopped {
			http.Error(w, ServerStopErr, http.StatusInternalServerError)
			return
		}

		// 将 rpc 放入调用 Channel
		server.rpcCh <- rpc

		// 阻塞等待 RPC 的响应
		rpcResp := <-rpc.ResponseChan
		resp, err := rpcResp.Response, rpcResp.Err
		if err != nil || resp == nil {
			http.Error(w, "Server Failed", http.StatusInternalServerError)
			return
		}

		// 将 resp 序列化
		bs, err = sonic.Marshal(resp)
		if err != nil {
			http.Error(w, "Sonic Marshal Failed", http.StatusInternalServerError)
			return
		}

		w.Write(bs)
	}
}
