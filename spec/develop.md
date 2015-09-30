# 开发指引

## 基本
### 版本控制
Git

### Branch
本项目只使用一个 Branch：master。所有更改全部提交进 master，并确保 master 在任一时刻都是可编译可使用的。

### 引用其它项目
* golang
  * 产品代码只能使用 golang 的标准库，即名称不包含任何网址的包；
  * 测试代码可以使用 golang.org/x/... ；
  * 如需引用其它项目请事先创建 Issue 讨论；
* 其它
  * 只要不违反双方的协议（本项目为 MIT），且对项目有帮助的工具，都可以使用。
  

## 开发流程

### 写代码之前
发现任何问题，或对项目有任何想法，请立即[创建 Issue](https://github.com/V2Ray/v2ray-core/blob/master/spec/issue.md) 讨论之，以减少重复劳动和消耗在代码上的时间。

### 修改代码
* golang
  * 请参考 [Effective Go](https://golang.org/doc/effective_go.html)；
  * 每一次 commit 之前请运行： gofmt -w github.com/v2ray/v2ray-core/
  * 每一次 commit 之前请确保测试通过： go test github.com/v2ray/v2ray-core/...
  * 提交 PR 之前请确保新增代码有超过 60% 的代码覆盖率（code coverage）。
* 其它
  * 请注意代码的可读性
  
### Pull Request
提交 PR 之前请先运行 git pull 以确保 merge 可顺利进行。

## 对代码的修改
### 功能性问题
请提交至少一个测试用例（test case）来验证对现有功能的改动。

### 性能相关
请提交必要的测试数据来证明现有代码的性能缺陷，或是新增代码的性能提升。

### 其它
视具体情况而定。


