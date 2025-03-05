import { cn } from "@/lib/utils";
import React from "react";
export const TipBox = ({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) => {
  return (
    <div
      className={cn(
        "absolute top-0 left-2 right-2 w-[calc(100%-1rem)] p-2 bg-elevated border border-border rounded-t-md text-xs font-normal transform -translate-y-full",
        className,
      )}
    >
      {children}
    </div>
  );
};
