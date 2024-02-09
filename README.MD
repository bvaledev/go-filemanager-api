# GO File Manager

API para gerenciamento de arquivos em nuvem. Atualmente utiliza S3, e também é flexivel para a utilização de outros storages.

Neste projeto utilizei o AWS S3 GO V2 configurado para os spaces da DigitalOcean

## Funcionalidades

- Upload de multiplos arquivos

## Conceitos aplicados

- GO Routines, 
- Buffered Channels
- WaitGroup
- Testes
- Clean Architecture
- Standard Go Project Layout

## Rodando o aplicativo

Duplique o arquivo env `cp .env.example .env` e configure com os dados corretos.


```
$ go mod tidy

$ go run cmd/api/main.go

```