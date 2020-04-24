import { ApolloError } from 'apollo-server-express';
import { Args, FieldResolver, ID, Mutation, Query, Resolver, Root } from 'type-graphql';

import {
  addRecipe,
  deleteRecipe,
  getAllRecipes,
  getFavouriters,
  getRecipe,
  IRecipe,
  IUser,
  updateRecipe,
} from '../../database/models';
import { NOT_FOUND } from '../../errors';
import { CreateRecipeArgs, IdArg, Recipe, UpdateRecipeArgs, User } from '../types';

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
  async recipe(@Args() { id }: IdArg): Promise<IRecipe | null> {
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
   * Mutations
   */

  @Mutation(() => Recipe, { description: 'Add recipe' })
  async addRecipe(@Args() { userId, recipe }: CreateRecipeArgs) {
    return await addRecipe(userId, recipe.name, recipe.url);
  }

  @Mutation(() => Recipe, { description: 'Update recipe' })
  async updateRecipe(@Args() { id, recipe }: UpdateRecipeArgs) {
    return await updateRecipe(id, recipe.name, recipe.url).catch((e) => {
      console.error(e);
      throw new ApolloError(`Recipe not found with id ${id}`, NOT_FOUND);
    });
  }

  @Mutation(() => ID, { description: 'Delete recipe' })
  async deleteRecipe(@Args() { id }: IdArg) {
    await deleteRecipe(id).catch((e) => {
      console.error(e);
      throw new ApolloError(`Recipe not found with id ${id}`, NOT_FOUND);
    });
    return id;
  }

  /**
   * Fields
   */

  @FieldResolver((_type) => [User], { description: 'Users who have favourited the recipe', defaultValue: [] })
  async favouritedBy(@Root() recipe: IRecipe): Promise<IUser[]> {
    return await getFavouriters(recipe._id);
  }
}
