import React, { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import HelpMenu from "../components/HelpMenu";
import "../styles/Dashboard.css";
import { Chart, registerables } from "chart.js";
import MapModal from "../components/MapModal";
import { useState } from "react";

import { FaRoad, FaPlay, FaChartLine, FaPlus, FaMap } from "react-icons/fa";

Chart.register(...registerables);

interface Intersection {
  id: string;
  name: string;
  details: {
    address: string;
    city: string;
    province: string;
    latitude: number;
    longitude: number;
  };
}

interface ApiIntersection {
  id: string;
  name: string;
  created_at?: string;
  run_count?: number;
  status?: string;
  details?: {
    address?: string;
    city?: string;
    province?: string;
    latitude?: number;
    longitude?: number;
    [key: string]: unknown;
  };
  traffic_density?: string; // Added for traffic density
  [key: string]: unknown;
}

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const chartRef = useRef<HTMLCanvasElement | null>(null);
  const chartInstanceRef = useRef<Chart | null>(null);
  const [isMapModalOpen, setIsMapModalOpen] = useState(false);
  const [mapIntersections, setMapIntersections] = useState<Intersection[]>([]);
  const [isLoadingMap, setIsLoadingMap] = useState(false);
  const [mapError, setMapError] = useState<string | null>(null);

  const [recentIntersections, setRecentIntersections] = useState<
    ApiIntersection[]
  >([]);
  const [loadingRecent, setLoadingRecent] = useState(false);
  const [recentError, setRecentError] = useState<string | null>(null);

  const [totalIntersections, setTotalIntersections] = useState<number>(0);
  const [loadingTotal, setLoadingTotal] = useState(false);
  const [activeSimulations, setActiveSimulations] = useState<number>(0);
  const [loadingActiveSimulations, setLoadingActiveSimulations] =
    useState(false);
  const [totalSimulationsRun, setTotalSimulationsRun] = useState<number>(0);
  const [totalOptimizationsRun, setTotalOptimizationsRun] = useState<number>(0);

  const fetchAllIntersections = async () => {
    setLoadingTotal(true);
    // setLoadingActiveSimulations(true); // This line will be removed as activeSimulations is no longer used in the same way
    try {
      const token = localStorage.getItem("authToken");
      const response = await fetch("http://localhost:9090/intersections", {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      if (!response.ok) throw new Error("Failed to fetch intersections");
      const data = await response.json();
      const items: ApiIntersection[] = data.intersections || [];

      setTotalIntersections(items.length);

      let totalSims = 0;
      let totalOptimizations = 0;

      items.forEach(item => {
        const runCount = typeof item.run_count === "number" ? item.run_count : 0;
        totalSims += runCount;

        if (item.status === "INTERSECTION_STATUS_OPTIMISED") {
          totalOptimizations += runCount;
        }
      });

      setTotalSimulationsRun(totalSims);
      setTotalOptimizationsRun(totalOptimizations);

      updateChart(items);
    } catch (err: unknown) {
      console.error("Failed to fetch intersections:", err);
      setTotalIntersections(0);
      setTotalSimulationsRun(0); // Reset on error
      setTotalOptimizationsRun(0); // Reset on error
      if (chartInstanceRef.current) {
        chartInstanceRef.current.destroy();
        chartInstanceRef.current = null;
      }
    } finally {
      setLoadingTotal(false);
      setLoadingActiveSimulations(false); // Keep this to ensure loading state is cleared
    }
  };

  const updateChart = (intersections: ApiIntersection[]) => {
    if (!chartRef.current) return;

    if (chartInstanceRef.current) {
      chartInstanceRef.current.destroy();
    }

    const ctx = chartRef.current.getContext("2d");
    if (!ctx) return;

    const chartData = processTrafficDensityData(intersections);

    const isDarkMode = document.documentElement.classList.contains("dark");

    const gradient = ctx.createLinearGradient(0, 0, 0, 200);
    gradient.addColorStop(0, "rgba(37, 99, 235, 0.4)");
    gradient.addColorStop(1, "rgba(37, 99, 235, 0)");

    const labelColor = isDarkMode ? "#F0F6FC" : "#6B7280";
    const gridColor = isDarkMode ? "#30363D" : "#E5E7EB";

    chartInstanceRef.current = new Chart(ctx, {
      type: "line",
      data: {
        labels: chartData.labels,
        datasets: [
          {
            label: "Traffic Density",
            data: chartData.data,
            fill: true,
            backgroundColor: gradient,
            borderColor: "#0f5ba7",
            borderWidth: 2.5,
            pointRadius: 4,
            pointHoverRadius: 8,
            pointHoverBackgroundColor: "#0f5ba7",
            pointHoverBorderColor: "#0f5ba7",
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
              stepSize: 1,
              font: { size: 12, family: "'Inter', sans-serif" },
              padding: 10,
              callback: function (value) {
                const densityLabels = ["", "Low", "Medium", "High"];
                return densityLabels[value as number] || value;
              },
            },
            border: { display: false },
            min: 0,
            max: 4,
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
            callbacks: {
              title: function (context) {
                return `${context[0].label}`;
              },
              label: function (context) {
                const densityLabels = ["", "Low", "Medium", "High"];
                const densityValue =
                  densityLabels[context.parsed.y as number] || context.parsed.y;
                return `Traffic Density: ${densityValue}`;
              },
            },
          },
        },
      },
    });
  };

  const processTrafficDensityData = (intersections: ApiIntersection[]) => {
    if (intersections.length === 0) {
      return { labels: ["No Data"], data: [0] };
    }

    // Map traffic density categories to numeric values based on actual API response
    const densityMap: { [key: string]: number } = {
      TRAFFIC_DENSITY_LOW: 1,
      TRAFFIC_DENSITY_MEDIUM: 2,
      TRAFFIC_DENSITY_HIGH: 3,
    };

    // Process intersections and map their traffic density to numeric values
    const densityData = intersections.map((intersection) => {
      const density = intersection.traffic_density || "TRAFFIC_DENSITY_MEDIUM";
      const mappedValue = densityMap[density] || 2; // Default to medium (2) if density is unknown

      // Debug logging
      console.log(
        `Intersection: ${intersection.name}, Raw density: ${density}, Mapped value: ${mappedValue}`,
      );

      return mappedValue;
    });

    // Create labels for the x-axis (intersection numbers)
    const labels = intersections.map((_, index) => `Intersection ${index + 1}`);

    console.log("Final density data:", densityData);
    console.log("Final labels:", labels);

    return { labels, data: densityData };
  };

  useEffect(() => {
    fetchAllIntersections();
  }, []);

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

      const intersectionsWithCoords = (data.intersections || [])
        .map((intr: ApiIntersection) => {
          const name = intr.name ?? "";
          const match = name.match(/\[(-?\d+\.?\d*),(-?\d+\.?\d*)\]/);
          if (match) {
            const lat = parseFloat(match[1]);
            const lng = parseFloat(match[2]);
            return {
              id: String(intr.id),
              name: name.split(" [")[0],
              details: {
                address: intr.details?.address ?? "",
                city: intr.details?.city ?? "",
                province: intr.details?.province ?? "",
                latitude: lat,
                longitude: lng,
              },
            };
          }
          return null;
        })
        .filter((intr): intr is Intersection => intr !== null);

      setMapIntersections(intersectionsWithCoords);
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : "Unknown error";
      setMapError(errorMessage);
      setMapIntersections([]);
    } finally {
      setIsLoadingMap(false);
    }
  };

  const fetchRecentIntersections = async () => {
    setLoadingRecent(true);
    setRecentError(null);
    try {
      const token = localStorage.getItem("authToken");
      const response = await fetch("http://localhost:9090/intersections", {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      if (!response.ok) throw new Error("Failed to fetch intersections");
      const data = await response.json();
      const items: ApiIntersection[] = data.intersections || [];

      const sorted = items
        .slice()
        .sort((a, b) => {
          const at = a.created_at ? new Date(a.created_at).getTime() : 0;
          const bt = b.created_at ? new Date(b.created_at).getTime() : 0;
          return bt - at;
        })
        .slice(0, 5);

      setRecentIntersections(sorted);
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : "Unknown error";
      setRecentError(errorMessage);
      setRecentIntersections([]);
    } finally {
      setLoadingRecent(false);
    }
  };

  useEffect(() => {
    fetchRecentIntersections();
  }, []);

  const handleOpenMapModal = () => {
    setIsMapModalOpen(true);
    fetchMapIntersections();
  };
  const handleCloseMapModal = () => setIsMapModalOpen(false);

  const handleSimulateFromMap = (id: string, name: string) => {
    navigate("/simulation-results", {
      state: {
        name: `Simulation Results for ${name}`,
        description: `Viewing simulation results for ${name}`,
        intersections: [name],
        intersectionIds: [id],
        type: "simulations",
      },
    });
  };

  const handleViewDetails = (intersection: ApiIntersection) => {
    navigate("/simulation-results", {
      state: {
        name: `Simulation Results for ${intersection.name}`,
        description: `Viewing simulation results for ${intersection.name}`,
        intersections: [intersection.name],
        intersectionIds: [intersection.id],
        type: "simulations",
      },
    });
  };

  const getStatusStyles = (status?: string) => {
    switch (status) {
      case "Optimised":
        return {
          statusColor: "bg-statusGreen",
          textColor: "text-statusTextGreen",
        };
      case "Unoptimised":
        return {
          statusColor: "bg-statusYellow",
          textColor: "text-statusTextYellow",
        };
      case "Failed":
        return {
          statusColor: "bg-statusRed",
          textColor: "text-statusTextRed",
        };
      default:
        return {
          statusColor: "bg-statusYellow",
          textColor: "text-statusTextYellow",
        };
    }
  };

  const mapApiStatus = (status?: string): string => {
    switch (status) {
      case "INTERSECTION_STATUS_OPTIMISED":
        return "Optimised";
      case "unoptimised":
        return "Unoptimised";
      case "Failed":
        return "Failed";
      default:
        return "Unoptimised";
    }
  };

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
              <p className="card-p">
                {loadingTotal ? "..." : totalIntersections}
              </p>
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
              <p className="card-p">
                {loadingActiveSimulations ? "..." : totalSimulationsRun}
              </p>
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
              <p className="card-p">
                {loadingActiveSimulations ? "..." : totalOptimizationsRun}
              </p>
            </div>
          </div>
        </div>
        <div className="dashboard-main-grid">
          <div className="main-column">
            <div className="grid grid-cols-2 xl:grid-cols-3 gap-4 mb-6">
              <button
                className="newInt quick-action-button bg-[#0F5BA7] dark:bg-[#388BFD] text-white dark:text-[#E6EDF3] flex items-center justify-center gap-2 hover:bg-[#0D4A8A] dark:hover:bg-[#2D7BD8] transition-colors duration-200 shadow-md hover:shadow-lg transform hover:scale-105"
                onClick={() => navigate("/intersections")}
                title="Go to Intersections page to create new intersections"
              >
                <FaPlus /> New Intersection
              </button>
              <button
                className="runSim quick-action-button bg-[#2B9348] dark:bg-[#2DA44E] text-white dark:text-[#E6EDF3] flex items-center justify-center gap-2 hover:bg-[#237A3A] dark:hover:bg-[#238636] transition-colors duration-200 shadow-md hover:shadow-lg transform hover:scale-105"
                onClick={() => navigate("/simulations")}
                title="Go to Simulations page to run new simulations"
              >
                <FaPlay /> Run Simulation
              </button>
              <button
                className="viewMap quick-action-button border-2 border-[#0F5BA7] dark:border-[#388BFD] text-[#0F5BA7] dark:text-[#388BFD] bg-white dark:bg-[#0D1117] hover:bg-[#e6f1fa] transition flex items-center justify-center gap-2 col-span-2 xl:col-span-1"
                onClick={handleOpenMapModal}
              >
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
                      <th className="p-2">#</th>
                      <th className="p-2">Location</th>
                      <th className="p-2">Status</th>
                      <th className="p-2">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {loadingRecent ? (
                      <tr>
                        <td
                          colSpan={4}
                          className="p-4 text-center text-gray-500 dark:text-gray-400"
                        >
                          Loading...
                        </td>
                      </tr>
                    ) : recentError ? (
                      <tr>
                        <td
                          colSpan={4}
                          className="p-4 text-center text-red-500"
                        >
                          Error: {recentError}
                        </td>
                      </tr>
                    ) : recentIntersections.length === 0 ? (
                      <tr>
                        <td
                          colSpan={4}
                          className="p-4 text-center text-gray-500 dark:text-gray-400"
                        >
                          No recent simulations found.
                        </td>
                      </tr>
                    ) : (
                      recentIntersections.map((intr, index) => {
                        const displayStatus = mapApiStatus(intr.status);
                        const styles = getStatusStyles(displayStatus);
                        return (
                          <tr
                            key={intr.id}
                            className="border-t dark:border-[#30363D]"
                          >
                            <td className="p-2 text-gray-700 dark:text-[#F0F6FC]">
                              {index + 1}
                            </td>
                            <td className="p-2 text-gray-700 dark:text-[#F0F6FC]">
                              {(intr.details?.address || intr.name || "Unnamed Intersection").split(',')[0]}
                            </td>
                            <td className="p-2">
                              <span
                                className={`status px-2 py-1 rounded-full text-xs ${styles.statusColor} ${styles.textColor}`}
                              >
                                {displayStatus}
                              </span>
                            </td>
                            <td className="p-2">
                              <button
                                className="view-details-button text-[#0F5BA7] dark:text-[#388BFD] hover:underline border-none"
                                onClick={() => handleViewDetails(intr)}
                              >
                                View Details
                              </button>
                            </td>
                          </tr>
                        );
                      })
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          <div className="side-column">
            <div className="graph-card bg-white dark:bg-gray-800 rounded-lg shadow-md">
              <div className="graph-card-header">
                <h2 className="text-xl font-semibold text-gray-800 dark:text-[#E6EDF3] text-align-center">
                  Traffic Density Distribution
                </h2>
              </div>
              <div className="traffic-chart-container">
                <canvas ref={chartRef}></canvas>
              </div>
            </div>

            {/* === START: UPDATED RECENT INTERSECTIONS CARD === */}
            <div className="inter-card bg-white dark:bg-gray-800 rounded-lg shadow-md flex flex-col h-full">
              {/* Card Header */}
              <div className="p-4 border-b border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-800 dark:text-[#E6EDF3]">
                  Recent Intersections
                </h3>
              </div>

              {/* Intersection List */}
              <div className="intersection-list flex-grow overflow-y-auto">
                {/* Loading State */}
                {loadingRecent && (
                  <div className="flex items-center justify-center h-full p-4">
                    <p className="text-gray-500 dark:text-gray-400">
                      Loading...
                    </p>
                  </div>
                )}

                {/* Error State */}
                {recentError && (
                  <div className="flex items-center justify-center h-full p-4">
                    <p className="text-red-500">Error: {recentError}</p>
                  </div>
                )}

                {/* Empty State */}
                {!loadingRecent &&
                  !recentError &&
                  recentIntersections.length === 0 && (
                    <div className="flex items-center justify-center h-full p-4">
                      <p className="text-gray-500 dark:text-gray-400">
                        No intersections found.
                      </p>
                    </div>
                  )}

                {/* Data List */}
                {!loadingRecent && !recentError && (
                  <ul className="divide-y divide-gray-200 dark:divide-gray-700">
                    {recentIntersections.map((intr) => {
                      const displayName = (intr.name || "Unnamed Intersection").split(' [')[0];
                      const address =
                        (intr.details?.address ||
                        [intr.details?.city, intr.details?.province]
                          .filter(Boolean)
                          .join(", ") ||
                        "No address provided").split(',')[0];

                      const displayStatus = mapApiStatus(intr.status);
                      const styles = getStatusStyles(displayStatus);
                      const formattedDate = intr.created_at
                        ? new Date(intr.created_at).toLocaleDateString(
                            "en-US",
                            {
                              month: "short",
                              day: "numeric",
                            },
                          )
                        : "";

                      return (
                        <li
                          key={intr.id}
                          className="p-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors duration-150 cursor-pointer"
                          onClick={() => handleViewDetails(intr)}
                        >
                          <div className="flex items-center justify-between gap-4">
                            <div className="flex items-center gap-3 min-w-0">
                              <div className="flex-shrink-0 h-10 w-10 flex items-center justify-center bg-blue-100 dark:bg-blue-900/50 rounded-lg">
                                <FaRoad className="h-5 w-5 text-[#0F5BA7] dark:text-[#388BFD]" />
                              </div>
                              <div className="min-w-0">
                                <p className="text-sm font-semibold text-gray-900 dark:text-gray-100 truncate">
                                  {displayName}
                                </p>
                                <p className="text-xs text-gray-500 dark:text-gray-400 truncate">
                                  {address}
                                </p>
                              </div>
                            </div>
                            <div className="text-right flex-shrink-0">
                              <div className="flex items-center justify-end gap-2">
                                <div
                                  className={`h-2 w-2 rounded-full ${styles.statusColor}`}
                                  title={displayStatus}
                                />
                                <p
                                  className={`hidden sm:block text-xs font-medium capitalize ${styles.textColor}`}
                                >
                                  {displayStatus}
                                </p>
                              </div>
                              <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">
                                {formattedDate}
                              </p>
                            </div>
                          </div>
                        </li>
                      );
                    })}
                  </ul>
                )}
              </div>

              {/* Card Footer */}
              <div className="inter-card-footer p-2 border-t border-gray-200 dark:border-gray-700 mt-auto">
                <button
                  className="w-full text-center text-sm font-medium text-[#0F5BA7] dark:text-[#388BFD] bg-white dark:bg-transparent hover:underline p-2 rounded-md transition-colors hover:bg-blue-50 dark:hover:bg-blue-900/20 border border-gray-200 dark:border-transparent"
                  onClick={() => navigate("/intersections")}
                >
                  View All ({totalIntersections})
                </button>
              </div>
            </div>
            {/* === END: UPDATED RECENT INTERSECTIONS CARD === */}
          </div>
        </div>
      </main>
      <Footer />
      <HelpMenu />
      <MapModal
        isOpen={isMapModalOpen}
        onClose={handleCloseMapModal}
        intersections={mapIntersections}
        onSimulate={handleSimulateFromMap}
      />
      {isLoadingMap && isMapModalOpen && (
        <div className="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-30">
          <div className="bg-white dark:bg-gray-800 p-6 rounded shadow text-center">
            Loading map data...
          </div>
        </div>
      )}
      {mapError && isMapModalOpen && (
        <div className="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-30">
          <div className="bg-white dark:bg-gray-800 p-6 rounded shadow text-center text-red-600">
            {mapError}
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;
