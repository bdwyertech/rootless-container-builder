#!/busybox/sh
# Magic to Provision the Container
# Brian Dwyer - Intelligent Digital Services

# Workaround for GitLab ENTRYPOINT double execution (issue: 1380)
if [ ! -e '/kaniko/.gitlab-runner.lock' ]; then
	touch /kaniko/.gitlab-runner.lock
	# Docker Configuration Helper Utility
	helper-utility
fi

# Passthrough
exec "$@"
