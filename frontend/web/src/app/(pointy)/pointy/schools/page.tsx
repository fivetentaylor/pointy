import {
  GraduationCapIcon,
  PenLineIcon,
  MailIcon,
  FileTextIcon,
  CoinsIcon,
  LightbulbIcon,
} from "lucide-react";
import { MarketingNav } from "@/components/marketing-pointy/MarketingNav";
import { CTASection } from "@/components/marketing-pointy/CTASection";
import { UseCasesSection } from "@/components/marketing-pointy/UseCasesSection";
import { HowItWorksSection } from "@/components/marketing-pointy/HowItWorksSection";
import { HeroSection } from "@/components/marketing-pointy/HeroSection";
import { APP_HOST } from "@/lib/urls";
import { Button } from "@/components/ui/button";
import { CenteredLayout } from "@/components/marketing-pointy/CenteredLayout";

const schoolHowItWorks = [
  {
    title: "Trust student work again",
    description:
      "With full visibility into your students' writing process, research, and AI usage, you can focus on developing their ideas and skills.",
    image: "/images/how-it-works/show-your-work.png",
    height: 38,
    imageStyle: {
      width: 1029,
      height: 836,
    },
  },
  {
    title: "Review sources",
    description:
      "See how students incorporate course materials and external sources. Review citation accuracy and ensure proper attribution of both AI and reference materials.",
    height: 38,
    image: "/images/how-it-works/review-sources.png",
    imageStyle: {
      width: 1010,
      height: 756,
    },
  },
  {
    title: "Collaborate effectively",
    description:
      "Give targeted feedback exactly where it's needed. Add inline comments, suggest revisions, and engage with students' thought process.",
    height: 38,
    image: "/images/how-it-works/collaborate-effectively-schools.png",
    imageStyle: {
      width: 976,
      height: 802,
    },
  },
  {
    title: "See student process",
    description:
      "Watch how students develop their ideas and arguments in real-time. See their research process, writing iterations, and AI interactions all in one place.",
    height: 38,
    image: "/images/how-it-works/see-student-process.png",
    imageStyle: {
      width: 978,
      height: 780,
    },
  },
];

const schoolUseCases = [
  {
    title: "Letters of reference",
    description:
      "Craft compelling recommendations that highlight student achievements and potential.",
    icon: PenLineIcon,
  },
  {
    title: "Student communications",
    description:
      "Write clear, professional emails that maintain effective dialogue with students and faculty.",
    icon: MailIcon,
  },
  {
    title: "Internal memos",
    description:
      "Develop concise departmental communications that keep faculty aligned and informed.",
    icon: FileTextIcon,
  },
  {
    title: "Grant applications",
    description:
      "Create persuasive proposals that demonstrate research value and secure funding.",
    icon: CoinsIcon,
  },
  {
    title: "Thought leadership",
    description:
      "Shape academic discourse with insights that establish expertise in your field.",
    icon: LightbulbIcon,
  },
  {
    title: "Academic research",
    description:
      "Transform complex findings into clear, impactful academic publications.",
    icon: GraduationCapIcon,
  },
];

export const metadata = {
  pathname: "/pointy/schools",
};

export default function SchoolsPage() {
  return (
    <>
      <MarketingNav pathname={"/pointy/schools"} />
      <main>
        <CenteredLayout>
          <HeroSection
            video="/videos/landing/schools.webm"
            title="See how students think & work"
            description="Transform AI from a black box to a transparent learning tool. Track and guide your students' process every step of the way."
          />
          <HowItWorksSection items={schoolHowItWorks} />
          <div className="text-center my-[1.72rem]">
            <Button
              asChild
              size="lg"
              className="bg-[#7C3AED] hover:bg-[#6D28D9] text-white px-8 my-8"
            >
              <a href={`${APP_HOST}/login`}>Join our private beta</a>
            </Button>
          </div>
        </CenteredLayout>
        <UseCasesSection items={schoolUseCases} type="academic" />
        <CTASection />
      </main>
    </>
  );
}
