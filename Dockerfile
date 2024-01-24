# Final stage
FROM alpine:latest
WORKDIR /
COPY ./fileserver .
COPY ./uploads ./uploads
EXPOSE 8080
USER 1000:1000
ENTRYPOINT ["/fileserver"]
