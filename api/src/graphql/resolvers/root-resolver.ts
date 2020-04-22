import { IResolvers } from 'graphql-tools';

const resolverMap: IResolvers = {
  Query: {
    user(_: void, __: void): string {
      return `👋 Hello world! 👋`;
    }
  }
};

export { resolverMap };
