import { useState, useEffect } from "react";
import Navbar from "../components/Navbar";
import { Search, X } from "lucide-react";
import IntersectionCard from "../components/IntersectionCard";
import "../styles/Intersections.css";
import Footer from "../components/Footer";
import HelpMenu from "../components/HelpMenu";

export interface IntersectionFormData {
  name: string;
  traffic_density: "low" | "medium" | "high";
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
}
interface CreateIntersectionModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: IntersectionFormData) => void;
  isLoading: boolean;
  error: string | null;
}

const CreateIntersectionModal: React.FC<CreateIntersectionModalProps> = ({
  isOpen,
  onClose,
  onSubmit,
  isLoading,
  error,
}) => {
  const [formData, setFormData] = useState<IntersectionFormData>({
    name: "",
    traffic_density: "low",
    details: {
      address: "",
      city: "Pretoria",
      province: "Gauteng",
    },
    default_parameters: {
      green: 10,
      yellow: 3,
      red: 5,
      speed: 60,
      seed: Math.floor(Math.random() * 10000000000),
      intersection_type: "traffic light",
    },
  });

  if (!isOpen) return null;

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    const keys = name.split('.');

    if (keys.length > 1) {
      const [parentKey, childKey] = keys as [keyof IntersectionFormData, string];
      if (parentKey === 'details' || parentKey === 'default_parameters') {
        setFormData(prev => ({
          ...prev,
          [parentKey]: {
            ...prev[parentKey],
            [childKey]: e.target.type === 'number' ? parseInt(value, 10) : value,
          },
        }));
      }
    } else {
      setFormData(prev => ({
        ...prev,
        [name]: value,
      }));
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-60 flex justify-center items-center z-50">
      <div className="bg-white dark:bg-gray-800 p-8 rounded-lg shadow-xl w-full max-w-2xl relative">
        <button onClick={onClose} className="absolute top-4 right-4 text-gray-500 hover:text-gray-800 dark:hover:text-gray-200">
          <X size={24} />
        </button>
        <h2 className="text-2xl font-bold mb-6 text-center text-gray-800 dark:text-white">Create New Intersection</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="name" className="block text-sm font-medium text-gray-700 dark:text-gray-300">Intersection Name</label>
            <input type="text" name="name" id="name" required className="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm text-black dark:text-white" value={formData.name} onChange={handleChange} />
          </div>
          <div>
            <label htmlFor="details.address" className="block text-sm font-medium text-gray-700 dark:text-gray-300">Address</label>
            <input type="text" name="details.address" id="details.address" required className="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm text-black dark:text-white" value={formData.details.address} onChange={handleChange} />
          </div>
          <div>
              <label htmlFor="traffic_density" className="block text-sm font-medium text-gray-700 dark:text-gray-300">Traffic Density</label>
              <select name="traffic_density" id="traffic_density" required className="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm text-black dark:text-white" value={formData.traffic_density} onChange={handleChange}>
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
              </select>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
                <label htmlFor="default_parameters.green" className="block text-sm font-medium text-gray-700 dark:text-gray-300">Green Light (s)</label>
                <input type="number" name="default_parameters.green" id="default_parameters.green" required className="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm text-black dark:text-white" value={formData.default_parameters.green} onChange={handleChange} />
            </div>
            <div>
                <label htmlFor="default_parameters.yellow" className="block text-sm font-medium text-gray-700 dark:text-gray-300">Yellow Light (s)</label>
                <input type="number" name="default_parameters.yellow" id="default_parameters.yellow" required className="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm text-black dark:text-white" value={formData.default_parameters.yellow} onChange={handleChange} />
            </div>
            <div>
                <label htmlFor="default_parameters.red" className="block text-sm font-medium text-gray-700 dark:text-gray-300">Red Light (s)</label>
                <input type="number" name="default_parameters.red" id="default_parameters.red" required className="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm text-black dark:text-white" value={formData.default_parameters.red} onChange={handleChange} />
            </div>
          </div>
          {error && <p className="text-red-500 text-sm text-center">{error}</p>}
          <div className="flex justify-end space-x-4 pt-4">
            <button type="button" onClick={onClose} className="px-4 py-2 bg-gray-300 text-gray-800 rounded-md hover:bg-gray-400 dark:bg-gray-600 dark:text-white dark:hover:bg-gray-500">Cancel</button>
            <button type="submit" disabled={isLoading} className="px-4 py-2 bg-red-700 text-white rounded-md hover:bg-red-800 disabled:bg-red-400 disabled:cursor-not-allowed">
              {isLoading ? "Creating..." : "Create Intersection"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const API_BASE_URL = "/api";

const getAuthToken = () => {
  return localStorage.getItem('authToken');
};

interface Intersection {
  id: string;
  name: string;
  details: {
    address: string;
    city: string;
    province: string;
  };
  default_parameters: {
    intersection_type: string;
  };
  image?: string;
}

const Intersections = () => {
  const [searchQuery, setSearchQuery] = useState("");
  const [intersections, setIntersections] = useState<Intersection[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [createError, setCreateError] = useState<string | null>(null);

  const fetchIntersections = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_BASE_URL}/intersections`, {
        headers: { 'Authorization': `Bearer ${getAuthToken()}` }
      });
      if (!response.ok) throw new Error(`Failed to fetch intersections: ${response.statusText}`);
      const data = await response.json();
      setIntersections(data.intersections || []);
    } catch (err: unknown) {
      if (err instanceof Error) {
          setError(err.message);
      } else {
          setError("An unexpected error occurred.");
      }
      setIntersections([]);
    } finally {
      setIsLoading(false);
    }
  };
  const searchIntersectionById = async (id: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_BASE_URL}/intersections/${id}`, {
         headers: { 'Authorization': `Bearer ${getAuthToken()}` }
      });
      if (response.status === 404) {
        setIntersections([]);
        return;
      }
      if (!response.ok) throw new Error(`Failed to find intersection with ID ${id}: ${response.statusText}`);
      const data = await response.json();
      setIntersections(data ? [data] : []);
    } catch (err: unknown) {
        if (err instanceof Error) {
            setError(err.message);
        } else {
            setError("An unexpected error occurred while searching.");
        }
      setIntersections([]);
    } finally {
      setIsLoading(false);
    }
  };
  const handleCreateIntersection = async (formData: IntersectionFormData) => {
    setIsCreating(true);
    setCreateError(null);
    try {
      const response = await fetch(`${API_BASE_URL}/intersections`, {
          method: 'POST',
          headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${getAuthToken()}`
          },
          body: JSON.stringify(formData),
      });
      if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || 'Failed to create intersection');
      }
      setIsModalOpen(false);
      fetchIntersections();
    } catch (err: unknown) {
      if (err instanceof Error) {
        setCreateError(err.message);
      } else {
        setCreateError("An unexpected error occurred during creation.");
      }
    } finally {
        setIsCreating(false);
    }
  };
  useEffect(() => {
    fetchIntersections();
  }, []);
  useEffect(() => {
    const handler = setTimeout(() => {
      if (searchQuery.trim() === '') {
        fetchIntersections();
      } else if (!isNaN(Number(searchQuery))) {
        searchIntersectionById(searchQuery);
      }
    }, 500);
    return () => clearTimeout(handler);
  }, [searchQuery]);

  const filteredIntersections = searchQuery && isNaN(Number(searchQuery))
    ? intersections.filter(intersection =>
        intersection.name.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : intersections;

  return (
    <>
      <div className="intersectionBody flex flex-col min-h-screen bg-gray-100">
        <Navbar />
        <div className="main-content flex-grow w-full">
          <div className="max-w-6xl mx-auto w-full px-4 py-8 pb-24">
            <div className="topBar flex justify-between items-center mb-6 gap-x-4">
              <div className="searchContainer relative w-full max-w-md">
                <input
                  type="text"
                  placeholder="Search by Name or ID..."
                  className="searchBar w-full pl-4 pr-10 py-2 border border-gray-300 rounded-full focus:outline-none focus:ring-2 focus:ring-red-500"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
                <div className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500">
                  <Search size={20} />
                </div>
              </div>
              <button
                onClick={() => setIsModalOpen(true)}
                className="addIntersectionBtn flex-shrink-0 bg-red-700 hover:bg-red-800 text-white font-medium py-2 px-4 rounded-md"
              >
                Add Intersection
              </button>
            </div>
            <div className="intersections space-y-6 pr-2">
              {isLoading ? (
                <p className="text-center text-gray-500 dark:text-gray-400">Loading intersections...</p>
              ) : error ? (
                <p className="text-center text-red-500">{error}</p>
              ) : filteredIntersections.length > 0 ? (
                filteredIntersections.map((intersection) => (
                  <IntersectionCard
                    key={intersection.id}
                    id={intersection.id}
                    name={intersection.name}
                    location={`${intersection.details.address}`}
                    lanes={intersection.default_parameters.intersection_type}
                    onSimulate={(id) => console.log(`Simulate ${id}`)}
                    onEdit={(id) => console.log(`Edit ${id}`)}
                    onDelete={(id) => console.log(`Delete ${id}`)}
                  />
                ))
              ) : (
                <p className="text-center text-gray-500 dark:text-gray-400">No intersections found.</p>
              )}
            </div>
          </div>
        </div>
        <Footer />
        <HelpMenu />
      </div>

      <CreateIntersectionModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateIntersection}
        isLoading={isCreating}
        error={createError}
      />
    </>
  );
};

export default Intersections;
