import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { useParams } from "react-router";

export function HostLobby({
  participants,
  onStart,
}: {
  participants: string[];
  onStart: () => void;
}) {
  const { code } = useParams();

  return (
    <Card className="mt-4">
      <CardHeader className="pb-2">
        <div className="flex items-center gap-3 text-sm">
          <span className="text-muted-foreground">Game Code:</span>
          <span className="font-mono bg-secondary px-3 py-1 rounded-md font-medium">
            {code}
          </span>
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="space-y-1.5">
            <h2 className="text-2xl font-semibold tracking-tight">
              Players in Lobby
            </h2>
            <p className="text-sm text-muted-foreground">
              {participants.length} player{participants.length !== 1 ? "s" : ""}{" "}
              joined
            </p>
          </div>
          <Button
            onClick={onStart}
            size="lg"
            className="px-8"
            disabled={participants.length < 1}
          >
            Start Game
          </Button>
        </div>

        <ul className="flex flex-wrap gap-2">
          {participants.map((participant, index) => (
            <li
              key={index}
              className="px-4 py-2 bg-secondary hover:bg-secondary/80 rounded-md
                         transition-colors duration-200 text-sm font-medium"
            >
              {participant}
            </li>
          ))}
        </ul>
      </CardContent>
    </Card>
  );
}
