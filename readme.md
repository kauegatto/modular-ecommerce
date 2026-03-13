# Projeto Ecommerce Modular

* Toda a infraestrutura necessária está contida em um docker-compose:
  * Nats.io
  * Postgres
  * PgAdmin
* Requisitos:
  * Docker
  * SQLc
  * Goose
  * Go 1.22+

## Observações

Esse projeto tem como foco aprendizado, estou usando diferentes ferramentas, práticas e padrões que já ouvi falar bem mas nunca tive a oportunidade de pôr em prática em um projeto real em grande escala. Esse projeto tem como foco ser grande, escalável, apesar disso, ainda é focado em conhecimento. Por enquanto só o módulo de order, configuração do nats, db e publicação e recebimento de eventos está desenvolvido.

## Modelos de Domínio, diagramas de Arquitetura

Você pode encontrar diagramas de arquitetura à nível de solução e módulo (por enquanto, só o módulo de order no e-commerce.drawio na raíz do projeto. Lá poderá encontar diagramas de classes, máquinas de estados e alguns outros diagramas representando os objetos de domínio em uma abordagem orientada à domínio (DDD), deixando claro aggregates, entidades e vo's


## Event Driven (Com NATS.io)

[Post no dev.to](https://dev.to/kauegatto/wip-arquiteturas-orientadas-a-eventos-e-monolitos-modulares-3ac2-temp-slug-5623860?preview=5be5a5733061cd124a999f8373fb107687897a9e3b03fb92fc01952737f53d3682f42925b507aa4bb02858c7fc797539b2686e03222e024f271ddb42) sobre arquitetura orientada à eventos, especialmente em sistemas monolíticos modularizados e microserviços:

## Camada de DBModel com SQLc

Para agilizar e promover mais segurança no desenvolvimento, optei por utilizar [sqlc](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html), ferramenta que gera objetos que serão deserializados e serializados no banco de dados, structs de request/response, abstrações para uso de transações, entre outras coisas, precisamos somente prover o estado do schema e as queries

> verificar `src//Order/query.sql`

## Atualizar camada de infra

```bash
sqlc generate
```

## Migrations (Com Goose)

Para cuidar do schema do banco de dados, faremos migrações, que gerenciam o estado do nosso banco de dados em diferentes momentos.

Somos capazes de voltar para pontos prefixados do passado, voltar apenas um ou uma quantidade específica de versões, entre outras vantagens como seeding.

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
❯ export GOOSE_MIGRATION_DIR=./migrations
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

## To-dos

* [X] Integration Events rename
* [X] Conectar com banco
* [X] Use natsConfig
* [ ] Jetstream
* [ ] Inicialização e Atribuição melhor dos módulos
* [X] Decidir entre switch ou parse de evento para eventos
* [X] Separar melhor o mapping de tipos internos para tipos externos de IntegrationEvents
* [X] Inicialização do módulo melhorada, principalmente dos subscribers
* [ ] Add otel & logging
* [x] Introduce sqlc
* [X] Integrate withe eRede!

### Para minha facilidade: PGCli

`pgcli 'postgres://admin%40pgadmin.com:admin@localhost:5432/postgres'`

## Problemas conhecidos

* N+1 Query, diversos lugares, diversas vezes
* NATS não implementa jetstream e persistência
* Falta de retry de eventos em todas as operações
* Falta de implementação do outbox pattern
