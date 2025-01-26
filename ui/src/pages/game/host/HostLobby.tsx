import { Button } from "@/components/ui/button";
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
    <div className="mt-4">
      <div className="flex items-center gap-2 mb-8">
        <span>Game Code:</span>
        <span className="font-mono font-medium">{code}</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div className="space-y-1">
          <h2 className="text-2xl font-semibold">Players in Lobby</h2>
          <p className="text-sm text-muted-foreground">
            {participants.length} player{participants.length !== 1 ? "s" : ""}{" "}
            joined
          </p>
        </div>
        <Button onClick={onStart} size="lg" disabled={participants.length < 1}>
          Start Game
        </Button>
      </div>

      <ul className="flex flex-wrap gap-3">
        {participants.map((participant, index) => (
          <li
            key={index}
            className="px-5 py-2.5 bg-secondary rounded-lg shadow-sm"
          >
            {participant}
          </li>
        ))}
      </ul>
    </div>
  );
}
