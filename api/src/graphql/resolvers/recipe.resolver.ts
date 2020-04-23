import { ApolloError } from 'apollo-server-express';
import { IResolvers } from 'graphql-tools';

import { getFavouritedUsers, getRecipe, IRecipe } from '../../database/models';
import { NOT_FOUND } from '../../errors';

const recipeResolvers: IResolvers = {
  // Type
  Recipe: {
    id: (obj: IRecipe) => obj._id,
    name: (obj: IRecipe) => obj.Name,
    url: (obj: IRecipe) => obj.URL,
    userId: (obj: IRecipe) => obj.UserID,
    added: (obj: IRecipe) => obj.Added,
    favourited: async (obj: IRecipe) => await getFavouritedUsers(obj._id)
  },

  // Queries
  Query: {
    // recipes: async () => {
    //   return await getAllRecipes();
    // },
    recipe: async (_obj: IRecipe, args: { id: string }) => {
      return await getRecipe(args.id).catch((e) => {
        console.error(e);
        throw new ApolloError(`Recipe not found with id ${args.id}`, NOT_FOUND);
      });
    }
  }
};

export { recipeResolvers };
