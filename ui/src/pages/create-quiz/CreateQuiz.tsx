import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useNavigate } from "react-router";
import { useMutation } from "@tanstack/react-query";
import { request } from "@/lib/axios";
import { Loader2, Plus, Save, Trash2, MoveUp, MoveDown } from "lucide-react";
import { Label } from "@/components/ui/label";
import { z } from "zod";
import { useFieldArray, useForm, Control } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Checkbox } from "@/components/ui/checkbox";

const questionSchema = z.object({
  question: z.string().min(1, "Question is required"),
  answers: z
    .array(
      z.object({
        text: z.string().min(1, "Answer text is required"),
        isCorrect: z.boolean(),
      })
    )
    .min(2, "At least 2 answers are required")
    .refine(
      (answers) => answers.some((a) => a.isCorrect),
      "At least one answer must be marked as correct"
    ),
  timeLimitSeconds: z.number().min(5, "Time limit must be at least 5 seconds"),
});

const formSchema = z.object({
  title: z.string().min(1, "Title is required"),
  questions: z
    .array(questionSchema)
    .min(1, "At least one question is required"),
});

type FormData = z.infer<typeof formSchema>;

export function CreateQuiz() {
  const navigate = useNavigate();

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      title: "",
      questions: [
        {
          question: "",
          answers: [
            { text: "", isCorrect: true },
            { text: "", isCorrect: false },
          ],
          timeLimitSeconds: 30,
        },
      ],
    },
  });

  const { fields, append, remove, move } = useFieldArray({
    name: "questions",
    control: form.control,
  });

  const createQuizMutation = useMutation({
    mutationFn: async (data: FormData) => {
      const response = await request({
        url: "/quizzes",
        method: "POST",
        data,
      });
      return response.data;
    },
    onSuccess: () => {
      navigate("/");
    },
  });

  return (
    <div className="container mx-auto px-4 py-8">
      <Card className="max-w-2xl mx-auto">
        <CardHeader>
          <CardTitle>Create New Quiz</CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit((data) =>
                createQuizMutation.mutate(data)
              )}
              className="space-y-6"
            >
              <FormField
                control={form.control}
                name="title"
                render={({ field }) => (
                  <FormItem>
                    <Label>Quiz Title</Label>
                    <FormControl>
                      <Input placeholder="Enter quiz title" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="space-y-6">
                {fields.map((field, qIndex) => (
                  <QuestionForm
                    key={field.id}
                    control={form.control}
                    qIndex={qIndex}
                    move={move}
                    remove={remove}
                  />
                ))}

                <Button
                  type="button"
                  variant="outline"
                  onClick={() =>
                    append({
                      question: "",
                      answers: [
                        { text: "", isCorrect: false },
                        { text: "", isCorrect: false },
                      ],
                      timeLimitSeconds: 30,
                    })
                  }
                >
                  <Plus className="h-4 w-4 mr-2" />
                  Add Question
                </Button>
              </div>

              <Button
                type="submit"
                className="w-full"
                disabled={createQuizMutation.isPending}
              >
                {createQuizMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                <Save className="mr-2 h-4 w-4" />
                Save Quiz
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}

interface QuestionFormProps {
  control: Control<FormData>;
  qIndex: number;
  move: (from: number, to: number) => void;
  remove: (index: number) => void;
}

function QuestionForm({ control, qIndex, move, remove }: QuestionFormProps) {
  const {
    fields: answerFields,
    append: appendAnswer,
    remove: removeAnswer,
  } = useFieldArray({
    name: `questions.${qIndex}.answers`,
    control,
  });

  return (
    <div className="space-y-4 p-4 border rounded-lg relative">
      <div className="absolute right-2 top-2 flex gap-1">
        {qIndex > 0 && (
          <Button
            type="button"
            variant="ghost"
            size="icon"
            onClick={() => move(qIndex, qIndex - 1)}
          >
            <MoveUp className="h-4 w-4" />
          </Button>
        )}
        {qIndex < answerFields.length - 1 && (
          <Button
            type="button"
            variant="ghost"
            size="icon"
            onClick={() => move(qIndex, qIndex + 1)}
          >
            <MoveDown className="h-4 w-4" />
          </Button>
        )}
        <Button
          type="button"
          variant="ghost"
          size="icon"
          onClick={() => remove(qIndex)}
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>

      <FormField
        control={control}
        name={`questions.${qIndex}.question`}
        render={({ field }) => (
          <FormItem>
            <Label>Question {qIndex + 1}</Label>
            <FormControl>
              <Input placeholder="Enter question" {...field} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={control}
        name={`questions.${qIndex}.timeLimitSeconds`}
        render={({ field }) => (
          <FormItem>
            <div className="flex items-center gap-2">
              <Label>Time Limit (seconds):</Label>
              <FormControl>
                <Input
                  type="number"
                  min="5"
                  className="w-24"
                  {...field}
                  onChange={(e) => field.onChange(parseInt(e.target.value))}
                />
              </FormControl>
            </div>
            <FormMessage />
          </FormItem>
        )}
      />

      <div className="space-y-2">
        <Label>Answers (mark correct answers)</Label>
        {answerFields.map((answer, aIndex) => (
          <AnswerForm
            key={answer.id}
            control={control}
            qIndex={qIndex}
            aIndex={aIndex}
            removeAnswer={removeAnswer}
          />
        ))}
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => appendAnswer({ text: "", isCorrect: false })}
        >
          <Plus className="h-4 w-4 mr-2" />
          Add Answer
        </Button>
      </div>
    </div>
  );
}

interface AnswerFormProps {
  control: Control<FormData>;
  qIndex: number;
  aIndex: number;
  removeAnswer: (index: number) => void;
}

function AnswerForm({
  control,
  qIndex,
  aIndex,
  removeAnswer,
}: AnswerFormProps) {
  return (
    <div className="flex items-center gap-2">
      <FormField
        control={control}
        name={`questions.${qIndex}.answers.${aIndex}.isCorrect`}
        render={({ field: correctField }) => (
          <FormItem className="flex items-center space-x-2">
            <FormControl>
              <Checkbox
                checked={correctField.value}
                onCheckedChange={(checked: any) => {
                  correctField.onChange(checked);
                }}
              />
            </FormControl>
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name={`questions.${qIndex}.answers.${aIndex}.text`}
        render={({ field: answerField }) => (
          <FormItem className="flex-1">
            <FormControl>
              <Input placeholder={`Answer ${aIndex + 1}`} {...answerField} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      <Button
        type="button"
        variant="ghost"
        size="icon"
        onClick={() => removeAnswer(aIndex)}
      >
        <Trash2 className="h-4 w-4" />
      </Button>
    </div>
  );
}
