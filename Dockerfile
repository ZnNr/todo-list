FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY *.db ./
COPY cmd/*.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /todo-list

ENV TODO_PORT="7540"
ENV TODO_DBFILE="./todolist.db"



EXPOSE 7540

CMD ["/todo-list"]