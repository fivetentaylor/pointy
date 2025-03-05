import PointyLogo from "@/components/ui/PointyLogo";
import Link from "next/link";
import { Twitter, Linkedin } from "lucide-react";

export function MarketingFooter() {
  return (
    <footer>
      <div className="max-w-[68.75rem] mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <PointyLogo />
      </div>
      <div className="border-t">
        <div className="max-w-[68.75rem] mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="text-sm text-muted-foreground">
            <div className="flex gap-4">
              <Link href="/pointy/legal/privacy" className="hover:underline">
                Privacy Policy
              </Link>
              <Link
                href="/pointy/legal/terms-of-service"
                className="hover:underline"
              >
                Terms of Service
              </Link>
            </div>
          </div>
          <div className="flex gap-4 mt-6">
            <Link
              href="https://twitter.com/writewithreviso"
              className="text-muted-foreground hover:text-foreground transition-colors"
              aria-label="Follow Pointy on X"
            >
              <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <g clip-path="url(#clip0_313_5595)">
                  <path
                    d="M13.3076 10.4643L20.3808 2H18.7046L12.563 9.34942L7.65769 2H2L9.41779 13.1136L2 21.9897H3.67621L10.1619 14.2285L15.3423 21.9897H21L13.3072 10.4643H13.3076ZM11.0118 13.2115L10.2602 12.1049L4.28017 3.29901H6.85474L11.6807 10.4056L12.4323 11.5123L18.7054 20.7498H16.1309L11.0118 13.212V13.2115Z"
                    fill="#71717A"
                  />
                </g>
                <defs>
                  <clipPath id="clip0_313_5595">
                    <rect
                      width="19"
                      height="20"
                      fill="white"
                      transform="translate(2 2)"
                    />
                  </clipPath>
                </defs>
              </svg>
            </Link>
            <Link
              href="https://linkedin.com/company/writewithreviso"
              className="text-muted-foreground hover:text-foreground transition-colors"
              aria-label="Follow Pointy on LinkedIn"
            >
              <Linkedin size={24} />
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
