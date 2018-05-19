FROM vxlabs/glide as builder

WORKDIR $GOPATH/src/github.com/jbonachera/packer-scaleway-plugin
COPY glide* ./
RUN glide install -v
COPY . ./
RUN go test $(glide nv) && \
    go build -buildmode=exe -a -o /bin/packer-builder-scaleway-volumesurrogate ./cmd/scaleway-volumesurrogate

FROM alpine
COPY --from=builder /bin/packer-builder-scaleway-volumesurrogate /usr/bin/packer-builder-scaleway-volumesurrogate

