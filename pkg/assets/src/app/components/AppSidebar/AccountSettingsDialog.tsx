import React, { useState, useRef, useEffect } from "react";
import { AlertTriangle, Edit2 } from "lucide-react";
import {
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Switch } from "@/components/ui/switch";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn, getInitials } from "@/lib/utils";
import { useQuery, useMutation, useLazyQuery } from "@apollo/client";
import { analytics } from "@/lib/segment";
import {
  USER_SETTINGS_CLICK_AVATAR,
  USER_SETTINGS_UPDATE_AVATAR,
  USER_SETTINGS_UPDATE_DISPLAY_NAME,
  USER_SETTINGS_UPDATE_NAME,
} from "@/lib/events";
import {
  GetMe,
  getMyPreference,
  updateMe,
  updateMyPreference,
} from "@/queries/user";
import { useCurrentUserContext } from "@/contexts/CurrentUserContext";
import { useErrorToast } from "@/hooks/useErrorToast";
import { useToast } from "../ui/use-toast";

type AccountSettingsDialogContainerProps = {
  showOpen?: boolean;
  onClose: () => void;
};
type AccountSettingsDialogProps = {
  currentUser: User | null;
  onSave: (fullName: string, displayName: string) => void;
  onUploadFile: (file: File) => Promise<{ result: any; error: any }>;
} & AccountSettingsDialogContainerProps;

