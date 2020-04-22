import compression from 'compression';
import cors from 'cors';
import * as dotenv from 'dotenv';
import express from 'express';
import helmet from 'helmet';
import { createServer } from 'http';

import { ConfigureMongoDB } from './database';

// Map .env
dotenv.config();

const port = process.env.PORT || 3000;

// Configure express server
const app = express();
app.use(helmet());
app.disable('x-powered-by');
app.use('*', cors());
app.use(compression());

app.get('/', (_req, res) => res.send('Hello World!'));

// Setup DB and server
ConfigureMongoDB();
const httpServer = createServer(app);

httpServer.on('error', (e) => console.error(e));

httpServer.listen({ port: port }, () => {
  console.log(`🚀 Test api is running on port ${port}`);
});
