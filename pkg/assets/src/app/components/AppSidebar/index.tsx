import React from "react";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { DocumentGroups } from "./DocumentGroups";
import User from "./User";
import { Button } from "../ui/button";
import { EditIcon, FolderPlusIcon } from "lucide-react";
import { useWsDisconnect } from "@/hooks/useWsDisconnect";
import { useDocumentContext } from "@/contexts/DocumentContext";
import { DragProvider } from "@/contexts/DragContext";

export function AppSidebar() {
  return (
    <Sidebar collapsible="icon" className="pt-6 bg-sidebar">
      <SidebarHeader>
        <SidebarTrigger className="mb-3 ml-2" />
        <HeaderNewButtons />
      </SidebarHeader>
      <SidebarContent>
        <DragProvider>
          <DocumentGroups />
        </DragProvider>
      </SidebarContent>
      <SidebarFooter>
        <User />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}

const HeaderNewButtons = () => {
  const { createNewDocument, docLoading } = useDocumentContext();
  const { createNewFolder, loadingCreateFolder } = useDocumentContext();
  const { isDisconnected } = useWsDisconnect();

  return (
    <div className="flex grid grid-cols-4 gap-1">
      <Button
        variant="outline"
        className="w-full group-data-[collapsible=icon]:border-none col-span-3"
        disabled={docLoading || isDisconnected}
        onClick={createNewDocument}
      >
        <EditIcon className="w-4 h-4 min-w-4" />
        <span className="ml-2 truncate group-data-[collapsible=icon]:hidden">
          New Draft
        </span>
      </Button>
      <Button
        variant="outline"
        className="w-full group-data-[collapsible=icon]:hidden col-span-1"
        disabled={loadingCreateFolder || isDisconnected}
        onClick={createNewFolder}
      >
        <FolderPlusIcon className="w-4 h-4 min-w-4" />
      </Button>
    </div>
  );
};
