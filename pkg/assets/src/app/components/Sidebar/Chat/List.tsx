import React, { useRef, useEffect } from "react";
import { MessageFieldsFragment } from "@/__generated__/graphql";
import { RevisoUserID } from "@/constants";

import Message from "./Message";
import { ScrollArea } from "@/components/ui/scroll-area";

type ListProps = {
  loading: boolean;
  messages: MessageFieldsFragment[];
};

const List = ({ messages }: ListProps) => {
  const messagesEndRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView();
    }
  }, [messages.length, messagesEndRef]);

  const findPreviousMessageFromUser = (
    messages: MessageFieldsFragment[],
    currentIndex: number,
    userId: string,
  ) => {
    for (let i = currentIndex - 1; i >= 0; i--) {
      if (messages[i].user.id === userId) {
        return messages[i];
      }
    }
    return null;
  };

  return (
    <ScrollArea className="MessageContainerScrollArea flex-1 overflow-y-auto mr-[-0.52rem] pr-3">
      <div className="flex-grow flex flex-col h-full">
        {messages.map((message, idx) => {
          const previousMessage = findPreviousMessageFromUser(
            messages,
            idx,
            RevisoUserID,
          );

          if (idx === messages.length - 1) {
            return (
              <Message
                key={message.id}
                message={message}
                previousRevisoMessage={previousMessage}
                ref={messagesEndRef}
              />
            );
          }
          return (
            <Message
              key={message.id}
              message={message}
              previousRevisoMessage={previousMessage}
            />
          );
        })}
      </div>
    </ScrollArea>
  );
};

export default List;
