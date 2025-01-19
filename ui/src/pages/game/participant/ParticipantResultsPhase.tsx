import { Card } from "@/components/ui/card";

export function ParticipantResultsPhase() {
  return (
    <div className="max-w-4xl mx-auto mt-8">
      <Card className="p-8 text-center">
        <h2 className="text-2xl font-semibold mb-4">Results</h2>
        <p className="text-lg text-muted-foreground">
          Please look at the host's screen to see the current standings.
        </p>
      </Card>
    </div>
  );
}
