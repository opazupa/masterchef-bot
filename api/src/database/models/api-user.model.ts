import mongoose, { Document, Schema } from 'mongoose';

import { API_USER, API_USER_COLLECTION } from '../collections';

const APIROLE = {
  ADMIN: 'ADMIN',
  USER: 'USER'
};

type ApiRole = 'ADMIN' | 'USER';

/**
 * IApiUser document
 *
 * @interface IRecipe
 * @extends {Document}
 */
interface IApiUser extends Document {
  Email: string;
  Password: string;
  Roles: ApiRole[];
}

/**
 * ApiUser schema
 */
const apiUserSchema: Schema = new Schema({
  UserName: { type: Schema.Types.ObjectId, required: true },
  Password: { type: Schema.Types.String, required: true },
  Roles: [{ type: Schema.Types.String, required: true }]
});

// ApiUser mongoose model
const ApiUsers = mongoose.model<IApiUser>(API_USER, apiUserSchema, API_USER_COLLECTION);

/**
 * Get user by username
 *
 * @param {string} userName
 * @returns {(Promise<IRecipe | null>)}
 */
const getApiUser = async (userName: string): Promise<IApiUser | null> => {
  return await ApiUsers.findOne({ UserName: userName });
};

export { ApiRole, APIROLE, getApiUser };
