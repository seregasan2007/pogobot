FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
RUN go build -o main . 
CMD ["/app/main"]