load("//infra/bazel:build.bzl", "foreign_go_binary")

def gen_targets(matrix):
  pkg = "./main"
  output = "v2ray"

  for (os, arch, ver) in matrix:

    if arch in ["arm"]:
      bin_name = "v2ray_" + os + "_" + arch + "_" + ver
      foreign_go_binary(
        name = bin_name,
        pkg = pkg,
        output = output,
        os = os,
        arch = arch,
        ver = ver,
        arm = ver,
      )

      if os in ["windows"]:
        bin_name = "v2ray_" + os + "_" + arch + "_" + ver + "_nowindow"
        foreign_go_binary(
          name = bin_name,
          pkg = pkg,
          output = "w" + output,
          os = os,
          arch = arch,
          ver = ver,
          arm = ver,
          ld = "-H windowsgui",
        )

    else:
      bin_name = "v2ray_" + os + "_" + arch
      foreign_go_binary(
        name = bin_name,
        pkg = pkg,
        output = output,
        os = os,
        arch = arch,
        ver = ver,
      )

      if os in ["windows"]:
        bin_name = "v2ray_" + os + "_" + arch + "_nowindow"
        foreign_go_binary(
          name = bin_name,
          pkg = pkg,
          output = "w" + output,
          os = os,
          arch = arch,
          ver = ver,
          ld = "-H windowsgui",
        )

      if arch in ["mips", "mipsle"]:
        bin_name = "v2ray_" + os + "_" + arch + "_softfloat"
        foreign_go_binary(
          name = bin_name,
          pkg = pkg,
          output = output + "_softfloat",
          os = os,
          arch = arch,
          ver = ver,
          mips = "softfloat",
        )
