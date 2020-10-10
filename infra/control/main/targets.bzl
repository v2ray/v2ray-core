load("//infra/bazel:build.bzl", "foreign_go_binary")

def gen_targets(matrix):
  pkg = "./infra/control/main"
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
