FROM golang:1.20.3 as BUILDER

WORKDIR /app

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /linebotngrok

FROM golang:1.20.3
COPY --from=BUILDER /linebotngrok /linbot

EXPOSE 8080

CMD ["/linebot"]