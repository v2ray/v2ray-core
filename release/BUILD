package(default_visibility = ["//visibility:public"])

load("//infra/bazel:zip.bzl", "pkg_zip")
load("//release:mapping.bzl", "gen_mappings")

filegroup(
    name = "config_json",
    srcs = [
        "config/config.json",
        "config/vpoint_socks_vmess.json",
        "config/vpoint_vmess_freedom.json",
    ],
)

filegroup(
    name = "systemd",
    srcs = [
        "config/systemd/system/v2ray.service",
        "config/systemd/system/v2ray@.service",
    ],
)

filegroup(
    name = "doc",
    srcs = glob(["doc/*.md"]),
)

filegroup(
    name = "geodata",
    srcs = [
        "config/geoip.dat",
        "config/geosite.dat",
    ],
)

pkg_zip(
    name = "v2ray_darwin_amd64_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        "//infra/control/main:v2ctl_darwin_amd64",
        "//infra/control/main:v2ctl_darwin_amd64_sig",
        "//main:v2ray_darwin_amd64",
        "//main:v2ray_darwin_amd64_sig",
    ],
    out = "v2ray-macos-64.zip",
    mappings = gen_mappings("darwin", "amd64", "0"),
)

pkg_zip(
    name = "v2ray_windows_amd64_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        "//infra/control/main:v2ctl_windows_amd64",
        "//infra/control/main:v2ctl_windows_amd64_sig",
        "//main:v2ray_windows_amd64",
        "//main:v2ray_windows_amd64_sig",
        "//main:v2ray_windows_amd64_nowindow",
        "//main:v2ray_windows_amd64_nowindow_sig",
    ],
    out = "v2ray-windows-64.zip",
    mappings = gen_mappings("windows", "amd64", "0"),
)

pkg_zip(
    name = "v2ray_windows_x86_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        "//infra/control/main:v2ctl_windows_386",
        "//infra/control/main:v2ctl_windows_386_sig",
        "//main:v2ray_windows_386",
        "//main:v2ray_windows_386_sig",
        "//main:v2ray_windows_386_nowindow",
        "//main:v2ray_windows_386_nowindow_sig",
    ],
    out = "v2ray-windows-32.zip",
    mappings = gen_mappings("windows", "386", "0"),
)

pkg_zip(
    name = "v2ray_windows_armv7_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        "//infra/control/main:v2ctl_windows_arm_7",
        "//infra/control/main:v2ctl_windows_arm_7_sig",
        "//main:v2ray_windows_arm_7",
        "//main:v2ray_windows_arm_7_sig",
        "//main:v2ray_windows_arm_7_nowindow",
        "//main:v2ray_windows_arm_7_nowindow_sig",
    ],
    out = "v2ray-windows-arm32-v7a.zip",
    mappings = gen_mappings("windows", "arm", "7"),
)

pkg_zip(
    name = "v2ray_freebsd_amd64_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        "//infra/control/main:v2ctl_freebsd_amd64",
        "//infra/control/main:v2ctl_freebsd_amd64_sig",
        "//main:v2ray_freebsd_amd64",
        "//main:v2ray_freebsd_amd64_sig",
    ],
    out = "v2ray-freebsd-64.zip",
    mappings = gen_mappings("freebsd", "amd64", "0"),
)

pkg_zip(
    name = "v2ray_freebsd_x86_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        "//infra/control/main:v2ctl_freebsd_386",
        "//infra/control/main:v2ctl_freebsd_386_sig",
        "//main:v2ray_freebsd_386",
        "//main:v2ray_freebsd_386_sig",
    ],
    out = "v2ray-freebsd-32.zip",
    mappings = gen_mappings("freebsd", "386", "0"),
)

pkg_zip(
    name = "v2ray_openbsd_amd64_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        "//infra/control/main:v2ctl_openbsd_amd64",
        "//infra/control/main:v2ctl_openbsd_amd64_sig",
        "//main:v2ray_openbsd_amd64",
        "//main:v2ray_openbsd_amd64_sig",
    ],
    out = "v2ray-openbsd-64.zip",
    mappings = gen_mappings("openbsd", "amd64", "0"),
)

pkg_zip(
    name = "v2ray_openbsd_x86_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        "//infra/control/main:v2ctl_openbsd_386",
        "//infra/control/main:v2ctl_openbsd_386_sig",
        "//main:v2ray_openbsd_386",
        "//main:v2ray_openbsd_386_sig",
    ],
    out = "v2ray-openbsd-32.zip",
    mappings = gen_mappings("openbsd", "386", "0"),
)

pkg_zip(
    name = "v2ray_dragonfly_amd64_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        "//infra/control/main:v2ctl_dragonfly_amd64",
        "//infra/control/main:v2ctl_dragonfly_amd64_sig",
        "//main:v2ray_dragonfly_amd64",
        "//main:v2ray_dragonfly_amd64_sig",
    ],
    out = "v2ray-dragonfly-64.zip",
    mappings = gen_mappings("dragonfly", "amd64", "0"),
)

pkg_zip(
    name = "v2ray_linux_amd64_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_amd64",
        "//infra/control/main:v2ctl_linux_amd64_sig",
        "//main:v2ray_linux_amd64",
        "//main:v2ray_linux_amd64_sig",
    ],
    out = "v2ray-linux-64.zip",
    mappings = gen_mappings("linux", "amd64", "0"),
)

