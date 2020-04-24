import { ArgsType, Field, ID } from 'type-graphql';

import { RecipeInputType } from '../inputs';

/**
 * Recipe args for creation
 *
 * @class CreateRecipeArgs
 */
@ArgsType()
export class CreateRecipeArgs {
  @Field(() => ID, { description: 'User Id' })
  userId!: string;

  @Field(() => RecipeInputType, { description: 'Recipe info' })
  recipe!: RecipeInputType;
}
