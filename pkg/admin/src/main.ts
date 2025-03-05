import 'htmx.org';
import './wasm_exec'

import './rogueId';

declare global {
  var RevisoVersion: string;
  var process: {
    env: {
      APP_HOST: string;
      WS_HOST: string;
      NODE_ENV: string;
      WEB_HOST: string;
      SEGMENT_KEY: string;
      PUBLIC_POSTHOG_KEY: string;
      PUBLIC_POSTHOG_HOST: string;
      IMAGE_TAG: string;
    };
  };
}
