FROM golang

ENV GOPATH=/go
ENV PROJECT_REPO=gitlab.com/rockship/payment-gateway
ENV PROJECT_DIR=${GOPATH}/src/${PROJECT_REPO}/
ENV PROJECT_PORT=3000

RUN apt-get update -qq && apt-get install -qq -y curl netcat
RUN curl -s https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR ${PROJECT_DIR}
COPY Gopkg.* ${PROJECT_DIR}
RUN dep ensure --vendor-only -v

COPY . ${PROJECT_DIR}

EXPOSE ${PROJECT_PORT}

RUN go build -o payment-gateway
CMD sh wait-for.sh db_payment:5432 -- ./payment-gateway