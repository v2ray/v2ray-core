load("//bazel:build.bzl", "foreign_go_binary")
load("//bazel:gpg.bzl", "gpg_sign")

def gen_targets(matrix):
  output = "v2ctl"
  pkg = "v2ray.com/core/infra/control/main"

  for (os, arch) in matrix:
    bin_name = "v2ctl_" + os + "_" + arch
    foreign_go_binary(
      name = bin_name,
      pkg = pkg,
      output = output,
      os = os,
      arch = arch,
      gotags = "confonly",
    )

    gpg_sign(
      name = bin_name + "_sig",
      base = ":" + bin_name,
    )

    if arch in ["mips", "mipsle"]:
      bin_name = "v2ctl_" + os + "_" + arch + "_softfloat"
      foreign_go_binary(
        name = bin_name,
        pkg = pkg,
        output = output + "_softfloat",
        os = os,
        arch = arch,
        mips = "softfloat",
        gotags = "confonly",
      )

      gpg_sign(
        name = bin_name + "_sig",
        base = ":" + bin_name,
      )
    
    if arch in ["arm"]:
      bin_name = "v2ctl_" + os + "_" + arch + "_armv7"
      foreign_go_binary(
        name = bin_name,
        pkg = pkg,
        output = output + "_armv7",
        os = os,
        arch = arch,
        arm = "7",
        gotags = "confonly",
      )

      gpg_sign(
        name = bin_name + "_sig",
        base = ":" + bin_name,
      )
