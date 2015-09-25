# 错误信息

## 简介
在日志中可以看到 [Error XXXXXX] 的信息，其中 XXXXXX 表示错误代码，已知的错误代码和解释如下：


## 0x0001 Authentication Error
* 原因：未认证用户。
* 解决：请检查客户端和服务器的用户数据。

## 0x0002 Protocol Version Error
* 原因：客户端使用了不正确的协议
* 解决：
  * 如果错误信息为 Invalid version 67 （或 71、80），则表示你的浏览器使用了 HTTP 代理，而 V2Ray 只接受 Socks 代理。
  * 请检查客户端配置。

## 0x0003 Corrupted Packet Error
* 原因：网络数据损坏
* 解决：极有可能你的网络连接被劫持，请更换网络线路或 IP。


## 0x0004 IP Format Error
* 原因：不正确的 IP 地址
* 解决：请检查客户端软件，如浏览器的配置

## 0x0005 Configuration Error
* 原因：配置文件不能正常读取
* 解决：请检查配置文件是否存在，权限是否合适，内容是否正常

## 0x0006 Invalid Operation Error
* 原因：不正确的操作


## 0x03E8 Socks Version 4
* 原因：客户端使用了 SOCKS 4 协议
* 解决：升级客户端软件

