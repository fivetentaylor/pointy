import { RevisoLogo } from "@/components/ui/RevisoLogo";
import { Button } from "@/components/ui/button";
import { SearchXIcon } from "lucide-react";
import Link from "next/link";

export default function NotFound() {
  return (
    <>
      <section className="relative flex max-sm:flex-col w-full justify-between items-center border-b py-4 sm:pt-7 sm:pb-[1.875rem] pl-8 pr-8">
        <div className="scale-[0.85] mt-[-0.25rem]">
          <RevisoLogo />
        </div>
        <h2 className="absolute left-1/2 transform -translate-x-1/2 max-sm:mt-4 text-2xl leading-[2.375rem] font-bold cursor-pointer text-center w-full px-8">
          Not Found
        </h2>
      </section>
      <div className="h-screen w-full flex items-center justify-center mt-[-3rem]">
        <div className="flex flex-col mx-auto">
          <div className="w-12 h-12 mb-4 bg-zinc-100 rounded-full mx-auto flex items-center justify-center">
            <SearchXIcon className="w-6 h-6 stroke-zinc-500" />
          </div>

          <p className="text-center text-foreground text-base font-semibold leading-normal">
            We&apos;re sorry, that page was not found.
          </p>
          <Button
            asChild
            variant="secondary"
            size="sm"
            className="mt-4 mx-auto"
          >
            <Link href="/documents">Go back to home</Link>
          </Button>
        </div>
      </div>
    </>
  );
}
