import { HeroSection } from "@/components/marketing-pointy/HeroSection";
import { HowItWorksSection } from "@/components/marketing-pointy/HowItWorksSection";
import { UseCasesSection } from "@/components/marketing-pointy/UseCasesSection";
import { CTASection } from "@/components/marketing-pointy/CTASection";
import {
  PencilIcon,
  BookOpenIcon,
  BrainIcon,
  MessageCircleIcon,
  GraduationCapIcon,
  StarIcon,
} from "lucide-react";
import { MarketingNav } from "@/components/marketing-pointy/MarketingNav";
import { Button } from "@/components/ui/button";
import { APP_HOST } from "@/lib/urls";
import { CenteredLayout } from "@/components/marketing-pointy/CenteredLayout";

const studentHowItWorks = [
  {
    title: "Learn from your materials",
    description:
      "Transform lecture notes, readings, and research into polished assignments.",
    image: "/images/how-it-works/upload-your-knowledge.png",
    height: 38,
    centerImage: true,
    imageStyle: {
      width: 1019,
      height: 516,
    },
  },
  {
    title: "Show your work",
    description:
      "Share your Pointy document with your professors and give them full and accurate visibility into your work.",
    image: "/images/how-it-works/show-your-work.png",
    height: 36.1875,
    imageStyle: {
      width: 1029,
      height: 836,
    },
  },
  {
    title: "Write in conversation",
    description:
      "Start with rough ideas and let Pointy help shape them into clear, academic prose through natural dialog.",
    image: "/images/how-it-works/write-in-conversation-students.png",
    height: 38,
    imageStyle: {
      width: 1032,
      height: 660,
    },
  },
  {
    title: "Maintain your voice",
    description:
      "Keep your unique perspective while getting expert writing guidance.",
    image: "/images/how-it-works/maintain-your-voice-students.png",
    height: 38,
    imageStyle: {
      width: 1072,
      height: 908,
    },
  },
  {
    title: "Edit quickly",
    description: "Perfect your writing instantly by highlighting any text. ",
    image: "/images/how-it-works/edit-quickly-students.png",
    height: 38,
    imageStyle: {
      width: 1048,
      height: 948,
    },
  },
  {
    title: "Talk to your document",
    description:
      "Use your voice to speak to AI and get feedback, make edits, and more.",
    image: "/images/how-it-works/talk-to-your-doc.png",
    height: 39.8375,
    topAlign: true,
    imageStyle: {
      height: 535,
      width: 748,
      maxWidthRem: 28,
      marginLeft: true,
    },
  },
];

const studentUseCases = [
  {
    title: "Essays",
    description:
      "Develop thoughtful, well-structured essays that showcase your critical thinking and analysis.",
    icon: PencilIcon,
  },
  {
    title: "Literature reviews",
    description:
      "Synthesize complex sources into comprehensive, organized reviews that advance scholarly discourse.",
    icon: BookOpenIcon,
  },
  {
    title: "Research papers",
    description:
      "Transform your research into clear, compelling academic arguments backed by thorough evidence.",
    icon: BrainIcon,
  },
  {
    title: "Lab reports",
    description:
      "Create precise, methodical reports that effectively communicate experimental findings.",
    icon: MessageCircleIcon,
  },
  {
    title: "Thesis dissertations",
    description:
      "Craft substantial academic works that contribute meaningful insights to your field of study.",
    icon: GraduationCapIcon,
  },
  {
    title: "Personal statements",
    description:
      "Write authentic, compelling narratives that effectively communicate your goals and experiences.",
    icon: StarIcon,
  },
];

export default function StudentsPage() {
  return (
    <>
      <MarketingNav pathname={"/pointy/students"} />
      <main>
        <CenteredLayout>
          <HeroSection
            video="/videos/landing/students.webm"
            title="Show your process, own your work"
            description={`Write papers with nothing to hide from professors.
Pointy turns AI from ghostwriter to academic thought partner.`}
          />
          <HowItWorksSection items={studentHowItWorks} />
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
        <UseCasesSection items={studentUseCases} type="academic" />
        <CTASection />
      </main>
    </>
  );
}
