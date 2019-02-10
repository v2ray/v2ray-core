load("//infra/bazel:build.bzl", "foreign_go_binary")
load("//infra/bazel:gpg.bzl", "gpg_sign")
load("//infra/bazel:matrix.bzl", "SUPPORTED_MATRIX")
load("//main:targets.bzl", "gen_targets")

package(default_visibility=["//visibility:public"])

gen_targets(SUPPORTED_MATRIX)
