# We link in all the transitive dependencies to ensure bazel puts them all on
# the rpath so we can run tests on machines without all the dependencies
# installed on the system.
#
# Things we deliberately get from the host system or sysroot:
#  * glibc: libc6, libm, librt, libdl, libpthread, libresolv
#  * libstdc++
#  * libgcc_s (it's tied to glibc and libstdc++)
#  * things that integrate tightly with glibc: libnsl, libtirpc,
#      libgssapi_krb5, libkrb5, libk5crypto, libkrb5support
#  * libssl (it has lots of conf files that need to match): libcrypto, libssl
#  * libselinux (it needs to be compatible with the kernel)
#  * libunwind (it's used by libstdc++)
#
# We want to keep our versions of opencv on all platforms exactly matching to
# minimize the maintenance burden. We use the same toolchain with similar rootfs
# on all of them, so if you find yourself tempted to make them different
# consider carefully whether the problem you're solving is really unique to that
# platform or will apply to all of them.

_common_srcs_list = [
    "usr/lib/%s/libopencv_core.so.4.5",
    "usr/lib/%s/libopencv_features2d.so.4.5",
    "usr/lib/%s/libopencv_imgproc.so.4.5",
    "usr/lib/%s/libopencv_flann.so.4.5",
    "usr/lib/%s/libopencv_highgui.so.4.5",
    "usr/lib/%s/libopencv_videoio.so.4.5",
    "usr/lib/%s/libopencv_aruco.so.4.5",
    "usr/lib/%s/libopencv_imgcodecs.so.4.5",
    "usr/lib/%s/libopencv_ml.so.4.5",
    "usr/lib/%s/libopencv_calib3d.so.4.5",
    "usr/lib/%s/libtbb.so.2",
    "usr/lib/%s/libgtk-3.so.0",
    "usr/lib/%s/libgdk-3.so.0",
    "usr/lib/%s/libpangocairo-1.0.so.0",
    "usr/lib/%s/libpango-1.0.so.0",
    "usr/lib/%s/libatk-1.0.so.0",
    "usr/lib/%s/libgdcmDICT.so.3.0",
    "usr/lib/%s/libgdcmCommon.so.3.0",
    "usr/lib/%s/libgdcmIOD.so.3.0",
    "usr/lib/%s/libgdcmMSFF.so.3.0",
    "usr/lib/%s/libavutil.so.56",
    "usr/lib/%s/libswscale.so.5",
    "usr/lib/%s/libcairo-gobject.so.2",
    "usr/lib/%s/libcairo.so.2",
    "usr/lib/%s/libgdk_pixbuf-2.0.so.0",
    "usr/lib/%s/libgio-2.0.so.0",
    "usr/lib/%s/libgobject-2.0.so.0",
    "usr/lib/%s/libglib-2.0.so.0",
    "usr/lib/%s/libgthread-2.0.so.0",
    "usr/lib/%s/libdc1394.so.25",
    "usr/lib/%s/libgphoto2.so.6",
    "usr/lib/%s/libgphoto2_port.so.12",
    "usr/lib/%s/libavcodec.so.58",
    "usr/lib/%s/libavformat.so.58",
    "usr/lib/%s/libjpeg.so.62",
    "usr/lib/%s/libwebp.so.6",
    "usr/lib/%s/libpng16.so.16",
    "usr/lib/%s/libtiff.so.5",
    "usr/lib/%s/libImath-2_5.so.25",
    "usr/lib/%s/libIlmImf-2_5.so.25",
    "usr/lib/%s/libIex-2_5.so.25",
    "usr/lib/%s/libHalf-2_5.so.25",
    "usr/lib/%s/libIlmThread-2_5.so.25",
    "usr/lib/libgdal.so.28",
    "usr/lib/%s/libgdcmDSED.so.3.0",
    "usr/lib/%s/libgmodule-2.0.so.0",
    "usr/lib/%s/libX11.so.6",
    "usr/lib/%s/libXi.so.6",
    "usr/lib/%s/libXcomposite.so.1",
    "usr/lib/%s/libXdamage.so.1",
    "usr/lib/%s/libXfixes.so.3",
    "usr/lib/%s/libatk-bridge-2.0.so.0",
    "usr/lib/%s/libxkbcommon.so.0",
    "usr/lib/%s/libwayland-cursor.so.0",
    "usr/lib/%s/libwayland-egl.so.1",
    "usr/lib/%s/libwayland-client.so.0",
    "usr/lib/%s/libepoxy.so.0",
    "usr/lib/%s/libharfbuzz.so.0",
    "usr/lib/%s/libpangoft2-1.0.so.0",
    "usr/lib/%s/libfontconfig.so.1",
    "usr/lib/%s/libfreetype.so.6",
    "usr/lib/%s/libXinerama.so.1",
    "usr/lib/%s/libXrandr.so.2",
    "usr/lib/%s/libXcursor.so.1",
    "usr/lib/%s/libXext.so.6",
    "usr/lib/%s/libthai.so.0",
    "usr/lib/%s/libfribidi.so.0",
    "lib/%s/libexpat.so.1",
    "usr/lib/%s/libgdcmjpeg8.so.3.0",
    "usr/lib/%s/libgdcmjpeg12.so.3.0",
    "usr/lib/%s/libgdcmjpeg16.so.3.0",
    "usr/lib/%s/libopenjp2.so.7",
    "usr/lib/%s/libCharLS.so.2",
    "usr/lib/%s/libuuid.so.1",
    "usr/lib/%s/libjson-c.so.5",
    "usr/lib/%s/libva-drm.so.2",
    "usr/lib/%s/libva.so.2",
    "usr/lib/%s/libva-x11.so.2",
    "usr/lib/%s/libvdpau.so.1",
    "usr/lib/%s/libdrm.so.2",
    "usr/lib/%s/libpixman-1.so.0",
    "usr/lib/%s/libxcb-shm.so.0",
    "usr/lib/%s/libxcb.so.1",
    "usr/lib/%s/libxcb-render.so.0",
    "usr/lib/%s/libXrender.so.1",
    "usr/lib/%s/libffi.so.7",
    "usr/lib/%s/libraw1394.so.11",
    "usr/lib/%s/libusb-1.0.so.0",
    "usr/lib/%s/libltdl.so.7",
    "usr/lib/%s/libexif.so.12",
    "usr/lib/%s/libswresample.so.3",
    "usr/lib/%s/libvpx.so.6",
    "usr/lib/%s/libwebpmux.so.3",
    "usr/lib/%s/librsvg-2.so.2",
    "usr/lib/%s/libzvbi.so.0",
    "usr/lib/%s/libsnappy.so.1",
    "usr/lib/%s/libaom.so.0",
    "usr/lib/%s/libcodec2.so.0.9",
    "usr/lib/%s/libgsm.so.1",
    "usr/lib/%s/libmp3lame.so.0",
    "usr/lib/%s/libopus.so.0",
    "usr/lib/%s/libshine.so.3",
    "usr/lib/%s/libspeex.so.1",
    "usr/lib/%s/libtheoraenc.so.1",
    "usr/lib/%s/libtheoradec.so.1",
    "usr/lib/%s/libtwolame.so.0",
    "usr/lib/%s/libvorbis.so.0",
    "usr/lib/%s/libvorbisenc.so.2",
    "usr/lib/%s/libwavpack.so.1",
    "usr/lib/%s/libx264.so.160",
    "usr/lib/%s/libx265.so.192",
    "usr/lib/%s/libxvidcore.so.4",
    "usr/lib/%s/libxml2.so.2",
    "usr/lib/%s/libgme.so.0",
    "usr/lib/%s/libopenmpt.so.0",
    "usr/lib/%s/libchromaprint.so.1",
    "usr/lib/%s/libbluray.so.2",
    "usr/lib/%s/libgnutls.so.30",
    "usr/lib/%s/libssh-gcrypt.so.4",
    "usr/lib/%s/libzstd.so.1",
    "usr/lib/%s/libjbig.so.0",
    "usr/lib/libarmadillo.so.10",
    "usr/lib/%s/libproj.so.19",
    "usr/lib/%s/libpoppler.so.102",
    "usr/lib/%s/libfreexl.so.1",
    "usr/lib/%s/libqhull.so.8.0",
    "usr/lib/%s/libgeos_c.so.1",
    "usr/lib/%s/libepsilon.so.1",
    "usr/lib/%s/libodbc.so.2",
    "usr/lib/%s/libodbcinst.so.2",
    "usr/lib/%s/libkmlbase.so.1",
    "usr/lib/%s/libkmldom.so.1",
    "usr/lib/%s/libkmlengine.so.1",
    "usr/lib/%s/libxerces-c-3.2.so",
    "usr/lib/%s/libnetcdf.so.18",
    "usr/lib/%s/libhdf5_serial_hl.so.100",
    "usr/lib/%s/libsz.so.2",
    "usr/lib/%s/libhdf5_serial.so.103",
    "usr/lib/libmfhdfalt.so.0",
    "usr/lib/libdfalt.so.0",
    "usr/lib/libogdi.so.4.1",
    "usr/lib/%s/libgif.so.7",
    "usr/lib/%s/libgeotiff.so.5",
    "usr/lib/%s/libpq.so.5",
    "usr/lib/%s/libdapclient.so.6",
    "usr/lib/%s/libdap.so.27",
    "usr/lib/%s/libspatialite.so.7",
    "usr/lib/%s/libcurl-gnutls.so.4",
    "usr/lib/%s/libfyba.so.0",
    "usr/lib/%s/libfygm.so.0",
    "usr/lib/%s/libfyut.so.0",
    "usr/lib/%s/libmariadb.so.3",
    "lib/%s/libdbus-1.so.3",
    "usr/lib/%s/libatspi.so.0",
    "usr/lib/%s/libgraphite2.so.3",
    "usr/lib/%s/libdatrie.so.1",
    "usr/lib/%s/libXau.so.6",
    "usr/lib/%s/libXdmcp.so.6",
    "usr/lib/%s/libblkid.so.1",
    "usr/lib/%s/libsoxr.so.0",
    "usr/lib/%s/libogg.so.0",
    "usr/lib/%s/libicui18n.so.67",
    "usr/lib/%s/libicuuc.so.67",
    "usr/lib/%s/libicudata.so.67",
    "usr/lib/%s/libmpg123.so.0",
    "usr/lib/%s/libvorbisfile.so.3",
    "usr/lib/%s/libp11-kit.so.0",
    "usr/lib/%s/libidn2.so.0",
    "usr/lib/%s/libunistring.so.2",
    "usr/lib/%s/libtasn1.so.6",
    "usr/lib/%s/libnettle.so.8",
    "usr/lib/%s/libhogweed.so.6",
    "usr/lib/%s/libgmp.so.10",
    "usr/lib/%s/libgcrypt.so.20",
    "usr/lib/%s/blas/libblas.so.3",
    "usr/lib/%s/lapack/liblapack.so.3",
    "usr/lib/%s/libarpack.so.2",
    "usr/lib/%s/libsuperlu.so.5",
    "usr/lib/%s/libnss3.so",
    "usr/lib/%s/libsmime3.so",
    "usr/lib/%s/libnspr4.so",
    "usr/lib/%s/liblcms2.so.2",
    "usr/lib/%s/libgeos-3.9.0.so",
    "usr/lib/%s/libminizip.so.1",
    "usr/lib/%s/liburiparser.so.1",
    "usr/lib/%s/libaec.so.0",
    "usr/lib/%s/libssl3.so",
    "usr/lib/%s/libldap_r-2.4.so.2",
    "usr/lib/%s/libsqlite3.so.0",
    "usr/lib/%s/libnghttp2.so.14",
    "usr/lib/%s/librtmp.so.1",
    "usr/lib/%s/libssh2.so.1",
    "usr/lib/%s/libpsl.so.5",
    "usr/lib/%s/liblber-2.4.so.2",
    "usr/lib/%s/libsystemd.so.0",
    "usr/lib/%s/libbsd.so.0",
    "lib/%s/libgpg-error.so.0",
    "usr/lib/%s/libgfortran.so.5",
    "usr/lib/%s/libnssutil3.so",
    "usr/lib/%s/libplc4.so",
    "usr/lib/%s/libplds4.so",
    "usr/lib/%s/libsasl2.so.2",
    "usr/lib/%s/liblz4.so.1",
    "lib/%s/libpcre.so.3",
    "usr/lib/%s/libgomp.so.1",
    "usr/lib/%s/libcharls.so.2",
    "usr/lib/%s/libcfitsio.so.9",
    "usr/lib/%s/librttopo.so.1",
    "usr/lib/%s/libgstbase-1.0.so.0",
    "usr/lib/%s/libgstreamer-1.0.so.0",
    "usr/lib/%s/libgstapp-1.0.so.0",
    "usr/lib/%s/libgstriff-1.0.so.0",
    "usr/lib/%s/libgstpbutils-1.0.so.0",
    "usr/lib/%s/libOpenCL.so.1",
    "usr/lib/%s/libmount.so.1",
    "usr/lib/%s/libdav1d.so.4",
    "usr/lib/%s/librabbitmq.so.4",
    "usr/lib/%s/libdeflate.so.0",
    "usr/lib/%s/libheif.so.1",
    "usr/lib/%s/libbrotlidec.so.1",
    "usr/lib/%s/libzmq.so.5",
    "usr/lib/%s/libsrt-gnutls.so.1.4",
    "usr/lib/%s/libudev.so.1",
    "usr/lib/%s/libudfread.so.0",
    "usr/lib/%s/libmd.so.0",
    "usr/lib/%s/libdw.so.1",
    "usr/lib/%s/libgstaudio-1.0.so.0",
    "usr/lib/%s/libgsttag-1.0.so.0",
    "usr/lib/%s/libgstvideo-1.0.so.0",
    "usr/lib/%s/libde265.so.0",
    "usr/lib/%s/libbrotlicommon.so.1",
    "usr/lib/%s/libsodium.so.23",
    "usr/lib/%s/libpgm-5.3.so.0",
    "usr/lib/%s/libnorm.so.1",
    "usr/lib/%s/libelf.so.1",
    "usr/lib/%s/liborc-0.4.so.0",
]

cc_library(
    name = "opencv",
    srcs = select({
        "@platforms//cpu:x86_64": [s % "x86_64-linux-gnu" if "%" in s else s for s in _common_srcs_list] + [
            "usr/lib/x86_64-linux-gnu/libmfx.so.1",
            "usr/lib/x86_64-linux-gnu/libquadmath.so.0",
            "usr/lib/x86_64-linux-gnu/libnuma.so.1",
        ],
        "@platforms//cpu:armv7": [s % "arm-linux-gnueabihf" if "%" in s else s for s in _common_srcs_list],
        "@platforms//cpu:arm64": [s % "aarch64-linux-gnu" if "%" in s else s for s in _common_srcs_list],
    }),
    hdrs = glob([
        "usr/include/opencv4/**",
    ]),
    includes = [
        "usr/include",
        "usr/include/opencv4",
    ],
    linkopts = [
        "-ldl",
        "-lresolv",
    ],
    target_compatible_with = select({
        "@platforms//cpu:x86_64": [
            "@platforms//os:linux",
        ],
        "@platforms//cpu:armv7": [
            "@platforms//os:linux",
        ],
        "@platforms//cpu:arm64": [
            "@platforms//os:linux",
        ],
    }),
    visibility = ["//visibility:public"],
)
