import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { useParams } from "react-router";

export function ParticipantLobby() {
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
      <CardContent className="py-12">
        <div className="text-center space-y-2">
          <h2 className="text-xl font-semibold tracking-tight">
            Waiting for host
          </h2>
          <p className="text-sm text-muted-foreground">
            The game will begin shortly...
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
