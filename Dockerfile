FROM golang:1.13-alpine
WORKDIR /go/src/github.com/bdwyertech/kaniko-gitlab/helper-utility
COPY helper-utility/ .
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM gcr.io/kaniko-project/executor:debug

COPY --from=0 /go/src/github.com/bdwyertech/kaniko-gitlab/helper-utility/helper-utility /kaniko/.

ARG BUILD_DATE
ARG VCS_REF

LABEL org.opencontainers.image.title="kaniko-gitlab" \
      org.opencontainers.image.authors="bdwyertech@github.com" \
      org.opencontainers.image.source="https://github.com/bdwyertech/docker-kaniko-gitlab.git" \
      org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.created=$BUILD_DATE

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]
