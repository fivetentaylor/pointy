import { ApolloClient, InMemoryCache, split, from } from "@apollo/client";
import { onError } from "@apollo/client/link/error";
import { getMainDefinition } from "@apollo/client/utilities";
import { GraphQLWsLink } from "@apollo/client/link/subscriptions";
import { createClient } from "graphql-ws";
import { loadErrorMessages, loadDevMessages } from "@apollo/client/dev";
import { wsEventEmitter } from "./wsEmitter";
import { WS_HOST } from "@/lib/urls";
import createUploadLink from "apollo-upload-client/createUploadLink.mjs";

if (process.env.NODE_ENV === "development") {
  loadDevMessages();
  loadErrorMessages();
}

const retryWait = (attempts: number) => {
  // Calculate the base delay
  const baseDelay = 250 * 2 ** attempts;
  // Add jitter by introducing a random factor
  const jitter = Math.random() * baseDelay * 0.1;
  // Calculate the total delay with jitter
  const totalDelay = baseDelay + jitter;
  console.log("Attempting to reconnect to GQL WebSocket in ", totalDelay);
  return new Promise<void>((resolve) => setTimeout(resolve, totalDelay));
};

let websocketOpen = false;
let emitCloseTimeout: NodeJS.Timeout;

const wsLink =
  typeof window !== "undefined"
    ? new GraphQLWsLink(
        createClient({
          url: WS_HOST + "/graphql/query",
          retryAttempts: 100, // Maximum retry attempts
          retryWait, // Function to determine retry delay
          shouldRetry: (error) => {
            return true;
          },
          on: {
            closed: () => {
              websocketOpen = false;
              // delay the emit to avoid hiccups
              emitCloseTimeout = setTimeout(() => {
                if (!websocketOpen) {
                  wsEventEmitter.emit("wsStatus", { open: websocketOpen });
                }
              }, 2000);
            },
            opened: () => {
              websocketOpen = true;
              clearTimeout(emitCloseTimeout);
              wsEventEmitter.emit("wsStatus", { open: websocketOpen });
            },
            error: () => {
              websocketOpen = false;
              // delay the emit to avoid hiccups
              emitCloseTimeout = setTimeout(() => {
                if (!websocketOpen) {
                  wsEventEmitter.emit("wsStatus", { open: websocketOpen });
                }
              }, 2000);
            },
          },
        }),
      )
    : null;

const httpLink = createUploadLink({
  uri: "/graphql/query",
  credentials: "include",
});

const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (typeof window !== "undefined" && wsLink != null) {
    if (graphQLErrors) {
      const hasAuthError =
        graphQLErrors.filter(
          ({ message, locations, path }) =>
            message === "we could not find you" || message === "please login",
        ).length > 0;

      if (hasAuthError) {
        window.location.href = "/login";
      }
    }
  }
});

const link =
  typeof window !== "undefined" && wsLink != null
    ? split(
        ({ query }) => {
          const def = getMainDefinition(query);
          return (
            def.kind === "OperationDefinition" &&
            def.operation === "subscription"
          );
        },
        wsLink,
        httpLink,
      )
    : httpLink;

const cache = new InMemoryCache({
  typePolicies: {
    TLMessageV1: {
      keyFields: ["eventId"],
    },
  },
});

export const client = new ApolloClient({
  connectToDevTools: process.env.NODE_ENV === "development",
  cache:
    typeof window !== "undefined"
      ? cache.restore((window as any).__APOLLO_STATE__)
      : cache,
  link: from([errorLink, link]),
});
