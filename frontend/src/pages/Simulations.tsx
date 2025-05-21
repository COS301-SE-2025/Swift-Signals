import React, { useEffect, useRef, useState } from 'react';
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';
import { Bar } from 'react-chartjs-2';
import { Chart, registerables } from 'chart.js';
import '../styles/Simulations.css';

// Register Chart.js components
Chart.register(...registerables);

// Sample simulation data for two tables with added status
const simulationsTable1 = [
  { id: "SIM001", intersection: "Main St & 1st Ave", avgWaitTime: 45.2, vehicleThroughput: 1200, status: "Complete" },
  { id: "SIM002", intersection: "Broadway & 5th St", avgWaitTime: 30.8, vehicleThroughput: 1500, status: "Running" },
  { id: "SIM003", intersection: "Elm St & Park Rd", avgWaitTime: 52.1, vehicleThroughput: 900, status: "Failed" },
  { id: "SIM004", intersection: "Main St & 1st Ave", avgWaitTime: 45.2, vehicleThroughput: 1200, status: "Complete" },
  { id: "SIM005", intersection: "Broadway & 5th St", avgWaitTime: 30.8, vehicleThroughput: 1500, status: "Running" },
  { id: "SIM006", intersection: "Elm St & Park Rd", avgWaitTime: 52.1, vehicleThroughput: 900, status: "Failed" },
];

const simulationsTable2 = [
  { id: "SIM004", intersection: "Oak Ave & Central Blvd", avgWaitTime: 28.5, vehicleThroughput: 1800, status: "Complete" },
  { id: "SIM005", intersection: "Pine St & River Dr", avgWaitTime: 47.3, vehicleThroughput: 1100, status: "Running" },
  { id: "SIM006", intersection: "Maple Rd & 2nd Ave", avgWaitTime: 35.6, vehicleThroughput: 1300, status: "Failed" },
  { id: "SIM012", intersection: "Oak Ave & Central Blvd", avgWaitTime: 28.5, vehicleThroughput: 1800, status: "Complete" },
  { id: "SIM020", intersection: "Pine St & River Dr", avgWaitTime: 47.3, vehicleThroughput: 1100, status: "Running" },
  { id: "SIM007", intersection: "Maple Rd & 2nd Ave", avgWaitTime: 35.6, vehicleThroughput: 1300, status: "Failed" },
];

