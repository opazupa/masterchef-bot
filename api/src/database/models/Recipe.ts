import mongoose, { Document, Schema } from 'mongoose';

import { RECIPE, RECIPE_COLLECTION, USER, USER_COLLECTION } from './collections';
import { IUser } from './user';

/**
 * IRecipe document
 *
 * @export
 * @interface IRecipe
 * @extends {Document}
 */
export interface IRecipe extends Document {
  UserID: string;
  Name: string;
  URL: string;
  Added: Date;
}

const recipeSchema: Schema = new Schema({
  UserID: { type: Schema.Types.ObjectId, ref: USER, required: true },
  Name: { type: Schema.Types.String, required: true },
  URL: { type: Schema.Types.String, required: true },
  Added: { type: Schema.Types.Date, required: true, default: Date.now }
});

const Recipe = mongoose.model<IRecipe>(RECIPE, recipeSchema, RECIPE_COLLECTION);

const getAllRecipes = async (): Promise<IRecipe[]> => {
  return await Recipe.find();
};

const getRecipe = async (id: string): Promise<IRecipe | null> => {
  return await Recipe.findById(id);
};

const getUserRecipes = async (userId: string): Promise<IRecipe[]> => {
  return await Recipe.find({ UserID: userId });
};

const getFavouritedUsers = async (recipeId: string): Promise<IUser[]> => {
  return await Recipe.aggregate([
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

export { getAllRecipes, getRecipe, getUserRecipes, getFavouritedUsers };
