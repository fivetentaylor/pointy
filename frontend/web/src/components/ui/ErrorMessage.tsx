import { AlertTriangleIcon } from "lucide-react";

export const ErrorMessage = function ({
  area,
}: {
  area: "document" | "chat" | "page";
  className?: string;
}) {
  return (
    <div className="flex flex-col mx-auto">
      <div className="w-12 h-12 mb-4 bg-rose-100 rounded-full mx-auto flex items-center justify-center">
        <AlertTriangleIcon className="w-6 h-6 stroke-rose-500" />
      </div>

      <p className="text-center text-foreground text-base font-semibold leading-normal">
        {area === "page" ? "This" : "Your"} {area} couldn&apos;t be loaded due
        to an error.
        <br />
        Our team has been notified and will investigate.
      </p>

      <p className="mt-2 text-center text-foreground text-xs">
        Still having issues? Email us at{" "}
        <a href="mailto:support@revi.so">support@revi.so</a>
      </p>
    </div>
  );
};
