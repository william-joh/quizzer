import { useEffect, useState } from "react";
import { QuizInfo } from "../Game";
import { GameInfo } from "../GameInfo";
import { ParticipantQuestionPhase } from "./ParticipantQuestionPhase";
import { Participant } from "../GamePage";
import { ParticipantResultsPhase } from "./ParticipantResultsPhase";
import { ParticipantLobby } from "./ParticipantLobby";

interface ParticipantGameProps {
  ws: WebSocket;
  initialQuizInfo: QuizInfo;
  participant: Participant;
}

export function ParticipantGame({
  ws,
  initialQuizInfo,
  participant,
}: ParticipantGameProps) {
  const [quizInfo, setQuizInfo] = useState(initialQuizInfo);
  const [phase, setPhase] = useState("lobby");
  const [options, setOptions] = useState<string[]>([]);

  useEffect(() => {
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setPhase(data.phase);

      switch (data.phase) {
        case "lobby":
          setQuizInfo({
            title: data.quizTitle,
            hostName: data.hostName,
            isHost: data.isHost,
            participants: data.participantNames,
          });
          break;
        case "question":
          setOptions(data.options);
          break;
        case "results":
          break;
        default:
          throw new Error("Unknown phase");
      }
    };

    // return () => {
    //   ws.close();
    // };
  }, [ws]);

  const answerQuestion = (answer: string) => {
    ws.send(
      JSON.stringify({
        type: "AnswerQuestion",
        data: { answer, id: participant.id },
      })
    );
  };

  return (
    <div>
      <GameInfo quizInfo={quizInfo} />

      {phase == "lobby" && <ParticipantLobby />}

      {phase == "question" && (
        <ParticipantQuestionPhase
          options={options}
          onSelectQuestion={answerQuestion}
        />
      )}

      {phase == "results" && <ParticipantResultsPhase />}
    </div>
  );
}
