FROM golang:1.15-buster as build
ARG ARCH=amd64
ENV CGO_ENABLED=0
WORKDIR /work
RUN mkdir -p /etc/ct-exporter /var/log/ct-exporter \
    && chown nobody:nogroup /etc/ct-exporter /var/log/ct-exporter
COPY . .
RUN make ARCH=${ARCH}

FROM scratch
ARG ARCH=amd64
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build --chown=nobody:nogroup /etc/ct-exporter /etc/ct-exporter
COPY --from=build --chown=nobody:nogroup /var/log/ct-exporter /var/log/ct-exporter
COPY --from=build /tmp/ct-exporter/artifacts/ct-exporter-linux-${ARCH} /ct-exporter

USER 65534:65534
EXPOSE 9809
ENTRYPOINT [ "/ct-exporter" ]
