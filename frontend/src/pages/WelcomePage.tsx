import React from "react";
import "../styles/WelcomePage.css"; // Make sure this path is correct
import Carousel from "../components/Carousel";

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
          <div style={{ height: "300px", position: "relative" }}>
            <Carousel
              baseWidth={350}
              autoplay={true}
              autoplayDelay={5000}
              pauseOnHover={true}
              loop={true}
              round={false}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default WelcomePage;
