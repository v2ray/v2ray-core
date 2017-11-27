# V2Ray 用户支持 (User Support)

**English reader please skip to the [English section](#way-to-get-support) below**

## 获得帮助信息的途径

您可以从以下渠道获取帮助：

1. 官方网站：[v2ray.com](https://www.v2ray.com)
1. Github：[Issues](https://github.com/v2ray/v2ray-core/issues)
1. Telegram：[主群](https://t.me/projectv2ray)

## Github Issue 规则

1. 请按模板填写 issue；
1. 配置文件内容使用格式化代码段进行修饰（见下面的解释）；
1. 在提交 issue 前尝试减化配置文件，比如删除不必要 inbound / outbound 模块；
1. 在提交 issue 前尝试确定问题所在，比如将 socks 代理换成 http 再次观察问题是否能重现；
1. 配置文件必须结构完整，即除了必要的隐私信息之外，配置文件可以直接拿来运行。

**不按模板填写的 issue 将直接被关闭**

## 格式化代码段

在配置文件上下加入 Markdown 特定的修饰符，如下：

\`\`\`javascript

{
  // 配置文件内容
}

\`\`\`

## Way to Get Support

You may get help in the following ways:

1. Office Site: [v2ray.com](https://www.v2ray.com)
1. Github: [Issues](https://github.com/v2ray/v2ray-core/issues)
1. Telegram: [Main Group](https://t.me/projectv2ray)

## Github Issue Rules

1. Please fill in the issue template.
1. Decorate config file with Markdown formatter (See below).
1. Try to simplify config file before submitting the issue, such as removing unnecessary inbound / outbound blocks.
1. Try to determine the cause of the issue, for example, replacing socks inbound with http inbound to see if the issue still exists.
1. Config file must be structurally complete.

**Any issue not following the issue template will be closed immediately.**

## Code formatter

Add the following Markdown decorator to config file content:

\`\`\`javascript

{
  // config file
}

\`\`\`
