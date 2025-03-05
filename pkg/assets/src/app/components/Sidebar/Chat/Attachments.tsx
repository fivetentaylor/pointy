import React, { useState, useEffect, useRef, useMemo } from "react";
import {
  AlertTriangleIcon,
  FileIcon,
  FileTextIcon,
  FolderClosedIcon,
  LaptopIcon,
  LinkIcon,
  PlusCircleIcon,
  SearchIcon,
  XIcon,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { FetchResult, useQuery } from "@apollo/client";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuPortal,
} from "@/components/ui/dropdown-menu";
import { useDocumentContext } from "@/contexts/DocumentContext";
import { Spinner } from "@/components/ui/spinner";
import { useCurrentUserContext } from "@/contexts/CurrentUserContext";
import { ScrollArea } from "@/components/ui/scroll-area";
import { ListUsersAttachments } from "@/queries/attachments";
import { analytics } from "@/lib/segment";
import {
  AI_CLICK_ADD_CONTEXT,
  AI_CLICK_CONTEXT_BUTTON,
  AI_CONTEXT_ADD_DRAFT,
  AI_CONTEXT_ADD_FILE,
  AI_CONTEXT_CLICK_ADD_URL,
  AI_CONTEXT_CLICK_REMOVE_ITEM,
  AI_CONTEXT_SAVE_URL,
  AI_CONTEXT_SEARCH,
  AI_CONTEXT_UPLOAD_FILE,
} from "@/lib/events";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { ContentTypeDisplayInfoMap } from "./ContentTypeDislpay";
import { useErrorToast } from "@/hooks/useErrorToast";

const acceptTypes = [
  "text/markdown",
  ".md",
  "text/csv",
  ".csv",
  "application/msword",
  "application/vnd.ms-word",
  ".doc",
  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
  ".docx",
  "application/vnd.openxmlformats-officedocument.presentationml.presentation",
  ".pptx",
  "application/vnd.oasis.opendocument.text",
  ".odt",
  // "application/vnd.apple.pages",
  // "application/x-iwork-pages-sffpages",
  // ".pages",
  "application/pdf",
  ".pdf",
  "application/rtf",
  "application/x-rtf",
  "text/rtf",
  "text/richtext",
  ".rtf",
  "text/html",
  ".html",
  ".htm",
  "text/url",
  ".url",
  "text/xml",
  "application/xml",
  ".xml",
  "text/plain",
  ".txt",
].join(",");

export type AttachmentType = {
  id: string;
  type: "file" | "draft" | "selection" | "url";
  name: string;
  status?: "success" | "error" | "uploading";
  contentType?: string;
};

type AttachmentsProps = {
  selectedHtml: string;
  uploadAttachment: (
    file: File,
  ) => Promise<FetchResult<UploadAttachmentMutation>>;
  onClearSelection: () => void;
  activeAttachments: AttachmentType[];
  setActiveAttachments: (value: AttachmentType[]) => void;
  setActiveSelection: (value: string) => void;
};

