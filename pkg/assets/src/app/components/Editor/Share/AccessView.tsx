import React from "react";
import { ChevronDown } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { DocEditors } from ".";
import { EmailInput } from "./EmailInput";

/* eslint-disable @next/next/no-img-element */
export const AccessView = function ({
  editors,
  email,
  sharedLinks,
  onEmailChange,
  onEnter,
  onRemoveEditor,
  onRemoveLink,
}: {
  editors: DocEditors | undefined;
  email: string;
  onEmailChange: (email: string) => void;
  onEnter: () => void;
  sharedLinks: SharedDocumentLinksQuery["sharedLinks"] | undefined;
  onRemoveLink: (link: SharedDocumentLinksQuery["sharedLinks"][0]) => void;
  onRemoveEditor: (editor: DocEditors[0]) => void;
}) {
  return (
    <div className="flex flex-col">
      <h3 className="text-foreground text-lg font-semibold">Who has access</h3>
      <div className="mt-4 mb-2">
        <EmailInput
          email={email}
          onEmailChange={onEmailChange}
          onEnter={onEnter}
        />
      </div>
      <div className="w-full">
        {editors?.map((editor, index) => (
          <div key={index} className="w-full mt-4 flex items-center">
            <div className="flex items-center flex-grow">
              {editor.picture ? (
                <img
                  src={editor.picture}
                  alt={editor.name}
                  title={`${editor.name} <${editor.email}>`}
                  className="w-8 h-8 rounded-full"
                />
              ) : (
                <div
                  className="w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center cursor-default"
                  title={editor.email}
                >
                  {editor.name.charAt(0).toUpperCase()}
                </div>
              )}
              <div className="flex flex-col ml-4">
                <span className="text-foreground text-sm">{editor.email}</span>
                <span className="text-foreground text-sm">{editor.name}</span>
              </div>
            </div>

            <ACLDropdown
              canTransferOwnership
              onClickRemove={() => onRemoveEditor(editor)}
            />
          </div>
        ))}
        {sharedLinks?.map((link, index) => (
          <div key={index} className="w-full mt-4 flex items-center">
            <div className="flex items-center flex-grow">
              <div
                className="w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center cursor-default"
                title={link.inviteeEmail}
              >
                {link.inviteeEmail.charAt(0).toUpperCase()}
              </div>
              <div className="flex flex-col ml-4">
                <span className="text-foreground text-sm">
                  {link.inviteeEmail}
                </span>
                <span className="text-foreground text-sm">
                  pending invitation
                </span>
              </div>
            </div>

            <ACLDropdown
              onClickRemove={() => onRemoveLink(link)}
              canTransferOwnership={false}
            />
          </div>
        ))}
      </div>
    </div>
  );
};

const ACLDropdown = function ({
  onClickRemove,
}: {
  canTransferOwnership: boolean;
  onClickRemove: () => void;
}) {
  return (
    <div>
      <DropdownMenu>
        <DropdownMenuTrigger>
          <div className="flex items-center cursor-pointer">
            <span>Writer</span>
            <ChevronDown className="ml-2 w-4 h-4 text-muted-icon" />
          </div>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem className="cursor-pointer" onClick={onClickRemove}>
            Remove Access
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};
