
   
ARG ARCH="arm64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="John Orthoefer <john@Orthoefer.org>"

ARG ARCH="arm64"
ARG OS="linux"
COPY pws_exporter  /bin/pws_exporter

EXPOSE      9874
ENTRYPOINT  [ "/bin/pws_exporter" ]
CMD         [ "--version" ]
