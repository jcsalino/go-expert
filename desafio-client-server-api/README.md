# Desafio Client-Server API

Desafio do curso **Go Expert**. O `server` busca a cotação USD-BRL na [AwesomeAPI](https://economia.awesomeapi.com.br/json/last/USD-BRL), grava no SQLite e devolve o `bid`. O `client` consome o server e salva a cotação em `cotacao.txt`.

## Requisitos

- Go 1.23+

## Como rodar

1. Suba o server (fica em `http://localhost:8080/cotacao`):

   ```sh
   go run server.go
   ```

2. Em outro terminal, rode o client:

   ```sh
   go run client.go
   ```

## Resultado

- `cotacao.txt` — gerado pelo client no formato `Dollar: <valor>`.
- `cotacao.db` — banco SQLite onde o server grava cada cotação.

## Timeouts

| Operação              | Timeout |
| --------------------- | ------- |
| Chamada à API externa | 200ms   |
| Gravação no banco     | 10ms    |
| Requisição do client  | 300ms   |
