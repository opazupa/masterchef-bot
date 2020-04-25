import { ApolloError } from 'apollo-server-express';
import { Args, FieldResolver, Query, Resolver, Root } from 'type-graphql';

import { getFavouriteRecipes, getUser, getUserRecipes, getUsers, IRecipe, IUser } from '../../database/models';
import { NOT_FOUND } from '../../errors';
import { IdArg, Recipe, User } from '../types';

/**
 * User resolver
 *
 * @export
 * @class UserResolver
 */
@Resolver((_of) => User)
export class UserResolver {
  /**
   * Queries
   */

  @Query((_returns) => User, { nullable: true, description: 'Get user by id' })
  async user(@Args() { id }: IdArg): Promise<IUser | null> {
    return await getUser(id).catch((e) => {
      console.error(e);
      throw new ApolloError(`Recipe not found with id ${id}`, NOT_FOUND);
    });
  }

  @Query((_returns) => [User], { description: 'Get all users' })
  async users(): Promise<IUser[]> {
    return await getUsers();
  }

  /**
   * Fields
   */

  @FieldResolver((_type) => [Recipe], { defaultValue: [] })
  async favourites(@Root() user: IUser): Promise<IRecipe[]> {
    return await getFavouriteRecipes(user._id);
  }

  /**
   * Fields
   */

  @FieldResolver((_type) => [Recipe], { defaultValue: [] })
  async recipes(@Root() user: IUser): Promise<IRecipe[]> {
    return await getUserRecipes(user._id);
  }
}
