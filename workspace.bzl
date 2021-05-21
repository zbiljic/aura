"""Workspace initialization. Consult the WORKSPACE on how to use it."""

# Import third party repository rules.

# Import external repository rules.
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _initialize_third_party():
    """ Load third party repositories.  See above load() statements. """
    pass

# Define all external repositories.
def _repositories():
    """All external dependencies."""

    # To update any of the dependencies bellow:
    # a) update URL and strip_prefix to the new git commit hash
    # b) get the sha256 hash of the commit by running:
    #    curl -L <url> | sha256sum
    # and update the sha256 with the result.

    pass

def workspace():
    # Import third party repositories.
    _initialize_third_party()

    # Import all other repositories. This should happen before initializing
    # any external repositories, because those come with their own
    # dependencies.
    _repositories()
