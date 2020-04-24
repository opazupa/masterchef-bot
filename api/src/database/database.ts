import mongoose, { ConnectionOptions } from 'mongoose';

import { configuration } from '../configuration';

// Mongo configuration
const mongoConfig = {
  URI: configuration.databaseConnection,
  OPTIONS: <ConnectionOptions>{
    dbName: configuration.databaseName,
    useNewUrlParser: true
  }
};

/**
 * Configure and connect to mongo DB
 *
 * @returns
 */
export const configureMongoDB = async () => {
  await mongoose.connect(mongoConfig.URI, mongoConfig.OPTIONS, (err) => {
    console.log('Mongodb connected', err ? `with 💥 💥 💥 : ${err}` : 'successfully 👍.');
  });
};
