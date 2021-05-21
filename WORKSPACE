workspace(name = "com_github_zbiljic_aura")

######################
# Common
######################

#
# Bazel infrastructure
#
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

#
# Skylib
#
http_archive(
    name = "bazel_skylib",
    sha256 = "1c531376ac7e5a180e0237938a2536de0c54d93f5c278634818e0efc952dd56c",
    urls = [
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.0.3/bazel-skylib-1.0.3.tar.gz",
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.0.3/bazel-skylib-1.0.3.tar.gz",
    ],
)

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

#
# Buildtools for bazel
#
http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "c28eef4d30ba1a195c6837acf6c75a4034981f5b4002dda3c5aa6e48ce023cf1",
    strip_prefix = "buildtools-4.0.1",
    urls = ["https://github.com/bazelbuild/buildtools/archive/4.0.1.tar.gz"],
)

######################
# Go support
######################

GO_VERSION = "1.16.4"

#
# Go Rules
#
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "69de5c704a05ff37862f7e0f5534d4f479418afc21806c887db544a316f3cb6b",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.27.0/rules_go-v0.27.0.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.27.0/rules_go-v0.27.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(go_version = GO_VERSION)

load("@io_bazel_rules_go//extras:embed_data_deps.bzl", "go_embed_data_dependencies")

go_embed_data_dependencies()

#
# Gazelle (creates Bazel rules from standard Go builds)
#
http_archive(
    name = "bazel_gazelle",
    sha256 = "62ca106be173579c0a167deb23358fdfe71ffa1e4cfdddf5582af26520f1c66f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

#
# Go external tools
#

http_archive(
    name = "com_github_ash2k_bazel_tools",
    sha256 = "4db4de3839534a036cc9f7505f24b75133b0d4542a993e6ccc1ac0c91ea1f50d",
    strip_prefix = "bazel-tools-e8c576d666521401ea1471eb6fd5444583c6702c",
    urls = ["https://github.com/ash2k/bazel-tools/archive/e8c576d666521401ea1471eb6fd5444583c6702c.zip"],
)

load("@com_github_ash2k_bazel_tools//buildozer:deps.bzl", "buildozer_dependencies")
load("@com_github_ash2k_bazel_tools//goimports:deps.bzl", "goimports_dependencies")
load("@com_github_ash2k_bazel_tools//golangcilint:deps.bzl", "golangcilint_dependencies")

buildozer_dependencies()

goimports_dependencies()

golangcilint_dependencies()

######################
# Protobuf support
######################

http_archive(
    name = "rules_proto",
    sha256 = "9fc210a34f0f9e7cc31598d109b5d069ef44911a82f507d5a88716db171615a8",
    strip_prefix = "rules_proto-f7a30f6f80006b591fa7c437fe5a951eb10bcbcf",
    urls = [
        "https://github.com/bazelbuild/rules_proto/archive/f7a30f6f80006b591fa7c437fe5a951eb10bcbcf.tar.gz",
    ],
)

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")

rules_proto_dependencies()

rules_proto_toolchains()

######################
# External dependencies
######################

# Initialize the dependencies.
load("//:workspace.bzl", "workspace")

workspace()

load("//third_party:go_workspace.bzl", "go_dependencies")

# gazelle:repository_macro third_party/go_workspace.bzl%go_dependencies
go_dependencies()
