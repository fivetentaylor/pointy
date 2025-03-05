import React, {
  useEffect,
  useRef,
  KeyboardEvent as ReactKeyboardEvent,
} from "react";
import { Editor, useEditor } from "@tiptap/react";
import Placeholder from "@tiptap/extension-placeholder";
import Document from "@tiptap/extension-document";
import Mention from "@tiptap/extension-mention";
import Paragraph from "@tiptap/extension-paragraph";
import Text from "@tiptap/extension-text";
import { EditorContent } from "@tiptap/react";

import { makeMentionListConfig } from "./MentionList";
import { useCurrentUserContext } from "@/contexts/CurrentUserContext";
import { useDocumentContext } from "@/contexts/DocumentContext";

type TipTapJSONContent = {
  content?: TipTapJSONContent[];
  type: string;
  text?: string;
  attrs?: {
    id: string;
    label: string;
  };
};

type SimpleUser = Omit<User, "isAdmin">;

const convertJSONTOText = (content: TipTapJSONContent[]) => {
  let text = "";
  content.forEach((item) => {
    if (item.content) {
      text += "\n";
      text += convertJSONTOText(item.content);
    }
    if (item.type === "mention" && item.attrs) {
      const mention = item.attrs;
      text += "@@" + btoa(`:user:${mention.id}:${mention.label}`) + "@@";
    } else if (item.type === "text") {
      text += item.text;
    } else if (item.type === "paragraph" && !item.content) {
      text += "\n";
    }
  });

  return text;
};

const convertTextToTipTapJSON = (text: string): any => {
  // Trim any leading or trailing whitespace, including newlines
  text = text.trim();

  const content: any[] = [];
  const regex = /@@([A-Za-z0-9+/=]+)@@/g;
  let lastIndex = 0;
  let match;

  while ((match = regex.exec(text)) !== null) {
    if (match.index > lastIndex) {
      content.push({
        type: "text",
        text: text.slice(lastIndex, match.index),
      });
    }

    const decodedMention = atob(match[1]);
    const [, id, label] = decodedMention.match(/:user:(.+):(.+)/) || [];

    if (id && label) {
      content.push({
        type: "mention",
        attrs: { id, label },
      });
    }

    lastIndex = regex.lastIndex;
  }

  if (lastIndex < text.length) {
    content.push({
      type: "text",
      text: text.slice(lastIndex),
    });
  }

  // Ensure there's always at least one paragraph, even if empty
  return {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: content.length > 0 ? content : [],
      },
    ],
  };
};

type MentionEditorProps = {
  autoFocus?: boolean;
  initialContent?: string;
  disabled?: boolean;
  mentionContainer?: HTMLElement | null;
  mentionsEnabled?: boolean;
  placeholder: string;
  onChange: (content: string) => void;
  onEnter: () => void;
  onLoaded: (editor: Editor) => void;
  onKeyDown?: (event: ReactKeyboardEvent<HTMLDivElement>) => void;
  onKeyUp?: (event: ReactKeyboardEvent<HTMLDivElement>) => void;
};

