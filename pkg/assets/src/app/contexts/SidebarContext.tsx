import { useQuery } from "@apollo/client";
import React, { useEffect, createContext, useContext, useState } from "react";
import { useParams } from "react-router-dom";
import { useDocumentContext } from "./DocumentContext";

type SidebarContextState = ReturnType<typeof useSetupSidebar>;

export type SidebarMode = "timeline" | "chat";

type SidebarContextProviderProps = {
  children: React.ReactNode;
};

const SidebarContext = createContext<SidebarContextState | undefined>(
  undefined,
);

export const useSetupSidebar = () => {
  const [sidebarMode, _setSidebarMode] = useState<SidebarMode>();
  const [sidebarLoading, setSidebarLoading] = useState<boolean>(true);
  const [showSubscriptions, setShowSubscriptions] = useState<boolean>(false);
  const { docData, docLoading } = useDocumentContext();

  // on doc loading, check editor count. If more than one editor, set to timeline. Else chat.
  useEffect(() => {
    if (docLoading) {
      return;
    }

    if (!docLoading && sidebarLoading) {
      setSidebarLoading(false);
    }

    const params = new URLSearchParams(location.search);
    const sbParam = params.get("sb");
    const editors = docData?.editors;

    if (sbParam === "t") {
      _setSidebarMode("timeline");
    } else if (sbParam === "c") {
      _setSidebarMode("chat");
    } else if (editors && editors.length > 1) {
      _setSidebarMode("timeline");
    } else {
      _setSidebarMode("chat");
    }
  }, [_setSidebarMode, docData, docLoading]);

  const setSidebarMode = (mode: SidebarMode) => {
    const params = new URLSearchParams(window.location.search);
    params.set("sb", mode === "timeline" ? "t" : "c");

    window.history.pushState(
      {},
      "",
      `${window.location.pathname}?${params.toString()}`,
    );

    _setSidebarMode(mode);
  };

  return {
    setSidebarMode,
    sidebarLoading,
    sidebarMode,
    showSubscriptions,
    setShowSubscriptions,
  };
};

export const SidebarContextProvider = function ({
  children,
}: SidebarContextProviderProps) {
  const state = useSetupSidebar();

  return (
    <SidebarContext.Provider value={state}>{children}</SidebarContext.Provider>
  );
};

export const useSidebarContext = () => {
  const context = useContext(SidebarContext);
  if (context === undefined) {
    throw new Error(
      "useSidebarContext must be used within a SidebarContextProvider",
    );
  }
  return context;
};
