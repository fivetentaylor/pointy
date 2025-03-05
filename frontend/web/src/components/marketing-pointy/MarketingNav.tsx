"use client";

import { Button } from "@/components/ui/button";
import PointyLogo from "@/components/ui/PointyLogo";
import { APP_HOST } from "@/lib/urls";
import Link from "next/link";
import { cn } from "@/lib/utils";
import { Menu, X } from "lucide-react";
import { useState } from "react";
import { CenteredLayout } from "./CenteredLayout";

interface MobileMenuProps {
  pathname: string;
}

function MobileMenu({ pathname }: MobileMenuProps) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <button className="sm:hidden p-2" onClick={() => setIsOpen(true)}>
        <Menu className="h-6 w-6" />
      </button>

      {isOpen && (
        <div className="fixed inset-0 bg-white/80 backdrop-blur-sm z-50">
          <div className="fixed inset-x-4 top-4 bg-white rounded-xl border border-[hsla(240,6%,90%,1)] shadow-[0px_10px_15px_-3px_hsla(0,0%,0%,0.1),0px_4px_6px_-4px_hsla(220,43%,11%,0.1)] p-6">
            <div className="flex justify-between items-center mb-8">
              <PointyLogo />
              <button onClick={() => setIsOpen(false)} className="p-2">
                <X className="h-6 w-6" />
              </button>
            </div>
            <div className="flex flex-col gap-6">
              {[
                { label: "For students", href: "/pointy/students" },
                { label: "For schools", href: "/pointy/schools" },
                { label: "For creatives", href: "/pointy/creatives" },
              ].map(({ label, href }) => (
                <Link
                  key={href}
                  href={href}
                  className="text-xl leading-[1.5125rem] font-medium"
                  onClick={() => setIsOpen(false)}
                >
                  {label}
                </Link>
              ))}
            </div>
            <div className="flex flex-col gap-4 mt-8">
              <Link
                href="https://twitter.com/writewithreviso"
                className="text-base leading-[1.28rem] font-normal"
                onClick={() => setIsOpen(false)}
              >
                Follow Pointy on X
              </Link>
              <Link
                href="https://linkedin.com/company/writewithreviso"
                className="text-base leading-[1.28rem] font-normal"
                onClick={() => setIsOpen(false)}
              >
                Follow Pointy on LinkedIn
              </Link>
            </div>
            <div className="mt-auto pt-8">
              <Button
                className="w-full mb-4 bg-primary hover:bg-primary/90"
                asChild
              >
                <a href={`${APP_HOST}/login`}>Sign up</a>
              </Button>
              <div className="text-center">
                Existing user?{" "}
                <a href={`${APP_HOST}/login`} className="text-primary">
                  Sign in
                </a>
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

// Server Component
interface MarketingNavProps {
  pathname: string;
}

export function MarketingNav({ pathname }: MarketingNavProps) {
  const isActive = (path: string) => {
    if (path === "/pointy/creatives") {
      return pathname === "/pointy/creatives" || pathname === "/";
    }
    return pathname === path;
  };

  return (
    <div className="sticky top-0 z-50 w-full bg-white">
      <CenteredLayout variant="wide">
        <nav className="flex items-center py-3 sm:py-6 px-4 sm:px-0 -mx-4 sm:mx-0 border-b sm:border-none">
          <div className="flex-1">
            <Link href="/pointy">
              <PointyLogo />
            </Link>
          </div>
          <div className="flex-2 items-center justify-center gap-4 hidden sm:flex">
            <Button variant="ghost" asChild>
              <Link
                href="/pointy/students"
                className={cn(
                  "px-4 py-2",
                  isActive("/pointy/students") && "bg-muted rounded-full",
                )}
              >
                For Students
              </Link>
            </Button>
            <Button variant="ghost" asChild>
              <Link
                href="/pointy/schools"
                className={cn(
                  "px-4 py-2",
                  isActive("/pointy/schools") && "bg-muted rounded-full",
                )}
              >
                For Schools
              </Link>
            </Button>
            <Button variant="ghost" asChild>
              <Link
                href="/pointy/creatives"
                className={cn(
                  "px-4 py-2",
                  isActive("/pointy/creatives") && "bg-muted rounded-full",
                )}
              >
                For Creatives
              </Link>
            </Button>
          </div>
          <div className="flex-1 justify-end gap-4 hidden sm:flex">
            <Button variant="outline" asChild>
              <a href={`${APP_HOST}/login`}>Sign in</a>
            </Button>
            <Button asChild className="bg-primary hover:bg-primary/90">
              <a href={`${APP_HOST}/login`}>Sign up</a>
            </Button>
          </div>
          <MobileMenu pathname={pathname} />
        </nav>
      </CenteredLayout>
    </div>
  );
}
