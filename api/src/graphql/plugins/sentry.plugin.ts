import * as Sentry from '@sentry/node';
import { PluginDefinition } from 'apollo-server-core';
import { ApolloError } from 'apollo-server-express';
import { GraphQLRequestContextDidEncounterErrors } from 'apollo-server-plugin-base';

import { IContext } from '../../context';

export const sentryPlugin = {
  requestDidStart(_: any) {
    /* Within this returned object, define functions that respond
    to request-specific lifecycle events. */
    return {
      didEncounterErrors(ctx: GraphQLRequestContextDidEncounterErrors<IContext>) {
        // If we couldn't parse the operation, don't
        // do anything here
        if (!ctx.operation) {
          return;
        }

        for (const err of ctx.errors) {
          // Only report internal server errors,
          // all errors extending ApolloError should be user-facing
          if (err.originalError instanceof ApolloError) {
            continue;
          }

          // Add scoped report details and send to Sentry
          Sentry.withScope((scope) => {
            // Annotate whether failing operation was query/mutation/subscription
            scope.setTag('kind', ctx.operation!.operation);

            // Log query and variables as extras (make sure to strip out sensitive data!)
            if (ctx.context.user) {
              scope.setExtra('user', ctx.context.user.userName);
            }
            scope.setExtra('query', ctx.request.query);
            scope.setExtra('variables', ctx.request.variables);

            if (err.path) {
              // We can also add the path as breadcrumb
              scope.addBreadcrumb({
                category: 'query-path',
                message: err.path.join(' > '),
                level: Sentry.Severity.Debug
              });
            }

            Sentry.captureException(err);
          });
        }
      }
    };
  }
} as PluginDefinition;
