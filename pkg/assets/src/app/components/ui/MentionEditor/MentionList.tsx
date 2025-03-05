import React, {
  forwardRef,
  useEffect,
  useImperativeHandle,
  useState,
} from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { cn, getInitials } from "@/lib/utils";
import { ReactRenderer } from "@tiptap/react";
import { SuggestionKeyDownProps, SuggestionProps } from "@tiptap/suggestion";
import tippy from "tippy.js";
import { UserCircleIcon } from "lucide-react";
import { analytics } from "@/lib/segment";
import { COMMENT_AT_MENTION } from "@/lib/events";

export type MentionListRef = {
  onKeyDown: (props: { event: KeyboardEvent }) => boolean;
};

export type MentionListProps = {
  items: (User | string)[];
  command: ({ id, label }: { id: string; label: string }) => void;
  query: string; // Add this line
};

const isValidEmail = (email: string) => {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
};

export const makeMentionListConfig = ({
  onToggleMentions,
  mentionContainer,
}: {
  onToggleMentions: (showing: boolean) => void;
  mentionContainer?: HTMLElement | null;
}) => {
  return {
    decorationClass: "bg-accent",

    render: () => {
      let component: ReactRenderer<MentionListRef, MentionListProps> | null =
        null;
      let popup: any;

      return {
        onStart: (props: SuggestionProps) => {
          onToggleMentions(true);
          component = new ReactRenderer(MentionList, {
            props: {
              ...props,
              query: props.query, // Add this line
            },
            editor: props.editor,
          });

          popup = tippy("body", {
            getReferenceClientRect: () => {
              const emptyRect = new DOMRect();
              if (!props.clientRect) {
                return emptyRect;
              }
              const rect = props.clientRect();
              if (!rect) {
                return emptyRect;
              }
              return rect;
            },
            appendTo: () => mentionContainer || document.body,
            content: component.element,
            showOnCreate: true,
            interactive: true,
            trigger: "manual",
            placement: "bottom-start",
            render(instance) {
              const popper = document.createElement("div");
              const box = instance.props.content as HTMLElement;
              popper.appendChild(box);
              return { popper };
            },
          });
        },

        onUpdate(props: SuggestionProps) {
          component?.updateProps(props);

          if (!props.clientRect) {
            return;
          }

          popup[0].setProps({
            getReferenceClientRect: props.clientRect,
          });
        },

        onKeyDown(props: SuggestionKeyDownProps) {
          if (props.event.key === "Escape") {
            popup[0].hide();

            return true;
          }

          if (component && component.ref) {
            return component.ref.onKeyDown(props);
          }

          return false;
        },

        onExit() {
          onToggleMentions(false);
          popup[0].destroy();
          component?.destroy();
        },
      };
    },
  };
};

const INVITE_ID = "new-user-id";

export const MentionList = forwardRef<MentionListRef, MentionListProps>(
  ({ items, command, query }, ref) => {
    // Add query to the destructured props
    const [selectedIndex, setSelectedIndex] = useState(0);

    const selectItem = (index: number) => {
      const item = items[index];

      if (typeof item === "string") {
        analytics.track(COMMENT_AT_MENTION, { isInvite: true });
        command({ id: INVITE_ID, label: item });
      } else {
        analytics.track(COMMENT_AT_MENTION, { isInvite: false });
        command({ id: item.id, label: item.displayName });
      }
    };

    const listHandler = (delta: number) => {
      const newIndex = (selectedIndex + delta) % items.length;
      const item = items[newIndex];
      if (typeof item === "string" && isValidEmail(item)) {
        setSelectedIndex(items.length - 1);
      } else {
        setSelectedIndex(newIndex);
      }
    };

    const enterHandler = () => {
      selectItem(selectedIndex);
    };

    useEffect(() => {
      const startIndex = 0;
      const item = items[startIndex];
      if (typeof item === "string" && !isValidEmail(item)) {
        setSelectedIndex(-1);
        return;
      }

      setSelectedIndex(startIndex);
    }, [items]);

    useImperativeHandle(ref, () => ({
      onKeyDown: ({ event }) => {
        if (event.key === "ArrowUp") {
          listHandler(-1);
          return true;
        }

        if (event.key === "ArrowDown") {
          listHandler(1);
          return true;
        }

        if (event.key === "Enter") {
          enterHandler();
          return true;
        }

        if (event.key === "Tab") {
          enterHandler();
          return true;
        }

        return false;
      },
    }));

    if (!items || items.length === 0) {
      return null;
    }

    const users = items.filter((item) => typeof item !== "string");
    const newUserEntry = items.find((item) => typeof item === "string");

    return (
      <div className="bg-background rounded-md shadow-lg border border-border overflow-hidden p-2 w-[18.25rem]">
        <p className="px-2 text-sm text-muted-foreground">People</p>
        <ul className="mt-2">
          {users.map((user) => (
            <li
              key={user.id}
              onClick={() => command({ id: user.id, label: user.displayName })}
              className={cn(
                "flex items-center hover:elevated cursor-pointer text-sm text-foreground p-2",
                selectedIndex === items.indexOf(user) ? "bg-elevated" : "",
              )}
            >
              <Avatar className="w-6 h-6 mr-2">
                <AvatarImage
                  alt="Profile icon"
                  src={user.picture || undefined}
                />
                <AvatarFallback className="text-background bg-primary">
                  {getInitials(user.name)}
                </AvatarFallback>
              </Avatar>

              <ListName user={user} />
            </li>
          ))}
          {newUserEntry && (
            <li
              key="add-person"
              onClick={() => command({ id: INVITE_ID, label: newUserEntry })}
              className={cn(
                "flex items-center hover:elevated cursor-pointer text-sm text-foreground p-2",
                selectedIndex === items.indexOf(newUserEntry)
                  ? "bg-elevated"
                  : "",
              )}
            >
              <UserCircleIcon className="w-6 h-6 mr-2 min-w-6" />
              <span className="break-all">
                {`${!isValidEmail(newUserEntry) ? "Add valid email" : "Invite"}: ${newUserEntry}
                `}
              </span>
            </li>
          )}
        </ul>
      </div>
    );
  },
);

const ListName = ({ user }: { user: User }) => {
  const namesAreDifferent = user.displayName !== user.name;
  return (
    <span className="text-ellipsis overflow-hidden text-nowrap">
      {namesAreDifferent
        ? `${user.displayName} • ${user.name}`
        : user.displayName}
    </span>
  );
};

MentionList.displayName = "MentionList";
