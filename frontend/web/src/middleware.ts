import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

function matchesPath(pathname: string, patterns: string[]): boolean {
  return patterns.some((pattern) => {
    // Exact match
    if (pathname === pattern) return true;

    // Pattern ends with /* for wildcard matching
    if (pattern.endsWith("/*")) {
      const basePattern = pattern.slice(0, -2);
      return pathname.startsWith(basePattern);
    }

    return false;
  });
}

export function middleware(request: NextRequest) {
  const hostname = request.headers.get("host");
  const { pathname } = request.nextUrl;

  const pointyOnlyPaths = ["/pointy/*"];
  const revisoOnlyPaths = ["/work", "/schools", "/students", "/legal/*"];

  if (hostname?.includes("localhost")) {
    return NextResponse.next();
  }

  if (pathname === "/" && hostname?.includes("pointy.ai")) {
    const url = request.nextUrl.clone();
    url.pathname = "/pointy";
    return NextResponse.redirect(url);
  }

  if (matchesPath(pathname, pointyOnlyPaths)) {
    if (!hostname?.includes("pointy.ai")) {
      return new NextResponse("Not authorized", { status: 403 });
    }
  }

  if (matchesPath(pathname, revisoOnlyPaths)) {
    if (!hostname?.includes("revi")) {
      return new NextResponse("Not authorized", { status: 403 });
    }
  }

  return NextResponse.next();
}
