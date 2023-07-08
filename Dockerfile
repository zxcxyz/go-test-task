# syntax=docker/dockerfile:1
####https://docs.docker.com/language/golang/build-images/
FROM golang:1.20

#envs
ENV SERVER_PORT=8080
ENV METRICS_ENDPOINT=/metrics
# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /appserv

# Run
CMD ["/appserv"]
