import { GraphQLSchema } from 'graphql';
import { buildSchema } from 'type-graphql';

import { authChecker } from '../../auth';
import { AuthResolver, RecipeResolver, UserResolver } from '../resolvers/';

/**
 * Generates graphql schema
 *
 * @returns {Promise<GraphQLSchema>}
 */
const createSchema: Promise<GraphQLSchema> = buildSchema({
  resolvers: [UserResolver, RecipeResolver, AuthResolver],
  emitSchemaFile: true,
  validate: true,
  authChecker: authChecker
});

export { createSchema };
