package(default_visibility = ["//visibility:public"])

load(
    "@envoy//bazel:envoy_build_system.bzl",
    "envoy_cc_binary",
)

envoy_cc_binary(
    name = "envoy-static",
    repository = "@envoy",
    deps = [
        "//source/extensions/filters/listener/text_protocol_inspector:config",
        "//source/extensions/filters/http/protocol_annotator:config",
        "@envoy//source/exe:envoy_main_entry_lib",
    ],
)
