import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router";
import { GameInfo } from "./GameInfo";
import { HostLobby } from "./HostLobby";
import { Participant } from "./GamePage";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";

export interface QuizInfo {
  title: string;
  hostName: string;
  isHost: boolean;
}

export function Game({ participant }: { participant: Participant }) {
  console.log("Participant", participant);
  const { code } = useParams();

  const [quizInfo, setQuizInfo] = useState<QuizInfo | null>(null);
  const [phase, setPhase] = useState("");
  const [participants, setParticipants] = useState<string[]>([]);
  const [connectionError, setConnectionError] = useState(false);

  const ws = useRef<WebSocket>(
    new WebSocket(`ws://127.0.0.1:8000/game/${code}`)
  );

  // Create websocket to connect to the game
  useEffect(() => {
    console.log("Code", code);
    ws.current.onopen = () => {
      console.log("Connected to game");

      ws.current.send(
        `{ "type": "Join", "data": { "username": "${participant.username}", "id": "${participant.id}"} }`
      );
    };
    ws.current.onmessage = (event) => {
      console.log("Got message", event.data);
      const data = JSON.parse(event.data);
      console.log("Data", data);
      setPhase(data.phase);

      if (data.phase == "lobby") {
        setParticipants(data.participantNames);
        setQuizInfo({
          title: data.quizTitle,
          hostName: data.hostName,
          isHost: data.isHost,
        });
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
      <Alert variant="destructive" className="mt-4">
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>Connection Error</AlertTitle>
        <AlertDescription>
          Lost connection to the game. Please try refreshing the page.
        </AlertDescription>
      </Alert>
    );
  }

  if (!quizInfo || phase == "") {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <GameInfo quizInfo={quizInfo} />

      <h1>Game</h1>

      <div>Phase: {phase}</div>

      {quizInfo.isHost ? (
        <>
          {phase == "lobby" && (
            <HostLobby
              participants={participants}
              onStart={() => ws.current.send(`{ "type": "Start" }`)}
            />
          )}
        </>
      ) : (
        <>{phase == "lobby" && <div>Waiting for host to start game...</div>}</>
      )}
    </div>
  );
}
