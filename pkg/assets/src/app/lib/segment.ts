import { AnalyticsBrowser } from "@segment/analytics-next";
import posthog from "posthog-js";

export const analytics = AnalyticsBrowser.load({
  writeKey: process.env.SEGMENT_KEY!,
});

analytics.ready(() => {
  posthog.init(process.env.PUBLIC_POSTHOG_KEY, {
    api_host: process.env.PUBLIC_POSTHOG_HOST,
    segment: analytics.instance! as any,
    capture_pageview: false,
    loaded: () => analytics.page(),
  });
});
