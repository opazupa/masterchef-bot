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
  jwtSecret: string;
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
  databaseConnection: process.env.DATABASE_CONNECTION as string,
  databaseName: process.env.DATABASE_NAME as string,
  enablePlayground: process.env.ENABLE_PLAYGROUND === 'true',
  jwtSecret: process.env.JWT_SECRET as string
};
