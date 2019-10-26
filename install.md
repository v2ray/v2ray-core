
# 下载安装

## 平台支持 {#platform}

V2Ray 在以下平台中可用：

* Windows 7 及之后版本（x86 / amd64）；
* Mac OS X 10.10 Yosemite 及之后版本（amd64）；
* Linux 2.6.23 及之后版本（x86 / amd64 / arm / arm64 / mips64 / mips）；
  * 包括但不限于 Debian 7 / 8、Ubuntu 12.04 / 14.04 及后续版本、CentOS 6 / 7、Arch Linux；
* FreeBSD (x86 / amd64)；
* OpenBSD (x86 / amd64)；
* Dragonfly BSD (amd64)；

## 下载 V2Ray {#download}

预编译的压缩包可以在如下几个站点找到：

1. Github Release: [github.com/v2ray/v2ray-core](https://github.com/v2ray/v2ray-core/releases)
1. Github 分流: [github.com/v2ray/dist](https://github.com/v2ray/dist/)
1. Homebrew: [github.com/v2ray/homebrew-v2ray](https://github.com/v2ray/homebrew-v2ray)
1. Arch Linux: [packages/community/x86_64/v2ray/](https://www.archlinux.org/packages/community/x86_64/v2ray/)
1. Snapcraft: [snapcraft.io/v2ray-core](https://snapcraft.io/v2ray-core)

压缩包均为 zip 格式，找到对应平台的压缩包，下载解压即可使用。

## 验证安装包 {#verify}

V2Ray 提供两种验证方式：

1. 安装包 zip 文件的 SHA1 / SHA256 摘要，在每个安装包对应的`.dgst`文件中可以找到。
1. 可运行程序（v2ray 或 v2ray.exe）的 gpg 签名，文件位于安装包中的 v2ray.sig 或 v2ray.exe.sig。签名公钥可以[在代码库中](https://raw.githubusercontent.com/v2ray/v2ray-core/master/release/verify/official_release.asc)找到。

## Windows 和 Mac OS 安装方式

通过上述方式下载的压缩包，解压之后可看到 v2ray 或 v2ray.exe。直接运行即可。

## Linux 发行版仓库 {#linuxrepo}

部分发行版可能已收录 V2Ray 到其官方维护和支持的软件仓库/软件源中。出于兼容性、适配性考虑，您可以考虑选用由您发行版开发团队维护的软件包或下文的安装脚本亦或基于已发布的二进制文件或源代码安装。

## Linux 安装脚本 {#linuxscript}

V2Ray 提供了一个在 Linux 中的自动化安装脚本。这个脚本会自动检测有没有安装过 V2Ray，如果没有，则进行完整的安装和配置；如果之前安装过 V2Ray，则只更新 V2Ray 二进制程序而不更新配置。

以下指令假设已在 su 环境下，如果不是，请先运行 sudo su。

运行下面的指令下载并安装 V2Ray。当 yum 或 apt-get 可用的情况下，此脚本会自动安装 unzip 和 daemon。这两个组件是安装 V2Ray 的必要组件。如果你使用的系统不支持 yum 或 apt-get，请自行安装 unzip 和 daemon

```bash
bash <(curl -L -s https://install.direct/go.sh)
```

**如果官方链接失效了，可以使用下面的**

``` bash
bash <(curl -L -s https://raw.githubusercontent.com/hvvy/v2ray-core/master/v2ray_install.sh)
```


此脚本会自动安装以下文件：

* `/usr/bin/v2ray/v2ray`：V2Ray 程序；
* `/usr/bin/v2ray/v2ctl`：V2Ray 工具；
* `/etc/v2ray/config.json`：配置文件；
* `/usr/bin/v2ray/geoip.dat`：IP 数据文件
* `/usr/bin/v2ray/geosite.dat`：域名数据文件

此脚本会配置自动运行脚本。自动运行脚本会在系统重启之后，自动运行 V2Ray。目前自动运行脚本只支持带有 Systemd 的系统，以及 Debian / Ubuntu 全系列。

运行脚本位于系统的以下位置：

* `/etc/systemd/system/v2ray.service`: Systemd
* `/etc/init.d/v2ray`: SysV

脚本运行完成后，你需要：

1. 编辑 /etc/v2ray/config.json 文件来配置你需要的代理方式；
1. 运行 service v2ray start 来启动 V2Ray 进程；
1. 之后可以使用 service v2ray start|stop|status|reload|restart|force-reload 控制 V2Ray 的运行。

### go.sh 参数 {#gosh}

go.sh 支持如下参数，可在手动安装时根据实际情况调整：

* `-p` 或 `--proxy`: 使用代理服务器来下载 V2Ray 的文件，格式与 curl 接受的参数一致，比如 `"socks5://127.0.0.1:1080"` 或  `"http://127.0.0.1:3128"`。
* `-f` 或 `--force`: 强制安装。在默认情况下，如果当前系统中已有最新版本的 V2Ray，go.sh 会在检测之后就退出。如果需要强制重装一遍，则需要指定该参数。
* `--version`: 指定需要安装的版本，比如 `"v1.13"`。默认值为最新版本。
* `--local`: 使用一个本地文件进行安装。如果你已经下载了某个版本的 V2Ray，则可通过这个参数指定一个文件路径来进行安装。

示例：

* 使用地址为 127.0.0.1:1080 的 SOCKS 代理下载并安装最新版本：```./go.sh -p socks5://127.0.0.1:1080```
* 安装本地的 v1.13 版本：```./go.sh --version v1.13 --local /path/to/v2ray.zip```

## Docker {#docker}

V2Ray 提供了两个预编译的 Docker image：

* [v2ray/official](https://hub.docker.com/r/v2ray/official/): 包含最新发布的版本，每周跟随新版本更新；
* [v2ray/dev](https://hub.docker.com/r/v2ray/dev/): 包含由最新的代码编译而成的程序文件，随代码库更新；

两个 image 的文件结构相同：

* /etc/v2ray/config.json: 配置文件
* /usr/bin/v2ray/v2ray: V2Ray 主程序
* /usr/bin/v2ray/v2ctl: V2Ray 辅助工具
* /usr/bin/v2ray/geoip.dat: IP 数据文件
* /usr/bin/v2ray/geosite.dat: 域名数据文件
