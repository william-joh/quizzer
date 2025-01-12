import { useCurrentUser, User } from "@/contexts/userContext";
import { useState } from "react";
import { Game } from "./Game";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export interface Participant {
  username: string;
  id: string;
}

export function GamePage() {
  const { currentUser } = useCurrentUser();
  const [currentParticipant, setCurrentParticipant] =
    useState<Participant | null>(getCurrentParticipant(currentUser));
  console.log("Current Participant", currentParticipant);

  if (!currentParticipant) {
    return (
      <UserNameForm
        onSubmit={(username) =>
          setCurrentParticipant({
            username,
            id: crypto.randomUUID(),
          })
        }
      />
    );
  }

  return <Game participant={currentParticipant} />;
}

function UserNameForm({ onSubmit }: { onSubmit: (username: string) => void }) {
  const [username, setUsername] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (username.trim()) {
      onSubmit(username.trim());
    }
  };

  return (
    <Card className="max-w-sm mx-auto mt-8">
      <CardHeader>
        <CardTitle>Enter your name</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="Your name"
            required
          />
          <Button type="submit">Join Game</Button>
        </form>
      </CardContent>
    </Card>
  );
}

function getCurrentParticipant(currentUser: User | null): Participant | null {
  if (!currentUser) {
    return null;
  }

  return { username: currentUser.username, id: currentUser.id };
}
