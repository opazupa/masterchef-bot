import { IResolvers } from 'graphql-tools';
import { merge } from 'lodash';

import { recipeResolvers } from './recipe.resolver';
import { userResolvers } from './user.resolver';

const resolverMap: IResolvers = merge({}, recipeResolvers, userResolvers);

export { resolverMap };
