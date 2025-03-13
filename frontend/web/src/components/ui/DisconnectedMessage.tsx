import { CloudOffIcon } from "lucide-react";

export const DisconnectedMessage = function ({}: {}) {
  return (
    <div className="flex flex-col mx-auto">
      <div className="w-12 h-12 mb-4 bg-zinc-100 rounded-full mx-auto flex items-center justify-center">
        <CloudOffIcon className="w-6 h-6 stroke-zinc-500" />
      </div>

      <p className="text-center text-foreground text-base font-semibold leading-normal">
        No connnection.
      </p>
      <p className="text-center text-foreground text-base leading-normal">
        Go online to get the latest messages.
      </p>
      <p className="mt-2 text-center text-foreground text-xs">
        Still having issues? Email us at{" "}
        <a href="mailto:taylor@pointy.ai">taylor@pointy.ai</a>
      </p>
    </div>
  );
};
