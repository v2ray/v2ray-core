def gen_mappings(os, arch, ver):
  return {
    "v2ray_core/release/config": "",
    "v2ray_core/main/" + os + "/" + arch + "/" + ver: "",
    "v2ray_core/infra/control/main/" + os + "/" + arch + "/" + ver : "",
  }
