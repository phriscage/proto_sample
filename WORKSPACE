workspace(name = "com_github_phriscage_proto_sample")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

PROTO_GRPC_RULE_SHA = "9ba7299c5eb6ec45b6b9a0ceb9916d0ab96789ac8218269322f0124c0c0d24e2"

PROTO_GRPC_RULE_VERSION = "4.5.0"

http_archive(
    name = "rules_proto_grpc",
    sha256 = PROTO_GRPC_RULE_SHA,
    strip_prefix = "rules_proto_grpc-%s" % PROTO_GRPC_RULE_VERSION,
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/releases/download/%s/rules_proto_grpc-%s.tar.gz" % (PROTO_GRPC_RULE_VERSION, PROTO_GRPC_RULE_VERSION)],
)

load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_repos", "rules_proto_grpc_toolchains")
rules_proto_grpc_toolchains()
rules_proto_grpc_repos()

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")
rules_proto_dependencies()
rules_proto_toolchains()

load("@rules_proto_grpc//:repositories.bzl", "bazel_gazelle", "io_bazel_rules_go")  # buildifier: disable=same-origin-load

io_bazel_rules_go()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

bazel_gazelle()

load("@rules_proto_grpc//grpc-gateway:repositories.bzl", rules_proto_grpc_gateway_repos = "gateway_repos")

rules_proto_grpc_gateway_repos()

load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

go_register_toolchains(
    version = "1.20.10",
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()