// Simulation Table Component
const SimulationTable: React.FC<{ simulations: Array<{ id: string; intersection: string; avgWaitTime: number; vehicleThroughput: number; status: string }>, currentPage: number, setCurrentPage: (page: number) => void }> = ({ simulations, currentPage, setCurrentPage }) => {
  const rowsPerPage = 4;
  const totalPages = Math.ceil(simulations.length / rowsPerPage);
  const startIndex = currentPage * rowsPerPage;
  const endIndex = startIndex + rowsPerPage;
  const paginatedSimulations = simulations.slice(startIndex, endIndex);

  // Chart options for a modern look
  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: {
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        cornerRadius: 8,
        padding: 10,
        titleFont: { size: 12, weight: 'bold' as const },
        bodyFont: { size: 12 },
        displayColors: false,
      },
    },
    scales: {
      x: { display: false },
      y: { beginAtZero: true, display: false },
    },
    animation: {
      duration: 1000,
      easing: 'easeOutQuart' as const,
    },
    elements: {
      bar: {
        borderRadius: 6,
        borderWidth: 0,
      },
    },
  };

  const handleViewResults = (simId: string) => {
    alert(`Viewing results for simulation ${simId}`);
    // Replace with actual logic, e.g., navigate to results page
  };

  const handleDelete = (simId: string) => {
    alert(`Deleting simulation ${simId}`);
    // Replace with actual delete logic, e.g., API call
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  // Function to determine status class
  const statusClass = (status: string) => {
    switch (status) {
      case 'Complete':
        return 'bg-green-200 text-green-800 border-green-300';
      case 'Running':
        return 'bg-yellow-200 text-yellow-800 border-yellow-300';
      case 'Failed':
        return 'bg-red-200 text-red-800 border-red-300';
      default:
        return 'bg-gray-200 text-gray-800 border-gray-300';
    }
  };

  return (
    <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg overflow-hidden table-fixed-height relative">
      <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead className="bg-gray-50 dark:bg-gray-700">
          <tr>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Simulation ID</th>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Intersection</th>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Avg Wait Time</th>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Throughput</th>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Graph</th>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Status</th>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
          {paginatedSimulations.map((sim) => {
            // Create gradient for bar colors
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');
            const gradientWaitTime = ctx?.createLinearGradient(0, 0, 0, 100);
            gradientWaitTime?.addColorStop(0, '#34D399');
            gradientWaitTime?.addColorStop(1, '#10B981');

            const gradientThroughput = ctx?.createLinearGradient(0, 0, 0, 100);
            gradientThroughput?.addColorStop(0, '#3B82F6');
            gradientThroughput?.addColorStop(1, '#2563EB');

            // Chart data for the row
            const chartData = {
              labels: ['Wait', 'Throughput'],
              datasets: [
                {
                  data: [sim.avgWaitTime, sim.vehicleThroughput / 10],
                  backgroundColor: [gradientWaitTime, gradientThroughput],
                  hoverBackgroundColor: ['#6EE7B7', '#60A5FA'],
                  borderWidth: 0,
                },
              ],
            };

            return (
              <tr key={sim.id}>
                <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-gray-200">{sim.id}</td>
                <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-gray-200">{sim.intersection}</td>
                <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-gray-200">{sim.avgWaitTime.toFixed(1)}</td>
                <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-gray-200">{sim.vehicleThroughput}</td>
                <td className="px-4 py-3 whitespace-nowrap text-sm">
                  <div className="h-16 w-24">
                    <Bar data={chartData} options={chartOptions} />
                  </div>
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-sm">
                  <span className={`inline-flex items-center px-3 py-1 rounded-full border ${statusClass(sim.status)}`}>
                    {sim.status}
                  </span>
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-sm">
                  <div className="flex flex-col space-y-2">
                    <button
                      onClick={() => handleViewResults(sim.id)}
                      className="viewBtn text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300 text-sm font-medium w-full text-center"
                    >
                      View
                    </button>
                    <button
                      onClick={() => handleDelete(sim.id)}
                      className="deleteBtn text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300 text-sm font-medium w-full text-center"
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
      {simulations.length > rowsPerPage && (
        <div className="absolute bottom-0 left-0 right-0 flex justify-center items-center p-4 space-x-2 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
          <button
            onClick={() => handlePageChange(currentPage - 1)}
            disabled={currentPage === 0}
            className={`px-3 py-1 rounded-full text-sm font-medium bg-gradient-to-r from-indigo-500 to-indigo-600 text-white hover:from-indigo-600 hover:to-indigo-700 dark:from-indigo-400 dark:to-indigo-500 dark:hover:from-indigo-500 dark:hover:to-indigo-600 transition-all duration-300 ${currentPage === 0 ? 'opacity-50 cursor-not-allowed' : ''}`}
          >
            Previous
          </button>
          {Array.from({ length: totalPages }, (_, index) => (
            <button
              key={index}
              onClick={() => handlePageChange(index)}
              className={`px-3 py-1 rounded-full text-sm font-medium ${currentPage === index ? 'bg-indigo-600 text-white dark:bg-indigo-500' : 'bg-gray-200 text-gray-700 dark:bg-gray-600 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-500'} transition-all duration-300`}
            >
              {index + 1}
            </button>
          ))}
          <button
            onClick={() => handlePageChange(currentPage + 1)}
            disabled={currentPage === totalPages - 1}
            className={`px-3 py-1 rounded-full text-sm font-medium bg-gradient-to-r from-indigo-500 to-indigo-600 text-white hover:from-indigo-600 hover:to-indigo-700 dark:from-indigo-400 dark:to-indigo-500 dark:hover:from-indigo-500 dark:hover:to-indigo-600 transition-all duration-300 ${currentPage === totalPages - 1 ? 'opacity-50 cursor-not-allowed' : ''}`}
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
};

const Simulations: React.FC = () => {
  const [filter1, setFilter1] = useState<string>('All Intersections');
  const [filter2, setFilter2] = useState<string>('All Intersections');
  const [page1, setPage1] = useState<number>(0);
  const [page2, setPage2] = useState<number>(0);

  // Filter simulations based on selected intersection
  const filteredSimulations1 = filter1 === 'All Intersections' ? simulationsTable1 : simulationsTable1.filter(sim => sim.intersection === filter1);
  const filteredSimulations2 = filter2 === 'All Intersections' ? simulationsTable2 : simulationsTable2.filter(sim => sim.intersection === filter2);

  const handleNewSimulation = (table: 'simulations' | 'optimizations') => {
    alert(`Creating new ${table === 'simulations' ? 'simulation' : 'optimization'}`);
    // Replace with actual logic, e.g., navigate to a form or open a modal
  };

  return (
    <div className="min-h-screen bg-gray-100 dark:bg-gray-900">
      <Navbar />
      <div className="sim-main-content flex-grow p-6">
        <div className="simGrid grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <div className="flex items-center justify-between mb-4">
              <h1 className="text-3xl font-bold text-gray-800 dark:text-gray-200">Recent Simulations</h1>
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => handleNewSimulation('simulations')}
                  className="px-4 py-2 rounded-md text-sm font-medium bg-gradient-to-r from-green-500 to-green-600 text-white hover:from-green-600 hover:to-green-700 dark:from-green-400 dark:to-green-500 dark:hover:from-green-500 dark:hover:to-green-600 transition-all duration-300 shadow-md hover:shadow-lg"
                >
                  New Simulation
                </button>
                <select
                  value={filter1}
                  onChange={(e) => { setFilter1(e.target.value); setPage1(0); }}
                  className="w-48 p-2 rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                  {['All Intersections', ...new Set(simulationsTable1.map(sim => sim.intersection))].map((intersection) => (
                    <option key={intersection} value={intersection}>
                      {intersection}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <SimulationTable simulations={filteredSimulations1} currentPage={page1} setCurrentPage={setPage1} />
          </div>
          <div>
            <div className="flex items-center justify-between mb-4">
              <h1 className="text-3xl font-bold text-gray-800 dark:text-gray-200">Recent Optimizations</h1>
              <div className="flex items-center space-x-2">
                {/* <button
                  onClick={() => handleNewSimulation('optimizations')}
                  className="px-4 py-2 rounded-md text-sm font-medium bg-gradient-to-r from-green-500 to-green-600 text-white hover:from-green-600 hover:to-green-700 dark:from-green-400 dark:to-green-500 dark:hover:from-green-500 dark:hover:to-green-600 transition-all duration-300 shadow-md hover:shadow-lg"
                >
                  New Optimization
                </button> */}
                <select
                  value={filter2}
                  onChange={(e) => { setFilter2(e.target.value); setPage2(0); }}
                  className="w-48 p-2 rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                  {['All Intersections', ...new Set(simulationsTable2.map(sim => sim.intersection))].map((intersection) => (
                    <option key={intersection} value={intersection}>
                      {intersection}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <SimulationTable simulations={filteredSimulations2} currentPage={page2} setCurrentPage={setPage2} />
          </div>
        </div>
      </div>
      <Footer />
    </div>
  );
};

export default Simulations;