const AttachmentMenuContents = ({
  searchText,
  activeAttachments,
  documentAttachments,
  setActiveAttachments,
  onAddURL,
  onUploadFromComputer,
}: {
  documentAttachments: AttachmentType[];
  onAddURL: () => void;
  onUploadFromComputer: () => void;
  searchText: string;
  removeAttachment: (id: string) => void;
} & Pick<
  AttachmentsProps,
  "uploadAttachment" | "activeAttachments" | "setActiveAttachments"
>) => {
  const { documents } = useDocumentContext();
  const { currentUser } = useCurrentUserContext();
  const { myDocuments, sharedDocuments } = documents.reduce<{
    myDocuments: typeof documents;
    sharedDocuments: typeof documents;
  }>(
    (acc, doc) => {
      if (doc.ownedBy?.id === currentUser?.id) {
        acc.myDocuments.push(doc);
      } else {
        acc.sharedDocuments.push(doc);
      }
      return acc;
    },
    { myDocuments: [], sharedDocuments: [] },
  );

  const filteredDocuments = useMemo(() => {
    return {
      myDocuments: myDocuments.filter((doc) =>
        doc.title.toLowerCase().includes(searchText.toLowerCase()),
      ),
      sharedDocuments: sharedDocuments.filter((doc) =>
        doc.title.toLowerCase().includes(searchText.toLowerCase()),
      ),
    };
  }, [searchText, myDocuments, sharedDocuments]);

  const filteredAttachments = useMemo(() => {
    return documentAttachments.filter((attachment) =>
      attachment.name.toLowerCase().includes(searchText.toLowerCase()),
    );
  }, [searchText, documentAttachments]);

  const handleDraftClick = (doc: (typeof documents)[0]) => {
    analytics.track(AI_CONTEXT_ADD_DRAFT);
    setActiveAttachments([
      ...activeAttachments,
      {
        id: doc.id,
        type: "draft",
        name: doc.title,
        status: "success",
      },
    ]);
  };

  const handleSourceClick = (attachment: AttachmentType) => {
    analytics.track(AI_CONTEXT_ADD_FILE);
    setActiveAttachments([
      ...activeAttachments,
      {
        id: attachment.id,
        status: "success",
        name: attachment.name,
        contentType: attachment.contentType,
        type: "file",
      },
    ]);
  };

  return (
    <DropdownMenuGroup>
      {searchText.length > 0 &&
        filteredDocuments.myDocuments.length < 1 &&
        filteredDocuments.sharedDocuments.length < 1 &&
        filteredAttachments.length < 1 && (
          <DropdownMenuItem>
            <div className="min-h-28">
              <div className="text-muted-foreground text-xs font-medium leading-none">
                No results found
              </div>
            </div>
          </DropdownMenuItem>
        )}
      {searchText.length > 0 &&
        (filteredDocuments.myDocuments.length > 0 ||
          filteredDocuments.sharedDocuments.length > 0 ||
          filteredAttachments.length > 0) && (
          <ScrollArea className="max-w-56 h-52 max-h-52 min-h-52">
            {filteredDocuments.myDocuments.length > 0 && (
              <>
                {filteredDocuments.myDocuments.length > 0 && (
                  <div className="my-2 ml-2 text-muted-foreground text-xs font-medium leading-none">
                    Your drafts
                  </div>
                )}
                {filteredDocuments.myDocuments.map((doc) => (
                  <DocumentDropdownItem
                    key={doc.id}
                    title={doc.title}
                    onClick={() => handleDraftClick(doc)}
                  />
                ))}
              </>
            )}
            {filteredDocuments.sharedDocuments.length > 0 && (
              <>
                <div className="my-2 ml-2 text-muted-foreground text-xs font-medium leading-none">
                  Shared with you
                </div>
                {filteredDocuments.sharedDocuments.map((doc) => (
                  <DocumentDropdownItem
                    key={doc.id}
                    title={doc.title}
                    onClick={() => handleDraftClick(doc)}
                  />
                ))}
              </>
            )}
            {filteredAttachments.length > 0 && (
              <>
                <div className="my-2 ml-2 text-muted-foreground text-xs font-medium leading-none">
                  Your sources
                </div>
                {filteredAttachments?.map((attachment) => (
                  <SourceDropdownItem
                    key={attachment.id}
                    contentType={attachment.contentType}
                    name={attachment.name}
                    type={attachment.type}
                    onClick={() => handleSourceClick(attachment)}
                  />
                ))}
              </>
            )}
          </ScrollArea>
        )}

      {searchText.length < 1 && (
        <>
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>
              <FileIcon className="w-4 h-4 mr-2" />
              <span>Your drafts</span>
            </DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent sideOffset={10}>
                <ScrollArea
                  className="max-w-56 max-h-52"
                  style={{
                    height: `${(myDocuments.length + sharedDocuments.length) * 2 + 3.5}rem`, // add 3.5 to account for the Your/Shared labels
                  }}
                >
                  {myDocuments.length > 0 && (
                    <div className="my-2 ml-2 text-muted-foreground text-xs font-medium leading-none">
                      Your drafts
                    </div>
                  )}
                  {myDocuments.map((doc) => (
                    <DocumentDropdownItem
                      key={doc.id}
                      title={doc.title}
                      onClick={() => handleDraftClick(doc)}
                    />
                  ))}
                  {sharedDocuments.length > 0 && (
                    <div className="my-2 ml-2 text-muted-foreground text-xs font-medium leading-none">
                      Shared with you
                    </div>
                  )}
                  {sharedDocuments?.map((doc) => (
                    <DocumentDropdownItem
                      key={doc.id}
                      title={doc.title}
                      onClick={() => handleDraftClick(doc)}
                    />
                  ))}
                </ScrollArea>
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>

          {documentAttachments?.length > 0 && (
            <DropdownMenuSub>
              <DropdownMenuSubTrigger>
                <FolderClosedIcon className="w-4 h-4 mr-2" />
                <span>All sources</span>
              </DropdownMenuSubTrigger>
              <DropdownMenuPortal>
                <DropdownMenuSubContent sideOffset={10}>
                  <ScrollArea
                    className="max-w-56 max-h-52"
                    style={{
                      height: `${documentAttachments.length * 2}rem`,
                    }}
                  >
                    {documentAttachments?.map((attachment) => (
                      <SourceDropdownItem
                        key={attachment.id}
                        contentType={attachment.contentType}
                        name={attachment.name}
                        type={attachment.type}
                        onClick={() => handleSourceClick(attachment)}
                      />
                    ))}
                  </ScrollArea>
                </DropdownMenuSubContent>
              </DropdownMenuPortal>
            </DropdownMenuSub>
          )}

          <DropdownMenuItem onClick={onAddURL}>
            <>
              <LinkIcon className="w-4 h-4 mr-2" />
              <span>Add text from URL</span>
            </>
          </DropdownMenuItem>

          <DropdownMenuItem onClick={onUploadFromComputer}>
            <>
              <LaptopIcon className="w-4 h-4 mr-2" />
              <span>Upload from computer</span>
            </>
          </DropdownMenuItem>
        </>
      )}
    </DropdownMenuGroup>
  );
};

