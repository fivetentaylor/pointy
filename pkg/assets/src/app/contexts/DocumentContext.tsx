import React, {
  createContext,
  useState,
  useContext,
  useMemo,
  useEffect,
} from "react";
import { useNavigate, useParams } from "react-router-dom";
import {
  useQuery,
  useMutation,
  useSubscription,
  useApolloClient,
} from "@apollo/client";
import { useErrorToast } from "@/hooks/useErrorToast";
import {
  CreateDocument,
  CreateFlaggedVersion,
  CreateFolder,
  DeleteDocument,
  DeleteFlaggedVersion,
  DocumentInsertedSubscription,
  DocumentUpdatedSubscription,
  EditFlaggedVersion,
  GetBaseDocuments,
  GetDocument,
  GetDocuments,
  GetFolderDocuments,
  GetSharedDocuments,
  MoveDocument,
  SharedDocumentLinks,
  ShareDocument,
  UnshareDocument,
  UpdateDocumentPreferences,
  UpdateDocumentTitle,
  UpdateDocumentVisibility,
  UpdateSharedLink,
} from "@/queries/document";
import { useCurrentUserContext } from "./CurrentUserContext";
import { analytics } from "@/lib/segment";
import {
  DRAFTS_CREATE_NEW,
  DRAFTS_CREATE_NEW_FOLDER,
  DRAFTS_DELETE,
} from "@/lib/events";
import { gql } from "@apollo/client";

type DocumentContextState = ReturnType<typeof useSetupDocument>;

type DocumentContextProviderProps = {
  children: React.ReactNode;
};

export const DOCUMENT_LIMIT = 50;

const DocumentContext = createContext<DocumentContextState | undefined>(
  undefined,
);

interface DocumentEdge {
  node: DocumentFieldsFragment;
  cursor: string;
}

