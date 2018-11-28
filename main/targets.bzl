load("@v2ray_ext//bazel:build.bzl", "foreign_go_binary")
load("@v2ray_ext//bazel:gpg.bzl", "gpg_sign")
load("@v2ray_ext//bazel:plugin.bzl", "PLUGIN_SUPPORTED_OS")

def gen_targets(matrix):
  pkg = "v2ray.com/core/main"
  output = "v2ray"

  for (os, arch) in matrix:

    cgo_enabled = "0"
    if os in PLUGIN_SUPPORTED_OS:
      cgo_enabled = "1"

    bin_name = "v2ray_" + os + "_" + arch
    foreign_go_binary(
      name = bin_name,
      pkg = pkg,
      output = output,
      os = os,
      arch = arch,
      cgo_enabled = cgo_enabled,
    )

    gpg_sign(
      name = bin_name + "_sig",
      base = ":" + bin_name,
    )

    if os in ["windows"]:
      bin_name = "v2ray_" + os + "_" + arch + "_nowindow"
      foreign_go_binary(
        name = bin_name,
        pkg = pkg,
        output = "w" + output,
        os = os,
        arch = arch,
        ld = "-H windowsgui",
      )

      gpg_sign(
        name = bin_name + "_sig",
        base = ":" + bin_name,
      )

    if arch in ["mips", "mipsle"]:
      bin_name = "v2ray_" + os + "_" + arch + "_softfloat"
      foreign_go_binary(
        name = bin_name,
        pkg = pkg,
        output = output+"_softfloat",
        os = os,
        arch = arch,
        mips = "softfloat",
      )

      gpg_sign(
        name = bin_name + "_sig",
        base = ":" + bin_name,
      )
    
    if arch in ["arm"]:
      bin_name = "v2ray_" + os + "_" + arch + "_armv7"
      foreign_go_binary(
        name = bin_name,
        pkg = pkg,
        output = output+"_armv7",
        os = os,
        arch = arch,
        arm = "7",
      )

      gpg_sign(
        name = bin_name + "_sig",
        base = ":" + bin_name,
      )
