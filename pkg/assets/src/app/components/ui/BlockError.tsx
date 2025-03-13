import React from "react";
import { AlertTriangleIcon } from "lucide-react";

export const BlockError = ({ text }: { text: string }) => {
  return (
    <section className="flex flex-col items-center justify-center h-[calc(100dvh-6.2rem)] text-center">
      <div className="mb-4">
        <div className="w-12 h-12 bg-rose-100 rounded-full flex items-center justify-center">
          <AlertTriangleIcon className="w-6 h-6 stroke-rose-500" />
        </div>
      </div>
      <p className="text-[1rem] leading-1.5rem] font-semibold mb-4">
        {text}
        <br />
        Our team has been notified and will investigate.
      </p>
      <p className="mt-2 text-center text-foreground text-xs">
        Still having issues? Email us at{" "}
        <a href="mailto:taylor@pointy.ai">taylor@pointy.ai</a>
      </p>
    </section>
  );
};
