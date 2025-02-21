package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	//"strings"
	"time"

	pb "grpc_chat/grpc_chat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client     pb.ChatServiceClient
	username   string
	token      string
	chatStream pb.ChatService_ChatClient
}

func newClient() *Client {
	conn, err := grpc.Dial("118.31.248.178:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接服务器失败: %v", err)
	}

	return &Client{
		client: pb.NewChatServiceClient(conn),
	}
}

func (c *Client) register(username, password string) error {
	resp, err := c.client.Register(context.Background(), &pb.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("注册失败: %v", err)
	}

	if !resp.Success {
		return fmt.Errorf("注册失败: %s", resp.Message)
	}

	fmt.Println("注册成功!")
	return nil
}

func (c *Client) login(username, password string) error {
	resp, err := c.client.Login(context.Background(), &pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("登录失败: %v", err)
	}

	if !resp.Success {
		return fmt.Errorf("登录失败: %s", resp.Message)
	}

	c.username = username
	c.token = resp.Token
	fmt.Println("登录成功!")
	return nil
}

func (c *Client) startChat() error {
	if c.username == "" || c.token == "" {
		return fmt.Errorf("请先登录")
	}

	stream, err := c.client.Chat(context.Background())
	if err != nil {
		return fmt.Errorf("建立聊天连接失败: %v", err)
	}
	c.chatStream = stream

	// 发送第一条消息以标识用户
	err = stream.Send(&pb.ChatMessage{
		Username:  c.username,
		Content:   "",
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("发送身份信息失败: %v", err)
	}

	// 启动接收消息的goroutine
	go c.receiveMessages()

	// 发送消息
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "/quit" {
			break
		}

		err := stream.Send(&pb.ChatMessage{
			Username:  c.username,
			Content:   msg,
			Timestamp: time.Now().Unix(),
		})
		if err != nil {
			log.Printf("发送消息失败: %v", err)
		}
	}

	return nil
}

func (c *Client) receiveMessages() {
	for {
		msg, err := c.chatStream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Printf("接收消息失败: %v", err)
			return
		}

		time := time.Unix(msg.Timestamp, 0).Format("15:04:05")
		fmt.Printf("[%s] %s: %s\n", time, msg.Username, msg.Content)
	}
}

func main() {
	client := newClient()

	for {
		fmt.Println("\n请选择操作:")
		fmt.Println("1. 注册")
		fmt.Println("2. 登录")
		fmt.Println("3. 进入聊天室")
		fmt.Println("4. 退出")

		var choice string
		fmt.Print("请输入选项: ")
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			var username, password string
			fmt.Print("请输入用户名: ")
			fmt.Scanln(&username)
			fmt.Print("请输入密码: ")
			fmt.Scanln(&password)

			err := client.register(username, password)
			if err != nil {
				fmt.Println(err)
			}

		case "2":
			var username, password string
			fmt.Print("请输入用户名: ")
			fmt.Scanln(&username)
			fmt.Print("请输入密码: ")
			fmt.Scanln(&password)

			err := client.login(username, password)
			if err != nil {
				fmt.Println(err)
			}

		case "3":
			fmt.Println("进入聊天室 (输入 /quit 退出)")
			err := client.startChat()
			if err != nil {
				fmt.Println(err)
			}

		case "4":
			fmt.Println("再见!")
			return

		default:
			fmt.Println("无效的选项，请重试")
		}
	}
}
