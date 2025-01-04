package quizzer

type Question struct {
	ID                    string   `json:"id"`
	QuizID                string   `json:"quiz_id"`
	Question              string   `json:"question"`
	Index                 int      `json:"index"`
	TimeLimitSeconds      uint64   `json:"time_limit_seconds"`
	Answers               []string `json:"answers"`
	CorrectAnswers        []string `json:"correct_answers"`
	VideoURL              *string  `json:"video_url,omitempty"`
	VideoStartTimeSeconds *uint64  `json:"video_start_time_seconds,omitempty"`
	VideoEndTimeSeconds   *uint64  `json:"video_end_time_seconds,omitempty"`
}
