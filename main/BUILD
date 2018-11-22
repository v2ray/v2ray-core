load("@v2ray_ext//bazel:build.bzl", "foreign_go_binary")
load("@v2ray_ext//bazel:gpg.bzl", "gpg_sign")
load("@v2ray_ext//bazel:matrix.bzl", "SUPPORTED_MATRIX")
load("//main:targets.bzl", "gen_targets")

package(default_visibility=["//visibility:public"])

gen_targets(SUPPORTED_MATRIX)
