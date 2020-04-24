import { Field, ID, ObjectType } from 'type-graphql';

/**
 * Recipe type
 *
 * Property names and casing must match to ones from datasource!
 * @export
 * @class User
 */
@ObjectType({ description: 'The User model' })
export class User {
  @Field(() => ID, { name: 'id' })
  _id!: string;

  @Field(() => String, { name: 'userName' })
  UserName!: string;

  @Field(() => Date, { name: 'registered' })
  Registered!: Date;
}
