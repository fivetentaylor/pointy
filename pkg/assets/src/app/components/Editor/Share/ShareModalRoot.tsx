import React, { useMemo, useState } from "react";
import { Button } from "@/components/ui/button";

import { APP_HOST, WEB_HOST } from "@/lib/urls";
import {
  ChevronRightIcon,
  DownloadIcon,
  FilesIcon,
  GlobeIcon,
  PlusIcon,
} from "lucide-react";
import { Switch } from "@/components/ui/switch";
import { DocEditors } from ".";
import { EmailInput } from "./EmailInput";
import { WithTooltip } from "@/components/ui/FloatingTooltip";

export const ShareModalRoot = function ({
  email,
  isDocPublic,
  onEmailChange,
  onAccessChange,
  docId,
  isDocumentUpdating,
  sharedLinks,
  editors,
  onEnter,
  onEditAccess,
}: {
  email: string;
  docId: string;
  isDocPublic: boolean;
  isDocumentUpdating: boolean;
  sharedLinks: SharedDocumentLinksQuery["sharedLinks"] | undefined;
  editors: DocEditors | undefined;
  onAccessChange: (isPublic: boolean) => void;
  onEmailChange: (email: string) => void;
  onEnter: () => void;
  onEditAccess: () => void;
}) {
  const [copyClicked, setCopyClicked] = useState(false);
  const showAccessList = useMemo(() => {
    return (
      (editors && editors.length > 0) || (sharedLinks && sharedLinks.length > 0)
    );
  }, [editors, sharedLinks]);

  return (
    <div className="flex flex-col pt-4">
      <EmailInput
        email={email}
        onEmailChange={onEmailChange}
        onEnter={onEnter}
      />

      {/* Access List */}
      {showAccessList ? (
        <div className="flex mt-2">
          <div className="flex flex-grow items-center">
            {editors?.map((editor, index) => (
              <div key={`editor-${index}`} className="flex">
                {editor.picture ? (
                  <img
                    src={editor.picture}
                    alt={editor.name}
                    title={`${editor.name} <${editor.email}>`}
                    className="w-6 h-6 rounded-full"
                  />
                ) : (
                  <div
                    className="w-6 h-6 rounded-full bg-primary text-primary-foreground flex items-center justify-center cursor-default"
                    title={editor.email}
                  >
                    {editor.name.charAt(0).toUpperCase()}
                  </div>
                )}
              </div>
            ))}
            {sharedLinks?.map((link, index) => (
              <div key={`link-${index}`} className="flex">
                <div
                  className="w-6 h-6 rounded-full bg-primary text-primary-foreground flex items-center justify-center cursor-default"
                  title={link.inviteeEmail}
                >
                  {link.inviteeEmail.charAt(0).toUpperCase()}
                </div>
              </div>
            ))}

            <Button
              className="rounded-full bg-background text-muted-icon flex items-center justify-center cursor-pointer p-0 m-0 border-dashed border-[1px] border-border ml-2 hover:bg-background/90 w-[calc(1.5rem+2px)] h-[calc(1.5rem+2px)]"
              title="Add more writers"
              onClick={onEnter}
            >
              <PlusIcon className="w-4 h-4" />
            </Button>
          </div>
          <div className="items-end">
            <Button
              variant="link"
              className="text-primary"
              onClick={onEditAccess}
            >
              Edit access
            </Button>
          </div>
        </div>
      ) : (
        <div className="mt-2" />
      )}

      <hr className="text-border mt-2 mb-4" />

      <div className="flex flex-col">
        <h4 className="font-semibold text-base">Collaboration link</h4>
        <div className="flex mt-[1.125rem] items-center">
          <GlobeIcon className="w-4 h-4 text-foreground mr-[0.6875rem]" />
          <span className="text-md leading-[1.25rem] text-foreground flex-grow">
            Anyone with the link can view
          </span>
          <Switch
            className="flex-end"
            disabled={isDocumentUpdating}
            checked={isDocPublic}
            onCheckedChange={onAccessChange}
          />
        </div>
      </div>

      <div className="w-full">
        <Button
          className="w-full mt-4"
          variant="outline"
          title={`Copy ${APP_HOST}/drafts/${docId} to clipboard`}
          onClick={() => {
            navigator.clipboard.writeText(`${APP_HOST}/drafts/${docId}`);
            setCopyClicked(true);
          }}
        >
          {copyClicked ? "Link copied to clipboard!" : "Copy link"}
        </Button>
      </div>

      <hr className="text-border my-4" />

      <WithTooltip tooltipText="Coming soon!">
        <div
          className={`h-[3.125rem] flex items-center text-secondary-foreground rounded-lg
          ${/*cursor-pointer hover:bg-border/30 - replace when ready */ "cursor-not-allowed"}`}
        >
          <FilesIcon className="w-6 h-6 mr-[1.4375rem] opacity-50" />
          <span className="flex-grow text-[0.9375rem] leading-[1.46875rem] opacity-50">
            Share a copy
          </span>
          <ChevronRightIcon className="w-6 h-6 opacity-50" />
        </div>
      </WithTooltip>

      <WithTooltip tooltipText="Coming soon!">
        <div
          className={`h-[3.125rem] flex items-center text-secondary-foreground rounded-lg
          ${/*cursor-pointer hover:bg-border/30 - replace when ready */ "cursor-not-allowed"}`}
        >
          <DownloadIcon className="w-6 h-6 mr-[1.4375rem] opacity-50" />
          <span className="flex-grow text-[0.9375rem] leading-[1.46875rem] opacity-50">
            Download a copy
          </span>
          <ChevronRightIcon className="w-6 h-6 opacity-50" />
        </div>
      </WithTooltip>
    </div>
  );
};
