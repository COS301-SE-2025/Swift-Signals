import React, { useEffect, useRef } from 'react';
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';
import '../styles/Simulations.css';

const carouselItems = [
  { id: 1, title: 'Simulation #1', intersection: 'Main St & 1st Ave', avgWaitTime: '15s', throughput: '120 vehicles/min', status: 'Completed' },
  { id: 2, title: 'Simulation #2', intersection: 'Church St & Park Rd', avgWaitTime: '--', throughput: '--', status: 'Pending' },
  { id: 3, title: 'Simulation #3', intersection: 'University Rd & Lynnwood Rd', avgWaitTime: '10s', throughput: '100 vehicles/min', status: 'Completed' },
  { id: 4, title: 'Simulation #4', intersection: 'Main St & 1st Ave', avgWaitTime: '--', throughput: '--', status: 'Pending' },
  { id: 5, title: 'Simulation #5', intersection: 'Main St & 1st Ave', avgWaitTime: '--', throughput: '--', status: 'Failed' },
];

const Carousel: React.FC = () => {
  const [currentIndex, setCurrentIndex] = React.useState(0);
  const carouselRef = useRef<HTMLDivElement>(null);

  // Ensure the carousel loops around
  const getDisplayItems = () => {
    const items = [];
    const totalItems = carouselItems.length;
    for (let i = 0; i < 3; i++) {
      const index = (currentIndex + i) % totalItems;
      items.push(carouselItems[index]);
    }
    return items;
  };

  const handlePrev = () => {
    setCurrentIndex((prev) => {
      const newIndex = prev - 1;
      return newIndex < 0 ? carouselItems.length - 1 : newIndex;
    });
  };

  const handleNext = () => {
    setCurrentIndex((prev) => (prev + 1) % carouselItems.length);
  };

  // Determine status color based on status field
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'Completed':
        return 'bg-green-300';
      case 'Pending':
        return 'bg-yellow-300';
      case 'Failed':
        return 'bg-red-300';
      default:
        return 'bg-gray-400'; // Fallback for unknown status
    }
  };

  return (
    <div className="carContainer relative w-full max-w-8xl mx-auto py-0 flex items-center">
      <button
        onClick={handlePrev}
        className="mr-6 bg-gradient-to-r from-gray-800 to-gray-700 dark:from-gray-700 dark:to-gray-600 text-white p-4 rounded-full shadow-md hover:from-gray-700 hover:to-gray-600 dark:hover:from-gray-600 dark:hover:to-gray-500 focus:outline-none transition-all duration-300 ease-in-out hover:scale-110"
      >
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 19l-7-7 7-7"></path>
        </svg>
      </button>
      <div className="flex-1 overflow-hidden">
        <div
          ref={carouselRef}
          className="flex transition-transform duration-500 ease-in-out"
          style={{ transform: `translateX(-${(currentIndex % carouselItems.length) * (100 / 3)}%)` }}
        >
          {carouselItems.concat(carouselItems).map((item, index) => (
            <div
              key={`${item.id}-${index}`}
              className="w-1/3 flex-shrink-0 px-4"
            >
              <div className="c-card bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden relative">
                {/* Status Circle */}
                <div
                  className={`absolute top-4 right-4 w-8 h-8 rounded-full ${getStatusColor(item.status)}`}
                  title={item.status}
                ></div>
                <div className="p-6 text-left">
                  <h3 className="title text-xl font-bold text-gray-900 dark:text-white">
                    {item.title}
                  </h3>
                  <p className="mt-2 text-xl text-gray-600 dark:text-gray-300">
                    Intersection: {item.intersection}
                  </p>
                  <p className="mt-2 text-xl text-gray-600 dark:text-gray-300">
                    Avg Wait Time: {item.avgWaitTime}
                  </p>
                  <p className="mt-2 text-xl text-gray-600 dark:text-gray-300">
                    Throughput: {item.throughput}
                  </p>
                  {/* Buttons */}
                  <div className="mt-4 flex space-x-4">
                    <button
                      className="view-results-btn px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center space-x-2"
                      onClick={() => console.log(`View Results for ${item.title}`)}
                    >
                      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      </svg>
                      <span>View Results</span>
                    </button>
                    <button
                      className="delete-btn px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors flex items-center space-x-2"
                      onClick={() => console.log(`Delete ${item.title}`)}
                    >
                      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5-4h4a1 1 0 011 1v1H9V4a1 1 0 011-1zm-7 4h18" />
                      </svg>
                      <span>Delete</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
      <button
        onClick={handleNext}
        className="ml-6 bg-gradient-to-r from-gray-800 to-gray-700 dark:from-gray-700 dark:to-gray-600 text-white p-4 rounded-full shadow-md hover:from-gray-700 hover:to-gray-600 dark:hover:from-gray-600 dark:hover:to-gray-500 focus:outline-none transition-all duration-300 ease-in-out hover:scale-110"
      >
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5l7 7-7 7"></path>
        </svg>
      </button>
    </div>
  );
};

const Simulations: React.FC = () => {
  // State for dropdown selection
  const [selectedIntersection, setSelectedIntersection] = React.useState('All Intersections');

  // Handler for dropdown change
  const handleIntersectionChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setSelectedIntersection(event.target.value);
  };

  // Handler for running new simulation
  const handleRunSimulation = () => {
    console.log(`Running new simulation for: ${selectedIntersection}`);
  };
    return (
    <div className="min-h-screen bg-gray-100 dark:bg-gray-900">
      <Navbar />
      <div className="main-content flex-grow">
        <section className="carousel py-12">
          <Carousel />
          <div className="controls mt-8 flex justify-center items-center space-x-4">
            <label className="lbl text-lg font-semibold text-gray-800 dark:text-gray-200">
              Filter By:
            </label>
            <div className="dropdown relative inline-flex">
              <select
                value={selectedIntersection}
                onChange={handleIntersectionChange}
                className="appearance-none bg-maroon text-white px-6 py-2 rounded-full border border-black focus:outline-none w-164"
              >
                <option value="All Intersections">All Intersections</option>
                {carouselItems.map((item) => (
                  <option key={item.id} value={item.intersection}>
                    {item.intersection}
                  </option>
                ))}
              </select>
              <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-white">
                <svg
                  className="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M19 9l-7 7-7-7"
                  ></path>
                </svg>
              </div>
            </div>
            <button
              onClick={handleRunSimulation}
              className="run-btn px-6 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
            >
              Run New Simulation
            </button>
          </div>
        </section>
      </div>
      <Footer />
    </div>
  );
};

export default Simulations;