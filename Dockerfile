FROM golang:1.20

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make
RUN cp bin/pm-web /usr/local/bin/app && cp bin/pm-cli /usr/local/bin/pm-cli && cp bin/pm-discord /usr/local/bin/pm-discord

CMD ["app"]