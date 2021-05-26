"""Provides the repository macro to import golink."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def repo():
    """Imports golink."""

    GOLINK_VERSION = "1.1.0"
    GOLINK_SHA256 = "c505a82b7180d4315bbaf05848e9b7d2683e80f1b16159af51a0ecae6fb2d54d"

    http_archive(
        name = "golink",
        patches = ["//third_party/golink:golink.patch"],
        sha256 = GOLINK_SHA256,
        strip_prefix = "golink-{version}".format(version = GOLINK_VERSION),
        urls = [
            "https://github.com/nikunjy/golink/archive/v{version}.tar.gz".format(version = GOLINK_VERSION),
        ],
    )
