import React, { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Slider } from "@/components/ui/slider";
import { ErrorBoundary } from "@/components/ui/ErrorBoundary";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { Undo, ChevronLeft, ChevronRight } from "lucide-react";

export const Scrub = ({
  sliderMax,
  onSlider,
  onRevert,
  currentValue,
}: {
  sliderMax: number;
  onSlider: (n: number) => void;
  onRevert: () => void;
  currentValue: number;
}) => {
  const [sliderValue, setSliderValue] = useState([currentValue]);

  useEffect(() => {
    setSliderValue([currentValue]);
  }, [currentValue]);

  const onValueChange = (n: [number]) => {
    setSliderValue(n);
    onSlider(n[0]);
  };

  const stepLeft = () => {
    const newValue = Math.max(0, sliderValue[0] - 1);
    onValueChange([newValue]);
  };

  const stepRight = () => {
    const newValue = Math.min(sliderMax, sliderValue[0] + 1);
    onValueChange([newValue]);
  };

  return (
    <div className="relative p-2 w-full">
      <div className="flex items-center space-x-2 min-w-[19.3125rem]">
        <Button
          variant="outline"
          className="h-8 px-2 py-2 text-xs leading-none"
          onClick={stepLeft}
          disabled={sliderValue[0] === 0}
        >
          <ChevronLeft className="h-4 w-4" />
        </Button>
        <div className="flex-grow min-w-[14rem]">
          <Slider
            value={sliderValue}
            onValueChange={onValueChange}
            max={sliderMax}
            step={1}
            className="w-full"
          />
        </div>
        <Button
          variant="outline"
          className="h-8 px-2 py-2 text-xs leading-none"
          onClick={stepRight}
          disabled={sliderValue[0] === sliderMax}
        >
          <ChevronRight className="h-4 w-4" />
        </Button>
        <Button
          className="bg-primary text-white hover:bg-primary/90 h-8 px-2 py-2 text-xs leading-none shrink-0"
          onClick={onRevert}
        >
          <Undo className="h-4 w-4 mr-2" />
          Revert
        </Button>
      </div>
    </div>
  );
};

const ScrubWrapper = ({ container }: { container: HTMLElement | null }) => {
  const { editor } = useRogueEditorContext();
  const [sliderMax, setSliderMax] = useState(0);
  const [currentValue, setCurrentValue] = useState(0);

  useEffect(() => {
    const n = editor?.scrubInit(false);
    if (n) {
      setSliderMax(n);
      setCurrentValue(n);
    }
    return () => {
      editor?.scrubExit();
    };
  }, [editor]);

  const onSlider = (n: number) => {
    setCurrentValue(n);
    editor?.scrubTo(n);
  };

  const onRevert = () => {
    editor?.scrubRevert();
  };

  return (
    <ErrorBoundary fallback={<div>Error</div>}>
      <Scrub
        sliderMax={sliderMax}
        onSlider={onSlider}
        onRevert={onRevert}
        currentValue={currentValue}
      />
    </ErrorBoundary>
  );
};

export default ScrubWrapper;
