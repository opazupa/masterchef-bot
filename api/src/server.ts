import 'reflect-metadata';

import { ApolloServer } from 'apollo-server-express';
import compression from 'compression';
import cors from 'cors';
import express from 'express';
import depthLimit from 'graphql-depth-limit';
import helmet from 'helmet';
import { createServer } from 'http';

import { configuration } from './configuration';
import { configureMongoDB } from './database';
import { createSchema } from './graphql';

const bootstrap = async () => {
  // Configure express server
  const app = express();
  app.use(helmet());
  app.disable('x-powered-by');
  app.use('*', cors());
  app.use(compression());

  // Setup GraphQL server
  const server = new ApolloServer({
    schema: await createSchema,
    validationRules: [depthLimit(7)],
    introspection: configuration.enablePlayground,
    playground: configuration.enablePlayground,
    subscriptions: {
      path: '/subscriptions'
    }
  });
  server.applyMiddleware({ app, path: '/graphql' });

  // Setup DB
  configureMongoDB();

  const httpServer = createServer(app);
  server.installSubscriptionHandlers(httpServer);
  httpServer.on('error', (e) => console.error(e));
  httpServer.listen({ port: configuration.port }, () => {
    console.log(`🚀 Test api is running on port ${configuration.port}`);
  });
};

bootstrap().catch((err) => console.log(err));
