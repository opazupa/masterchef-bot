import DataLoader from 'dataloader';

import { getFavouriteRecipes, getFavouriters, getUserRecipes, getUsers, IRecipe, IUser } from '../database/models';

/**
 * BatchLoaders in context
 *
 * @export
 * @interface IBatchLoaders
 */
interface IBatchLoaders {
  user: {
    favourites: DataLoader<string, IRecipe[], string>;
    recipes: DataLoader<string, IRecipe[], string>;
  };
  recipe: {
    favouriters: DataLoader<string, IUser[], string>;
    user: DataLoader<string, IUser | null, string>;
  };
}

/**
 * Create group batch loader for a function call which
 *
 * @param {readonly string[]} keys
 * @param {(keys: string[]) => Promise<Map<string, any[]>>} func
 * @returns
 */
const createGroupBatcher = async <T>(keys: readonly string[], func: (keys: string[]) => Promise<Map<string, T[]>>) => {
  const results = await func(keys as string[]);
  return keys.map((_key: string) => results.get(_key) || []);
};

/**
 * Create single batch loader for a function call which
 *
 * @param {readonly string[]} keys
 * @param {(keys: string[]) => Promise<Map<string, any>>} func
 * @returns
 */
const createSingleBatcher = async <T>(keys: readonly string[], func: (keys: string[]) => Promise<Map<string, T>>) => {
  const results = await func(keys as string[]);
  return keys.map((_key: string) => results.get(_key) || null);
};

/**
 * Create configured batch loaders
 *
 * @returns {IBatchLoaders}
 */
const createBatchLoaders = (): IBatchLoaders => {
  return {
    user: {
      favourites: new DataLoader((keys: readonly string[]) => createGroupBatcher(keys, getFavouriteRecipes)),
      recipes: new DataLoader((keys: readonly string[]) => createGroupBatcher(keys, getUserRecipes))
    },
    recipe: {
      favouriters: new DataLoader((keys: readonly string[]) => createGroupBatcher(keys, getFavouriters)),
      user: new DataLoader((keys: readonly string[]) => createSingleBatcher(keys, getUsers))
    }
  };
};

export { IBatchLoaders, createBatchLoaders };
