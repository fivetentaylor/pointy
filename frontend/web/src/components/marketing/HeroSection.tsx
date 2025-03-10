import { APP_HOST } from "@/lib/urls";
import { Button } from "../ui/button";
import { ArrowRight } from "lucide-react";

interface HeroSectionProps {
  video: string;
  title: string;
  description: string;
}

function AnnouncementBanner() {
  return (
    <div className="flex justify-center mb-8">
      <a
        href="/winding-down"
        className="group inline-flex items-center gap-2 px-4 sm:px-5 py-2 bg-foreground text-background rounded-full hover:bg-foreground/90 transition-colors"
      >
        <span className="text-base font-normal">
          Reviso is winding&nbsp;down:
          A&nbsp;letter&nbsp;from&nbsp;our&nbsp;team
        </span>
        <ArrowRight className="w-4 h-4 shrink-0 transition-transform group-hover:translate-x-0.5" />
      </a>
    </div>
  );
}

export function HeroSection({ video, title, description }: HeroSectionProps) {
  return (
    <section className="pt-6 sm:pt-0 sm:py-8 md:pt-12 text-center mb-2.5">
      <AnnouncementBanner />
      <h1 className="text-[clamp(1.5rem,calc(1.5rem+((1vw-0.2rem)*3.869)),3.125rem)] leading-[1.1] tracking-[-0.09563rem] font-bold mb-4">
        {title}
      </h1>
      <p className="text-[clamp(1rem,calc(1rem+((1vw-0.2rem)*0.893)),1.375rem)] leading-[1.3] font-light text-foreground max-w-2xl mx-auto whitespace-pre-line">
        {description}
      </p>
      <div className="aspect-video bg-muted rounded-lg lg:mx-[-5%] xl:mx-[-10rem]">
        <div className="w-full h-full flex items-center justify-center text-muted-foreground">
          <video src={video} autoPlay loop muted playsInline />
        </div>
      </div>
    </section>
  );
}
