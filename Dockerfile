FROM golang

WORKDIR /app

COPY . /app

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 3000

CMD ["go", "run", "main.go"]