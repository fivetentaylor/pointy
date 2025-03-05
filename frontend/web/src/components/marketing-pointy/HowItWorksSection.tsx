import { cn } from "@/lib/utils";
import Image from "next/image";

interface HowItWorksItem {
  title: string;
  description: string;
  image: string;
  height?: number;
  imageStyle?: {
    width?: number;
    height?: number;
    maxWidthRem?: number;
    marginLeft?: boolean;
  };
  centerImage?: boolean;
  topAlign?: boolean;
}

interface HowItWorksSectionProps {
  items: HowItWorksItem[];
}

export function HowItWorksSection({ items }: HowItWorksSectionProps) {
  // Split items into left and right columns
  const leftColumnItems = items.filter((_, i) => i % 2 === 0);
  const rightColumnItems = items.filter((_, i) => i % 2 === 1);

  const ItemCard = ({ item }: { item: HowItWorksItem }) => (
    <div
      className={`item-card bg-zinc-50 rounded-lg flex flex-col w-full h-auto min-h-[24.125rem] lg:min-h-0 `}
      data-custom-height={item.height ? "true" : undefined}
      style={
        item.height
          ? ({ "--custom-height": `${item.height}rem` } as React.CSSProperties)
          : undefined
      }
    >
      <div
        className={`px-6 sm:px-10 pt-6 sm:pt-10 ${
          item.centerImage ? "flex flex-col flex-1" : ""
        }`}
      >
        <div className={item.centerImage ? "mb-4 sm:mb-6" : "mb-4 sm:mb-0"}>
          <h3 className="text-[clamp(1rem,calc(1rem+((1vw-0.2rem)*1.19)),1.5rem)] leading-[1.3] font-semibold mb-2 sm:mb-2 tracking-[-0.045rem]">
            {item.title}
          </h3>
          <p className="text-muted-foreground text-[clamp(0.875rem,calc(0.875rem+((1vw-0.2rem)*0.893)),1.25rem)] leading-[1.4] font-light">
            {item.description}
          </p>
        </div>
        {item.image && item.centerImage && (
          <div className="flex-1 flex items-center justify-center">
            <Image
              src={item.image}
              alt={item.title}
              width={item.imageStyle?.width || 400}
              height={item.imageStyle?.height || 400}
              className="rounded-lg w-full h-auto"
              style={{ maxWidth: `${item.imageStyle?.maxWidthRem}rem` }}
            />
          </div>
        )}
      </div>
      {item.image && !item.centerImage && (
        <div className="flex-1 flex flex-col">
          <div
            className={cn(
              item.topAlign ? "mt-4" : "mt-auto",
              item.imageStyle?.marginLeft ? "ml-8" : "",
            )}
          >
            <Image
              src={item.image}
              alt={item.title}
              width={item.imageStyle?.width || 400}
              height={item.imageStyle?.height || 400}
              className="rounded-lg w-full h-auto"
              style={{ maxWidth: `${item.imageStyle?.maxWidthRem}rem` }}
            />
          </div>
        </div>
      )}
    </div>
  );

  return (
    <section className="" id="how-it-works">
      <h2 className="text-[clamp(1.25rem,calc(1.25rem+((1vw-0.2rem)*2.976)),2.5rem)] leading-[1.1] tracking-[-0.075rem] font-bold mb-1 sm:mb-2 text-center">
        How it works
      </h2>

      {/* Mobile: Stack all items vertically */}
      <div className="grid grid-cols-1 gap-6 md:hidden">
        {items.map((item, i) => (
          <ItemCard key={i} item={item} />
        ))}
      </div>

      {/* Desktop: Two flex columns */}
      <div className="hidden md:flex gap-6">
        {/* Left column */}
        <div className="flex-1 flex flex-col gap-6">
          {leftColumnItems.map((item, i) => (
            <ItemCard key={i} item={item} />
          ))}
        </div>
        {/* Right column */}
        <div className="flex-1 flex flex-col gap-6">
          {rightColumnItems.map((item, i) => (
            <ItemCard key={i} item={item} />
          ))}
        </div>
      </div>
    </section>
  );
}
