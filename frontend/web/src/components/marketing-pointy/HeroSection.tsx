import { APP_HOST } from "@/lib/urls";
import { Button } from "../ui/button";

interface HeroSectionProps {
  video: string;
  title: string;
  description: string;
}

export function HeroSection({ video, title, description }: HeroSectionProps) {
  return (
    <section className="pt-6 sm:pt-0 sm:py-12 md:pt-20 text-center mb-2.5">
      <h1 className="text-[clamp(1.5rem,calc(1.5rem+((1vw-0.2rem)*3.869)),3.125rem)] leading-[1.1] tracking-[-0.09563rem] font-bold mb-4">
        {title}
      </h1>
      <p className="text-[clamp(1rem,calc(1rem+((1vw-0.2rem)*0.893)),1.375rem)] leading-[1.3] font-light text-foreground max-w-2xl mx-auto whitespace-pre-line">
        {description}
      </p>
      <div>
        <Button
          asChild
          size="lg"
          className="bg-primary hover:bg-primary/90 text-white px-8 my-5 sm:my-10"
        >
          <a href={`${APP_HOST}/login`}>Join our private beta</a>
        </Button>
      </div>
      <div className="aspect-video bg-muted rounded-lg lg:mx-[-5%] xl:mx-[-10rem]">
        <div className="w-full h-full flex items-center justify-center text-muted-foreground">
          <video src={video} autoPlay loop muted playsInline />
        </div>
      </div>
    </section>
  );
}
