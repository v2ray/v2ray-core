package(default_visibility=["//visibility:public"])

load("@v2ray_ext//bazel:zip.bzl", "pkg_zip")
load("//release:mapping.bzl", "gen_mappings")

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
  name = "v2ray_darwin_amd64_package",
  srcs = [
    "//main:v2ray_darwin_amd64",
    "//main:v2ray_darwin_amd64_sig",
    "@v2ray_ext//tools/control/main:v2ctl_darwin_amd64",
    "@v2ray_ext//tools/control/main:v2ctl_darwin_amd64_sig",
    ":config_json",
    ":doc",
    ":geodata",
  ],
  out = "v2ray-macos-64.zip",
  mappings = gen_mappings("darwin", "amd64"),
)

pkg_zip(
  name = "v2ray_windows_amd64_package",
  srcs = [
    "//main:v2ray_windows_amd64",
    "//main:v2ray_lwindows_amd64_sig",
    "//main:v2ray_windows_amd64_nowindow",
    "//main:v2ray_lwindows_amd64_nowindow_sig",
    "@v2ray_ext//tools/control/main:v2ctl_windows_amd64",
    "@v2ray_ext//tools/control/main:v2ctl_windows_amd64_sig",
    ":config_json",
    ":doc",
    ":geodata",
  ],
  out = "v2ray-windows-64.zip",
)

pkg_zip(
  name = "v2ray_windows_x86_package",
  srcs = [
    "//main:v2ray_windows_x86",
    "//main:v2ray_windows_x86_sig",
    "//main:v2ray_windows_x86_nowindow",
    "//main:v2ray_windows_x86_nowindow_sig",
    "@v2ray_ext//tools/control/main:v2ctl_windows_x86",
    "@v2ray_ext//tools/control/main:v2ctl_windows_x86_sig",
    ":config_json",
    ":doc",
    ":geodata",
  ],
  out = "v2ray-windows-32.zip",
)

pkg_zip(
  name = "v2ray_linux_amd64_package",
  srcs = [
    "//main:v2ray_linux_amd64",
    "//main:v2ray_linux_amd64_sig",
    "@v2ray_ext//tools/control/main:v2ctl_linux_amd64",
    "@v2ray_ext//tools/control/main:v2ctl_linux_amd64_sig",
    ":config_all",
    ":systemd",
    ":doc",
    ":geodata",
  ],
  out = "v2ray-linux-64.zip",
)

pkg_zip(
  name = "v2ray_linux_x86_package",
  srcs = [
    "//main:v2ray_linux_x86",
    "//main:v2ray_linux_x86_sig",
    "@v2ray_ext//tools/control/main:v2ctl_linux_x86",
    "@v2ray_ext//tools/control/main:v2ctl_linux_x86_sig",
    ":config_all",
    ":systemd",
    ":doc",
    ":geodata",
  ],
  out = "v2ray-linux-32.zip",
)
