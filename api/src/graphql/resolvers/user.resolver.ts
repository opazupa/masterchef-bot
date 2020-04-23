import { ApolloError } from 'apollo-server-express';
import { Arg, Query, Resolver } from 'type-graphql';

import { getUser, getUsers, IUser } from '../../database/models';
import { NOT_FOUND } from '../../errors';
import { User } from '../types/user';

@Resolver((_of) => User)
export class UserResolver {
  @Query((_returns) => User, { nullable: true, description: 'Get user by id' })
  async user(@Arg('id') id: string): Promise<IUser | null> {
    return await getUser(id).catch((e) => {
      console.error(e);
      throw new ApolloError(`Recipe not found with id ${id}`, NOT_FOUND);
    });
  }

  @Query((_returns) => [User], { description: 'Get all users' })
  async users(): Promise<IUser[]> {
    return await getUsers();
  }
}
