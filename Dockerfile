FROM golang:latest

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o código-fonte do seu projeto para o contêiner
COPY . .

# Compile sua aplicação Go
RUN go build -o main .

# Exponha a porta em que sua aplicação vai rodar
EXPOSE 8080

# Comando para executar sua aplicação
CMD ["./main"]
