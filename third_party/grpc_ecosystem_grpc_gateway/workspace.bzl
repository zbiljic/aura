"""Provides the repository macro to import grpc-gateway."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def repo():
    """Imports grpc-gateway."""

    GRPC_GATEWAY_VERSION = "2.7.3"
    GRPC_GATEWAY_SHA256 = "851e202014f1a086fb0c4465e68a7ca978967d4a1682e52310eea3107ff46eed"

    http_archive(
        name = "com_github_grpc_ecosystem_grpc_gateway_v2",
        patches = ["//third_party/grpc_ecosystem_grpc_gateway:rules.patch"],
        sha256 = GRPC_GATEWAY_SHA256,
        strip_prefix = "grpc-gateway-{version}".format(version = GRPC_GATEWAY_VERSION),
        urls = [
            "https://github.com/grpc-ecosystem/grpc-gateway/archive/refs/tags/v{version}.zip".format(version = GRPC_GATEWAY_VERSION),
        ],
    )