const DocumentDropdownItem = ({
  title,
  onClick,
}: {
  title: string;
  onClick: () => void;
}) => {
  return (
    <DropdownMenuItem
      onClick={onClick}
      className="hover:bg-accent hover:text-accent-foreground"
    >
      <FileIcon className="w-4 h-4 min-w-4 mr-2" />
      <span className="truncate">{title}</span>
    </DropdownMenuItem>
  );
};

const SourceDropdownItem = ({
  name,
  type,
  contentType,
  onClick,
}: {
  name: string;
  contentType: AttachmentType["contentType"];
  type: AttachmentType["type"];
  onClick: () => void;
}) => {
  const displayInfo = ContentTypeDisplayInfoMap[contentType || ""] || {
    label: "",
    color: "#000000",
    icon: <FileTextIcon className="w-4 h-4 min-w-4 text-white" />,
  };
  console.log(type, displayInfo, ContentTypeDisplayInfoMap);
  return (
    <DropdownMenuItem
      onClick={onClick}
      className="hover:bg-accent hover:text-accent-foreground flex items-center"
    >
      <div
        className={cn(
          "w-4 min-w-4 h-4 rounded-sm justify-center items-center inline-flex mr-2",
        )}
      >
        {displayInfo.iconFG}
      </div>
      <span className="truncate">{name}</span>
    </DropdownMenuItem>
  );
};

