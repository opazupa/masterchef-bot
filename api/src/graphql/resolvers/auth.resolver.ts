import { ApolloError } from 'apollo-server-express';
import { Arg, Mutation, Resolver } from 'type-graphql';

import { createTokenForUser, verifyUser } from '../../auth';
import { getApiUser } from '../../database/models';
import { NOT_FOUND, UN_AUTHORIZED } from '../../errors';
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
    const user = await getApiUser(login.userName).catch((e) => {
      console.error(e);
      throw new ApolloError(`User not found with username ${login.userName}`, NOT_FOUND);
    });

    // Verify hashed pasword
    if (!(await verifyUser(user!, login.password))) {
      throw new ApolloError(`Credentieals doesn't match`, UN_AUTHORIZED);
    }

    // Create token for the user
    const { tokenType, token, expiresIn } = createTokenForUser(user!);
    return <Auth>{ userName: user!.UserName, tokenType, token, expiresIn };
  }
}
