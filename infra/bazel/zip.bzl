# Copied from google/nomulus project as we don't want to import the whole repository.

ZIPPER = "@bazel_tools//tools/zip:zipper"

def long_path(ctx, file_):
    """Constructs canonical runfile path relative to TEST_SRCDIR.
    Args:
      ctx: A Skylark rule context.
      file_: A File object that should appear in the runfiles for the test.
    Returns:
      A string path relative to TEST_SRCDIR suitable for use in tests and
      testing infrastructure.
    """
    if file_.short_path.startswith("../"):
        return file_.short_path[3:]
    if file_.owner and file_.owner.workspace_root:
        return file_.owner.workspace_root + "/" + file_.short_path
    return ctx.workspace_name + "/" + file_.short_path

def collect_runfiles(targets):
    """Aggregates runfiles from targets.
    Args:
      targets: A list of Bazel targets.
    Returns:
      A list of Bazel files.
    """
    data = depset()
    for target in targets:
        if hasattr(target, "runfiles"):
            data += target.runfiles.files
            continue
        if hasattr(target, "data_runfiles"):
            data += target.data_runfiles.files
        if hasattr(target, "default_runfiles"):
            data += target.default_runfiles.files
    return data

def _get_runfiles(target, attribute):
    runfiles = getattr(target, attribute, None)
    if runfiles:
        return runfiles.files
    return []

def _zip_file(ctx):
    """Implementation of zip_file() rule."""
    for s, d in ctx.attr.mappings.items():
        if (s.startswith("/") or s.endswith("/") or
            d.startswith("/") or d.endswith("/")):
            fail("mappings should not begin or end with slash")
    srcs = depset(transitive = [depset(ctx.files.srcs),depset(ctx.files.data),depset(collect_runfiles(ctx.attr.data))])
    # srcs += ctx.files.srcs
    # srcs += ctx.files.data
    # srcs += collect_runfiles(ctx.attr.data)
    mapped = _map_sources(ctx, srcs, ctx.attr.mappings)
    cmd = [
        "#!/bin/sh",
        "set -e",
        'repo="$(pwd)"',
        'zipper="${repo}/%s"' % ctx.file._zipper.path,
        'archive="${repo}/%s"' % ctx.outputs.out.path,
        'tmp="$(mktemp -d "${TMPDIR:-/tmp}/zip_file.XXXXXXXXXX")"',
        'cd "${tmp}"',
    ]
    cmd += [
        '"${zipper}" x "${repo}/%s"' % dep.zip_file.path
        for dep in ctx.attr.deps
    ]
    cmd += ["rm %s" % filename for filename in ctx.attr.exclude]
    cmd += [
        'mkdir -p "${tmp}/%s"' % zip_path
        for zip_path in depset(
            [
                zip_path[:zip_path.rindex("/")]
                for _, zip_path in mapped
                if "/" in zip_path
            ],
        ).to_list()
    ]
    cmd += [
        'ln -sf "${repo}/%s" "${tmp}/%s"' % (path, zip_path)
        for path, zip_path in mapped
    ]
    cmd += [
        ("find . | sed 1d | cut -c 3- | LC_ALL=C sort" +
         ' | xargs "${zipper}" cC "${archive}"'),
        'cd "${repo}"',
        'rm -rf "${tmp}"',
    ]
    script = ctx.actions.declare_file("%s/%s.sh" % (ctx.bin_dir, ctx.label.name))
    ctx.actions.write(output = script, content = "\n".join(cmd), is_executable = True)
    inputs = [ctx.file._zipper]
    inputs += [dep.zip_file for dep in ctx.attr.deps]
    inputs += list(srcs.to_list())
    ctx.actions.run(
        inputs = inputs,
        outputs = [ctx.outputs.out],
        executable = script,
        mnemonic = "zip",
        progress_message = "Creating zip with %d inputs %s" % (
            len(inputs),
            ctx.label,
        ),
    )
    return struct(files = depset([ctx.outputs.out]), zip_file = ctx.outputs.out)

def _map_sources(ctx, srcs, mappings):
    """Calculates paths in zip file for srcs."""

    # order mappings with more path components first
    mappings = sorted([
        (-len(source.split("/")), source, dest)
        for source, dest in mappings.items()
    ])

    # get rid of the integer part of tuple used for sorting
    mappings = [(source, dest) for _, source, dest in mappings]
    mappings_indexes = range(len(mappings))
    used = {i: False for i in mappings_indexes}
    mapped = []
    for file_ in srcs.to_list():
        run_path = long_path(ctx, file_)
        zip_path = None
        for i in mappings_indexes:
            source = mappings[i][0]
            dest = mappings[i][1]
            if not source:
                if dest:
                    zip_path = dest + "/" + run_path
                else:
                    zip_path = run_path
            elif source == run_path:
                if dest:
                    zip_path = dest
                else:
                    zip_path = run_path
            elif run_path.startswith(source + "/"):
                if dest:
                    zip_path = dest + run_path[len(source):]
                else:
                    zip_path = run_path[len(source) + 1:]
            else:
                continue
            used[i] = True
            break
        if not zip_path:
            fail("no mapping matched: " + run_path)
        mapped += [(file_.path, zip_path)]
    for i in mappings_indexes:
        if not used[i]:
            fail('superfluous mapping: "%s" -> "%s"' % mappings[i])
    return mapped

pkg_zip = rule(
    implementation = _zip_file,
    attrs = {
        "out": attr.output(mandatory = True),
        "srcs": attr.label_list(allow_files = True),
        "data": attr.label_list(allow_files = True),
        "deps": attr.label_list(providers = ["zip_file"]),
        "exclude": attr.string_list(),
        "mappings": attr.string_dict(),
        "_zipper": attr.label(default = Label(ZIPPER), allow_single_file = True),
    },
)
