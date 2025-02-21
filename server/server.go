package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	pb "grpc_chat/grpc_chat/proto"
)

type server struct {
	pb.UnimplementedChatServiceServer
	users    map[string]string // username -> password
	clients  map[string]pb.ChatService_ChatServer
	mutex    sync.RWMutex
}

func newServer() *server {
	return &server{
		users:   make(map[string]string),
		clients: make(map[string]pb.ChatService_ChatServer),
	}
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[req.Username]; exists {
		return &pb.RegisterResponse{
			Success: false,
			Message: "用户名已存在",
		}, nil
	}

	s.users[req.Username] = req.Password
	return &pb.RegisterResponse{
		Success: true,
		Message: "注册成功",
	}, nil
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	password, exists := s.users[req.Username]
	if !exists || password != req.Password {
		return &pb.LoginResponse{
			Success: false,
			Message: "用户名或密码错误",
		}, nil
	}

	return &pb.LoginResponse{
		Success: true,
		Message: "登录成功",
		Token:   "dummy-token", // 在实际应用中应该生成真实的token
	}, nil
}

func (s *server) Chat(stream pb.ChatService_ChatServer) error {
	// 等待第一条消息以获取用户信息
	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	username := msg.Username

	// 注册客户端
	s.mutex.Lock()
	s.clients[username] = stream
	s.mutex.Unlock()

	// 清理函数
	defer func() {
		s.mutex.Lock()
		delete(s.clients, username)
		s.mutex.Unlock()
	}()

	// 广播用户加入消息
	s.broadcast(&pb.ChatMessage{
		Username:  "System",
		Content:   fmt.Sprintf("%s 加入了聊天室", username),
		Timestamp: time.Now().Unix(),
	})

	// 处理消息
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// 设置消息时间戳
		msg.Timestamp = time.Now().Unix()

		// 广播消息给所有客户端
		s.broadcast(msg)
	}
}

func (s *server) broadcast(msg *pb.ChatMessage) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, client := range s.clients {
		client.Send(msg)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, newServer())

	log.Printf("服务器启动在 :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}