const Attachments = ({
  selectedHtml,
  uploadAttachment,
  onClearSelection,
  activeAttachments,
  setActiveAttachments,
  setActiveSelection,
}: AttachmentsProps) => {
  const [fileSearchInput, setFileSearchInput] = useState("");
  const [attachmentsOpen, setAttachmentsOpen] = useState(false);
  const [isShowingURLDialog, setIsShowingURLDialog] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);
  const removeAttachment = (id: string) => {
    analytics.track(AI_CONTEXT_CLICK_REMOVE_ITEM);
    setActiveAttachments(activeAttachments.filter((a) => a.id !== id));
  };
  const { data: documentAttachmentsData } = useQuery(ListUsersAttachments, {
    fetchPolicy: "cache-and-network",
  });
  const documentAttachments: AttachmentType[] = useMemo(() => {
    return (
      documentAttachmentsData?.listUsersAttachments.map((da) => {
        return {
          id: da.id,
          name: da.filename,
          type: "file",
          contentType: da.contentType,
        };
      }) || []
    );
  }, [documentAttachmentsData]);

  const showErrorToast = useErrorToast();

  // there's a bug with radix dialog + alert interaction that can leave doc in modal state
  useEffect(() => {
    if (!isShowingURLDialog) {
      setTimeout(() => {
        if (window.getComputedStyle(document.body).pointerEvents === "none") {
          document.body.style.pointerEvents = "auto";
        }
      }, 250);
    }
  }, [isShowingURLDialog]);

  const handleUploadFromComputer = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  const handleAddURL = () => {
    analytics.track(AI_CONTEXT_CLICK_ADD_URL);
    setIsShowingURLDialog(true);
  };

  const handleSaveURL = async (url: string) => {
    analytics.track(AI_CONTEXT_SAVE_URL);
    setIsShowingURLDialog(false);
    const uploadingAttachment = createPlaceholderAttachment({
      name: url,
      type: "url",
      contentType: "text/html",
    });

    setActiveAttachments([...activeAttachments, uploadingAttachment]);
    try {
      const urlEscaped = encodeURI(url);
      const file = new File([url], urlEscaped, { type: "text/url" });
      const response = await uploadAttachment(file);
      if (!response || !response.data || !response.data.uploadAttachment) {
        showErrorToast(
          "We couldn’t download content from the URL you provided. Some websites restrict access to services like ours, which might be causing this issue. Please double-check the URL or try using a different one. ",
        );
        removeAttachment("uploading");
        return;
      }

      const newAttachment = {
        id: response.data.uploadAttachment.id,
        type: "file",
        name: response.data.uploadAttachment.filename,
        status: "success",
        contentType: file.type,
      } as AttachmentType;

      removeAttachment("uploading");
      setActiveAttachments([...activeAttachments, newAttachment]);
    } catch (err) {
      removeAttachment("uploading");
      console.error("Upload error:", err);
    }
  };

  const createPlaceholderAttachment = ({
    name,
    contentType,
    type,
  }: {
    name: string;
    contentType: string;
    type: AttachmentType["type"];
  }) =>
    ({
      id: "uploading",
      status: "uploading",
      type,
      name,
      contentType,
    }) as AttachmentType;

  const handleFileChange = async (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    if (!event.target.files) return;
    analytics.track(AI_CONTEXT_UPLOAD_FILE);
    const file = event.target.files[0];

    if (file) {
      const uploadingAttachment = createPlaceholderAttachment({
        name: file.name,
        type: "file",
        contentType: file.type,
      });

      setActiveAttachments([...activeAttachments, uploadingAttachment]);
      try {
        const response = await uploadAttachment(file);
        if (!response || !response.data || !response.data.uploadAttachment) {
          removeAttachment("uploading");
          return;
        }

        const newAttachment = {
          id: response.data.uploadAttachment.id,
          type: "file",
          name: response.data.uploadAttachment.filename,
          status: "success",
          contentType: file.type,
        } as AttachmentType;

        removeAttachment("uploading");
        setActiveAttachments([...activeAttachments, newAttachment]);

        event.target.value = "";
      } catch (err) {
        removeAttachment("uploading");
        console.error("Upload error:", err);
      }
    }
  };

  return (
    <>
      <AddURLDialog
        isOpen={isShowingURLDialog}
        setIsOpen={(open: boolean) => {
          setIsShowingURLDialog(open);
        }}
        onSave={handleSaveURL}
      />

      <DropdownMenu onOpenChange={setAttachmentsOpen} open={attachmentsOpen}>
        <DropdownMenuTrigger asChild>
          <Button
            variant="icon"
            size="icon"
            className="h-6 w-6 p-1 justify-start cursor-pointer text-muted-foreground hover:text-foreground hover:bg-elevated focus-visible:ring-0 focus-visible:ring-offset-0 focus:ring-0"
            onClick={() => {
              analytics.track(AI_CLICK_CONTEXT_BUTTON);
            }}
          >
            <PlusCircleIcon className="w-4 min-w-4 h-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="start" side="top">
          <DropdownMenuItem
            onClick={(evt) => {
              evt.stopPropagation();
              evt.preventDefault();
            }}
          >
            <div className="max-w-md w-full h-8">
              <div className="relative">
                <SearchIcon className="absolute left-2 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                <input
                  type="text"
                  ref={searchInputRef}
                  placeholder="Search..."
                  value={fileSearchInput}
                  onKeyDown={(e) => {
                    if (e.key !== "Escape") {
                      e.stopPropagation();
                    }
                  }}
                  onBlurCapture={(e) => {
                    const text = e.target.value.trim();
                    if (text.length > 0) {
                      e.preventDefault();
                      e.stopPropagation();
                      searchInputRef.current?.focus();
                    }
                  }}
                  onChange={(evt) => {
                    if (evt.target.value.length > 0) {
                      analytics.track(AI_CONTEXT_SEARCH);
                    }
                    setFileSearchInput(evt.target.value);
                  }}
                  className="w-full h-8 pl-8 pr-4 py-2 border border-border rounded-md focus:outline-none focus:ring-0 focus-visible:ring-0"
                />
              </div>
            </div>
          </DropdownMenuItem>
          <AttachmentMenuContents
            searchText={fileSearchInput}
            uploadAttachment={uploadAttachment}
            activeAttachments={activeAttachments}
            documentAttachments={documentAttachments}
            setActiveAttachments={setActiveAttachments}
            removeAttachment={removeAttachment}
            onUploadFromComputer={handleUploadFromComputer}
            onAddURL={handleAddURL}
          />
        </DropdownMenuContent>
      </DropdownMenu>

      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        accept={acceptTypes}
        style={{ display: "none" }}
      />

      {selectedHtml && (
        <AttachmentButton
          attachment={{
            id: "",
            type: "selection",
            name: "Selected text",
            status: "success",
          }}
          onClearSelection={onClearSelection}
          onSelect={(isSelected) => {
            if (isSelected) {
              setActiveSelection(selectedHtml);
            } else {
              setActiveSelection("");
            }
          }}
        />
      )}
      {activeAttachments.map((attachment) => (
        <AttachmentButton
          key={attachment.id}
          attachment={attachment}
          onClearSelection={() => removeAttachment(attachment.id)}
          onSelect={() => {}}
        />
      ))}
      {activeAttachments.length === 0 && !selectedHtml && (
        <div
          className="text-muted-foreground text-sm font-medium leading-none pt-[1px] cursor-pointer"
          onClick={() => {
            analytics.track(AI_CLICK_ADD_CONTEXT);
            setAttachmentsOpen(true);
          }}
        >
          Add context
        </div>
      )}
    </>
  );
};

