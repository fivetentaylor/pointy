import * as Sentry from "@sentry/browser";

Sentry.init({
  dsn: "https://081b5fe2700396ad052aef08d2fc223b@o4506560295993344.ingest.sentry.io/4506620883894272",

  // Alternatively, use `process.env.npm_package_version` for a dynamic release version
  // if your build tool supports it.
  release: `reviso-app@${process.env.IMAGE_TAG}`,
  integrations: [
    Sentry.browserTracingIntegration(),
    Sentry.replayIntegration(),
  ],

  // Set tracesSampleRate to 1.0 to capture 100%
  // of transactions for tracing.
  // We recommend adjusting this value in production
  tracesSampleRate: 1.0,

  // Set `tracePropagationTargets` to control for which URLs trace propagation should be enabled
  tracePropagationTargets: [
    "localhost",
    /^https:\/\/app\.reviso\.biz/,
    /^https:\/\/app\.reviso\.dev/,
    /^https:\/\/app\.revi\.so/,
  ],

  // Capture Replay for 10% of all sessions,
  // plus for 100% of sessions with an error
  replaysSessionSampleRate: 0.1,
  replaysOnErrorSampleRate: 1.0,
});
