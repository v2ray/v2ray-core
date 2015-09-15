# V2Ray 简明教程

## 工作机制

你需要至少两个 Point Server（设为 A、B）才可以正常穿墙。以网页浏览为例，你的浏览器和 A 以 Socks 5 协议通信，B 和目标网站之间以 HTTP 协议通信，A 和 B 之间使用 V2Ray 的自有协议 [VMess](https://github.com/V2Ray/v2ray-core/blob/master/spec/vmess.md) 通信，如下图：

![](https://github.com/V2Ray/v2ray-core/blob/master/spec/v2ray.png)

通常 Point A 运行在你自己的电脑，Point B 运行在一台海外的 VPS 中。

## 安装 V2Ray Point Server
[安装 V2Ray](https://github.com/V2Ray/v2ray-core/blob/master/spec/install.md)

## 配置 V2Ray Point Server
### Point A
示例配置保存于 vpoint_socks_vmess.json 文件中，格式如下：
```javascript
{
  "port": 1080, // 监听端口
  "inbound": {
    "protocol": "socks",  // 传入数据所用协议
    "file": "in_socks.json" // socks 配置文件
  },
  "outbound": {
    "protocol": "vmess", // 中继协议
    "file": "out_vmess.json" // vmess 配置文件
  }
}
```

另外还需要两个文件，保存于同一文件夹下：

```javascript
// in_socks.json
{
  "auth": "noauth" // 认证方式，暂时只支持匿名
}
```

```javascript
// out_vmess.json
{
  "vnext": [
    {
      "address": "127.0.0.1", // Point B 的 IP 地址
      "port": 27183, // Point B 的监听端口
      "users": [
        {"id": "ad937d9d-6e23-4a5a-ba23-bce5092a7c51"}  // 用户 ID，必须包含在 Point B 的配置文件中
      ]
    }
  ]
}
```

### Point B
示例配置保存于 vpoint_vmess_freedom.json 文件中，格式如下：
```javascript
{
  "port": 27183, // 监听端口
  "inbound": {
    "protocol": "vmess", // 中继协议
    "file": "in_vmess.json" // vmess 配置文件
  },
  "outbound": {
    "protocol": "freedom", // 出口协议，暂时只有这一个，不用改
    "file": "" // 暂无配置
  }
}
```

另外还需要 in_vmess.json：
```javascript
// in_vmess.json
{
  "clients": [
    {"id": "ad937d9d-6e23-4a5a-ba23-bce5092a7c51"}  // 认可的用户 ID
  ]
}
```

## 运行

Point Server A

./server --config="vpoint_socks_vmess.json 的绝对路径"

Point Server B

./server --config="vpoint_vmess_freedom.json 的绝对路径"

## 测试服务器可用性：

curl -v --socks5-hostname 127.0.0.1:1080 https://www.google.com/

