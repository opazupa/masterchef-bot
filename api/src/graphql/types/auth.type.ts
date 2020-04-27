import { Field, ObjectType } from 'type-graphql';

/**
 * Auth type
 *
 * Property names and casing must match to ones from datasource!
 * @export
 * @class Auth
 */
@ObjectType({ description: 'The login auth model' })
export class Auth {
  @Field(() => String, { description: 'UserName' })
  userName!: string;

  @Field(() => String, { description: 'Token' })
  token!: string;

  @Field(() => String, { description: 'Token type' })
  tokenType!: string;

  @Field(() => Number, { description: 'Token expiration' })
  expiresIn!: number;
}
