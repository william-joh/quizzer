import { Card, CardContent, CardHeader } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Trophy } from "lucide-react";

interface HostResultsPhaseProps {
  nrQuestionsCompleted: number;
  totalQuestions: number;
  results: { name: string; nrCorrect: number }[];
  onContinue: () => void;
}

export function HostResultsPhase({
  nrQuestionsCompleted,
  totalQuestions,
  results,
  onContinue,
}: HostResultsPhaseProps) {
  const sortedResults = [...results].sort((a, b) => b.nrCorrect - a.nrCorrect);
  const isLastQuestion = nrQuestionsCompleted === totalQuestions;

  return (
    <Card className="mt-4">
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <div className="space-y-1.5">
          <h2 className="text-2xl font-semibold tracking-tight">Results</h2>
          <p className="text-sm text-muted-foreground">Current standings</p>
        </div>
        <div className="bg-secondary px-4 py-2 rounded-md font-medium text-sm">
          Question {nrQuestionsCompleted} of {totalQuestions}
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-16">Rank</TableHead>
              <TableHead>Player</TableHead>
              <TableHead className="text-right">Score</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {sortedResults.map((result, index) => (
              <TableRow key={result.name}>
                <TableCell className="font-medium">
                  {index === 0 && (
                    <Trophy className="h-4 w-4 text-yellow-500 inline mr-1" />
                  )}
                  {index + 1}
                </TableCell>
                <TableCell>{result.name}</TableCell>
                <TableCell className="text-right font-medium">
                  {result.nrCorrect} / {nrQuestionsCompleted}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <div className="flex justify-end pt-4">
          <Button
            size="lg"
            onClick={onContinue}
            variant={isLastQuestion ? "secondary" : "default"}
          >
            {isLastQuestion ? "End Quiz" : "Next Question"}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
