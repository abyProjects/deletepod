FROM golang

RUN apt-get update -y && apt-get install -y curl

WORKDIR /app

COPY . /app/

RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN mv kubectl /usr/local/bin

RUN go mod download

ARG PODNAME
ENV PODNAME nginx

ARG NAMESPACE
ENV NAMESPACE default

ARG TOKEN
ENV TOKEN abcd

RUN go build -o /deletepod

CMD /deletepod pod --name ${PODNAME} --namespace ${NAMESPACE} --token ${TOKEN}