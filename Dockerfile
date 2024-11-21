# 打包依赖阶段使用golang作为基础镜像
FROM golang:1.23 as builder

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY  go.mod .
COPY  go.sum .
RUN go mod download

COPY . .

# 指定OS等，并go build
RUN CGO_ENABLED=0 go build -o code_challenge1  ./

FROM debian:stable-slim

WORKDIR /app
RUN apt-get update && apt-get -y install ca-certificates tzdata gosu tini locales && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \ 
    && echo "${TZ}" > /etc/timezone && apt-get clean && sed -i '/en_US.UTF-8/s/^# //g' /etc/locale.gen && locale-gen

ENV TZ="Asia/Shanghai" GIN_MODE=release PUID=0 PGID=0 UMASK=0
ENV LANG=en_US.UTF-8 LANGUAGE=en_US:en LC_ALL=en_US.UTF-8

# 将上一个阶段publish文件夹下的所有文件复制进来
COPY --from=builder /app/code_challenge1 .
COPY --from=builder /app/entrypoint.sh .
RUN chmod +x /app/entrypoint.sh

VOLUME /app/data
EXPOSE 8080

ENTRYPOINT ["tini","./entrypoint.sh"]