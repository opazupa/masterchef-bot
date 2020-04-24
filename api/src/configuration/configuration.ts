import * as dotenv from 'dotenv';

/**
 * Env configuration interface
 *
 * @interface IConfiguration
 */
interface IConfiguration {
  port: string | number;
  debugMode: boolean;
  databaseConnection: string;
  databaseName: string;
  enablePlayground: boolean;
}

// Apply .env
dotenv.config();

/**
 * Env configuration
 *
 * @interface IConfiguration
 */
export const configuration: IConfiguration = {
  port: process.env.PORT || 3000,
  debugMode: process.env.DEBUG_MODE === 'true',
  databaseConnection: <string>process.env.DATABASE_CONNECTION,
  databaseName: <string>process.env.DATABASE_NAME,
  enablePlayground: process.env.ENABLE_PLAYGROUND === 'true'
};
