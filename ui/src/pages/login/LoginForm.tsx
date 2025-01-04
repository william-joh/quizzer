import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useMutation } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { AxiosResponse } from "axios";
import { request } from "@/lib/axios";
import { useCurrentUser } from "@/contexts/userContext";
import { useNavigate } from "react-router";

const formSchema = z.object({
  username: z.string().min(2).max(50),
  password: z.string().min(8).max(50),
});

export function LoginForm() {
  const navigate = useNavigate();
  const { fetchCurrentUser } = useCurrentUser();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      password: "",
    },
  });

  const loginMutation = useMutation({
    mutationFn: doLogin,
    onSuccess: async () => {
      console.log("login success, fetching current user");
      try {
        await fetchCurrentUser();
      } catch (error) {
        console.error("Error fetching current user", error);
        return;
      }

      console.log("login success, done fetching current user");
      navigate("/");
    },
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle>Login</CardTitle>
      </CardHeader>

      <Form {...form}>
        <form
          onSubmit={form.handleSubmit((data) => loginMutation.mutate(data))}
          className="space-y-8"
        >
          <CardContent className="space-y-2">
            <FormField
              control={form.control}
              name="username"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Username</FormLabel>
                  <FormControl>
                    <Input placeholder="" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input placeholder="" type="password" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </CardContent>

          <CardFooter className="block">
            {loginMutation.isError && (
              <div
                className="block p-4 mb-4 text-sm text-red-800 rounded-lg bg-red-50 dark:bg-gray-800 dark:text-red-400"
                role="alert"
              >
                <span>{loginMutation.error.message}</span>
              </div>
            )}

            <Button type="submit" disabled={loginMutation.isPending}>
              {loginMutation.isPending && <Loader2 className="animate-spin" />}
              Login
            </Button>
          </CardFooter>
        </form>
      </Form>
    </Card>
  );
}

async function doLogin(user: {
  username: string;
  password: string;
}): Promise<AxiosResponse<any, any>> {
  console.log("do login", user);

  return request({
    url: "/auth",
    method: "POST",
    withCredentials: true,
    withXSRFToken: true,
    data: {
      username: user.username,
      password: user.password,
    },
  });
}
