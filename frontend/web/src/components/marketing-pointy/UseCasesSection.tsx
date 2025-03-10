import { LucideIcon } from "lucide-react";
import { CenteredLayout } from "./CenteredLayout";

interface UseCaseItem {
  title: string;
  description: string;
  icon: LucideIcon;
}

interface UseCasesSectionProps {
  items: UseCaseItem[];
  type: "academic" | "professional";
}

export function UseCasesSection({ items, type }: UseCasesSectionProps) {
  return (
    <section
      className="bg-zinc-50 max-w-[100vw] w-[100vw] pb-8 pt-[1.625rem] sm:pt-[3.8125rem] sm:pb-[3.25rem]"
      id="use-cases"
    >
      <CenteredLayout>
        <h2 className="text-[clamp(1.25rem,calc(1.25rem+((1vw-0.2rem)*2.976)),2.5rem)] leading-[1.1] font-bold mb-3 sm:mb-2 text-center">
          10x your {type} writing
        </h2>
        <p className="text-[clamp(1rem,calc(1rem+((1vw-0.2rem)*0.893)),1.375rem)] leading-[1.3] font-light text-center text-muted-foreground sm:mb-12">
          Just some of what Pointy can help you with
        </p>
        <div className="grid md:grid-cols-3 sm:gap-8">
          {items.map(({ title, description, icon: Icon }) => (
            <div
              key={title}
              className="pt-6 p-0 md:pl-6 text-center md:text-left mx-auto md:max-w-auto"
            >
              <div className="mb-4 flex justify-center md:justify-start">
                <Icon className="w-9 h-9 text-primary" />
              </div>
              <h3 className="text-[clamp(1rem,calc(1rem+((1vw-0.2rem)*1.19)),1.5rem)] leading-[1.3] font-semibold mb-2">
                {title}
              </h3>
              <p className="text-[clamp(0.875rem,calc(0.875rem+((1vw-0.2rem)*0.893)),1.25rem)] leading-[1.4] font-light text-muted-foreground">
                {description}
              </p>
            </div>
          ))}
        </div>
        <p className="text-[clamp(1rem,calc(1rem+((1vw-0.2rem)*0.893)),1.375rem)] leading-[1.3] text-center text-muted-foreground mt-8 font-normal">
          &hellip;and so much more
        </p>
      </CenteredLayout>
    </section>
  );
}
