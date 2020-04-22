import compression from 'compression';
import cors from 'cors';
import * as dotenv from 'dotenv';
import express from 'express';
import helmet from 'helmet';
import { createServer } from 'http';

dotenv.config();

const port = process.env.PORT || 3000;
const app = express();

app.use(helmet());
app.disable('x-powered-by');
app.use('*', cors());
app.use(compression());

app.get('/', (_req, res) => res.send('Hello World!'));

const httpServer = createServer(app);

httpServer.listen({ port: port }, (): void => {
  console.log(`🚀 Test api is running on port ${port}`);
});
