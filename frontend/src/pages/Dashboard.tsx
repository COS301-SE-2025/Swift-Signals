import React, { useEffect, useRef } from "react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import HelpMenu from "../components/HelpMenu"; // MODIFICATION: Import the new component
import "../styles/Dashboard.css";
import { Chart, registerables } from "chart.js";

// Icons used in this component
import {
  FaRoad,
  FaPlay,
  FaChartLine,
  FaPlus,
  FaMap,
  FaArrowRight,
} from "react-icons/fa";

// Register Chart.js components
Chart.register(...registerables);

// Data remains the same
const simulations = [
  {
    id: "#1234",
    intersection: "Main St & 5th Ave",
    status: "Complete",
    statusColor: "bg-statusGreen",
    textColor: "text-statusTextGreen",
  },
  {
    id: "#1233",
    intersection: "Broadway & 7th St",
    status: "Running",
    statusColor: "bg-statusYellow",
    textColor: "text-statusTextYellow",
  },
  {
    id: "#1232",
    intersection: "Park Ave & 3rd St",
    status: "Failed",
    statusColor: "bg-statusRed",
    textColor: "text-statusTextRed",
  },
  {
    id: "#1231",
    intersection: "Broadway & 7th St",
    status: "Running",
    statusColor: "bg-statusYellow",
    textColor: "text-statusTextYellow",
  },
];

const topIntersections = [
  { name: "Main St & 5th Ave", volume: 15000, volumeText: "15,000 vehicles" },
  { name: "Broadway & 7th St", volume: 13500, volumeText: "13,500 vehicles" },
  { name: "Park Ave & 3rd St", volume: 12000, volumeText: "12,000 vehicles" },
];

