import React, { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { Eye, Trash2, ChevronDown } from "lucide-react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import "../styles/Simulations.css";
import "@fortawesome/fontawesome-free/css/all.min.css";
import HelpMenu from "../components/HelpMenu";

const API_BASE_URL = "http://localhost:9090";

const getAuthToken = () => {
  return localStorage.getItem("authToken");
};

// #region API and Data Types
interface ApiIntersection {
  id: string;
  name: string;
  details: {
    address: string;
    city: string;
    province: string;
  };
  default_parameters: {
    // This object contains the nested simulation data
    optimisation_type: string;
    simulation_parameters: {
      green: number;
      yellow: number;
      red: number;
      speed: number;
      seed: number;
      intersection_type: string;
    };
  };
  traffic_density: string;
  status?: string;
  created_at?: string;
  run_count?: number;
  last_run_at?: string;
}

//  UPDATED: Changed the data structure for the table
interface SimulationData {
  id: number;
  backendId: string;
  intersection: string;
  trafficDensity: string; // New: "High", "Medium", or "Low"
  speed: number; // New: Vehicle speed in km/h
  status: string;
}

const NewSimulationModal: React.FC<{
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: {
    name: string;
    description: string;
    intersections: string[];
  }) => void;
  intersections: string[];
  type: "simulations" | "optimizations";
}> = ({ isOpen, onClose, onSubmit, intersections, type }) => {
  const [simulationName, setSimulationName] = useState("");
  const [simulationDescription, setSimulationDescription] = useState("");
  const [selectedIntersections, setSelectedIntersections] = useState<string[]>(
    [],
  );

  const handleAddIntersection = (intersection: string) => {
    if (intersection && !selectedIntersections.includes(intersection)) {
      setSelectedIntersections([...selectedIntersections, intersection]);
    }
  };

  const handleRemoveIntersection = (intersection: string) => {
    setSelectedIntersections(
      selectedIntersections.filter((item) => item !== intersection),
    );
  };

  const handleSubmit = () => {
    if (!simulationName || selectedIntersections.length === 0) {
      alert(
        "Please provide a simulation name and select at least one intersection.",
      );
      return;
    }
    const simulationData = {
      name: simulationName,
      description: simulationDescription,
      intersections: selectedIntersections,
    };
    onSubmit(simulationData);
    setSimulationName("");
    setSimulationDescription("");
    setSelectedIntersections([]);
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div className="simulation-modal-content bg-white dark:bg-[#161B22] rounded-lg shadow-xl w-full max-w-md p-6 relative max-h-[90vh] overflow-y-auto">
        <button
          onClick={onClose}
          className="crossBtn absolute top-4 right-4 text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-100"
        ></button>
        <h2 className="text-xl font-bold text-gray-800 dark:text-gray-200 mb-4">
          New {type === "simulations" ? "Simulation" : "Optimization"}
        </h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Simulation Name
            </label>
            <input
              type="text"
              value={simulationName}
              onChange={(e) => setSimulationName(e.target.value)}
              className="simulation-name-input w-full p-2 rounded-md border-2 border-gray-300 dark:border-[#30363D] bg-white dark:bg-[#161B22] text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="Enter simulation name"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Simulation Description
            </label>
            <textarea
              value={simulationDescription}
              onChange={(e) => setSimulationDescription(e.target.value)}
              className="w-full p-2 rounded-md border-2 border-gray-300 dark:border-[#30363D] bg-white dark:bg-[#161B22] text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="Enter simulation description"
              rows={3}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Intersections
            </label>
            <div className="flex flex-wrap gap-2 mb-3">
              {selectedIntersections.map((intersection) => (
                <div
                  key={intersection}
                  className="intersection-pill flex items-center px-2 py-0.5 rounded-full bg-indigo-100 text-indigo-800 dark:bg-[#161B22] dark:text-indigo-100 text-xs"
                >
                  {intersection}
                  <button
                    onClick={() => handleRemoveIntersection(intersection)}
                    className="ml-1 text-indigo-600 hover:text-indigo-800 dark:text-indigo-300 dark:hover:text-indigo-100 remove-cross"
                  ></button>
                </div>
              ))}
            </div>
            <select
              value=""
              onChange={(e) => handleAddIntersection(e.target.value)}
              className="w-full p-2 rounded-md border-2 border-gray-300 dark:border-[#30363D] bg-white dark:bg-[#161B22] text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-[#2B9348]"
            >
              <option value="">Select an intersection</option>
              {intersections.map((intersection) => (
                <option key={intersection} value={intersection}>
                  {intersection}
                </option>
              ))}
            </select>
          </div>
        </div>
        <div className="mt-6 flex justify-end space-x-2">
          <button
            onClick={onClose}
            className="px-4 py-2 rounded-md text-sm font-medium bg-gray-200 text-gray-700 hover:bg-gray-300 dark:bg-[#161B22] dark:text-gray-200 dark:border-2 dark:border-[#DA3633] dark:hover:bg-[#DA3633] transition-all duration-300"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            className="create-simulation-submit-btn px-4 py-2 rounded-md text-sm font-medium bg-[#0F5BA7] dark:bg-[#388BFD] text-white hover:from-green-600 hover:to-green-700 dark:from-green-400 dark:to-green-500 dark:hover:from-green-500 dark:hover:to-green-600 transition-all duration-300"
          >
            Create
          </button>
        </div>
      </div>
    </div>
  );
};