const AddURLDialog = ({
  isOpen,
  setIsOpen,
  onSave,
}: {
  isOpen: boolean;
  setIsOpen: (val: boolean) => void;
  onSave: (url: string) => void;
}) => {
  const [url, setURL] = useState("");
  const [hasEdited, setHasEdited] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  const prependProtocol = (url: string) => {
    if (url.startsWith("http://") || url.startsWith("https://")) {
      return url;
    }
    return `https://${url}`;
  };

  useEffect(() => {
    if (hasEdited) {
      return;
    }
    if (url.trim().length > 0) {
      setHasEdited(true);
    }
  }, [url]);

  const isValidURL = useMemo(() => {
    if (!url.trim()) {
      return false;
    }
    const urlToCheck = prependProtocol(url);

    try {
      const parsedUrl = new URL(urlToCheck);
      // Check for minimum domain requirements
      const hostParts = parsedUrl.hostname.split(".");
      // Must have at least 2 parts (domain and TLD)
      // Last part (TLD) must be at least 2 chars
      console.log(parsedUrl.host);
      return (
        hostParts.length >= 2 && hostParts[hostParts.length - 1].length >= 2
      );
    } catch {
      console.log("invalid url");
      return false;
    }
  }, [url, isEditing]);

  const showError = !isValidURL && !isEditing && hasEdited;

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add text from URL</DialogTitle>
          <DialogDescription>
            Only visible website text can be imported, and paid articles are not
            currently supported.
          </DialogDescription>
          <div>
            <div className="mt-2 flex flex-col">
              <Input
                className="w-full"
                onChange={(e) => setURL(e.target.value)}
                onFocus={() => setIsEditing(true)}
                onBlur={() => {
                  setURL(prependProtocol(url));
                  setIsEditing(false);
                }}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    onSave(url);
                  }
                }}
                value={url}
                placeholder="https://example.com"
              />
              {showError && (
                <p className="flex items-center text-destructive text-xs text-normal leading-none mt-1">
                  <AlertTriangleIcon className="w-4 h-4 inline mr-1" />
                  <span>Please enter a valid URL</span>
                </p>
              )}
            </div>
            <div className={cn("flex justify-end", showError ? "" : "mt-4")}>
              <Button
                size="sm"
                className="mr-2"
                variant="outline"
                onClick={() => {
                  setURL("");
                  setIsOpen(false);
                }}
              >
                Cancel
              </Button>
              <Button
                size="sm"
                className="bg-primary hover:bg-primary/90"
                disabled={!url.trim() || !isValidURL}
                onClick={() => {
                  onSave(url);
                }}
              >
                Add text from URL
              </Button>
            </div>
          </div>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};

