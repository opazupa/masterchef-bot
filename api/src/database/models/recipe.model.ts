import mongoose, { Document, Schema } from 'mongoose';

import { RECIPE, RECIPE_COLLECTION, USER, USER_COLLECTION } from './collections';

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
}

/**
 * Recipe schema
 */
const recipeSchema: Schema = new Schema({
  UserID: { type: Schema.Types.ObjectId, ref: USER, required: true },
  Name: { type: Schema.Types.String, required: true },
  URL: { type: Schema.Types.String, required: true },
  Added: { type: Schema.Types.Date, required: true, default: Date.now }
});

// Recipe mongoose model
const Recipes = mongoose.model<IRecipe>(RECIPE, recipeSchema, RECIPE_COLLECTION);

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
 * Get recipes by user
 *
 * @param {string} userId
 * @returns {Promise<IRecipe[]>}
 */
const getUserRecipes = async (userId: string): Promise<IRecipe[]> => {
  return await Recipes.find({ UserID: userId });
};

/**
 * Get users who have favourited the recipe
 *
 * @param {string} recipeId
 * @returns {Promise<any[]>}
 */
const getFavouriters = async (recipeId: string): Promise<any[]> => {
  return await Recipes.aggregate([
    {
      $match: {
        _id: recipeId
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
      $unwind: {
        path: '$Favourited',
        preserveNullAndEmptyArrays: false
      }
    },
    {
      $project: {
        _id: '$Favourited._id',
        UserName: '$Favourited.UserName',
        Registered: '$Favourited.Registered'
      }
    },
    {
      $sort: {
        Name: -1
      }
    }
  ]);
};

export { IRecipe, getAllRecipes, getRecipe, getUserRecipes, getFavouriters };