const Dashboard: React.FC = () => {
  const chartRef = useRef<HTMLCanvasElement | null>(null);
  const chartInstanceRef = useRef<Chart | null>(null);

  useEffect(() => {
    if (chartRef.current) {
      if (chartInstanceRef.current) {
        chartInstanceRef.current.destroy();
      }
      const ctx = chartRef.current.getContext("2d");
      if (!ctx) return;
      const gradient = ctx.createLinearGradient(0, 0, 0, 200);
      gradient.addColorStop(0, "rgba(153, 25, 21, 0.4)");
      gradient.addColorStop(1, "rgba(153, 25, 21, 0)");
      chartInstanceRef.current = new Chart(ctx, {
        type: "line",
        data: {
          labels: ["6 AM", "7 AM", "8 AM", "9 AM", "10 AM"],
          datasets: [
            {
              label: "Traffic Volume",
              data: [5000, 10000, 8000, 12000, 9000],
              fill: true,
              backgroundColor: gradient,
              borderColor: "#991915",
              borderWidth: 2.5,
              pointRadius: 0,
              pointHoverRadius: 8,
              pointHoverBackgroundColor: "#991915",
              pointHoverBorderColor: "#fff",
              tension: 0.4,
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          interaction: { mode: "index", intersect: false },
          scales: {
            x: {
              grid: { display: false },
              ticks: {
                color: "#6B7280",
                font: { size: 12, family: "'Inter', sans-serif" },
              },
              border: { display: false },
            },
            y: {
              grid: { color: "#E5E7EB", drawTicks: false },
              ticks: {
                color: "#6B7280",
                stepSize: 2500,
                font: { size: 12, family: "'Inter', sans-serif" },
                padding: 10,
              },
              border: { display: false },
            },
          },
          plugins: {
            legend: { display: false },
            tooltip: {
              enabled: true,
              backgroundColor: "#111827",
              titleColor: "#F9FAFB",
              bodyColor: "#E5E7EB",
              cornerRadius: 8,
              padding: 12,
              titleFont: {
                weight: "bold",
                size: 14,
                family: "'Inter', sans-serif",
              },
              bodyFont: { size: 12, family: "'Inter', sans-serif" },
              displayColors: false,
              caretPadding: 10,
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

  const maxVolume = Math.max(...topIntersections.map((i) => i.volume), 0);

  return (
    <div className="dashboard-screen min-h-screen bg-gray-100 dark:bg-gray-900">
      <Navbar />
      <main className="main-content">
        <div className="card-grid">
          <div className="card">
            <div className="card-icon-1">
              <span className="text-blue-600">
                <FaRoad />
              </span>
            </div>
            <div>
              <h3 className="card-h3">Total Intersections</h3>
              <p className="card-p">24</p>
            </div>
          </div>
          <div className="card">
            <div className="card-icon-2">
              <span className="text-green-600">
                <FaPlay />
              </span>
            </div>
            <div>
              <h3 className="card-h3">Active Simulations</h3>
              <p className="card-p">8</p>
            </div>
          </div>
          <div className="card">
            <div className="card-icon-3">
              <span className="text-purple-600">
                <FaChartLine />
              </span>
            </div>
            <div>
              <h3 className="card-h3">Optimization Runs</h3>
              <p className="card-p">156</p>
            </div>
          </div>
        </div>
        <div className="dashboard-main-grid">
          <div className="main-column">
            {/* --- MODIFICATION START --- */}
            {/* Replaced 'flex flex-wrap' with a responsive grid layout */}
            <div className="grid grid-cols-2 xl:grid-cols-3 gap-4 mb-6">
              <button className="quick-action-button bg-customIndigo text-white flex items-center justify-center gap-2">
                <FaPlus /> New Intersection
              </button>
              <button className="quick-action-button bg-customGreen text-white flex items-center justify-center gap-2">
                <FaPlay /> Run Simulation
              </button>
              {/* This button now spans 2 columns on smaller screens and 1 on extra-large screens */}
              <button className="quick-action-button bg-customPurple text-white flex items-center justify-center gap-2 col-span-2 xl:col-span-1">
                <FaMap /> View Map
              </button>
            </div>
            {/* --- MODIFICATION END --- */}
            <div className="recent-simulations-tab bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md">
              <h2 className="text-xl font-semibold text-gray-800 dark:text-white mb-4">
                Recent Simulations
              </h2>
              <div className="overflow-x-auto">
                <table className="table-auto w-full text-left">
                  <thead>
                    <tr className="text-gray-600 dark:text-gray-400">
                      <th className="p-2">ID</th>
                      <th className="p-2">Intersection</th>
                      <th className="p-2">Status</th>
                      <th className="p-2">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {simulations.map((sim) => (
                      <tr
                        key={sim.id}
                        className="border-t dark:border-gray-700"
                      >
                        <td className="p-2 text-gray-700 dark:text-gray-200">
                          {sim.id}
                        </td>
                        <td className="p-2 text-gray-700 dark:text-gray-200">
                          {sim.intersection}
                        </td>
                        <td className="p-2">
                          <span
                            className={`status px-2 py-1 rounded-full text-xs ${sim.statusColor} ${sim.textColor}`}
                          >
                            {sim.status}
                          </span>
                        </td>
                        <td className="p-2">
                          <button className="view-details-button text-blue-600 hover:underline">
                            View Details
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          <div className="side-column">
            <div className="graph-card bg-white dark:bg-gray-800 rounded-lg shadow-md">
              <div className="graph-card-header">
                <h2 className="text-xl font-semibold text-gray-800 dark:text-white">
                  Traffic Volume
                </h2>
                <button className="view-report-button">
                  View Report <FaArrowRight />
                </button>
              </div>
              <div className="traffic-chart-container">
                <canvas ref={chartRef}></canvas>
              </div>
            </div>

            <div className="inter-card bg-white dark:bg-gray-800 rounded-lg shadow-md">
              <div className="inter-card-header">
                <h3 className="text-xl font-semibold text-gray-800 dark:text-white">
                  Top Intersections
                </h3>
              </div>
              <div className="intersection-list">
                {topIntersections.map((intersection, index) => {
                  const percentage =
                    maxVolume > 0 ? (intersection.volume / maxVolume) * 100 : 0;
                  return (
                    <div key={index} className="intersection-item">
                      <div className="intersection-details">
                        <span className="intersection-name">
                          {intersection.name}
                        </span>
                        <span className="intersection-volume">
                          {intersection.volumeText}
                        </span>
                      </div>
                      <div className="progress-bar-container">
                        <div
                          className="progress-bar"
                          style={{ width: `${percentage}%` }}
                        ></div>
                      </div>
                    </div>
                  );
                })}
              </div>
              <div className="inter-card-footer">
                <span>Avg Daily Volume:</span>
                <span className="font-semibold">12,000 vehicles</span>
              </div>
            </div>
          </div>
        </div>
      </main>
      <Footer />
      <HelpMenu />
    </div>
  );
};

export default Dashboard;
