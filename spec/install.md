# V2Ray 安装方式

目前 V2Ray 还在早期测试阶段，暂不提供预编译的运行文件。请使用下面的方式下载源文件并编译。

## 编译源文件
1. 安装 Git： sudo apt-get install git -y
2. 安装 golang：
  1. curl -o go_latest.tar.gz https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
  2. sudo tar -C /usr/local -xzf go_latest.tar.gz
  3. export PATH=$PATH:/usr/local/go/bin
  4. export GOPATH=$HOME/work
3. go get github.com/v2ray/v2ray-core
4. go build github.com/v2ray/v2ray-core/release/server

### Debian / Ubuntu
sudo bash <(curl -s https://raw.githubusercontent.com/v2ray/v2ray-core/master/release/install.sh)

此脚本会自动安装 git 和 golan 1.5 （如果系统上没有的话），然后把 v2ray 编译到 $GOPATH/bin/v2ray，新装的 golang 会把 GOPATH 设定到 /v2ray。

## 配置和运行
[链接](https://github.com/V2Ray/v2ray-core/blob/master/spec/guide.md)
