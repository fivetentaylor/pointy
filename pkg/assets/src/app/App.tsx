import React from "react";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import Draft from "./pages/Draft";
import GraphQL from "./contexts/GraphQL";
import posthog from "posthog-js";
import { PostHogProvider } from "posthog-js/react";
import { CurrentUserProvider } from "./contexts/CurrentUserContext";
import { Toaster } from "./components/ui/toaster";

const router = createBrowserRouter([
  {
    path: "/drafts/:draftId/:threadId?",
    element: <Draft />,
  },
]);

function App() {
  return (
    <PostHogProvider client={posthog}>
      <GraphQL>
        <CurrentUserProvider>
          <Toaster />
          <RouterProvider router={router} />
        </CurrentUserProvider>
      </GraphQL>
    </PostHogProvider>
  );
}

export default App;