pkg_zip(
    name = "v2ray_linux_x86_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_386",
        "//infra/control/main:v2ctl_linux_386_sig",
        "//main:v2ray_linux_386",
        "//main:v2ray_linux_386_sig",
    ],
    out = "v2ray-linux-32.zip",
    mappings = gen_mappings("linux", "386", "0"),
)

pkg_zip(
    name = "v2ray_linux_arm64_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_arm64",
        "//infra/control/main:v2ctl_linux_arm64_sig",
        "//main:v2ray_linux_arm64",
        "//main:v2ray_linux_arm64_sig",
    ],
    out = "v2ray-linux-arm64-v8a.zip",
    mappings = gen_mappings("linux", "arm64", "0"),
)

pkg_zip(
    name = "v2ray_linux_armv7_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_arm_7",
        "//infra/control/main:v2ctl_linux_arm_7_sig",
        "//main:v2ray_linux_arm_7",
        "//main:v2ray_linux_arm_7_sig",
    ],
    out = "v2ray-linux-arm32-v7a.zip",
    mappings = gen_mappings("linux", "arm", "7"),
)

pkg_zip(
    name = "v2ray_linux_armv6_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_arm_6",
        "//infra/control/main:v2ctl_linux_arm_6_sig",
        "//main:v2ray_linux_arm_6",
        "//main:v2ray_linux_arm_6_sig",
    ],
    out = "v2ray-linux-arm32-v6.zip",
    mappings = gen_mappings("linux", "arm", "6"),
)

pkg_zip(
    name = "v2ray_linux_armv5_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_arm_5",
        "//infra/control/main:v2ctl_linux_arm_5_sig",
        "//main:v2ray_linux_arm_5",
        "//main:v2ray_linux_arm_5_sig",
    ],
    out = "v2ray-linux-arm32-v5.zip",
    mappings = gen_mappings("linux", "arm", "5"),
)

pkg_zip(
    name = "v2ray_linux_mips32_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_mips",
        "//infra/control/main:v2ctl_linux_mips_sig",
        "//infra/control/main:v2ctl_linux_mips_softfloat",
        "//infra/control/main:v2ctl_linux_mips_softfloat_sig",
        "//main:v2ray_linux_mips",
        "//main:v2ray_linux_mips_sig",
        "//main:v2ray_linux_mips_softfloat",
        "//main:v2ray_linux_mips_softfloat_sig",
    ],
    out = "v2ray-linux-mips32.zip",
    mappings = gen_mappings("linux", "mips", "0"),
)

pkg_zip(
    name = "v2ray_linux_mips32le_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_mipsle",
        "//infra/control/main:v2ctl_linux_mipsle_sig",
        "//infra/control/main:v2ctl_linux_mipsle_softfloat",
        "//infra/control/main:v2ctl_linux_mipsle_softfloat_sig",
        "//main:v2ray_linux_mipsle",
        "//main:v2ray_linux_mipsle_sig",
        "//main:v2ray_linux_mipsle_softfloat",
        "//main:v2ray_linux_mipsle_softfloat_sig",
    ],
    out = "v2ray-linux-mips32le.zip",
    mappings = gen_mappings("linux", "mipsle", "0"),
)

pkg_zip(
    name = "v2ray_linux_mips64_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_mips64",
        "//infra/control/main:v2ctl_linux_mips64_sig",
        "//main:v2ray_linux_mips64",
        "//main:v2ray_linux_mips64_sig",
    ],
    out = "v2ray-linux-mips64.zip",
    mappings = gen_mappings("linux", "mips64", "0"),
)

pkg_zip(
    name = "v2ray_linux_mips64le_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_mips64le",
        "//infra/control/main:v2ctl_linux_mips64le_sig",
        "//main:v2ray_linux_mips64le",
        "//main:v2ray_linux_mips64le_sig",
    ],
    out = "v2ray-linux-mips64le.zip",
    mappings = gen_mappings("linux", "mips64le", "0"),
)

pkg_zip(
    name = "v2ray_linux_riscv64_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_riscv64",
        "//infra/control/main:v2ctl_linux_riscv64_sig",
        "//main:v2ray_linux_riscv64",
        "//main:v2ray_linux_riscv64_sig",
    ],
    out = "v2ray-linux-riscv64.zip",
    mappings = gen_mappings("linux", "riscv64", "0"),
)

pkg_zip(
    name = "v2ray_linux_s390x_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_s390x",
        "//infra/control/main:v2ctl_linux_s390x_sig",
        "//main:v2ray_linux_s390x",
        "//main:v2ray_linux_s390x_sig",
    ],
    out = "v2ray-linux-s390x.zip",
    mappings = gen_mappings("linux", "s390x", "0"),
)

pkg_zip(
    name = "v2ray_linux_ppc64_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_ppc64",
        "//infra/control/main:v2ctl_linux_ppc64_sig",
        "//main:v2ray_linux_ppc64",
        "//main:v2ray_linux_ppc64_sig",
    ],
    out = "v2ray-linux-ppc64.zip",
    mappings = gen_mappings("linux", "ppc64", "0"),
)

pkg_zip(
    name = "v2ray_linux_ppc64le_package",
    srcs = [
        ":config_json",
        ":doc",
        ":geodata",
        ":systemd",
        "//infra/control/main:v2ctl_linux_ppc64le",
        "//infra/control/main:v2ctl_linux_ppc64le_sig",
        "//main:v2ray_linux_ppc64le",
        "//main:v2ray_linux_ppc64le_sig",
    ],
    out = "v2ray-linux-ppc64le.zip",
    mappings = gen_mappings("linux", "ppc64le", "0"),
)
