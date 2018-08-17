licenses(["notice"])  # Apache 2

# Text protocol filter

load(
    "@envoy//bazel:envoy_build_system.bzl",
    "envoy_cc_library",
    "envoy_package",
)

envoy_package()

envoy_cc_library(
    name = "text_protocol_inspector_lib",
    srcs = ["text_protocol_inspector.cc"],
    hdrs = ["text_protocol_inspector.h"],
    repository = "@envoy",
    deps = [
        "@envoy//include/envoy/event:dispatcher_interface",
        "@envoy//include/envoy/event:timer_interface",
        "@envoy//include/envoy/network:filter_interface",
        "@envoy//include/envoy/network:listen_socket_interface",
        "@envoy//source/common/api:os_sys_calls_lib",
        "@envoy//source/common/common:assert_lib",
        "@envoy//source/common/common:minimal_logger_lib",
        "@envoy//source/common/common:utility_lib",
        "@envoy//source/extensions/transport_sockets:well_known_names",
    ],
)

envoy_cc_library(
    name = "config",
    srcs = ["config.cc"],
    repository = "@envoy",
    deps = [
        "@envoy//include/envoy/registry",
        "@envoy//include/envoy/server:filter_config_interface",
        "//source/extensions/filters/listener/text_protocol_inspector:text_protocol_inspector_lib",
    ],
)
