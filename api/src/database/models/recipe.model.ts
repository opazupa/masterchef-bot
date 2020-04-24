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

export { IRecipe, addRecipe, getAllRecipes, getRecipe, getUserRecipes, getFavouriters, updateRecipe, deleteRecipe };