const SimulationTable: React.FC<{
  simulations: SimulationData[];
  currentPage: number;
  setCurrentPage: (page: number) => void;
  onViewResults: (backendId: string, intersectionName: string) => void;
}> = ({ simulations, currentPage, setCurrentPage, onViewResults }) => {
  const rowsPerPage = 4;
  const totalPages = Math.ceil(simulations.length / rowsPerPage);
  const startIndex = currentPage * rowsPerPage;
  const endIndex = startIndex + rowsPerPage;
  const paginatedSimulations = simulations.slice(startIndex, endIndex);

  const handleDelete = (backendId: string) => {
    alert(`Deleting simulation with backend ID ${backendId}`);
    // Replace with actual delete logic
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const statusClass = (status: string) => {
    switch (status) {
      case "optimised":
        return "bg-green-200 text-green-800 border-green-300";
      case "unoptimised":
        return "bg-yellow-200 text-yellow-800 border-yellow-300";
      case "Failed":
        return "bg-red-200 text-red-800 border-red-300";
      default:
        return "bg-gray-200 text-gray-800 border-gray-300";
    }
  };

  return (
    <div className="simTable bg-white dark:bg-[#161B22] shadow-md rounded-lg overflow-hidden table-fixed-height relative">
      {simulations.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-16 px-4">
          <div className="text-gray-400 dark:text-gray-500 mb-4">
            <svg
              className="w-16 h-16 mx-auto"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="1.5"
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-gray-600 dark:text-gray-400 mb-2">
            No Data to Display
          </h3>
          <p className="text-sm text-gray-500 dark:text-gray-500 text-center max-w-sm">
            There are no simulations available at the moment. Create a new
            simulation to get started.
          </p>
        </div>
      ) : (
        <>
          <table className="simulationTable min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className="simTableHead bg-gray-50 dark:bg-[#161B22]">
              <tr>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  No.
                </th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Intersection
                </th>
                {/*  UPDATED: Changed table headers */}
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Traffic Density
                </th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Speed
                </th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-[#161B22] divide-y divide-gray-200 dark:divide-gray-700">
              {paginatedSimulations.map((sim) => (
                <tr key={sim.backendId}>
                  <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-gray-200">
                    #{sim.id}
                  </td>
                  <td className="intersectionCell px-4 py-3 whitespace-wrap text-sm text-gray-900 dark:text-gray-200">
                    {sim.intersection}
                  </td>
                  {/*  UPDATED: Changed table cells to display new data */}
                  <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-gray-200">
                    {sim.trafficDensity}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-gray-200">
                    {sim.speed} km/h
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm">
                    <span
                      className={`sim-status inline-flex items-center px-3 py-1 rounded-md border ${statusClass(
                        sim.status,
                      )}`}
                    >
                      {sim.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm">
                    <div className="flex flex-col space-y-2">
                      <button
                        onClick={() =>
                          onViewResults(sim.backendId, sim.intersection)
                        }
                        className="viewBtn text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300 text-sm font-medium w-full text-center"
                        title="View Results"
                      >
                        <Eye size={18} strokeWidth={2} />
                      </button>
                      <button
                        onClick={() => handleDelete(sim.backendId)}
                        className="deleteBtn text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300 text-sm font-medium w-full text-center"
                        title="Delete Simulation"
                      >
                        <Trash2 size={18} strokeWidth={2} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {simulations.length > rowsPerPage && (
            <div className="pagination absolute bottom-0 left-0 right-0 flex justify-center items-center p-4 space-x-2 bg-white dark:bg-[#161B22] border-t border-gray-200 dark:border-gray-700">
              <button
                onClick={() => handlePageChange(currentPage - 1)}
                disabled={currentPage === 0}
                className={`px-3 py-1 rounded-full text-sm font-medium bg-[#0F5BA7] dark:bg-[#388BFD] text-white hover:from-indigo-600 hover:to-indigo-700 dark:from-indigo-400 dark:to-indigo-500 dark:hover:from-indigo-500 dark:hover:to-indigo-600 transition-all duration-300 ${
                  currentPage === 0 ? "opacity-50 cursor-not-allowed" : ""
                }`}
              >
                Prev
              </button>
              {Array.from({ length: totalPages }, (_, index) => (
                <button
                  key={index}
                  onClick={() => handlePageChange(index)}
                  className={`px-3 py-1 rounded-full text-sm font-medium ${
                    currentPage === index
                      ? "bg-[#0F5BA7] text-white dark:bg-[#388BFD]"
                      : "bg-gray-200 text-gray-700 dark:bg-gray-600 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-500"
                  } transition-all duration-300`}
                >
                  {index + 1}
                </button>
              ))}
              <button
                onClick={() => handlePageChange(currentPage + 1)}
                disabled={currentPage === totalPages - 1}
                className={`px-3 py-1 rounded-full text-sm font-medium bg-[#0F5BA7] dark:bg-[#388BFD] text-white hover:from-indigo-600 hover:to-indigo-700 dark:from-indigo-400 dark:to-indigo-500 dark:hover:from-indigo-500 dark:hover:to-indigo-600 transition-all duration-300 ${
                  currentPage === totalPages - 1
                    ? "opacity-50 cursor-not-allowed"
                    : ""
                }`}
              >
                Next
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
};

const Simulations: React.FC = () => {
  const navigate = useNavigate();
  const [filter1, setFilter1] = useState<string>("All Intersections");
  const [filter2, setFilter2] = useState<string>("All Intersections");
  const [page1, setPage1] = useState<number>(0);
  const [page2, setPage2] = useState<number>(0);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const modalType: "simulations" | "optimizations" = "simulations";
  const [simulations, setSimulations] = useState<SimulationData[]>([]);
  const [optimizations, setOptimizations] = useState<SimulationData[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchIntersections = async (): Promise<ApiIntersection[]> => {
    try {
      const res = await fetch(`${API_BASE_URL}/intersections`, {
        headers: { Authorization: `Bearer ${getAuthToken()}` },
      });
      if (!res.ok)
        throw new Error(`Failed to fetch intersections: ${res.statusText}`);
      const data = await res.json();
      return data.intersections || [];
    } catch (err: unknown) {
      console.error("Error fetching intersections:", err);
      throw err;
    }
  };

  const createIntersection = async (intersectionData: {
    name: string;
    traffic_density: string;
    details: {
      address: string;
      city: string;
      province: string;
    };
    default_parameters: {
      green: number;
      yellow: number;
      red: number;
      speed: number;
      seed: number;
      intersection_type: string;
    };
  }) => {
    try {
      const res = await fetch(`${API_BASE_URL}/intersections`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${getAuthToken()}`,
        },
        body: JSON.stringify(intersectionData),
      });
      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.message || "Failed to create intersection");
      }
      return await res.json();
    } catch (err: unknown) {
      console.error("Error creating intersection:", err);
      throw err;
    }
  };

  const runSimulation = async (intersectionId: string) => {
    try {
      const res = await fetch(
        `${API_BASE_URL}/intersections/${intersectionId}/simulate`,
        {
          headers: { Authorization: `Bearer ${getAuthToken()}` },
        },
      );
      if (!res.ok)
        throw new Error(`Failed to run simulation: ${res.statusText}`);
      return await res.json();
    } catch (err: unknown) {
      console.error("Error running simulation:", err);
      throw err;
    }
  };

  const runOptimization = async (intersectionId: string) => {
    try {
      const res = await fetch(
        `${API_BASE_URL}/intersections/${intersectionId}/optimise`,
        {
          method: "POST",
          headers: { Authorization: `Bearer ${getAuthToken()}` },
        },
      );
      if (!res.ok)
        throw new Error(`Failed to run optimization: ${res.statusText}`);
      return await res.json();
    } catch (err: unknown) {
      console.error("Error running optimization:", err);
      throw err;
    }
  };

  const convertToSimulationData = (
    intersections: ApiIntersection[],
  ): { sims: SimulationData[]; opts: SimulationData[] } => {
    //  ADDED: Helper function to format traffic density strings
    const formatTrafficDensity = (density?: string): string => {
      if (!density) return "Medium";
      const lowerCaseDensity = density.toLowerCase();
      if (lowerCaseDensity.includes("high")) return "High";
      if (lowerCaseDensity.includes("low")) return "Low";
      return "Medium";
    };

    const mapApiStatus = (apiStatus?: string): string => {
      switch (apiStatus) {
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

    const sortedIntersections = [...intersections].sort((a, b) => {
      const dateA = a.created_at ? new Date(a.created_at).getTime() : 0;
      const dateB = b.created_at ? new Date(b.created_at).getTime() : 0;
      return dateB - dateA;
    });

    const sims: SimulationData[] = [];
    const opts: SimulationData[] = [];

    sortedIntersections.forEach((intersection, index) => {
      const commonData = {
        id: index + 1,
        backendId: intersection.id,
        intersection: intersection.name,
        trafficDensity: formatTrafficDensity(intersection.traffic_density),
        speed: intersection.default_parameters.simulation_parameters.speed,
        status: mapApiStatus(intersection.status),
      };

      if (intersection.run_count && intersection.run_count > 0) {
        sims.push(commonData);
      }

      if (intersection.status === "optimised") {
        opts.push(commonData);
      }
    });

    return { sims, opts };
  };

  const handleCreateSimulation = async (data: {
    name: string;
    description: string;
    intersections: string[];
  }) => {
    try {
      const allIntersections = await fetchIntersections();
      const selectedIntersection = allIntersections.find(
        (i) => i.name === data.intersections[0],
      );

      let intersectionId: string;

      if (selectedIntersection) {
        intersectionId = selectedIntersection.id;
      } else {
        const newIntersectionData = {
          name: data.intersections[0],
          traffic_density: "medium",
          details: {
            address: data.intersections[0],
            city: "Pretoria",
            province: "Gauteng",
          },
          default_parameters: {
            green: 30,
            yellow: 3,
            red: 27,
            speed: 60,
            seed: Math.floor(Math.random() * 10000000000),
            intersection_type: "trafficlight",
          },
        };
        const newIntersection = await createIntersection(newIntersectionData);
        intersectionId = newIntersection.id;
      }

      if (modalType === "simulations") {
        await runSimulation(intersectionId);
      } else {
        await runOptimization(intersectionId);
      }

      setIsModalOpen(false);
      fetchData();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "An unexpected error occurred");
    }
  };

  const fetchData = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const intersections = await fetchIntersections();
      const { sims, opts } = convertToSimulationData(intersections);
      setSimulations(sims);
      setOptimizations(opts);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "An unexpected error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleViewResults = (backendId: string, intersectionName: string) => {
    navigate("/simulation-results", {
      state: {
        name: `Simulation Results for ${intersectionName}`,
        description: `Viewing simulation results for ${intersectionName}`,
        intersections: [intersectionName],
        intersectionIds: [backendId],
        type: "simulations",
      },
    });
  };

  const filteredSimulations = simulations.filter(
    (sim) =>
      filter1 === "All Intersections" || sim.intersection.includes(filter1),
  );

  const filteredOptimizations = optimizations.filter(
    (opt) =>
      filter2 === "All Intersections" || opt.intersection.includes(filter2),
  );

  return (
    <div className="simulationsBody bg-gray-50 dark:bg-[#0D1117] min-h-screen">
      <Navbar />
      <main className="main-content w-full px-4 py-8 pb-24">
        <div className="max-w-6xl mx-auto">
          <div className="flex justify-between items-center mb-6">
            <h1 className="text-3xl font-bold text-gray-800 dark:text-gray-200">
              Simulations & Optimizations
            </h1>
            <button
              onClick={() => setIsModalOpen(true)}
              className="create-simulation-btn bg-[#0F5BA7] dark:bg-[#238636] hover:bg-blue-700 dark:hover:bg-[#2DA44E] text-white font-medium py-2 px-6 rounded-lg transition-colors shadow-sm"
            >
              Create New
            </button>
          </div>

          {isLoading && (
            <div className="text-center text-gray-500 dark:text-gray-400">
              Loading data...
            </div>
          )}
          {error && <div className="text-center text-red-500">{error}</div>}

          {!isLoading && !error && (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              <div className="simulations-section">
                <div className="flex justify-between items-center mb-4">
                  <h2 className="text-xl font-semibold text-gray-700 dark:text-gray-300">
                    Simulations
                  </h2>
                  <select
                    value={filter1}
                    onChange={(e) => setFilter1(e.target.value)}
                    className="filter-dropdown bg-white dark:bg-[#161B22] border border-gray-300 dark:border-gray-600 rounded-md px-3 py-1 text-sm"
                  >
                    <option>All Intersections</option>
                    {[...new Set(simulations.map((s) => s.intersection))].map(
                      (name) => (
                        <option key={name} value={name}>
                          {name}
                        </option>
                      ),
                    )}
                  </select>
                </div>
                <SimulationTable
                  simulations={filteredSimulations}
                  currentPage={page1}
                  setCurrentPage={setPage1}
                  onViewResults={handleViewResults}
                />
              </div>

              <div className="optimizations-section">
                <div className="flex justify-between items-center mb-4">
                  <h2 className="text-xl font-semibold text-gray-700 dark:text-gray-300">
                    Optimizations
                  </h2>
                  <select
                    value={filter2}
                    onChange={(e) => setFilter2(e.target.value)}
                    className="filter-dropdown bg-white dark:bg-[#161B22] border border-gray-300 dark:border-gray-600 rounded-md px-3 py-1 text-sm"
                  >
                    <option>All Intersections</option>
                    {[...new Set(optimizations.map((o) => o.intersection))].map(
                      (name) => (
                        <option key={name} value={name}>
                          {name}
                        </option>
                      ),
                    )}
                  </select>
                </div>
                <SimulationTable
                  simulations={filteredOptimizations}
                  currentPage={page2}
                  setCurrentPage={setPage2}
                  onViewResults={handleViewResults}
                />
              </div>
            </div>
          )}
        </div>
      </main>
      <Footer />
      <HelpMenu />
      <NewSimulationModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateSimulation}
        intersections={[...new Set(simulations.map((s) => s.intersection))]}
        type={modalType}
      />
    </div>
  );
};

export default Simulations;