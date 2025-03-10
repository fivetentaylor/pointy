"use client";
import { Typewriter } from "react-simple-typewriter";

export const TypeWriterCopy = function () {
  return (
    <p className="max-w-[42rem] leading-normal text-gray sm:text-xl sm:leading-8">
      <Typewriter
        words={["Break through the blank page"]}
        cursor
        cursorBlinking
      />
    </p>
  );
};
