import { useEffect, useState } from "react";
import { QuizInfo } from "../Game";
import { GameInfo } from "../GameInfo";
import { HostLobby } from "./HostLobby";
import { HostQuestionPhase } from "./HostQuestionPhase";
import { HostResultsPhase } from "./HostResultsPhase";

interface HostGameProps {
  ws: WebSocket;
  initialQuizInfo: QuizInfo;
}

export interface HostQuestion {
  question: string;
  options: string[];
  timeLimit: number;
}

export interface HostResults {
  nrQuestionsCompleted: number;
  totalQuestions: number;
  results: { name: string; nrCorrect: number }[];
}

export function HostGame({ ws, initialQuizInfo }: HostGameProps) {
  const [quizInfo, setQuizInfo] = useState(initialQuizInfo);
  const [phase, setPhase] = useState("lobby");
  const [question, setQuestion] = useState<HostQuestion | null>(null);
  const [results, setResults] = useState<HostResults | null>(null);

  useEffect(() => {
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log("raw data", data);
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
          setQuestion({
            question: data.question,
            options: data.options,
            timeLimit: data?.timeLimit || 5,
          });
          break;
        case "results":
          setResults({
            nrQuestionsCompleted: data.nrQuestionsCompleted,
            totalQuestions: data.totalQuestions,
            results: data.results,
          });
          console.log(data.results);
          break;

        default:
          throw new Error("Unknown phase" + data);
      }
    };

    // return () => {
    //   ws.close();
    // };
  }, [ws]);

  const finishQuestion = () => {
    ws.send(`{ "type": "FinishQuestion" }`);
  };

  const nextQuestion = () => {
    ws.send(`{ "type": "NextQuestion" }`);
  };

  if (phase == "question" && !question)
    throw new Error("No question when in question phase");
  if (phase == "results" && !results)
    throw new Error("No results when in results phase");

  return (
    <div>
      <GameInfo quizInfo={quizInfo} />

      {phase == "lobby" && (
        <HostLobby
          participants={quizInfo.participants}
          onStart={() => ws.send(`{ "type": "Start" }`)}
        />
      )}

      {phase == "question" && question && (
        <HostQuestionPhase question={question} onTimeUp={finishQuestion} />
      )}

      {phase == "results" && results && (
        <HostResultsPhase {...results} onContinue={nextQuestion} />
      )}
    </div>
  );
}