const AttachmentButton = ({
  attachment,
  onSelect,
  onClearSelection,
}: {
  attachment: AttachmentType;
  onSelect: (isSelected: boolean) => void;
  onClearSelection: () => void;
}) => {
  const [isSelected, setSelected] = useState(false);

  useEffect(() => {
    onSelect(isSelected);
  }, [isSelected]);

  const contentTypeDisplayInfo = ContentTypeDisplayInfoMap[
    attachment.contentType || ""
  ] || {
    label: "",
    color: "#000000",
    icon: <FileTextIcon className="w-4 h-4 text-white" />,
  };

  return (
    <Button
      variant="outline"
      onClick={() => {
        setSelected(!isSelected);
      }}
      className={cn(
        "flex gap-0.5 py-2 px-1.5 h-6 bg-transparent hover:bg-transparent max-w-44",
        isSelected ? "border border-primary" : "",
      )}
    >
      {attachment.status === "uploading" && (
        <Spinner className="w-3 min-w-3 h-3 mr-1" />
      )}
      {attachment.status === "success" && (
        <>
          {attachment.type === "file" && (
            <>
              <div
                className={cn(
                  "w-3 min-w-3 h-3 rounded-sm justify-center items-center inline-flex mr-1",
                  contentTypeDisplayInfo.color,
                )}
              >
                {contentTypeDisplayInfo.lilIcon}
              </div>
            </>
          )}
          {attachment.type === "draft" && (
            <div className="w-3 min-w-3 h-3 bg-reviso rounded-sm justify-center items-center inline-flex mr-1">
              <FileIcon className="w-2 min-w-2 h-2 text-background" />
            </div>
          )}
        </>
      )}
      <span className="truncate">{attachment.name}</span>
      {attachment.status === "success" && (
        <Button
          variant="icon"
          size="icon"
          onClick={() => onClearSelection()}
          className="h-4 w-4 max-h-4 max-w-4 ml-1"
        >
          <XIcon className="w-4 min-w-4 h-4" />
        </Button>
      )}
    </Button>
  );
};

export default Attachments;
