import 'reflect-metadata';

import { ApolloServer } from 'apollo-server-express';
import compression from 'compression';
import cors from 'cors';
import express from 'express';
import jwt from 'express-jwt';
import depthLimit from 'graphql-depth-limit';
import helmet from 'helmet';
import { createServer } from 'http';
import path from 'path';
import favicon from 'serve-favicon';

import { getUserFromToken, getUserFromWSParams } from './auth';
import { configuration } from './configuration';
import { IContext } from './context';
import { configureMongoDB } from './database';
import { createBatchLoaders } from './dataloaders';
import { createSchema } from './graphql';

const API_PATH = '/graphql';
const WS_PATH = '/subscriptions';

const bootstrap = async () => {
  // Configure express server
  const app = express();
  app.use(helmet());
  app.disable('x-powered-by');
  app.use('*', cors());
  app.use(compression());

  // Add home page
  app.use(favicon(__dirname + '/public/favicon.ico'));
  app.get('/', (_req, res) => {
    res.sendFile(path.join(__dirname + '/public/index.html'));
  });

  // Apply JWT middleware
  app.use(
    API_PATH,
    jwt({
      secret: configuration.jwtSecret,
      credentialsRequired: false
    })
  );

  // Setup GraphQL server
  const server = new ApolloServer({
    schema: await createSchema,
    validationRules: [depthLimit(7)],
    introspection: configuration.enablePlayground,
    playground: configuration.enablePlayground,
    subscriptions: {
      path: WS_PATH,
      onConnect: (connectionParams) => getUserFromWSParams(connectionParams)
    },
    context: async ({ req, connection }) => {
      return {
        loaders: createBatchLoaders(),
        // Check if it's a websocket connection
        user: connection ? connection.context.user : getUserFromToken(req.user)
      } as IContext;
    }
  });
  server.applyMiddleware({ app, path: API_PATH });

  // Setup DB
  configureMongoDB();
  const httpServer = createServer(app);
  server.installSubscriptionHandlers(httpServer);

  // tslint:disable: no-console
  httpServer.on('error', (e) => console.error(e));
  httpServer.listen({ port: configuration.port }, () => {
    console.log(`🚀 Test api is running on port ${configuration.port}`);
  });
};

bootstrap().catch((err) => console.log(err));
