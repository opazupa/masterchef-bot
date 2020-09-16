# Masterchef GraphQL API :atom:

This is the node `GraphQL` API for **Masterchef** recipes.

## Components
- Typescript
- Express
- Apollo Server - GraphQL
    - WS subscriptions
- Mongoose
- JWT bearer auth
- Sentry.io for error tracking

## [Schema](./schema.gql)


## Configuration 🔧
| Key | Value |
| --- | --- |
| `PORT`                | <port_number> for api to run on                           |
| `DATABASE_CONNECTION` | <your_mongo_connection> to connect to mongo server        |
| `DATABASE_NAME`       | <your_db_name> to use in mongo server                     | 
| `ENABLE_PLAYGROUND`   | `true/false` to enable graphql playground in `/grahpql`   |
| `JWT_SECRET`          | <jwt_secret> for authentication                           |
| `SENTRY_DSN`          | <your_dsn> to use with sentry                             | 
| `SENTRY_ENVIRONMENT`  | <your_env> to use with sentry                             |