import { ApiRole } from '../database/models';
import { IBatchLoaders } from '../dataloaders';

export interface IContextUser {
  userName: string;
  roles: ApiRole[];
}

/**
 * Interface for server request context
 *
 * @export
 * @interface IContext
 */
export interface IContext {
  loaders: IBatchLoaders;
  user: IContextUser;
}
