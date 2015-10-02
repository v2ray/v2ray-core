# V2Ray 安装方式

## 预编译程序
发布于 [Release](https://github.com/v2ray/v2ray-core/releases) 中，每周更新，[更新周期见此](https://github.com/V2Ray/v2ray-core/blob/master/spec/roadmap.md)。

## 编译源文件

大概流程，请根据实际情况修改

1. 安装 Git： sudo apt-get install git -y
2. 安装 golang：
  1. 下载安装文件：
    1. 64位：curl -o go_latest.tar.gz https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
    2. 32位：curl -o go_latest.tar.gz https://storage.googleapis.com/golang/go1.5.1.linux-386.tar.gz
  2. sudo tar -C /usr/local -xzf go_latest.tar.gz
  3. export PATH=$PATH:/usr/local/go/bin
  4. export GOPATH=$HOME/work
3. go get github.com/v2ray/v2ray-core
4. go build github.com/v2ray/v2ray-core/release/server

### Arch Linux
1. 安装 Git： sudo pacman -S git
2. 安装 golang：sudo pacman -S go
   1. export GOPATH=$HOME/work
3. go get github.com/v2ray/v2ray-core
4. go build -o $GOPATH/bin/v2ray -compiler gc github.com/v2ray/v2ray-core/release/server

### Debian / Ubuntu
bash <(curl -s https://raw.githubusercontent.com/v2ray/v2ray-core/master/release/install.sh)

此脚本会自动安装 git 和 golan 1.5 （如果系统上没有的话，并且需要 root 权限），然后把 v2ray 编译到 $GOPATH/bin/v2ray，新装的 golang 会把 GOPATH 设定到 /v2ray。


## 配置和运行
[链接](https://github.com/V2Ray/v2ray-core/blob/master/spec/guide.md)
