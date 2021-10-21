FROM golang:1.17-alpine as helper
WORKDIR /go/src/github.com/bdwyertech/kaniko-gitlab/helper-utility
COPY helper-utility/ .
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -ldflags="-s -w" -trimpath .

FROM gcr.io/kaniko-project/executor:debug

COPY --from=helper /go/src/github.com/bdwyertech/kaniko-gitlab/helper-utility/helper-utility /kaniko/.

ARG BUILD_DATE
ARG VCS_REF

LABEL org.opencontainers.image.title="bdwyertech/kaniko-gitlab" \
      org.opencontainers.image.description="For running Kaniko within a GitLab CI Environment" \
      org.opencontainers.image.authors="Brian Dwyer <bdwyertech@github.com>" \
      org.opencontainers.image.url="https://hub.docker.com/r/bdwyertech/kaniko-gitlab" \
      org.opencontainers.image.source="https://github.com/bdwyertech/docker-kaniko-gitlab.git" \
      org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.created=$BUILD_DATE \
      org.label-schema.name="bdwyertech/kaniko-gitlab" \
      org.label-schema.description="For running Kaniko within a GitLab CI Environment" \
      org.label-schema.url="https://hub.docker.com/r/bdwyertech/kaniko-gitlab" \
      org.label-schema.vcs-url="https://github.com/bdwyertech/docker-kaniko-gitlab.git"\
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.build-date=$BUILD_DATE

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]
