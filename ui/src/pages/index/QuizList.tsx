import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useEffect, useState } from "react";

interface Quiz {
  id: string;
  title: string;
  createdBy: string;
  createdAt: string;
  questionCount: number;
}

const mockedQuizzes: Quiz[] = [
  {
    id: "1",
    title: "General Knowledge",
    createdBy: "User1",
    createdAt: "2023-01-01",
    questionCount: 10,
  },
  {
    id: "2",
    title: "Science Quiz",
    createdBy: "User2",
    createdAt: "2023-02-01",
    questionCount: 15,
  },
  {
    id: "3",
    title: "History Quiz",
    createdBy: "User3",
    createdAt: "2023-03-01",
    questionCount: 20,
  },
];

export function QuizList() {
  const [quizzes, setQuizzes] = useState<Quiz[]>([]);

  useEffect(() => {
    // Simulate fetching data
    setQuizzes(mockedQuizzes);
  }, []);

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {quizzes.map((quiz) => (
          <Card key={quiz.id} className="min-w-[250px]">
            <CardHeader>
              <CardTitle>{quiz.title}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <div>
                <span className="text-sm text-gray-700">Created by: </span>
                <span className="text-sm font-semibold text-gray-900">
                  {quiz.createdBy}
                </span>
              </div>
              <div>
                <span className="text-sm text-gray-700">Created at: </span>
                <span className="text-sm font-semibold text-gray-900">
                  {quiz.createdAt}
                </span>
              </div>
              <div>
                <span className="text-sm text-gray-700">Questions: </span>
                <span className="text-sm font-semibold text-gray-900">
                  {quiz.questionCount}
                </span>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
