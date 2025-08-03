import React, { useEffect, useRef } from "react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import HelpMenu from "../components/HelpMenu";
import "../styles/Dashboard.css";
import { Chart, registerables } from "chart.js";
import MapModal from "../components/MapModal";
import { useState } from "react";

import {
  FaRoad,
  FaPlay,
  FaChartLine,
  FaPlus,
  FaMap,
  FaArrowRight,
} from "react-icons/fa";

Chart.register(...registerables);

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
  const [isMapModalOpen, setIsMapModalOpen] = useState(false);
  const [mapIntersections, setMapIntersections] = useState<any[]>([]);
  const [isLoadingMap, setIsLoadingMap] = useState(false);
  const [mapError, setMapError] = useState<string | null>(null);

  useEffect(() => {
    if (chartRef.current) {
      if (chartInstanceRef.current) {
        chartInstanceRef.current.destroy();
      }
      const ctx = chartRef.current.getContext("2d");
      if (!ctx) return;

      const isDarkMode = document.documentElement.classList.contains('dark');
      
      const gradient = ctx.createLinearGradient(0, 0, 0, 200);
      gradient.addColorStop(0, "rgba(15, 91, 167, 0.4)");
      gradient.addColorStop(1, "rgba(15, 91, 167, 0)");

      const labelColor = isDarkMode ? "#F0F6FC" : "#6B7280";
      const gridColor = isDarkMode ? "#30363D" : "#E5E7EB";
      
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
              borderColor: "#0F5BA7",
              borderWidth: 2.5,
              pointRadius: 0,
              pointHoverRadius: 8,
              pointHoverBackgroundColor: "#0F5BA7",
              pointHoverBorderColor: "#0066CC",
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
                color: labelColor,
                font: { size: 12, family: "'Inter', sans-serif" },
              },
              border: { display: false },
            },
            y: {
              grid: { color: gridColor, drawTicks: false },
              ticks: {
                color: labelColor,
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

  // Fetch intersections for map modal
  const fetchMapIntersections = async () => {
    setIsLoadingMap(true);
    setMapError(null);
    try {
      const token = localStorage.getItem("authToken");
      const response = await fetch("http://localhost:9090/intersections", {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      if (!response.ok) throw new Error("Failed to fetch intersections");
      const data = await response.json();
      const intersectionsWithCoords = (data.intersections || []).map((intr: any, idx: number) => ({
        ...intr,
        details: {
          ...intr.details,
          latitude: intr.details.latitude ?? (-25.7479 + 0.01 * idx),
          longitude: intr.details.longitude ?? (28.2293 + 0.01 * idx),
        },
      }));
      setMapIntersections(intersectionsWithCoords);
    } catch (err: any) {
      setMapError(err.message || "Unknown error");
      setMapIntersections([]);
    } finally {
      setIsLoadingMap(false);
    }
  };

  const handleOpenMapModal = () => {
    setIsMapModalOpen(true);
    fetchMapIntersections();
  };
  const handleCloseMapModal = () => setIsMapModalOpen(false);

  return (
    <div className="dashboard-screen min-h-screen bg-gray-100 dark:bg-gray-900">
      <Navbar />
      <main className="main-content">
        <div className="card-grid">
          <div className="card">
            <div className="card-icon-1">
              <span className="text-[#0F5BA7]">
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
              <span className="text-[#0F5BA7]">
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
              <span className="text-[#0F5BA7]">
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
            <div className="grid grid-cols-2 xl:grid-cols-3 gap-4 mb-6">
              <button className="newInt quick-action-button bg-[#0F5BA7] dark:bg-[#388BFD] text-white dark:text-[#E6EDF3] flex items-center justify-center gap-2">
                <FaPlus /> New Intersection
              </button>
              <button className="runSim quick-action-button bg-[#2B9348] dark:bg-[#2DA44E] text-white dark:text-[#E6EDF3] flex items-center justify-center gap-2">
                <FaPlay /> Run Simulation
              </button>
              <button className="viewMap quick-action-button border-2 border-[#0F5BA7] dark:border-[#388BFD] text-[#0F5BA7] dark:text-[#388BFD] bg-white dark:bg-[#0D1117] hover:bg-[#e6f1fa] transition flex items-center justify-center gap-2 col-span-2 xl:col-span-1" onClick={handleOpenMapModal}>
                <FaMap /> View Map
              </button>
            </div>
            <div className="recent-simulations-tab bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md">
              <h2 className="text-xl font-semibold text-gray-800 dark:text-[#E6EDF3] mb-4">
                Recent Simulations
              </h2>
              <div className="overflow-x-auto">
                <table className="table-auto w-full text-left">
                  <thead>
                    <tr className="text-gray-600 dark:text-[#C9D1D9]">
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
                        className="border-t dark:border-[#30363D]"
                      >
                        <td className="p-2 text-gray-700 dark:text-[#F0F6FC]">
                          {sim.id}
                        </td>
                        <td className="p-2 text-gray-700 dark:text-[#F0F6FC]">
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
                          <button className="view-details-button text-[#0F5BA7] dark:text-[#388BFD] hover:underline border-none">
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
                <h2 className="text-xl font-semibold text-gray-800 dark:text-[#E6EDF3]">
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
                <h3 className="text-xl font-semibold text-gray-800 dark:text-[#E6EDF3]">
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
      <MapModal isOpen={isMapModalOpen} onClose={handleCloseMapModal} intersections={mapIntersections} />
      {isLoadingMap && isMapModalOpen && (
        <div className="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-30">
          <div className="bg-white dark:bg-gray-800 p-6 rounded shadow text-center">Loading map data...</div>
        </div>
      )}
      {mapError && isMapModalOpen && (
        <div className="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-30">
          <div className="bg-white dark:bg-gray-800 p-6 rounded shadow text-center text-red-600">{mapError}</div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;