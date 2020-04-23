import { ApolloServer } from 'apollo-server-express';
import compression from 'compression';
import cors from 'cors';
import express from 'express';
import depthLimit from 'graphql-depth-limit';
import helmet from 'helmet';
import { createServer } from 'http';

import { configuration } from './configuration';
import { configureMongoDB } from './database';
import { schema } from './graphql/schema';

// Configure express server
const app = express();
app.use(helmet());
app.disable('x-powered-by');
app.use('*', cors());
app.use(compression());

// Setup GraphQL server
const server = new ApolloServer({
  schema,
  validationRules: [depthLimit(7)],
  introspection: true,
  playground: true
});
server.applyMiddleware({ app, path: '/graphql' });

// Setup DB
configureMongoDB();

const httpServer = createServer(app);
httpServer.on('error', (e) => console.error(e));
httpServer.listen({ port: configuration.port }, () => {
  console.log(`🚀 Test api is running on port ${configuration.port}`);
});
