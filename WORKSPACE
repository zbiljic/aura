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
    sha256 = "af87959afe497dc8dfd4c6cb66e1279cb98ccc84284619ebfec27d9c09a903de",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.2.0/bazel-skylib-1.2.0.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.2.0/bazel-skylib-1.2.0.tar.gz",
    ],
)

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

#
# Buildtools for bazel
#
http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "ae34c344514e08c23e90da0e2d6cb700fcd28e80c02e23e4d5715dddcb42f7b3",
    strip_prefix = "buildtools-4.2.2",
    urls = ["https://github.com/bazelbuild/buildtools/archive/4.2.2.tar.gz"],
)

######################
# Go support
######################

GO_VERSION = "1.17.2"

#
# Go Rules
#
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "2b1641428dff9018f9e85c0384f03ec6c10660d935b750e3fa1492a281a53b0f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.29.0/rules_go-v0.29.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.29.0/rules_go-v0.29.0.zip",
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
    sha256 = "de69a09dc70417580aabf20a28619bb3ef60d038470c7cf8442fafcf627c21cb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
    ],
)

#
# Go external tools
#

http_archive(
    name = "com_github_ash2k_bazel_tools",
    sha256 = "5fb439f08ad365258008e9a2265cda0d1e8a1cebf0ab2811c50ec2c9d301e266",
    strip_prefix = "bazel-tools-f8b27b99cae951099385655e0bb0fc9cc1c7baa4",
    urls = ["https://github.com/ash2k/bazel-tools/archive/f8b27b99cae951099385655e0bb0fc9cc1c7baa4.zip"],
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
    sha256 = "83c8798f5a4fe1f6a13b5b6ae4267695b71eed7af6fbf2b6ec73a64cf01239ab",
    strip_prefix = "rules_proto-b22f78685bf62775b80738e766081b9e4366cdf0",
    urls = [
        "https://github.com/bazelbuild/rules_proto/archive/b22f78685bf62775b80738e766081b9e4366cdf0.tar.gz",
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

# Load Gazelle dependencies after others
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
