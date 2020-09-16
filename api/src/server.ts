import 'reflect-metadata';

import { configuration, initApp } from './configuration';
import { configureMongoDB } from './database';

const bootstrap = async () => {
  // Configure express server
  const server = await initApp();
  // Connect DB
  configureMongoDB();

  // tslint:disable: no-console
  server.on('error', (e) => console.error(e));
  server.listen({ port: configuration.port }, () => {
    console.log(`🚀 Api is running on port ${configuration.port}`);
  });
};

bootstrap().catch((err) => console.log(err));
