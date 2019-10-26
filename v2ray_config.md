## 配置文件示范

### WebSocket+TLS+Web

服务器配置
这次 TLS 的配置将写入 Nginx / Caddy / Apache 配置中，由这些软件来监听 443 端口（443 比较常用，并非 443 不可），然后将流量转发到 V2Ray 的 WebSocket 所监听的内网端口（本例是 10000），V2Ray 服务器端不需要配置 TLS。

服务器 V2Ray 配置

``` json
{
  "inbounds": [
    {
      "port": 10000,
      "listen":"127.0.0.1",//只监听 127.0.0.1，避免除本机外的机器探测到开放了 10000 端口
      "protocol": "vmess",
      "settings": {
        "clients": [
          {
            "id": "b831381d-6324-4d53-ad4f-8cda48b30811",
            "alterId": 64
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
        "path": "/ray"
        }
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom",
      "settings": {}
    }
  ]
}
```

客户端配置

``` json
{
  "inbounds": [
    {
      "port": 1080,
      "listen": "127.0.0.1",
      "protocol": "socks",
      "sniffing": {
        "enabled": true,
        "destOverride": ["http", "tls"]
      },
      "settings": {
        "auth": "noauth",
        "udp": false
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "mydomain.me",
            "port": 443,
            "users": [
              {
                "id": "b831381d-6324-4d53-ad4f-8cda48b30811",
                "alterId": 64
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "security": "tls",
        "wsSettings": {
          "path": "/ray"
        }
      }
    }
  ]
}
```

### Nginx 配置

配置中使用的是域名和证书使用 TLS 小节的举例，请替换成自己的。

``` bash
server {
  listen 443 ssl;
  ssl on;
  ssl_certificate       /etc/v2ray/v2ray.crt;
  ssl_certificate_key   /etc/v2ray/v2ray.key;
  ssl_protocols         TLSv1 TLSv1.1 TLSv1.2;
  ssl_ciphers           HIGH:!aNULL:!MD5;
  server_name           mydomain.me;
    location /ray { # 与 V2Ray 配置中的 path 保持一致
      proxy_redirect off;
      proxy_pass http://127.0.0.1:10000; # 假设WebSocket监听在环回地址的10000端口上
      proxy_http_version 1.1;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection "upgrade";
      proxy_set_header Host $host;
      # Show real IP in v2ray access.log
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### Caddy 配置

因为 Caddy 会自动申请证书并自动更新，所以使用 Caddy 不用指定证书、密钥。

``` bash
mydomain.me
{
  log ./caddy.log
  proxy /ray localhost:10000 {
    websocket
    header_upstream -Origin
  }
}
```

### Apache 配置

同样地，配置中使用的是域名和证书使用 TLS 小节的举例，请替换成自己的。

``` bash
<VirtualHost *:443>
  ServerName mydomain.me
  SSLCertificateFile /etc/v2ray/v2ray.crt
  SSLCertificateKeyFile /etc/v2ray/v2ray.key
  
  SSLProtocol -All +TLSv1 +TLSv1.1 +TLSv1.2
  SSLCipherSuite HIGH:!aNULL
  
  <Location "/ray/">
    ProxyPass ws://127.0.0.1:10000/ray/ upgrade=WebSocket
    ProxyAddHeaders Off
    ProxyPreserveHost On
    RequestHeader append X-Forwarded-For %{REMOTE_ADDR}s
  </Location>
