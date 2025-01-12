import { QuizInfo } from "./Game";

export function GameInfo({ quizInfo }: { quizInfo: QuizInfo }) {
  if (!quizInfo) return null;

  return (
    <div className="flex justify-between items-center py-4">
      <h1 className="text-2xl font-bold">{quizInfo.title}</h1>
      <div className="text-muted-foreground">
        Host: <span className="font-medium">{quizInfo.hostName}</span>
      </div>
    </div>
  );
}
