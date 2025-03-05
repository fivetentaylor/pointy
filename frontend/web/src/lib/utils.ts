import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { format, getYear, parseISO } from "date-fns";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const formatHumanReadableDate = (dateString: string) => {
  const parsedDate = parseISO(dateString);
  const currentYear = getYear(new Date());
  const dateYear = getYear(parsedDate);
  const currentMonth = getYear(new Date());
  const dateMonth = getYear(parsedDate);

  // Check if the year of the date matches the current year
  if (dateYear === currentYear) {
    if (dateMonth === currentMonth) {
      return format(parsedDate, "h:mma");
    }
    return format(parsedDate, "MMM dd, h:mma");
  } else {
    return format(parsedDate, "MMM dd, yyyy");
  }
};

export const abbreviateName = (name: string) => {
  // If the name is empty or only consists of whitespace, return "?"
  if (!name.trim()) return "?";

  // Filter words with only alphanumeric characters and transform them
  let alphaWords = name.split(/\s+/).filter((word) => /^[a-z]+$/i.test(word));

  // If there's only one valid word, return its first letter capitalized
  if (alphaWords.length === 1) return alphaWords[0][0].toUpperCase();

  // If there's more than one valid word, return the first letter of the first
  // word and the first letter of the last word, both capitalized
  if (alphaWords.length >= 2) {
    return (
      alphaWords[0][0].toUpperCase() +
      alphaWords[alphaWords.length - 1][0].toUpperCase()
    );
  }

  // If there are no valid alpha words, return "?"
  return "?";
};

export const formatShortDate = (dateString: string) => {
  const inputDate = new Date(dateString);
  const currentDate = new Date();

  function getMonthName(monthIndex: number) {
    const months = [
      "Jan",
      "Feb",
      "Mar",
      "Apr",
      "May",
      "Jun",
      "Jul",
      "Aug",
      "Sep",
      "Oct",
      "Nov",
      "Dec",
    ];
    return months[monthIndex];
  }

  function isSameDay(date1: Date, date2: Date) {
    return (
      date1.getDate() === date2.getDate() &&
      date1.getMonth() === date2.getMonth() &&
      date1.getFullYear() === date2.getFullYear()
    );
  }

  function formatTime(date: Date) {
    return date.toLocaleTimeString([], { hour: "numeric", minute: "2-digit" });
  }

  if (isSameDay(inputDate, currentDate)) {
    return formatTime(inputDate);
  } else {
    return `${getMonthName(inputDate.getMonth())} ${inputDate.getDate()}`;
  }
};

export const timeAgo = (isoDate: string): string => {
  const eventDate = new Date(isoDate);
  const currentDate = new Date();
  const diffInSeconds = Math.floor(
    (currentDate.getTime() - eventDate.getTime()) / 1000,
  );

  if (diffInSeconds < 60) {
    return "just now";
  } else if (diffInSeconds < 3600) {
    return `${Math.floor(diffInSeconds / 60)}m`;
  } else if (diffInSeconds < 86400) {
    return `${Math.floor(diffInSeconds / 3600)}h`;
  } else if (diffInSeconds < 2592000) {
    return `${Math.floor(diffInSeconds / 86400)}d`;
  } else if (diffInSeconds < 31536000) {
    return `${Math.floor(diffInSeconds / 2592000)}mo`;
  } else {
    return `${Math.floor(diffInSeconds / 31536000)}y`;
  }
};

export const timeAgoLong = (isoDate: string): string => {
  const eventDate = new Date(isoDate);
  const currentDate = new Date();
  const diffInSeconds = Math.floor(
    (currentDate.getTime() - eventDate.getTime()) / 1000,
  );

  if (diffInSeconds < 60) {
    return "just now";
  } else if (diffInSeconds < 3600) {
    const val = Math.floor(diffInSeconds / 60);
    return `${val} ${val === 1 ? "minute" : "minutes"} ago`;
  } else if (diffInSeconds < 86400) {
    const val = Math.floor(diffInSeconds / 3600);
    return `${val} ${val === 1 ? "hour" : "hours"} ago`;
  } else if (diffInSeconds < 2592000) {
    const val = Math.floor(diffInSeconds / 86400);
    return `${val} ${val === 1 ? "day" : "days"} ago`;
  } else if (diffInSeconds < 31536000) {
    const val = Math.floor(diffInSeconds / 2592000);
    return `${val} ${val === 1 ? "month" : "months"} ago`;
  } else {
    const val = Math.floor(diffInSeconds / 31536000);
    return `${val} ${val === 1 ? "year" : "years"} ago`;
  }
};

export function getInitials(name: string) {
  return name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase();
}

export function getFirstName(name: string) {
  return name.split(" ")[0];
}

export const validateEmail = function (email: string): boolean {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return re.test(email.toLowerCase());
};

export const isSafari = function (): boolean {
  const userAgent = navigator.userAgent;
  const safari = /^((?!chrome|android).)*safari/i.test(userAgent);
  return safari;
};
