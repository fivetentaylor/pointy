import React, { useState } from "react";
import { Button } from "@/components/ui/button";

interface ConfirmButtonProps {
  children: React.ReactNode;
  onConfirm: () => void;
}

export const ConfirmButton = ({ children, onConfirm }: ConfirmButtonProps) => {
  const [showConfirmation, setShowConfirmation] = useState(false);

  return (
    <div>
      {showConfirmation ? (
        <div className="animate-in slide-in-from-right">
          <span className="px-2">Are you sure?</span>
          <Button
            className="mx-2"
            variant="outline"
            onClick={(event) => {
              event.stopPropagation();
              setShowConfirmation(false);
            }}
          >
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={(event) => {
              event.stopPropagation();
              // 2. Invoke onConfirm when "Yes" is clicked
              onConfirm();
              setShowConfirmation(false);
            }}
          >
            Yes
          </Button>
        </div>
      ) : (
        <Button
          variant="destructive"
          onClick={(event) => {
            event.stopPropagation();
            setShowConfirmation(true);
          }}
        >
          {children}
        </Button>
      )}
    </div>
  );
};
