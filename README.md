# gRPC 聊天应用

这是一个基于 gRPC 的简单聊天应用，支持用户注册、登录和实时聊天功能。

## 项目架构

### 目录结构
```
.
├── proto/           # Protocol Buffers 定义文件
│   └── chat.proto   # 聊天服务接口定义
├── server/          # 服务端实现
│   └── server.go    # 服务端代码
├── client/          # 客户端实现
│   └── client.go    # 客户端代码
└── grpc_chat/       # 生成的 gRPC 代码
    └── proto/       # 编译后的 Protocol Buffers 代码
```

### 核心功能

1. **用户管理**
   - 用户注册：新用户可以注册账号
   - 用户登录：已注册用户可以登录系统

2. **聊天功能**
   - 实时消息：支持用户之间的实时消息交换
   - 系统通知：显示用户加入聊天室的系统消息
   - 时间戳：每条消息都带有发送时间

3. **通信机制**
   - 使用 gRPC 双向流式通信
   - 支持多客户端同时连接
   - 消息广播机制

## 使用流程

### 启动服务器

1. 进入项目目录：
   ```bash
   cd grpc_test
   ```

2. 启动服务器：
   ```bash
   go run server/server.go
   ```
   服务器将在 50051 端口启动

### 启动客户端

1. 在新的终端窗口中启动客户端：
   ```bash
   go run client/client.go
   ```

2. 客户端操作流程：
   - 选择 "1" 进行用户注册
   - 选择 "2" 进行用户登录
   - 选择 "3" 进入聊天室
   - 在聊天室中输入消息并按回车发送
   - 输入 "/quit" 退出聊天室
   - 选择 "4" 退出程序

## 测试指南

### 功能测试

1. **注册测试**
   - 测试注册新用户
   - 测试重复用户名注册

2. **登录测试**
   - 测试正确的用户名和密码
   - 测试错误的用户名或密码

3. **聊天测试**
   - 打开多个终端，运行多个客户端
   - 使用不同账号登录
   - 测试消息发送和接收
   - 测试用户加入提醒
   - 测试用户退出

### 测试步骤示例

1. **多用户聊天测试**
   ```bash
   # 终端 1
   go run client/client.go
   # 注册并登录用户1
   
   # 终端 2
   go run client/client.go
   # 注册并登录用户2
   
   # 两个用户进入聊天室后互相发送消息
   ```

2. **错误处理测试**
   - 测试服务器未启动时的客户端行为
   - 测试未登录直接进入聊天室
   - 测试网络断开时的行为

## 部署注意事项

1. **环境要求**
   - Go 1.16 或更高版本
   - 正确设置 GOPATH 和 GOROOT
   - 安装必要的依赖：
     ```bash
     go mod tidy
     ```

2. **网络配置**
   - 确保 50051 端口可用
   - 如需跨网络部署，修改客户端连接地址
   - 考虑添加 TLS 加密（当前版本使用不安全连接）

3. **性能考虑**
   - 监控内存使用情况
   - 注意并发连接数限制
   - 考虑添加消息队列机制

4. **安全建议**
   - 实现真实的用户认证机制
   - 添加密码加密存储
   - 实现真实的 token 生成和验证
   - 添加消息内容验证

## 扩展建议

1. **功能扩展**
   - 添加私聊功能
   - 实现消息历史记录
   - 添加文件传输功能
   - 实现用户在线状态

2. **性能优化**
   - 添加消息缓存
   - 实现消息压缩
   - 优化并发处理

3. **安全增强**
   - 添加 TLS 加密
   - 实现完整的认证机制
   - 添加消息加密

## Protocol Buffers 定义说明

### 协议文件结构

项目使用 Protocol Buffers 定义服务接口和消息结构。主要协议文件 `proto/chat.proto` 包含以下部分：

1. **包声明和选项设置**
   ```protobuf
   syntax = "proto3";
   package chat;
   option go_package = "grpc_chat/proto";
   ```

2. **服务接口定义**
   ```protobuf
   service ChatService {
     rpc Register(RegisterRequest) returns (RegisterResponse) {}
     rpc Login(LoginRequest) returns (LoginResponse) {}
     rpc Chat(stream ChatMessage) returns (stream ChatMessage) {}
   }
   ```

3. **消息类型定义**
   - 注册相关消息
     ```protobuf
     message RegisterRequest {
       string username = 1;
       string password = 2;
     }

     message RegisterResponse {
       bool success = 1;
       string message = 2;
     }
     ```
   - 登录相关消息
     ```protobuf
     message LoginRequest {
       string username = 1;
       string password = 2;
     }

     message LoginResponse {
       bool success = 1;
       string message = 2;
       string token = 3;
     }
     ```
   - 聊天消息
     ```protobuf
     message ChatMessage {
       string username = 1;
       string content = 2;
       int64 timestamp = 3;
     }
     ```

