package all

import (
	// The following are necessary as they register handlers in their init functions.
	_ "v2ray.com/core/app/dispatcher/impl"
	_ "v2ray.com/core/app/dns/server"
	_ "v2ray.com/core/app/proxy"
	_ "v2ray.com/core/app/proxyman/outbound"
	_ "v2ray.com/core/app/router"

	_ "v2ray.com/core/proxy/blackhole"
	_ "v2ray.com/core/proxy/dokodemo"
	_ "v2ray.com/core/proxy/freedom"
	_ "v2ray.com/core/proxy/http"
	_ "v2ray.com/core/proxy/shadowsocks"
	_ "v2ray.com/core/proxy/socks"
	_ "v2ray.com/core/proxy/vmess/inbound"
	_ "v2ray.com/core/proxy/vmess/outbound"

	_ "v2ray.com/core/transport/internet/kcp"
	_ "v2ray.com/core/transport/internet/tcp"
	_ "v2ray.com/core/transport/internet/tls"
	_ "v2ray.com/core/transport/internet/udp"
	_ "v2ray.com/core/transport/internet/websocket"

	_ "v2ray.com/core/transport/internet/headers/http"
	_ "v2ray.com/core/transport/internet/headers/noop"
	_ "v2ray.com/core/transport/internet/headers/srtp"
	_ "v2ray.com/core/transport/internet/headers/utp"
	_ "v2ray.com/core/transport/internet/headers/wechat"
)
