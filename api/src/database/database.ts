import mongoose, { ConnectionOptions } from 'mongoose';

import { configuration } from '../configuration';

// Mongo configuration
const mongoConfig = {
  URI: configuration.databaseConnection,
  OPTIONS: {
    dbName: configuration.databaseName,
    useNewUrlParser: true
  } as ConnectionOptions
};

/**
 * Configure and connect to mongo DB
 *
 * @returns
 */
export const configureMongoDB = async () => {
  await mongoose.connect(mongoConfig.URI, mongoConfig.OPTIONS, (err) => {
    // tslint:disable-next-line: no-console
    console.log('Mongodb connected', err ? `with 💥 💥 💥 : ${err}` : 'successfully 👍.');
  });
};
