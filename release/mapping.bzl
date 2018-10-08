def gen_mappings(os, arch):
  return {
    "v2ray_core/release/doc": "doc",
    "v2ray_core/release/config": "",
    "v2ray_core/main/" + os + "/" + arch: "",
    "v2ray_ext/tools/control/main/" + os + "/" + arch: "",
  }
