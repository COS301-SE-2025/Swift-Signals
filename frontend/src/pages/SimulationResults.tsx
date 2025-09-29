import { Chart, registerables } from "chart.js";
import type { ChartConfiguration } from "chart.js";
import React, { useEffect, useState, useRef } from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom";

import Footer from "../components/Footer";
import HelpMenu from "../components/HelpMenu";
import Navbar from "../components/Navbar";
import { API_BASE_URL } from "../config";
import "../styles/SimulationResults.css";

if (Chart.register) {
  Chart.register(...registerables);
}

// #region API Integration
// const API_BASE_URL = "http://localhost:9090";

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
}

type Position = { time: number; x: number; y: number; speed: number };
type Vehicle = { id: string; positions: Position[] };

// Matches the 'output' part of the API's SimulationResponse
type SimulationOutput = {
  vehicles: Vehicle[];
  intersection?: {
    trafficLights?: {
      phases: { duration?: number }[];
    }[];
  };
};

// Full API response for simulation/optimisation endpoints
interface ApiSimulationResponse {
  output: SimulationOutput;
  results: ApiSimulationResults;
}

// Full API response for a single intersection
interface ApiIntersection {
  name: string;
  traffic_density: string;
  status?: string; // Added status property
  // Add other fields from the full intersection object if needed
}
// #endregion

// #region Loading Component
const LoadingAnimation: React.FC = () => (
  <div className="flex flex-col items-center justify-center min-h-[400px] space-y-6">
    <div className="relative">
      {/* Outer spinning ring */}
      <div className="w-20 h-20 border-4 border-gray-600 border-t-teal-500 rounded-full animate-spin"></div>
      {/* Inner pulsing dot */}
      <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2">
        <div className="w-4 h-4 bg-gradient-to-r from-teal-400 to-emerald-500 rounded-full animate-pulse"></div>
      </div>
    </div>

    {/* Loading text with animated dots */}
    <div className="text-center">
      <div className="text-xl font-semibold text-gray-300 mb-2">
        Loading Simulation Data
      </div>
      <div className="flex items-center justify-center space-x-1">
        <div
          className="w-2 h-2 bg-teal-400 rounded-full animate-bounce"
          style={{ animationDelay: "0ms" }}
        ></div>
        <div
          className="w-2 h-2 bg-teal-400 rounded-full animate-bounce"
          style={{ animationDelay: "150ms" }}
        ></div>
        <div
          className="w-2 h-2 bg-teal-400 rounded-full animate-bounce"
          style={{ animationDelay: "300ms" }}
        ></div>
      </div>
    </div>

    {/* Progress bar */}
    <div className="w-64 h-2 bg-gray-700 rounded-full overflow-hidden">
      <div className="h-full bg-gradient-to-r from-teal-400 to-emerald-500 rounded-full animate-pulse"></div>
    </div>
  </div>
);
// #endregion

// #region Error Component
const ErrorDisplay: React.FC<{ error: string; onRetry: () => void }> = ({
  error,
  onRetry,
}) => (
  <div className="flex flex-col items-center justify-center min-h-[400px] space-y-6 text-center">
    {/* Error icon */}
    <div className="w-16 h-16 bg-red-500/20 rounded-full flex items-center justify-center">
      <svg
        className="w-8 h-8 text-red-500"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
    </div>

    {/* Error message */}
    <div className="max-w-md">
      <h3 className="text-xl font-semibold text-red-400 mb-2">
        Failed to Load Data
      </h3>
      <p className="text-gray-400 mb-6">{error}</p>
    </div>

    {/* Retry button */}
    <button
      onClick={onRetry}
      className="px-8 py-3 bg-gradient-to-r from-red-600 to-red-700 text-white font-semibold rounded-xl 
                 hover:from-red-700 hover:to-red-800 transform transition-all duration-300 ease-in-out 
                 hover:scale-105 focus:outline-none focus:ring-4 focus:ring-red-500/50 
                 hover:shadow-xl hover:shadow-red-500/40 flex items-center space-x-2"
    >
      <svg
        className="w-5 h-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
        />
      </svg>
      <span>Retry</span>
    </button>
  </div>
);
// #endregion

