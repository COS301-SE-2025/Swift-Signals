import React, { useEffect, useState, useRef } from 'react';
import { useLocation } from 'react-router-dom';
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';
import '../styles/SimulationResults.css';
import { Chart, registerables } from 'chart.js';

// Register Chart.js components
Chart.register(...registerables);

const SimulationResults: React.FC = () => {
  const location = useLocation();
  const [simulationData, setSimulationData] = useState<{ name: string; description: string; intersections: string[] } | null>(null);
  const [params, setParams] = useState<{ speed: number; density: number }>({ speed: 50, density: 30 });
  const [optimizedParams, setOptimizedParams] = useState<{ speed: number; density: number }>({ speed: 70, density: 20 });
  const [isRunningSim, setIsRunningSim] = useState(false);
  const [isRunningOpt, setIsRunningOpt] = useState(false);
  const [showFooter, setShowFooter] = useState(false);
  const [fullScreenChart, setFullScreenChart] = useState<'simulation' | 'optimized' | null>(null);
  const simCanvasRef = useRef<HTMLCanvasElement | null>(null);
  const optCanvasRef = useRef<HTMLCanvasElement | null>(null);

  useEffect(() => {
    if (location.state) {
      setSimulationData(location.state);
    }

    const handleScroll = () => {
      const scrollPosition = window.scrollY + window.innerHeight;
      const documentHeight = document.documentElement.scrollHeight;
      setShowFooter(scrollPosition >= documentHeight - 10);
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, [location.state]);

  const handleParamChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setParams(prev => ({ ...prev, [name]: Number(value) }));
  };

  const handleOptimizedParamChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setOptimizedParams(prev => ({ ...prev, [name]: Number(value) }));
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: {
        backgroundColor: 'rgba(0, 0, 0, 0.9)',
        cornerRadius: 10,
        padding: 12,
        bodyFont: { size: 14, weight: 'bold' as const },
      },
    },
    scales: {
      x: { display: true, title: { display: true, text: 'Parameter', font: { size: 14, weight: 'bold' as const } }, grid: { color: 'rgba(255, 255, 255, 0.1)' } },
      y: { beginAtZero: true, max: 100, title: { display: true, text: 'Value', font: { size: 14, weight: 'bold' as const } }, grid: { color: 'rgba(255, 255, 255, 0.1)' } },
    },
    animation: {
      duration: 1200,
      easing: 'easeOutQuart' as const,
    },
  };

  const simulationChartData = {
    labels: ['Speed', 'Density'],
    datasets: [
      {
        label: 'Simulation',
        data: [params.speed, params.density],
        backgroundColor: ['#34D399', '#10B981'],
        borderColor: ['#34D399', '#10B981'],
        borderWidth: 2,
        barThickness: 20,
      },
    ],
  };

  const optimizedChartData = {
    labels: ['Speed', 'Density'],
    datasets: [
      {
        label: 'Optimized',
        data: [optimizedParams.speed, optimizedParams.density],
        backgroundColor: ['#3B82F6', '#2563EB'],
        borderColor: ['#3B82F6', '#2563EB'],
        borderWidth: 2,
        barThickness: 20,
      },
    ],
  };

  useEffect(() => {
    let simChart: Chart | undefined;
    if (simCanvasRef.current) {
      simChart = new Chart(simCanvasRef.current, {
        type: 'bar',
        data: simulationChartData,
        options: chartOptions,
      });
    }
    return () => simChart?.destroy();
  }, [params, isRunningSim]);

  useEffect(() => {
    let optChart: Chart | undefined;
    if (optCanvasRef.current) {
      optChart = new Chart(optCanvasRef.current, {
        type: 'bar',
        data: optimizedChartData,
        options: chartOptions,
      });
    }
    return () => optChart?.destroy();
  }, [optimizedParams, isRunningOpt]);

  useEffect(() => {
    let fullScreenChartInstance: Chart | undefined;
    const canvas = document.getElementById('fullScreenChart') as HTMLCanvasElement | null;
    if (canvas && fullScreenChart) {
      const data = fullScreenChart === 'simulation' ? simulationChartData : optimizedChartData;
      fullScreenChartInstance = new Chart(canvas, {
        type: 'bar',
        data,
        options: {
          ...chartOptions,
          plugins: {
            ...chartOptions.plugins,
            title: {
              display: true,
              text: fullScreenChart === 'simulation' ? 'Simulation Visualization' : 'Optimized Visualization',
              font: { size: 18, weight: 'bold' as const },
              color: '#fff',
            },
          },
        },
      });
    }
    return () => fullScreenChartInstance?.destroy();
  }, [fullScreenChart]);

  const handleRunSimulation = () => {
    setIsRunningSim(true);
    setTimeout(() => setIsRunningSim(false), 2000);
  };

  const handleRunOptimized = () => {
    setIsRunningOpt(true);
    setTimeout(() => setIsRunningOpt(false), 2000);
  };

  const handleOptimize = () => {
    setOptimizedParams({ speed: Math.min(params.speed + 20, 100), density: Math.max(params.density - 10, 0) });
  };

  const openFullScreen = (chartType: 'simulation' | 'optimized') => {
    setFullScreenChart(chartType);
  };

  const closeFullScreen = () => {
    setFullScreenChart(null);
  };

  if (!simulationData) return <div className="text-center text-gray-700 dark:text-gray-300 py-10">Loading...</div>;

  return (
    <div className="simulation-results-page bg-gradient-to-br from-gray-900 via-gray-800 to-black text-gray-100">
      <Navbar />
      <div className="simRes-main-content py-8 px-6 overflow-y-auto">
        <div className="results max-w-7xl mx-auto">
          <h1 className="simName text-4xl font-extrabold bg-gradient-to-r from-teal-400 to-emerald-500 bg-clip-text text-transparent mb-4">{simulationData.name}</h1>
          <p className="simDesc text-lg text-gray-400 mb-6 leading-relaxed">{simulationData.description}</p>
          <div className="flex flex-wrap gap-3 mb-8">
            {simulationData.intersections.map((intersection, index) => (
              <span key={index} className="simInt inline-flex items-center px-4 py-2 bg-white/10 backdrop-blur-md rounded-full text-sm font-medium text-teal-300 border border-teal-500/30 hover:bg-white/20 transition-all duration-300">
                {intersection}
              </span>
            ))}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-60">
            {/* Simulation Visualization Section */}
            <div className="visSection bg-white/5 backdrop-blur-md p-6 rounded-xl shadow-lg border border-gray-800/50 hover:shadow-xl transition-all duration-300">
              <h2 className="text-2xl font-semibold bg-gradient-to-r from-teal-400 to-emerald-500 bg-clip-text text-transparent mb-4">Simulation Visualization</h2>
              <div className="flex flex-col md:flex-row gap-6">
                <div className="flex-2 visualization-box h-96 relative bg-gray-900/50 border border-gray-700 rounded-lg overflow-hidden">
                  <div onClick={() => openFullScreen('simulation')} className="chart-container cursor-pointer">
                    <canvas id="simulationChart" ref={simCanvasRef} className="w-full h-full" />
                  </div>
                  <div className="absolute bottom-6 right-6 flex space-x-4">
                    <button
                      onClick={handleRunSimulation}
                      disabled={isRunningSim}
                      className={`px-6 py-3 rounded-lg text-sm font-semibold ${isRunningSim ? 'bg-gray-600 cursor-not-allowed' : 'bg-gradient-to-r from-green-500 to-emerald-600 hover:from-green-600 hover:to-emerald-700'} text-white shadow-md hover:shadow-lg transition-all duration-300`}
                    >
                      {isRunningSim ? 'Running...' : 'Run'}
                    </button>
                    <button
                      onClick={handleOptimize}
                      className="px-6 py-3 rounded-lg text-sm font-semibold bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 text-white shadow-md hover:shadow-lg transition-all duration-300"
                    >
                      Optimize
                    </button>
                  </div>
                </div>
                <div className="parameters md:w-1/3 bg-gray-800/70 p-6 rounded-xl border border-gray-700">
                  <h3 className="text-xl font-medium bg-gradient-to-r from-teal-400 to-emerald-500 bg-clip-text text-transparent mb-4">Parameters</h3>
                  <div className="space-y-5">
                    <div>
                      <label className="block text-sm font-medium text-gray-500 dark:text-gray-300 mb-2">Speed</label>
                      <input
                        type="number"
                        name="speed"
                        value={params.speed}
                        onChange={handleParamChange}
                        className="paramInput w-full p-3 rounded-lg bg-gray-700/50 border border-gray-600 text-gray-100 focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all duration-300"
                        min="0"
                        max="100"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-500 dark:text-gray-300 mb-2">Density</label>
                      <input
                        type="number"
                        name="density"
                        value={params.density}
                        onChange={handleParamChange}
                        className="paramInput w-full p-3 rounded-lg bg-gray-700/50 border border-gray-600 text-gray-100 focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all duration-300"
                        min="0"
                        max="100"
                      />
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Optimized Visualization Section */}
            <div className="visSection bg-white/5 backdrop-blur-md p-6 rounded-xl shadow-lg border border-gray-800/50 hover:shadow-xl transition-all duration-300">
              <h2 className="text-2xl font-semibold bg-gradient-to-r from-blue-400 to-indigo-500 bg-clip-text text-transparent mb-4">Optimized Visualization</h2>
              <div className="flex flex-col md:flex-row gap-6">
                <div className="flex-2 visualization-box h-96 relative bg-gray-900/50 border border-gray-700 rounded-lg overflow-hidden">
                  <div onClick={() => openFullScreen('optimized')} className="chart-container cursor-pointer">
                    <canvas id="optimizedChart" ref={optCanvasRef} className="w-full h-full" />
                  </div>
                  <div className="absolute bottom-6 right-6">
                    <button
                      onClick={handleRunOptimized}
                      disabled={isRunningOpt}
                      className={`px-6 py-3 rounded-lg text-sm font-semibold ${isRunningOpt ? 'bg-gray-600 cursor-not-allowed' : 'bg-gradient-to-r from-green-500 to-emerald-600 hover:from-green-600 hover:to-emerald-700'} text-white shadow-md hover:shadow-lg transition-all duration-300`}
                    >
                      {isRunningOpt ? 'Running...' : 'Run'}
                    </button>
                  </div>
                </div>
                <div className="parameters md:w-1/3 bg-gray-800/70 p-6 rounded-xl border border-gray-700">
                  <h3 className="text-xl font-medium bg-gradient-to-r from-blue-400 to-indigo-500 bg-clip-text text-transparent mb-4">Optimized Parameters</h3>
                  <div className="space-y-5">
                    <div>
                      <label className="block text-sm font-medium text-gray-500 dark:text-gray-300 mb-2">Speed</label>
                      <input
                        type="number"
                        name="speed"
                        value={optimizedParams.speed}
                        onChange={handleOptimizedParamChange}
                        className="paramInput w-full p-3 rounded-lg bg-gray-700/50 border border-gray-600 text-gray-100 focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-300"
                        min="0"
                        max="100"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-500 dark:text-gray-300 mb-2">Density</label>
                      <input
                        type="number"
                        name="density"
                        value={optimizedParams.density}
                        onChange={handleOptimizedParamChange}
                        className="paramInput w-full p-3 rounded-lg bg-gray-700/50 border border-gray-600 text-gray-100 focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-300"
                        min="0"
                        max="100"
                      />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Full-Screen Modal */}
          {fullScreenChart && (
            <div className="fullscreen-modal fixed inset-0 bg-black bg-opacity-80 flex items-center justify-center z-50">
              <div className="relative w-full h-full max-w-7xl max-h-[90vh] p-4">
                <button
                  onClick={closeFullScreen}
                  className="absolute top-4 right-4 text-white bg-gray-800/70 rounded-full p-2 hover:bg-gray-700 transition-all duration-300"
                >
                  ✕
                </button>
                <canvas id="fullScreenChart" className="w-full h-full" />
              </div>
            </div>
          )}

          {/* Additional Content to Ensure Scrolling */}
          <div className="mt-12">
            <h2 className="statTitle text-2xl font-semibold bg-gradient-to-r from-teal-400 to-emerald-500 bg-clip-text text-transparent mb-4">Statistics</h2>
            <p className="text-gray-400 mb-6 leading-relaxed">
              Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
            </p>
            <p className="text-gray-400 mb-6 leading-relaxed">
              Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt.
            </p>
            <p className="text-gray-400 mb-6 leading-relaxed">
              Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
            </p>
            <p className="text-gray-400 mb-6 leading-relaxed">
              Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt.
            </p>
          </div>
        </div>
      </div>
      {showFooter && <Footer />}
    </div>
  );
};

export default SimulationResults;