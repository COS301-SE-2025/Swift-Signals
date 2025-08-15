import React, { useEffect, useState, useRef } from "react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import "../styles/SimulationResults.css";
import HelpMenu from "../components/HelpMenu";
import { Chart, registerables } from "chart.js";
import { useLocation } from "react-router-dom";
import type { ChartConfiguration } from "chart.js";

Chart.register(...registerables);

// #region API Integration
const API_BASE_URL = "http://localhost:9090";

const getAuthToken = () => {
  return localStorage.getItem("authToken");
};

// Types based on the API Swagger definition
interface ApiSimulationResults {
  average_speed: number;
  average_travel_time: number;
  average_waiting_time: number;
  total_vehicles: number;
  total_travel_time: number;
  // Other fields from model.SimulationResults can be added here
}

type Position = { time: number; x: number; y: number; speed: number };
type Vehicle = { id: string; positions: Position[] };

// This type now matches the 'output' part of the API's SimulationResponse
type SimulationOutput = {
  name?: string;
  description?: string;
  intersections?: string[];
  vehicles: Vehicle[];
  intersection?: {
    trafficLights?: {
      phases: { duration?: number }[];
    }[];
  };
};
// #endregion

// #region Helper Functions (Unchanged)
function computeStats(vehicles: Vehicle[]) {
  let totalSpeed = 0,
    maxSpeed = -Infinity,
    minSpeed = Infinity,
    speedCount = 0;
  let totalDistance = 0;
  const finalSpeeds: number[] = [];

  vehicles.forEach((veh: Vehicle) => {
    let prevPos: Position | null = null;
    veh.positions.forEach((pos: Position) => {
      totalSpeed += pos.speed;
      speedCount++;
      if (pos.speed > maxSpeed) maxSpeed = pos.speed;
      if (pos.speed < minSpeed) minSpeed = pos.speed;

      if (prevPos) {
        const dx = pos.x - prevPos.x;
        const dy = pos.y - prevPos.y;
        totalDistance += Math.sqrt(dx * dx + dy * dy);
      }
      prevPos = pos;
    });

    if (veh.positions.length > 0) {
      finalSpeeds.push(veh.positions[veh.positions.length - 1].speed);
    }
  });

  return {
    avgSpeed: speedCount ? totalSpeed / speedCount : 0,
    maxSpeed: maxSpeed === -Infinity ? 0 : maxSpeed,
    minSpeed: minSpeed === Infinity ? 0 : minSpeed,
    totalDistance,
    vehicleCount: vehicles.length,
    finalSpeeds,
  };
}

function getAllTimes(vehicles: Vehicle[]): number[] {
  const timeSet = new Set<number>();
  vehicles.forEach((veh) =>
    veh.positions.forEach((pos) => timeSet.add(pos.time)),
  );
  return Array.from(timeSet).sort((a, b) => a - b);
}

function getAverageSpeedOverTime(
  vehicles: Vehicle[],
  allTimes: number[],
): number[] {
  return allTimes.map((t) => {
    let sum = 0,
      count = 0;
    vehicles.forEach((veh) => {
      const pos = veh.positions.find((p) => p.time === t);
      if (pos) {
        sum += pos.speed;
        count++;
      }
    });
    return count ? sum / count : 0;
  });
}

function getVehicleCountOverTime(
  vehicles: Vehicle[],
  allTimes: number[],
): number[] {
  return allTimes.map((t) => {
    let count = 0;
    vehicles.forEach((veh) => {
      if (veh.positions.some((p) => p.time === t)) {
        count++;
      }
    });
    return count;
  });
}

function getTotalDistancePerVehicle(vehicles: Vehicle[]): number[] {
  return vehicles.map((veh) => {
    let dist = 0;
    let prevPos: Position | null = null;
    veh.positions.forEach((pos) => {
      if (prevPos) {
        const dx = pos.x - prevPos.x;
        const dy = pos.y - prevPos.y;
        dist += Math.sqrt(dx * dx + dy * dy);
      }
      prevPos = pos;
    });
    return dist;
  });
}

function getHistogramData(
  data: number[],
  binSize: number,
  maxVal: number,
): { counts: number[]; labels: string[] } {
  const bins = Math.ceil(maxVal / binSize);
  const counts: number[] = Array(bins).fill(0);
  const labels: string[] = [];

  data.forEach((v: number) => {
    if (v >= 0) {
      const idx = Math.min(Math.floor(v / binSize), bins - 1);
      counts[idx]++;
    }
  });

  for (let i = 0; i < bins; i++) {
    labels.push(
      `${(i * binSize).toFixed(0)}-${((i + 1) * binSize).toFixed(0)}`,
    );
  }

  return { counts, labels };
}

