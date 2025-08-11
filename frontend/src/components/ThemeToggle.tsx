import React, { useState, useEffect } from "react";
import "./ThemeToggle.css";

const ThemeToggle: React.FC = () => {
  const [isDark, setIsDark] = useState(false);

  useEffect(() => {
    // Check for saved theme preference or default to light mode
    const savedTheme = localStorage.getItem("theme");
    const prefersDark = window.matchMedia(
      "(prefers-color-scheme: dark)",
    ).matches;

    if (savedTheme === "dark" || (!savedTheme && prefersDark)) {
      setIsDark(true);
      document.documentElement.classList.add("dark");
    }
  }, []);

  const toggleTheme = () => {
    const newTheme = !isDark;
    setIsDark(newTheme);

    if (newTheme) {
      document.documentElement.classList.add("dark");
      localStorage.setItem("theme", "dark");
    } else {
      document.documentElement.classList.remove("dark");
      localStorage.setItem("theme", "light");
    }
  };

  return (
    <button
      onClick={toggleTheme}
      className="theme-toggle"
      aria-label={`Switch to ${isDark ? "light" : "dark"} mode`}
      title={`Switch to ${isDark ? "light" : "dark"} mode`}
    >
      <div className="toggle-track">
        <div className="toggle-thumb">
          <div className="toggle-icon">
            {isDark ? (
              <svg
                width="12"
                height="12"
                viewBox="0 0 24 24"
                fill="currentColor"
              >
                <path d="M21.64,13a1,1,0,0,0-1.05-.14,8.05,8.05,0,0,1-3.37.73A8.15,8.15,0,0,1,9.08,5.49a8.59,8.59,0,0,1,.25-2A1,1,0,0,0,8,2.36,10.14,10.14,0,1,0,22,14.05,1,1,0,0,0,21.64,13Zm-9.5,6.69A8.14,8.14,0,0,1,7.08,5.22v.27A10.15,10.15,0,0,0,17.22,15.63a9.79,9.79,0,0,0,2.1-.22A8.11,8.11,0,0,1,12.14,19.73Z" />
              </svg>
            ) : (
              <svg
                width="12"
                height="12"
                viewBox="0 0 24 24"
                fill="currentColor"
              >
                <path d="M12,2.25a.75.75,0,0,0-.75.75v1.5a.75.75,0,0,0,1.5,0V3A.75.75,0,0,0,12,2.25ZM21,12a.75.75,0,0,0-.75-.75H18.75a.75.75,0,0,0,0,1.5h1.5A.75.75,0,0,0,21,12ZM12,18a6,6,0,1,1,6-6,6,6,0,0,1-6,6Zm0-10.5A4.5,4.5,0,1,0,16.5,12,4.5,4.5,0,0,0,12,7.5ZM5.99,4.22A.75.75,0,1,0,4.93,5.28L6,6.34A.75.75,0,1,0,7.06,5.28ZM3,12a.75.75,0,0,0,.75.75H5.25a.75.75,0,0,0,0-1.5H3.75A.75.75,0,0,0,3,12ZM4.93,18.72A.75.75,0,1,0,5.99,19.78L7.06,18.72A.75.75,0,1,0,6,17.66ZM12,21.75a.75.75,0,0,0,.75-.75V19.5a.75.75,0,0,0-1.5,0V21A.75.75,0,0,0,12,21.75ZM18.01,19.78a.75.75,0,1,0,1.06-1.06L17.94,17.66a.75.75,0,1,0-1.06,1.06ZM19.78,5.28A.75.75,0,1,0,18.72,4.22L17.66,5.28A.75.75,0,1,0,18.72,6.34Z" />
              </svg>
            )}
          </div>
        </div>
      </div>
    </button>
  );
};

export default ThemeToggle;
