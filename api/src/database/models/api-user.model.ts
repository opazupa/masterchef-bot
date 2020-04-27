import mongoose, { Document, Schema } from 'mongoose';

import { API_USER, API_USER_COLLECTION } from '../collections';

// API User roles
const APIROLE = {
  ADMIN: 'ADMIN',
  USER: 'USER'
};

// Type for api role
type ApiRole = 'ADMIN' | 'USER';

/**
 * IApiUser document
 *
 * @interface IApiUser
 * @extends {Document}
 */
interface IApiUser extends Document {
  UserName: string;
  Password: string;
  Roles: ApiRole[];
  Added: Date;
  Updated: Date;
}

/**
 * ApiUser schema
 */
const apiUserSchema: Schema = new Schema(
  {
    UserName: { type: Schema.Types.String, required: true },
    Password: { type: Schema.Types.String, required: true },
    Roles: [{ type: Schema.Types.String }],
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

/**
 * Add ApiUser
 *
 * @param {string} userName
 * @param {string} password
 * @param {ApiRole[]} [roles=[]]
 */
const addApiUser = async (userName: string, password: string, roles: ApiRole[] = []) => {
  await ApiUsers.create({ UserName: userName, Password: password, Roles: roles } as IApiUser);
};

export { IApiUser, ApiRole, APIROLE, getApiUser, addApiUser };
