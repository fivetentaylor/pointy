import React, {
  useLayoutEffect,
  useState,
  useRef,
  useMemo,
  useEffect,
} from "react";
import {
  DeleteIcon,
  Edit2Icon,
  EllipsisIcon,
  FileIcon,
  FolderIcon,
  FolderOpenIcon,
  LinkIcon,
} from "lucide-react";
import { useToast } from "../ui/use-toast";
import {
  SidebarMenuButton,
  SidebarMenuAction,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { DOCUMENT_LIMIT, useDocumentContext } from "@/contexts/DocumentContext";
import { DeleteDraftAlert } from "../ui/DeleteDraftAlert";
import { cn } from "@/lib/utils";
import { GetFolderDocuments } from "@/queries/document";
import { useQuery } from "@apollo/client";
import { DeleteFolderAlert } from "../ui/DeleteFolderAlert";
import { useDragContext } from "@/contexts/DragContext";
import { useIsMobile } from "@/hooks/use-mobile";

type MoveFunction = (
  docId: string,
  newFolderId: string | null | undefined,
  oldFolderId: string | null | undefined,
) => void;
type DeleteFunction = (
  docId: string,
  folderId?: string | null,
  deleteChildren?: boolean,
) => void;

type DocumentMenuItemProps = {
  doc: DocumentFieldsFragment;
  onNavigate: (id: string) => void;
  onDelete: DeleteFunction;
  onRename: (id: string, title: string) => void;
  isDisconnected: boolean;
  canDelete: boolean;
  onMoveDocumentToFolder: MoveFunction;
};

export function DocumentMenuItem({
  doc,
  onNavigate,
  onDelete,
  onRename,
  isDisconnected,
  canDelete,
  onMoveDocumentToFolder,
}: DocumentMenuItemProps) {
  const { docData: currentDocument } = useDocumentContext();

  return (
    <>
      {doc.isFolder ? (
        <Folder
          doc={doc}
          currentDocument={currentDocument}
          onRename={onRename}
          isDisconnected={isDisconnected}
          canDelete={canDelete}
          onMoveDocumentToFolder={onMoveDocumentToFolder}
          onDelete={onDelete}
          onNavigate={onNavigate}
        />
      ) : (
        <Doc
          doc={doc}
          currentDocument={currentDocument}
          onNavigate={onNavigate}
          isDisconnected={isDisconnected}
          canDelete={canDelete}
          onDelete={onDelete}
          onMoveDocumentToFolder={onMoveDocumentToFolder}
        />
      )}
    </>
  );
}

type FolderProps = {
  doc: DocumentFieldsFragment;
  currentDocument: DocumentFieldsFragment | null;
  onRename: (id: string, title: string) => void;
  isDisconnected: boolean;
  canDelete: boolean;
  onMoveDocumentToFolder: MoveFunction;
  onNavigate: (id: string) => void;
  onDelete: DeleteFunction;
};

export function Folder({
  doc,
  currentDocument,
  onRename,
  isDisconnected,
  canDelete,
  onMoveDocumentToFolder,
  onDelete,
  onNavigate,
}: FolderProps) {
  const [isOpen, setIsOpen] = useState<boolean | null>(null);
  const [isRenaming, setIsRenaming] = useState(false);
  const [renameValue, setRenameValue] = useState(doc.title);
  const [showDeleteDialog, _setShowDeleteDialog] = useState(false);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [isDraggedOver, setIsDraggedOver] = useState(false);
  const initialOpenStateRef = useRef<boolean | null>(null);
  const dragTimeoutRef = useRef<number | null>(null);
  const folderRef = useRef<HTMLDivElement>(null);
  const childrenRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // if the current document is null, we shouldn't mess with open state
    if (!currentDocument && isOpen !== null) {
      return;
    }

    const newOpenState = doc.id === currentDocument?.folderID;
    setIsOpen(newOpenState);
    if (initialOpenStateRef.current === null) {
      initialOpenStateRef.current = newOpenState;
    }
  }, [doc, currentDocument]);

  const resetDragState = () => {
    setIsDraggedOver(false);
    setIsOpen(initialOpenStateRef.current);
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();

    // Clear any existing timeout
    if (dragTimeoutRef.current !== null) {
      window.clearTimeout(dragTimeoutRef.current);
      dragTimeoutRef.current = null;
    }

    if (!isDraggedOver) {
      setIsDraggedOver(true);
      setIsOpen(true);
    }
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();

    // Check if we're still within the folder's children
    const relatedTarget = e.relatedTarget as Node;
    if (folderRef.current?.contains(relatedTarget)) {
      // If we're entering the children container, don't do anything
      if (childrenRef.current?.contains(relatedTarget)) {
        return;
      }
      return;
    }

    // Only process drag leave after a short delay to prevent flicker
    dragTimeoutRef.current = window.setTimeout(() => {
      resetDragState();
      dragTimeoutRef.current = null;
    }, 50);
  };

  const handleChildrenDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();

    const relatedTarget = e.relatedTarget as Node;
    // If we're moving back to the folder item or staying within children, don't reset
    if (folderRef.current?.contains(relatedTarget)) {
      return;
    }

    // Only process drag leave after a short delay to prevent flicker
    dragTimeoutRef.current = window.setTimeout(() => {
      resetDragState();
      dragTimeoutRef.current = null;
    }, 50);
  };

  // Clean up timeout on unmount
  useEffect(() => {
    return () => {
      if (dragTimeoutRef.current !== null) {
        window.clearTimeout(dragTimeoutRef.current);
      }
    };
  }, []);

  // Make sure that the dialog is closed before updating the state of the delete dialog
  const setShowDeleteDialog = (value: boolean) => {
    setDropdownOpen(false);
    _setShowDeleteDialog(value);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();

    // Clear any existing timeout
    if (dragTimeoutRef.current !== null) {
      window.clearTimeout(dragTimeoutRef.current);
      dragTimeoutRef.current = null;
    }

    const draggedDocIdJson = e.dataTransfer.getData("text/plain");
    const input = JSON.parse(draggedDocIdJson);
    const draggedDocId = input.id;
    const draggedFolderId = input.folderID;

    if (draggedDocId && draggedDocId !== doc.id) {
      // Move the dragged doc into this folder
      onMoveDocumentToFolder(draggedDocId, doc.id, draggedFolderId);
    }
    setIsDraggedOver(false);
  };

  const renameInputRef = useRef<HTMLInputElement | null>(null);

  const finishRename = () => {
    const trimmed = renameValue.trim();
    if (trimmed && trimmed !== doc.title) {
      onRename(doc.id, trimmed);
    }
    setIsRenaming(false);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      e.preventDefault();
      e.stopPropagation();
      finishRename();
    } else if (e.key === "Escape") {
      setRenameValue(doc.title);
      setIsRenaming(false);
    }
  };

  useLayoutEffect(() => {
    if (isRenaming && renameInputRef.current) {
      renameInputRef.current.focus();
      renameInputRef.current.select();
    }
  }, [isRenaming]);

  const { data: folderData, loading: folderLoading } = useQuery(
    GetFolderDocuments,
    {
      variables: { folderId: doc.id, offset: 0, limit: DOCUMENT_LIMIT },
      skip: !isOpen && !showDeleteDialog,
    },
  );

  const folderDocuments = useMemo(
    () => folderData?.folderDocuments?.edges?.map((edge) => edge.node) ?? [],
    [folderData],
  );

  const handleClick = async () => {
    const newOpenState = !isOpen;
    setIsOpen(newOpenState);
    initialOpenStateRef.current = newOpenState;
  };

  const isMobile = useIsMobile();

  return (
    <div ref={folderRef}>
      <DeleteFolderAlert
        showDeleteDialog={showDeleteDialog}
        setShowDeleteDialog={setShowDeleteDialog}
        onClickDelete={(deleteChildren: boolean) =>
          onDelete(doc.id, doc.folderID, deleteChildren)
        }
        draftTitle={doc.title}
        hasDrafts={folderDocuments.length > 0}
        loading={folderLoading}
      />

      <SidebarMenuItem
        className={cn(
          "group/sidebar-menu-item",
          isDraggedOver && "border border-primary rounded-md",
        )}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        {isRenaming ? (
          <div className="flex items-center gap-2 px-2 py-1.5 bg-muted">
            <FolderIcon className="h-4 w-4" />
            <input
              ref={renameInputRef}
              className="border-b border-gray-300 bg-transparent outline-none"
              value={renameValue}
              onChange={(e) => setRenameValue(e.target.value)}
              onBlur={finishRename}
              onKeyDown={handleKeyDown}
            />
          </div>
        ) : (
          <>
            <SidebarMenuButton
              asChild
              disabled={isDisconnected}
              onClick={!isMobile ? handleClick : undefined}
              onTouchStart={isMobile ? handleClick : undefined}
              className="hover:has-[button:hover]:bg-accent"
            >
              <button>
                {isOpen ? (
                  <FolderOpenIcon className="h-4 w-4" />
                ) : (
                  <FolderIcon className="h-4 w-4" />
                )}
                <span>{doc.title}</span>
              </button>
            </SidebarMenuButton>
            <DropdownMenu open={dropdownOpen} onOpenChange={setDropdownOpen}>
              <DropdownMenuTrigger asChild>
                <SidebarMenuAction
                  className={cn(
                    "group-hover/sidebar-menu-item:opacity-100 hover:bg-transparent",
                    dropdownOpen ? "opacity-100" : "opacity-0",
                  )}
                >
                  <EllipsisIcon className="h-4 w-4" />
                </SidebarMenuAction>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-52" align="start" side="right">
                <DropdownMenuGroup>
                  <DropdownMenuItem
                    onClick={() => {
                      setIsRenaming(true);
                    }}
                  >
                    <Edit2Icon className="mr-2 h-4 w-4" />
                    <span>Rename</span>
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => setShowDeleteDialog(true)}
                    disabled={!canDelete}
                  >
                    <DeleteIcon className="mr-2 h-4 w-4" />
                    <span>Delete</span>
                  </DropdownMenuItem>
                </DropdownMenuGroup>
              </DropdownMenuContent>
            </DropdownMenu>
          </>
        )}
      </SidebarMenuItem>

      {isOpen && (folderLoading || folderDocuments.length > 0) && (
        <div
          ref={childrenRef}
          className="pl-4"
          onDragLeave={handleChildrenDragLeave}
        >
          {folderLoading ? (
            <div className="text-muted">Loading...</div>
          ) : (
            folderDocuments.map((folderDoc) => (
              <DocumentMenuItem
                key={folderDoc.id}
                doc={folderDoc as DocumentFieldsFragment}
                currentDocument={currentDocument}
                onRename={onRename}
                isDisconnected={isDisconnected}
                canDelete={canDelete}
                onNavigate={onNavigate}
                onDelete={onDelete}
                onMoveDocumentToFolder={onMoveDocumentToFolder}
              />
            ))
          )}
        </div>
      )}
    </div>
  );
}

