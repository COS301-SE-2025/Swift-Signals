import React, { useEffect, useState, useRef } from "react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import "../styles/SimulationResults.css";
import HelpMenu from "../components/HelpMenu";
import { Chart, registerables } from "chart.js";
import { useLocation } from "react-router-dom";

Chart.register(...registerables);

type Position = { time: number; x: number; y: number; speed: number };
type Vehicle = { id: string; positions: Position[] };

function computeStats(vehicles: Vehicle[]) {
  let totalSpeed = 0, maxSpeed = -Infinity, minSpeed = Infinity, speedCount = 0;
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
  vehicles.forEach((veh) => veh.positions.forEach((pos) => timeSet.add(pos.time)));
  return Array.from(timeSet).sort((a, b) => a - b);
}

function getAverageSpeedOverTime(vehicles: Vehicle[], allTimes: number[]): number[] {
  return allTimes.map((t) => {
    let sum = 0, count = 0;
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

function getVehicleCountOverTime(vehicles: Vehicle[], allTimes: number[]): number[] {
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

function getVehicleIds(vehicles: Vehicle[]): string[] {
  return vehicles.map((veh: Vehicle) => veh.id);
}

function getHistogramData(data: number[], binSize: number, maxVal: number): { counts: number[], labels: string[] } {
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
        labels.push(`${(i * binSize).toFixed(0)}-${((i + 1) * binSize).toFixed(0)}`);
    }

    return { counts, labels };
}

function downsampleData(labels: any[], data: any[], maxPoints: number) {
  if (labels.length <= maxPoints) {
    return { downsampledLabels: labels, downsampledData: data };
  }

  const downsampledLabels = [];
  const downsampledData = [];
  const step = Math.ceil(labels.length / maxPoints);

  for (let i = 0; i < labels.length; i += step) {
    downsampledLabels.push(labels[i]);
    downsampledData.push(data[i]);
  }

  return { downsampledLabels, downsampledData };
}

const SimulationResults: React.FC = () => {
  const [simData, setSimData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const location = useLocation();
  const simInfo = location.state || {};
  
  const simName = simInfo.name || simData?.name || "Simulation";
  const simDesc = simInfo.description || simData?.description || "";
  const simIntersections = simInfo.intersections || simData?.intersections || [];

  const chartInstances = useRef<Chart[]>([]);
  
  const chartRefs = {
    simAvgSpeedRef: useRef<HTMLCanvasElement | null>(null),
    simVehCountRef: useRef<HTMLCanvasElement | null>(null),
    simFinalSpeedHistRef: useRef<HTMLCanvasElement | null>(null),
    simTotalDistHistRef: useRef<HTMLCanvasElement | null>(null),
    optAvgSpeedRef: useRef<HTMLCanvasElement | null>(null),
    optVehCountRef: useRef<HTMLCanvasElement | null>(null),
    optFinalSpeedHistRef: useRef<HTMLCanvasElement | null>(null),
    optTotalDistHistRef: useRef<HTMLCanvasElement | null>(null),
  };

  useEffect(() => {
    // Fetching the simulation data from the provided JSON file
    fetch("/simulation_output (1).json")
      .then((res) => {
        if (!res.ok) {
          throw new Error(`HTTP error! status: ${res.status}`);
        }
        return res.json();
      })
      .then((data) => {
        setSimData(data);
        setLoading(false);
      })
      .catch((err) => {
        setError("Failed to load simulation data. Please check the file path and format.");
        setLoading(false);
        console.error(err);
      });
  }, []);

  useEffect(() => {
    // Cleanup previous chart instances
    chartInstances.current.forEach((c) => c?.destroy());
    chartInstances.current = [];

    if (!simData || !simData.vehicles) return;

    // --- Prepare data for charts ---
    const stats = computeStats(simData.vehicles);
    const allTimes = getAllTimes(simData.vehicles);
    const avgSpeedOverTime = getAverageSpeedOverTime(simData.vehicles, allTimes);
    const vehCountOverTime = getVehicleCountOverTime(simData.vehicles, allTimes);
    const totalDistPerVeh = getTotalDistancePerVehicle(simData.vehicles);

    // Downsample time-series data for cleaner line charts
    const MAX_TIME_POINTS = 100;
    const { downsampledLabels: timeLabels, downsampledData: downsampledAvgSpeed } = downsampleData(allTimes, avgSpeedOverTime, MAX_TIME_POINTS);
    const { downsampledData: downsampledVehCount } = downsampleData(allTimes, vehCountOverTime, MAX_TIME_POINTS);
    
    // Prepare histogram data
    const { counts: finalSpeedHist, labels: finalSpeedHistLabels } = getHistogramData(stats.finalSpeeds, 2, 40);
    const maxDist = totalDistPerVeh.length > 0 ? Math.max(...totalDistPerVeh) : 0;
    const { counts: totalDistHist, labels: totalDistHistLabels } = getHistogramData(totalDistPerVeh, 50, Math.ceil(maxDist / 50) * 50);

    // --- Chart Rendering ---
    const baseOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
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
            title: { display: true, color: "#fff", font: { size: 14 } } 
        },
        y: { 
            grid: { color: "rgba(255,255,255,0.1)" }, 
            ticks: { color: "#ccc" }, 
            title: { display: true, color: "#fff", font: { size: 14 } },
            beginAtZero: true
        },
      },
    };
    
    const createChart = (ref: React.RefObject<HTMLCanvasElement | null>, config: any) => {
        if (ref.current) {
            const chart = new Chart(ref.current, config);
            chartInstances.current.push(chart);
        }
    };

    // Simulation Charts
    createChart(chartRefs.simAvgSpeedRef, {
        type: "line",
        data: { labels: timeLabels, datasets: [{ label: "Average Speed", data: downsampledAvgSpeed, borderColor: "#2B9348", backgroundColor: "#48ac4d33", fill: true, tension: 0.3 }] },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Average Speed Over Time", color: "#fff", font: {size: 18} } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { ...baseOptions.scales.x.title, text: "Time (s)" } }, y: { ...baseOptions.scales.y, title: { ...baseOptions.scales.y.title, text: "Speed (m/s)" } } } },
    });
    createChart(chartRefs.simVehCountRef, {
        type: "line",
        data: { labels: timeLabels, datasets: [{ label: "Vehicle Count", data: downsampledVehCount, borderColor: "#0F5BA7", backgroundColor: "#60A5FA33", fill: true, tension: 0.3 }] },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Vehicle Count Over Time", color: "#fff", font: {size: 18} } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { ...baseOptions.scales.x.title, text: "Time (s)" } }, y: { ...baseOptions.scales.y, title: { ...baseOptions.scales.y.title, text: "Count" } } } },
    });
    createChart(chartRefs.simFinalSpeedHistRef, {
        type: "bar",
        data: { labels: finalSpeedHistLabels, datasets: [{ label: "Final Speed Distribution", data: finalSpeedHist, backgroundColor: "#0F5BA7" }] },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Histogram of Final Speeds", color: "#fff", font: {size: 18} } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { ...baseOptions.scales.x.title, text: "Speed (m/s)" } }, y: { ...baseOptions.scales.y, title: { ...baseOptions.scales.y.title, text: "Number of Vehicles" } } } },
    });
    createChart(chartRefs.simTotalDistHistRef, {
        type: "bar",
        data: { labels: totalDistHistLabels, datasets: [{ label: "Total Distance Distribution", data: totalDistHist, backgroundColor: "#2B9348" }] },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Histogram of Total Distance", color: "#fff", font: {size: 18} } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { ...baseOptions.scales.x.title, text: "Distance (m)" } }, y: { ...baseOptions.scales.y, title: { ...baseOptions.scales.y.title, text: "Number of Vehicles" } } } },
    });

    // Optimized Charts (using same data for demonstration)
    createChart(chartRefs.optAvgSpeedRef, {
        type: "line",
        data: { labels: timeLabels, datasets: [{ label: "Average Speed (Optimized)", data: downsampledAvgSpeed, borderColor: "#2B9348", backgroundColor: "#48ac4d33", fill: true, tension: 0.3 }] },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Average Speed (Optimized)", color: "#fff", font: {size: 18} } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { ...baseOptions.scales.x.title, text: "Time (s)" } }, y: { ...baseOptions.scales.y, title: { ...baseOptions.scales.y.title, text: "Speed (m/s)" } } } },
    });
    createChart(chartRefs.optVehCountRef, {
        type: "line",
        data: { labels: timeLabels, datasets: [{ label: "Vehicle Count (Optimized)", data: downsampledVehCount, borderColor: "#0F5BA7", backgroundColor: "#60A5FA33", fill: true, tension: 0.3 }] },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Vehicle Count (Optimized)", color: "#fff", font: {size: 18} } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { ...baseOptions.scales.x.title, text: "Time (s)" } }, y: { ...baseOptions.scales.y, title: { ...baseOptions.scales.y.title, text: "Count" } } } },
    });
    createChart(chartRefs.optFinalSpeedHistRef, {
        type: "bar",
        data: { labels: finalSpeedHistLabels, datasets: [{ label: "Final Speed Distribution (Optimized)", data: finalSpeedHist, backgroundColor: "#0F5BA7" }] },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Final Speeds Hist. (Optimized)", color: "#fff", font: {size: 18} } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { ...baseOptions.scales.x.title, text: "Speed (m/s)" } }, y: { ...baseOptions.scales.y, title: { ...baseOptions.scales.y.title, text: "Number of Vehicles" } } } },
    });
    createChart(chartRefs.optTotalDistHistRef, {
        type: "bar",
        data: { labels: totalDistHistLabels, datasets: [{ label: "Total Distance Distribution (Optimized)", data: totalDistHist, backgroundColor: "#2B9348" }] },
        options: { ...baseOptions, plugins: { ...baseOptions.plugins, title: { display: true, text: "Total Distance Hist. (Optimized)", color: "#fff", font: {size: 18} } }, scales: { ...baseOptions.scales, x: { ...baseOptions.scales.x, title: { ...baseOptions.scales.x.title, text: "Distance (m)" } }, y: { ...baseOptions.scales.y, title: { ...baseOptions.scales.y.title, text: "Number of Vehicles" } } } },
    });

    // Cleanup function
    return () => {
      chartInstances.current.forEach((c) => c?.destroy());
      chartInstances.current = [];
    };
  }, [simData]); // Re-render charts when simData changes

  if (loading) return <div className="text-center text-gray-700 dark:text-gray-300 py-10">Loading simulation data...</div>;
  if (error) return <div className="text-center text-red-500 py-10">{error}</div>;
  if (!simData) return <div className="text-center text-gray-700 dark:text-gray-300 py-10">No simulation data found.</div>;

  const stats = computeStats(simData.vehicles);
  const { numPhases, totalCycle } = simData.intersection?.trafficLights?.[0] 
    ? {
        numPhases: simData.intersection.trafficLights[0].phases.length,
        totalCycle: simData.intersection.trafficLights[0].phases.reduce((sum: number, p: any) => sum + (p.duration || 0), 0)
      }
    : { numPhases: 0, totalCycle: 0 };
    
  const handleViewComparison = () => {
    window.location.href = '/comparison-rendering';
  };
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
                  <span key={idx} className="px-4 py-2 bg-white/10 dark:bg-[#161B22] backdrop-blur-md rounded-full text-sm font-medium text-[#0F5BA7] border-2 border-[#0F5BA7] hover:bg-white/20 transition-all duration-300">
                    {intersection}
                  </span>
                ))}
              </div>
            )}
          </div>
          <div className="flex flex-col space-y-24">
            {/* Simulation Results Section */}
            <section className="visSection simulation-section bg-white/5 backdrop-blur-md p-6 px-8 md:px-12 rounded-xl shadow-lg border border-gray-800/50 w-full max-w-full mx-auto text-center">
              <h2 className="text-2xl font-semibold mb-4 bg-[#0F5BA7] bg-clip-text text-transparent">Simulation Results</h2>
              <div className="flex flex-col md:flex-row gap-8">
                {/* Stats column */}
                <div className="flex flex-row md:flex-col gap-4 md:gap-6 mb-2 md:mb-0 md:min-w-[180px] md:max-w-[220px]">
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1">Average Speed</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{stats ? stats.avgSpeed.toFixed(2) : "..."} <span className="text-base text-[#0F5BA7] font-normal">m/s</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1">Max Speed</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{stats ? stats.maxSpeed.toFixed(2) : "..."} <span className="text-base text-[#0F5BA7] font-normal">m/s</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1">Min Speed</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{stats ? stats.minSpeed.toFixed(2) : "..."} <span className="text-base text-[#0F5BA7] font-normal">m/s</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1">Total Distance</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{stats ? stats.totalDistance.toFixed(2) : "..."} <span className="text-base text-[#0F5BA7] font-normal">m</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-teal-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1"># Vehicles</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{stats ? stats.vehicleCount : "..."}</div>
                  </div>
                  {/* Traffic Light Stat Cards */}
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-yellow-400/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1"># TL Phases</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{numPhases}</div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-yellow-400/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1">TL Cycle Duration</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{totalCycle} <span className="text-base text-[#0F5BA7] font-normal">s</span></div>
                  </div>
                </div>
                {/* Graphs grid */}
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 md:gap-12">
                  <div className="bg-white dark:bg-[#161B22] outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={chartRefs.simAvgSpeedRef} className="w-full h-full" /></div>
                  <div className="bg-white dark:bg-[#161B22] outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={chartRefs.simVehCountRef} className="w-full h-full" /></div>
                  <div className="bg-white dark:bg-[#161B22] outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={chartRefs.simFinalSpeedHistRef} className="w-full h-full" /></div>
                  <div className="bg-white dark:bg-[#161B22] outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={chartRefs.simTotalDistHistRef} className="w-full h-full" /></div>
                </div>
              </div>
              {/* Action buttons */}
              <div className="flex flex-col sm:flex-row gap-4 justify-center mt-12">
                <button className="px-8 py-3 text-base font-bold text-white bg-gradient-to-r from-blue-500 to-blue-600 rounded-xl shadow-lg shadow-blue-500/50 transform transition-all duration-300 ease-in-out hover:scale-105 hover:shadow-xl hover:shadow-blue-500/60 focus:outline-none focus:ring-4 focus:ring-blue-300">
                  Optimize
                </button>
                <button 
                  onClick={handleViewComparison} 
                  className="px-8 py-3 text-base font-bold text-[#0F5BA7] bg-transparent border-2 border-[#0F5BA7] rounded-xl transform transition-all duration-300 ease-in-out hover:bg-[#0F5BA7] hover:text-white hover:scale-105 focus:outline-none focus:ring-4 focus:ring-gray-600"
                >
                  View Comparison Rendering
                </button>
              </div>
            </section>
            {/* Optimized Results Section */}
            <section className="visSection optimized-section bg-white/5 backdrop-blur-md p-6 px-8 md:px-12 rounded-xl shadow-lg border border-gray-800/50 w-full max-w-full mx-auto mt-12 text-center">
              <h2 className="text-2xl font-semibold mb-4 bg-[#0F5BA7] bg-clip-text text-transparent">Optimized Results</h2>
              <div className="flex flex-col md:flex-row gap-8">
                {/* Stats column */}
                <div className="flex flex-row md:flex-col gap-4 md:gap-6 mb-2 md:mb-0 md:min-w-[180px] md:max-w-[220px]">
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-blue-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1">Average Speed</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{stats ? stats.avgSpeed.toFixed(2) : "..."} <span className="text-base text-[#0F5BA7] font-normal">m/s</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-blue-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1">Max Speed</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{stats ? stats.maxSpeed.toFixed(2) : "..."} <span className="text-base text-[#0F5BA7] font-normal">m/s</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-blue-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1">Min Speed</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{stats ? stats.minSpeed.toFixed(2) : "..."} <span className="text-base text-[#0F5BA7] font-normal">m/s</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-blue-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1">Total Distance</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{stats ? stats.totalDistance.toFixed(2) : "..."} <span className="text-base text-[#0F5BA7] font-normal">m</span></div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-blue-500/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1"># Vehicles</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{stats ? stats.vehicleCount : "..."}</div>
                  </div>
                  {/* Traffic Light Stat Cards */}
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-yellow-400/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1"># TL Phases</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{numPhases}</div>
                  </div>
                  <div className="stat-cube bg-white dark:bg-[#161B22] border border-yellow-400/30 outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-xl p-6 text-center shadow-md">
                    <div className="text-sm font-bold text-gray-600 mb-1">TL Cycle Duration</div>
                    <div className="text-2xl font-bold text-[#0F5BA7]">{totalCycle} <span className="text-base text-[#0F5BA7] font-normal">s</span></div>
                  </div>
                </div>
                {/* Graphs grid */}
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 md:gap-12">
                  <div className="bg-white dark:bg-[#161B22] outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={chartRefs.optAvgSpeedRef} className="w-full h-full" /></div>
                  <div className="bg-white dark:bg-[#161B22] outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={chartRefs.optVehCountRef} className="w-full h-full" /></div>
                  <div className="bg-white dark:bg-[#161B22] outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={chartRefs.optFinalSpeedHistRef} className="w-full h-full" /></div>
                  <div className="bg-white dark:bg-[#161B22] outline outline-2 outline-gray-300 dark:outline-[#388BFD] rounded-2xl p-6 h-[28rem] min-w-[400px] max-w-[900px] w-full flex items-center justify-center"><canvas ref={chartRefs.optTotalDistHistRef} className="w-full h-full" /></div>
                </div>
              </div>
              {/* Action button */}
              <div className="flex flex-row gap-4 justify-center mt-12">
                {/* <button 
                  onClick={handleViewRendering} 
                  className="px-8 py-3 text-base font-bold text-[#0F5BA7] bg-transparent border-2 border-[#0F5BA7] rounded-xl transform transition-all duration-300 ease-in-out hover:bg-[#0F5BA7] hover:text-white hover:scale-105 focus:outline-none focus:ring-4 focus:ring-gray-600"
                >
                  View Rendering
                </button> */}
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