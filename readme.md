# Projeto

## Event Driven - NATS

## Migrations c/ Goose

## Camada de dbmodel com sqlc


## Atualizar camada de infra

```bash
    sqlc generate
```

## Migrations

### Integração entre sqlc e migrations

Da documentação
> sqlc does not perform database migrations for you. However, sqlc is able to differentiate between up and down migrations. sqlc ignores down migrations when parsing SQL files.

> sqlc supports parsing migrations from the following tools:
>
> * atlas
> * dbmate
> * golang-migrate
> * goose
> * sql-migrate
> * tern

Só temos que colocar no sqlc.yml o schema dentro de uma pasta de migrations

```yml
schema: "db/migrations"
```

### Configurar o goose

com env:
```bash
GOOSE_DRIVER=DRIVER
GOOSE_DBSTRING=DBSTRING
GOOSE_MIGRATION_DIR=MIGRATION_DIR
```

exemplo:

```bash
❯ export GOOSE_DRIVER=postgres
❯ export GOOSE_DBSTRING="host=localhost port=5432 user=admin@pgadmin.com password=admin dbname=postgres sslmode=disable"
❯ GOOSE_MIGRATION_DIR=./migrations
```

### Aplicar migrations

```bash
goose up
```

```bash
$ goose up
$ OK    001_basics.sql
$ OK    002_next.sql
$ OK    003_and_again.go
```

### Criar migrations

```bash
$ goose create add_some_column sql
$ Created new file: 20170506082420_add_some_column.sql

$ goose -s create add_some_column sql
$ Created new file: 00001_add_some_column.sql
```

### Resto dos comandos
https://github.com/pressly/goose?tab=readme-ov-file#up-to

## 🧾 To-dos

* [X] Integration Events rename
* [X] Conectar com banco
* [-] Use natsConfig
* [ ] Jetstream
* [ ] Inicialização e Atribuição melhor dos módulos
* [X] Decidir entre switch ou parse de evento para eventos
* [ ] Separar melhor o mapping de tipos internos para tipos externos de IntegrationEvents
* [X] Inicialização do módulo melhorada, principalmente dos subscribers
* [ ] Add otel & logging
* [x] Introduce sqlc
