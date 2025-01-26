import { QuizInfo } from "./Game";
import { Card, CardContent } from "@/components/ui/card";

export function GameInfo({ quizInfo }: { quizInfo: QuizInfo }) {
  if (!quizInfo) return null;

  return (
    <Card className="mt-4">
      <CardContent className="flex justify-between items-center py-6">
        <h1 className="text-2xl font-bold">{quizInfo.title}</h1>
        <div className="text-muted-foreground">
          Host: <span className="font-medium">{quizInfo.hostName}</span>
        </div>
      </CardContent>
    </Card>
  );
}
