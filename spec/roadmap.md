# V2Ray 开发计划

## 版本号
V2Ray 的版本号形如 X.Y.Z，其中 X 表示 Milestone，Y 表示 Release，如 2.3 表示第二个 Milestone 的第三个 Release；Z 仅作为修复紧急 Bug 之后的发布使用，一般不出现。

## 周期
V2Ray 将在每周一发布一个 [Release](https://github.com/v2ray/v2ray-core/releases)，每 12 周左右完成一个 Milestone。

## Milestones

### Milestone 0
**目标：可用**

M0 将提供一个可用的 V2Ray Point Server，包含 Windows、Mac OS 和 Linux（Debian 为主）的预编译文件。主要功能和限制如下：
* SOCKS 4 / 5 协议，仅提供 TCP 代理；
* 使用自有 [VMess 协议](https://github.com/V2Ray/v2ray-core/blob/master/spec/vmess.md)做隧道；
* 服务器端支持多用户；
* 客户端支持多服务器，暂不支持负载平衡，多服务器时随机选择服务器；

### Milestone 1
**目标：兼容**

M1 将完成在服务器端对兼容 Shadowsocks 和 GoAgent 协议的兼容，为用户提供多种选择。期望的功能如下：
* SOCKS 协议的 UDP 代理；
* 可选择路由，不必要的网站不使用代理；
* WebSocket 代理；

### Milestone 2
**目标：交互**

M2 将提供必要的 API 供第三方程序调用。主要功能如下：
* 可以查询 V2Ray 进程的当前状态；
* 可以远程和 V2Ray 进程通信并控制；
* 可以动态管理用户和服务器列表；
