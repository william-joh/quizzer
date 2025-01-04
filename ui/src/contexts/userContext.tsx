import { request } from "@/lib/axios";
import { createContext, useState, useContext } from "react";

interface User {
  id: string;
  signupDate: string;
  username: string;
}

export const CurrentUserContext = createContext<{
  currentUser: User | null;
  setCurrentUser: (user: User | null) => void;
  fetchCurrentUser: () => Promise<void>;
}>({
  currentUser: null,
  setCurrentUser: () => {},
  fetchCurrentUser: () => {
    return new Promise(() => {});
  },
});

export const CurrentUserProvider = ({ children }: { children: any }) => {
  const [currentUser, setCurrentUser] = useState<User | null>(null);

  const fetchCurrentUser = async () => {
    try {
      const resp = await request({
        url: "/current-user",
        method: "GET",
      });

      console.log("Current user resp", resp);

      setCurrentUser(resp.data);
    } catch (error) {
      console.error("Error fetching current user", error);
      setCurrentUser(null);
      return;
    }
  };

  return (
    <CurrentUserContext.Provider
      value={{
        currentUser,
        setCurrentUser: (user) => setCurrentUser(user),
        fetchCurrentUser,
      }}
    >
      {children}
    </CurrentUserContext.Provider>
  );
};

export const useCurrentUser = () => useContext(CurrentUserContext);
