import { ArgsType, Field, ID } from 'type-graphql';

import { RecipeInputType } from '../inputs';

/**
 * Recipe args for updating
 *
 * @class UpdateRecipeArgs
 */
@ArgsType()
export class UpdateRecipeArgs {
  @Field(() => ID, { description: 'Id' })
  id!: string;

  @Field(() => RecipeInputType, { description: 'Recipe info' })
  recipe!: RecipeInputType;
}
