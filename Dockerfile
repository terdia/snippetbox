FROM golang:1.17-alpine as builder

WORKDIR /go/src/snippetbox

COPY go.mod go.sum ./
RUN go mod download
RUN go get gotest.tools/gotestsum

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/web -x /go/src/snippetbox/cmd/web

EXPOSE 4000

CMD ["/go/bin/web"]

#FROM alpine:latest
##RUN apk --no-cache add ca-certificates

#WORKDIR /snippetbox

#COPY --from=builder /go/bin/ .

#COPY . .

#CMD ["./web"]


