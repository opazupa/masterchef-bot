import mongoose, { Document, Schema } from 'mongoose';

import { RECIPE, RECIPE_COLLECTION, USER } from './collections';

/**
 * IRecipe document
 *
 * @export
 * @interface IRecipe
 * @extends {Document}
 */
export interface IRecipe extends Document {
  UserID: number;
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

export const Recipe = mongoose.model<IRecipe>(RECIPE, recipeSchema, RECIPE_COLLECTION);
