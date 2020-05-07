import { ApolloError } from 'apollo-server-express';
import {
  Args,
  Authorized,
  Ctx,
  FieldResolver,
  ID,
  Mutation,
  Publisher,
  PubSub,
  Query,
  Resolver,
  ResolverFilterData,
  Root,
  Subscription,
} from 'type-graphql';

import { IContext } from '../../context';
import {
  addRecipe,
  APIROLE,
  deleteRecipe,
  getAllRecipes,
  getRecipe,
  IRecipe,
  IUser,
  updateRecipe,
} from '../../database/models';
import { NOT_FOUND } from '../../errors';
import { CreateRecipeArgs, IdArg, Recipe, UpdateRecipeArgs, User } from '../types';
import { RecipeTopicArgs, Topics } from '../types/topics';

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
      // tslint:disable-next-line: no-console
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
  @Authorized()
  @Mutation(() => Recipe, { description: 'Add recipe' })
  async addRecipe(
    @Args() { userId, recipe }: CreateRecipeArgs,
    @PubSub(Topics.NewRecipe) newRecipeNotification: Publisher<Recipe>
  ) {
    const added = await addRecipe(userId, recipe.name, recipe.url);
    await newRecipeNotification(added);
    return added;
  }

  @Authorized()
  @Mutation(() => Recipe, { description: 'Update recipe' })
  async updateRecipe(
    @Args() { id, recipe }: UpdateRecipeArgs,
    @PubSub(Topics.UpdatedRecipe) updatedRecipeNotification: Publisher<Recipe>
  ) {
    const updated = await updateRecipe(id, recipe.name, recipe.url).catch((e) => {
      // tslint:disable-next-line: no-console
      console.error(e);
      throw new ApolloError(`Recipe not found with id ${id}`, NOT_FOUND);
    });
    await updatedRecipeNotification(updated!);
    return updated;
  }

  @Authorized(APIROLE.ADMIN)
  @Mutation(() => ID, { description: 'Delete recipe' })
  async deleteRecipe(
    @Args() { id }: IdArg,
    @PubSub(Topics.DeletedRecipe) deletedRecipeNotification: Publisher<string>
  ) {
    await deleteRecipe(id).catch((e) => {
      // tslint:disable-next-line: no-console
      console.error(e);
      throw new ApolloError(`Recipe not found with id ${id}`, NOT_FOUND);
    });
    await deletedRecipeNotification(id);
    return id;
  }

  /**
   * Subscriptions
   */
  @Authorized()
  @Subscription({
    topics: Topics.NewRecipe,
    filter: ({ payload, args }: ResolverFilterData<IRecipe, RecipeTopicArgs>) => {
      if (args.userId) {
        return args.userId.toString() === payload.UserID.toString();
      }
      if (args.name) {
        return payload.Name.includes(args.name);
      }
      if (args.url) {
        return payload.URL.includes(args.url);
      }
      return true;
    },
    description: 'Notification on new recipe'
  })
  newRecipe(@Root() recipe: IRecipe, @Args() _args: RecipeTopicArgs): Recipe {
    return recipe;
  }

  @Authorized()
  @Subscription({
    topics: Topics.UpdatedRecipe,
    filter: ({ payload, args }: ResolverFilterData<IRecipe, RecipeTopicArgs>) => {
      if (args.userId) {
        return args.userId.toString() === payload.UserID.toString();
      }
      if (args.name) {
        return payload.Name.includes(args.name);
      }
      if (args.url) {
        return payload.URL.includes(args.url);
      }
      return true;
    },
    description: 'Notification on recipe update'
  })
  updatedRecipe(@Root() recipe: IRecipe, @Args() _args: RecipeTopicArgs): Recipe {
    return recipe;
  }

  @Authorized()
  @Subscription({
    topics: Topics.DeletedRecipe,
    description: 'Notification on deleted recipes'
  })
  deletedRecipe(@Root() id: string): string {
    return id;
  }

  /**
   * Fields
   */

  @FieldResolver((_type) => [User], {
    description: 'Users who have favourited the recipe',
    defaultValue: [],
    nullable: true
  })
  async favouritedBy(@Root() recipe: IRecipe, @Ctx() ctx: IContext): Promise<IUser[]> {
    return await ctx.loaders.recipe.favouriters.load(recipe._id);
  }

  @FieldResolver((_type) => User, { description: 'User who added the recipe' })
  async user(@Root() recipe: IRecipe, @Ctx() ctx: IContext): Promise<IUser | null> {
    return await ctx.loaders.recipe.user.load(recipe.UserID).catch((e) => {
      // tslint:disable-next-line: no-console
      console.error(e);
      throw new ApolloError(`User not found with id ${recipe.UserID}`, NOT_FOUND);
    });
  }
}
