# VMess 设计
## 摘要
* 版本：1

## 格式
### 数据请求
认证部分：
* 16 字节：基于时间的 hash(用户 [ID](https://github.com/V2Ray/v2ray-core/blob/master/spec/id.md))，见下文

指令部分：
* 1 字节：版本号，目前为 0x1
* 16 字节：请求数据 IV
* 16 字节：请求数据 Key
* 4 字节：认证信息 V
* 1 字节：指令
  * 0x00：保留
  * 0x01：TCP 请求
  * 0x02：UDP 请求
* 2 字节：目标端口
* 1 字节：目标类型
  * 0x01：IPv4
  * 0x02：域名
  * 0x03：IPv6
* 目标地址：
  * 4 字节：IPv4
  * 1 字节长度 + 域名
  * 16 字节：IPv6
* 4 字节：指令部分前面所有内容的 FNV1a hash

数据部分
* N 字节：请求数据

其中指令部分经过 AES-128-CFB 加密：
* Key：md5(用户 ID + 'c48619fe-8f02-49e0-b9e9-edf763e17e21')
* IV：md5(X + X + X + X)，X = []byte(UserHash 生成的时间) (8 字节, Big Endian)

数据部分使用 AES-128-CFB 加密，Key 和 IV 在请求数据中

### 数据应答
数据部分
* 4 字节：认证信息 V
* N 字节：应答数据

其中数据部分使用 AES-128-CFB 加密，IV 为 md5(请求数据 IV)，Key 为 md5(请求数据 Key)

## 基于时间的用户 ID Hash

* H = MD5
* K = 用户 ID (16 字节)
* M = UTC 时间，精确到秒，取值为当前时间的前后 30 秒随机值(8 字节, Big Endian)
* Hash = HMAC(H, K, M)
