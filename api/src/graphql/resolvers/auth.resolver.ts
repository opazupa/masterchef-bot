import { AuthenticationError } from 'apollo-server-express';
import { Arg, Mutation, Resolver } from 'type-graphql';

import { createTokenForUser, verifyUser } from '../../auth';
import { getApiUser } from '../../database/models';
import { Auth, LoginInputType } from '../types';

/**
 * Auth resolver
 *
 * @export
 * @class AuthResolver
 */
@Resolver()
export class AuthResolver {
  @Mutation(() => Auth, { description: 'Login' })
  async login(@Arg('user') login: LoginInputType): Promise<Auth> {
    const user = await getApiUser(login.userName);

    // Verify user and hashed password
    if (!user || !(await verifyUser(user!, login.password))) {
      throw new AuthenticationError(`No user found or wrong password`);
    }

    // Create token for the user
    const { tokenType, token, expiresIn } = createTokenForUser(user!);
    return { userName: user!.UserName, tokenType, token, expiresIn } as Auth;
  }
}
