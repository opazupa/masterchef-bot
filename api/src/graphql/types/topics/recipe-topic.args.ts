import { ArgsType, Field, ID } from 'type-graphql';

/**
 * Recipe subscription topic args
 *
 * @export
 * @class RecipeTopicArgs
 */
@ArgsType()
export class RecipeTopicArgs {
  @Field(() => ID, { description: 'User Id filter', nullable: true })
  userId?: string;

  @Field(() => String, { description: 'Recipe name like-filter', nullable: true })
  name?: string;

  @Field(() => String, { description: 'Recipe url like-filter', nullable: true })
  url?: string;
}
