import React, { useEffect, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Checkbox } from "@/components/ui/checkbox";
import { DocumentPreference } from "@/__generated__/graphql";
import { GetDocument, UpdateDocumentPreferences } from "@/queries/document";
import { useMutation, useQuery } from "@apollo/client";
import { useParams } from "react-router-dom";
import { useDocumentContext } from "@/contexts/DocumentContext";
import { NOTIFICATIONS_SETTINGS_UPDATE } from "@/lib/events";
import { analytics } from "@/lib/segment";

type TimelineNotificationSettingsDialogProps = {
  showOpen?: boolean;
  onClose: () => void;
};

// Example of a reusable CheckboxItem component
const CheckboxItem = ({
  disabled,
  id,
  label,
  checked,
  onChange,
}: {
  disabled: boolean;
  id: string;
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
}) => (
  <div className="items-top flex space-x-2">
    <Checkbox
      id={id}
      onCheckedChange={onChange}
      checked={checked}
      disabled={disabled}
    />
    <label
      htmlFor={id}
      className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
    >
      {label}
    </label>
  </div>
);

const TimelineNotificationSettingsDialog = function ({
  showOpen,
  onClose,
}: TimelineNotificationSettingsDialogProps) {
  const [preferences, setPreferences] = useState<DocumentPreference | null>(
    null,
  );
  const { docData, updatePreferences, savingPreferences } =
    useDocumentContext();

  useEffect(() => {
    if (docData) {
      setPreferences(docData.preferences);
    }
  }, [docData]);

  const handlePreferenceChange = (preferenceKey: string, newValue: boolean) => {
    analytics.track(NOTIFICATIONS_SETTINGS_UPDATE, {
      notificationPreference: preferenceKey,
      notificationValue: newValue,
    });

    const newPreferences: DocumentPreference = {
      ...preferences!,
      [preferenceKey]: newValue,
    };

    setPreferences(newPreferences);

    updatePreferences({
      variables: {
        documentId: docData?.id,
        input: newPreferences,
      },
    });
  };

  return (
    <Dialog open={showOpen} onOpenChange={onClose}>
      <DialogContent className="min-w-[425px]">
        <DialogHeader>
          <DialogTitle>Notifications</DialogTitle>
          <DialogDescription>
            Select how you&#39;d like to be updated about this document.
          </DialogDescription>
        </DialogHeader>
        {preferences && (
          <div className="flex flex-col space-y-4 pt-2">
            <CheckboxItem
              id="enableFirstOpenNotifications"
              label="First time someone opens the document"
              disabled={savingPreferences}
              checked={preferences.enableFirstOpenNotifications}
              onChange={(checked) =>
                handlePreferenceChange("enableFirstOpenNotifications", checked)
              }
            />
            <CheckboxItem
              id="enableMentionNotifications"
              label="Anytime someone mentions you"
              disabled={savingPreferences}
              checked={preferences.enableMentionNotifications}
              onChange={(checked) =>
                handlePreferenceChange("enableMentionNotifications", checked)
              }
            />
            <CheckboxItem
              id="enableAllCommentNotifications"
              label="All comments left on the document"
              disabled={savingPreferences}
              checked={preferences.enableAllCommentNotifications}
              onChange={(checked) =>
                handlePreferenceChange("enableAllCommentNotifications", checked)
              }
            />
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
};

export default TimelineNotificationSettingsDialog;
