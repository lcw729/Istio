FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/hcp-loadbalancer \
    USER_UID=root \
    USER_NAME=hcp-loadbalancer

COPY _output/bin/hcp-loadbalancer ${OPERATOR}
COPY bin /usr/local/bin

# RUN mkdir /usr/local/bin/user_setup
RUN /usr/local/bin/user_setup
RUN chmod +x /usr/local/bin/entrypoint

ENTRYPOINT [ "/usr/local/bin/entrypoint" ]

USER ${USER_UID}