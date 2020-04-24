import { Field, ID, ObjectType } from 'type-graphql';

/**
 * Recipe type
 *
 * Property names and casing must match to ones from datasource!
 *
 * @export
 * @class Recipe
 */
@ObjectType({ description: 'The Recipe model' })
export class Recipe {
  @Field(() => ID, { name: 'id' })
  _id!: string;

  @Field(() => String, { name: 'name' })
  Name!: string;

  @Field(() => String, { name: 'url' })
  URL!: string;

  @Field(() => String, { name: 'userId' })
  UserID!: string;

  @Field(() => Date, { name: 'added' })
  Added!: Date;

  @Field(() => Date, { name: 'updated', nullable: true })
  Updated?: Date;
}