const useSetupDocument = () => {
  const { draftId } = useParams();
  const navigate = useNavigate();
  const showErrorToast = useErrorToast();
  const { currentUser } = useCurrentUserContext();
  const apolloClient = useApolloClient();
  const [myDocumentsPage] = useState(0);

  if (!draftId) {
    window.location.href = "/error";
  }

  const {
    data: documentData,
    refetch: refetchDoc,
    loading: docLoading,
  } = useQuery(GetDocument, {
    variables: { id: draftId || "" },
  });

  const { data: sharedLinkData, refetch: refetchSharedLinks } = useQuery(
    SharedDocumentLinks,
    { variables: { id: draftId || "" } },
  );

  const {
    loading: loadingDocuments,
    data: documentsData,
    subscribeToMore: subscribeToMoreDocuments,
  } = useQuery(GetDocuments, {
    variables: {
      offset: myDocumentsPage * DOCUMENT_LIMIT,
      limit: DOCUMENT_LIMIT,
    },
    notifyOnNetworkStatusChange: true, // reflect cache updates
  });

  const {
    loading: loadingBaseDocuments,
    data: documentsBaseData,
    subscribeToMore: subscribeToMoreBaseDocuments,
  } = useQuery(GetBaseDocuments, {
    variables: {
      offset: myDocumentsPage * DOCUMENT_LIMIT,
      limit: DOCUMENT_LIMIT,
    },
    notifyOnNetworkStatusChange: true, // reflect cache updates
  });

  const {
    loading: loadingSharedDocuments,
    data: documentsSharedData,
    subscribeToMore: subscribeToMoreSharedDocuments,
  } = useQuery(GetSharedDocuments, {
    variables: {
      offset: myDocumentsPage * DOCUMENT_LIMIT,
      limit: DOCUMENT_LIMIT,
    },
    notifyOnNetworkStatusChange: true, // reflect cache updates
  });

  const getDocuments = async (offset = 0, limit = DOCUMENT_LIMIT) => {
    try {
      const { data } = await apolloClient.query({
        query: GetDocuments,
        variables: {
          offset,
          limit,
        },
        fetchPolicy: "network-only", // Always fetch fresh data
      });

      return data?.documents?.edges.map((edge) => edge.node) || [];
    } catch (error) {
      console.error(error);
      showErrorToast(`Failed to fetch folder documents}`);
      return [];
    }
  };

  const getFolderDocuments = async (
    folderId: string,
    offset = 0,
    limit = DOCUMENT_LIMIT,
  ) => {
    try {
      const { data } = await apolloClient.query({
        query: GetFolderDocuments,
        variables: {
          folderId,
          offset,
          limit,
        },
        fetchPolicy: "network-only", // Always fetch fresh data
      });

      return data?.folderDocuments?.edges.map((edge) => edge.node) || [];
    } catch (error) {
      console.error(error);
      showErrorToast(`Failed to fetch folder documents}`);
      return [];
    }
  };

  const documents = useMemo(
    () =>
      documentsData?.documents.edges
        .map((edge) => edge.node)
        .sort((a, b) => b.updatedAt.localeCompare(a.updatedAt)) || [],
    [documentsData],
  );

  const baseDocuments = useMemo(
    () =>
      documentsBaseData?.baseDocuments.edges
        .map((edge) => edge.node)
        .sort((a, b) => b.updatedAt.localeCompare(a.updatedAt)) || [],
    [documentsBaseData],
  );

  const sharedDocuments = useMemo(
    () =>
      documentsSharedData?.sharedDocuments.edges
        .map((edge) => edge.node)
        .sort((a, b) => b.updatedAt.localeCompare(a.updatedAt)) || [],
    [documentsSharedData],
  );

  const docData = useMemo(
    () => documentData?.document as DocumentFieldsFragment,
    [documentData],
  );
  const isDocPublic = useMemo(() => docData?.isPublic, [docData]);
  const editors = useMemo(
    () => docData?.editors.filter((e) => e.id !== currentUser?.id),
    [docData, currentUser?.id],
  );
  const sharedLinks = useMemo(
    () => sharedLinkData?.sharedLinks.filter((sl) => sl.isActive),
    [sharedLinkData],
  );

  const [createDocumentMutation, { loading: loadingCreateDocument }] =
    useMutation(CreateDocument, {
      refetchQueries: [{ query: GetDocuments }],
    });

  const [createFolderMutation, { loading: loadingCreateFolder }] = useMutation(
    CreateFolder,
    {
      refetchQueries: [{ query: GetDocuments }],
    },
  );

  const createNewDocument = async () => {
    analytics.track(DRAFTS_CREATE_NEW);
    const { data } = await createDocumentMutation();
    if (!data) return;
    const { createDocument } = data;
    navigate(`/drafts/${createDocument.id}`);
  };

  const createNewFolder = async () => {
    analytics.track(DRAFTS_CREATE_NEW_FOLDER);
    await createFolderMutation();
  };

  const [updateDocumentMutation, { loading: savingDocument }] =
    useMutation<any>(UpdateDocumentTitle);

  const [deleteDocumentMutation, { loading: deletingDocument }] =
    useMutation<any>(DeleteDocument, {
      refetchQueries: [GetDocument],
    });

  const updateDocument = async (input: DocumentInput) => {
    const { data, errors } = await updateDocumentMutation({
      variables: {
        input,
        id: draftId || "",
      },
    });

    return { data, errors };
  };

  const renameFolder = async (folderId: string, title: string) => {
    await updateDocumentMutation({
      variables: {
        input: {
          title,
        },
        id: folderId,
      },
      refetchQueries: [{ query: GetDocuments }],
    });
  };

  const deleteDocument = async (
    id: string,
    folderId?: string | null,
    deleteChildren?: boolean,
  ) => {
    analytics.track(DRAFTS_DELETE);

    // Stage 1: Immediate cache update
    apolloClient.cache.modify({
      fields: {
        documents: (existing, { readField }) => ({
          ...existing,
          edges: existing.edges.filter(
            (edge: DocumentEdge) => readField("id", edge.node) !== id,
          ),
        }),
        baseDocuments: (existing, { readField }) => ({
          ...existing,
          edges: existing.edges.filter(
            (edge: DocumentEdge) => readField("id", edge.node) !== id,
          ),
        }),
        folderDocuments: (existing, { readField }) => ({
          ...existing,
          edges: existing.edges.filter(
            (edge: DocumentEdge) => readField("id", edge.node) !== id,
          ),
        }),
      },
    });

    // Stage 2: Background processing
    if (!deleteChildren) {
      // Fetch and process folder documents in the background
      getFolderDocuments(id).then((folderDocs) => {
        if (folderDocs.length > 0) {
          apolloClient.cache.modify({
            fields: {
              documents: (existing) => ({
                ...existing,
                edges: [
                  ...existing.edges,
                  ...folderDocs.map((doc) => ({
                    node: { ...doc, folderID: null },
                    cursor: doc.id,
                  })),
                ],
              }),
              baseDocuments: (existing) => ({
                ...existing,
                edges: [
                  ...existing.edges,
                  ...folderDocs.map((doc) => ({
                    node: { ...doc, folderID: null },
                    cursor: doc.id,
                  })),
                ],
              }),
            },
          });
        }
      });
    }

    // Stage 2: Server deletion
    await deleteDocumentMutation({
      variables: {
        id,
        deleteChildren,
      },
      onError: () => {
        showErrorToast("Error deleting document");
        // Optionally: revert the cache changes here if the mutation fails
      },
      refetchQueries: [{ query: GetBaseDocuments }, { query: GetDocuments }],
    });

    // Navigation logic
    if (folderId) {
      const documents = await getFolderDocuments(folderId);
      if (documents.length === 0) {
        window.location.href = "/drafts";
      } else {
        navigate(`/drafts/${documents[0].id}`);
      }
    } else {
      const documents = await getDocuments();
      if (documents.length === 0) {
        window.location.href = "/drafts";
      } else {
        if (documents[0].isFolder) {
          window.location.href = "/drafts";
        } else {
          navigate(`/drafts/${documents[0].id}`);
        }
      }
    }
  };

  const [updatePreferences, { loading: savingPreferences }] = useMutation(
    UpdateDocumentPreferences,
  );

  const [shareDocumentMutation] = useMutation(ShareDocument, {
    variables: {
      id: draftId || "",
      emails: [],
      message: "",
    },
  });
  const shareDocument = async (emails: string[], message?: string) => {
    const { data, errors } = await shareDocumentMutation({
      variables: {
        id: draftId || "",
        emails,
        message,
      },
    });
    return { data, errors };
  };

  const [updateSharedLinkMutation] = useMutation(UpdateSharedLink);
  const updateShareLink = async (inviteLink: string, isActive: boolean) => {
    const { data, errors } = await updateSharedLinkMutation({
      variables: {
        inviteLink,
        isActive,
      },
    });
    return { data, errors };
  };

  const [unshareDocumentMutation] = useMutation(UnshareDocument);
  const unshareDocument = async (userID: string) => {
    const { data, errors } = await unshareDocumentMutation({
      variables: {
        docId: draftId || "",
        editorId: userID,
      },
    });
    return { data, errors };
  };

  const [updateDocumentVisibilityMutation] = useMutation(
    UpdateDocumentVisibility,
  );
  const updateDocumentVisibility = async (isPublic: boolean) => {
    const { data, errors } = await updateDocumentVisibilityMutation({
      variables: {
        id: draftId || "",
        input: {
          isPublic,
        },
      },
    });
    return { data, errors };
  };

  const [moveDocumentMutation] = useMutation(MoveDocument, {
    onError: (error) => {
      showErrorToast(`Failed to move document: ${error.message}`);
    },
    refetchQueries: [GetDocuments],
  });
  const moveDocument = async (
    documentId: string,
    newFolderId: string | null | undefined,
    oldFolderId: string | null | undefined,
  ) => {
    // Get the document from cache before mutation
    const document = apolloClient.cache.readFragment<{
      folderID: string | null;
    }>({
      id: `Document:${documentId}`,
      fragment: gql`
        fragment MovedDocument on Document {
          id
          folderID
        }
      `,
    });

    // Optimistically update the cache
    apolloClient.cache.modify({
      id: `Document:${documentId}`,
      fields: {
        folderID: () => newFolderId ?? null,
      },
    });

    try {
      const { data, errors } = await moveDocumentMutation({
        variables: {
          documentId,
          folderId: newFolderId,
        },
      });

      if (errors) {
        // Revert the cache on error
        apolloClient.cache.modify({
          id: `Document:${documentId}`,
          fields: {
            folderID: () => document?.folderID ?? null,
          },
        });
        return { data, errors };
      }

      // Refetch queries after successful move
      if (oldFolderId) {
        await getFolderDocuments(oldFolderId);
      }
      if (newFolderId) {
        await getFolderDocuments(newFolderId);
      }
      if (!oldFolderId || !newFolderId) {
        await apolloClient.refetchQueries({ include: [GetBaseDocuments] });
      }

      return { data, errors: null };
    } catch (error) {
      // Revert the cache on error
      apolloClient.cache.modify({
        id: `Document:${documentId}`,
        fields: {
          folderID: () => document?.folderID ?? null,
        },
      });
      throw error;
    }
  };

  const [createFlaggedVersionMutation] = useMutation(CreateFlaggedVersion, {
    onError: (error) => {
      showErrorToast(`Failed to create version: ${error.message}`);
    },
    variables: {
      documentId: draftId || "",
      input: {
        name: "",
        updateID: "",
      },
    },
  });

  const createFlaggedVersion = async (input: {
    name: string;
    updateID: string;
  }) => {
    return createFlaggedVersionMutation({
      variables: {
        documentId: draftId || "",
        input,
      },
    });
  };

  const [editFlaggedVersionMutation] = useMutation(EditFlaggedVersion, {
    onError: (error) => {
      showErrorToast(`Failed to edit version: ${error.message}`);
    },
  });

  const [deleteFlaggedVersionMutation] = useMutation(DeleteFlaggedVersion, {
    onError: (error) => {
      showErrorToast(`Failed to delete version: ${error.message}`);
    },
  });

  const editFlaggedVersion = async (
    flaggedVersionId: string,
    input: {
      name: string;
      updateID: string;
    },
  ) => {
    return editFlaggedVersionMutation({
      variables: {
        flaggedVersionId,
        input,
      },
    });
  };

  const deleteFlaggedVersion = async (
    flaggedVersionId: string,
    timelineEventId: string,
  ) => {
    return deleteFlaggedVersionMutation({
      variables: {
        flaggedVersionId,
        timelineEventId,
      },
    });
  };

  useSubscription(DocumentUpdatedSubscription, {
    variables: { documentId: draftId || "" },
  });

  useEffect(() => {
    if (!currentUser) return;

    const unsubscribe = subscribeToMoreDocuments({
      document: DocumentInsertedSubscription,
      variables: {
        userId: currentUser.id,
      },
      updateQuery: (prev, { subscriptionData }) => {
        const newDocument = subscriptionData?.data
          ?.documentInserted as DocumentFieldsFragment;

        if (
          newDocument &&
          !prev.documents.edges.some((edge) => edge.node.id === newDocument.id)
        ) {
          return {
            documents: {
              ...prev.documents,
              edges: [
                { node: newDocument, cursor: newDocument.id },
                ...prev.documents.edges,
              ],
            },
          };
        }

        return prev;
      },
    });

    const unsubscribeBase = subscribeToMoreBaseDocuments({
      document: DocumentInsertedSubscription,
      variables: {
        userId: currentUser.id,
      },
      updateQuery: (prev, { subscriptionData }) => {
        const newDocument = subscriptionData?.data
          ?.documentInserted as DocumentFieldsFragment;

        if (
          newDocument &&
          newDocument.access === "owner" &&
          newDocument.folderID === null &&
          !prev.baseDocuments.edges.some(
            (edge) => edge.node.id === newDocument.id,
          )
        ) {
          return {
            baseDocuments: {
              ...prev.baseDocuments,
              edges: [
                { node: newDocument, cursor: newDocument.id },
                ...prev.baseDocuments.edges,
              ],
            },
          };
        }

        return prev;
      },
    });

    const unsubscribeShared = subscribeToMoreSharedDocuments({
      document: DocumentInsertedSubscription,
      variables: {
        userId: currentUser.id,
      },
      updateQuery: (prev, { subscriptionData }) => {
        const newDocument = subscriptionData?.data
          ?.documentInserted as DocumentFieldsFragment;

        if (
          newDocument &&
          newDocument.access !== "owner" &&
          !prev.sharedDocuments.edges.some(
            (edge) => edge.node.id === newDocument.id,
          )
        ) {
          return {
            sharedDocuments: {
              ...prev.sharedDocuments,
              edges: [
                { node: newDocument, cursor: newDocument.id },
                ...prev.sharedDocuments.edges,
              ],
            },
          };
        }

        return prev;
      },
    });

    return () => {
      unsubscribe();
      unsubscribeBase();
      unsubscribeShared();
    };
  }, [currentUser, subscribeToMoreDocuments]);

  return {
    baseDocuments,
    createFlaggedVersion,
    createNewDocument,
    createNewFolder,
    deleteDocument,
    deleteFlaggedVersion,
    deletingDocument,
    docData,
    docLoading,
    documents,
    draftId,
    editFlaggedVersion,
    editors,
    isDocPublic,
    loadingCreateDocument,
    getFolderDocuments,
    loadingBaseDocuments,
    loadingCreateFolder,
    loadingDocuments,
    loadingSharedDocuments,
    moveDocument,
    refetchDoc,
    refetchSharedLinks,
    renameFolder,
    savingDocument,
    savingPreferences,
    shareDocument,
    sharedDocuments,
    sharedLinks,
    unshareDocument,
    updateDocument,
    updateDocumentVisibility,
    updatePreferences,
    updateShareLink,
  };
};

// Create a provider component
export const DocumentContextProvider = ({
  children,
}: DocumentContextProviderProps) => {
  const state = useSetupDocument();

  return (
    <DocumentContext.Provider value={state}>
      {children}
    </DocumentContext.Provider>
  );
};

// Custom hook to use the DocumentContext
export const useDocumentContext = () => {
  const context = useContext(DocumentContext);
  if (context === undefined) {
    throw new Error(
      "useDocumentContext must be used within a DocumentProvider",
    );
  }
  return context;
};
