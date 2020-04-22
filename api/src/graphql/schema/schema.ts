import 'graphql-import-node';

import { makeExecutableSchema } from 'apollo-server-express';
import { GraphQLSchema } from 'graphql';

import { typeDefs } from '../types';
import { resolverMap as resolvers } from './../resolvers';
import * as schemaDef from './schema.graphql';

const schema: GraphQLSchema = makeExecutableSchema({
  typeDefs: [schemaDef, ...typeDefs],
  resolvers
});

export { schema };
