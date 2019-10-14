FROM gcr.io/kaniko-project/executor:debug

ARG BUILD_DATE
ARG VCS_REF

LABEL org.opencontainers.image.title="kaniko-gitlab" \
      org.opencontainers.image.authors="bdwyertech@github.com" \
      org.opencontainers.image.source="https://github.com/bdwyertech/docker-kaniko-gitlab.git" \
      org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.created=$BUILD_DATE

COPY docker-manifest/config.json /kaniko/.docker/config.json

ENTRYPOINT ["/busybox/sh", "-c"]
