def _gpg_sign_impl(ctx):
  output_file = ctx.actions.declare_file(ctx.file.base.basename + ctx.attr.suffix, sibling = ctx.file.base)
  if not ctx.configuration.default_shell_env.get("GPG_PASS"):
    ctx.actions.write(output_file, "")
  else:
    command = "echo ${GPG_PASS} | gpg --pinentry-mode loopback --digest-algo SHA512 --passphrase-fd 0 --output %s --detach-sig %s" % (output_file.path, ctx.file.base.path)
    ctx.actions.run_shell(
      command = command,
      use_default_shell_env = True,
      inputs = [ctx.file.base],
      outputs = [output_file],
      progress_message = "Signing binary",
      mnemonic = "gpg",
    )
  return [DefaultInfo(files = depset([output_file]))]

gpg_sign = rule(
  implementation = _gpg_sign_impl,
  attrs = {
    "base": attr.label(allow_single_file=True),
    "suffix": attr.string(default=".sig"),
  },
)
