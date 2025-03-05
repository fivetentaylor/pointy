import React, { useEffect } from "react";
import { MoreHorizontalIcon } from "lucide-react";
import { useCursorContext } from "@/contexts/CursorContext";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { useCurrentUserContext } from "@/contexts/CurrentUserContext";
import { AuthorInfo } from "@/../rogueEditor";
import { useTimelineContext } from "./TimelineContext";
import { useParams } from "react-router-dom";

type TimelineActiveWritersProps = {
  onMount: () => void;
};

const TimelineActiveWriters: React.FC<TimelineActiveWritersProps> = ({
  onMount,
}) => {
  const { currentUser } = useCurrentUserContext();
  const { editing, cursors } = useCursorContext();

  const editingCursors = cursors.filter((cursor) => cursor.editing);

  return (
    <>
      {editing && currentUser && (
        <TimelineActiveWriter
          allowSummarize
          name={"You"}
          title="are editing this document"
          userId={currentUser.id}
          onMount={onMount}
        />
      )}
      {editingCursors.map((cursor) => (
        <TimelineActiveWriter
          allowSummarize={false}
          key={cursor.userID}
          name={cursor.name}
          title="is editing this document"
          userId={cursor.userID}
          cursor={cursor}
          onMount={onMount}
        />
      ))}
    </>
  );
};

type TimelineActiveWriterProps = {
  allowSummarize: boolean;
  name: string;
  title: string;
  userId: string;
  cursor?: AuthorInfo;
  onMount: () => void;
};

const TimelineActiveWriter: React.FC<TimelineActiveWriterProps> = ({
  allowSummarize,
  name,
  title,
  userId,
  cursor,
  onMount,
}) => {
  const { currentUser } = useCurrentUserContext();
  const { editor } = useRogueEditorContext();
  const { draftId } = useParams();
  const { forceTimelineUpdateSummary } = useTimelineContext();

  useEffect(() => {
    onMount();
  }, []);

  const handleSummarize = () => {
    if (!editor) {
      return;
    }

    if (currentUser?.id === userId) {
      editor.editing = false;
    } else {
      if (cursor) {
        cursor.editing = false;
        editor.updateCursor(cursor);
      }
    }

    forceTimelineUpdateSummary({
      variables: {
        documentId: draftId || "",
        userId: userId,
      },
    });
  };

  return (
    <div className="flex items-start mb-4 first:mt-auto px-4">
      <div className="w-4 h-4 flex items-center justify-start mr-2">
        <MoreHorizontalIcon className="w-4 h-4 text-muted-foreground" />
      </div>
      <div className="flex-grow mt-[0.125rem]">
        <div className="text-xs font-normal">
          <span>{name} </span>
          <span>{title}</span>
          {allowSummarize && (
            <>
              <span className="ml-1 text-muted-foreground">•</span>
              <span
                className="ml-1 text-reviso hover:underline cursor-pointer"
                onClick={handleSummarize}
              >
                Finish editing & summarize
              </span>
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default TimelineActiveWriters;
