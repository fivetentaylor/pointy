"use client";

import React, { ReactNode } from "react";
import { Button } from "@/components/ui/button";
import { CopyIcon } from "lucide-react";

interface CopyButtonProps {
  textToCopy: string;
  children: ReactNode;
}

// Use the prop types
export const CopyButton: React.FC<CopyButtonProps> = ({
  textToCopy,
  children,
}) => {
  const handleCopyClick = () => {
    navigator.clipboard
      .writeText(textToCopy)
      .then(() => {
        alert("Text copied to clipboard");
      })
      .catch((err: Error) => {
        alert("Failed to copy text: " + err);
      });
  };

  return (
    <Button variant="ghost" onClick={handleCopyClick}>
      <CopyIcon className="w-4 h-4" />
      {children}
    </Button>
  );
};
