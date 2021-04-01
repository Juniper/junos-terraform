FROM golang:3.12-alpine

ENV TERRAFORM_VERSION=0.13.4

VOLUME ["/data"]

WORKDIR /data

ENTRYPOINT [ "/bin/bash" ]

RUN apk update && \
    apk add curl jq python bash ca-certificates git openssl unzip wget bash && \
    cd /tmp && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin && \
    wget https://dl.google.com/dl/cloudsdk/channels/rapid/google-cloud-sdk.zip -O /tmp/google-cloud-sdk.zip && \
    cd /usr/local && unzip /tmp/google-cloud-sdk.zip && \
    google-cloud-sdk/install.sh --usage-reporting=false --path-update=true --bash-completion=true && \
    google-cloud-sdk/bin/gcloud config set --installation component_manager/disable_update_check true && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

ENV PATH = $PATH:/usr/local/google-cloud-sdk/bin/
ENV PATH $PATH:/go/bin/

WORKDIR /jtaf

## Copy project inside the container
COPY *.go /jtaf/
COPY generator/jtaf/generator
COPY terraform_providers /jtaf/terraform_providers
