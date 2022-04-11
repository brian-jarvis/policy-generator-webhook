FROM registry.redhat.io/rhel8/go-toolset:1.16 as gobuilder

WORKDIR /policy-generator-webhook

COPY . .

RUN go build -mod readonly -o work-bin/policy-generator-webhook ./


FROM registry.access.redhat.com/ubi8/ubi-micro:8.5

WORKDIR /root

COPY --from=gobuilder /policy-generator-webhook/work-bin/policy-generator-webhook /root/
