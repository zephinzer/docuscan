FROM golang:1.12.0-alpine3.9 AS build
RUN apk update --no-cache && apk upgrade --no-cache
RUN apk add --no-cache bash curl git g++ jq leptonica poppler-utils tesseract-ocr-dev
# c
RUN wget -q -P /usr/share/tessdata https://github.com/tesseract-ocr/tessdata/raw/master/chi_sim.traineddata
RUN wget -q -P /usr/share/tessdata https://github.com/tesseract-ocr/tessdata/raw/master/chi_tra.traineddata
# m
RUN wget -q -P /usr/share/tessdata https://github.com/tesseract-ocr/tessdata/raw/master/msa.traineddata
# i
RUN wget -q -P /usr/share/tessdata https://github.com/tesseract-ocr/tessdata/raw/master/hin.traineddata
RUN wget -q -P /usr/share/tessdata https://github.com/tesseract-ocr/tessdata/raw/master/tam.traineddata
# oh?
WORKDIR /app
ENV GO111MODULE=on
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download && go mod vendor
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=1

COPY . .

RUN go build -a -o /go/bin/app

ENTRYPOINT ["/go/bin/app"]
