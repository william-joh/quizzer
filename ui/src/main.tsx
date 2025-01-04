import { StrictMode } from "react";
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

const queryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <CurrentUserProvider>
        <BrowserRouter>
          <Routes>
            <Route element={<Navbar />}>
              <Route element={<CenterLayout />}>
                <Route index element={<Index />} />

                <Route path="login" element={<Login />} />
                <Route path="/quizzes/create" element={<CreateQuiz />} />
              </Route>
            </Route>
          </Routes>
        </BrowserRouter>
      </CurrentUserProvider>
    </QueryClientProvider>
  </StrictMode>
);
