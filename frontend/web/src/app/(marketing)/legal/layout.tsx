import { RevisoLogo } from "@/components/ui/RevisoLogo";

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex flex-col pt-4 px-9">
      <nav className="flex items-center">
        <div className="flex-1">
          <a href="/">
            <RevisoLogo />
          </a>
        </div>
      </nav>
      <section className="flex flex-col items-center justify-center text-center pt-[8rem] max-sm:pt-[4rem]">
        <div className="text-left max-w-[43rem]">{children}</div>
      </section>
    </div>
  );
}
