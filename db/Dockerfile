FROM arigaio/atlas:latest-alpine

COPY migrations /migrations

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["atlas migrate apply --dir file://migrations --url $DATABASE_URL"]