//import { useState } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
//import Navbar from './components/Navbar';
import Footer from './components/Footer';
import Login from './pages/Login';
import SignUp from './pages/SignUp';
import Dashboard from './pages/Dashboard';
import Simulations from './pages/Simulations';
import ErrorBoundary from './components/ErrorBoundary';
//import reactLogo from './assets/react.svg';
//import viteLogo from '/vite.svg';
import './App.css';

function App() {
  return (
    <Router>
      {/* <Navbar /> */}
      <ErrorBoundary>
        <Routes>
          <Route path="/" element={<Login />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path='/logout' element={<Login />} />
          <Route path="/signup" element={<SignUp />} />
          <Route path="/simulations" element={<Simulations />} />
          {/* Add more routes as needed */}
        </Routes>
      </ErrorBoundary>
      <Footer />
    </Router>
  );
}

export default App;