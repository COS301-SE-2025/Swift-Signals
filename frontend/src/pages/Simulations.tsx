
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';
import { Bar } from 'react-chartjs-2';
import { Chart, registerables } from 'chart.js';
import { MapContainer, TileLayer, Marker, useMapEvents } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import type { LatLng } from 'leaflet';
import '../styles/Simulations.css';
import '@fortawesome/fontawesome-free/css/all.min.css';

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

const LocationMarker: React.FC<{ setSelectedLocation: (location: string) => void; setCoordinates: (coords: string) => void }> = ({ setSelectedLocation, setCoordinates }) => {
  const [position, setPosition] = useState<LatLng | null>(null);
  useMapEvents({
    click(e) {
      setPosition(e.latlng);
      const coordinates = `${e.latlng.lat.toFixed(4)}, ${e.latlng.lng.toFixed(4)}`;
      setSelectedLocation(coordinates);
      setCoordinates(coordinates);
      console.log(`Marker placed at coordinates: ${coordinates}`);
    },
  });

  return position === null ? null : <Marker position={position} />;
};

// New Simulation Modal Component
const NewSimulationModal: React.FC<{
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: { name: string; description: string; intersections: string[] }) => void;
  intersections: string[];
  type: 'simulations' | 'optimizations';
}> = ({ isOpen, onClose, onSubmit, intersections, type }) => {
  const navigate = useNavigate();
  const [simulationName, setSimulationName] = useState('');
  const [simulationDescription, setSimulationDescription] = useState('');
  const [selectedIntersections, setSelectedIntersections] = useState<string[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [activeTab, setActiveTab] = useState<'List' | 'Search' | 'Map'>('List');
  const [coordinates, setCoordinates] = useState<string | null>(null);

  // Handle adding an intersection from the List or Search tab
  const handleAddIntersection = (intersection: string) => {
    if (intersection && !selectedIntersections.includes(intersection)) {
      setSelectedIntersections([...selectedIntersections, intersection]);
    }
  };

  // Handle removing an intersection
  const handleRemoveIntersection = (intersection: string) => {
    setSelectedIntersections(selectedIntersections.filter((item) => item !== intersection));
  };

  // Handle search action
  const handleSearch = () => {
    if (searchQuery && !selectedIntersections.includes(searchQuery)) {
      console.log(`Searching for location: ${searchQuery}`);
      setSelectedIntersections([...selectedIntersections, searchQuery]);
      setSearchQuery('');
    }
  };

  // Handle map click to add coordinates as an intersection
  const handleMapSelection = (location: string) => {
    if (!selectedIntersections.includes(location)) {
      setSelectedIntersections([...selectedIntersections, location]);
    }
  };

  const handleSubmit = () => {
    if (!simulationName || selectedIntersections.length === 0) {
      alert('Please provide a simulation name and select at least one intersection.');
      return;
    }
    const simulationData = { name: simulationName, description: simulationDescription, intersections: selectedIntersections };
    onSubmit(simulationData);
    setSimulationName('');
    setSimulationDescription('');
    setSelectedIntersections([]);
    setSearchQuery('');
    setCoordinates(null);
    navigate('/simulation-results', { state: simulationData });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-md p-6 relative">
        {/* Close Button */}
        <button
          onClick={onClose}
          className="crossBtn absolute top-4 right-4 text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-100"
        >
          ✕
        </button>

        {/* Modal Header */}
        <h2 className="text-xl font-bold text-gray-800 dark:text-gray-200 mb-4">
          New {type === 'simulations' ? 'Simulation' : 'Optimization'}
        </h2>

        {/* Form */}
        <div className="space-y-4">
          {/* Simulation Name */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Simulation Name
            </label>
            <input
              type="text"
              value={simulationName}
              onChange={(e) => setSimulationName(e.target.value)}
              className="w-full p-2 rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="Enter simulation name"
            />
          </div>

          {/* Simulation Description */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Simulation Description
            </label>
            <textarea
              value={simulationDescription}
              onChange={(e) => setSimulationDescription(e.target.value)}
              className="w-full p-2 rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="Enter simulation description"
              rows={3}
            />
          </div>

          {/* Intersection Selection */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Intersections
            </label>

            {/* Selected Intersections as Pills */}
            <div className="flex flex-wrap gap-2 mb-3">
              {selectedIntersections.map((intersection) => (
                <div
                  key={intersection}
                  className="intersection-pill flex items-center px-2 py-0.5 rounded-full bg-indigo-100 text-indigo-800 dark:bg-indigo-700 dark:text-indigo-100 text-xs"
                >
                  {intersection}
                  <button
                    onClick={() => handleRemoveIntersection(intersection)}
                    className="ml-1 text-indigo-600 hover:text-indigo-800 dark:text-indigo-300 dark:hover:text-indigo-100 remove-cross"
                  >
                    ✕
                  </button>
                </div>
              ))}
            </div>

            {/* Tabs */}
            <div className="flex space-x-2 mb-3">
              <button
                onClick={() => setActiveTab('List')}
                className={`px-3 py-1 rounded-md text-sm font-medium ${activeTab === 'List' ? 'bg-indigo-600 text-white dark:bg-indigo-500' : 'bg-gray-200 text-gray-700 dark:bg-gray-600 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-500'} transition-all duration-300`}
              >
                List
              </button>
              <button
                onClick={() => setActiveTab('Search')}
                className={`px-3 py-1 rounded-md text-sm font-medium ${activeTab === 'Search' ? 'bg-indigo-600 text-white dark:bg-indigo-500' : 'bg-gray-200 text-gray-700 dark:bg-gray-600 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-500'} transition-all duration-300`}
              >
                Search
              </button>
              <button
                onClick={() => setActiveTab('Map')}
                className={`px-3 py-1 rounded-md text-sm font-medium ${activeTab === 'Map' ? 'bg-indigo-600 text-white dark:bg-indigo-500' : 'bg-gray-200 text-gray-700 dark:bg-gray-600 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-500'} transition-all duration-300`}
              >
                Map
              </button>
            </div>

            {/* Tab Content */}
            {activeTab === 'List' && (
              <select
                value=""
                onChange={(e) => handleAddIntersection(e.target.value)}
                className="w-full p-2 rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
                <option value="">Select an intersection</option>
                {intersections.map((intersection) => (
                  <option key={intersection} value={intersection}>
                    {intersection}
                  </option>
                ))}
              </select>
            )}

            {activeTab === 'Search' && (
              <div className="flex space-x-2">
                <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-full p-2 rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="Search for a location"
                />
                <button
                  onClick={handleSearch}
                  className="px-4 py-2 rounded-md text-sm font-medium bg-indigo-600 text-white hover:bg-indigo-700 dark:bg-indigo-500 dark:hover:bg-indigo-600 transition-all duration-300"
                >
                  Add
                </button>
              </div>
            )}

            {activeTab === 'Map' && (
              <div>
                <MapContainer center={[-26.2041, 28.0473]} zoom={6} style={{ height: '200px', width: '100%' }}>
                  <TileLayer
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    attribution='© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                  />
                  <LocationMarker setSelectedLocation={handleMapSelection} setCoordinates={setCoordinates} />
                </MapContainer>
                {coordinates && (
                  <p className="mt-2 text-sm text-gray-700 dark:text-gray-300">
                    Last Clicked Coordinates: <span className="font-medium">{coordinates}</span>
                  </p>
                )}
              </div>
            )}
          </div>
        </div>

        {/* Modal Footer */}
        <div className="mt-6 flex justify-end space-x-2">
          <button
            onClick={onClose}
            className="px-4 py-2 rounded-md text-sm font-medium bg-gray-200 text-gray-700 hover:bg-gray-300 dark:bg-gray-600 dark:text-gray-200 dark:hover:bg-gray-500 transition-all duration-300"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            className="px-4 py-2 rounded-md text-sm font-medium bg-gradient-to-r from-green-500 to-green-600 text-white hover:from-green-600 hover:to-green-700 dark:from-green-400 dark:to-green-500 dark:hover:from-green-500 dark:hover:to-green-600 transition-all duration-300"
          >
            Create
          </button>
        </div>
      </div>
    </div>
  );
};

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
    <div className="simTable bg-white dark:bg-gray-800 shadow-md rounded-lg overflow-hidden table-fixed-height relative">
      <table className="simulationTable min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead className="simTableHead bg-gray-50 dark:bg-gray-700">
          <tr>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Simulation ID</th>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Intersection</th>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Avg Wait Time</th>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Throughput</th>
            <th className="graphTHead px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Graph</th>
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
                <td className="intersectionCell px-4 py-3 whitespace-wrap text-sm text-gray-900 dark:text-gray-200">{sim.intersection}</td>
                <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-gray-200">{sim.avgWaitTime.toFixed(1)}</td>
                <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-gray-200">{sim.vehicleThroughput}</td>
                <td className="chartCell px-4 py-3 whitespace-nowrap text-sm">
                  <div className="h-16 w-24">
                    <Bar data={chartData} options={chartOptions} />
                  </div>
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-sm">
                  <span className={`sim-status inline-flex items-center px-3 py-1 rounded-full border ${statusClass(sim.status)}`}>
                    {sim.status}
                  </span>
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-sm">
                  <div className="flex flex-col space-y-2">
                    <button
                      onClick={() => handleViewResults(sim.id)}
                      className="viewBtn text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300 text-sm font-medium w-full text-center" title="View Results"
                    >
                      <i className="fas fa-eye"></i>
                    </button>
                    <button
                      onClick={() => handleDelete(sim.id)}
                      className="deleteBtn text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300 text-sm font-medium w-full text-center" title="Delete Simulation"
                    >
                      <i className="fas fa-trash"></i>
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
            Prev
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
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalType, setModalType] = useState<'simulations' | 'optimizations'>('simulations');

  // Filter simulations based on selected intersection
  const filteredSimulations1 = filter1 === 'All Intersections' ? simulationsTable1 : simulationsTable1.filter(sim => sim.intersection === filter1);
  const filteredSimulations2 = filter2 === 'All Intersections' ? simulationsTable2 : simulationsTable2.filter(sim => sim.intersection === filter2);

  // Get all unique intersections for the dropdown
  const allIntersections = Array.from(new Set([...simulationsTable1, ...simulationsTable2].map(sim => sim.intersection)));

  const handleNewSimulation = (type: 'simulations' | 'optimizations') => {
    setModalType(type);
    setIsModalOpen(true);
  };

  const handleModalSubmit = (data: { name: string; description: string; intersections: string[] }) => {
    console.log(`New ${modalType === 'simulations' ? 'Simulation' : 'Optimization'} Created:`, data);
    // Replace with actual logic, e.g., API call to save the new simulation/optimization
  };

  return (
    <div className="simulationBody min-h-screen bg-gray-100 dark:bg-gray-900">
      <Navbar />
      <div className="sim-main-content flex-grow p-6">
        <div className="simGrid grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="simTableContainer">
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
          <div className="simTableContainer">
            <div className="flex items-center justify-between mb-4">
              <h1 className="text-3xl font-bold text-gray-800 dark:text-gray-200">Recent Optimizations</h1>
              <div className="flex items-center space-x-2">
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
      <NewSimulationModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleModalSubmit}
        intersections={allIntersections}
        type={modalType}
      />
    </div>
  );
};

export default Simulations;

