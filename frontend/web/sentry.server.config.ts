// This file configures the initialization of Sentry on the server.
// The config you add here will be used whenever the server handles a request.
// https://docs.sentry.io/platforms/javascript/guides/nextjs/

import * as Sentry from "@sentry/nextjs";

Sentry.init({
  dsn: "https://081b5fe2700396ad052aef08d2fc223b@o4506560295993344.ingest.sentry.io/4506620883894272",

  // Adjust this value in production, or use tracesSampler for greater control
  tracesSampleRate: 1,

  // Setting this option to true will print useful information to the console while you're setting up Sentry.
  debug: false,

  // Disable the Undici integration which is causing issues
  integrations: (integrations) => {
    return integrations.filter((integration) => integration.name !== "Undici");
  },
});
