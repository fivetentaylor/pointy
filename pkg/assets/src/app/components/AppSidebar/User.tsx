import React, { useState } from "react";
import { User } from "@/__generated__/graphql";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import {
  ExternalLinkIcon,
  LogOut,
  Settings2Icon,
  ReceiptTextIcon,
  ArrowRightIcon,
} from "lucide-react";
import { cn, getInitials } from "@/lib/utils";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
} from "@/components/ui/dropdown-menu";
import {
  PROFILE_CLICK_AVATAR,
  PROFILE_CLICK_COLOR_THEME,
  PROFILE_OPEN_SETTINGS,
  PAYMENTS_CLICK_CHECKOUT,
} from "@/lib/events";
import { analytics } from "@/lib/segment";
import AccountSettingsDialog from "./AccountSettingsDialog";
import { Dialog } from "../ui/dialog";
import { Skeleton } from "../ui/skeleton";
import { useWsDisconnect } from "@/hooks/useWsDisconnect";
import {
  colorThemeService,
  ThemePreference,
} from "@/lib/service/ColorThemeService";
import { useCurrentUserContext } from "@/contexts/CurrentUserContext";
import { Button } from "../ui/button";
import SubscriptionDialog from "./SubscribeDialog";
import { BillingPortalSession } from "@/queries/payments";
import { useMutation } from "@apollo/client";
import { Spinner } from "../ui/spinner";
import { TipBox } from "../ui/TipBox";
import posthog from "posthog-js";
import { useSignals } from "@preact/signals-react/runtime";
import { WEB_HOST } from "@/lib/urls";

type UserProps = {
  me: User | null;
  loading: boolean;
  onClickTheme: (theme: "light" | "dark" | "system") => void;
  onClickLogout: () => void;
  theme: string;
};

