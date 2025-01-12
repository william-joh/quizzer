import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { request } from "@/lib/axios";
import { useMutation, useQuery } from "@tanstack/react-query";
import moment from "moment";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { MoreVertical, Play } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Navigate, useNavigate } from "react-router";

interface Quiz {
  id: string;
  title: string;
  createdByName: string;
  createdAt: string;
  nrQuestions: number;
}

export function QuizList() {
  const navigate = useNavigate();

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

  const startQuizMutation = useMutation({
    mutationFn: async (quizId: string) => {
      const response = await request({
        url: `/quizzes/${quizId}/start`,
        method: "POST",
      });
      console.log("start quiz response", response);

      return response.data;
    },
    onSuccess: (data) => {
      console.log("Quiz started", data.code);
      navigate("/game/" + data.code);
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
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle>{quiz.title}</CardTitle>

              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" className="h-8 w-8 p-0">
                    <span className="sr-only">Open menu</span>
                    <MoreVertical className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem
                    className="cursor-pointer"
                    onClick={() => {
                      startQuizMutation.mutate(quiz.id);
                    }}
                  >
                    <Play className="mr-2 h-4 w-4" />
                    <span>Start Quiz</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
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
