import React, { createContext, useState, useEffect } from "react";
import { API_BASE_URL } from "../config";

interface User {
  username: string;
  role: string;
}

interface UserContextType {
  user: User | null;
  setUser: React.Dispatch<React.SetStateAction<User | null>>;
  logout: () => void;
  refetchUser: () => void;
  isLoading: boolean; // Add isLoading state
}

export const UserContext = createContext<UserContextType | undefined>(
  undefined,
);

export const UserProvider: React.FC<React.PropsWithChildren> = ({
  children,
}) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true); // Initialize as true

  const fetchUser = async () => {
    setIsLoading(true); // Set loading to true when fetching starts
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
          const userRole = data.is_admin ? "admin" : "regular";
          setUser({ username: data.username, role: userRole });
        } else {
          setUser(null); // Clear user if token is invalid or expired
          localStorage.removeItem("authToken");
        }
      } catch (error) {
        console.error("Failed to fetch user:", error);
        setUser(null); // Clear user on fetch error
        localStorage.removeItem("authToken");
      } finally {
        setIsLoading(false); // Set loading to false when fetching ends
      }
    } else {
      setUser(null); // Ensure user is null if no token
      setIsLoading(false); // Set loading to false if no token
    }
  };

  useEffect(() => {
    fetchUser();
  }, []);

  const refetchUser = () => {
    fetchUser();
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem("authToken");
  };

  return (
    <UserContext.Provider
      value={{ user, setUser, logout, refetchUser, isLoading }}
    >
      {children}
    </UserContext.Provider>
  );
};
