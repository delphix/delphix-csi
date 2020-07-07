FROM alpine
LABEL maintainers="Daniel Stolf"
LABEL description="Delphix Kubernetes Driver"

# Add util-linux to get a new version of losetup.
RUN apk add util-linux nfs-utils rpcbind openrc
RUN rc-update add nfsmount
COPY ./bin/delphixplugin /delphixplugin
ENTRYPOINT ["/delphixplugin"]
