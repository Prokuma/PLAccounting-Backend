FROM golang:1.21.0-alpine
RUN apk update && apk add git
RUN mkdir -p /go/src/github.com/Prokuma
WORKDIR /go/src/github.com/Prokuma
#RUN git clone https://github.com/Prokuma/PLAccounting-Backend.git
ADD . /go/src/github.com/Prokuma/PLAccounting-Backend
WORKDIR /go/src/github.com/Prokuma/PLAccounting-Backend
#RUN git checkout tags/alpha-v0.1
RUN go build -o main .
EXPOSE 3000
CMD ["/go/src/github.com/Prokuma/PLAccounting-Backend/main"]