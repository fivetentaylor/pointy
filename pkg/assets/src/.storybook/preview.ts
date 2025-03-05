import type { Preview } from "@storybook/react";
import "../style/main.css";

import { MockedProvider } from "@apollo/client/testing";

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
    apolloClient: {
      MockedProvider,
      // any props you want to pass to MockedProvider on every story
    },
  },
};

export default preview;
