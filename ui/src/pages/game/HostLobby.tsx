import { Button } from "@/components/ui/button";

export function HostLobby({
  participants,
  onStart,
}: {
  participants: string[];
  onStart: () => void;
}) {
  return (
    <div className="mt-4">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-semibold">Players in Lobby</h2>
        <Button onClick={onStart} size="lg">
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
