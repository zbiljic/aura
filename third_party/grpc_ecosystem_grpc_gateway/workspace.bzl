"""Provides the repository macro to import grpc-gateway."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def repo():
    """Imports grpc-gateway."""

    GRPC_GATEWAY_VERSION = "2.16.0"
    GRPC_GATEWAY_SHA256 = "26f3983a6db6844d843e3a701e38f561d862e46a6081ee90c14f6a23a46c2dbc"

    http_archive(
        name = "com_github_grpc_ecosystem_grpc_gateway_v2",
        patches = ["//third_party/grpc_ecosystem_grpc_gateway:rules.patch"],
        sha256 = GRPC_GATEWAY_SHA256,
        strip_prefix = "grpc-gateway-{version}".format(version = GRPC_GATEWAY_VERSION),
        urls = [
            "https://github.com/grpc-ecosystem/grpc-gateway/archive/refs/tags/v{version}.zip".format(version = GRPC_GATEWAY_VERSION),
        ],
    )
