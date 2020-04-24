import mongoose, { Document, Schema } from 'mongoose';

import { RECIPE, RECIPE_COLLECTION, USER, USER_COLLECTION } from './collections';
import { IRecipe } from './recipe.model';

/**
 * IUser document
 *
 * @interface IUser
 * @extends {Document}
 */
interface IUser extends Document {
  TelegramID: number;
  UserName: string;
  Registered: Date;
  Favourites: { RecipeID: string }[];
}

/**
 * User schema
 */
const userSchema: Schema = new Schema({
  TelegramID: { type: Schema.Types.Number, required: true },
  UserName: { type: Schema.Types.String, required: true },
  Registered: { type: Schema.Types.Date, required: true, default: Date.now },
  Favourites: [
    new Schema({
      RecipeID: { type: Schema.Types.ObjectId, ref: RECIPE }
    })
  ]
});

// User mongoose model
const Users = mongoose.model<IUser>(USER, userSchema, USER_COLLECTION);

/**
 * Get users
 *
 * @returns {Promise<IUser[]>}
 */
const getUsers = async (): Promise<IUser[]> => {
  return await Users.find();
};

/**
 * Get user by id
 *
 * @param {string} id
 * @returns {(Promise<IUser | null>)}
 */
const getUser = async (id: string): Promise<IUser | null> => {
  return await Users.findById(id);
};

/**
 * Get favourite recipes by user
 *
 * @param {string} userId
 * @returns {Promise<IRecipe[]>}
 */
const getFavouriteRecipes = async (userId: string): Promise<IRecipe[]> => {
  return await Users.aggregate([
    {
      $match: {
        _id: userId
      }
    },
    {
      $lookup: {
        from: RECIPE_COLLECTION,
        localField: 'Favourites.RecipeID',
        foreignField: '_id',
        as: 'FavouriteRecipe'
      }
    },
    {
      $unwind: {
        path: '$FavouriteRecipe',
        preserveNullAndEmptyArrays: false
      }
    },
    {
      $project: {
        _id: '$FavouriteRecipe._id',
        Name: '$FavouriteRecipe.Name',
        URL: '$FavouriteRecipe.URL',
        Added: '$FavouriteRecipe.Added',
        UserID: '$FavouriteRecipe.UserID'
      }
    },
    {
      $sort: {
        Name: -1
      }
    }
  ]);
};

export { IUser, getUsers, getUser, getFavouriteRecipes };
