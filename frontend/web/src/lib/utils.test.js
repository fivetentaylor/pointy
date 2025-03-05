import { expect, describe, it } from "vitest";
import { timeAgo, abbreviateName } from "./utils";

describe("timeAgo", () => {
  const now = new Date();

  it('should return "Just now" for dates less than a minute ago', () => {
    const date = new Date(now.getTime() - 10 * 1000); // 10 seconds ago
    expect(timeAgo(date)).toBe("just now");
  });

  it("should return in minutes format", () => {
    const date = new Date(now.getTime() - 20 * 60 * 1000); // 20 minutes ago
    expect(timeAgo(date)).toBe("20m");
  });

  it("should handle singular minute correctly", () => {
    const date = new Date(now.getTime() - 1 * 60 * 1000); // 1 minute ago
    expect(timeAgo(date)).toBe("1m");
  });

  it("should return in hours format", () => {
    const date = new Date(now.getTime() - 3 * 60 * 60 * 1000); // 3 hours ago
    expect(timeAgo(date)).toBe("3h");
  });

  it("should handle singular hour correctly", () => {
    const date = new Date(now.getTime() - 1 * 60 * 60 * 1000); // 1 hour ago
    expect(timeAgo(date)).toBe("1h");
  });

  it("should return in days format", () => {
    const date = new Date(now.getTime() - 5 * 24 * 60 * 60 * 1000); // 5 days ago
    expect(timeAgo(date)).toBe("5d");
  });

  it("should handle singular day correctly", () => {
    const date = new Date(now.getTime() - 1 * 24 * 60 * 60 * 1000); // 1 day ago
    expect(timeAgo(date)).toBe("1d");
  });

  it("should return in months format", () => {
    const date = new Date(now.getTime() - 2 * 30 * 24 * 60 * 60 * 1000); // 2 months ago
    expect(timeAgo(date)).toBe("2mo");
  });

  it("should handle singular month correctly", () => {
    const date = new Date(now.getTime() - 1 * 30 * 24 * 60 * 60 * 1000); // 1 month ago
    expect(timeAgo(date)).toBe("1mo");
  });

  it("should return in years format", () => {
    const date = new Date(now.getTime() - 3 * 365 * 24 * 60 * 60 * 1000); // 3 years ago
    expect(timeAgo(date)).toBe("3y");
  });

  it("should handle singular year correctly", () => {
    const date = new Date(now.getTime() - 1 * 365 * 24 * 60 * 60 * 1000); // 1 year ago
    expect(timeAgo(date)).toBe("1y");
  });
});

describe("abbreviateName", () => {
  it('should return "?" for empty or whitespace-only names', () => {
    expect(abbreviateName("")).toBe("?");
    expect(abbreviateName("    ")).toBe("?");
  });

  it("should return the capitalized first letter for single-word names", () => {
    expect(abbreviateName("John")).toBe("J");
    expect(abbreviateName("alice")).toBe("A");
  });

  it("should return the capitalized first letter of the first and last word for multi-word names", () => {
    expect(abbreviateName("John Doe")).toBe("JD");
    expect(abbreviateName("alice wonderland")).toBe("AW");
    expect(abbreviateName("first middle last")).toBe("FL");
  });

  it("should skip non-alphanumeric words", () => {
    expect(abbreviateName("123 John")).toBe("J");
    expect(abbreviateName("John 456")).toBe("J");
    expect(abbreviateName("123 John 789")).toBe("J");
    expect(abbreviateName("123 John Doe")).toBe("JD");
    expect(abbreviateName("John 123 Doe")).toBe("JD");
  });

  it('should return "?" if the name consists of non-alphabetic words only', () => {
    expect(abbreviateName("123 456")).toBe("?");
    expect(abbreviateName("!!! ???")).toBe("?");
  });
});
