import { ArgsType, Field, ID } from 'type-graphql';

/**
 * Id arg
 *
 * @class IdArg
 */
@ArgsType()
export class IdArg {
  @Field(() => ID, { description: 'Id' })
  id!: string;
}
