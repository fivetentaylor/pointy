"use client";

import React, { useContext, useEffect, useState } from "react";
import { User } from "@/__generated__/graphql";
import { analytics } from "@/lib/segment";
import { useQuery } from "@apollo/client";
import { GetMe, GetMessagingLimit, GetUsersInMyDomain } from "@/queries/user";

type CurrentUserState = {
  currentUser: User | null;
  messagingLimit: {
    total: number;
    used: number;
    startingAt: string;
    endingAt: string;
  } | null;
  usersInDomain: User[];
  loading: boolean;
  loadingUsersInDomain: boolean;
};

type CurrentUserContextProviderProps = {
  children: React.ReactNode;
};

const CurrentUserContext = React.createContext<CurrentUserState | undefined>(
  undefined,
);

let tokenRefreshed = false;

export const CurrentUserProvider: React.FC<CurrentUserContextProviderProps> = ({
  children,
}) => {
  const { loading: loadingMe, data: meData } = useQuery(GetMe);
  const { loading: loadingMessageLimit, data: messageLimit } =
    useQuery(GetMessagingLimit);
  const { loading: loadingUsersInDomain, data: usersInDomainData } =
    useQuery(GetUsersInMyDomain);
  const [usersInDomain, setUsersInDomain] = useState<User[]>([]);

  useEffect(() => {
    if (!tokenRefreshed) {
      fetch("/auth/refresh_token", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .then(() => {
          tokenRefreshed = true;
        })
        .catch((error) => {
          console.error("Error refreshing token:", error);
        });
    }
  }, []);

  useEffect(() => {
    const identifyUser = async () => {
      if (meData && meData.me) {
        // check if the user has been identified
        const identifiedUser = await analytics.user();
        const { email, name } = identifiedUser.traits() || {};

        if (email === meData.me.email && name === meData.me.name) {
          return;
        }
        console.log("Identifying user", meData.me);

        analytics.identify(meData.me.id, {
          email: meData.me.email,
          name: meData.me.name,
        });
      }
    };

    identifyUser();
  }, [meData]);

  useEffect(() => {
    if (usersInDomainData && usersInDomainData.usersInMyDomain) {
      setUsersInDomain(usersInDomainData.usersInMyDomain);
    }
  }, [usersInDomainData]);

  const currentUser = meData?.me || null;
  const messagingLimit = messageLimit?.getMessagingLimits;

  return (
    <CurrentUserContext.Provider
      value={{
        currentUser,
        messagingLimit,
        usersInDomain,
        loading: loadingMe || loadingMessageLimit,
        loadingUsersInDomain,
      }}
    >
      {children}
    </CurrentUserContext.Provider>
  );
};

export const useCurrentUserContext = () => {
  const context = useContext(CurrentUserContext);
  if (context === undefined) {
    throw new Error(
      "useCurrentUserContext must be used within a CurrentUserProvider",
    );
  }
  return context;
};
