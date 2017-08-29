FROM scratch

LABEL maintainer="estafette.io" \
      description="The estafette-extension-bitbucket-status component is an Estafette extension to update build status in Bitbucket for builds handled by Estafette CI"

COPY ca-certificates.crt /etc/ssl/certs/
COPY estafette-extension-bitbucket-status /

ENTRYPOINT ["/estafette-extension-bitbucket-status"]