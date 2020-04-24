import { GraphQLSchema } from 'graphql';
import { buildSchema } from 'type-graphql';

import { RecipeResolver, UserResolver } from '../resolvers/';

/**
 * Generates graphql schema
 *
 * @returns {Promise<GraphQLSchema>}
 */
const createSchema: Promise<GraphQLSchema> = buildSchema({
  resolvers: [UserResolver, RecipeResolver],
  emitSchemaFile: true,
  validate: true
});

export { createSchema };
