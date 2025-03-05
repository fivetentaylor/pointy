import React from "react";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "./alert-dialog";
import { Spinner } from "./spinner";
import { Checkbox } from "./checkbox";

export const DeleteFolderAlert = function ({
  showDeleteDialog,
  setShowDeleteDialog,
  onClickDelete,
  draftTitle,
  hasDrafts,
  loading,
}: {
  showDeleteDialog: boolean;
  setShowDeleteDialog: (show: boolean) => void;
  onClickDelete: (deleteChildren: boolean) => void;
  draftTitle: string;
  hasDrafts: boolean;
  loading: boolean;
}) {
  const [deleteChildren, setDeleteChildren] = React.useState(false);

  const handleDelete = () => {
    onClickDelete(deleteChildren);
  };

  return (
    <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete &quot;{draftTitle}?&quot;</AlertDialogTitle>
          <AlertDialogDescription>
            Confirm you want to delete the folder
            {hasDrafts ? " and the drafts" : ""}. You will no longer have access
            to the folder
            {hasDrafts && (
              <div className="mt-2 flex items-center space-x-2">
                <Checkbox
                  id="delete-drafts"
                  checked={deleteChildren}
                  onCheckedChange={(v: boolean) => setDeleteChildren(v)}
                />
                <label
                  htmlFor="delete-drafts"
                  className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Delete drafts in the folder as well
                </label>
              </div>
            )}
          </AlertDialogDescription>
        </AlertDialogHeader>
        {loading && (
          <AlertDialogFooter>
            <Spinner />
          </AlertDialogFooter>
        )}
        {!loading && (
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>

            <AlertDialogAction
              onClick={handleDelete}
              className="bg-destructive hover:bg-destructive/80"
            >
              Delete folder{deleteChildren ? " and drafts" : ""}
            </AlertDialogAction>
          </AlertDialogFooter>
        )}
      </AlertDialogContent>
    </AlertDialog>
  );
};
