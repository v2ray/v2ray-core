load("//infra/bazel:build.bzl", "foreign_go_binary")
load("//infra/bazel:gpg.bzl", "gpg_sign")

def gen_targets(matrix):
  pkg = "v2ray.com/core/infra/control/main"
  output = "v2ctl"

  for (os, arch, ver) in matrix:

    if arch in ["arm"]:
      bin_name = "v2ctl_" + os + "_" + arch + "_" + ver
      foreign_go_binary(
        name = bin_name,
        pkg = pkg,
        output = output,
        os = os,
        arch = arch,
        ver = ver,
        arm = ver,
        gotags = "confonly",
      )

      gpg_sign(
        name = bin_name + "_sig",
        base = ":" + bin_name,
      )

    else:
      bin_name = "v2ctl_" + os + "_" + arch
      foreign_go_binary(
        name = bin_name,
        pkg = pkg,
        output = output,
        os = os,
        arch = arch,
        ver = ver,
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
          ver = ver,
          mips = "softfloat",
          gotags = "confonly",
        )

        gpg_sign(
          name = bin_name + "_sig",
          base = ":" + bin_name,
        )
