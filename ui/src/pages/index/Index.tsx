import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Link, useNavigate } from "react-router";
import { useState } from "react";

export function Index() {
  const navigate = useNavigate();
  const [code, setCode] = useState("");

  const handleJoinQuiz = () => {
    if (code.trim()) {
      navigate(`/game/${code.trim()}`);
    }
  };

  return (
    <div className="container mx-auto h-[calc(100vh-3.5rem)] flex items-center overflow-hidden">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 w-full max-w-4xl mx-auto">
        {/* Join Quiz Section */}
        <Card className="w-full">
          <CardHeader>
            <CardTitle>Join a Quiz</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col space-y-4">
              <Input
                placeholder="Enter quiz code"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleJoinQuiz()}
              />
              <Button onClick={handleJoinQuiz}>Join Quiz</Button>
            </div>
          </CardContent>
        </Card>

        {/* Create Quiz Section */}
        <Card className="w-full">
          <CardHeader>
            <CardTitle>Create & Manage Quizzes</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col space-y-4">
              <p className="text-muted-foreground">
                Create, manage, view and start all your quizzes.
              </p>
              <Button asChild>
                <Link to="/quizzes">Go to Quizzes</Link>
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
