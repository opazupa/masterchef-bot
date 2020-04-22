import mongoose, { ConnectionOptions } from 'mongoose';

const MongoConfig = {
  URI: <string>process.env.DATABASE_CONNECTION,
  OPTIONS: <ConnectionOptions>{
    dbName: process.env.DATABASE_NAME,
    useNewUrlParser: true
  }
};

/**
 * Configure and connect to mongo DB
 *
 * @returns
 */
export const ConfigureMongoDB = () => {
  mongoose.connect(MongoConfig.URI, MongoConfig.OPTIONS, (err) => {
    console.log('Mongodb connected', err ? `with 💥 💥 💥 : ${err}` : 'successfully 👍.');
  });
};
