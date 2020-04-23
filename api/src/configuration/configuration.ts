import * as dotenv from 'dotenv';

interface IConfiguration {
  port: string | number;
  debugMode: Boolean;
  databaseConnection: string;
  databaseName: string;
  enablePlayground: Boolean;
}

dotenv.config();

export const configuration: IConfiguration = {
  port: process.env.PORT || 3000,
  debugMode: process.env.DEBUG_MODE === 'true',
  databaseConnection: <string>process.env.DATABASE_CONNECTION,
  databaseName: <string>process.env.DATABASE_NAME,
  enablePlayground: process.env.ENABLE_PLAYGROUND === 'true'
};
