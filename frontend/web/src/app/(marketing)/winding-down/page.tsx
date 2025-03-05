import { CenteredLayout } from "@/components/marketing/CenteredLayout";
import { HeroSection } from "@/components/marketing/HeroSection";

import { MarketingNav } from "@/components/marketing/MarketingNav";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";

const LetterContent = () => {
  return (
    <article className="max-w-2xl mx-auto py-12 px-4 font-sans">
      <Link
        href="/"
        className="inline-flex items-center gap-2 text-muted-foreground hover:text-foreground mb-8"
      >
        <ArrowLeft className="w-4 h-4" />
        Back to home
      </Link>

      <h1 className="text-4xl font-bold tracking-tight mb-8 !font-sans">
        The Next Chapter for Reviso
      </h1>

      <div className="prose prose-lg prose-slate max-w-none !font-sans">
        <p>
          We are excited to share some big news about how Reviso will be
          evolving.
        </p>

        <p>
          First, the bittersweet: Reviso will be winding down in its current
          form at the end of February, 2025.
        </p>

        <p>
          But not to bury the lede – we will be open-sourcing Reviso, and the
          current app will be migrating to a new home at pointy.ai
        </p>

        <h2 className="!font-sans">The Story So Far</h2>

        <p>
          A little over a year ago, I joined together with Colleen, Taylor, and
          James to build a new tool enabling the confident expression of ideas
          in writing. It&apos;s been a privilege to serve our users in that
          journey, and to have seen the ways in which you&apos;ve used Reviso to
          accelerate brainstorms, make painful writing less painful, tackle
          academic projects, and even explore creative fiction.
        </p>

        <p>
          As we developed Reviso, we discovered a more specific challenge that
          speaks to our backgrounds in both product development and customer
          support. To take on that challenge, Colleen and I are shifting
          attention to a new product,{" "}
          <a href="http://kestrel.app" className="text-primary hover:underline">
            Kestrel.app
          </a>
          .
        </p>

        <p>
          Meanwhile, Taylor will be starting a new company, pointy.ai, to
          further develop Reviso&apos;s product. PointyAI will continue to be
          the writing app you&apos;ve come to know and (hopefully) love. Taylor
          plans to keep developing PointyAI with new features and special
          attention to helping students learn to write and find their voice.
        </p>

        <p>
          And in the spirit of learning – we&apos;ve learned <em>a lot</em> in
          building Reviso. Bringing an AI product to life is a major technical
          challenge. In order to share our learnings, we&apos;ll be
          open-sourcing the entire Reviso codebase. We hope that the technical
          innovations that have enabled Reviso can empower other builders to
          continue to push what&apos;s possible with AI experiences.
        </p>

        <h2 className="!font-sans">What this Means for Our Users</h2>

        <ul>
          <li>
            Your account and all of your existing work in Reviso will{" "}
            <em>automatically</em> transfer to PointyAI. You do not need to take
            any action to facilitate this transfer.
          </li>
          <li>
            Starting on February 28th, logins to Reviso will automatically
            redirect to PointyAI.
          </li>
          <li>
            If you want to opt out of this automatic transfer, please email us
            at{" "}
            <a
              href="mailto:founders@revi.so"
              className="text-primary hover:underline"
            >
              founders@revi.so
            </a>{" "}
            before February 26th
          </li>
        </ul>

        <h2 className="!font-sans">A Moment of Gratitude</h2>

        <p>
          We could not be more grateful for the support of you, Reviso&apos;s
          earliest users, who were willing to take on the bugs of a new product
          and give us your honest feedback.
        </p>

        <p>
          When we set out on this journey, we had no idea where it would lead
          us. We certainly didn&apos;t expect to be inventing a new editor from
          scratch! We&apos;ve learned so much on the way, and are excited to
          take those learnings forward with us into Kestrel and PointyAI.
        </p>

        <p>
          I&apos;ve signed off my product update emails with &quot;happy
          writing&quot;. My hope is that Reviso has helped make your writing,
          your thinking, and your communicating just a little bit happier.
        </p>

        <p className="mt-8">
          See you around,
          <br />
          Justin, on behalf of team Reviso (Colleen, Taylor, and James)
        </p>
      </div>
    </article>
  );
};

export default function WindingDownPage() {
  return (
    <>
      <MarketingNav pathname="" />
      <main>
        <CenteredLayout>
          <LetterContent />
        </CenteredLayout>
      </main>
    </>
  );
}
