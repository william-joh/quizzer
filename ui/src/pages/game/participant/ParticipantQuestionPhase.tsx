import { Button } from "@/components/ui/button";
import { useState } from "react";

export function ParticipantQuestionPhase({
  options,
  onSelectQuestion,
}: {
  options: string[];
  onSelectQuestion: (option: string) => void;
}) {
  const [selectedOption, setSelectedOption] = useState<string | null>(null);

  const handleOptionSelect = (option: string) => {
    setSelectedOption(option);
    onSelectQuestion(option);
  };

  return (
    <div className="max-w-4xl mx-auto mt-8">
      <div className="flex justify-between items-center mb-8">
        <h2 className="text-3xl font-bold">Select An Option</h2>
      </div>

      <div className="grid grid-cols-2 gap-4">
        {options.map((option, index) => (
          <Button
            key={index}
            onClick={() => handleOptionSelect(option)}
            variant={selectedOption === option ? "default" : "secondary"}
            className="p-6 h-auto text-xl"
            disabled={selectedOption !== null}
          >
            {option}
          </Button>
        ))}
      </div>
    </div>
  );
}
