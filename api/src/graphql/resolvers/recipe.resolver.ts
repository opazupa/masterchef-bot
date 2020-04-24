import { ApolloError } from 'apollo-server-express';
import { Arg, FieldResolver, Query, Resolver, Root } from 'type-graphql';

import { getAllRecipes, getFavouriters, getRecipe, IRecipe, IUser } from '../../database/models';
import { NOT_FOUND } from '../../errors';
import { Recipe, User } from '../types';

/**
 * Recipe resolver
 *
 * @export
 * @class RecipeResolver
 */
@Resolver((_of) => Recipe)
export class RecipeResolver {
  /**
   * Queries
   */

  @Query((_returns) => Recipe, { description: 'Get Recipe by id', nullable: true })
  async recipe(@Arg('id') id: string): Promise<IRecipe | null> {
    return await getRecipe(id).catch((e) => {
      console.error(e);
      throw new ApolloError(`Recipe not found with id ${id}`, NOT_FOUND);
    });
  }

  @Query((_returns) => [Recipe], { description: 'Get all Recipes' })
  async recipes(): Promise<IRecipe[]> {
    return await getAllRecipes();
  }

  /**
   * Fields
   */

  @FieldResolver((_type) => [User], { description: 'Users who have favourited the recipe', defaultValue: [] })
  async favouritedBy(@Root() recipe: IRecipe): Promise<IUser[]> {
    return await getFavouriters(recipe._id);
  }
}
