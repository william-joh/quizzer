import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { request } from "@/lib/axios";
import { useQuery } from "@tanstack/react-query";
import moment from "moment";

interface Quiz {
  id: string;
  title: string;
  createdByName: string;
  createdAt: string;
  nrQuestions: number;
}

export function QuizList() {
  const { isPending, error, data } = useQuery({
    queryKey: ["quizzes"],
    queryFn: async () => {
      const response = await request({
        url: "/quizzes",
        method: "GET",
      });
      console.log("list quizzes response", response);

      return response.data;
    },
  });

  if (error) {
    return <div>Error: {error.message}</div>;
  }
  if (isPending) {
    return <div>Loading...</div>;
  }

  const quizzes = data as Quiz[];

  if (!quizzes || quizzes.length === 0) {
    return <div>No quizzes found</div>;
  }

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
                  {quiz.createdByName}
                </span>
              </div>
              <div>
                <span className="text-sm text-gray-700">Questions: </span>
                <span className="text-sm font-semibold text-gray-900">
                  {quiz.nrQuestions}
                </span>
              </div>
              <div>
                <span className="text-sm text-gray-700">Created at: </span>
                <span className="text-sm font-semibold text-gray-900">
                  {moment(quiz.createdAt).format("YYYY-MM-DD")}{" "}
                  {/* Format the date */}
                </span>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
