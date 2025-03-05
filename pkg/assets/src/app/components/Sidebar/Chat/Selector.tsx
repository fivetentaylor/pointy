import React from "react";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useCurrentUserContext } from "@/contexts/CurrentUserContext";
import { Thread } from "@/__generated__/graphql";

type SelectorProps = {
  threadId: string;
  loading: boolean;
  threads: Thread[] | null;
  onSelectThread: (threadId: string) => void;
};

interface GroupedThreads {
  user: User;
  threads: Thread[];
}

const Selector = ({
  threadId,
  loading,
  threads,
  onSelectThread,
}: SelectorProps) => {
  const handleValueChange = (value: string) => {
    if (value) {
      onSelectThread(value);
    }
  };

  const { currentUser, loading: loadingMe } = useCurrentUserContext();

  if (loadingMe || (loading && !threads)) {
    return (
      <div className="mr-3 mb-2">
        <Select disabled>
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Loading" />
          </SelectTrigger>
          <SelectContent></SelectContent>
        </Select>
      </div>
    );
  }

  if (!threads) {
    return (
      <div className="mr-3 mb-2">
        <Select disabled>
          <SelectTrigger className="w-full">
            <SelectValue placeholder="No Threads" />
          </SelectTrigger>
          <SelectContent></SelectContent>
        </Select>
      </div>
    );
  }

  const groupedByUser: Record<string, GroupedThreads> = threads.reduce(
    (acc, thread) => {
      const userId = thread.user.id;
      if (!acc[userId]) {
        acc[userId] = {
          user: thread.user,
          threads: [],
        };
      }
      acc[userId].threads.push(thread);
      return acc;
    },
    {} as Record<string, GroupedThreads>,
  );

  const showingOthersThreads = Object.keys(groupedByUser).some(
    (userId: string) => userId !== currentUser?.id,
  );

  return (
    <div className="mr-3 mb-2 flex gap-2">
      <Select value={threadId} onValueChange={handleValueChange}>
        <SelectTrigger className="w-full [&>span]:text-left [&>span]:truncate">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {Object.values(groupedByUser).map(({ user, threads }) => (
            <div key={user.id}>
              {showingOthersThreads && (
                <SelectItem key={user.id} value={user.id} disabled>
                  <div className="flex items-center text-xs">
                    {user.id === currentUser?.id
                      ? "Your Topics"
                      : user.name + "'s Topics"}
                  </div>
                </SelectItem>
              )}
              {threads.map((thread) => (
                <SelectItem key={thread.id} value={thread.id}>
                  <div className="flex items-center">{thread.title}</div>
                </SelectItem>
              ))}
            </div>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
};

export default Selector;
