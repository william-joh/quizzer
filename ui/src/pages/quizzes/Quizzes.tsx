import { Button } from "@/components/ui/button";
import { Link } from "react-router";
import { QuizList } from "./QuizList";

export function Quizzes() {
  return (
    <div className="mt-8">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Quizzes</h1>
        <Button asChild>
          <Link to="/quizzes/create">Create Quiz</Link>
        </Button>
      </div>
      <QuizList />
    </div>
  );
}
