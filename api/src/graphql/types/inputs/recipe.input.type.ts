import { Field, InputType } from 'type-graphql';

/**
 * Recipe input type
 *
 * @export
 * @class RecipeInputType
 */
@InputType()
export class RecipeInputType {
  @Field()
  name!: string;

  @Field()
  url!: string;
}
