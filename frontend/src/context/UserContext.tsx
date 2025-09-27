import React, { createContext, useState, useEffect } from "react";
import { API_BASE_URL } from "../config";

interface User {
  username: string;
  role: string;
}

interface UserContextType {
  user: User | null;
  setUser: React.Dispatch<React.SetStateAction<User | null>>;
}

export const UserContext = createContext<UserContextType | undefined>(
  undefined,
);

export const UserProvider: React.FC<React.PropsWithChildren> = ({
  children,
}) => {
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const fetchUser = async () => {
      const token = localStorage.getItem("authToken");
      if (token) {
        try {
          const res = await fetch(`${API_BASE_URL}/me`, {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });
          if (res.ok) {
            const data = await res.json();
            setUser({ username: data.username, role: data.role });
          }
        } catch (error) {
          console.error("Failed to fetch user:", error);
        }
      }
    };

    fetchUser();
  }, []);

  return (
    <UserContext.Provider value={{ user, setUser }}>
      {children}
    </UserContext.Provider>
  );
};
