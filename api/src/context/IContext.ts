import { IBatchLoaders } from '../dataloaders';

/**
 * Interface for server request context
 *
 * @export
 * @interface IContext
 */
export interface IContext {
  loaders: IBatchLoaders;
}
