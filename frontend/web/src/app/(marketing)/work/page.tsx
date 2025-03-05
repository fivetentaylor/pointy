import { HeroSection } from "@/components/marketing/HeroSection";
import { HowItWorksSection } from "@/components/marketing/HowItWorksSection";
import { UseCasesSection } from "@/components/marketing/UseCasesSection";
import { CTASection } from "@/components/marketing/CTASection";
import {
  PencilRulerIcon,
  BookOpenIcon,
  LightbulbIcon,
  MessageCircleHeartIcon,
  HandshakeIcon,
  MicVocalIcon,
} from "lucide-react";
import { MarketingNav } from "@/components/marketing/MarketingNav";
import { Button } from "@/components/ui/button";
import { APP_HOST } from "@/lib/urls";
import { CenteredLayout } from "@/components/marketing/CenteredLayout";

const workHowItWorks = [
  {
    title: "Upload your knowledge",
    description:
      "Seamlessly integrate your research, documents, and references to power your writing.",
    image: "/images/how-it-works/upload-your-knowledge.png",
    height: 38,
    centerImage: true,
    imageStyle: {
      width: 1019,
      height: 516,
    },
  },
  {
    title: "Write in conversation",
    description:
      "Ask Reviso to make the changes you want using your own terms.",
    image: "/images/how-it-works/write-in-conversation.png",
    height: 36.1875,
    imageStyle: {
      width: 1032,
      height: 738,
    },
  },
  {
    title: "Edit quickly",
    description: "Perfect your writing instantly by highlighting any text.",
    image: "/images/how-it-works/edit-quickly.png",
    height: 38,
    imageStyle: {
      width: 1048,
      height: 962,
    },
  },
  {
    title: "Maintain your voice",
    description:
      "Keep your unique voice while getting whole document suggestions.",
    image: "/images/how-it-works/maintain-your-voice.png",
    height: 38,
    imageStyle: {
      width: 908,
      height: 1072,
    },
  },
  {
    title: "Talk to your document",
    description:
      "Use your voice to speak to AI and get feedback, make edits, and more.",
    image: "/images/how-it-works/talk-to-your-doc.png",
    height: 38,
    topAlign: true,
    imageStyle: {
      height: 535,
      width: 748,
      maxWidthRem: 28,
      marginLeft: true,
    },
  },
  {
    title: "Collaborate with others",
    description:
      "Give targeted feedback exactly where it's needed. Add inline comments, suggest revisions, and more.",
    image: "/images/how-it-works/collaborate-with-others.png",
    imageStyle: {
      width: 922,
      height: 976,
    },
    height: 39.8375,
  },
];

const workUseCases = [
  {
    title: "Strategic plans",
    description:
      "Create comprehensive roadmaps and strategies that align teams.",
    icon: PencilRulerIcon,
  },
  {
    title: "Marketing content",
    description:
      "Craft compelling messages and campaigns that resonate with your audience.",
    icon: MicVocalIcon,
  },
  {
    title: "Product documentation",
    description:
      "Develop clear, user-friendly documentation that helps customers get the most from your products.",
    icon: BookOpenIcon,
  },
  {
    title: "Performance reviews",
    description:
      "Take the busy work out of your reviews and get that next promotion.",
    icon: HandshakeIcon,
  },
  {
    title: "Internal communications",
    description:
      "Build clear, effective messages that keep teams aligned, informed, and motivated.",
    icon: MessageCircleHeartIcon,
  },
  {
    title: "Thought leadership",
    description:
      "Shape industry conversations with insightful articles and perspectives that establish your expertise.",
    icon: LightbulbIcon,
  },
];

export default function WorkPage() {
  return (
    <>
      <MarketingNav pathname={"/work"} />
      <main>
        <CenteredLayout>
          <HeroSection
            video="/videos/landing/work.webm"
            title="Transform ideas into impact"
            description="Write with confidence using AI that helps you get the right message across and make your ideas stand out."
          />
          <HowItWorksSection items={workHowItWorks} />
        </CenteredLayout>
        <UseCasesSection items={workUseCases} type="professional" />
      </main>
    </>
  );
}
