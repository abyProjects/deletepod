FROM golang:alpine

WORKDIR /app

COPY . /app/

RUN go mod download

ARG PODNAME
ENV PODNAME="nginx"

ARG NAMESPACE
ENV NAMESPACE="default"

ARG TOKEN
ENV TOKEN="abcd"

RUN go build -o /deletepod

CMD [ "/deletepod pod --name ${PODNAME} --namespace ${NAMESPACE} --token ${TOKEN} ]