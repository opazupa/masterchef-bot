# Masterchef Telegram Bot and Recipe Collection :fire:

This is the repository for the **Masterchef bot** features.

## Features

| Component | *path* |
| --- | --- |
| Telegram bot                      | [`/bot`](./bot/)                                 |
| Express GraphQL API               | [`/api`](./api/)                                 |   
| Mongo DB                          | [`/mongo_db`](./mongo_db/)                       |
| [Sentry.io](https://sentry.io)    | Configured for [`bot`](./bot/) & [`api`](./api/) |

## Development 🚨

### Start up :rocket:

Add configurations with `.env` file by taking a look from [`.example.env`](./.example.env) and from
- [api configurations](./api/README.md##Configuration)
- [bot configurations](./bot/README.md##Configuration)

### Running with Make 😽

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

