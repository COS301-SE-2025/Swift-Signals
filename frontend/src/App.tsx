//import { useState } from 'react';
// import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
//import Navbar from './components/Navbar';
// import Footer from './components/Footer';
import Login from "./pages/Login";
import SignUp from "./pages/SignUp";
import Dashboard from "./pages/Dashboard";
import Simulations from "./pages/Simulations";
import Intersections from "./pages/Intersections";
import Users from "./pages/Users";
import SimulationResults from "./pages/SimulationResults";
import ErrorBoundary from "./components/ErrorBoundary";
import WelcomePage from "./pages/WelcomePage";
import ComparisonView from "./pages/ComparisonView";
//import reactLogo from './assets/react.svg';
//import viteLogo from '/vite.svg';
import "./App.css";

function App() {
  return (
    <Router>
      {/* <Navbar /> */}
      <ErrorBoundary>
        <Routes>
          <Route path="/" element={<WelcomePage />} />
          <Route path="/login" element={<Login />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/logout" element={<Login />} />
          <Route path="/signup" element={<SignUp />} />
          <Route path="/simulations" element={<Simulations />} />
          <Route path="/intersections" element={<Intersections />} />
          <Route path="/Users" element={<Users />} />
          <Route path="/simulation-results" element={<SimulationResults />} />
          <Route path="/comparison-rendering" element={<ComparisonView />} />
          {/* Add more routes as needed */}
          {/* <Route path="/about" element={<About />} /> */}
          {/* <Route path="/contact" element={<Contact />} /> */}
          {/* Add more routes as needed */}
        </Routes>
      </ErrorBoundary>
      {/* <Footer /> */}
    </Router>
  );
}

export default App;
