FROM registry.redhat.io/rhel8/go-toolset:1.16 as gobuilder

WORKDIR /policy-generator-webhook
USER root

COPY . .

RUN mkdir work-bin
RUN go build -mod readonly -o work-bin/policy-generator-webhook ./


FROM registry.access.redhat.com/ubi8/ubi-micro:8.5

WORKDIR /root

COPY --from=gobuilder /policy-generator-webhook/work-bin/policy-generator-webhook .

USER 65324:65324

ENTRYPOINT ["/root/policy-generator-webhook"]