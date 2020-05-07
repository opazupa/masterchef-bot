import { ApolloError } from 'apollo-server-express';
import { Args, Ctx, FieldResolver, Query, Resolver, Root } from 'type-graphql';

import { IContext } from '../../context';
import { getAllUsers, getUser, IRecipe, IUser } from '../../database/models';
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
      // tslint:disable-next-line: no-console
      console.error(e);
      throw new ApolloError(`Recipe not found with id ${id}`, NOT_FOUND);
    });
  }

  @Query((_returns) => [User], { description: 'Get all users' })
  async users(): Promise<IUser[]> {
    return await getAllUsers();
  }

  /**
   * Fields
   */

  @FieldResolver((_type) => [Recipe], { description: 'Favourite recipes', defaultValue: [], nullable: true })
  async favourites(@Root() user: IUser, @Ctx() ctx: IContext): Promise<IRecipe[]> {
    return ctx.loaders.user.favourites.load(user._id);
  }

  @FieldResolver((_type) => [Recipe], { description: 'Added recipes', defaultValue: [], nullable: true })
  async recipes(@Root() user: IUser, @Ctx() ctx: IContext): Promise<IRecipe[]> {
    return await ctx.loaders.user.recipes.load(user._id);
  }
}
