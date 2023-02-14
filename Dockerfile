FROM golang:1.20.1-buster as build
WORKDIR /work
RUN mkdir -p /etc/ct-exporter /var/log/ct-exporter \
    && chown nobody:nogroup /etc/ct-exporter /var/log/ct-exporter


FROM scratch
LABEL org.opencontainers.image.authors="Hsn723" \
      org.opencontainers.image.title="ct-exporter" \
      org.opencontainers.image.source="https://github.com/hsn723/ct-exporter"
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build --chown=nobody:nogroup /etc/ct-exporter /etc/ct-exporter
COPY --from=build --chown=nobody:nogroup /var/log/ct-exporter /var/log/ct-exporter
COPY LICENSE /LICENSE
COPY ct-exporter /

USER 65534:65534
EXPOSE 9809
ENTRYPOINT [ "/ct-exporter" ]
