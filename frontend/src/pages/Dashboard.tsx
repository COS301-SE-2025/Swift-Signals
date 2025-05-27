import React, { useEffect, useRef } from 'react';
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';
import '../styles/Dashboard.css';
import { Chart, registerables } from 'chart.js';

//icons
import { FaRoad } from "react-icons/fa";
import { FaPlay } from "react-icons/fa";
import { FaChartLine } from "react-icons/fa6";
import { FaPlus } from "react-icons/fa";
import { FaMap } from "react-icons/fa";

// Register Chart.js components
Chart.register(...registerables);

const simulations = [
  { id: '#1234', intersection: 'Main St & 5th Ave', status: 'Complete', statusColor: 'bg-statusGreen', textColor: 'text-statusTextGreen' },
  { id: '#1233', intersection: 'Broadway & 7th St', status: 'Running', statusColor: 'bg-statusYellow', textColor: 'text-statusTextYellow' },
  { id: '#1232', intersection: 'Park Ave & 3rd St', status: 'Failed', statusColor: 'bg-statusRed', textColor: 'text-statusTextRed' },
  { id: '#1231', intersection: 'Broadway & 7th St', status: 'Running', statusColor: 'bg-statusYellow', textColor: 'text-statusTextYellow' }
];

const topIntersections = [
  { name: 'Main St & 5th Ave', volume: '15,000 vehicles' },
  { name: 'Broadway & 7th St', volume: '13,500 vehicles' },
  { name: 'Park Ave & 3rd St', volume: '12,000 vehicles' },
];

