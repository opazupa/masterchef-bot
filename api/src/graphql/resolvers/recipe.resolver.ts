import { ApolloError } from 'apollo-server-express';
import { Arg, Query, Resolver } from 'type-graphql';

import { getAllRecipes, getRecipe, IRecipe } from '../../database/models';
import { NOT_FOUND } from '../../errors';
import { Recipe } from '../types/recipe';

@Resolver((_of) => Recipe)
export class RecipeResolver {
  @Query((_returns) => Recipe, { nullable: true, description: 'Get Recipe by id' })
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
}
