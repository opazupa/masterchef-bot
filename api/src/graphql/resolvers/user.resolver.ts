import { ApolloError } from 'apollo-server-express';
import { IResolvers } from 'graphql-tools';

import { getFavouriteRecipes, getUser, getUsers, IUser } from '../../database/models';
import { NOT_FOUND } from '../../errors';

const userResolvers: IResolvers = {
  // Type
  User: {
    id: (obj: IUser) => obj._id,
    userName: (obj: IUser) => obj.UserName,
    registered: (obj: IUser) => obj.Registered,
    favourites: async (obj: IUser) => await getFavouriteRecipes(obj._id)
  },

  // Queries
  Query: {
    users: async () => {
      return await getUsers();
    },
    user: async (_obj: IUser, args: { id: string }) => {
      return await getUser(args.id).catch((e) => {
        console.error(e);
        throw new ApolloError(`User not found with id ${args.id}`, NOT_FOUND);
      });
    }
  }
};

export { userResolvers };
