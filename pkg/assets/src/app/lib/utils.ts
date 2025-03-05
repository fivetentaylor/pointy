import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function getInitials(name: string) {
  return name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase();
}

export function parseId(id: string): [string, number] {
  const parts = id.split("_");
  if (parts.length !== 2) {
    throw new Error("Invalid id");
  }
  const lamport = parseInt(parts[1]);
  return [parts[0], lamport];
}

export const formatShortDate = (inputDate: Date) => {
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

export function timeAgo(timestamp: string): string {
  const now: Date = new Date();
  const then: Date = new Date(timestamp);
  const seconds: number = Math.floor((now.getTime() - then.getTime()) / 1000);

  let interval: number = Math.floor(seconds / 31536000);

  if (interval > 1) {
    return interval + "y";
  }
  interval = Math.floor(seconds / 2592000);
  if (interval > 1) {
    return interval + "mo";
  }
  interval = Math.floor(seconds / 86400);
  if (interval >= 1) {
    return interval + "d";
  }
  interval = Math.floor(seconds / 3600);
  if (interval >= 1) {
    return interval + "h";
  }
  interval = Math.floor(seconds / 60);
  if (interval >= 1) {
    return interval + "m";
  }
  return "<1m";
}

export const timeSpecific = (isoDate: string): string => {
  const date = new Date(isoDate);

  // Options to format the date and time
  const options: Intl.DateTimeFormatOptions = {
    weekday: "short", // Thu
    year: "numeric", // 2024
    month: "short", // Aug
    day: "numeric", // 22
    hour: "numeric", // 4
    minute: "numeric", // 01
    hour12: true, // AM/PM format
  };

  // Getting the time zone abbreviation
  const timeZoneOptions: Intl.DateTimeFormatOptions = {
    timeZoneName: "short",
  };

  const timeZoneStr = date
    .toLocaleString("en-US", timeZoneOptions)
    .split(" ")
    .pop();

  // Formatting the date and time with the given options
  const formattedDate = date.toLocaleString("en-US", options);

  // Returning the formatted date with timezone abbreviation
  return `${formattedDate} ${timeZoneStr}`;
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

export function arrayToSentence(arr: string[]): string {
  if (arr.length === 0) return "";
  if (arr.length === 1) return arr[0];
  if (arr.length === 2) return arr.join(" and ");

  const arrCopy = [...arr];
  const lastItem = arrCopy.pop();
  return arrCopy.join(", ") + " and " + lastItem;
}

export function formatLocalTime(
  isoDate: string = new Date().toISOString(),
): string {
  const date = new Date(isoDate);
  return date
    .toLocaleString("en-US", {
      hour: "numeric",
      minute: "2-digit",
      hour12: true,
    })
    .toLowerCase();
}
