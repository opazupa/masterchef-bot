import mongoose, { Document, Schema } from 'mongoose';

import { RECIPE, USER, USER_COLLECTION } from './collections';

/**
 * IUser document
 *
 * @export
 * @interface IUser
 * @extends {Document}
 */
export interface IUser extends Document {
  TelegramID: number;
  UserName: string;
  Registered: Date;
  Favourites: { RecipeID: any }[];
}

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

export default mongoose.model<IUser>(USER, userSchema, USER_COLLECTION);