function downsampleData<TLabel, TData>(
  labels: TLabel[],
  data: TData[],
  maxPoints: number,
): { downsampledLabels: TLabel[]; downsampledData: TData[] } {
  if (labels.length <= maxPoints) {
    return { downsampledLabels: labels, downsampledData: data };
  }

  const downsampledLabels: TLabel[] = [];
  const downsampledData: TData[] = [];
  const step = Math.ceil(labels.length / maxPoints);

  for (let i = 0; i < labels.length; i += step) {
    downsampledLabels.push(labels[i]);
    downsampledData.push(data[i]);
  }

  return { downsampledLabels, downsampledData };
}
// #endregion

const SimulationResults: React.FC = () => {
  const [simData, setSimData] = useState<SimulationOutput | null>(null);
  const [optimizedData, setOptimizedData] = useState<SimulationOutput | null>(
    null,
  );
  const [apiResults, setApiResults] = useState<ApiSimulationResults | null>(
    null,
  );
  const [optimizedApiResults, setOptimizedApiResults] =
    useState<ApiSimulationResults | null>(null);

  const [showOptimized, setShowOptimized] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [canBeOptimized, setCanBeOptimized] = useState(false);

  const location = useLocation();
  const simInfo = location.state || {};

  // Extract details passed from the Simulations page
  const intersectionId = simInfo.intersectionIds?.[0];
  const simName = simInfo.name || "Simulation";
  const simDesc = simInfo.description || "";
  const simIntersections = simInfo.intersections || [];

  const chartInstances = useRef<Chart[]>([]);
  const chartRefs = {
    avgSpeedRef: useRef<HTMLCanvasElement | null>(null),
    vehCountRef: useRef<HTMLCanvasElement | null>(null),
    finalSpeedHistRef: useRef<HTMLCanvasElement | null>(null),
    totalDistHistRef: useRef<HTMLCanvasElement | null>(null),
  };

  // Fetch data from API on component mount
  useEffect(() => {
    if (!intersectionId) {
      setError("No intersection ID provided to display results.");
      setLoading(false);
      return;
    }

    const fetchData = async () => {
      try {
        setLoading(true);
        const authToken = getAuthToken();
        const headers = { Authorization: `Bearer ${authToken}` };

        // Fetch standard simulation data
        const simRes = await fetch(
          `${API_BASE_URL}/intersections/${intersectionId}/simulate`,
          { headers },
        );
        if (!simRes.ok) {
          throw new Error(
            `Failed to fetch simulation data: ${simRes.statusText}`,
          );
        }
        const simResponseData = await simRes.json();
        setSimData(simResponseData.output);
        setApiResults(simResponseData.results);

        // Check for and fetch optimized data
        const optRes = await fetch(
          `${API_BASE_URL}/intersections/${intersectionId}/optimise`,
          { headers },
        );
        if (optRes.ok) {
          const optResponseData = await optRes.json();
          setOptimizedData(optResponseData.output);
          setOptimizedApiResults(optResponseData.results);
          setCanBeOptimized(true); // Optimized data is available
          // If the user came here from running an optimization, show it by default
          if (simInfo.type === "optimizations") {
            setShowOptimized(true);
          }
        } else {
          setCanBeOptimized(false);
        }
      } catch (err: any) {
        setError(err.message || "Failed to load data from the API.");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [intersectionId, simInfo.type]);

  const handleOptimize = () => {
    if (canBeOptimized) {
      setShowOptimized(!showOptimized);
    }
  };

  // Render charts when data changes
  useEffect(() => {
    chartInstances.current.forEach((c) => c?.destroy());
    chartInstances.current = [];

    if (!simData || !simData.vehicles) return;

    const stats = computeStats(simData.vehicles);
    const optStats =
      showOptimized && optimizedData && optimizedData.vehicles
        ? computeStats(optimizedData.vehicles)
        : null;

    const allTimes = getAllTimes(simData.vehicles);
    const avgSpeedOverTime = getAverageSpeedOverTime(
      simData.vehicles,
      allTimes,
    );
    const vehCountOverTime = getVehicleCountOverTime(
      simData.vehicles,
      allTimes,
    );
    const totalDistPerVeh = getTotalDistancePerVehicle(simData.vehicles);

    const optAvgSpeedOverTime =
      showOptimized && optimizedData && optimizedData.vehicles
        ? getAverageSpeedOverTime(optimizedData.vehicles, allTimes)
        : [];
    const optVehCountOverTime =
      showOptimized && optimizedData && optimizedData.vehicles
        ? getVehicleCountOverTime(optimizedData.vehicles, allTimes)
        : [];
    const optTotalDistPerVeh =
      showOptimized && optimizedData && optimizedData.vehicles
        ? getTotalDistancePerVehicle(optimizedData.vehicles)
        : [];

    const MAX_TIME_POINTS = 100;
    const {
      downsampledLabels: timeLabels,
      downsampledData: downsampledAvgSpeed,
    } = downsampleData(allTimes, avgSpeedOverTime, MAX_TIME_POINTS);
    const { downsampledData: downsampledVehCount } = downsampleData(
      allTimes,
      vehCountOverTime,
      MAX_TIME_POINTS,
    );
    const { downsampledData: downsampledOptAvgSpeed } = downsampleData(
      allTimes,
      optAvgSpeedOverTime,
      MAX_TIME_POINTS,
    );
    const { downsampledData: downsampledOptVehCount } = downsampleData(
      allTimes,
      optVehCountOverTime,
      MAX_TIME_POINTS,
    );

    const { counts: finalSpeedHist, labels: finalSpeedHistLabels } =
      getHistogramData(stats.finalSpeeds, 2, 40);
    const maxDist =
      totalDistPerVeh.length > 0 ? Math.max(...totalDistPerVeh) : 0;
    const { counts: totalDistHist, labels: totalDistHistLabels } =
      getHistogramData(totalDistPerVeh, 50, Math.ceil(maxDist / 50) * 50);

    const optFinalSpeedHist = optStats
      ? getHistogramData(optStats.finalSpeeds, 2, 40)
      : null;
    const optTotalDistHist =
      showOptimized && optTotalDistPerVeh.length > 0
        ? getHistogramData(
            optTotalDistPerVeh,
            50,
            Math.ceil(Math.max(...optTotalDistPerVeh) / 50) * 50,
          )
        : null;

    const baseOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: showOptimized,
          labels: {
            color: "#fff",
            font: { size: 12 },
          },
        },
        tooltip: {
          backgroundColor: "rgba(0,0,0,0.8)",
          titleFont: { size: 16 },
          bodyFont: { size: 14 },
          padding: 12,
          cornerRadius: 4,
        },
      },
      scales: {
        x: {
          grid: { color: "rgba(255,255,255,0.1)" },
          ticks: { color: "#ccc", maxTicksLimit: 10 },
          title: { display: true, color: "#fff", font: { size: 14 } },
        },
        y: {
          grid: { color: "rgba(255,255,255,0.1)" },
          ticks: { color: "#ccc" },
          title: { display: true, color: "#fff", font: { size: 14 } },
          beginAtZero: true,
        },
      },
    };

    const createChart = (
      ref: React.RefObject<HTMLCanvasElement | null>,
      config: ChartConfiguration,
    ) => {
      if (ref.current) {
        const chart = new Chart(ref.current, config);
        chartInstances.current.push(chart);
      }
    };

    const avgSpeedDatasets = [
      {
        label: "Original Average Speed",
        data: downsampledAvgSpeed,
        borderColor: "#0F5BA7",
        backgroundColor: "#60A5FA33",
        fill: true,
        tension: 0.3,
      },
    ];

    const vehCountDatasets = [
      {
        label: "Original Vehicle Count",
        data: downsampledVehCount,
        borderColor: "#0F5BA7",
        backgroundColor: "#60A5FA33",
        fill: true,
        tension: 0.3,
      },
    ];

    const finalSpeedDatasets = [
      {
        label: "Original Final Speed Distribution",
        data: finalSpeedHist,
        backgroundColor: "#0F5BA7",
      },
    ];

    const totalDistDatasets = [
      {
        label: "Original Total Distance Distribution",
        data: totalDistHist,
        backgroundColor: "#0F5BA7",
      },
    ];

    if (showOptimized && downsampledOptAvgSpeed.length > 0) {
      avgSpeedDatasets.push({
        label: "Optimized Average Speed",
        data: downsampledOptAvgSpeed,
        borderColor: "#2B9348",
        backgroundColor: "#48ac4d33",
        fill: true,
        tension: 0.3,
      });
    }

    if (showOptimized && downsampledOptVehCount.length > 0) {
      vehCountDatasets.push({
        label: "Optimized Vehicle Count",
        data: downsampledOptVehCount,
        borderColor: "#2B9348",
        backgroundColor: "#48ac4d33",
        fill: true,
        tension: 0.3,
      });
    }

    if (showOptimized && optFinalSpeedHist) {
      finalSpeedDatasets.push({
        label: "Optimized Final Speed Distribution",
        data: optFinalSpeedHist.counts,
        backgroundColor: "#2B9348",
      });
    }

    if (showOptimized && optTotalDistHist) {
      totalDistDatasets.push({
        label: "Optimized Total Distance Distribution",
        data: optTotalDistHist.counts,
        backgroundColor: "#2B9348",
      });
    }

    createChart(chartRefs.avgSpeedRef, {
      type: "line",
      data: { labels: timeLabels, datasets: avgSpeedDatasets },
      options: {
        ...baseOptions,
        plugins: {
          ...baseOptions.plugins,
          title: {
            display: true,
            text: "Average Speed Over Time",
            color: "#fff",
            font: { size: 18 },
          },
        },
        scales: {
          ...baseOptions.scales,
          x: {
            ...baseOptions.scales.x,
            title: { ...baseOptions.scales.x.title, text: "Time (s)" },
          },
          y: {
            ...baseOptions.scales.y,
            title: { ...baseOptions.scales.y.title, text: "Speed (m/s)" },
          },
        },
      },
    });

    createChart(chartRefs.vehCountRef, {
      type: "line",
      data: { labels: timeLabels, datasets: vehCountDatasets },
      options: {
        ...baseOptions,
        plugins: {
          ...baseOptions.plugins,
          title: {
            display: true,
            text: "Vehicle Count Over Time",
            color: "#fff",
            font: { size: 18 },
          },
        },
        scales: {
          ...baseOptions.scales,
          x: {
            ...baseOptions.scales.x,
            title: { ...baseOptions.scales.x.title, text: "Time (s)" },
          },
          y: {
            ...baseOptions.scales.y,
            title: { ...baseOptions.scales.y.title, text: "Count" },
          },
        },
      },
    });

    createChart(chartRefs.finalSpeedHistRef, {
      type: "bar",
      data: { labels: finalSpeedHistLabels, datasets: finalSpeedDatasets },
      options: {
        ...baseOptions,
        plugins: {
          ...baseOptions.plugins,
          title: {
            display: true,
            text: "Histogram of Final Speeds",
            color: "#fff",
            font: { size: 18 },
          },
        },
        scales: {
          ...baseOptions.scales,
          x: {
            ...baseOptions.scales.x,
            title: { ...baseOptions.scales.x.title, text: "Speed (m/s)" },
          },
          y: {
            ...baseOptions.scales.y,
            title: {
              ...baseOptions.scales.y.title,
              text: "Number of Vehicles",
            },
          },
        },
      },
    });

    createChart(chartRefs.totalDistHistRef, {
      type: "bar",
      data: { labels: totalDistHistLabels, datasets: totalDistDatasets },
      options: {
        ...baseOptions,
        plugins: {
          ...baseOptions.plugins,
          title: {
            display: true,
            text: "Histogram of Total Distance",
            color: "#fff",
            font: { size: 18 },
          },
        },
        scales: {
          ...baseOptions.scales,
          x: {
            ...baseOptions.scales.x,
            title: { ...baseOptions.scales.x.title, text: "Distance (m)" },
          },
          y: {
            ...baseOptions.scales.y,
            title: {
              ...baseOptions.scales.y.title,
              text: "Number of Vehicles",
            },
          },
        },
      },
    });

    return () => {
      chartInstances.current.forEach((c) => c?.destroy());
      chartInstances.current = [];
    };
  }, [simData, showOptimized, optimizedData]);

  if (loading)
    return (
      <div className="text-center text-gray-700 dark:text-gray-300 py-10">
        Loading simulation data from API...
      </div>
    );
  if (error)
    return <div className="text-center text-red-500 py-10">{error}</div>;
  if (!simData)
    return (
      <div className="text-center text-gray-700 dark:text-gray-300 py-10">
        No simulation data found.
      </div>
    );

  // Use a mix of API results (for accuracy) and locally computed stats (for graphs/details)
  const stats = computeStats(simData.vehicles);
  const optStats =
    showOptimized && optimizedData && optimizedData.vehicles
      ? computeStats(optimizedData.vehicles)
      : null;

  const { numPhases, totalCycle } = simData.intersection?.trafficLights?.[0]
    ? {
        numPhases: simData.intersection.trafficLights[0].phases.length,
        totalCycle: simData.intersection.trafficLights[0].phases.reduce(
          (sum: number, p: { duration?: number }) => sum + (p.duration ?? 0),
          0,
        ),
      }
    : { numPhases: 0, totalCycle: 0 };

  const handleViewComparison = () => {
    window.location.href = "/comparison-rendering";
  };

  return (
    <div className="simulation-results-page bg-gradient-to-br from-gray-900 via-gray-800 to-black text-gray-100 min-h-screen">
      <Navbar />
      <div className="simRes-main-content py-8 px-6 overflow-y-auto">
        <div className="results max-w-full mx-auto">
          {/* Simulation meta info */}
          <div className="mb-10 flex flex-col lg:flex-row lg:justify-between lg:items-start gap-6">
            <div className="flex-1 text-left">
              <h1 className="simName text-4xl font-extrabold bg-gradient-to-r from-teal-400 to-emerald-500 bg-clip-text text-transparent mb-2 text-left">
                {simName}
              </h1>
              {simDesc && (
                <p className="simDesc text-lg text-gray-400 mb-4 leading-relaxed text-left">
                  {simDesc}
                </p>
              )}
              {simIntersections && simIntersections.length > 0 && (
                <div className="flex flex-wrap gap-2 mb-2 justify-start">
                  {simIntersections.map((intersection: string, idx: number) => (
                    <span
                      key={idx}
                      className="px-4 py-2 bg-white/10 dark:bg-[#161B22] backdrop-blur-md rounded-full text-sm font-medium text-[#0F5BA7] border-2 border-[#0F5BA7] hover:bg-white/20 transition-all duration-300"
                    >
                      {intersection}
                    </span>
                  ))}
                </div>
              )}
            </div>

            <div className="flex flex-col gap-3 lg:min-w-[280px]">
              <button
                onClick={handleViewComparison}
                className="px-8 py-3 text-base font-bold text-white bg-[#0F5BA7] border-2 border-[#0F5BA7] rounded-xl transform transition-all duration-300 ease-in-out hover:scale-105 focus:outline-none focus:ring-4 focus:ring-[#0F5BA7]/50 hover:shadow-xl hover:shadow-[#0F5BA7]/40"
              >
                View Rendering
              </button>
              <button
                onClick={handleOptimize}
                disabled={!canBeOptimized}
                className={`px-8 py-3 text-base font-bold text-white rounded-xl shadow-lg transform transition-all duration-300 ease-in-out focus:outline-none focus:ring-4 ${
                  !canBeOptimized
                    ? "bg-gray-600 cursor-not-allowed"
                    : "bg-gradient-to-r from-green-600 to-green-700 shadow-green-500/50 hover:scale-105 hover:shadow-xl hover:shadow-green-500/60 focus:ring-green-300"
                }`}
              >
                {!canBeOptimized
                  ? "Not Optimized"
                  : showOptimized
                    ? "Hide Optimization"
                    : "Show Optimization"}
              </button>
            </div>
          </div>

          {/* Simulation Results Section */}
          <section className="visSection simulation-section bg-white/5 backdrop-blur-md p-8 rounded-xl shadow-lg border border-gray-800/50 w-full text-center">
            <h2 className="text-2xl font-semibold mb-8">
              {showOptimized ? (
                <>
                  <span className="bg-[#0F5BA7] bg-clip-text text-transparent">
                    Simulation Results
                  </span>
                  <span className="text-gray-400"> vs </span>
                  <span className="bg-[#2B9348] bg-clip-text text-transparent">
                    Optimized Results
                  </span>
                </>
              ) : (
                <span className="bg-[#0F5BA7] bg-clip-text text-transparent">
                  Simulation Results
                </span>
              )}
            </h2>

            <div className="statGrid grid grid-cols-2 md:grid-cols-4 xl:grid-cols-7 gap-2 mb-8 justify-items-center">
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[180px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  Average Speed
                </div>
                <div className="text-xl font-bold text-[#0F5BA7]">
                  {apiResults ? apiResults.average_speed.toFixed(2) : "..."}{" "}
                  <span className="text-sm text-[#0F5BA7] font-normal">
                    m/s
                  </span>
                </div>
                {showOptimized && optimizedApiResults && (
                  <div className="text-lg font-semibold text-[#2B9348] mt-1">
                    {optimizedApiResults.average_speed.toFixed(2)}{" "}
                    <span className="text-xs text-[#2B9348] font-normal">
                      m/s
                    </span>
                  </div>
                )}
              </div>
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[180px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  Max Speed
                </div>
                <div className="text-xl font-bold text-[#0F5BA7]">
                  {stats ? stats.maxSpeed.toFixed(2) : "..."}{" "}
                  <span className="text-sm text-[#0F5BA7] font-normal">
                    m/s
                  </span>
                </div>
                {showOptimized && optStats && (
                  <div className="text-lg font-semibold text-[#2B9348] mt-1">
                    {optStats.maxSpeed.toFixed(2)}{" "}
                    <span className="text-xs text-[#2B9348] font-normal">
                      m/s
                    </span>
                  </div>
                )}
              </div>
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[180px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  Min Speed
                </div>
                <div className="text-xl font-bold text-[#0F5BA7]">
                  {stats ? stats.minSpeed.toFixed(2) : "..."}{" "}
                  <span className="text-sm text-[#0F5BA7] font-normal">
                    m/s
                  </span>
                </div>
                {showOptimized && optStats && (
                  <div className="text-lg font-semibold text-[#2B9348] mt-1">
                    {optStats.minSpeed.toFixed(2)}{" "}
                    <span className="text-xs text-[#2B9348] font-normal">
                      m/s
                    </span>
                  </div>
                )}
              </div>
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[180px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  Total Distance
                </div>
                <div className="text-xl font-bold text-[#0F5BA7]">
                  {stats ? stats.totalDistance.toFixed(2) : "..."}{" "}
                  <span className="text-sm text-[#0F5BA7] font-normal">m</span>
                </div>
                {showOptimized && optStats && (
                  <div className="text-lg font-semibold text-[#2B9348] mt-1">
                    {optStats.totalDistance.toFixed(2)}{" "}
                    <span className="text-xs text-[#2B9348] font-normal">
                      m
                    </span>
                  </div>
                )}
              </div>
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[180px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  # Vehicles
                </div>
                <div className="text-xl font-bold text-[#0F5BA7]">
                  {apiResults ? apiResults.total_vehicles : "..."}
                </div>
                {showOptimized && optimizedApiResults && (
                  <div className="text-lg font-semibold text-[#2B9348] mt-1">
                    {optimizedApiResults.total_vehicles}
                  </div>
                )}
              </div>
              {/* Traffic Light Stat Cards */}
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-yellow-400/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[180px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  # TL Phases
                </div>
                <div className="text-xl font-bold text-[#0F5BA7]">
                  {numPhases}
                </div>
              </div>
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-yellow-400/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[180px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  TL Cycle Duration
                </div>
                <div className="text-xl font-bold text-[#0F5BA7]">
                  {totalCycle}{" "}
                  <span className="text-sm text-[#0F5BA7] font-normal">s</span>
                </div>
              </div>
            </div>

            {/* Graphs grid */}
            <div className="graphGrid grid grid-cols-1 lg:grid-cols-2 gap-8">
              <div className="bg-white dark:bg-[#161B22] outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-2xl p-6 h-80 w-full flex items-center justify-center">
                <canvas ref={chartRefs.avgSpeedRef} className="w-full h-full" />
              </div>
              <div className="bg-white dark:bg-[#161B22] outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-2xl p-6 h-80 w-full flex items-center justify-center">
                <canvas ref={chartRefs.vehCountRef} className="w-full h-full" />
              </div>
              <div className="bg-white dark:bg-[#161B22] outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-2xl p-6 h-80 w-full flex items-center justify-center">
                <canvas
                  ref={chartRefs.finalSpeedHistRef}
                  className="w-full h-full"
                />
              </div>
              <div className="bg-white dark:bg-[#161B22] outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-2xl p-6 h-80 w-full flex items-center justify-center">
                <canvas
                  ref={chartRefs.totalDistHistRef}
                  className="w-full h-full"
                />
              </div>
            </div>
          </section>
        </div>
      </div>
      <Footer />
      <HelpMenu />
    </div>
  );
};

export default SimulationResults;