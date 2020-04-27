import 'reflect-metadata';

import { ApolloServer } from 'apollo-server-express';
import compression from 'compression';
import cors from 'cors';
import express from 'express';
import jwt from 'express-jwt';
import depthLimit from 'graphql-depth-limit';
import helmet from 'helmet';
import { createServer } from 'http';
import jsonwebtoken from 'jsonwebtoken';

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
      onConnect: (connectionParams: any) => {
        if (connectionParams.Authorization) {
          return <IContext>{
            user: jsonwebtoken.verify(connectionParams.Authorization, configuration.jwtSecret)
          };
        }
        return null;
      }
    },
    context: async ({ req, connection }) => {
      return <IContext>{
        loaders: createBatchLoaders(),
        user: connection ? connection.context.user : req.user
      };
    }
  });
  server.applyMiddleware({ app, path: API_PATH });

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
