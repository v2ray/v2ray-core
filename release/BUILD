package(default_visibility=["//visibility:public"])

load("@v2ray_ext//bazel:build.bzl", "foreign_go_binary")
load("@v2ray_ext//bazel:zip.bzl", "pkg_zip")

package(default_visibility=["//visibility:public"])

foreign_go_binary(
  name = "v2ray_linux_amd64",
  pkg = "v2ray.com/core/main",
  output = "v2ray",
  os = "linux",
  arch = "amd64",
)

filegroup(
  name = "config_json",
  srcs = ["config/config.json"],
)

filegroup(
  name = "config_all",
  srcs = ["config/vpoint_socks_vmess.json", "config/vpoint_vmess_freedom.json"],
)

filegroup(
  name = "systemd",
  srcs = ["config/systemd/v2ray.service", "config/systemv/v2ray"],
)

filegroup(
  name = "doc",
  srcs = glob(["doc/*.md"]),
)

pkg_zip(
  name = "v2ray_linux_amd64_package",
  srcs = [
    ":v2ray_linux_amd64",
    "@v2ray_ext//tools/control/main:v2ctl_linux_amd64",
    ":config_all",
    ":systemd",
    ":doc",
  ],
  out = "v2ray-linux-64.zip",
)
