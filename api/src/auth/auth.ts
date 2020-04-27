import bcrypt from 'bcryptjs';
import jsonwebtoken from 'jsonwebtoken';
import { AuthChecker } from 'type-graphql';

import { configuration } from '../configuration';
import { IContext, IContextUser } from '../context';
import { IApiUser } from '../database/models';

const AUTHORIZATION = 'Authorization';
const BEARER_TOKEN_TYPE = 'Bearer';

/**
 * Interface for JWT token data
 *
 * @interface ITokenData
 */
interface ITokenData {
  user: IContextUser;
  exp: number;
  iat?: number;
}

/**
 * Load user from authorization params
 *
 * @param {*} connectionParams
 * @returns
 */
const getUserFromWSParams = (connectionParams: any) => {
  // Check for Authorization property
  if (connectionParams[AUTHORIZATION]) {
    // Load the received JWT token
    const { user } = jsonwebtoken.verify(
      connectionParams[AUTHORIZATION].split(' ')[1],
      configuration.jwtSecret
    ) as ITokenData;

    return <IContext>{ user };
  }
  return null;
};

/**
 * Load user from request user data
 *
 * @param {*} connectionParams
 * @returns
 */
const getUserFromToken = (data: ITokenData): IContextUser | undefined => {
  return data?.user;
};

/**
 * Create JWT token for the user
 *
 * @param {IApiUser} apiUser
 * @returns {{ tokenType: string; token: string; expiresIn: number }}
 */
const createTokenForUser = (apiUser: IApiUser): { tokenType: string; token: string; expiresIn: number } => {
  const expiresIn = Math.floor(Date.now() / 1000) + 60 * 60;
  const token = jsonwebtoken.sign(
    {
      user: <IContextUser>{
        userName: apiUser.UserName,
        roles: apiUser.Roles
      },
      exp: expiresIn
    } as ITokenData,
    configuration.jwtSecret
  );

  return { tokenType: BEARER_TOKEN_TYPE, token, expiresIn };
};

/**
 * Verify user with password hashes
 *
 * @param {IApiUser} apiUser
 * @param {string} password
 * @returns {Promise<boolean>}
 */
const verifyUser = async (apiUser: IApiUser, password: string): Promise<boolean> => {
  return await bcrypt.compare(password, apiUser.Password);
};

/**
 * Hash user password
 *
 * @param {string} password
 * @returns {Promise<string>}
 */
const hashPassword = async (password: string): Promise<string> => {
  return await bcrypt.hash(password, 10);
};

/**
 * Auth check for @Authorized() directives
 *
 * @param {context: IContext}
 * @param {*} roles
 * @returns
 */
const authChecker: AuthChecker<IContext> = ({ context: { user } }, roles) => {
  // No user
  if (!user) {
    return false;
  }

  if (roles.length === 0 || user.roles.some((role) => roles.includes(role))) {
    // Grant access if the roles overlap
    return true;
  }
  // No roles matched
  return false;
};

export { authChecker, getUserFromWSParams, getUserFromToken, createTokenForUser, verifyUser, hashPassword };
