FROM golang:1.18-alpine AS build

ENV GOOS=darwin
WORKDIR /go/src/foxyshot
COPY . .
RUN CGO_ENABLED=0 go build -o /app/foxyshot

FROM scratch AS export
COPY --from=build /app /