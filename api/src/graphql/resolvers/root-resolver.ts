import { IResolvers } from 'graphql-tools';

const resolverMap: IResolvers = {
  Query: {
    hello(_obj, _args, _context, _info): string {
      return `👋 Hello world! 👋`;
    }
  }
};

export { resolverMap };
