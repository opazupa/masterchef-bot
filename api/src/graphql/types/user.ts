import { Field, FieldResolver, ID, ObjectType, Root } from 'type-graphql';

import { getFavouriteRecipes, IRecipe } from '../../database/models';
import { Recipe } from './recipe';

@ObjectType({ description: 'The User model' })
export class User {
  @Field(() => ID)
  id!: string;

  @Field(() => String, { name: 'userName' })
  UserName?: String;

  @Field(() => Date, { name: 'registered' })
  Registered?: Date;

  @FieldResolver((_type) => [Recipe], { defaultValue: [] })
  async favourites(@Root() user: User): Promise<IRecipe[]> {
    return await getFavouriteRecipes(user.id);
  }
}
