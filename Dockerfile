FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
WORKDIR /bin

RUN microdnf install --nodocs tar gcc gzip git bind-utils findutils sudo \
    && curl -Lo ./oc.tar.gz https://mirror.openshift.com/pub/openshift-v4/clients/ocp/stable/openshift-client-linux.tar.gz \
    && tar -xf oc.tar.gz \
    && rm -f oc.tar.gz \
    && curl -Lo ./golang.tar.gz https://go.dev/dl/go1.20.3.linux-amd64.tar.gz \
    && tar -xf golang.tar.gz -C / \
    && rm -f golang.tar.gz \
    && microdnf update \
    && microdnf clean all

ENV GOROOT=/go
ENV GOPATH=/root/go
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH

ENV OCP_API_URL ${OCP_API_URL}
ENV OCP_CRED_USR ${OCP_CRED_USR}
ENV OCP_CRED_PSW ${OCP_CRED_PSW}
ENV OCP_TOKEN ${OCP_TOKEN}

ENV TEST_GROUP ${TEST_GROUP}
ENV TEST_CASE ${TEST_CASE}

ENV OCP_ARCH ${OCP_ARCH}
ENV NIGHTLY ${NIGHTLY}
ENV ROSA ${ROSA}
ENV MUST_GATHER_TAG ${MUST_GATHER_TAG}

COPY . /opt/maistra-test-tool
WORKDIR /opt/maistra-test-tool

RUN go install gotest.tools/gotestsum@latest && go mod download

ENTRYPOINT ["scripts/runtests.sh"]
