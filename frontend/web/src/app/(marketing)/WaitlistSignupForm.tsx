"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { FormEvent, useState } from "react";

import API from "@/lib/api";
import { validateEmail } from "@/lib/utils";
import { AlertTriangle } from "lucide-react";
import posthog from "posthog-js";
import { WAITLIST_SIGNUP } from "@/lib/events";
import { analytics } from "@/lib/segment";

export const WaitlistSignupForm = function () {
  const [emailError, setEmailError] = useState("");
  const [email, setEmail] = useState("");
  const [isSignedUp, setIsSignedUp] = useState(false);
  const onSubmit = async function (event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!validateEmail(email)) {
      setEmailError("Invalid email format");
      return;
    }

    setEmailError("");

    try {
      let response = await API.post("/waitlist", {
        email,
      });

      if (response.status !== 200) {
        console.error("Network response was not ok");
      }
      analytics.track(WAITLIST_SIGNUP, { email });
      setIsSignedUp(true);
      window.location.href = `https://tally.so/r/3No8kG?email=${email}`;
    } catch (error: any) {
      console.error(error.message);
    }
  };

  const runEmailCheck = (value: string) => {
    if (value !== "" && !validateEmail(value)) {
      setEmailError("Invalid email format");
    } else {
      setEmailError("");
    }
  };

  return (
    <form
      className="pt-4 sm:pt-6 flex max-sm:flex-col max-sm:items-center"
      onSubmit={onSubmit}
    >
      <div className="flex flex-col flex-grow sm:min-w-[23.5rem] max-w-[23.5rem]">
        <Input
          id="email"
          type="email"
          placeholder="name@example.com"
          value={email}
          className={`max-sm:mb-4  ${emailError && "border-red-500"}`}
          onBlur={() => {
            runEmailCheck(email);
          }}
          onChange={(e) => {
            setEmail(e.target.value);

            if (emailError !== "") {
              runEmailCheck(e.target.value);
            }
          }}
        />
        {emailError && (
          <div className="text-red-500 text-left pt-1 flex align-middle">
            <AlertTriangle className="inline-block h-5 w-5 mr-2" />
            {emailError}
          </div>
        )}
      </div>

      <Button
        className="ml-3 bg-primary hover:bg-primary/90 min-w-[8.375rem]"
        type="submit"
      >
        {isSignedUp ? "On the list!" : "Join waitlist"}
      </Button>
    </form>
  );
};
