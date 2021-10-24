"""Provides the repository macro to import grpc-gateway."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def repo():
    """Imports grpc-gateway."""

    GRPC_GATEWAY_VERSION = "2.6.0"
    GRPC_GATEWAY_SHA256 = "8d7f101db6c458f3d263c823da224a4df05a413847673d1255d2fcce9deddd1f"

    http_archive(
        name = "com_github_grpc_ecosystem_grpc_gateway_v2",
        patches = ["//third_party/grpc_ecosystem_grpc_gateway:rules.patch"],
        sha256 = GRPC_GATEWAY_SHA256,
        strip_prefix = "grpc-gateway-{version}".format(version = GRPC_GATEWAY_VERSION),
        urls = [
            "https://github.com/grpc-ecosystem/grpc-gateway/archive/refs/tags/v{version}.zip".format(version = GRPC_GATEWAY_VERSION),
        ],
    )