const UserComponent = ({
  me,
  loading,
  onClickTheme,
  onClickLogout,
  theme,
}: UserProps) => {
  const [profileSettingsOpen, setProfileSettingsOpen] = useState(false);
  const [subscribeDialogOpen, setSubscribeDialogOpen] = useState(false);
  const { isDisconnected } = useWsDisconnect();

  const [billingPortalSession, { loading: loadingBillingPortal }] =
    useMutation(BillingPortalSession);

  if (loading || !me) {
    return (
      <div className="my-4 mx-4 flex items-center">
        <Avatar className="w-6 h-6">
          <AvatarFallback className="text-background bg-elevated"></AvatarFallback>
        </Avatar>
        <Skeleton className="ml-2 h-6 w-20 bg-elevated" />
      </div>
    );
  }

  const handleBillingPortal = async (e: React.SyntheticEvent) => {
    e.preventDefault();
    if (me.subscriptionStatus !== "active") {
      setSubscribeDialogOpen(true);
      analytics.track(PAYMENTS_CLICK_CHECKOUT);
      return;
    }

    const { data } = await billingPortalSession();
    if (data?.billingPortalSession) {
      window.location.href = data.billingPortalSession.url;
    }
  };

  return (
    <div className="relative">
      <Dialog open={profileSettingsOpen} onOpenChange={setProfileSettingsOpen}>
        <AccountSettingsDialog
          showOpen={profileSettingsOpen}
          onClose={() => setProfileSettingsOpen(false)}
        />
      </Dialog>

      <Dialog open={subscribeDialogOpen} onOpenChange={setSubscribeDialogOpen}>
        <SubscriptionDialog
          showOpen={subscribeDialogOpen}
          onClose={() => setSubscribeDialogOpen(false)}
          onSubscribe={() => analytics.track(PAYMENTS_CLICK_CHECKOUT)}
        />
      </Dialog>

      {!posthog.isFeatureEnabled("stripe") && (
        <TipBox className="w-[calc(100%-1.5rem)] py-1 group-data-[collapsible=icon]:hidden pl-8 truncate">
          {me.subscriptionStatus === "active" && (
            <>
              <div>Professional plan</div>
            </>
          )}
          {me.subscriptionStatus !== "active" && (
            <>
              <span>Free plan</span>
              <Button
                variant="highlight"
                className="p-0 pl-1 text-xs text-primary m-0 mt-[-1rem] h-4"
                onClick={() => {
                  setSubscribeDialogOpen(true);
                  analytics.track(PAYMENTS_CLICK_CHECKOUT);
                }}
              >
                Upgrade
              </Button>
            </>
          )}
        </TipBox>
      )}

      <div className="flex justify-center mb-3 group-data-[collapsible=icon]:hidden ml-[-0.5rem] truncate">
        <a
          href={`${WEB_HOST}/winding-down`}
          target="_blank"
          rel="noreferrer"
          className="group inline-flex items-center px-3 py-2 bg-foreground text-background rounded-full hover:bg-foreground/90 transition-colors"
        >
          <span className="text-sm font-normal">Reviso is winding down</span>
          <ArrowRightIcon className="w-4 h-4 shrink-0 transition-transform group-hover:translate-x-0.5" />
        </a>
      </div>

      <DropdownMenu modal={false}>
        <DropdownMenuTrigger
          asChild
          className={cn(
            "w-full rounded",
            isDisconnected ? "cursor-not-allowed opacity-50" : "cursor-pointer",
          )}
          onClick={() => analytics.track(PROFILE_CLICK_AVATAR)}
          disabled={isDisconnected}
        >
          <div className="my-2 mx-2 flex items-center group-data-[collapsible=icon]:mx-1">
            <Avatar className="w-6 h-6">
              <AvatarImage src={me.picture || undefined} />
              {!me.picture && (
                <AvatarFallback className="text-background bg-foreground">
                  {getInitials(me.name || "")}
                </AvatarFallback>
              )}
            </Avatar>
            <div className="truncate ml-2 group-data-[collapsible=icon]:hidden">
              {me.name}
            </div>
          </div>
        </DropdownMenuTrigger>

        <DropdownMenuContent className="w-56">
          <DropdownMenuGroup>
            <DropdownMenuItem
              onClick={() => {
                analytics.track(PROFILE_OPEN_SETTINGS);
                setProfileSettingsOpen(true);
              }}
            >
              <Settings2Icon className="mr-2 h-4 w-4" />
              <span>Account settings</span>
            </DropdownMenuItem>

            {posthog.isFeatureEnabled("stripe") && (
              <DropdownMenuItem
                onClick={handleBillingPortal}
                disabled={loadingBillingPortal}
              >
                {loadingBillingPortal ? (
                  <Spinner />
                ) : (
                  <ReceiptTextIcon className="mr-2 h-4 w-4" />
                )}
                <span className="flex items-center justify-between">
                  Billing{" "}
                  <ExternalLinkIcon className="ml-2 h-4 w-4 text-muted-foreground" />
                </span>
              </DropdownMenuItem>
            )}
          </DropdownMenuGroup>
          <DropdownMenuGroup>
            <DropdownMenuRadioGroup value={theme}>
              <DropdownMenuRadioItem
                value="light"
                onClick={() => onClickTheme("light")}
              >
                Light
              </DropdownMenuRadioItem>
              <DropdownMenuRadioItem
                value="dark"
                onClick={() => onClickTheme("dark")}
              >
                Dark
              </DropdownMenuRadioItem>
              <DropdownMenuRadioItem
                value="system"
                onClick={() => onClickTheme("system")}
              >
                System
              </DropdownMenuRadioItem>
            </DropdownMenuRadioGroup>
          </DropdownMenuGroup>
          <DropdownMenuSeparator className="bg-border" />
          <DropdownMenuGroup>
            <DropdownMenuItem onClick={onClickLogout}>
              <LogOut className="mr-2 h-4 w-4" />
              <span>Log out {me?.email}</span>
            </DropdownMenuItem>
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};

const UserWrapper = () => {
  const { currentUser, loading: loadingMe } = useCurrentUserContext();
  useSignals();

  if (!loadingMe && !currentUser) {
    window.location.href = "/login";
  }
  const handleLogout = () => {
    window.sessionStorage.clear();
    window.location.href = "/logout";
  };

  return (
    <UserComponent
      me={currentUser}
      loading={loadingMe}
      theme={colorThemeService.themePreference.value}
      onClickTheme={(theme: ThemePreference) => {
        analytics.track(PROFILE_CLICK_COLOR_THEME, { theme });
        colorThemeService.setPreferredTheme(theme);
      }}
      onClickLogout={handleLogout}
    />
  );
};

export default UserWrapper;
