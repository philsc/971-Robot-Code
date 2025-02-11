load("//tools/build_rules:js.bzl", "rollup_bundle", "ts_project")
load("@com_github_google_flatbuffers//:build_defs.bzl", "flatbuffer_cc_library")
load("@com_github_google_flatbuffers//:typescript.bzl", "flatbuffer_ts_library")
load("//aos:config.bzl", "aos_config")

package(default_visibility = ["//visibility:public"])

cc_binary(
    name = "py_log_reader.so",
    srcs = ["py_log_reader.cc"],
    linkshared = True,
    target_compatible_with = ["@platforms//os:linux"],
    deps = [
        "//aos:configuration",
        "//aos:json_to_flatbuffer",
        "//aos/events:shm_event_loop",
        "//aos/events:simulated_event_loop",
        "//aos/events/logging:log_reader",
        "//third_party/python",
        "@com_github_google_glog//:glog",
    ],
)

py_test(
    name = "log_reader_test",
    srcs = ["log_reader_test.py"],
    data = [
        ":py_log_reader.so",
        "@sample_logfile//file",
    ],
    target_compatible_with = ["@platforms//os:linux"],
    deps = ["//aos:configuration_fbs_python"],
)

ts_project(
    name = "plot_index",
    srcs = ["plot_index.ts"],
    target_compatible_with = ["@platforms//os:linux"],
    deps = [
        "//aos:configuration_ts_fbs",
        "//aos/network/www:demo_plot",
        "//aos/network/www:proxy",
        "//frc971/control_loops/drivetrain:down_estimator_plotter",
        "//frc971/control_loops/drivetrain:drivetrain_plotter",
        "//frc971/control_loops/drivetrain:robot_state_plotter",
        "//frc971/control_loops/drivetrain:spline_plotter",
        "//frc971/wpilib:imu_plotter",
        "//y2020/control_loops/drivetrain:localizer_plotter",
        "//y2020/control_loops/superstructure:accelerator_plotter",
        "//y2020/control_loops/superstructure:finisher_plotter",
        "//y2020/control_loops/superstructure:hood_plotter",
        "//y2020/control_loops/superstructure:turret_plotter",
        "//y2021_bot3/control_loops/superstructure:superstructure_plotter",
        "//y2022/control_loops/superstructure:catapult_plotter",
        "//y2022/control_loops/superstructure:climber_plotter",
        "//y2022/control_loops/superstructure:intake_plotter",
        "//y2022/control_loops/superstructure:superstructure_plotter",
        "//y2022/control_loops/superstructure:turret_plotter",
        "//y2022/localizer:localizer_plotter",
        "//y2022/vision:vision_plotter",
    ],
)

rollup_bundle(
    name = "plot_index_bundle",
    entry_point = "plot_index.ts",
    target_compatible_with = ["@platforms//os:linux"],
    deps = [
        ":plot_index",
    ],
)

genrule(
    name = "copy_css",
    srcs = [
        "//aos/network/www:styles.css",
    ],
    outs = ["styles.css"],
    cmd = "cp $< $@",
)

filegroup(
    name = "plotter_files",
    srcs = [
        "index.html",
        "plot_index_bundle.min.js",
        "styles.css",
    ],
)

sh_binary(
    name = "web_plotter",
    srcs = ["web_plotter.sh"],
    data = [
        ":plotter_files",
        "//aos/network:log_web_proxy_main",
    ],
    target_compatible_with = ["@platforms//os:linux"],
)

sh_binary(
    name = "live_web_plotter_demo",
    srcs = ["live_web_plotter_demo.sh"],
    data = [
        ":plotter_files",
        "//aos/network:web_proxy_main",
        "//aos/network/www:test_config",
    ],
    target_compatible_with = ["@platforms//os:linux"],
)

flatbuffer_cc_library(
    name = "plot_data_fbs",
    srcs = [
        "plot_data.fbs",
    ],
    gen_reflections = 1,
    target_compatible_with = ["@platforms//os:linux"],
)

flatbuffer_ts_library(
    name = "plot_data_ts_fbs",
    srcs = [
        "plot_data.fbs",
    ],
    target_compatible_with = ["@platforms//os:linux"],
)

ts_project(
    name = "plot_data_utils",
    srcs = ["plot_data_utils.ts"],
    visibility = ["//visibility:public"],
    deps = [
        ":plot_data_ts_fbs",
        "//aos:configuration_ts_fbs",
        "//aos/network/www:aos_plotter",
        "//aos/network/www:plotter",
        "//aos/network/www:proxy",
        "@com_github_google_flatbuffers//reflection:reflection_ts_fbs",
        "@com_github_google_flatbuffers//ts:flatbuffers_ts",
    ],
)

aos_config(
    name = "plotter",
    src = "plotter_config.json",
    flatbuffers = [":plot_data_fbs"],
    target_compatible_with = ["@platforms//os:linux"],
    deps = ["//aos/events:aos_config"],
)

cc_library(
    name = "in_process_plotter",
    srcs = ["in_process_plotter.cc"],
    hdrs = ["in_process_plotter.h"],
    data = [
        ":plotter",
        "//frc971/analysis/cpp_plot:cpp_plot_files",
    ],
    deps = [
        ":plot_data_fbs",
        "//aos/events:simulated_event_loop",
        "//aos/network:web_proxy",
    ],
)

cc_binary(
    name = "in_process_plotter_demo",
    srcs = ["in_process_plotter_demo.cc"],
    deps = [
        ":in_process_plotter",
        "//aos:init",
    ],
)

cc_binary(
    name = "log_to_match",
    srcs = [
        "log_to_match.cc",
    ],
    target_compatible_with = ["@platforms//os:linux"],
    deps = [
        "//aos:init",
        "//aos/events:simulated_event_loop",
        "//aos/events/logging:log_reader",
        "//frc971/input:joystick_state_fbs",
    ],
)

cc_binary(
    name = "local_foxglove",
    srcs = ["local_foxglove.cc"],
    data = ["@foxglove_studio"],
    deps = [
        "//aos:init",
        "//aos/network:gen_embedded",
        "//aos/seasocks:seasocks_logger",
        "//third_party/seasocks",
    ],
)