export const MentionEditor = ({
  autoFocus = false,
  initialContent = "",
  editors,
  usersInDomain,
  placeholder,
  onChange,
  disabled,
  mentionContainer,
  mentionsEnabled = true,
  onLoaded,
  onEnter,
}: MentionEditorProps & {
  editors: SimpleUser[] | undefined;
  usersInDomain: SimpleUser[];
}) => {
  const editorsRef = useRef<SimpleUser[]>([]);
  const usersInDomainRef = useRef<SimpleUser[]>([]);
  const mentionActiveRef = useRef(false);

  const extensions = [
    Document,
    Paragraph,
    Text,
    Placeholder.configure({
      placeholder,
      emptyEditorClass:
        "cursor-text before:content-[attr(data-placeholder)] before:absolute before:top-0 before:left-1 before:text-muted-foreground  before:pointer-events-none",
    }),
  ];

  if (mentionsEnabled) {
    extensions.push(
      Mention.configure({
        HTMLAttributes: {
          class: "bg-reviso-highlight text-primary",
        },
        suggestion: {
          ...makeMentionListConfig({
            mentionContainer,
            onToggleMentions: (showing) => {
              mentionActiveRef.current = showing;
            },
          }),
          items: ({ query }: { query: string }): (SimpleUser | string)[] => {
            const matchingEditors = editorsRef.current.filter((item) =>
              itemMatchesQuery(item, query),
            );

            const matchingUsersInDomain = usersInDomainRef.current
              .filter((item) => itemMatchesQuery(item, query))
              .slice(0, 5); // Limit to first 5 results

            // Combine and remove duplicates by comparing id
            const combinedResults = [
              ...matchingEditors,
              ...matchingUsersInDomain,
            ].reduce((acc, current) => {
              const x = acc.find((item) => item.id === current.id);
              if (!x) {
                return acc.concat([current]);
              } else {
                return acc;
              }
            }, [] as SimpleUser[]);

            return [...combinedResults, query];
          },
          findSuggestionMatch: (config) => {
            const { $position } = config;
            let start = $position.pos;
            let text = "";

            // Look backwards from the cursor position
            while (start > 0) {
              const char = $position.doc.textBetween(start - 1, start);
              if (char === " " || char === "\n") {
                break;
              }
              text = char + text;
              start--;
            }

            // Check if the text starts with '@' and has content after it
            // or if it's at the start of the content
            if (
              (text.startsWith("@") && text.length > 0) ||
              (text === "@" && start === 0)
            ) {
              const result = {
                range: { from: Math.max(start, 1), to: $position.pos },
                query: text.slice(1), // Remove the @ symbol
                text: text,
              };

              return result;
            }

            return null;
          },
        },
      }),
    );
  }

  const editor = useEditor({
    editable: !disabled,
    extensions,
    content: convertTextToTipTapJSON(initialContent),
    editorProps: {
      attributes: {
        class:
          "max-h-[16rem] focus:outline-none leading-[1.3125rem] overflow-y-auto overflow-x-hidden mt-1 px-1 break-anywhere-important",
      },
      handleDOMEvents: {
        keydown: (_, event) => {
          if (
            !mentionActiveRef.current &&
            event.key === "Enter" &&
            !event.shiftKey
          ) {
            event.preventDefault();
            onEnter();
            return false;
          }

          if (event.key === "Enter" && event.shiftKey) {
            event.preventDefault();
            editor?.commands.enter();
            return false;
          }
        },
      },
    },
    onUpdate: ({ editor }) => {
      const json = editor.getJSON();
      if (json.content) {
        onChange(convertJSONTOText(json.content as TipTapJSONContent[]));
      } else {
        onChange(editor.getText());
      }
    },
  });

  useEffect(() => {
    if (editor !== null && placeholder !== "") {
      const extensions = editor.extensionManager.extensions.filter(
        (extension) => extension.name === "placeholder",
      );
      if (extensions.length > 0) {
        extensions[0].options["placeholder"] = placeholder;
        editor.view.dispatch(editor.state.tr);
      }
    }
  }, [editor, placeholder]);

  useEffect(() => {
    usersInDomainRef.current = usersInDomain || [];
  }, [usersInDomain]);

  useEffect(() => {
    editorsRef.current = editors || [];
    if (autoFocus) {
      editor?.commands.focus();
    }
  }, [editors]);

  useEffect(() => {
    if (editor) {
      onLoaded(editor);
    }
  }, [editor]);

  if (!editor) {
    return null;
  }

  return <EditorContent editor={editor} />;
};

// Helper function to check if an item matches the query
const itemMatchesQuery = (item: SimpleUser, query: string): boolean => {
  const lowercaseQuery = query.toLowerCase();
  return (
    item.displayName.toLowerCase().startsWith(lowercaseQuery) ||
    item.name.toLowerCase().startsWith(lowercaseQuery)
  );
};

const MentionEditorWrapper = (props: MentionEditorProps) => {
  const { usersInDomain } = useCurrentUserContext();
  const { editors } = useDocumentContext();

  return (
    <MentionEditor {...props} editors={editors} usersInDomain={usersInDomain} />
  );
};

export default MentionEditorWrapper;
