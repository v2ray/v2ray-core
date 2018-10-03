package(default_visibility=["//visibility:public"])

load("@v2ray_ext//bazel:build.bzl", "foreign_go_binary")
load("@v2ray_ext//bazel:zip.bzl", "pkg_zip")
load("@v2ray_ext//bazel:gpg.bzl", "gpg_sign")

foreign_go_binary(
  name = "v2ray_linux_amd64",
  pkg = "v2ray.com/core/main",
  output = "v2ray",
  os = "linux",
  arch = "amd64",
)

gpg_sign(
  name = "v2ray_linux_amd64_sig",
  base = ":v2ray_linux_amd64",
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
  srcs = ["config/systemd_v2ray.service", "config/systemv_v2ray"],
)

filegroup(
  name = "doc",
  srcs = glob(["doc/*.md"]),
)

filegroup(
  name = "geodata",
  srcs = ["config/geoip.dat", "config/geosite.dat"],
)

pkg_zip(
  name = "v2ray_linux_amd64_package",
  srcs = [
    ":v2ray_linux_amd64",
    ":v2ray_linux_amd64_sig",
    "@v2ray_ext//tools/control/main:v2ctl_linux_amd64",
    "@v2ray_ext//tools/control/main:v2ctl_linux_amd64_sig",
    ":config_all",
    ":systemd",
    ":doc",
    ":geodata",
  ],
  out = "v2ray-linux-64.zip",
)
