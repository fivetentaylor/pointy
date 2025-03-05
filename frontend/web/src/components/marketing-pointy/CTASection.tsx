import { Button } from "@/components/ui/button";
import { APP_HOST } from "@/lib/urls";
import { CenteredLayout } from "./CenteredLayout";

export function CTASection() {
  return (
    <CenteredLayout variant="wide">
      <section className="relative py-24 mt-[3.25rem] -mx-4 sm:rounded-[2rem] overflow-hidden lg:h-[32.625rem] lg:bg-cover bg-center sm:mb-24 bg-[linear-gradient(98.67deg,_#DCFCE7_5.68%,_#8B5CF6_154.19%)] sm:mx-0 lg:bg-[url('/images/ready-to-get-started.png')]">
        <div className="max-w-[68.75rem] mx-auto px-4 lg:px-8">
          <div className="lg:max-w-md text-center lg:text-left lg:mt-[3.3125rem]">
            <h2 className="text-[clamp(1.25rem,calc(1.25rem+((1vw-0.2rem)*2.976)),2.5rem)] leading-[1.1] font-bold mb-4">
              Ready to get started?
            </h2>
            <p className="text-[clamp(1rem,calc(1rem+((1vw-0.2rem)*0.893)),1.375rem)] leading-[1.3] font-light mb-8">
              Join our private beta for a limited time.
            </p>
            <Button
              asChild
              size="lg"
              className="bg-[#7C3AED] hover:bg-[#6D28D9] text-white px-8"
            >
              <a href={`${APP_HOST}/login`}>Join our private beta</a>
            </Button>
          </div>
        </div>
      </section>
    </CenteredLayout>
  );
}
