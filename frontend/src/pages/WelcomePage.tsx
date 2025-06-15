import React from "react";
import "../styles/WelcomePage.css"; // Make sure this path is correct

const WelcomePage: React.FC = () => {
  const handleLoginClick = () => {
    console.log("Login button clicked!");
    // Add navigation logic here, e.g., history.push('/login');
  };

  const handleRegisterClick = () => {
    console.log("Register button clicked!");
    // Add navigation logic here, e.g., history.push('/register');
  };

  return (
    <div className="welcome-page">
      <div className="glass-block">
        <div className="glass-block-left">
          {/* Replace with your actual logo path */}
          <img src="../../src/assets/logo.png" alt="Logo" />
          <h1>Welcome to Swift Signals!</h1>
          <div className="auth-buttons">
            <button className="modern-button" onClick={handleLoginClick}>
              Login
            </button>
            <button
              className="modern-button register"
              onClick={handleRegisterClick}
            >
              Register
            </button>
          </div>
        </div>
        <div className="glass-block-right">
          <p>
            Swift Signals is a simulation-driven traffic light optimization
            platform developed with Southern Cross Solutions to combat urban
            congestion. Designed for municipal traffic departments, it uses
            machine learning and historical traffic data to analyse
            intersections and optimize signal timing. With traffic congestion
            costing South Africa an estimated R1 billion annually in lost
            productivity, Swift Signals offers a scalable, modular web platform
            that simulates real-world traffic patterns and adjusts signal phases
            dynamically. Built using microservices, containerization, and CI/CD
            pipelines, it ensures long-term maintainability and deployment
            efficiency.
          </p>
        </div>
      </div>
    </div>
  );
};

export default WelcomePage;
