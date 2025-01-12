import { BrowserRouter, Routes, Route } from "react-router";
import { createRoot } from "react-dom/client";
import "./index.css";
import { Index } from "./pages/index/Index.tsx";
import { Login } from "./pages/login/Login.tsx";
import { CenterLayout } from "./layouts/CenterLayout.tsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { CurrentUserProvider } from "./contexts/userContext.tsx";
import { Navbar } from "./pages/navbar/Navbar.tsx";
import { CreateQuiz } from "./pages/create-quiz/CreateQuiz.tsx";
import { AuthLayout } from "./layouts/AuthLayout.tsx";
import { Quizzes } from "./pages/quizzes/Quizzes.tsx";
import { GamePage } from "./pages/game/GamePage.tsx";

const queryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
  <QueryClientProvider client={queryClient}>
    <CurrentUserProvider>
      <BrowserRouter>
        <Routes>
          <Route element={<Navbar />}>
            <Route element={<CenterLayout />}>
              <Route index element={<Index />} />

              <Route path="game/:code" element={<GamePage />} />

              <Route path="login" element={<Login />} />

              <Route element={<AuthLayout />}>
                <Route path="quizzes" element={<Quizzes />} />
                <Route path="/quizzes/create" element={<CreateQuiz />} />
              </Route>
            </Route>
          </Route>
        </Routes>
      </BrowserRouter>
    </CurrentUserProvider>
  </QueryClientProvider>
);
