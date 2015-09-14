# 当前状态

**注释**
* A：稳定可用
* B：默认开启，还在测试期，可能有各种问题
* C：需要手动开启
* D：正在调试期，暂不可用
* E：计划中，欢迎 Pull Request

**概况**

| 功能       | 状态 | 备注 |
| --------- | ---- | ---- |
| 多用户支持  | B  |  |
| 多服务器支持  | B  |  |
| 负载均衡 | E | |
| 多种加密方式 | E | 暂时只支持 AES-128 |
| 选择性路由 | E | |
| 自定义 DNS 解析 | E | |

**平台支持**

| 平台       | 状态 | 备注 |
| --------- | ---- | ---- |
| golang 编译  | B  |  |
| Windows  | E  |  |
| Mac OS | E | |
| Ubuntu | E | |
| Redhat | E | |
| OpenWRT | E | |

**Socks 5 协议**

| 功能       | 状态 | 备注 |
| --------- | ---- | ---- |
| TCP 连接    | B |  |
| UDP 连接    | E | [Issue #3](https://github.com/v2ray/v2ray-core/issues/3) |
| FTP 支持    | E | [Issue #2](https://github.com/v2ray/v2ray-core/issues/2) |

**[VMess 协议](https://github.com/V2Ray/v2ray-core/blob/master/spec/vmess.md)**

| 功能       | 状态 | 备注 |
| --------- | ---- | ---- |
| 单一连接    | B |  |
| 连接复用    | E |  |

**ShadowSocks 协议**

| 功能       | 状态 | 备注 |
| --------- | ---- | ---- |
| 单一连接    | E |  |