// #region Helper Functions
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

// Helper function to extract street name from a string that might contain coordinates
const getStreetName = (fullName: string | undefined | null): string => {
  if (!fullName) return "Simulation Results";
  // Remove 'Simulation Results for ' prefix
  let cleanedName = fullName.replace(/^Simulation Results for\s*/, "");
  // Remove coordinates in square brackets, e.g., ' [-25.757139,28.1936006]'
  cleanedName = cleanedName.replace(/\s*\[[^\]]*\]$/, "");
  return cleanedName.trim();
};

// Helper function to clean the description string
const cleanDescription = (
  desc: string | undefined | null,
): string | undefined => {
  if (!desc) return undefined;
  // Only remove coordinates in square brackets, e.g., ' [-25.757139,28.1936006]'
  const cleanedDesc = desc.replace(/\s*\[[^\]]*\]$/, "");
  return cleanedDesc.trim();
};

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
  const [intersectionData, setIntersectionData] =
    useState<ApiIntersection | null>(null);

  const [showOptimized, setShowOptimized] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isOptimized, setIsOptimized] = useState(false);
  const [isOptimizing, setIsOptimizing] = useState(false);
  const [optimizationStatus, setOptimizationStatus] = useState<string>("");

  const location = useLocation();
  const navigate = useNavigate();
  const params = useParams();
  const { intersectionIds, name, description, type } = location.state || {};

  console.log("name from location.state:", name);
  console.log("intersectionData?.name:", intersectionData?.name);

  // Get intersectionId from URL params first, then fall back to location.state
  const intersectionId = params.intersectionId || intersectionIds?.[0];

  const [isDarkMode, setIsDarkMode] = useState(
    document.documentElement.classList.contains("dark"),
  );

  useEffect(() => {
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.attributeName === "class") {
          setIsDarkMode(document.documentElement.classList.contains("dark"));
        }
      });
    });

    observer.observe(document.documentElement, { attributes: true });

    return () => observer.disconnect();
  }, []);

  const chartInstances = useRef<Chart[]>([]);
  const chartRefs = {
    avgSpeedRef: useRef<HTMLCanvasElement | null>(null),
    vehCountRef: useRef<HTMLCanvasElement | null>(null),
    finalSpeedHistRef: useRef<HTMLCanvasElement | null>(null),
    totalDistHistRef: useRef<HTMLCanvasElement | null>(null),
  };

  // Function to run optimization
  const runOptimization = async () => {
    if (!intersectionId) {
      alert("No intersection ID available for optimization.");
      return;
    }

    setIsOptimizing(true);
    setOptimizationStatus("Running optimization...");

    // Save optimization status to localStorage
    localStorage.setItem(
      `optimizationStatus_${intersectionId}`,
      JSON.stringify({ status: "optimising" }),
    );

    try {
      const authToken = getAuthToken();
      if (!authToken) {
        throw new Error("Authentication token not found. Please log in again.");
      }

      // Step 1: Run optimization
      const optResponse = await fetch(
        `${API_BASE_URL}/intersections/${intersectionId}/optimise`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${authToken}`,
            "Content-Type": "application/json",
          },
        },
      );

      if (!optResponse.ok) {
        if (optResponse.status === 401) {
          throw new Error("Authentication failed. Please log in again.");
        } else if (optResponse.status === 404) {
          throw new Error("Intersection not found for optimization.");
        } else {
          throw new Error(
            `Failed to run optimization: ${optResponse.statusText}`,
          );
        }
      }

      const optResult = await optResponse.json();
      console.log("Optimization result:", optResult);

      setOptimizationStatus(
        "Optimization completed successfully! Fetching optimized data...",
      );

      // Step 2: Fetch optimized simulation data
      const optDataResponse = await fetch(
        `${API_BASE_URL}/intersections/${intersectionId}/optimise`,
        {
          headers: { Authorization: `Bearer ${authToken}` },
        },
      );

      if (!optDataResponse.ok) {
        throw new Error(
          `Failed to fetch optimized data: ${optDataResponse.statusText}`,
        );
      }

      const optData: ApiSimulationResponse = await optDataResponse.json();

      if (!optData.output) {
        throw new Error(
          "Invalid optimized simulation data received from server",
        );
      }

      // Process traffic lights for optimized data if they exist
      if (
        optData.output.intersection &&
        optData.output.intersection.trafficLights
      ) {
        const processedTrafficLights = processTrafficLights(
          optData.output.intersection.trafficLights,
        );
        const newOptData = {
          ...optData.output,
          intersection: {
            ...optData.output.intersection,
            trafficLights: processedTrafficLights,
          },
        };
        setOptimizedData(newOptData);
      } else {
        setOptimizedData(optData.output);
      }

      setOptimizedApiResults(optData.results);
      setIsOptimized(true);
      setShowOptimized(true);
      setOptimizationStatus("Optimization completed successfully!");

      // Save optimized data to localStorage
      localStorage.setItem(
        `optimizationStatus_${intersectionId}`,
        JSON.stringify({
          status: "optimised",
          data: optData,
        }),
      );

      // Update the intersection status to "optimised" in the backend
      await updateIntersectionStatus(
        intersectionId,
        "INTERSECTION_STATUS_OPTIMISED",
      );

      setTimeout(() => {
        setOptimizationStatus("");
      }, 3000);
    } catch (error) {
      console.error("Error running optimization:", error);
      setOptimizationStatus(
        `Optimization failed: ${error instanceof Error ? error.message : "Unknown error"}`,
      );
      setTimeout(() => {
        setOptimizationStatus("");
      }, 5000);
    } finally {
      setIsOptimizing(false);
    }
  };

  // Helper function to process traffic lights (extracted from existing code)
  const processTrafficLights = (
    trafficLights: { phases?: { duration?: number }[] }[],
  ) => {
    // This would need to be updated based on your actual connection structure
    // For now, using a simplified approach

    const maxSignalIndex = 11; // Default value, adjust as needed
    const stateArrayLength = maxSignalIndex >= 0 ? maxSignalIndex + 1 : 12;

    const newPhases = [
      { duration: 30, state: "G".repeat(stateArrayLength) }, // North-South green
      { duration: 5, state: "y".repeat(stateArrayLength) }, // North-South yellow
      { duration: 30, state: "G".repeat(stateArrayLength) }, // East-West green
      { duration: 5, state: "y".repeat(stateArrayLength) }, // East-West yellow
    ];

    return trafficLights.map((light) => {
      let time = 0;
      const newStates = newPhases.map((phase) => {
        const state = { time: time, state: phase.state };
        time += phase.duration;
        return state;
      });
      newStates.push({ time: time, state: newPhases[0].state });
      return { ...light, phases: newPhases, states: newStates };
    });
  };

  // Function to update intersection status
  const updateIntersectionStatus = async (id: string, status: string) => {
    try {
      const authToken = getAuthToken();
      if (!authToken) return;

      await fetch(`${API_BASE_URL}/intersections/${id}`, {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${authToken}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ status }),
      });
    } catch (error) {
      console.error("Failed to update intersection status:", error);
    }
  };

  const fetchData = async () => {
    if (!intersectionId) {
      setError("No intersection ID provided.");
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const authToken = getAuthToken();
      const headers = { Authorization: `Bearer ${authToken}` };

      // Fetch simulation and intersection data in parallel
      const [simRes, intersectionRes] = await Promise.all([
        fetch(`${API_BASE_URL}/intersections/${intersectionId}/simulate`, {
          headers,
        }),
        fetch(`${API_BASE_URL}/intersections/${intersectionId}`, { headers }),
      ]);

      // Process simulation data
      if (!simRes.ok)
        throw new Error(
          `Failed to fetch simulation data: ${simRes.statusText}`,
        );
      const simResponseData: ApiSimulationResponse = await simRes.json();
      setSimData(simResponseData.output);
      setApiResults(simResponseData.results);

      // Process intersection details
      if (!intersectionRes.ok)
        throw new Error(
          `Failed to fetch intersection details: ${intersectionRes.statusText}`,
        );
      const intersectionResponseData: ApiIntersection =
        await intersectionRes.json();
      setIntersectionData(intersectionResponseData);

      // If intersection is optimized, fetch optimization data
      if (intersectionResponseData.status === "INTERSECTION_STATUS_OPTIMISED") {
        setIsOptimized(true);
        const optRes = await fetch(
          `${API_BASE_URL}/intersections/${intersectionId}/optimise`,
          {
            headers,
          },
        );

        if (optRes.ok) {
          const optResponseData: ApiSimulationResponse = await optRes.json();
          setOptimizedData(optResponseData.output);
          setOptimizedApiResults(optResponseData.results);
          // If the user came from the "Optimizations" table, show the comparison by default
          if (type === "optimizations") {
            setShowOptimized(true);
          }
        }
      } else {
        setIsOptimized(false);
      }
    } catch (err: unknown) {
      setError(
        err instanceof Error
          ? err.message
          : "Failed to load data from the API.",
      );
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const savedStatus = localStorage.getItem(
      `optimizationStatus_${intersectionId}`,
    );
    if (savedStatus) {
      const { status, data } = JSON.parse(savedStatus);
      if (status === "optimised") {
        setOptimizedData(data.output);
        setOptimizedApiResults(data.results);
        setIsOptimized(true);
        setShowOptimized(true);
      } else if (status === "optimising") {
        setIsOptimizing(true);
        // You might want to add polling logic here to check for completion
      }
    }
  }, [intersectionId]);

  useEffect(() => {
    fetchData();
  }, [intersectionId]);

  // Chart creation and updates
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
            color: isDarkMode ? "#fff" : "#777777ff",
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
          grid: {
            color: isDarkMode ? "rgba(255,255,255,0.1)" : "rgba(0,0,0,0.1)",
          },
          title: {
            display: true,
            color: isDarkMode ? "#fff" : "#333",
            font: { size: 14 },
          },
        },
        y: {
          grid: {
            color: isDarkMode ? "rgba(255,255,255,0.1)" : "rgba(0,0,0,0.1)",
          },
          ticks: { color: isDarkMode ? "#ccc" : "#666" },
          title: {
            display: true,
            color: isDarkMode ? "#fff" : "#777777ff",
            font: { size: 14 },
          },
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
            color: isDarkMode ? "#fff" : "#777777ff",
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
            color: isDarkMode ? "#fff" : "#777777ff",
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
            color: isDarkMode ? "#fff" : "#777777ff",
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
            color: isDarkMode ? "#fff" : "#777777ff",
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
  }, [simData, showOptimized, optimizedData, isDarkMode]);

  const handleViewRendering = () => {
    if (intersectionId) {
      navigate("/comparison-rendering", {
        state: {
          originalIntersectionId: intersectionId,
          originalIntersectionName: intersectionData?.name || "Simulation",
          simulationData: simData,
          optimizedData: optimizedData,
        },
      });
    } else {
      alert("No intersection ID available for rendering.");
    }
  };

  if (loading) {
    return (
      <div className="simulation-results-page bg-gradient-to-br from-gray-900 via-gray-800 to-black text-gray-100 min-h-screen">
        <Navbar />
        <div className="simRes-main-content py-8 px-6 overflow-y-auto">
          <LoadingAnimation />
        </div>
        <Footer />
        <HelpMenu />
      </div>
    );
  }

  if (error) {
    return (
      <div className="simulation-results-page bg-gradient-to-br from-gray-900 via-gray-800 to-black text-gray-100 min-h-screen">
        <Navbar />
        <div className="simRes-main-content py-8 px-6 overflow-y-auto">
          <ErrorDisplay error={error} onRetry={fetchData} />
        </div>
        <Footer />
        <HelpMenu />
      </div>
    );
  }

  if (!simData) {
    return (
      <div className="simulation-results-page bg-gradient-to-br from-gray-900 via-gray-800 to-black text-gray-100 min-h-screen">
        <Navbar />
        <div className="simRes-main-content py-8 px-6 overflow-y-auto">
          <div className="text-center text-gray-700 dark:text-gray-300 py-10">
            No simulation data found.
          </div>
        </div>
        <Footer />
        <HelpMenu />
      </div>
    );
  }

  const { numPhases, totalCycle } = simData.intersection?.trafficLights?.[0]
    ? {
        numPhases: simData.intersection.trafficLights[0].phases.length,
        totalCycle: simData.intersection.trafficLights[0].phases.reduce(
          (sum: number, p: { duration?: number }) => sum + (p.duration ?? 0),
          0,
        ),
      }
    : { numPhases: undefined, totalCycle: undefined };

  // Define trafficDensityLabel from intersectionData
  const trafficDensityLabel = intersectionData?.traffic_density
    ? intersectionData.traffic_density.replace(/_/g, " ").toLowerCase()
    : "unknown";

  const displayedName = getStreetName(name || intersectionData?.name);
  const displayedDescription = cleanDescription(description);

  return (
    <div className="simulation-results-page bg-gradient-to-br from-gray-900 via-gray-800 to-black text-gray-100 min-h-screen">
      <Navbar />
      <div className="simRes-main-content py-8 px-6 overflow-y-auto">
        <div className="results max-w-full mx-auto">
          <div className="mb-10 flex flex-col lg:flex-row lg:justify-between lg:items-start gap-6">
            <div className="flex-1 text-left">
              <h1 className="simName text-4xl font-extrabold bg-gradient-to-r from-teal-400 to-emerald-500 bg-clip-text text-transparent mb-2 text-left">
                {displayedName}
              </h1>
              {displayedDescription && (
                <p className="simDesc text-lg text-gray-400 mb-4 leading-relaxed text-left">
                  {displayedDescription}
                </p>
              )}
            </div>
            <div className="flex flex-col gap-3 lg:min-w-[280px]">
              <button
                onClick={handleViewRendering}
                className="px-8 py-3 text-base font-bold text-white bg-[#0F5BA7] border-2 border-[#0F5BA7] rounded-xl transform transition-all duration-300 ease-in-out hover:scale-105 focus:outline-none focus:ring-4 focus:ring-[#0F5BA7]/50 hover:shadow-xl hover:shadow-[#0F5BA7]/40"
              >
                View Rendering
              </button>

              {/* Optimization Button */}
              {!isOptimized && (
                <button
                  onClick={runOptimization}
                  disabled={isOptimizing}
                  className={`px-8 py-3 text-base font-bold text-white rounded-xl shadow-lg transform transition-all duration-300 ease-in-out focus:outline-none focus:ring-4 ${
                    isOptimizing
                      ? "bg-gray-600 cursor-not-allowed"
                      : "bg-gradient-to-r from-orange-600 to-orange-700 shadow-orange-500/50 hover:scale-105 hover:shadow-xl hover:shadow-orange-500/60 focus:ring-orange-300"
                  }`}
                >
                  {isOptimizing ? (
                    <div className="flex items-center justify-center space-x-2">
                      <div className="animate-spin inline-block w-4 h-4 border-2 border-current border-t-transparent rounded-full"></div>
                      <span>Optimizing...</span>
                    </div>
                  ) : (
                    "Run Optimization"
                  )}
                </button>
              )}

              {/* Show Optimization Results Button */}
              {isOptimized && (
                <button
                  onClick={() => setShowOptimized(!showOptimized)}
                  className="px-8 py-3 text-base font-bold text-white rounded-xl shadow-lg transform transition-all duration-300 ease-in-out hover:scale-105 focus:outline-none focus:ring-4 bg-gradient-to-r from-green-600 to-green-700 shadow-green-500/50 hover:shadow-xl hover:shadow-green-500/60 focus:ring-green-300"
                >
                  Show/Hide Optimization
                </button>
              )}
            </div>
          </div>

          {/* Optimization Status Message */}
          {optimizationStatus && (
            <div
              className={`mb-6 p-4 rounded-lg text-center ${
                optimizationStatus.includes("failed") ||
                optimizationStatus.includes("No improvement")
                  ? "bg-red-500/20 border border-red-500/30 text-red-300"
                  : "bg-green-600 border border-green-700 text-white dark:bg-green-700/20 dark:border-green-700/30 dark:text-green-300"
              }`}
            >
              <p className="font-semibold">{optimizationStatus}</p>
            </div>
          )}

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

            <div className="statGrid grid grid-cols-2 md:grid-cols-4 xl:grid-cols-8 gap-2 mb-8 justify-items-center">
              {/* Stat cubes */}
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[160px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  Avg Speed
                </div>
                <div className="text-xl font-bold text-[#0F5BA7]">
                  {apiResults ? apiResults.average_speed.toFixed(2) : "..."}{" "}
                  <span className="text-sm font-normal">m/s</span>
                </div>
                {showOptimized && optimizedApiResults && (
                  <div className="text-lg font-semibold text-[#2B9348] mt-1">
                    {optimizedApiResults.average_speed.toFixed(2)}{" "}
                    <span className="text-xs font-normal">m/s</span>
                  </div>
                )}
              </div>
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[160px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  Avg Travel Time
                </div>
                <div className="text-xl font-bold text-[#0F5BA7]">
                  {apiResults
                    ? apiResults.average_travel_time.toFixed(2)
                    : "..."}{" "}
                  <span className="text-sm font-normal">s</span>
                </div>
                {showOptimized && optimizedApiResults && (
                  <div className="text-lg font-semibold text-[#2B9348] mt-1">
                    {optimizedApiResults.average_travel_time.toFixed(2)}{" "}
                    <span className="text-xs font-normal">s</span>
                  </div>
                )}
              </div>
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[160px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  Avg Wait Time
                </div>
                <div className="text-xl font-bold text-[#0F5BA7]">
                  {apiResults
                    ? apiResults.average_waiting_time.toFixed(2)
                    : "..."}{" "}
                  <span className="text-sm font-normal">s</span>
                </div>
                {showOptimized && optimizedApiResults && (
                  <div className="text-lg font-semibold text-[#2B9348] mt-1">
                    {optimizedApiResults.average_waiting_time.toFixed(2)}{" "}
                    <span className="text-xs font-normal">s</span>
                  </div>
                )}
              </div>
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[160px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  Throughput
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
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-yellow-400/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[160px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  # Light Phases
                </div>
                <div className="text-xl font-bold text-[#0F5BA7]">
                  {numPhases}
                </div>
              </div>
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-yellow-400/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[160px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  Light Cycle
                </div>
                <div className="text-xl font-bold text-[#0F5BA7]">
                  {totalCycle} <span className="text-sm font-normal">s</span>
                </div>
              </div>
              <div className="stat-cube bg-white dark:bg-[#161B22] border border-purple-400/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-4 text-center shadow-md min-w-[160px]">
                <div className="text-sm font-bold text-gray-600 mb-1">
                  Traffic Density
                </div>
                <div className="text-xl font-bold text-[#0F5BA7] capitalize">
                  {trafficDensityLabel}
                </div>
              </div>
            </div>

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
export {
  computeStats,
  getAllTimes,
  getAverageSpeedOverTime,
  getVehicleCountOverTime,
  getTotalDistancePerVehicle,
  getStreetName,
  cleanDescription,
};
