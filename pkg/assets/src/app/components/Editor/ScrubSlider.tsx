import React, { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Slider } from "@/components/ui/slider";
import { ErrorBoundary } from "@/components/ui/ErrorBoundary";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { Undo, ChevronLeft, ChevronRight } from "lucide-react"; // Import Chevron icons
import { useSignals } from "@preact/signals-react/runtime";

export const ScrubSlider = ({
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

  // Update slider when currentValue changes
  useEffect(() => {
    setSliderValue([currentValue]);
  }, [currentValue]);

  const onValueChange = (n: [number]) => {
    setSliderValue(n);
    onSlider(n[0]);
  };

  // Handle incrementing the slider value
  const incrementSlider = () => {
    const newValue = Math.min(sliderValue[0] + 1, sliderMax);
    setSliderValue([newValue]);
    onSlider(newValue);
  };

  // Handle decrementing the slider value
  const decrementSlider = () => {
    const newValue = Math.max(sliderValue[0] - 1, 0);
    setSliderValue([newValue]);
    onSlider(newValue);
  };

  return (
    <div className="relative p-2 w-full">
      <div className="flex items-center space-x-2 min-w-[19.3125rem]">
        {/* Left Button */}
        <Button
          variant="outline"
          size="icon"
          className="h-8 w-8"
          onClick={decrementSlider}
          disabled={sliderValue[0] <= 0} // Disable if value is at minimum
        >
          <ChevronLeft className="h-4 w-4" />
        </Button>

        {/* Slider */}
        <div className="flex-grow min-w-[16rem]">
          <Slider
            value={sliderValue}
            onValueChange={onValueChange}
            max={sliderMax}
            step={1}
            className="w-full"
          />
        </div>

        {/* Right Button */}
        <Button
          variant="outline"
          size="icon"
          className="h-8 w-8"
          onClick={incrementSlider}
          disabled={sliderValue[0] >= sliderMax} // Disable if value is at maximum
        >
          <ChevronRight className="h-4 w-4" />
        </Button>

        {/* Revert Button */}
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

const ScrubSliderWrapper = () => {
  useSignals();

  const { editor } = useRogueEditorContext();
  const [sliderMax, setSliderMax] = useState(0);
  const [currentValue, setCurrentValue] = useState(0);

  useEffect(() => {
    const n = editor?.scrubInit(true);
    if (n) {
      setSliderMax(n);
      setCurrentValue(n); // Initialize current value
    }

    return () => {
      editor?.scrubExit();
    };
  }, [editor]);

  const onSlider = (n: number) => {
    setCurrentValue(n); // Update current value
    editor?.scrubTo(n);
  };

  const onRevert = () => {
    editor?.scrubRevert();
  };

  return (
    <ErrorBoundary fallback={<div>Error</div>}>
      <ScrubSlider
        sliderMax={sliderMax}
        onSlider={onSlider}
        onRevert={onRevert}
        currentValue={currentValue}
      />
    </ErrorBoundary>
  );
};

export default ScrubSliderWrapper;
