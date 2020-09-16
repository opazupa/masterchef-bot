import * as Sentry from '@sentry/node';
import { ApolloServer } from 'apollo-server-express';
import compression from 'compression';
import cors from 'cors';
import express, { Express } from 'express';
import jwt from 'express-jwt';
import depthLimit from 'graphql-depth-limit';
import helmet from 'helmet';
import { createServer, Server } from 'http';
import path from 'path';
import favicon from 'serve-favicon';

import { getUserFromToken, getUserFromWSParams } from '../auth';
import { API_PATH, WS_PATH } from '../const/api';
import { IContext } from '../context';
import { createBatchLoaders } from '../dataloaders';
import { createSchema, sentryPlugin } from '../graphql';
import { configuration } from './configuration';

/**
 * Init the express app httpserver
 *
 * @returns {Promise<Server>}
 */
export const initApp = async (): Promise<Server> => {
  // Configure express server
  const app = express();
  app.use(helmet());
  app.disable('x-powered-by');
  app.use('*', cors());
  app.use(compression());

  // Add JWT authentication
  addJwtMiddleware(app);
  // Add home page
  addHomePage(app);
  // Add sentry for error tracking
  addSentry();

  // Create apollo server
  const server = await initApolloServer();

  // Attach apollo and create final httpServer
  server.applyMiddleware({ app, path: API_PATH });
  const httpServer = createServer(app);
  // Enable websockets
  server.installSubscriptionHandlers(httpServer);
  return httpServer;
};

/**
 * Add home page
 *
 * @param {Express} app
 */
const addHomePage = (app: Express) => {
  app.use(favicon(__dirname + '/../public/favicon.ico'));
  app.get('/', (_req, res) => {
    res.sendFile(path.join(__dirname + '/../public/index.html'));
  });
};

/**
 * Init sentry for error tracking
 *
 */
const addSentry = () => {
  Sentry.init({
    dsn: configuration.sentryDsn,
    environment: configuration.sentryEnv,
    debug: configuration.debugMode
  });
};

/**
 * Add JWT middleware
 *
 * @param {Express} app
 */
const addJwtMiddleware = (app: Express) => {
  app.use(
    API_PATH,
    jwt({
      secret: configuration.jwtSecret,
      credentialsRequired: false
    })
  );
};

/**
 * Setup GraphQL server
 *
 * @returns {Promise<ApolloServer>}
 */
const initApolloServer = async (): Promise<ApolloServer> => {
  const server = new ApolloServer({
    schema: await createSchema,
    validationRules: [depthLimit(7)],
    introspection: configuration.enablePlayground,
    playground: configuration.enablePlayground,
    plugins: [sentryPlugin],
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
  return server;
};
