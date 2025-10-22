//import { useState } from 'react';
import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import AuthGuard from "./components/AuthGuard";

//import Navbar from './components/Navbar';
// import Footer from './components/Footer';
import ErrorBoundary from "./components/ErrorBoundary";
import ComparisonView from "./pages/ComparisonView";
import Dashboard from "./pages/Dashboard";
import Intersections from "./pages/Intersections";
import Login from "./pages/Login";
import SignUp from "./pages/SignUp";
import SimulationResults from "./pages/SimulationResults";
import Simulations from "./pages/Simulations";
import Users from "./pages/Users";
import WelcomePage from "./pages/WelcomePage";
//import reactLogo from './assets/react.svg';
//import viteLogo from '/vite.svg';
import { UserProvider } from "./context/UserContext";
import "./App.css";

function App() {
  return (
    <UserProvider>
      <Router>
        {/* <Navbar /> */}
        <ErrorBoundary>
          <Routes>
            <Route path="/" element={<WelcomePage />} />
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<SignUp />} />
            <Route
              path="/dashboard"
              element={
                <AuthGuard>
                  <Dashboard />
                </AuthGuard>
              }
            />
            <Route
              path="/simulations"
              element={
                <AuthGuard>
                  <Simulations />
                </AuthGuard>
              }
            />
            <Route
              path="/intersections"
              element={
                <AuthGuard>
                  <Intersections />
                </AuthGuard>
              }
            />
            <Route
              path="/Users"
              element={
                <AuthGuard requiredRole="admin">
                  <Users />
                </AuthGuard>
              }
            />
            <Route
              path="/simulation-results"
              element={
                <AuthGuard>
                  <SimulationResults />
                </AuthGuard>
              }
            />
            <Route
              path="/simulation-results/:intersectionId"
              element={
                <AuthGuard>
                  <SimulationResults />
                </AuthGuard>
              }
            />
            <Route
              path="/comparison-rendering"
              element={
                <AuthGuard>
                  <ComparisonView />
                </AuthGuard>
              }
            />
            {/* Add more routes as needed */}
          </Routes>
        </ErrorBoundary>
        {/* <Footer /> */}
      </Router>
    </UserProvider>
  );
}

export default App;
