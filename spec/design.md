# V2Ray 设计概要

## 目标
* V2Ray 自身提供基础的翻墙功能，单独使用可满足日常需求；
* V2Ray 可被用于二次开发，可为用户提供个性化的翻墙体验，从而也达到难以统一屏蔽之效果；
* V2Ray 为模块化设计，模块之间互相独立。每个模块可单独使用，也可和其它模块搭配使用。

## 架构

### 术语
* Point：一个 V2Ray 服务器称为 Point Server
* Set：一组 Point Server，包含多个 Point 进程，由一个 Master 进程统一管理。
* SuperSet：多机环境中的多个 Set

### 工作流程
Point 可接收来自用户或其它 Point 的请求，并将请求转发至配置中的下一个 Point（或 Set 或 SuperSet） 或目标网站，然后将所得到的应答回复给请求来源。
Point 采用白名单机制，只接受已认证帐号的请求。

### 通信协议
* Point 之间默认使用自有 VMess 协议，或第三方自定义协议。
* Point 和客户端之间可使用以下协议：
  * HTTP Proxy
  * SOCKS Proxy
  * PPTP / L2TP / SSTP 等 VPN 隧道
  * 其它自定义协议
* Point 和目标网站之间使用以下协议：
  * HTTP / HTTPS
  * UDP (DNS)

#### VMess
VMess 为 V2Ray 的原生协议，设计用于两个 Point 之间的通信。[详细设计](https://github.com/V2Ray/v2ray-core/blob/master/spec/vmess.md)

### Point
* 每个 Point 有一个 ID，运行时生成
* 每个 Point 可使用独立的配置文件，或从 Set 继承
* 一个 Point 监听主机上的一个特定端口（可配置），用于接收和发送数据
* 一个 Point 运行于一个独立的进程，可设定其使用的系统帐户

### Set
TODO

### SuperSet
TODO

## Point 详细设计
一个 Point 包含五个部分：
* 配置文件处理：读取和解析配置文件
* 输入（Inbound）：负责与客户端建立连接（如 TCP），接收客户端的消息
* 输出（Outbound）：负责向客户端发送消息

### 配置文件
配置文件使用 JSON / ProtoBuf 兼容格式

## 编程语言
暂定为 golang。
