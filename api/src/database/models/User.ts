import mongoose, { Document, Schema } from 'mongoose';

import { RecipeModel, UserCollection, UserModel } from './collections';

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

const UserSchema: Schema = new Schema({
  TelegramID: { type: Schema.Types.Number, required: true },
  UserName: { type: Schema.Types.String, required: true },
  Registered: { type: Schema.Types.Date, required: true, default: Date.now },
  Favourites: [
    new Schema({
      RecipeID: { type: Schema.Types.ObjectId, ref: RecipeModel }
    })
  ]
});

export default mongoose.model<IUser>(UserModel, UserSchema, UserCollection);
