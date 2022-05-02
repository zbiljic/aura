"""Provides the repository macro to import grpc-gateway."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def repo():
    """Imports grpc-gateway."""

    GRPC_GATEWAY_VERSION = "2.10.0"
    GRPC_GATEWAY_SHA256 = "e9ab6d341171c616174f4a328c4315ed51f11db8ddff02ee72e10149a796e0ed"

    http_archive(
        name = "com_github_grpc_ecosystem_grpc_gateway_v2",
        patches = ["//third_party/grpc_ecosystem_grpc_gateway:rules.patch"],
        sha256 = GRPC_GATEWAY_SHA256,
        strip_prefix = "grpc-gateway-{version}".format(version = GRPC_GATEWAY_VERSION),
        urls = [
            "https://github.com/grpc-ecosystem/grpc-gateway/archive/refs/tags/v{version}.zip".format(version = GRPC_GATEWAY_VERSION),
        ],
    )
