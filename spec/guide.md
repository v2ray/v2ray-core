# V2Ray 简明教程

## 工作机制

你需要至少两个 Point Server（设为 A、B）才可以正常穿墙。以网页浏览为例，你的浏览器和 A 以 Socks 5 协议通信，B 和目标网站之间以 HTTP 协议通信，A 和 B 之间使用 V2Ray 的自有协议 [VMess](https://github.com/V2Ray/v2ray-core/blob/master/spec/vmess.md) 通信，如下图：

![](https://github.com/V2Ray/v2ray-core/blob/master/spec/v2ray.png)

通常 Point A 运行在你自己的电脑，Point B 运行在一台海外的 VPS 中。

## 安装 V2Ray Point Server
[安装 V2Ray](https://github.com/V2Ray/v2ray-core/blob/master/spec/install.md)

## 配置 V2Ray Point Server
### Point A
示例配置保存于 [vpoint_socks_vmess.json](https://github.com/v2ray/v2ray-core/blob/master/release/config/vpoint_socks_vmess.json) 文件中，格式如下：
```javascript
{
  "port": 1080, // 监听端口
  "log" : {
    "access": "" // 访问记录，目前只在服务器端有效，这里留空
  },
  "inbound": {
    "protocol": "socks",  // 传入数据所用协议
    "settings": {
      "auth": "noauth", // 认证方式，暂时只支持匿名
      "udp": false // 如果要使用 UDP 转发，请改成 true
    }
  },
  "outbound": {
    "protocol": "vmess", // 中继协议，暂时只有这个
    "settings": {
      "vnext": [
        {
          "address": "127.0.0.1", // Point B 的 IP 地址，IPv4 或 IPv6，不支持域名
          "port": 27183, // Point B 的监听端口，请更换成其它的值
          "users": [
            // 用户 ID，必须包含在 Point B 的配置文件中。此 ID 将被用于通信的认证，请自行更换随机的 ID，可以使用 https://www.uuidgenerator.net/ 来生成新的 ID。
            {"id": "ad937d9d-6e23-4a5a-ba23-bce5092a7c51"}
          ],
          "network": "tcp" // 如果要使用 UDP 转发，请改成 "tcp,udp"
        }
      ]
    }
  }
}
```

### Point B
示例配置保存于 [vpoint_vmess_freedom.json](https://github.com/v2ray/v2ray-core/blob/master/release/config/vpoint_vmess_freedom.json) 文件中，格式如下：
```javascript
{
  "port": 27183, // 监听端口，必须和 Point A 中指定的一致
  "log" : {
    "access": "access.log" // 访问记录
  },
  "inbound": {
    "protocol": "vmess", // 中继协议，不用改
    "settings": {
      "clients": [
          // 认可的用户 ID，必须包含 Point A 中的用户 ID
        {"id": "ad937d9d-6e23-4a5a-ba23-bce5092a7c51"}
      ],
      "udp": false // 如果要使用 UDP 转发，请改成 true
    }
  },
  "outbound": {
    "protocol": "freedom", // 出口协议，不用改
    "settings": {} // 暂无配置
  }
}
```

### 其它
* V2Ray 的用户验证基于时间，请确保 A 和 B 所在机器的系统时间误差在一分钟以内。
* json 配置文件实际上不支持注释（即“//”之后的部分，在使用时请务必删去）。

## 运行

Point Server A

./server --config="vpoint_socks_vmess.json 的绝对路径"

Point Server B

./server --config="vpoint_vmess_freedom.json 的绝对路径"

## 测试服务器可用性

curl -v --socks5-hostname 127.0.0.1:1080 https://www.google.com/

## 调试

使用过程中遇到任何问题，请参考[错误信息](https://github.com/V2Ray/v2ray-core/blob/master/spec/errors.md)。
