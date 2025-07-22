import React, { useEffect, useState, useRef } from "react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import "../styles/SimulationResults.css";
import HelpMenu from "../components/HelpMenu";
import { Chart, registerables } from "chart.js";
import { useLocation } from "react-router-dom";

Chart.register(...registerables);

// Helper functions for stats and chart data
type Position = { time: number; x: number; y: number; speed: number };
type Vehicle = { id: string; positions: Position[] };

function computeStats(vehicles: Vehicle[]) {
  let totalSpeed = 0, maxSpeed = -Infinity, minSpeed = Infinity, speedCount = 0;
  let totalDistance = 0;
  let vehicleCount = vehicles.length;
  const finalSpeeds: number[] = [];
  vehicles.forEach((veh: Vehicle) => {
    let prev: Position | null = null;
    veh.positions.forEach((pos: Position, idx: number) => {
      totalSpeed += pos.speed;
      speedCount++;
      if (pos.speed > maxSpeed) maxSpeed = pos.speed;
      if (pos.speed < minSpeed) minSpeed = pos.speed;
      if (prev) {
        const dx = pos.x - prev.x;
        const dy = pos.y - prev.y;
        totalDistance += Math.sqrt(dx * dx + dy * dy);
      }
      prev = pos;
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
    vehicleCount,
    finalSpeeds,
  };
}

function getAllTimes(vehicles: Vehicle[]): number[] {
  // Get all unique time points across all vehicles
  const timeSet = new Set<number>();
  vehicles.forEach((veh: Vehicle) => {
    veh.positions.forEach((pos: Position) => timeSet.add(pos.time));
  });
  return Array.from(timeSet).sort((a, b) => a - b);
}

function getAverageSpeedOverTime(vehicles: Vehicle[], allTimes: number[]): number[] {
  // For each time, compute average speed of all vehicles present at that time
  return allTimes.map((t: number) => {
    let sum = 0, count = 0;
    vehicles.forEach((veh: Vehicle) => {
      const pos = veh.positions.find((p: Position) => p.time === t);
      if (pos) {
        sum += pos.speed;
        count++;
      }
    });
    return count ? sum / count : 0;
  });
}

function getVehicleCountOverTime(vehicles: Vehicle[], allTimes: number[]): number[] {
  // For each time, count how many vehicles have a position at that time
  return allTimes.map((t: number) => {
    let count = 0;
    vehicles.forEach((veh: Vehicle) => {
      if (veh.positions.find((p: Position) => p.time === t)) count++;
    });
    return count;
  });
}

function getTotalDistancePerVehicle(vehicles: Vehicle[]): number[] {
  return vehicles.map((veh: Vehicle) => {
    let dist = 0;
    let prev: Position | null = null;
    veh.positions.forEach((pos: Position) => {
      if (prev) {
        const dx = pos.x - prev.x;
        const dy = pos.y - prev.y;
        dist += Math.sqrt(dx * dx + dy * dy);
      }
      prev = pos;
    });
    return dist;
  });
}

function getVehicleIds(vehicles: Vehicle[]): string[] {
  return vehicles.map((veh: Vehicle) => veh.id);
}

function getHistogramData(data: number[], binSize = 2, maxVal = 40): number[] {
  // Simple histogram for speeds
  const bins = Math.ceil(maxVal / binSize);
  const counts: number[] = Array(bins).fill(0);
  data.forEach((v: number) => {
    const idx = Math.min(Math.floor(v / binSize), bins - 1);
    counts[idx]++;
  });
  return counts;
}

const SimulationResults: React.FC = () => {
  const [simData, setSimData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const location = useLocation();

  // Extract simulation info from location.state if available
  const simInfo = location.state || {};
  const simName = simInfo.name || simData?.name || "Simulation";
  const simDesc = simInfo.description || simData?.description || "";
  const simIntersections = simInfo.intersections || simData?.intersections || [];

  // Extract traffic light stats from simData
  let numPhases = 0;
  let totalCycle = 0;
  if (simData && simData.intersection && simData.intersection.trafficLights && simData.intersection.trafficLights.length > 0) {
    const tl = simData.intersection.trafficLights[0];
    numPhases = tl.phases.length;
    totalCycle = tl.phases.reduce((sum: number, p: any) => sum + (p.duration || 0), 0);
  }

  // Chart refs
  const simAvgSpeedRef = useRef<HTMLCanvasElement | null>(null);
  const simVehCountRef = useRef<HTMLCanvasElement | null>(null);
  const simFinalSpeedHistRef = useRef<HTMLCanvasElement | null>(null);
  const simTotalDistRef = useRef<HTMLCanvasElement | null>(null);

  const optAvgSpeedRef = useRef<HTMLCanvasElement | null>(null);
  const optVehCountRef = useRef<HTMLCanvasElement | null>(null);
  const optFinalSpeedHistRef = useRef<HTMLCanvasElement | null>(null);
  const optTotalDistRef = useRef<HTMLCanvasElement | null>(null);

  // Chart instances for cleanup
  const chartInstances = useRef<any[]>([]);

  useEffect(() => {
    fetch("/simulation_output (1).json")
      .then((res) => res.json())
      .then((data) => {
        setSimData(data);
        setLoading(false);
      })
      .catch((err) => {
        setError("Failed to load simulation data");
        setLoading(false);
      });
  }, []);

  // Compute stats and chart data
  let stats = null;
  let allTimes: number[] = [];
  let avgSpeedOverTime: number[] = [];
  let vehCountOverTime: number[] = [];
  let totalDistPerVeh: number[] = [];
  let vehIds: string[] = [];
  let finalSpeedHist: number[] = [];
  let histLabels: string[] = [];

  if (simData && simData.vehicles) {
    stats = computeStats(simData.vehicles);
    allTimes = getAllTimes(simData.vehicles);
    avgSpeedOverTime = getAverageSpeedOverTime(simData.vehicles, allTimes);
    vehCountOverTime = getVehicleCountOverTime(simData.vehicles, allTimes);
    totalDistPerVeh = getTotalDistancePerVehicle(simData.vehicles);
    vehIds = getVehicleIds(simData.vehicles);
    finalSpeedHist = getHistogramData(stats.finalSpeeds, 2, 40);
    histLabels = finalSpeedHist.map((_, i) => `${i * 2}-${i * 2 + 2}`);
  }

  // Handler for opening the rendering page
  const handleViewRendering = () => {
    window.open('/traffic-simulation', '_blank');
  };

  // Chart rendering
  useEffect(() => {
    // Cleanup previous charts
    chartInstances.current.forEach((c) => c && c.destroy());
    chartInstances.current = [];
    if (!simData || !simData.vehicles) return;
    // Chart options
    const baseOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
        tooltip: { backgroundColor: "rgba(0,0,0,0.9)", cornerRadius: 10, padding: 12, bodyFont: { size: 14, weight: 'bold' as const } },
      },
      scales: {
        x: { grid: { color: "rgba(255,255,255,0.1)" }, ticks: { color: "#ccc" }, title: { color: "#fff" } },
        y: { grid: { color: "rgba(255,255,255,0.1)" }, ticks: { color: "#ccc" }, title: { color: "#fff" } },
      },
    };
    // Simulation Results
    if (simAvgSpeedRef.current) {
      chartInstances.current.push(new Chart(simAvgSpeedRef.current, {
        type: "line",
        data: {
          labels: allTimes,
          datasets: [{ label: "Avg Speed", data: avgSpeedOverTime, borderColor: "#34D399", backgroundColor: "#34D39933", fill: true }],
        },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Average Speed Over Time", color: "#fff" } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { display: true, text: "Time (s)", color: "#fff" } }, y: { ...baseOptions.scales.y, title: { display: true, text: "Speed (m/s)", color: "#fff" } } } },
      }));
    }
    if (simVehCountRef.current) {
      chartInstances.current.push(new Chart(simVehCountRef.current, {
        type: "line",
        data: {
          labels: allTimes,
          datasets: [{ label: "Vehicle Count", data: vehCountOverTime, borderColor: "#60A5FA", backgroundColor: "#60A5FA33", fill: true }],
        },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Vehicle Count Over Time", color: "#fff" } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { display: true, text: "Time (s)", color: "#fff" } }, y: { ...baseOptions.scales.y, title: { display: true, text: "Count", color: "#fff" } } } },
      }));
    }
    if (simFinalSpeedHistRef.current) {
      chartInstances.current.push(new Chart(simFinalSpeedHistRef.current, {
        type: "bar",
        data: {
          labels: histLabels,
          datasets: [{ label: "Final Speeds", data: finalSpeedHist, backgroundColor: "#F59E42" }],
        },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Histogram of Final Speeds", color: "#fff" } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { display: true, text: "Speed (m/s)", color: "#fff" } }, y: { ...baseOptions.scales.y, title: { display: true, text: "Count", color: "#fff" } } } },
      }));
    }
    if (simTotalDistRef.current) {
      chartInstances.current.push(new Chart(simTotalDistRef.current, {
        type: "bar",
        data: {
          labels: vehIds,
          datasets: [{ label: "Total Distance", data: totalDistPerVeh, backgroundColor: "#8B5CF6" }],
        },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Total Distance per Vehicle", color: "#fff" } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { display: true, text: "Vehicle ID", color: "#fff" } }, y: { ...baseOptions.scales.y, title: { display: true, text: "Distance (m)", color: "#fff" } } } },
      }));
    }
    // Optimized Results (same data for now)
    if (optAvgSpeedRef.current) {
      chartInstances.current.push(new Chart(optAvgSpeedRef.current, {
        type: "line",
        data: {
          labels: allTimes,
          datasets: [{ label: "Avg Speed", data: avgSpeedOverTime, borderColor: "#3B82F6", backgroundColor: "#3B82F633", fill: true }],
        },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Average Speed Over Time (Optimized)", color: "#fff" } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { display: true, text: "Time (s)", color: "#fff" } }, y: { ...baseOptions.scales.y, title: { display: true, text: "Speed (m/s)", color: "#fff" } } } },
      }));
    }
    if (optVehCountRef.current) {
      chartInstances.current.push(new Chart(optVehCountRef.current, {
        type: "line",
        data: {
          labels: allTimes,
          datasets: [{ label: "Vehicle Count", data: vehCountOverTime, borderColor: "#818CF8", backgroundColor: "#818CF833", fill: true }],
        },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Vehicle Count Over Time (Optimized)", color: "#fff" } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { display: true, text: "Time (s)", color: "#fff" } }, y: { ...baseOptions.scales.y, title: { display: true, text: "Count", color: "#fff" } } } },
      }));
    }
    if (optFinalSpeedHistRef.current) {
      chartInstances.current.push(new Chart(optFinalSpeedHistRef.current, {
        type: "bar",
        data: {
          labels: histLabels,
          datasets: [{ label: "Final Speeds", data: finalSpeedHist, backgroundColor: "#F472B6" }],
        },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Histogram of Final Speeds (Optimized)", color: "#fff" } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { display: true, text: "Speed (m/s)", color: "#fff" } }, y: { ...baseOptions.scales.y, title: { display: true, text: "Count", color: "#fff" } } } },
      }));
    }
    if (optTotalDistRef.current) {
      chartInstances.current.push(new Chart(optTotalDistRef.current, {
        type: "bar",
        data: {
          labels: vehIds,
          datasets: [{ label: "Total Distance", data: totalDistPerVeh, backgroundColor: "#06B6D4" }],
        },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Total Distance per Vehicle (Optimized)", color: "#fff" } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { display: true, text: "Vehicle ID", color: "#fff" } }, y: { ...baseOptions.scales.y, title: { display: true, text: "Distance (m)", color: "#fff" } } } },
      }));
    }
    return () => {
      chartInstances.current.forEach((c) => c && c.destroy());
      chartInstances.current = [];
    };
  }, [simData]);

  if (loading) {
    return <div className="text-center text-gray-700 dark:text-gray-300 py-10">Loading simulation data...</div>;
  }
  if (error) {
    return <div className="text-center text-red-500 py-10">{error}</div>;
  }
  if (!simData) {
    return <div className="text-center text-gray-700 dark:text-gray-300 py-10">No simulation data found.</div>;
  }

  return (
    <div className="simulation-results-page bg-gradient-to-br from-gray-900 via-gray-800 to-black text-gray-100 min-h-screen">
      <Navbar />
      <div className="simRes-main-content py-8 px-6 overflow-y-auto">
        <div className="results max-w-7xl mx-auto">
          {/* Simulation meta info */}
          <div className="mb-10">
            <h1 className="simName text-4xl font-extrabold bg-gradient-to-r from-teal-400 to-emerald-500 bg-clip-text text-transparent mb-2">
              {simName}
            </h1>
            {simDesc && (
              <p className="simDesc text-lg text-gray-400 mb-4 leading-relaxed">{simDesc}</p>
            )}
            {simIntersections && simIntersections.length > 0 && (
              <div className="flex flex-wrap gap-2 mb-2">
                {simIntersections.map((intersection: string, idx: number) => (
                  <span key={idx} className="px-4 py-2 bg-white/10 backdrop-blur-md rounded-full text-sm font-medium text-teal-300 border border-teal-500/30 hover:bg-white/20 transition-all duration-300">
                    {intersection}
                  </span>
                ))}
              </div>
            )}
          </div>
          <div className="flex flex-col space-y-24">
            {/* Simulation Results Section */}
            <section className="visSection bg-white/5 backdrop-blur-md p-6 px-8 md:px-12 rounded-xl shadow-lg border border-gray-800/50 w-full max-w-full mx-auto text-center">
              <h2 className="text-2xl font-semibold mb-4 bg-gradient-to-r from-teal-400 to-emerald-500 bg-clip-text text-transparent">Simulation Results</h2>
              <div className="flex flex-col md:flex-row gap-8">
                {/* Stats column */}
                <div className="flex flex-row md:flex-col gap-4 md:gap-6 mb-2 md:mb-0 md:min-w-[180px] md:max-w-[220px]">
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-teal-500/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1">Average Speed</div>
                    <div className="text-2xl font-bold text-teal-300">{stats ? stats.avgSpeed.toFixed(2) : "..."} <span className="text-base font-normal">m/s</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-teal-500/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1">Max Speed</div>
                    <div className="text-2xl font-bold text-teal-300">{stats ? stats.maxSpeed.toFixed(2) : "..."} <span className="text-base font-normal">m/s</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-teal-500/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1">Min Speed</div>
                    <div className="text-2xl font-bold text-teal-300">{stats ? stats.minSpeed.toFixed(2) : "..."} <span className="text-base font-normal">m/s</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-teal-500/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1">Total Distance</div>
                    <div className="text-2xl font-bold text-teal-300">{stats ? stats.totalDistance.toFixed(2) : "..."} <span className="text-base font-normal">m</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-teal-500/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1"># Vehicles</div>
                    <div className="text-2xl font-bold text-teal-300">{stats ? stats.vehicleCount : "..."}</div>
                  </div>
                  {/* Traffic Light Stat Cards */}
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-yellow-400/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1"># TL Phases</div>
                    <div className="text-2xl font-bold text-yellow-400">{numPhases}</div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-yellow-400/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1">TL Cycle Duration</div>
                    <div className="text-2xl font-bold text-yellow-400">{totalCycle} <span className="text-base font-normal">s</span></div>
                  </div>
                </div>
                {/* Graphs grid */}
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 md:gap-12">
                  <div className="bg-white dark:bg-gray-900/60 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={simAvgSpeedRef} className="w-full h-full" /></div>
                  <div className="bg-white dark:bg-gray-900/60 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={simVehCountRef} className="w-full h-full" /></div>
                  <div className="bg-white dark:bg-gray-900/60 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={simFinalSpeedHistRef} className="w-full h-full" /></div>
                  <div className="bg-white dark:bg-gray-900/60 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={simTotalDistRef} className="w-full h-full" /></div>
                </div>
              </div>
              {/* Action buttons */}
              <div className="flex flex-row gap-4 justify-center mt-8">
                <button className="px-6 py-3 rounded-lg text-sm font-semibold bg-blue-500 hover:bg-blue-600 text-white shadow-md transition-all duration-300">Optimize</button>
                <button onClick={handleViewRendering} className="px-6 py-3 rounded-lg text-sm font-semibold bg-gray-200 hover:bg-gray-300 text-gray-800 shadow-md transition-all duration-300 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-gray-100">View Rendering</button>
              </div>
            </section>
            {/* Optimized Results Section */}
            <section className="visSection bg-white/5 backdrop-blur-md p-6 px-8 md:px-12 rounded-xl shadow-lg border border-gray-800/50 w-full max-w-full mx-auto mt-12 text-center">
              <h2 className="text-2xl font-semibold mb-4 bg-gradient-to-r from-blue-400 to-indigo-500 bg-clip-text text-transparent">Optimized Results</h2>
              <div className="flex flex-col md:flex-row gap-8">
                {/* Stats column */}
                <div className="flex flex-row md:flex-col gap-4 md:gap-6 mb-2 md:mb-0 md:min-w-[180px] md:max-w-[220px]">
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-blue-500/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1">Average Speed</div>
                    <div className="text-2xl font-bold text-blue-300">{stats ? stats.avgSpeed.toFixed(2) : "..."} <span className="text-base font-normal">m/s</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-blue-500/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1">Max Speed</div>
                    <div className="text-2xl font-bold text-blue-300">{stats ? stats.maxSpeed.toFixed(2) : "..."} <span className="text-base font-normal">m/s</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-blue-500/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1">Min Speed</div>
                    <div className="text-2xl font-bold text-blue-300">{stats ? stats.minSpeed.toFixed(2) : "..."} <span className="text-base font-normal">m/s</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-blue-500/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1">Total Distance</div>
                    <div className="text-2xl font-bold text-blue-300">{stats ? stats.totalDistance.toFixed(2) : "..."} <span className="text-base font-normal">m</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-blue-500/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1"># Vehicles</div>
                    <div className="text-2xl font-bold text-blue-300">{stats ? stats.vehicleCount : "..."}</div>
                  </div>
                  {/* Traffic Light Stat Cards */}
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-yellow-400/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1"># TL Phases</div>
                    <div className="text-2xl font-bold text-yellow-400">{numPhases}</div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-gray-900/80 border border-yellow-400/30 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm text-gray-400 mb-1">TL Cycle Duration</div>
                    <div className="text-2xl font-bold text-yellow-400">{totalCycle} <span className="text-base font-normal">s</span></div>
                  </div>
                </div>
                {/* Graphs grid */}
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 md:gap-12">
                  <div className="bg-white dark:bg-gray-900/60 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={optAvgSpeedRef} className="w-full h-full" /></div>
                  <div className="bg-white dark:bg-gray-900/60 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={optVehCountRef} className="w-full h-full" /></div>
                  <div className="bg-white dark:bg-gray-900/60 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={optFinalSpeedHistRef} className="w-full h-full" /></div>
                  <div className="bg-white dark:bg-gray-900/60 outline outline-1 outline-gray-300 dark:outline-gray-700 rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={optTotalDistRef} className="w-full h-full" /></div>
                </div>
              </div>
              {/* Action button */}
              <div className="flex flex-row gap-4 justify-center mt-8">
                <button onClick={handleViewRendering} className="px-6 py-3 rounded-lg text-sm font-semibold bg-gray-200 hover:bg-gray-300 text-gray-800 shadow-md transition-all duration-300 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-gray-100">View Rendering</button>
              </div>
            </section>
          </div>
        </div>
      </div>
      <Footer />
      <HelpMenu />
    </div>
  );
};

export default SimulationResults;
