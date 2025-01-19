export function ParticipantLobby({ participants }: { participants: string[] }) {
  return (
    <div className="mt-4">
      <h2 className="text-xl font-semibold mb-2">Players in Lobby</h2>
      <ul className="flex flex-wrap gap-2">
        {participants.map((participant, index) => (
          <li key={index} className="px-4 py-2 bg-secondary rounded-lg">
            {participant}
          </li>
        ))}
      </ul>
    </div>
  );
}