</VirtualHost>
```


### mKCP

V2Ray 引入了 KCP 传输协议，并且做了一些不同的优化，称为 mKCP。如果你发现你的网络环境丢包严重，可以考虑一下使用 mKCP。由于快速重传的机制，相对于常规的 TCP 来说，mKCP 在高丢包率的网络下具有更大的优势，也正是因为此， mKCP 明显会比 TCP 耗费更多的流量，所以请酌情使用。要了解的一点是，mKCP 与 KCPTUN 同样是 KCP 协议，但两者并不兼容。

在此我想纠正一个概念。基本上只要提起 KCP 或者 UDP，大家总会说”容易被 Qos“。Qos 是一个名词性的短语，中文意为服务质量，试想一下，你跟人家说一句”我的网络又被服务质量了“是什么感觉。其次，哪怕名词可以动词化，这么使用也是不合适的，因为 Qos 区分网络流量优先级的，就像马路上划分人行道、非机动车道、快车道、慢车道一样，哪怕你牛逼到运营商送你一条甚至十条专线，是快车道中的快车道，这也是 Qos 的结果。

mKCP 的配置比较简单，只需在服务器的 inbounds 和 客户端的 outbounds 添加一个 streamSettings 并设置成 mkcp 即可

服务器配置

``` json
{
  "inbounds": [
    {
      "port": 16823,
      "protocol": "vmess",
      "settings": {
        "clients": [
          {
            "id": "b831381d-6324-4d53-ad4f-8cda48b30811",
            "alterId": 64
          }
        ]
      },
      "streamSettings": {
        "network": "mkcp", //此处的 mkcp 也可写成 kcp，两种写法是起同样的效果
        "kcpSettings": {
          "uplinkCapacity": 5,
          "downlinkCapacity": 100,
          "congestion": true,
          "header": {
            "type": "none"
          }
        }
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom",
      "settings": {}
    }
  ]
}
```

客户端配置

``` json
{
  "inbounds": [
    {
      "port": 1080,
      "protocol": "socks",
      "sniffing": {
        "enabled": true,
        "destOverride": ["http", "tls"]
    },
      "settings": {
        "auth": "noauth"
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "serveraddr.com",
            "port": 16823,
            "users": [
              {
                "id": "b831381d-6324-4d53-ad4f-8cda48b30811",
                "alterId": 64
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "mkcp",
        "kcpSettings": {
          "uplinkCapacity": 5,
          "downlinkCapacity": 100,
          "congestion": true,
          "header": {
            "type": "none"
          }
        }
      }
    }
  ]
}
```

### HTTP/2

简单地说 HTTP/2 是 HTTP/1.1 的升级版（目前大多数网页还是 HTTP/1.1），点击这里可以直观地体会到 HTTP/2 相比于 HTTP/1.1 的提升（不代表 V2Ray 中 HTTP/2 相对于 TCP 的提升就是这样的）。HTTP/2协议一般简称为h2。

在v2ray中使用h2， 经常被用户们用来跟websocket方式做比较。从理论上来说，HTTP/2在首次连接时候，不像websocket需完成upgrade请求；v2ray客户端和服务端之间一般直接通信，较少中间层代理。但是，在配合 CDN、Nginx/Caddy/Apache等服务组件作为前置分流代理的应用场景上，h2没有websocket方式灵活，因为很多代理并不提供h2协议的后端支持。实际使用中，websocket和h2的方式，在体验上很可能没有明显区别，用户可自行根据需要选择。


与其它的传输层协议一样在 streamSettings 中配置，不过要注意的是使用 HTTP/2 要开启 TLS。

服务器配置

``` json
{
  "inbounds": [
    {
      "port": 443,
      "protocol": "vmess",
      "settings": {
        "clients": [
          {
            "id": "b831381d-6324-4d53-ad4f-8cda48b30811",
            "alterId": 64
          }
        ]
      },
      "streamSettings": {
        "network": "h2", // h2 也可写成 http，效果一样
        "httpSettings": { //此项是关于 HTTP/2 的设置
          "path": "/ray"
        },
        "security": "tls", // 配置tls
        "tlsSettings": {
          "certificates": [
            {
              "certificateFile": "/etc/v2ray/v2ray.crt", // 证书文件，详见 tls 小节
              "keyFile": "/etc/v2ray/v2ray.key" // 密钥文件
            }
          ]
        }
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom",
      "settings": {}
    }
  ]
}
```

客户端配置

``` json
{
  "inbounds": [
    {
      "port": 1080,
      "listen": "127.0.0.1",
      "protocol": "socks",
      "sniffing": {
        "enabled": true,
        "destOverride": ["http", "tls"]
      },
      "settings": {
        "auth": "noauth",
        "udp": false
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "mydomain.me",
            "port": 443,
            "users": [
              {
                "id": "b831381d-6324-4d53-ad4f-8cda48b30811",
                "alterId": 64
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "h2",
        "httpSettings": { //此项是关于 HTTP/2 的设置
          "path": "/ray"
        },
        "security": "tls"
      }
    }
  ]
}
```

### Shadowsocks

本节讲述 Shadowsocks 的配置。

其实，作为一个代理工具集合，V2Ray 集成有 Shadowsocks 模块。用 V2Ray 配置成 Shadowsocks 服务器或者 Shadowsocks 客户端都是可以的，兼容 Shadowsocks-libev, go-shadowsocks2 等基于 Shadowsocks 协议的客户端。

配置与 VMess 大同小异，客户端服务器端都要有入口和出口，只不过是协议(protocol)和相关设置(settings)不同，不作过多说明，直接给配置，如果你配置过 Shadowsocks，对比之下就能够明白每个参数的意思(配置还有注释说明呢)。


客户端配置

``` json
{
  "inbounds": [
    {
      "port": 1080, // 监听端口
      "protocol": "socks", // 入口协议为 SOCKS 5
      "sniffing": {
        "enabled": true,
        "destOverride": ["http", "tls"]
      },
      "settings": {
        "auth": "noauth"  // 不认证
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "shadowsocks",
      "settings": {
        "servers": [
          {
            "address": "serveraddr.com", // Shadowsocks 的服务器地址
            "method": "aes-128-gcm", // Shadowsocks 的加密方式
            "ota": true, // 是否开启 OTA，true 为开启
            "password": "sspasswd", // Shadowsocks 的密码
            "port": 1024  
          }
        ]
      }
    }
  ]
}
```

服务器配置

``` json
{
  "inbounds": [
    {
      "port": 1024, // 监听端口
      "protocol": "shadowsocks",
      "settings": {
        "method": "aes-128-gcm",
        "ota": true, // 是否开启 OTA
        "password": "sspasswd"
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom",  
      "settings": {}
    }
  ]
}
```

注意事项

1. 因为协议漏洞，Shadowsocks 已放弃 OTA(一次认证) 转而使用 AEAD，V2Ray 的 Shadowsocks 协议已经跟进 AEAD，但是仍然兼容 OTA。建议使用 AEAD (method 为 aes-256-gcm、aes-128-gcm、chacha20-poly1305 即可开启 AEAD), 使用 AEAD 时 OTA 会失效；
2. Shadowsocks 已经弃用 simple-obfs，可使用基于 V2Ray 的新版混淆插件（但也可以使用 V2Ray 的 Websocket/http2 + TLS ）；
3. 可以使用 V2Ray 的传输层配置（详见高级篇），但如果这么设置了将与原版 Shadowsocks 不兼容（兼容 Shadowsocks 新增的 v2ray-plugin插件)。

更新历史

- 2018-02-09 AEAD 更新
- 2018-09-03 描述更新
- 2018-11-09 跟进 v4.0+ 的配置格式
- 2019-01-19 v2ray-plugin