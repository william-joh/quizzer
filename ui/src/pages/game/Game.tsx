import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router";
import { Participant } from "./GamePage";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import { HostGame } from "./host/HostGame";
import { ParticipantGame } from "./participant/ParticipantGame";

export interface QuizInfo {
  title: string;
  hostName: string;
  isHost: boolean;
  participants: string[];
}

export function Game({ participant }: { participant: Participant }) {
  console.log("Participant", participant);
  const { code } = useParams();

  const [quizInfo, setQuizInfo] = useState<QuizInfo | null>(null);
  const [connectionError, setConnectionError] = useState(false);

  const ws = useRef<WebSocket>(
    new WebSocket(`ws://127.0.0.1:8000/game/${code}`)
  );

  // Create websocket to connect to the game
  useEffect(() => {
    console.log("Code", code);
    ws.current.onopen = () => {
      ws.current.send(
        `{ "type": "Join", "data": { "username": "${participant.username}", "id": "${participant.id}"} }`
      );
    };
    ws.current.onmessage = (event) => {
      console.log("Got message", event.data);
      const data = JSON.parse(event.data);
      console.log("Data", data);

      if (data.phase == "lobby") {
        setQuizInfo({
          title: data.quizTitle,
          hostName: data.hostName,
          isHost: data.isHost,
          participants: data.participantNames,
        });
      } else {
        throw new Error("Unknown phase");
      }
    };
    ws.current.onclose = () => {
      console.log("Disconnected from game");
      setConnectionError(true);
    };

    // return () => {
    //   ws.close();
    // };
  }, [code]);

  if (connectionError) {
    return (
      <div style={{ marginLeft: "8px", marginRight: "8px" }}>
        <Alert variant="destructive" className="mt-4">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Connection Error</AlertTitle>
          <AlertDescription>
            Lost connection to the game. Please try refreshing the page.
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  if (!quizInfo) {
    return <div>Loading...</div>;
  }

  return (
    <div style={{ marginLeft: "8px", marginRight: "8px" }}>
      {quizInfo.isHost && (
        <HostGame ws={ws.current} initialQuizInfo={quizInfo} />
      )}

      {!quizInfo.isHost && (
        <ParticipantGame
          ws={ws.current}
          initialQuizInfo={quizInfo}
          participant={participant}
        />
      )}
    </div>
  );
}
