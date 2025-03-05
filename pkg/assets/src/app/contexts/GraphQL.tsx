import React from "react";

import { ApolloProvider } from "@apollo/client";

import { client } from "@/lib/graphql";

export default function GraphQL({ children }: { children: React.ReactNode }) {
  return <ApolloProvider client={client}>{children}</ApolloProvider>;
}
