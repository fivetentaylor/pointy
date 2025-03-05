import React, { useCallback, useEffect, useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

import { XIcon } from "lucide-react";
import { Textarea } from "@/components/ui/textarea";
import { cn } from "@/lib/utils";
import { analytics } from "@/lib/segment";
import {
  SHARE_ADD_CUSTOM_MESSAGE,
  SHARE_ADD_EMAIL,
  SHARE_REMOVE_EMAIL,
} from "@/lib/events";

export const cleanupEmail = function (email: string) {
  const emailMatch = email.match(/<(.+?)>$/);
  return emailMatch ? emailMatch[1] : email.trim();
};

export const isValidEmail = function (email: string) {
  return /.*@.*\..+/.test(email);
};

export const EmailView = function ({
  email,
  setEmail,
  emails,
  setEmails,
  loadingshareDocument,
  onSubmit,
}: {
  email: string;
  setEmail: (email: string) => void;
  emails: string[];
  setEmails: (emails: string[]) => void;
  loadingshareDocument: boolean;
  onSubmit: (message?: string) => void;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [message, setMessage] = useState("");

  const focusOnInput = useCallback(() => {
    inputRef?.current?.focus();
  }, [inputRef]);

  useEffect(() => {
    if (inputRef.current) {
      focusOnInput();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const hasValidEmail =
    isValidEmail(email) || emails.filter(isValidEmail).length > 0;

  const handleAddEmail = (newEmail: string) => {
    const cleanEmail = cleanupEmail(newEmail);
    if (cleanEmail.length > 0 && cleanEmail.length <= 254) {
      analytics.track(SHARE_ADD_EMAIL);
      setEmails([...emails, cleanEmail]);
      setEmail("");
    }
  };

  return (
    <div className="flex flex-col">
      <h3 className="text-foreground text-lg font-semibold">
        Invite others to write with you
      </h3>
      <span className="text-muted-foreground font-light">
        Invited writers are able to access the document, edit it, and leave
        comments.
      </span>
      <h4 className="text-foreground text-sm mt-4 mb-2">People with access</h4>
      <div className="w-full border-input border-[1px] min-h-10 rounded-md px-3 flex flex-wrap mb-4 pb-2">
        {emails.map((thisEmail, index) => (
          <div
            key={`email-${index}`}
            className={cn(
              "border-[1px] rounded-md border-border flex py-1 px-2 items-center cursor-pointer w-[fit-content] mr-2 mt-2 max-h-[1.625rem]",
              isValidEmail(thisEmail) ? "" : "bg-pink-200",
            )}
            onClick={() => {
              //if input text is empty, set it to the value of this email
              if (email === "") {
                setEmail(thisEmail);
              }
              analytics.track(SHARE_REMOVE_EMAIL);
              //remove this email from emails array by index
              setEmails(emails.filter((_, i) => i !== index));
              focusOnInput();
            }}
          >
            <span className="text-xs text-foreground font-medium">
              {thisEmail}
            </span>
            <XIcon
              className={cn(
                "w-3 h-3  ml-1",
                isValidEmail(thisEmail ? "text-muted-icon" : "text-primary"),
              )}
            />
          </div>
        ))}
        <Input
          id="email"
          ref={inputRef}
          type="email"
          className="min-w-[4rem] border-0 bg-transparent border-transparent focus:border-0 focus:ring-0 focus-visible:ring-0 focus-visible:ring-offset-0 h-7 p-0 mt-2 w-auto flex-grow text-xs text-foreground"
          data-1p-ignore
          spellCheck="false"
          value={email}
          onKeyDown={(e) => {
            if (e.key === "Enter" || e.key === " ") {
              if (email.length > 0) {
                handleAddEmail(email);
              } else {
                onSubmit(message);
              }
            }
            if (e.key === "Backspace" && email === "") {
              analytics.track(SHARE_REMOVE_EMAIL);
              const lastEmail = emails[emails.length - 1];
              setEmails(emails.slice(0, -1));
              setEmail(lastEmail);
            }
          }}
          onChange={(e) => {
            const value = e.target.value;
            const lastChar = value.charAt(value.length - 1);
            if (lastChar === ",") {
              handleAddEmail(email);
            } else {
              setEmail(value);
            }
          }}
          onPaste={(e) => {
            e.preventDefault();
            const pastedText = e.clipboardData.getData("text");
            const cleanedEmail = cleanupEmail(pastedText);

            // Only auto-add if it looks like an email
            if (isValidEmail(cleanedEmail)) {
              handleAddEmail(pastedText);
            } else {
              setEmail(cleanedEmail);
            }
          }}
        />
      </div>
      <Textarea
        className="w-full border-input border-[1px] rounded-md p-3 h-28 text-sm text-foreground"
        placeholder="Add message (optional)"
        onChange={(e) => {
          setMessage(e.target.value);
        }}
        value={message}
      />
      <Button
        className="w-full bg-primary hover:bg-primary/90 hover:text-primary-foreground text-primary-foreground mt-6"
        onClick={(e) => {
          e.preventDefault();
          if (message.length > 0) {
            analytics.track(SHARE_ADD_CUSTOM_MESSAGE);
          }
          onSubmit(message);
        }}
        disabled={loadingshareDocument || !hasValidEmail}
        type="submit"
      >
        Send
      </Button>
    </div>
  );
};
