import { cn } from "@/lib/utils";

export function CenteredLayout({
  children,
  variant = "default",
}: {
  children: React.ReactNode;
  variant?: "default" | "wide";
}) {
  return (
    <div
      className={cn(
        "px-4 sm:px-6 lg:px-8 mx-auto w-full max-w-[72rem]",
        variant === "wide" && "max-w-[81rem]",
      )}
    >
      {children}
    </div>
  );
}
