FROM golang:1.18.3-alpine3.16 AS build

ARG VERSION=latest
ARG OUTPUT=/iam-pump

WORKDIR /src
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct \
      && go mod tidy -compat=1.18 \
      && go build -ldflags "-X main.Version=${VERSION}" -o ${OUTPUT}/ ./... \
      && cp configs/iam-pump.yaml ${OUTPUT}/ \
      && rm -rf /src

# ================================

FROM alpine:3.16

ENV TZ Asia/Shanghai

RUN apk add tzdata && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
      && echo ${TZ} > /etc/timezone \
      && apk del tzdata

COPY --from=build /iam-pump/iam-pump /opt/iam/bin/
COPY --from=build /iam-pump/iam-pump.yaml /etc/iam/

EXPOSE 8010
ENTRYPOINT [ "/opt/iam/bin/iam-pump" ]
CMD [ "-c", "/etc/iam/iam-pump.yaml" ]
