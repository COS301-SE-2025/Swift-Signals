import React, { useEffect, useRef } from 'react';
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';
import '../styles/SimulationResults.css';
import { Chart, registerables } from 'chart.js';

// Register Chart.js components
Chart.register(...registerables);

const SimulationResults: React.FC = () => {
  const simulationName = "Traffic Flow Simulation";
  const simulationDescription = "A simulation analyzing traffic patterns and intersections.";
  const intersections = ["Intersection A", "Intersection B", "Intersection C"];

  return (
    <div className="min-h-screen bg-gray-100 dark:bg-gray-900">
      <Navbar />
      <div className="simRes-main-content flex-grow p-6">
        {/* Simulation Header */}
        <div className="mb-6 p-4 bg-white dark:bg-gray-700 rounded-lg shadow">
          <h1 className="text-2xl font-bold">{simulationName}</h1>
          <p className="text-gray-600 dark:text-gray-300">{simulationDescription}</p>
          <h3 className="text-lg font-semibold mt-2">Intersections:</h3>
          <ul className="list-disc list-inside">
            {intersections.map((intersection, index) => (
              <li key={index}>{intersection}</li>
            ))}
          </ul>
        </div>

        {/* Main Grid Layout */}
        <div className="simRes-grid grid grid-cols-2 gap-6">
          {/* Simulation Visualization */}
          <div className="results bg-white dark:bg-gray-700 rounded-lg shadow p-4 flex items-center justify-center h-64">
            <svg className="w-16 h-16 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 16l4-4m0 0l4 4m-4-4v8m-6-6h12m-6-6a2 2 0 100-4 2 2 0 000 4z"></path>
            </svg>
            <h2 className="text-xl font-semibold mt-2">Simulation Visualization</h2>
          </div>

          {/* Optimize Button */}
          <div className="flex items-center justify-center h-64">
            <button className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">Optimize</button>
          </div>

          {/* Optimized Visualization */}
          <div className="results bg-white dark:bg-gray-700 rounded-lg shadow p-4 flex items-center justify-center h-64">
            <svg className="w-16 h-16 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 16l4-4m0 0l4 4m-4-4v8m-6-6h12m-6-6a2 2 0 100-4 2 2 0 000 4z"></path>
            </svg>
            <h2 className="text-xl font-semibold mt-2">Optimized Visualization</h2>
          </div>

          {/* Parameters Section */}
          <div className="bg-gray-200 dark:bg-gray-700 rounded-lg shadow p-4">
            <h2 className="text-xl font-semibold mb-4">Parameters</h2>
            <div className="space-y-2">
              <button className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">Apply Traffic Volume</button>
              <button className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">Apply Signal Timing</button>
              <button className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">Apply Road Conditions</button>
            </div>
          </div>
        </div>

        {/* Statistics Section */}
        <div className="mt-6 bg-gray-200 dark:bg-gray-700 rounded-lg shadow p-4">
          <h2 className="text-xl font-semibold mb-4">Statistics</h2>
          <div className="h-64">
            {/* Placeholder for chart/graph */}
            <canvas id="statisticsChart"></canvas>
          </div>
          <p className="mt-2">Average Traffic Flow: 1200 vehicles/hour</p>
          <p>Congestion Index: 0.75</p>
        </div>
      </div>
      <Footer />
    </div>
  );
};

export default SimulationResults;