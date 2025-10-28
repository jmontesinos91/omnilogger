FROM golang:1.23 as builder

RUN go clean -modcache
WORKDIR /app

RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Copy the go module and sum files
COPY go.mod go.sum ./

# ask for the argument and set gitlab credentials that came frome the docker compose and the .env file
ARG GITLAB_TOKEN
RUN git config --global url."https://gitlab-ci-token:${GITLAB_TOKEN}@git.omnicloud.mx/".insteadOf "https://git.omnicloud.mx/"
RUN echo "machine git.omnicloud.mx login gitlab-ci-token password ${GITLAB_TOKEN}" > ~/.netrc

ENV GOPRIVATE=git.omnicloud.mx/omnicloud/development/go-modules/*

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main cmd/main.go

FROM debian:bullseye-slim

ARG DATABASE_DSN
ARG OMNIVIEW_SERVER

ENV DATABASE_DSN=$DATABASE_DSN
ENV OMNIVIEW_SERVER=$OMNIVIEW_SERVER

#update aptget and install pgtools to use later in the entrypoint
RUN apt-get update && apt-get install -y postgresql-client && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

COPY --from=builder /app/main main
COPY --from=builder /app/resources/config.yml resources/config.yml

RUN chmod +x /app/main

ENTRYPOINT ["/app/main"]