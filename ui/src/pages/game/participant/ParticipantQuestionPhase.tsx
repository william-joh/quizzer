import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
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
    <Card className="mt-4">
      <CardHeader>
        <h2 className="text-2xl font-semibold tracking-tight">
          Choose your answer
        </h2>
        <p className="text-sm text-muted-foreground">
          Select one of the options below
        </p>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {options.map((option, index) => (
            <Button
              key={index}
              onClick={() => handleOptionSelect(option)}
              variant={selectedOption === option ? "default" : "outline"}
              className={`p-8 h-auto text-lg font-medium transition-all ${
                selectedOption === option ? "ring-2 ring-primary" : ""
              }`}
              disabled={selectedOption !== null}
            >
              {option}
            </Button>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
