import React from "react";
import { ToastAction } from "@/components/ui/toast";
import { useToast } from "@/components/ui/use-toast";
import { AlertTriangleIcon } from "lucide-react";

export const useErrorToast = function () {
  const { toast } = useToast();
  return (message: string) => {
    toast({
      hideCloseButton: true,
      title: (
        <div className="flex gap-2 items-center">
          <AlertTriangleIcon className="w-4 h-4" />
          <span className="font-semibold text-sm"> {message} </span>
        </div>
      ),
      description: (
        <span className="text-muted-foreground text-xs">
          {" "}
          Please try again{" "}
        </span>
      ),
      action: <ToastAction altText="Dismiss"> Dismiss </ToastAction>,
    });
  };
};