### 生成 gRPC 代码

1. **安装必要工具**
   ```bash
   # 安装 protoc 编译器
   brew install protobuf  # MacOS
   apt-get install protobuf-compiler  # Ubuntu/Debian

   # 安装 Go 插件
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

2. **生成代码命令**
   ```bash
   protoc --go_out=. \
          --go_opt=paths=source_relative \
          --go-grpc_out=. \
          --go-grpc_opt=paths=source_relative \
          proto/chat.proto
   ```

3. **生成的文件说明**
   - `grpc_chat/proto/chat.pb.go`：消息类型的 Go 代码
   - `grpc_chat/proto/chat_grpc.pb.go`：gRPC 服务接口的 Go 代码

### 接口说明

1. **Register 方法**
   - 用途：新用户注册
   - 类型：一元 RPC
   - 请求参数：用户名和密码
   - 返回结果：注册成功与否的状态

2. **Login 方法**
   - 用途：用户登录
   - 类型：一元 RPC
   - 请求参数：用户名和密码
   - 返回结果：登录状态和认证令牌

3. **Chat 方法**
   - 用途：实时聊天
   - 类型：双向流式 RPC
   - 流数据：聊天消息（包含用户名、内容和时间戳）
   - 特点：支持实时消息收发

### 类型映射

| Protocol Buffers 类型 | Go 类型 |
|---------------------|----------|
| string | string |
| bool | bool |
| int64 | int64 |
| repeated | slice |
| message | struct |

### 最佳实践

1. **字段编号**
   - 为每个字段分配唯一编号
   - 1-15 编号占用 1 个字节，常用字段优先使用
   - 16-2047 编号占用 2 个字节

2. **消息设计**
   - 相关字段组织在同一消息中
   - 避免过深的消息嵌套
   - 考虑向后兼容性

3. **服务设计**
   - 根据业务场景选择适当的 RPC 类型
   - 单向通信用一元 RPC
   - 实时交互用双向流式 RPC

## 服务器部署指南

### 1. 准备文件
需要将以下文件和目录复制到服务器：
- `server/` 目录（服务器端代码）
- `proto/` 目录（协议文件）
- `grpc_chat/` 目录（生成的 gRPC 代码）
- `go.mod` 和 `go.sum`（依赖管理文件）

### 2. 环境配置
1. 确保服务器已安装 Go 1.22 或更高版本
2. 将项目文件放置在服务器的工作目录中
3. 执行依赖安装：
   ```bash
   go mod tidy
   ```

### 3. 配置网络
1. 确保服务器的 50051 端口已开放
2. 如果有防火墙，需要添加相应的端口规则
3. 如果需要外网访问，建议配置 SSL/TLS 证书

### 4. 启动服务
1. 进入项目目录
2. 启动服务器：
   ```bash
   go run server/server.go
   ```
   或者编译后运行：
   ```bash
   go build -o chat_server server/server.go
   ./chat_server
   ```

### 5. 注意事项
- 建议使用进程管理工具（如 systemd、supervisor 等）来管理服务
- 考虑配置日志记录和监控
- 在生产环境中实现真实的用户认证和 token 机制
- 定期备份用户数据

## 服务器部署指南

### 1. 准备文件
需要将以下文件和目录复制到服务器：
- `server/` 目录（服务器端代码）
- `proto/` 目录（协议文件）
- `grpc_chat/` 目录（生成的 gRPC 代码）
- `go.mod` 和 `go.sum`（依赖管理文件）

### 2. 环境配置
1. 确保服务器已安装 Go 1.22 或更高版本
2. 将项目文件放置在服务器的工作目录中
3. 执行依赖安装：
   ```bash
   go mod tidy
   ```

### 3. 配置网络
1. 确保服务器的 50051 端口已开放
2. 如果有防火墙，需要添加相应的端口规则
3. 如果需要外网访问，建议配置 SSL/TLS 证书

### 4. 启动服务
1. 进入项目目录
2. 启动服务器：
   ```bash
   go run server/server.go
   ```
   或者编译后运行：
   ```bash
   go build -o chat_server server/server.go
   ./chat_server
   ```

### 5. 注意事项
- 建议使用进程管理工具（如 systemd、supervisor 等）来管理服务
- 考虑配置日志记录和监控
- 在生产环境中实现真实的用户认证和 token 机制
- 定期备份用户数据
