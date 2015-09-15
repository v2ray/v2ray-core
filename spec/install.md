# V2Ray 安装方式

目前 V2Ray 还在早期测试阶段，暂不提供预编译的运行文件。请使用下面的方式下载源文件并编译。

## 编译源文件
1. 安装 Git： sudo apt-get install git -y
2. 安装 golang：
  1. curl -o go_latest.tar.gz https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
  2. sudo tar -C /usr/local -xzf go_latest.tar.gz
3. export PATH=$PATH:/usr/local/go/bin
4. export GOPATH=$HOME/work
5. go get github.com/v2ray/v2ray-core
6. go build github.com/v2ray/v2ray-core/release/server


