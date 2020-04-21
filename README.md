# Masterchef Telegram Bot and Recipe Collection :fire:

This is the repository for the **Masterchef bot** features.

## Features

| Component | *path* |
| --- | --- |
| Telegram bot          | [`/bot`](./bot/)            |
| Exrpess GraphQL API   | [`/api`](./api/)            |   
| Mongo DB              | [`/mongo_db`](./mongo_db/)  |

## Development

### Start up :rocket:

Add configurations with `.env` file by taking a look from [`.example.env`](./.example.env) and [configurations](##Configuration).

### Running with Make

`Build`

```
make build
make dev-build
```
Run with `Hot reload`
```
// Loads configuration from .env file
make dev-up
make dev-down
make dev-clean
```
Run in `Production`
```
// Requires env variables apart from .env file
make up
make down
```


## Configuration
| Key | Value |
| --- | --- |
| `API_KEY`             | <your_bot_api_key> from Telegram                   |
| `DEBUG_MODE`          | true/false to enable detailed logging              |
| `DATABASE_CONNECTION` | <your_mongo_connection> to connect to mongo server |
| `DATABASE_NAME`       | <your_db_name> to use in mongo server              | 

