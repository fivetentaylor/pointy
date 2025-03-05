import React, { useEffect, useState, useCallback } from "react";
import { Document } from "@/__generated__/graphql";
import { useNavigate } from "react-router-dom";
import { ChevronDown, FileIcon, FilesIcon, SearchIcon } from "lucide-react";
import { useWsDisconnect } from "@/hooks/useWsDisconnect";
import { useDocumentContext } from "@/contexts/DocumentContext";
import { useToast } from "../ui/use-toast";
import { analytics } from "@/lib/segment";
import { DRAFTS_OPEN } from "@/lib/events";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
} from "@/components/ui/sidebar";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { ScrollArea } from "../ui/scroll-area";
import { DocumentMenuItem } from "./Items";
import { useDebounce } from "use-debounce";
import { useQuery } from "@apollo/client";
import { SearchDocuments } from "@/queries/document";

type DocumentGroupsProps = {};

export function DocumentGroups({}: DocumentGroupsProps) {
  const navigate = useNavigate();
  const { isDisconnected } = useWsDisconnect();
  const {
    draftId,
    baseDocuments,
    sharedDocuments,
    deleteDocument,
    renameFolder,
    moveDocument,
  } = useDocumentContext();
  const { toast } = useToast();

  const myDocuments = (baseDocuments || []) as Document[];

  const handleNavigate = (docId: string) => {
    analytics.track(DRAFTS_OPEN);
    if (isDisconnected) {
      toast({
        title: "You are disconnected. Please reconnect to continue.",
      });
      return;
    }
    const params = new URLSearchParams(location.search);
    const sbParam = params.get("sb");
    let target = `/drafts/${docId}`;
    if (sbParam) {
      target += `?sb=${sbParam}`;
    }
    navigate(target);
  };

  const handleFolderRename = (id: string, name: string) => {
    renameFolder(id, name);
  };

  const onMoveDocumentToFolder = (
    docId: string,
    newFolderId: string | null | undefined,
    oldFolderId: string | null | undefined,
  ) => {
    console.log(
      `Moving document ${docId} from folder ${oldFolderId} to folder ${newFolderId}`,
    );
    moveDocument(docId, newFolderId, oldFolderId);
  };

  // Handle Meta + [ for navigation
  const handleKeydown = useCallback(
    (event: KeyboardEvent) => {
      if (event.metaKey && event.key === "[") {
        event.preventDefault(); // Prevent default browser behavior
        const currentDocIndex = myDocuments.findIndex(
          (doc) => doc.id === draftId,
        );
        const nextDocIndex =
          currentDocIndex + 1 < myDocuments.length ? currentDocIndex + 1 : 0; // Loop back to the first document
        const nextDoc = myDocuments[nextDocIndex];

        if (nextDoc) {
          handleNavigate(nextDoc.id);
        }
      }

      if (event.metaKey && event.key === "]") {
        event.preventDefault(); // Prevent default browser behavior
        const currentDocIndex = myDocuments.findIndex(
          (doc) => doc.id === draftId,
        );
        const nextDocIndex =
          currentDocIndex - 1 >= 0
            ? currentDocIndex - 1
            : myDocuments.length - 1; // Loop back to the last document
        const nextDoc = myDocuments[nextDocIndex];

        if (nextDoc) {
          handleNavigate(nextDoc.id);
        }
      }
    },
    [myDocuments, handleNavigate],
  );

  useEffect(() => {
    document.addEventListener("keydown", handleKeydown);
    return () => {
      document.removeEventListener("keydown", handleKeydown);
    };
  }, [handleKeydown]);

  return (
    <>
      <Collapsible
        defaultOpen
        className="group/collapsible  group-data-[collapsible=icon]:hidden"
      >
        <SidebarGroup className="py-0">
          <SidebarGroupLabel asChild>
            <CollapsibleTrigger className="flex w-full items-center justify-between">
              Your drafts
              <ChevronDown className="h-4 w-4 transition-transform group-data-[state=open]/collapsible:rotate-180" />
            </CollapsibleTrigger>
          </SidebarGroupLabel>
          <CollapsibleContent>
            <SidebarGroupContent>
              <SidebarMenu>
                {myDocuments.map((doc) => (
                  <DocumentMenuItem
                    key={doc.id}
                    doc={doc}
                    onNavigate={handleNavigate}
                    onDelete={deleteDocument}
                    onRename={handleFolderRename}
                    isDisconnected={isDisconnected}
                    canDelete={true}
                    onMoveDocumentToFolder={onMoveDocumentToFolder}
                  />
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </CollapsibleContent>
        </SidebarGroup>
      </Collapsible>

      <Collapsible
        defaultOpen
        className="group/collapsible  group-data-[collapsible=icon]:hidden"
      >
        <SidebarGroup className="py-0 mt-[-0.5rem]">
          <SidebarGroupLabel asChild>
            <CollapsibleTrigger className="flex w-full items-center justify-between">
              Shared with you
              <ChevronDown className="h-4 w-4 transition-transform group-data-[state=open]/collapsible:rotate-180" />
            </CollapsibleTrigger>
          </SidebarGroupLabel>
          <CollapsibleContent>
            <SidebarGroupContent>
              <SidebarMenu>
                {sharedDocuments.map((doc) => (
                  <DocumentMenuItem
                    key={doc.id}
                    doc={doc as Document}
                    onNavigate={handleNavigate}
                    onDelete={deleteDocument}
                    onRename={handleFolderRename}
                    isDisconnected={isDisconnected}
                    canDelete={false}
                    onMoveDocumentToFolder={onMoveDocumentToFolder} // Pass down
                  />
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </CollapsibleContent>
        </SidebarGroup>
      </Collapsible>

      <IconDocumentMenu onNavigate={handleNavigate} />
    </>
  );
}

function IconDocumentMenu({
  onNavigate,
}: {
  onNavigate: (id: string) => void;
}) {
  const [searchTerm, setSearchTerm] = useState("");
  const [isOpen, setIsOpen] = useState(false);
  const [debouncedSearchTerm] = useDebounce(searchTerm, 100);

  const { data, loading, error, refetch } = useQuery(SearchDocuments, {
    variables: {
      query: debouncedSearchTerm,
      offset: 0,
      limit: 10,
    },
    skip: !debouncedSearchTerm && !isOpen,
  });

  useEffect(() => {
    if (isOpen) {
      refetch({
        query: debouncedSearchTerm,
        offset: 0,
        limit: 10,
      });
    }
  }, [isOpen, refetch, debouncedSearchTerm]);

  const documents = data?.searchDocuments?.edges || [];

  return (
    <DropdownMenu onOpenChange={setIsOpen}>
      <DropdownMenuTrigger
        asChild
        className="invisible group-data-[collapsible=icon]:visible focus:ring-0 focus:outline-none focus:ring-offset-0"
      >
        <Button variant="icon" className="mx-2">
          <FilesIcon className="w-4 h-4 min-w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-64" align="start" side="right">
        <div className="p-2 relative">
          <SearchIcon className="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
          <Input
            className="w-full py-2 h-8 pl-8"
            type="search"
            placeholder="Search files..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            onKeyDown={(e) => {
              if (e.key !== "Escape") {
                e.stopPropagation();
              }
            }}
          />
        </div>
        <DropdownMenuGroup>
          <ScrollArea className="h-64">
            <div className="my-2 ml-2 text-muted-foreground text-xs font-medium leading-none">
              {loading ? "Searching..." : "Most recent"}
            </div>
            {error && (
              <div className="text-red-500 text-xs ml-2">
                Error loading files
              </div>
            )}
            {documents.map(({ node: doc }) => (
              <DropdownMenuItem
                key={doc.id}
                onSelect={() => onNavigate(doc.id)}
              >
                <span className="w-4 mr-2">
                  <FileIcon className="h-4 w-4" />
                </span>
                <span className="truncate">{doc.title}</span>
              </DropdownMenuItem>
            ))}
            {!loading && documents.length === 0 && debouncedSearchTerm && (
              <div className="text-xs text-muted-foreground ml-2">
                No results found
              </div>
            )}
          </ScrollArea>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
