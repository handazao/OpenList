# ---------- 构建阶段 ----------
ARG BASE_IMAGE_TAG=base
ARG TARGETOS
ARG TARGETARCH

FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder
WORKDIR /app/

# 安装编译依赖
RUN apk add --no-cache bash curl jq gcc git musl-dev

# 拉取依赖
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源码
COPY . .

# 编译目标架构二进制
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o bin/openlist -ldflags="-s -w" ./cmd/openlist

# ---------- 运行阶段 ----------
FROM openlistteam/openlist-base-image:${BASE_IMAGE_TAG}
WORKDIR /opt/openlist/

ARG USER=openlist
ARG UID=1001
ARG GID=1001

RUN addgroup -g ${GID} ${USER} && \
    adduser -D -u ${UID} -G ${USER} ${USER} && \
    mkdir -p /opt/openlist/data

COPY --from=builder --chmod=755 --chown=${UID}:${GID} /app/bin/openlist ./
COPY --chmod=755 --chown=${UID}:${GID} entrypoint.sh /entrypoint.sh

USER ${USER}
RUN /entrypoint.sh version

ENV UMASK=022 RUN_ARIA2=false
VOLUME /opt/openlist/data/
EXPOSE 5244 5245
CMD [ "/entrypoint.sh" ]
