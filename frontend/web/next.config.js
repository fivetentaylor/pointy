module.exports = {
  experimental: {
    instrumentationHook: true,
    // swcPlugins: [
    //   [
    //     "@preact-signals/safe-react/swc",
    //     {
    // you should use `auto` mode to track only components which uses `.value` access.
    // Can be useful to avoid tracking of server side components
    //       mode: "auto",
    //     } /* plugin options here */,
    //   ],
    //  ],
  },
  reactStrictMode: false,
  output: "standalone",
  redirects: async () => {
    return [
      {
        source: "/documents",
        destination: `${process.env.NEXT_PUBLIC_APP_HOST}/drafts`,
        permanent: true,
        basePath: false,
      },
      {
        source: "/documents/:id",
        destination: `${process.env.NEXT_PUBLIC_APP_HOST}/drafts/:id`,
        permanent: true,
      },
      {
        source: "/read/:id",
        destination: `${process.env.NEXT_PUBLIC_API_HOST}/read/:id`,
        permanent: true,
      },
    ];
  },
};

// Injected content via Sentry wizard below

const { withSentryConfig } = require("@sentry/nextjs");

module.exports = withSentryConfig(
  module.exports,
  {
    // For all available options, see:
    // https://github.com/getsentry/sentry-webpack-plugin#options

    // Suppresses source map uploading logs during build
    silent: true,
    org: "reviso-z7",
    project: "web",
  },
  {
    // For all available options, see:
    // https://docs.sentry.io/platforms/javascript/guides/nextjs/manual-setup/

    // Upload a larger set of source maps for prettier stack traces (increases build time)
    widenClientFileUpload: true,

    // Transpiles SDK to be compatible with IE11 (increases bundle size)
    transpileClientSDK: true,

    // Routes browser requests to Sentry through a Next.js rewrite to circumvent ad-blockers (increases server load)
    tunnelRoute: "/monitoring",

    // Hides source maps from generated client bundles
    hideSourceMaps: true,

    // Automatically tree-shake Sentry logger statements to reduce bundle size
    disableLogger: true,

    // Enables automatic instrumentation of Vercel Cron Monitors.
    // See the following for more information:
    // https://docs.sentry.io/product/crons/
    // https://vercel.com/docs/cron-jobs
    automaticVercelMonitors: true,
  },
);
