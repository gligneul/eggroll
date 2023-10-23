# syntax=docker.io/docker/dockerfile:1.4
FROM ubuntu:22.04 as build-stage

RUN <<EOF
apt update
apt install -y --no-install-recommends \
    build-essential=12.9ubuntu3 \
    ca-certificates \
    g++-riscv64-linux-gnu=4:11.2.0--1ubuntu1 \
    wget
EOF

ARG GOVERSION=1.21.1

WORKDIR /opt/build

RUN wget https://go.dev/dl/go${GOVERSION}.linux-$(dpkg --print-architecture).tar.gz && \
    tar -C /usr/local -xzf go${GOVERSION}.linux-$(dpkg --print-architecture).tar.gz

ENV GOOS=linux
ENV GOARCH=riscv64
ENV CGO_ENABLED=1
ENV CC=riscv64-linux-gnu-gcc
ENV PATH=/usr/local/go/bin:${PATH}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY eggeth eggeth
COPY eggtest eggtest
COPY eggtypes eggtypes
COPY internal internal
COPY wallets wallets
COPY client.go .
COPY codec.go .
COPY env.go .
COPY contract.go .

# runtime stage: produces final image that will be executed
FROM --platform=linux/riscv64 riscv64/ubuntu:22.04 as runtime

LABEL io.sunodo.sdk_version=0.2.0
LABEL io.cartesi.rollups.ram_size=128Mi

ARG MACHINE_EMULATOR_TOOLS_VERSION=0.12.0
RUN <<EOF
apt-get update
apt-get install -y --no-install-recommends busybox-static=1:1.30.1-7ubuntu3 ca-certificates=20230311ubuntu0.22.04.1 curl=7.81.0-1ubuntu1.14
curl -fsSL https://github.com/cartesi/machine-emulator-tools/releases/download/v${MACHINE_EMULATOR_TOOLS_VERSION}/machine-emulator-tools-v${MACHINE_EMULATOR_TOOLS_VERSION}.tar.gz \
  | tar -C / --overwrite -xvzf -
rm -rf /var/lib/apt/lists/*
EOF

ENV PATH="/opt/cartesi/bin:${PATH}"

WORKDIR /var/opt/cartesi-dapp
ENTRYPOINT ["rollup-init"]
CMD ["/var/opt/cartesi-dapp/dapp"]

#
# Example images
#

FROM build-stage as honeypot-build-stage
COPY examples/honeypot examples/honeypot
RUN go build ./examples/honeypot/honeypot-contract
FROM --platform=linux/riscv64 runtime as honeypot-contract
COPY --from=honeypot-build-stage /opt/build/honeypot-contract dapp

FROM build-stage as inspect-build-stage
COPY examples/inspect examples/inspect
RUN go build ./examples/inspect/inspect-contract
FROM --platform=linux/riscv64 runtime as inspect-contract
COPY --from=inspect-build-stage /opt/build/inspect-contract dapp

FROM build-stage as template-build-stage
COPY examples/template examples/template
RUN go build ./examples/template/template-contract
FROM --platform=linux/riscv64 runtime as template-contract
COPY --from=template-build-stage /opt/build/template-contract dapp

FROM build-stage as textbox-build-stage
COPY examples/textbox examples/textbox
RUN go build ./examples/textbox/textbox-contract
FROM --platform=linux/riscv64 runtime as textbox-contract
COPY --from=textbox-build-stage /opt/build/textbox-contract dapp
