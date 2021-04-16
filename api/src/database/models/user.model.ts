import mongoose, { Document, Schema } from 'mongoose';

import { RECIPE, RECIPE_COLLECTION, USER, USER_COLLECTION } from '../collections';
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
const getAllUsers = async (): Promise<IUser[]> => {
  return await Users.find();
};

/**
 * Get users by ids
 *
 * @param {string[]} ids
 * @returns {(Promise<Map<string, IUser>>)}
 */
const getUsers = async (ids: string[]): Promise<Map<string, IUser>> => {
  const users = await Users.find({ _id: { $in: ids } });
  return new Map(users.map((user) => [user._id.toString(), user]));
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
 * Get favourite recipes by users
 *
 * @param {string[]} userIds
 * @returns {Promise<IRecipe[]>}
 */
const getFavouriteRecipes = async (userIds: string[]): Promise<Map<string, IRecipe[]>> => {
  const recipes = (await Users.aggregate([
    {
      $match: {
        _id: { $in: userIds }
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
  ])) as IRecipe[];

  return new Map(userIds.map((userId) => [userId, recipes.filter((r) => r.UserID.toString() === userId.toString())]));
};

export { IUser, getAllUsers, getUsers, getUser, getFavouriteRecipes };
