licenses(["notice"])  # Apache 2

load(
    "@envoy//bazel:envoy_build_system.bzl",
    "envoy_cc_test",
)

envoy_cc_test(
    name = "protocol_annotator_test",
    srcs = ["protocol_annotator_test.cc"],
    repository = "@envoy",
    deps = [
        "//source/extensions/filters/http/protocol_annotator",
        "@envoy//test/mocks/http:http_mocks",
    ],
)

envoy_cc_test(
    name = "config_test",
    srcs = ["config_test.cc"],
    repository = "@envoy",
    deps = [
        "//source/extensions/filters/http/protocol_annotator:config",
        "@envoy//test/mocks/server:server_mocks",
    ],
)

envoy_cc_test(
    name = "protocol_annotator_integration_test",
    srcs = ["protocol_annotator_integration_test.cc"],
    repository = "@envoy",
    deps = [
        "//source/extensions/filters/http/protocol_annotator:config",
        "@envoy//test/integration:http_integration_lib",
    ],
)
