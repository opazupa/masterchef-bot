import { buildSchema } from 'type-graphql';

import { RecipeResolver, UserResolver } from '../resolvers/';

const createSchema = buildSchema({
  resolvers: [UserResolver, RecipeResolver],
  emitSchemaFile: true,
  validate: true
});

export { createSchema };
