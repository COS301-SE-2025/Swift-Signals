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
} from "react-icons/fa";

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
  [key: string]: unknown;
}

const Dashboard: React.FC = () => {
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

  const [allIntersections, setAllIntersections] = useState<ApiIntersection[]>(
    []
  );

  const fetchAllIntersections = async () => {
    setLoadingTotal(true);
    try {
      const token = localStorage.getItem("authToken");
      const response = await fetch("http://localhost:9090/intersections", {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      if (!response.ok) throw new Error("Failed to fetch intersections");
      const data = await response.json();
      const items: ApiIntersection[] = data.intersections || [];

      setAllIntersections(items);
      setTotalIntersections(items.length);

      updateChart(items);
    } catch (err: unknown) {
      console.error("Failed to fetch intersections:", err);
      setTotalIntersections(0);
      setAllIntersections([]);
      if (chartInstanceRef.current) {
        chartInstanceRef.current.destroy();
        chartInstanceRef.current = null;
      }
    } finally {
      setLoadingTotal(false);
    }
  };

  const updateChart = (intersections: ApiIntersection[]) => {
    if (!chartRef.current) return;

    if (chartInstanceRef.current) {
      chartInstanceRef.current.destroy();
    }

    const ctx = chartRef.current.getContext("2d");
    if (!ctx) return;

    const chartData = processRunCountData(intersections);

    const isDarkMode = document.documentElement.classList.contains("dark");

    const gradient = ctx.createLinearGradient(0, 0, 0, 200);
    gradient.addColorStop(0, "rgba(37, 99, 235, 0.4)");
    gradient.addColorStop(1, "rgba(37, 99, 235, 0)");

    const labelColor = isDarkMode ? "#F0F6FC" : "#6B7280";
    const gridColor = isDarkMode ? "#30363D" : "#E5E7EB";

    const maxY = Math.max(...chartData.data);
    const step = Math.max(1, Math.ceil(maxY / 5) || 1);

    chartInstanceRef.current = new Chart(ctx, {
      type: "line",
      data: {
        labels: chartData.labels,
        datasets: [
          {
            label: "Optimization Runs (Cumulative)",
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
              stepSize: step,
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
            callbacks: {
              title: function (context) {
                return `${context[0].label}`;
              },
              label: function (context) {
                return `Total Runs: ${context.parsed.y}`;
              },
            },
          },
        },
      },
    });
  };

  const processRunCountData = (intersections: ApiIntersection[]) => {
    const withDates = intersections.filter((i) => i.created_at);
    if (withDates.length === 0) {
      return { labels: ["No Data"], data: [0] };
    }

    const sorted = withDates.sort(
      (a, b) =>
        new Date(a.created_at!).getTime() - new Date(b.created_at!).getTime()
    );

    const perDay: Record<string, number> = {};
    sorted.forEach((i) => {
      const label = new Date(i.created_at!).toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
      });
      perDay[label] =
        (perDay[label] || 0) + (typeof i.run_count === "number" ? i.run_count : 0);
    });

    const labels = Object.keys(perDay);
    const data: number[] = [];
    let running = 0;
    labels.forEach((d) => {
      running += perDay[d];
      data.push(running);
    });

    const last = labels.slice(-5);
    const lastData = data.slice(-5);

    if (last.length === 0) return { labels: ["No Data"], data: [0] };

    while (last.length < 5) {
      const nextDate = new Date();
      nextDate.setDate(nextDate.getDate() + last.length);
      last.push(
        nextDate.toLocaleDateString("en-US", { month: "short", day: "numeric" })
      );
      lastData.push(lastData[lastData.length - 1]);
    }

    return { labels: last, data: lastData };
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
      const intersectionsWithCoords: Intersection[] = (
        data.intersections || []
      ).map((intr: ApiIntersection, idx: number) => ({
        id: String(intr.id),
        name: intr.name ?? "Unnamed",
        details: {
          address: intr.details?.address ?? "",
          city: intr.details?.city ?? "",
          province: intr.details?.province ?? "",
          latitude: intr.details?.latitude ?? -25.7479 + 0.01 * idx,
          longitude: intr.details?.longitude ?? 28.2293 + 0.01 * idx,
        },
      }));
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

  const getStatusStyles = (status?: string) => {
    switch (status) {
      case "optimised":
        return {
          statusColor: "bg-statusGreen",
          textColor: "text-statusTextGreen",
        };
      case "unoptimised":
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
      case "optimised":
        return "optimised";
      case "unoptimised":
        return "unoptimised";
      case "Failed":
        return "Failed";
      default:
        return "unoptimised";
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
                      <th className="p-2">Intersection</th>
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
                              {intr.details?.address ||
                                intr.name ||
                                "Unnamed Intersection"}
                            </td>
                            <td className="p-2">
                              <span
                                className={`status px-2 py-1 rounded-full text-xs ${styles.statusColor} ${styles.textColor}`}
                              >
                                {displayStatus}
                              </span>
                            </td>
                            <td className="p-2">
                              <button className="view-details-button text-[#0F5BA7] dark:text-[#388BFD] hover:underline border-none">
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
                  Optimization Runs Over Time
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
                      const address =
                        intr.details?.address ||
                        [intr.details?.city, intr.details?.province]
                          .filter(Boolean)
                          .join(", ") ||
                        "No address provided";

                      const displayStatus = mapApiStatus(intr.status);
                      const styles = getStatusStyles(displayStatus);
                      const formattedDate = intr.created_at
                        ? new Date(intr.created_at).toLocaleDateString(
                            "en-US",
                            {
                              month: "short",
                              day: "numeric",
                            }
                          )
                        : "";

                      return (
                        <li
                          key={intr.id}
                          className="p-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors duration-150 cursor-pointer"
                        >
                          <div className="flex items-center justify-between gap-4">
                            <div className="flex items-center gap-3 min-w-0">
                              <div className="flex-shrink-0 h-10 w-10 flex items-center justify-center bg-blue-100 dark:bg-blue-900/50 rounded-lg">
                                <FaRoad className="h-5 w-5 text-[#0F5BA7] dark:text-[#388BFD]" />
                              </div>
                              <div className="min-w-0">
                                <p className="text-sm font-semibold text-gray-900 dark:text-gray-100 truncate">
                                  {intr.name || "Unnamed Intersection"}
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
                <button className="w-full text-center text-sm font-medium text-[#0F5BA7] dark:text-[#388BFD] hover:underline p-2 rounded-md transition-colors hover:bg-blue-50 dark:hover:bg-blue-900/20">
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