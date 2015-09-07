# VID 的定义和使用

VID 等价于 [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier)，是一个 16 字节长的随机数，它的作用相当于一个令牌（Token）。

## 设计
一个 VID 形如：de305d54-75b4-431b-adb2-eb6b9e546014，几乎完全随机，可以使用任何的 UUID 生成器来生成，比如[这个](https://www.uuidgenerator.net/)。

## 使用
VID 在消息传递过程中用于验证客户端的有效性，只有当服务器认可当前 VID 时，才进行后续操作，否则关闭连接甚至加入黑名单。

在多用户环境中，用户帐号应与 VID 分开存放，即用户帐号和 VID 有一对一或一对多的关系，在 VPoint 系统中，只负责管理 VID，用户帐号（及权限、费用等）由另外的系统管理。

在后续版本中，VPoint 之间应有能力进行沟通而生成新的临时 VID，从而减少通讯的可探测性。

