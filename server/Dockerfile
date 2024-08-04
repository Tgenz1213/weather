FROM golang:1.18

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd/weather

RUN go build -o /weather_app

RUN chmod +x /weather_app

EXPOSE 8080

CMD [ "/weather_app" ]