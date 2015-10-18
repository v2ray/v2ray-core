# VMess 设计
## 摘要
* 版本：1

## 格式
### 数据请求
认证部分：
* 16 字节：基于时间的 hash(用户 ID)，见下文

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

## 用户 ID
ID 等价于 [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier)，是一个 16 字节长的随机数，它的作用相当于一个令牌（Token）。

一个 ID 形如：de305d54-75b4-431b-adb2-eb6b9e546014，几乎完全随机，可以使用任何的 UUID 生成器来生成，比如[这个](https://www.uuidgenerator.net/)。

ID 在消息传递过程中用于验证客户端的有效性，只有当服务器认可当前 ID 时，才进行后续操作，否则关闭连接甚至加入黑名单。

在多用户环境中，用户帐号应与 ID 分开存放，即用户帐号和 ID 有一对一或一对多的关系，在 V2Ray Server 中，只负责管理 ID，用户帐号（及权限、费用等）由另外的系统管理。

在后续版本中，V2Ray Server 之间应有能力进行沟通而生成新的临时 ID，从而减少通讯的可探测性。

## 基于时间的用户 ID Hash

* H = MD5
* K = 用户 ID (16 字节)
* M = UTC 时间，精确到秒，取值为当前时间的前后 30 秒随机值(8 字节, Big Endian)
* Hash = HMAC(H, K, M)
