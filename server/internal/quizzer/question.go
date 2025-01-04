package quizzer

type Question struct {
	ID                    string   `json:"id"`
	QuizID                string   `json:"quizId"`
	Question              string   `json:"question"`
	Index                 int      `json:"index"`
	TimeLimitSeconds      uint64   `json:"timeLimitSeconds"`
	Answers               []string `json:"answers"`
	CorrectAnswers        []string `json:"correctAnswers"`
	VideoURL              *string  `json:"videoUrl,omitempty"`
	VideoStartTimeSeconds *uint64  `json:"videoStartTimeSeconds,omitempty"`
	VideoEndTimeSeconds   *uint64  `json:"videoEndTimeSeconds,omitempty"`
}
