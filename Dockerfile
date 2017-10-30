FROM vxlabs/glide as builder

WORKDIR $GOPATH/src/github.com/jbonachera/packer-scaleway-plugin
RUN mkdir release
COPY glide* ./
RUN glide install -v
COPY . ./
RUN go test $(glide nv) && \
    go build -buildmode=exe -a -o /bin/packer-builder-scaleway-volumesurrogate ./cmd/scaleway-volumesurrogate

FROM alpine
EXPOSE 1883
ENTRYPOINT ["/usr/bin/server"]
RUN apk -U add ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=builder /bin/packer-builder-scaleway-volumesurrogate /usr/bin/packer-builder-scaleway-volumesurrogate

