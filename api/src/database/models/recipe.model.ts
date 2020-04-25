import mongoose, { Document, Schema } from 'mongoose';

import { RECIPE, RECIPE_COLLECTION, USER, USER_COLLECTION } from '../collections';
import { IUser } from './user.model';

/**
 * IRecipe document
 *
 * @interface IRecipe
 * @extends {Document}
 */
interface IRecipe extends Document {
  UserID: string;
  Name: string;
  URL: string;
  Added: Date;
  Updated: Date;
}

/**
 * Recipe schema
 */
const recipeSchema: Schema = new Schema(
  {
    UserID: { type: Schema.Types.ObjectId, ref: USER, required: true },
    Name: { type: Schema.Types.String, required: true },
    URL: { type: Schema.Types.String, required: true },
    Added: { type: Schema.Types.Date, required: true, default: Date.now },
    Updated: { type: Schema.Types.Date, required: false, default: null }
  },
  {
    timestamps: {
      createdAt: 'Added',
      updatedAt: 'Updated'
    }
  }
);

// Recipe mongoose model
const Recipes = mongoose.model<IRecipe>(RECIPE, recipeSchema, RECIPE_COLLECTION);

/**
 * Add new recipe
 *
 * @returns {Promise<IRecipe>}
 */
const addRecipe = async (userId: string, name: string, url: string): Promise<IRecipe> => {
  return await Recipes.create({
    UserID: userId,
    Name: name,
    URL: url
  } as IRecipe);
};

/**
 * Get all recipes
 *
 * @returns {Promise<IRecipe[]>}
 */
const getAllRecipes = async (): Promise<IRecipe[]> => {
  return await Recipes.find();
};

/**
 * Get recipe by id
 *
 * @param {string} id
 * @returns {(Promise<IRecipe | null>)}
 */
const getRecipe = async (id: string): Promise<IRecipe | null> => {
  return await Recipes.findById(id);
};

/**
 * Get recipes by ids
 *
 * @param {string[]} ids
 * @returns {(Promise<Map<string, IRecipe | null>>)}
 */
const getRecipes = async (ids: string[]): Promise<Map<string, IRecipe | null>> => {
  const recipes = await Recipes.find({ _id: { $in: ids } });
  return new Map(ids.map((id) => [id, recipes.find((r) => r._id.toString() === id.toString()) || null]));
};

/**
 * Get recipes by user
 *
 * @param {string[]} userIds
 * @returns {Promise<Map<string, IRecipe[]>>}
 */
const getUserRecipes = async (userIds: string[]): Promise<Map<string, IRecipe[]>> => {
  const recipes = await Recipes.find({ UserID: { $in: userIds } });
  return new Map(userIds.map((userId) => [userId, recipes.filter((r) => r.UserID.toString() === userId.toString())]));
};

/**
 * Get users who have favourited the recipe
 *
 * @param {string[]} recipeIds
 * @returns {Promise<Map<string, IUser[]>>}
 */
const getFavouriters = async (recipeIds: string[]): Promise<Map<string, IUser[]>> => {
  const results = (await Recipes.aggregate([
    {
      $match: {
        _id: { $in: recipeIds }
      }
    },
    {
      $lookup: {
        from: USER_COLLECTION,
        localField: '_id',
        foreignField: 'Favourites.RecipeID',
        as: 'Favourited'
      }
    },
    {
      $sort: {
        'Favourited.Name': -1
      }
    }
  ])) as (IRecipe & { Favourited: IUser[] })[];
  return new Map(
    recipeIds.map((recipeId) => [
      recipeId,
      results.filter((r) => r._id.toString() === recipeId.toString()).flatMap((r) => r.Favourited)
    ])
  );
};

/**
 * Update new recipe
 *
 * @returns {Promise<IRecipe | null>}
 */
const updateRecipe = async (id: string, name: string, url: string): Promise<IRecipe | null> => {
  await Recipes.findOneAndUpdate(
    { _id: id },
    {
      Name: name,
      URL: url
    }
  );
  return await Recipes.findById(id);
};

/**
 * Delete the recipe
 *
 * @returns {Promise<void>}
 */
const deleteRecipe = async (id: string): Promise<void | null> => {
  await Recipes.findOneAndDelete({ _id: id });
};

export {
  IRecipe,
  addRecipe,
  getAllRecipes,
  getRecipes,
  getRecipe,
  getUserRecipes,
  getFavouriters,
  updateRecipe,
  deleteRecipe
};
