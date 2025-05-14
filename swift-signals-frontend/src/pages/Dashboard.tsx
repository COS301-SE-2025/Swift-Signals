import React, { useEffect, useRef } from 'react';
import Navbar from '../components/Navbar';
import '../styles/Dashboard.css';
import { Chart, registerables } from 'chart.js';

// Register Chart.js components
Chart.register(...registerables);

const simulations = [
  { id: '#1234', intersection: 'Main St & 5th Ave', status: 'Complete', statusColor: 'bg-green-400' },
  { id: '#1233', intersection: 'Broadway & 7th St', status: 'Running', statusColor: 'bg-yellow-400' },
  { id: '#1232', intersection: 'Park Ave & 3rd St', status: 'Failed', statusColor: 'bg-red-400' },
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
      // Destroy the existing chart instance if it exists
      if (chartInstanceRef.current) {
        chartInstanceRef.current.destroy();
      }

      // Create a new chart instance
      chartInstanceRef.current = new Chart(chartRef.current, {
        type: 'line',
        data: {
          labels: ['6 AM', '7 AM', '8 AM', '9 AM', '10 AM'],
          datasets: [
            {
              label: 'Traffic Volume',
              data: [5000, 10000, 8000, 12000, 9000],
              borderColor: '#E57373',
              borderWidth: 2,
              fill: false,
              tension: 0.4,
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          scales: {
            y: {
              beginAtZero: true,
              ticks: { stepSize: 5000 },
            },
          },
          plugins: {
            legend: { display: false },
          },
        },
      });
    }

    // Cleanup function to destroy the chart when the component unmounts
    return () => {
      if (chartInstanceRef.current) {
        chartInstanceRef.current.destroy();
        chartInstanceRef.current = null;
      }
    };
  }, []);

  return (
    <div className="min-h-screen bg-gray-100">
      <Navbar />
      <div className="main-content">
        <h1 className="Dashboard-h1">Dashboard Overview</h1>
        <p className="Dashboard-p">Monitor and manage traffic signal operations</p>

        {/* Summary Cards */}
        <div className="card-grid">
          <div className="card">
            <div className="card-icon">
              <span className="text-blue-600">📍</span>
            </div>
            <div>
              <h3 className="card-h3">Total Intersections</h3>
              <p className="card-p">24</p>
            </div>
          </div>
          <div className="card">
            <div className="card-icon">
              <span className="text-green-600">▶</span>
            </div>
            <div>
              <h3 className="card-h3">Active Simulations</h3>
              <p className="card-p">8</p>
            </div>
          </div>
          <div className="card">
            <div className="card-icon">
              <span className="text-purple-600">📈</span>
            </div>
            <div>
              <h3 className="card-h3">Optimization Runs</h3>
              <p className="card-p">156</p>
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="quick-actions">
          <button className="quick-action-button bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700">
            + New Intersection
          </button>
          <button className="quick-action-button bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700">
            ▶ Run Simulation
          </button>
          <button className="quick-action-button bg-purple-600 text-white px-4 py-2 rounded-lg hover:bg-gray-700">
            ≡ View Map
          </button>
        </div>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {/* Recent Simulations */}
          <div className="bg-white p-4 rounded-lg shadow-md">
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
                      <span className={`px-2 py-1 rounded-full text-xs ${sim.statusColor}`}>
                        {sim.status}
                      </span>
                    </td>
                    <td className="p-2">
                      <button className="text-blue-600 hover:underline">View Details</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Traffic Volume Chart and Top Intersections */}
          <div className="space-y-4">
            <div className="bg-white p-4 rounded-lg shadow-md">
              <h2 className="text-lg font-semibold text-gray-800 mb-4">Traffic Volume</h2>
              <div className="h-48">
                <canvas ref={chartRef}></canvas>
              </div>
            </div>
            <div className="bg-white p-4 rounded-lg shadow-md">
              <h2 className="text-lg font-semibold text-gray-800 mb-4">Top Intersections</h2>
              {topIntersections.map((intersection, index) => (
                <div key={index} className="flex justify-between py-2 border-t">
                  <span className="text-gray-600">{intersection.name}</span>
                  <span className="text-gray-800 font-semibold">{intersection.volume}</span>
                </div>
              ))}
              <div className="flex justify-between py-2 border-t">
                <span className="text-gray-600 font-semibold">Avg Daily Volume:</span>
                <span className="text-gray-800 font-semibold">12,000 vehicles</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;