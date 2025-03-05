import Analytics from "@/components/Analytics";
import { MarketingFooter } from "@/components/marketing/MarketingFooter";

export default function MarketingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex flex-col min-h-screen">
      <div>{children}</div>
      <MarketingFooter />
      <Analytics />
    </div>
  );
}
