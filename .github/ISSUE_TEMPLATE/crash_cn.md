---
name: V2Ray 进程崩溃
about: "提交一个 V2Ray 的 panic 日志"
---

提交 Issue 之前请先阅读 [Issue 指引](https://github.com/v2ray/v2ray-core/blob/master/.github/SUPPORT.md)，然后回答下面的问题，谢谢。
除非特殊情况，请完整填写所有问题。不按模板发的 issue 将直接被关闭。
如果你遇到的问题不是 V2Ray 的 bug，比如你不清楚要如何配置，请使用[Discussion](https://github.com/v2ray/discussion/issues)进行讨论。

1) 你正在使用哪个版本的 V2Ray？（如果服务器和客户端使用了不同版本，请注明）

2) 你的使用场景是什么？比如使用 Chrome 通过 Socks/VMess 代理观看 YouTube 视频。

3) 请附上 panic 时的完整输出。

在 Linux (Systemd) 上可以通过 `journalctl -u v2ray` 查看最近的 panic 日志。

4) 请附上配置文件（提交 Issue 前请隐藏服务器端IP地址）。

```javascript
    // 在这里附上配置文件
```

请预览一下你填的内容再提交。
