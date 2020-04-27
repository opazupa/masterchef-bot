import { Field, InputType } from 'type-graphql';

/**
 * Login input type
 *
 * @export
 * @class LoginInputType
 */
@InputType()
export class LoginInputType {
  @Field()
  userName!: string;

  @Field()
  password!: string;
}
