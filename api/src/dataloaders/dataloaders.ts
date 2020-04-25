import DataLoader from 'dataloader';

import { getUserRecipes, IRecipe } from '../database/models';

/**
 * BatchLoaders in context
 *
 * @export
 * @interface IBatchLoaders
 */
interface IBatchLoaders {
  user: {
    recipes: DataLoader<string, IRecipe[], string>;
  };
}

/**
 * Create barch loader for a function call which
 *
 * @param {readonly} keys
 * @param {*} string
 * @param {*} []
 * @param {(keys: any) => Promise<Map<string, any[]>>} func
 * @returns
 */
const createGroupBatcher = async (keys: readonly string[], func: (keys: any) => Promise<Map<string, any[]>>) => {
  const results = await func(keys);
  return keys.map((_key: string) => results.get(_key) || []);
};

/**
 * Create configured batch loaders
 *
 * @returns {IBatchLoaders}
 */
const createBatchLoaders = (): IBatchLoaders => {
  return {
    user: {
      recipes: new DataLoader((keys: readonly string[]) => createGroupBatcher(keys, getUserRecipes))
    }
  };
};

export { IBatchLoaders, createBatchLoaders };
