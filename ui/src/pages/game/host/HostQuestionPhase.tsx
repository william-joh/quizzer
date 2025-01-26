import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { useEffect, useState } from "react";
import { HostQuestion } from "./HostGame";

export function HostQuestionPhase({
  question,
  onTimeUp,
}: {
  question: HostQuestion;
  onTimeUp: () => void;
}) {
  const [timeLeft, setTimeLeft] = useState(question.timeLimit);

  useEffect(() => {
    if (timeLeft <= 0) {
      onTimeUp();
      return;
    }

    const timer = setInterval(() => {
      setTimeLeft((prev) => prev - 1);
    }, 1000);

    return () => clearInterval(timer);
  }, [timeLeft, onTimeUp]);

  return (
    <Card className="mt-4">
      <CardHeader className="flex flex-row items-center justify-between">
        <div className="space-y-1.5">
          <h2 className="text-2xl font-semibold tracking-tight">Question</h2>
          <p className="text-sm text-muted-foreground">
            Time remaining to answer
          </p>
        </div>
        <div className="text-3xl font-mono font-medium bg-secondary px-4 py-2 rounded-md">
          {timeLeft}s
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="rounded-lg border bg-card p-6">
          <p className="text-xl">{question.question}</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {question.options.map((option, index) => (
            <div
              key={index}
              className="p-6 rounded-lg bg-secondary/50 text-lg font-medium"
            >
              {option}
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
