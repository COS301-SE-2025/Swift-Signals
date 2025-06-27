import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "../styles/Login.css";
import Footer from "../components/Footer";

const API_BASE_URL = "http://localhost:9090";

interface TrafficLightProps {
  redActive: boolean;
  yellowActive: boolean;
  greenActive: boolean;
}

const TrafficLight = ({
  redActive,
  yellowActive,
  greenActive,
}: TrafficLightProps) => {
  const baseLightClasses =
    "w-20 h-20 rounded-full border-[4px] flex items-center justify-center transition-all duration-300 ease-in-out relative overflow-hidden";
  const inactiveBorderColor = "border-neutral-700";
  const activeBorderColor = "border-neutral-500";
  const inactiveRed = "bg-red-900/50 shadow-inner";
  const inactiveYellow = "bg-yellow-900/50 shadow-inner";
  const inactiveGreen = "bg-green-900/50 shadow-inner";
  const activeRed = "bg-red-600 shadow-[0_0_40px_16px_rgba(239,68,68,0.6)]";
  const activeYellow =
    "bg-yellow-500 shadow-[0_0_40px_16px_rgba(250,204,21,0.6)]";
  const activeGreen = "bg-green-500 shadow-[0_0_40px_16px_rgba(34,197,94,0.6)]";
  const innerHighlightBase =
    "w-6 h-6 rounded-full absolute top-1/3 left-1/3 transform -translate-x-1/2 -translate-y-1/2 opacity-80 blur-[2px]";

  return (
    <div className="traffic-light bg-gradient-to-b from-neutral-800 to-neutral-900 via-neutral-900 p-4 rounded-xl shadow-2xl flex flex-col space-y-4 w-30 items-center border border-neutral-700/70">
      <div
        className={`${baseLightClasses} ${
          redActive ? activeBorderColor : inactiveBorderColor
        } ${redActive ? activeRed : inactiveRed}`}
      >
        {redActive && (
          <div className={`${innerHighlightBase} bg-red-300`}></div>
        )}
      </div>
      <div
        className={`${baseLightClasses} ${
          yellowActive ? activeBorderColor : inactiveBorderColor
        } ${yellowActive ? activeYellow : inactiveYellow}`}
      >
        {yellowActive && (
          <div className={`${innerHighlightBase} bg-yellow-200`}></div>
        )}
      </div>
      <div
        className={`${baseLightClasses} ${
          greenActive ? activeBorderColor : inactiveBorderColor
        } ${greenActive ? activeGreen : inactiveGreen}`}
      >
        {greenActive && (
          <div className={`${innerHighlightBase} bg-green-300`}></div>
        )}
      </div>
    </div>
  );
};

