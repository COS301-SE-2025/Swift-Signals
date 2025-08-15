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
  initialData?: IntersectionFormData | null;
  isEditing: boolean;
}

interface DeleteConfirmationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  intersectionName: string;
  isLoading: boolean;
}

const DeleteConfirmationModal: React.FC<DeleteConfirmationModalProps> = ({
  isOpen,
  onClose,
  onConfirm,
  intersectionName,
  isLoading,
}) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm flex justify-center items-center z-50 p-4">
      <div className="bg-white dark:bg-[#161B22] p-8 rounded-xl shadow-2xl w-full max-w-lg relative border border-gray-200 dark:border-[#30363D] transform transition-all duration-200 scale-100">
        <button
          onClick={onClose}
          disabled={isLoading}
          className="absolute top-5 right-5 text-gray-400 dark:text-[#7D8590] hover:text-gray-600 dark:hover:text-[#E6EDF3] transition-colors duration-150 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <X size={24} />
        </button>
        
        <div className="text-center mb-8">
          <div className="mx-auto mb-4 w-16 h-16 bg-red-100 dark:bg-red-900/20 rounded-full flex items-center justify-center">
            <svg
              className="w-8 h-8 text-red-600 dark:text-red-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
              />
            </svg>
          </div>
          
          <h2 className="text-2xl font-bold mb-3 text-gray-900 dark:text-[#E6EDF3]">
            Delete Intersection?
          </h2>
          
          <div className="space-y-2">
            <p className="text-gray-700 dark:text-[#C9D1D9]">
              You're about to permanently delete
            </p>
            <p className="text-lg font-semibold text-gray-900 dark:text-[#E6EDF3] bg-gray-50 dark:bg-[#21262D] px-3 py-2 rounded-lg border border-gray-200 dark:border-[#30363D]">
              "{intersectionName}"
            </p>
            <p className="text-sm text-red-600 dark:text-red-400 font-medium mt-3">
              ⚠️ This action cannot be undone
            </p>
          </div>
        </div>
        
        <div className="flex flex-col sm:flex-row gap-3 sm:gap-4">
          <button
            type="button"
            onClick={onClose}
            disabled={isLoading}
            className="flex-1 px-6 py-3 bg-gray-100 dark:bg-[#21262D] border-2 border-gray-300 dark:border-[#30363D] text-gray-700 dark:text-[#C9D1D9] rounded-lg font-medium hover:bg-gray-200 dark:hover:bg-[#30363D] focus:outline-none focus:ring-2 focus:ring-gray-300 dark:focus:ring-[#30363D] transition-colors duration-150 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={onConfirm}
            disabled={isLoading}
            className="flex-1 px-6 py-3 bg-red-600 dark:bg-[#DA3633] text-white rounded-lg font-medium hover:bg-red-700 dark:hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500 dark:focus:ring-red-400 transition-colors duration-150 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            {isLoading ? (
              <>
                <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Deleting...
              </>
            ) : (
              "Delete Intersection"
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

const CreateIntersectionModal: React.FC<CreateIntersectionModalProps> = ({
  isOpen,
  onClose,
  onSubmit,
  isLoading,
  error,
  initialData,
  isEditing,
}) => {
  const getDefaultFormData = (): IntersectionFormData => ({
    name: "",
    traffic_density: "low",
    details: { address: "", city: "Pretoria", province: "Gauteng" },
    default_parameters: {
      green: 10,
      yellow: 3,
      red: 5,
      speed: 60,
      seed: Math.floor(Math.random() * 10000000000),
      intersection_type: "traffic light",
    },
  });

  const [formData, setFormData] = useState<IntersectionFormData>(getDefaultFormData());

  useEffect(() => {
    if (isEditing && initialData) {
      setFormData(initialData);
    } else if (!isEditing) {
      setFormData(getDefaultFormData());
    }
  }, [initialData, isEditing, isOpen]);

  if (!isOpen) return null;

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    const keys = name.split(".");

    if (keys.length > 1) {
      const [parentKey, childKey] = keys as [
        keyof IntersectionFormData,
        string,
      ];
      if (parentKey === "details" || parentKey === "default_parameters") {
        setFormData((prev) => ({
          ...prev,
          [parentKey]: {
            ...prev[parentKey],
            [childKey]:
              e.target.type === "number" ? parseInt(value, 10) : value,
          },
        }));
      }
    } else {
      setFormData((prev) => ({
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
    <div className="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50 p-4">
      <div className="bg-white dark:bg-[#161B22] p-4 sm:p-8 rounded-lg shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto relative">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-500 dark:text-[#E6EDF3] hover:text-gray-800 dark:hover:text-gray-200"
        >
          <X size={24} />
        </button>
        <h2 className="text-2xl font-bold mb-6 text-center text-gray-800 dark:text-[#E6EDF3]">
          {isEditing ? "Edit Intersection" : "Create New Intersection"}
        </h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label
              htmlFor="name"
              className="block text-sm font-medium text-gray-700 dark:text-[#E6EDF3]"
            >
              Intersection Name
            </label>
            <input
              type="text"
              name="name"
              id="name"
              required
              className="mt-1 block w-full px-3 py-2 bg-white dark:bg-[#161B22] border-2 border-gray-300 dark:border-[#30363D] rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm text-black dark:text-white"
              value={formData.name}
              onChange={handleChange}
            />
          </div>
          <div>
            <label
              htmlFor="details.address"
              className="block text-sm font-medium text-gray-700 dark:text-[#E6EDF3]"
            >
              Address
            </label>
            <input
              type="text"
              name="details.address"
              id="details.address"
              required
              className="mt-1 block w-full px-3 py-2 bg-white dark:bg-[#161B22] border-2 border-gray-300 dark:border-[#30363D] rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm text-black dark:text-white"
              value={formData.details.address}
              onChange={handleChange}
            />
          </div>
          <div>
            <label
              htmlFor="traffic_density"
              className="block text-sm font-medium text-gray-700 dark:text-[#E6EDF3]"
            >
              Traffic Density
            </label>
            <select
              name="traffic_density"
              id="traffic_density"
              required
              className="mt-1 block w-full px-3 py-2 bg-white dark:bg-[#161B22] border-2 border-gray-300 dark:border-[#30363D] rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm text-black dark:text-white"
              value={formData.traffic_density}
              onChange={handleChange}
            >
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
            </select>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label
                htmlFor="default_parameters.green"
                className="block text-sm font-medium text-gray-700 dark:text-[#E6EDF3]"
              >
                Green Light (s)
              </label>
              <input
                type="number"
                name="default_parameters.green"
                id="default_parameters.green"
                required
                className="mt-1 block w-full px-3 py-2 bg-white dark:bg-[#161B22] border-2 border-gray-300 dark:border-[#30363D] rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm text-black dark:text-white"
                value={formData.default_parameters.green}
                onChange={handleChange}
              />
            </div>
            <div>
              <label
                htmlFor="default_parameters.yellow"
                className="block text-sm font-medium text-gray-700 dark:text-[#E6EDF3]"
              >
                Yellow Light (s)
              </label>
              <input
                type="number"
                name="default_parameters.yellow"
                id="default_parameters.yellow"
                required
                className="mt-1 block w-full px-3 py-2 bg-white dark:bg-[#161B22] border-2 border-gray-300 dark:border-[#30363D] rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm text-black dark:text-white"
                value={formData.default_parameters.yellow}
                onChange={handleChange}
              />
            </div>
            <div>
              <label
                htmlFor="default_parameters.red"
                className="block text-sm font-medium text-gray-700 dark:text-[#E6EDF3]"
              >
                Red Light (s)
              </label>
              <input
                type="number"
                name="default_parameters.red"
                id="default_parameters.red"
                required
                className="mt-1 block w-full px-3 py-2 bg-white dark:bg-[#161B22] border-2 border-gray-300 dark:border-[#30363D] rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm text-black dark:text-white"
                value={formData.default_parameters.red}
                onChange={handleChange}
              />
            </div>
          </div>
          {error && <p className="text-red-500 text-sm text-center">{error}</p>}
          <div className="flex justify-end space-x-4 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-300 dark:bg-[#161B22] dark:border-2 dark:border-[#DA3633] text-gray-800 rounded-md hover:bg-gray-400 dark:text-white dark:hover:bg-[#DA3633]"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className="px-4 py-2 bg-[#0F5BA7] dark:bg-[#388BFD] text-white rounded-md hover:bg-red-800 disabled:bg-red-400 disabled:cursor-not-allowed"
            >
              {isLoading
                ? isEditing
                  ? "Updating..."
                  : "Creating..."
                : isEditing
                  ? "Update Intersection"
                  : "Create Intersection"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const API_BASE_URL = "http://localhost:9090";

const getAuthToken = () => {
  return localStorage.getItem("authToken");
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

  const [isEditing, setIsEditing] = useState(false);
  const [selectedIntersectionId, setSelectedIntersectionId] = useState<
    string | null
  >(null);
  const [editData, setEditData] = useState<IntersectionFormData | null>(null);

  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [intersectionToDelete, setIntersectionToDelete] = useState<{
    id: string;
    name: string;
  } | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  const fetchIntersections = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await fetch(`${API_BASE_URL}/intersections`, {
        headers: { Authorization: `Bearer ${getAuthToken()}` },
      });
      if (!res.ok)
        throw new Error(`Failed to fetch intersections: ${res.statusText}`);
      const data = await res.json();
      setIntersections(data.intersections || []);
    } catch (err: any) {
      setError(err.message || "Unexpected error");
    } finally {
      setIsLoading(false);
    }
  };

  const searchIntersectionById = async (id: string) => {
    setIsLoading(true);
    try {
      const res = await fetch(`${API_BASE_URL}/intersections/${id}`, {
        headers: { Authorization: `Bearer ${getAuthToken()}` },
      });
      if (res.status === 404) {
        setIntersections([]);
        return;
      }
      if (!res.ok)
        throw new Error(`Failed to find intersection: ${res.statusText}`);
      const data = await res.json();
      setIntersections(data ? [data] : []);
    } catch (err: any) {
      setError(err.message || "Unexpected error");
      setIntersections([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateIntersection = async (formData: IntersectionFormData) => {
    setIsCreating(true);
    setCreateError(null);
    try {
      const res = await fetch(`${API_BASE_URL}/intersections`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${getAuthToken()}`,
        },
        body: JSON.stringify(formData),
      });
      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.message || "Failed to create intersection");
      }
      setIsModalOpen(false);
      fetchIntersections();
    } catch (err: any) {
      setCreateError(err.message || "Unexpected error");
    } finally {
      setIsCreating(false);
    }
  };

  const handleUpdateIntersection = async (formData: IntersectionFormData) => {
    if (!selectedIntersectionId) return;
    setIsCreating(true);
    setCreateError(null);
    try {
      const res = await fetch(
        `${API_BASE_URL}/intersections/${selectedIntersectionId}`,
        {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${getAuthToken()}`,
          },
          body: JSON.stringify({
            name: formData.name,
            details: formData.details,
          }),
        }
      );
      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.message || "Failed to update intersection");
      }
      setIsModalOpen(false);
      setIsEditing(false);
      setSelectedIntersectionId(null);
      fetchIntersections();
    } catch (err: any) {
      setCreateError(err.message || "Unexpected error");
    } finally {
      setIsCreating(false);
    }
  };

  const handleDeleteClick = (id: string) => {
    const intersection = intersections.find((i) => i.id === id);
    if (!intersection) return;
    
    setIntersectionToDelete({
      id: intersection.id,
      name: intersection.name,
    });
    setIsDeleteModalOpen(true);
  };

  const handleDeleteConfirm = async () => {
    if (!intersectionToDelete) return;
    
    setIsDeleting(true);
    try {
      const res = await fetch(
        `${API_BASE_URL}/intersections/${intersectionToDelete.id}`,
        {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${getAuthToken()}`,
          },
        }
      );
      
      if (!res.ok) {
        const errorData = await res.json().catch(() => ({}));
        throw new Error(errorData.message || "Failed to delete intersection");
      }

      setIsDeleteModalOpen(false);
      setIntersectionToDelete(null);

      if (searchQuery.trim() === intersectionToDelete.id) {
        setSearchQuery("");
      }
      
      fetchIntersections();
    } catch (err: any) {
      setError(err.message || "Failed to delete intersection");
    } finally {
      setIsDeleting(false);
    }
  };

  const handleEditClick = (id: string) => {
    const intersection = intersections.find((i) => i.id === id);
    if (!intersection) return;
    setEditData({
      name: intersection.name,
      traffic_density: "low",
      details: intersection.details,
      default_parameters: {
        green: 10,
        yellow: 3,
        red: 5,
        speed: 60,
        seed: 1,
        intersection_type: intersection.default_parameters.intersection_type,
      },
    });
    setSelectedIntersectionId(id);
    setIsEditing(true);
    setIsModalOpen(true);
  };

  useEffect(() => {
    fetchIntersections();
  }, []);

  useEffect(() => {
    const handler = setTimeout(() => {
      if (searchQuery.trim() === "") {
        fetchIntersections();
      } else if (!isNaN(Number(searchQuery))) {
        searchIntersectionById(searchQuery);
      }
    }, 500);
    return () => clearTimeout(handler);
  }, [searchQuery]);

  const filteredIntersections =
    searchQuery && isNaN(Number(searchQuery))
      ? intersections.filter((intersection) =>
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
                  placeholder="Search by Name..."
                  className="searchBar w-full pl-4 pr-10 py-2 border-2 border-gray-300 rounded-full focus:outline-none focus:ring-2 focus:ring-red-500"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
                <div className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500">
                  <Search size={20} />
                </div>
              </div>
              <button
                onClick={() => {
                  setIsEditing(false);
                  setEditData(null);
                  setSelectedIntersectionId(null);
                  setCreateError(null);
                  setIsModalOpen(true);
                }}
                className="addIntersectionBtn flex-shrink-0 bg-[#0F5BA7] dark:bg-[#388BFD] hover:bg-[#3DAEF0] text-white font-medium py-2 px-4 rounded-md"
              >
                Add Intersection
              </button>
            </div>
            <div className="intersections space-y-6 pr-2">
              {isLoading ? (
                <p className="text-center text-gray-500 dark:text-gray-400">
                  Loading intersections...
                </p>
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
                    onEdit={handleEditClick}
                    onDelete={handleDeleteClick}
                  />
                ))
              ) : (
                <p className="text-center text-gray-500 dark:text-gray-400">
                  No intersections found.
                </p>
              )}
            </div>
          </div>
        </div>
        <Footer />
        <HelpMenu />
      </div>

      <CreateIntersectionModal
        isOpen={isModalOpen}
        onClose={() => {
          setIsModalOpen(false);
          setCreateError(null);
          if (!isEditing) {
            setEditData(null);
          }
        }}
        onSubmit={
          isEditing ? handleUpdateIntersection : handleCreateIntersection
        }
        isLoading={isCreating}
        error={createError}
        initialData={editData}
        isEditing={isEditing}
      />

      <DeleteConfirmationModal
        isOpen={isDeleteModalOpen}
        onClose={() => {
          setIsDeleteModalOpen(false);
          setIntersectionToDelete(null);
        }}
        onConfirm={handleDeleteConfirm}
        intersectionName={intersectionToDelete?.name || ""}
        isLoading={isDeleting}
      />
    </>
  );
};

export default Intersections;