const Dashboard: React.FC = () => {
  const chartRef = useRef<HTMLCanvasElement | null>(null);
  const chartInstanceRef = useRef<Chart | null>(null);

  useEffect(() => {
  if (chartRef.current) {
    if (chartInstanceRef.current) {
      chartInstanceRef.current.destroy();
    }

    const ctx = chartRef.current.getContext('2d');
    if (!ctx) return;

    // Create gradient fill
    const gradient = ctx.createLinearGradient(0, 0, 0, 180);
    gradient.addColorStop(0, 'rgba(153, 25, 21, 0.3)');
    gradient.addColorStop(1, 'rgba(153, 25, 21, 0)');

    chartInstanceRef.current = new Chart(ctx, {
      type: 'line',
      data: {
        labels: ['6 AM', '7 AM', '8 AM', '9 AM', '10 AM'],
        datasets: [
          {
            label: 'Traffic Volume',
            data: [5000, 10000, 8000, 12000, 9000],
            fill: true,
            backgroundColor: gradient,
            borderColor: '#991915',
            borderWidth: 3,
            pointBackgroundColor: '#991915',
            pointBorderColor: '#fff',
            pointHoverRadius: 6,
            pointRadius: 4,
            pointHoverBackgroundColor: '#fff',
            pointHoverBorderColor: '#991915',
            tension: 0.4,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        layout: {
          padding: {
            top: 10,
            bottom: 10,
            left: 0,
            right: 0,
          },
        },
        scales: {
          x: {
            grid: {
              display: false,
            },
            ticks: {
              color: '#6B7280',
              font: {
                size: 14,
                weight: 500,
              },
            },
            border: {
              display: false,
            },
          },
          y: {
            grid: {
              color: '#E5E7EB',
              drawTicks: false,
            },
            ticks: {
              color: '#6B7280',
              stepSize: 2000,
              font: {
                size: 14,
                weight: 500,
              },
            },
            border: {
              display: false,
            },
          },
        },
        plugins: {
          legend: {
            display: false,
          },
          tooltip: {
            backgroundColor: '#111827', // Tailwind's gray-900
            titleColor: '#F9FAFB', // Tailwind's gray-50
            bodyColor: '#E5E7EB', // Tailwind's gray-200
            cornerRadius: 4,
            padding: 10,
            titleFont: {
              weight: 'bold',
              size: 14,
            },
            bodyFont: {
              size: 13,
            },
          },
        },
      },
    });
  }

  return () => {
    if (chartInstanceRef.current) {
      chartInstanceRef.current.destroy();
      chartInstanceRef.current = null;
    }
  };
}, []);


  return (
    <div className="dashboard-screen min-h-screen bg-gray-100 dark:bg-gray-900">
      <Navbar />
      <div className="main-content flex-grow">
        <h1 className="Dashboard-h1">Dashboard Overview</h1>
        {/* <p className="Dashboard-p">Monitor and manage traffic signal operations</p> */}

        {/* Summary Cards */}
        <div className="card-grid">
          <div className="card">
            <div className="card-icon-1">
              <span className="text-blue-600"><FaRoad /></span>
            </div>
            <div>
              <h3 className="card-h3">Total Intersections</h3>
              <p className="card-p">24</p>
            </div>
          </div>
          <div className="card">
            <div className="card-icon-2">
              <span className="text-green-600"><FaPlay /></span>
            </div>
            <div>
              <h3 className="card-h3">Active Simulations</h3>
              <p className="card-p">8</p>
            </div>
          </div>
          <div className="card">
            <div className="card-icon-3">
              <span className="text-purple-600"><FaChartLine /></span>
            </div>
            <div>
              <h3 className="card-h3">Optimization Runs</h3>
              <p className="card-p">156</p>
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="quick-actions">
          <button className="quick-action-button bg-customIndigo text-white px-4 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2">
            <FaPlus /> {/* Add Plus Icon */}
            New Intersection
          </button>
          <button className="quick-action-button bg-customGreen text-white px-4 py-2 rounded-lg hover:bg-green-700 flex items-center gap-2">
            <FaPlay /> {/* Add Play Icon */}
            Run Simulation
          </button>
          <button className="quick-action-button bg-customPurple text-white px-4 py-2 rounded-lg hover:bg-gray-700 flex items-center gap-2">
            <FaMap /> {/* Add Map Icon */}
            View Map
          </button>
        </div>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {/* Recent Simulations */}
          <div className="recent-simulations-tab bg-white p-4 rounded-lg shadow-md">
            <h2 className="text-lg font-semibold text-gray-800 mb-4">Recent Simulations</h2>
            <table className="table-auto w-full text-left">
              <thead>
                <tr className="text-gray-600">
                  <th className="p-2">ID</th>
                  <th className="p-2">Intersection</th>
                  <th className="p-2">Status</th>
                  <th className="p-2">Actions</th>
                </tr>
              </thead>
              <tbody>
                {simulations.map((sim) => (
                  <tr key={sim.id} className="border-t">
                    <td className="p-2">{sim.id}</td>
                    <td className="p-2">{sim.intersection}</td>
                    <td className="p-2">
                      <span className={`status px-2 py-1 rounded-full text-xs ${sim.statusColor} ${sim.textColor}`}>
                        {sim.status}
                      </span>
                    </td>
                    <td className="p-2">
                      <button className="view-details-button text-blue-600 hover:underline">View Details</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
        {/* Traffic Volume Chart and Top Intersections */}
        <div className="stats bg-white p-4 rounded-lg shadow-md">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">Traffic Volume</h2>
          <div className="traffic-graph mb-4">
            <div className="traffic-chart">
              <canvas ref={chartRef}></canvas>
            </div>
          </div>
          <div className="top-intersections">
            <h3 className="text-md font-semibold text-gray-700 mb-2">Top Intersections</h3>
            {topIntersections.map((intersection, index) => (
              <div key={index} className="flex justify-between py-2 border-t">
            <span className="text-gray-600 dark:text-gray-200">{intersection.name}</span>
            <span className="text-gray-800 font-semibold dark:text-gray-200">{intersection.volume}</span>
              </div>
            ))}
            <div className="total flex justify-between py-2 border-t">
              <span className="text-gray-600 font-bold dark:text-gray-100">Avg Daily Volume:</span>
              <span className="text-gray-800 font-bold dark:text-gray-100">12,000 vehicles</span>
            </div>
          </div>
        </div>
      </div>
      <Footer />
    </div>
  );
};

export default Dashboard;