import React from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export const EmailInput = function ({
  email,
  onEmailChange,
  onEnter,
}: {
  email: string;
  onEmailChange: (email: string) => void;
  onEnter: () => void;
}) {
  return (
    <div className="flex flex-col">
      <Label htmlFor="username" className="block mb-2">
        Invite writers
      </Label>
      <Input
        id="email"
        type="email"
        placeholder="Add emails"
        className="w-full"
        data-1p-ignore
        value={email}
        onChange={(e) => onEmailChange(e.target.value)}
        onKeyDown={(e) => {
          onEnter();
        }}
      />
    </div>
  );
};
