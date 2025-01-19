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
    <div className="max-w-4xl mx-auto mt-8">
      <div className="flex justify-between items-center mb-8">
        <h2 className="text-3xl font-bold">Question</h2>
        <div className="text-4xl font-mono bg-secondary px-6 py-3 rounded-lg">
          {timeLeft}s
        </div>
      </div>

      <div className="bg-card p-6 rounded-lg shadow-lg mb-8">
        <p className="text-2xl">{question.question}</p>
      </div>

      <div className="grid grid-cols-2 gap-4">
        {question.options.map((option, index) => (
          <div key={index} className="bg-secondary p-6 rounded-lg text-xl">
            {option}
          </div>
        ))}
      </div>
    </div>
  );
}
