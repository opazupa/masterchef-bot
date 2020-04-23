import { Field, ID, ObjectType } from 'type-graphql';

@ObjectType({ description: 'The Recipe model' })
export class Recipe {
  @Field(() => ID)
  id!: string;

  @Field(() => String, { name: 'name' })
  Name!: String;

  @Field(() => String, { name: 'url' })
  URL!: String;

  @Field(() => String, { name: 'userId' })
  UserID!: String;

  @Field(() => Date, { name: 'added' })
  Added!: Date;
}
