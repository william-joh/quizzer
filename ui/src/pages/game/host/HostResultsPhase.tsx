import { Card } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";

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
    <div className="max-w-4xl mx-auto mt-8">
      <div className="flex items-center justify-between mb-8">
        <h2 className="text-3xl font-bold">Results</h2>
        <Card className="px-6 py-3">
          <p className="text-xl">
            Question {nrQuestionsCompleted} of {totalQuestions}
          </p>
        </Card>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Rank</TableHead>
            <TableHead>Player</TableHead>
            <TableHead className="text-right">Score</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {sortedResults.map((result, index) => (
            <TableRow key={result.name}>
              <TableCell className="font-medium">{index + 1}</TableCell>
              <TableCell>{result.name}</TableCell>
              <TableCell className="text-right">
                {result.nrCorrect} / {nrQuestionsCompleted}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      <div className="mt-8 flex justify-end">
        {!isLastQuestion && (
          <Button size="lg" onClick={onContinue}>
            Next Question
          </Button>
        )}
        {isLastQuestion && (
          <Button size="lg" onClick={onContinue} variant="secondary">
            Finish Quiz
          </Button>
        )}
      </div>
    </div>
  );
}
