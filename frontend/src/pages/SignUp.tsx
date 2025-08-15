import React from "react";
import { useNavigate } from "react-router-dom";
import "../styles/SignUp.css";
import Footer from "../components/Footer";
import logo from "../../src/assets/logo.png"

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

const SignUp = () => {
  const [username, setUsername] = React.useState("");
  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const [successMessage, setSuccessMessage] = React.useState<string | null>(
    null,
  );

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setSuccessMessage(null);

    if (!username || !email || !password) {
      setError("Please fill in all fields.");
      return;
    }

    setIsLoading(true);

    try {
      const response = await fetch(`${API_BASE_URL}/register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, email, password }),
      });

      const responseText = await response.text();

      let data;
      if (responseText) {
        try {
          data = JSON.parse(responseText);
          console.log("JSON RESPONSE:", data);
        } catch (e) {
          console.error(e);
          console.error("Failed to parse JSON:", responseText);
          throw new Error(
            `An unexpected response was received from the server.`,
          );
        }
      }

      if (!response.ok) {
        const errorMessage =
          data?.message || "An unexpected error occurred. Please try again.";
        throw new Error(errorMessage);
      }

      console.log("Registration successful:", data);
      setSuccessMessage("Registration successful! Redirecting to login...");

      setTimeout(() => {
        navigate("/login");
      }, 2000);
    } catch (err: unknown) {
      console.error("Sign-up error:", err);
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unexpected sign-up error occurred.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen min-w-screen w-full h-full flex flex-col sm:flex-row items-center justify-center font-sans from-slate-100 to-sky-100 p-4">
      <div
        className="welcomeMessage absolute top-8 left-1/2 transform -translate-x-1/2 flex flex-col items-center space-y-2 z-10 animate-fade-in-down"
        style={{ minWidth: 350 }}
      >
        <img
          src={logo}
          alt="Swift Signals Logo"
          className="signupLogo h-20 w-20 object-contain drop-shadow-lg"
        />
        <span className="welcomeText text-xl md:text-4xl font-bold text-gray-800 dark:text-white flex items-center gap-2">
          Welcome to Swift Signals
        </span>
      </div>
      <div
        className="signUpContainer bg-white p-8 rounded-xl shadow-2xl w-full max-w-md"
        style={{
          boxShadow:
            "0 8px 40px rgba(0,0,0,0.18), 0 1.5px 6px rgba(0,0,0,0.12)",
        }}
      >
        <h1 className="signUpTitle text-4xl font-bold text-center text-gray-800 mb-6">
          Sign Up
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
          {successMessage && (
            <div
              className="bg-green-100 border-l-4 border-green-500 text-green-700 p-3 rounded-md text-sm"
              role="alert"
            >
              <p>{successMessage}</p>
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
              onChange={(e) => {
                e.preventDefault();
                setUsername(e.target.value);
              }}
              placeholder="Username"
              className="w-full px-4 py-3 border-2 border-[#388BFD] rounded-full bg-gray-100 text-gray-900 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
              required
              disabled={isLoading}
            />
          </div>
          <div>
            <label htmlFor="email" className="sr-only">
              Email
            </label>
            <input
              type="email"
              id="email"
              name="email"
              value={email}
              onChange={(e) => {
                e.preventDefault();
                setEmail(e.target.value);
              }}
              placeholder="Email"
              className="w-full px-4 py-3 border-2 border-[#388BFD] rounded-full bg-gray-100 text-gray-900 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
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
              onChange={(e) => {
                e.preventDefault();
                setPassword(e.target.value);
              }}
              placeholder="Password"
              className="w-full px-4 py-3 border-2 border-[#388BFD] rounded-full bg-gray-100 text-gray-900 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
              required
              disabled={isLoading}
            />
          </div>
          <div>
            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-indigo-600 dark:bg-[#388BFD] hover:bg-indigo-700 text-white font-semibold py-3 px-4 rounded-full focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-300 ease-in-out transform hover:scale-105 shadow-lg disabled:bg-indigo-400 disabled:cursor-not-allowed flex items-center justify-center"
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
              {isLoading ? "Registering..." : "Register"}
            </button>
          </div>
        </form>
        <p className="regLink mt-8 text-center text-sm text-gray-600">
          Already have an account?{" "}
          <button
            type="button"
            onClick={() => navigate("/login")}
            className="font-medium text-indigo-600 dark:text-[#388BFD] hover:text-indigo-800 hover:underline transition-colors bg-transparent border-none p-0 m-0 cursor-pointer"
            style={{ background: "none" }}
          >
            Login here
          </button>
        </p>
      </div>

      <div className="mt-8 sm:mt-0 sm:ml-12 flex justify-center">
        <TrafficLight
          redActive={username.length > 0}
          yellowActive={email.length > 0}
          greenActive={password.length > 0}
        />
      </div>
      <Footer />
    </div>
  );
};

export default SignUp;