const Login = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [resetEmail, setResetEmail] = useState("");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const navigate = useNavigate();

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [resetSuccessMessage, setResetSuccessMessage] = useState<string | null>(
    null,
  );
  const [resetError, setResetError] = useState<string | null>(null);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    if (!username || !password) {
      setError("Please fill in all fields.");
      return;
    }
    setIsLoading(true);

    try {
      const response = await fetch(`${API_BASE_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: username, password: password }),
      });
      
      const responseText = await response.text();
      if (!response.ok) {
        let serverMessage = `Status: ${response.status}`;
        try {
          const errorData = JSON.parse(responseText);
          serverMessage = errorData?.message || serverMessage;
        } catch (e) {
          console.error("Could not parse error response as JSON:", responseText);
        }
        throw new Error(`Login failed. Server says: "${serverMessage}"`);
      }
      
      const data = JSON.parse(responseText);

      if (data?.token) {
        localStorage.setItem("authToken", data.token);
        console.log("Login successful:", data.message);
        navigate("/dashboard");
      } else {
        throw new Error("Login failed: No authentication token received.");
      }
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unexpected login error occurred.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleForgotPasswordSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setResetError(null);
    setResetSuccessMessage(null);
    setIsLoading(true);

    try {
      const response = await fetch(`${API_BASE_URL}/reset-password`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: resetEmail }),
      });

      const responseText = await response.text();
      let data;
      if (responseText) {
        try {
          data = JSON.parse(responseText);
        } catch (e) {
          console.error("Failed to parse JSON from reset-password:", responseText);
          throw new Error("An unexpected response was received from the server.");
        }
      }

      if (!response.ok) {
        throw new Error(data?.message || "Failed to send reset link.");
      }

      setResetSuccessMessage(data?.message || "Password reset instructions sent to your email.");

      setTimeout(() => {
        setIsModalOpen(false);
        setResetEmail("");
        setResetSuccessMessage(null);
      }, 3000);
    } catch (err: unknown) {
      console.error("Password reset error:", err);
      if (err instanceof Error) {
        setResetError(err.message);
      } else {
        setResetError("An unexpected password reset error occurred.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="loginScreen min-h-screen min-w-screen w-full h-full flex flex-col sm:flex-row items-center justify-center font-sans from-slate-100 to-sky-100 p-4">
      <div
        className="welcomeMessage absolute top-8 left-1/2 transform -translate-x-1/2 flex flex-col items-center space-y-2 z-10 animate-fade-in-down"
        style={{ minWidth: 350 }}
      >
        <img
          src="/src/assets/logo.png"
          alt="Swift Signals Logo"
          className="loginLogo h-20 w-20 object-contain drop-shadow-lg"
        />
        <span className="welcomeText text-xl md:text-4xl font-bold text-gray-800 dark:text-white flex items-center gap-2">
          Welcome to Swift Signals
        </span>
      </div>
      <div
        className="loginContainer bg-white p-8 rounded-xl shadow-2xl w-full max-w-md"
        style={{
          boxShadow:
            "0 8px 40px rgba(0,0,0,0.18), 0 1.5px 6px rgba(0,0,0,0.12)",
        }}
      >
        <h1 className="loginTitle text-4xl font-bold text-center text-gray-800 mb-6">
          Login
        </h1>
        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div
              className="bg-red-100 border-l-4 border-red-500 text-red-700 p-3 rounded-md text-sm"
              role="alert"
            >
              <p>{error}</p>
            </div>
          )}
          <div>
            <label htmlFor="username" className="sr-only">
              Username
            </label>
            <input
              type="text"
              id="username"
              name="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Username or Email"
              className="w-full px-4 py-3 border border-blue-300 rounded-full bg-gray-100 text-gray-900 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
              required
              disabled={isLoading}
            />
          </div>
          <div>
            <label htmlFor="password" className="sr-only">
              Password
            </label>
            <input
              type="password"
              id="password"
              name="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Password"
              className="w-full px-4 py-3 border border-blue-300 rounded-full bg-gray-100 text-gray-900 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
              required
              disabled={isLoading}
            />
          </div>
          <div className="text-right">
            <a
              href="#"
              onClick={(e) => {
                e.preventDefault();
                setIsModalOpen(true);
              }}
              className="text-sm text-indigo-600 dark:text-indigo-500 hover:text-indigo-800 hover:underline transition-colors"
            >
              Forgot Password?
            </a>
          </div>
          <div>
            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-3 px-4 rounded-full focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-300 ease-in-out transform hover:scale-105 shadow-lg disabled:bg-indigo-400 disabled:cursor-not-allowed flex items-center justify-center"
            >
              {isLoading && (
                <svg
                  className="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  ></circle>
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
              )}
              {isLoading ? "Logging In..." : "Log Me In"}
            </button>
          </div>
        </form>
        <p className="regLink mt-8 text-center text-sm text-gray-600">
          Don't have an account?{" "}
          <button
            type="button"
            onClick={() => navigate("/signup")}
            className="font-medium text-indigo-600 dark:text-indigo-500 hover:text-indigo-800 hover:underline transition-colors bg-transparent border-none p-0 m-0 cursor-pointer"
            style={{ background: "none" }}
          >
            Register Here
          </button>
        </p>
      </div>
      {isModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-2xl w-full max-w-sm">
            <h2 className="text-2xl font-bold text-gray-800 dark:text-white mb-4">
              Reset Password
            </h2>
            <form onSubmit={handleForgotPasswordSubmit} className="space-y-4">
              {resetSuccessMessage && (
                <div className="bg-green-100 text-green-800 p-3 rounded-md text-sm">
                  {resetSuccessMessage}
                </div>
              )}
              {resetError && (
                <div className="bg-red-100 text-red-800 p-3 rounded-md text-sm">
                  {resetError}
                </div>
              )}
              <div>
                <label htmlFor="resetEmail" className="sr-only">
                  Email
                </label>
                <input
                  type="email"
                  id="resetEmail"
                  name="resetEmail"
                  value={resetEmail}
                  onChange={(e) => setResetEmail(e.target.value)}
                  placeholder="Enter your email"
                  className="w-full px-4 py-3 border border-blue-300 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
                  required
                  disabled={isLoading}
                />
              </div>
              <div className="flex justify-end space-x-2">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  disabled={isLoading}
                  className="px-4 py-2 bg-gray-300 dark:bg-gray-600 text-gray-800 dark:text-white rounded-full hover:bg-gray-400 dark:hover:bg-gray-500 transition-colors disabled:opacity-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isLoading}
                  className="px-4 py-2 bg-indigo-600 text-white rounded-full hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-colors disabled:bg-indigo-400 flex items-center justify-center"
                >
                  {isLoading && !resetSuccessMessage && (
                    <svg
                      className="animate-spin -ml-1 mr-2 h-5 w-5 text-white"
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <circle
                        className="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        strokeWidth="4"
                      ></circle>
                      <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                      ></path>
                    </svg>
                  )}
                  Send Reset Link
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
      <div className="traffic-light-container mt-8 sm:mt-0 sm:ml-12 flex justify-center">
        <TrafficLight
          redActive={username.length > 0}
          yellowActive={username.length > 0}
          greenActive={password.length > 0}
        />
      </div>
      <Footer />
    </div>
  );
};

export default Login;