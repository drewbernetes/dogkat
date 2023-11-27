# syntax=docker/dockerfile:1
FROM cgr.dev/chainguard/go:latest as go-builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd/
COPY pkg ./pkg/

RUN go mod download

RUN CGO_ENABLED=0 go build -o /tmp/dogkat /app/cmd/dogkat/main.go


FROM cgr.dev/chainguard/wolfi-base:latest

RUN apk update --no-cache
RUN apk add bash build-base jq curl git python3 py3-pip openssh-client

RUN echo "ServerAliveInterval 120" >> /etc/ssh/ssh_config

RUN echo "dogkat:x:1000:1000:Dogkat Non Root,,,:/home/dogkat:" >> /etc/passwd
RUN mkdir -p /home/dogkat/.config/openstack/
RUN chown dogkat: -R /home/dogkat

COPY --from=go-builder /tmp/dogkat /usr/bin/

RUN chmod 755 /usr/bin/dogkat
ENV PATH=${PATH}:/home/dogkat/.local/bin/

USER 1000

ENTRYPOINT ["bash"]
