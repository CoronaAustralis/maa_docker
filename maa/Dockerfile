FROM alpine:edge as builder

RUN apk add --no-cache wget 

ARG SUPERCRONIC_SHA1SUM_ARM64="512f6736450c56555e01b363144c3c9d23abed4c"
ARG SUPERCRONIC_SHA1SUM_AMD64="cd48d45c4b10f3f0bfdd3a57d054cd05ac96812b"
ARG SUPERCRONIC_VERSION="0.2.29"

ARG TARGETARCH
RUN case ${TARGETARCH} in \
    "amd64") SUPERCRONIC_SHA1SUM=${SUPERCRONIC_SHA1SUM_AMD64} ;; \
    "arm64") SUPERCRONIC_SHA1SUM=${SUPERCRONIC_SHA1SUM_ARM64} ;; \
    esac \
    && wget https://github.com/aptible/supercronic/releases/download/v${SUPERCRONIC_VERSION}/supercronic-linux-${TARGETARCH} -O supercronic \
    && echo "${SUPERCRONIC_SHA1SUM}" supercronic | sha1sum -c - \
    && chmod +x supercronic \
    && mv supercronic /usr/local/bin/supercronic

ARG MAA_VERSION="0.4.5"

RUN case ${TARGETARCH} in \
    "amd64") MAA_ARCH="x86_64" ;; \
    "arm64") MAA_ARCH="aarch64" ;; \
    esac \
    && wget https://github.com/MaaAssistantArknights/maa-cli/releases/download/v${MAA_VERSION}/maa_cli-v${MAA_VERSION}-${MAA_ARCH}-unknown-linux-gnu.tar.gz \
    && tar -xzf maa_cli-v${MAA_VERSION}-${MAA_ARCH}-unknown-linux-gnu.tar.gz \
    && mv maa_cli-${MAA_ARCH}-unknown-linux-gnu/maa /usr/local/bin/maa \
    && chmod +x /usr/local/bin/maa

FROM debian
COPY --from=builder /usr/local/bin/supercronic /usr/local/bin/supercronic
COPY --from=builder /usr/local/bin/maa /usr/local/bin/maa
COPY ./app /source/
COPY ./entrypoint.sh /
ENV TZ=Asia/Shanghai
RUN apt-get update && apt-get install -y ca-certificates tzdata git adb && echo "$TZ" > /etc/timezone \
&& mkdir /app
ENV MAA_CONFIG_DIR=/app
ENV LANG C.UTF-8
ENTRYPOINT [ "/entrypoint.sh" ]