type DocProps = {
  doc: DocumentFieldsFragment;
  currentDocument: DocumentFieldsFragment | null;
  onNavigate: (id: string) => void;
  isDisconnected: boolean;
  canDelete: boolean;
  onDelete: DeleteFunction;
  onMoveDocumentToFolder: MoveFunction;
};

export function Doc({
  doc,
  currentDocument,
  onNavigate,
  isDisconnected,
  canDelete,
  onDelete,
  onMoveDocumentToFolder,
}: DocProps) {
  const { toast } = useToast();
  const { draggedItem, setDraggedItem } = useDragContext();
  const [showDeleteDialog, _setShowDeleteDialog] = useState(false);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [showDropIndicator, setShowDropIndicator] = useState(false);

  // Make sure that the dialog is closed before updating the state of the delete dialog
  const setShowDeleteDialog = (value: boolean) => {
    setDropdownOpen(false);
    _setShowDeleteDialog(value);
  };

  const handleCopyLink = () => {
    navigator.clipboard.writeText(`${window.location.origin}/drafts/${doc.id}`);
    toast({
      title: "Copied link to clipboard",
    });
  };

  const handleDragStart = (e: React.DragEvent) => {
    const dragData = {
      id: doc.id,
      folderID: doc.folderID ?? null,
    };
    console.log("DragStart:", dragData);
    e.dataTransfer.setData("text/plain", JSON.stringify(dragData));
    setDraggedItem(dragData);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setShowDropIndicator(false);
    if (!draggedItem) return;

    const draggedDocId = draggedItem.id;
    const draggedFolderId = draggedItem.folderID;
    if (draggedDocId && draggedDocId !== doc.folderID) {
      onMoveDocumentToFolder(draggedDocId, doc.folderID, draggedFolderId);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    // Only show drop indicator if this is a top-level document (no folder)
    // and the dragged item has a folder ID
    console.log("DragOver:", {
      docFolderID: doc.folderID,
      draggedItem,
      showDropIndicator: doc.folderID === null && draggedItem?.folderID,
    });
    if (doc.folderID === null && draggedItem?.folderID) {
      setShowDropIndicator(true);
    }
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    console.log("DragLeave");
    setShowDropIndicator(false);
  };

  const isMobile = useIsMobile();

  return (
    <>
      <DeleteDraftAlert
        showDeleteDialog={showDeleteDialog}
        setShowDeleteDialog={setShowDeleteDialog}
        onClickDelete={() => onDelete(doc.id, doc.folderID)}
        draftTitle={doc.title}
      />

      <SidebarMenuItem
        className={cn(
          "group/sidebar-menu-item relative",
          showDropIndicator &&
            "before:content-[''] before:absolute before:-top-[2px] before:left-0 before:right-0 before:h-[2px] before:bg-primary before:z-10",
        )}
        draggable={true}
        onDragStart={handleDragStart}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        <SidebarMenuButton
          asChild
          isActive={currentDocument?.id === doc.id}
          disabled={isDisconnected}
          onClick={!isMobile ? () => onNavigate(doc.id) : undefined}
          onTouchStart={isMobile ? () => onNavigate(doc.id) : undefined}
        >
          <button>
            <FileIcon className="h-4 w-4" />
            <span>{doc.title}</span>
          </button>
        </SidebarMenuButton>

        <DropdownMenu open={dropdownOpen} onOpenChange={setDropdownOpen}>
          <DropdownMenuTrigger asChild>
            <SidebarMenuAction
              className={cn(
                "group-hover/sidebar-menu-item:opacity-100 hover:bg-transparent",
                dropdownOpen ? "opacity-100" : "opacity-0",
              )}
            >
              <EllipsisIcon className="h-4 w-4" />
            </SidebarMenuAction>
          </DropdownMenuTrigger>
          <DropdownMenuContent className="w-52" align="start" side="right">
            <DropdownMenuGroup>
              <DropdownMenuItem onClick={handleCopyLink}>
                <LinkIcon className="mr-2 h-4 w-4" />
                <span>Copy Link</span>
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => setShowDeleteDialog(true)}
                disabled={!canDelete}
              >
                <DeleteIcon className="mr-2 h-4 w-4" />
                <span>Delete</span>
              </DropdownMenuItem>
            </DropdownMenuGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </>
  );
}