function AccountSettingsDialog({
  showOpen,
  onClose,
  onSave,
  onUploadFile,
  currentUser,
}: AccountSettingsDialogProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [enableActivityNotifications, setEnableActivityNotifications] =
    useState<boolean | undefined>(undefined);

  const [fullName, setFullName] = useState(currentUser?.name);
  const [displayName, setDisplayName] = useState(currentUser?.displayName);
  const [displayNameErr, setDisplayNameErr] = useState<string | null>();
  const [fullNameErr, setFullNameErr] = useState<string | null>();
  const showErrorToast = useErrorToast();
  const { toast } = useToast();

  const [updateMyPreferenceMutation] = useMutation(updateMyPreference);
  const { data: myPreferenceData, loading } = useQuery(getMyPreference, {
    skip: !showOpen,
  });

  useEffect(() => {
    if (myPreferenceData) {
      setEnableActivityNotifications(
        myPreferenceData?.myPreference?.enableActivityNotifications,
      );
    }
  }, [myPreferenceData]);

  const handleAvatarClick = () => {
    analytics.track(USER_SETTINGS_CLICK_AVATAR);
    if (fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  const handleFileChange = async (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    if (!event.target.files) return;
    const file = event.target.files[0];
    if (!file) return;

    const { error } = await onUploadFile(file);

    if (error) {
      console.log("showing error toast");
      showErrorToast("Failed to upload avatar.");
      onClose();
    }
  };

  const handleSave = async () => {
    if (
      !currentUser ||
      (displayName == currentUser.displayName && fullName == currentUser.name)
    ) {
      return;
    }
    if (displayNameErr || fullNameErr || !displayName || !fullName) {
      return;
    }

    onSave(fullName, displayName);
  };

  const handleDisplayNameChange = (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    setDisplayName(event.target.value);
    analytics.track(USER_SETTINGS_UPDATE_DISPLAY_NAME);
    if (event.target.value == "") {
      setDisplayNameErr("Display name cannot be empty.");
    } else {
      setDisplayNameErr(null);
    }
  };

  const handleFullNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setFullName(event.target.value);
    analytics.track(USER_SETTINGS_UPDATE_NAME);
    if (event.target.value == "") {
      setFullNameErr("Full name cannot be empty.");
    } else {
      setFullNameErr(null);
    }
  };

  const handleEnableActivityNotificationsChange = (checked: boolean) => {
    setEnableActivityNotifications(checked);
    savePreference({
      enableActivityNotifications: checked,
    });
    toast({
      title: "Activity notifications updated",
    });
  };

  const savePreference = ({
    enableActivityNotifications,
  }: {
    enableActivityNotifications?: boolean;
  }) => {
    updateMyPreferenceMutation({
      variables: {
        input: {
          enableActivityNotifications: enableActivityNotifications,
        },
      },
      refetchQueries: ["getMyPreference"],
    });
  };

  return (
    <DialogContent className="min-w-[425px]">
      <DialogHeader>
        <DialogTitle>Account</DialogTitle>
        <DialogDescription>Manage your Pointy profile</DialogDescription>
      </DialogHeader>
      <div className="self-stretch flex-col justify-start items-start gap-4 flex">
        <div className="self-stretch justify-start items-center gap-4 inline-flex">
          <div
            className="relative w-16 h-16 cursor-pointer"
            onClick={handleAvatarClick}
          >
            <Avatar className="w-16 h-16">
              <AvatarImage src={currentUser?.picture || undefined} />
              <AvatarFallback className="text-background bg-foreground">
                {getInitials(currentUser?.name || "")}
              </AvatarFallback>
            </Avatar>

            <div className="absolute inset-0 bg-modal-background flex items-center justify-center opacity-0 hover:opacity-100">
              <Edit2 className="w-6 h-6 z-20 text-background items-center justify-center" />
            </div>
            <input
              type="file"
              ref={fileInputRef}
              className="hidden"
              onChange={handleFileChange}
              accept="image/jpeg,image/png"
            />
          </div>
          <div className="grow shrink basis-0 flex-col justify-start items-start gap-1 inline-flex">
            <Label className="mb-1 text-foreground text-sm font-medium  leading-tight">
              Display name
            </Label>
            <Input
              className={cn(
                "px-3 py-2",
                displayNameErr ? "border-red-500" : "",
              )}
              name="display_name"
              placeholder="Display name"
              value={displayName}
              onChange={handleDisplayNameChange}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  handleSave();
                }
              }}
              onBlur={handleSave}
            />
            {displayNameErr && (
              <p className="text-red-500">
                <AlertTriangle className="w-4 h-4 inline mr-1" />
                {displayNameErr}
              </p>
            )}
          </div>
        </div>
        <div className="self-stretch flex-col justify-start items-start gap-1 flex">
          <Label className="mb-1 text-foreground text-sm font-medium leading-tight">
            Full name
          </Label>
          <Input
            className={cn("px-3 py-2", fullNameErr ? "border-red-500" : "")}
            name="full_name"
            placeholder="James Joyce"
            value={fullName}
            onChange={handleFullNameChange}
            onBlur={handleSave}
          />
          {fullNameErr && (
            <p className="text-red-500">
              <AlertTriangle className="w-4 h-4 inline mr-1" />
              {fullNameErr}
            </p>
          )}
        </div>
      </div>
      <div className="self-stretch py-1 flex-col justify-center items-center flex">
        <div className="w-full h-px relative bg-white border-t border-border"></div>
      </div>
      {loading && (
        <div className="self-stretch py-1 flex-col justify-center items-center flex">
          Loading...
        </div>
      )}

      {!loading && (
        <>
          <div className="mb-1 self-stretch flex justify-start items-start gap-1">
            <div className="flex-grow">
              <h3 className="mb-1 text-foreground text-sm font-medium leading-tight">
                Send email notifications
              </h3>
              <p>
                Pointy will send email notification depending on your document
                preferences.
              </p>
            </div>
            <div className="flex items-center h-full">
              <Switch
                checked={enableActivityNotifications}
                onCheckedChange={handleEnableActivityNotificationsChange}
              />
            </div>
          </div>
        </>
      )}
      <div className="self-stretch py-1 flex-col justify-center items-center flex">
        <div className="w-full h-px relative bg-white border-t border-border"></div>
      </div>
      <div>
        <div className="mb-1 text-foreground text-sm font-medium leading-tight">
          Need help with your account?
        </div>
        <div>We’re here to help, reach out at taylor@pointy.ai</div>
      </div>
    </DialogContent>
  );
}

export default function AccountSettingsDialogContainer(
  props: AccountSettingsDialogContainerProps,
) {
  const { currentUser } = useCurrentUserContext();
  const [updateMeMutation] = useMutation(updateMe);
  const [getMeQuery] = useLazyQuery(GetMe, {
    fetchPolicy: "network-only",
  });
  const { toast } = useToast();

  const onSave = async (fullName: string, displayName: string) => {
    await updateMeMutation({
      variables: {
        input: {
          displayName,
          name: fullName,
        },
      },
    });
    toast({
      title: "Account settings saved",
    });
  };

  const onUploadFile = async (file: File) => {
    let result = null,
      error = null;

    console.log(file);
    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await fetch(`/api/v1/avatar`, {
        method: "PUT",
        body: formData,
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Failed to upload avatar.");
      }

      analytics.track(USER_SETTINGS_UPDATE_AVATAR);
      getMeQuery();
      toast({
        title: "Avatar updated",
      });
    } catch (err) {
      error = err;
      console.error("Error uploading avatar:", error);
    }

    return { result, error };
  };

  return (
    currentUser && (
      <AccountSettingsDialog
        currentUser={currentUser}
        onSave={onSave}
        onUploadFile={onUploadFile}
        {...props}
      />
    )
  );
}
