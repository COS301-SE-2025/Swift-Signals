import React from "react";
import { useNavigate } from "react-router-dom";
import "../styles/WelcomePage.css";
import Carousel from "../components/Carousel";
import logo from "../../src/assets/logo.png"

const WelcomePage: React.FC = () => {
  const navigate = useNavigate(); // Initialize useNavigate

  const handleLoginClick = () => {
    console.log("Login button clicked!");
    navigate("/login"); // Navigate to the Login page
  };

  const handleRegisterClick = () => {
    console.log("Register button clicked!");
    navigate("/signup"); // Navigate to the Signup page
  };

  return (
    <div className="welcome-page">
      <div className="glass-block">
        <div className="glass-block-left">
          {/* Replace with your actual logo path */}
          <img src={logo} alt="Logo" />
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
              autoplayDelay={9000}
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
