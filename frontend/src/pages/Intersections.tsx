import type { LatLng } from "leaflet";
import L from "leaflet";
import { Search, X, FileText, MapPin, TrafficCone } from "lucide-react";
import { useState, useEffect, useRef } from "react";
import {
  MapContainer,
  TileLayer,
  Marker,
  useMapEvents,
  Popup,
} from "react-leaflet";
import { useNavigate } from "react-router-dom";

import Footer from "../components/Footer";
import HelpMenu from "../components/HelpMenu";
import IntersectionCard from "../components/IntersectionCard";
import Navbar from "../components/Navbar";
import "../styles/Intersections.css";

import "leaflet/dist/leaflet.css";
import { API_BASE_URL } from "../config";
import { CHATBOT_BASE_URL } from "../config";

// =================================================================
// DATA STRUCTURES & INTERFACES
// =================================================================

type TrafficDensityUI = "low" | "medium" | "high";

export interface IntersectionFormData {
  name: string;
  traffic_density: TrafficDensityUI;
  details: {
    address: string;
    city: string;
    province: string;
    latitude?: number;
    longitude?: number;
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

interface SimulationParameters {
  intersection_type: string;
  green: number;
  yellow: number;
  red: number;
  speed: number;
  seed: number;
}

interface OptimisationParameters {
  optimisation_type: string;
  simulation_parameters: SimulationParameters;
}

interface Intersection {
  id: string;
  name: string;
  details: {
    address: string;
    city: string;
    province: string;
  };
  default_parameters: OptimisationParameters;
  traffic_density: string;
  image?: string;
}

// #region API and Data Types
// Type for Overpass API elements
interface OverpassElement {
  type: "node" | "way" | "relation";
  id: number;
  lat?: number;
  lon?: number;
  tags?: {
    [key: string]: string | undefined;
    name?: string;
    highway?: string;
    "addr:city"?: string;
    "addr:state"?: string;
    "addr:province"?: string;
  };
  nodes?: number[];
  bounds?: {
    minlat: number;
    minlon: number;
    maxlat: number;
    maxlon: number;
  };
}

// Type for intersection with calculated distance
type IntersectionWithDistance = Intersection & {
  distance: number;
  intersection: string;
  lat: number;
  lon: number;
};
// #endregion

// =================================================================
// MAPPING & MODAL COMPONENTS
// =================================================================

const LocationMarker: React.FC<{
  setSelectedLocation: (location: {
    address: string;
    city: string;
    province: string;
    lat: number;
    lng: number;
  }) => void;
  setCoordinates: (coords: string) => void;
  setIsSnapping?: (snapping: boolean) => void;
}> = ({ setSelectedLocation, setCoordinates, setIsSnapping }) => {
  const [position, setPosition] = useState<LatLng | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [snappedAddress, setSnappedAddress] = useState<{
    address: string;
    city: string;
    province: string;
  } | null>(null);
  const markerRef = useRef<L.Marker>(null);

  useEffect(() => {
    if (markerRef.current && snappedAddress) {
      markerRef.current.openPopup();
    }
  }, [snappedAddress]);

  // Function to find nearest intersection using Overpass API
  const findNearestIntersection = async (
    lat: number,
    lon: number,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ): Promise<any | null> => {
    try {
      setIsSnapping?.(true);
      setIsLoading(true);

      const radius = 0.005;
      const bbox = `${lat - radius},${lon - radius},${lat + radius},${
        lon + radius
      }`;

      const overpassQuery = `
    [out:json][timeout:30];
    (
     way["highway"~"^(primary|secondary|tertiary|residential|trunk|motorway|unclassified)$"](${bbox});
    );
    (._;>;);
    out geom;
   `;

      const overpassUrl = "https://overpass-api.de/api/interpreter";
      const response = await fetch(overpassUrl, {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: `data=${encodeURIComponent(overpassQuery)}`,
      });

      if (!response.ok) {
        console.warn("Overpass API failed:", response.status);
        return null;
      }

      const data = await response.json();

      if (!data.elements || data.elements.length === 0) {
        return null;
      }

      const ways = data.elements.filter(
        (el: OverpassElement) => el.type === "way" && el.tags?.name,
      );
      const nodes = data.elements.filter(
        (el: OverpassElement) => el.type === "node",
      );

      if (ways.length < 2) {
        return null;
      }

      const nodeMap = new Map<number, { lat: number; lon: number }>();
      nodes.forEach((node: OverpassElement) => {
        if (node.id && node.lat && node.lon) {
          nodeMap.set(node.id, { lat: node.lat, lon: node.lon });
        }
      });

      const nodeWaysMap = new Map<number, OverpassElement[]>();
      ways.forEach((way: OverpassElement) => {
        if (way.nodes) {
          way.nodes.forEach((nodeId: number) => {
            if (!nodeWaysMap.has(nodeId)) {
              nodeWaysMap.set(nodeId, []);
            }
            nodeWaysMap.get(nodeId)!.push(way);
          });
        }
      });

      const intersections: IntersectionWithDistance[] = [];

      for (const [nodeId, nodeWays] of nodeWaysMap.entries()) {
        if (nodeWays.length >= 2) {
          const uniqueRoads = [
            ...new Set(
              nodeWays
                .map((way) => way.tags?.name)
                .filter((name): name is string => !!name),
            ),
          ];

          if (uniqueRoads.length >= 2) {
            const nodeCoords = nodeMap.get(nodeId);
            if (nodeCoords) {
              const distance = Math.sqrt(
                Math.pow(nodeCoords.lat - lat, 2) +
                  Math.pow(nodeCoords.lon - lon, 2),
              );

              intersections.push({
                lat: nodeCoords.lat,
                lon: nodeCoords.lon,
                roads: uniqueRoads.slice(0, 2),
                intersection: `${uniqueRoads[0]} & ${uniqueRoads[1]}`,
                distance,
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
              } as any);
            }
          }
        }
      }

      intersections.sort((a, b) => a.distance - b.distance);

      if (intersections.length > 0) {
        const closestIntersection = intersections[0];
        const reverseGeocodeUrl = `${CHATBOT_BASE_URL}/reverse-geocode?lat=${closestIntersection.lat}&lon=${closestIntersection.lon}`;
        const reverseGeocodeResponse = await fetch(reverseGeocodeUrl);
        const reverseGeocodeData = await reverseGeocodeResponse.json();
        const address = reverseGeocodeData.address;
        const streetName = closestIntersection.intersection;
        const city = address.city || address.town || "";
        const province = address.state || address.province || "";

        return {
          ...closestIntersection,
          address: streetName,
          city: city,
          province: province,
        };
      }

      return null;
    } catch (error) {
      console.error("Error finding nearest intersection:", error);
      return null;
    } finally {
      setIsLoading(false);
      setIsSnapping?.(false);
    }
  };

  useMapEvents({
    async click(e) {
      console.log("Map clicked at:", e.latlng);

      setPosition(e.latlng);
      setSnappedAddress(null);

      try {
        const nearestIntersection = await findNearestIntersection(
          e.latlng.lat,
          e.latlng.lng,
        );

        if (nearestIntersection) {
          const snappedPosition = {
            lat: nearestIntersection.lat,
            lng: nearestIntersection.lon,
          } as LatLng;

          setPosition(snappedPosition);
          const newAddress = {
            address: nearestIntersection.address,
            city: nearestIntersection.city,
            province: nearestIntersection.province,
            lat: nearestIntersection.lat,
            lng: nearestIntersection.lon,
          };
          setSelectedLocation(newAddress);
          setSnappedAddress(newAddress);
          setCoordinates(
            `${nearestIntersection.lat.toFixed(
              6,
            )}, ${nearestIntersection.lon.toFixed(6)}`,
          );

          console.log(
            "Snapped to intersection:",
            nearestIntersection.intersection,
          );
        } else {
          const coordinates = `${e.latlng.lat.toFixed(
            6,
          )}, ${e.latlng.lng.toFixed(6)}`;
          setSelectedLocation({
            address: coordinates,
            city: "",
            province: "",
            lat: e.latlng.lat,
            lng: e.latlng.lng,
          });
          setCoordinates(coordinates);

          console.log("No intersection found, using clicked coordinates");
        }
      } catch (error) {
        console.error("Error processing map click:", error);
        const coordinates = `${e.latlng.lat.toFixed(
          6,
        )}, ${e.latlng.lng.toFixed(6)}`;
        setSelectedLocation({
          address: coordinates,
          city: "",
          province: "",
          lat: e.latlng.lat,
          lng: e.latlng.lng,
        });
        setCoordinates(coordinates);
      }
    },
  });

  return (
    <>
      {position && (
        <Marker ref={markerRef} position={position}>
          {snappedAddress && (
            <Popup>
              <b>{snappedAddress.address}</b>
              <br />
              {snappedAddress.city}, {snappedAddress.province}
            </Popup>
          )}
        </Marker>
      )}
      {isLoading && (
        <div className="absolute top-2 left-2 bg-white dark:bg-gray-800 px-3 py-1 rounded-md shadow-md z-10">
          <div className="flex items-center space-x-2">
            <div className="animate-spin inline-block w-4 h-4 border-2 border-current border-t-transparent rounded-full"></div>
            <span className="text-sm text-gray-700 dark:text-gray-300">
              Finding nearest intersection...
            </span>
          </div>
        </div>
      )}
    </>
  );
};

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
              You&apos;re about to permanently delete
            </p>
            <p className="text-lg font-semibold text-gray-900 dark:text-[#E6EDF3] bg-gray-50 dark:bg-[#21262D] px-3 py-2 rounded-lg border border-gray-200 dark:border-[#30363D]">
              &quot;{intersectionName}&quot;
            </p>
            <p className="text-sm text-red-600 dark:text-red-400 font-medium mt-3">
              This action cannot be undone
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
                <svg
                  className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  ></circle>
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
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

interface CreateIntersectionModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (formData: IntersectionFormData) => void;
  isLoading: boolean;
  error: string | null;
  initialData?: IntersectionFormData | null;
  isEditing: boolean;
}

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
    traffic_density: "medium",
    details: { address: "", city: "Pretoria", province: "Gauteng" },
    default_parameters: {
      green: 30,
      yellow: 3,
      red: 27,
      speed: 60,
      seed: Math.floor(Math.random() * 10000000000),
      intersection_type: "INTERSECTION_TYPE_TRAFFICLIGHT",
    },
  });

  const [formData, setFormData] =
    useState<IntersectionFormData>(getDefaultFormData());
  const [activeTab, setActiveTab] = useState<"Manual" | "Map">("Manual");
  const [coordinates, setCoordinates] = useState<string | null>(null);
  const [isSnapping, setIsSnapping] = useState(false);

  useEffect(() => {
    if (isOpen) {
      if (isEditing && initialData) {
        setFormData(initialData);
      } else {
        setFormData(getDefaultFormData());
      }
    }
  }, [initialData, isEditing, isOpen]);

  if (!isOpen) return null;

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;
    const keys = name.split(".");

    setFormData((prev) => {
      const newFormData = JSON.parse(JSON.stringify(prev));

      if (keys.length > 1) {
        let currentLevel: Record<string, unknown> = newFormData;
        for (let i = 0; i < keys.length - 1; i++) {
          currentLevel = currentLevel[keys[i]] as Record<string, unknown>;
        }
        const finalKey = keys[keys.length - 1];
        const isNumber = ["green", "yellow", "red", "speed", "seed"].includes(
          finalKey,
        );
        (currentLevel as Record<string, unknown>)[finalKey] = isNumber
          ? parseInt(value, 10) || 0
          : value;
      } else {
        (newFormData as Record<string, unknown>)[name] = value;
      }

      return newFormData;
    });
  };

  const handleMapSelection = (location: {
    address: string;
    city: string;
    province: string;
    lat: number;
    lng: number;
  }) => {
    setFormData((prev) => ({
      ...prev,
      name: location.address,
      details: {
        ...prev.details,
        address: `${location.address}, ${location.city}, ${location.province}`,
        city: location.city,
        province: location.province,
        latitude: location.lat,
        longitude: location.lng,
      },
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const newFormData = { ...formData };
    if (newFormData.details.latitude && newFormData.details.longitude) {
      newFormData.name = `${newFormData.name} [${newFormData.details.latitude},${newFormData.details.longitude}]`;
    }
    onSubmit(newFormData);
  };

  const inputClasses =
    "mt-1 block w-full px-3 py-2 bg-gray-50 dark:bg-[#0D1117] border-2 border-gray-300 dark:border-[#30363D] rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 dark:focus:ring-offset-[#161B22] focus:ring-[#2da44e] sm:text-sm text-gray-900 dark:text-[#C9D1D9]";

  const TrafficDensityButton = ({
    value,
    label,
  }: {
    value: TrafficDensityUI;
    label: string;
  }) => (
    <button
      type="button"
      onClick={() =>
        setFormData((prev) => ({ ...prev, traffic_density: value }))
      }
      className={`flex-1 px-3 py-2 text-sm font-medium rounded-md transition-all duration-150 focus:outline-none focus:ring-2 focus:ring-offset-2 dark:focus:ring-offset-[#161B22] focus:ring-[#2da44e] ${
        formData.traffic_density === value
          ? "bg-[#2da44e] text-white shadow-md"
          : "bg-gray-200 dark:bg-[#21262D] text-gray-700 dark:text-[#C9D1D9] hover:bg-gray-300 dark:hover:bg-[#30363D]"
      }`}
    >
      {label}
    </button>
  );

  return (
    <div className="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm flex justify-center items-center z-50 px-2 py-8">
      <div className="bg-white dark:bg-[#161B22] p-4 sm:p-8 md:p-6 rounded-xl shadow-2xl w-full sm:max-w-xl md:max-w-2xl lg:max-w-4xl create-intersection-modal relative border border-gray-200 dark:border-[#30363D] flex flex-col max-h-[90vh] overflow-y-auto [\@media(min-width:1025px)_and_\(max-width:1400px\)]\:max-w-\[850px\] [\@media(min-width:769px)_and_\(max-width:1024px\)]\:max-w-\[800px\]">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-400 dark:text-[#7D8590] hover:text-gray-600 dark:hover:text-[#E6EDF3] transition-colors duration-150"
        >
          <X size={24} />
        </button>
        <h2 className="text-3xl font-bold mb-8 text-center text-gray-900 dark:text-[#E6EDF3]">
          {isEditing ? "Edit Intersection" : "Create New Intersection"}
        </h2>
        <form
          onSubmit={handleSubmit}
          className="flex flex-col flex-grow overflow-y-auto"
        >
          <div
            className={`grid grid-cols-1 lg:grid-cols-2 gap-x-12 ${activeTab === "Map" ? "gap-y-4" : "gap-y-8"}`}
          >
            <div className="space-y-8">
              <div className="space-y-5">
                <h3 className="flex items-center gap-3 text-xl font-semibold text-gray-800 dark:text-[#E6EDF3] border-b border-gray-200 dark:border-[#30363D] pb-3">
                  <FileText size={20} className="text-[#0f5ba7]" />
                  General Information
                </h3>
                <div>
                  <label
                    htmlFor="name"
                    className="block text-sm font-medium text-gray-700 dark:text-[#C9D1D9] mb-1"
                  >
                    Intersection Name
                  </label>
                  <input
                    type="text"
                    name="name"
                    id="name"
                    required
                    className={inputClasses}
                    value={formData.name}
                    onChange={handleChange}
                    placeholder="e.g., Lynnwood & Atterbury"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-[#C9D1D9] mb-2">
                    Traffic Density
                  </label>
                  <div className="flex space-x-2 bg-gray-100 dark:bg-[#0D1117] p-1 rounded-lg">
                    <TrafficDensityButton value="low" label="Low" />
                    <TrafficDensityButton value="medium" label="Medium" />
                    <TrafficDensityButton value="high" label="High" />
                  </div>
                </div>
              </div>
              <div className="space-y-5">
                <h3 className="flex items-center gap-3 text-xl font-semibold text-gray-800 dark:text-[#E6EDF3] border-b border-gray-200 dark:border-[#30363D] pb-3">
                  <MapPin size={20} className="text-[#0f5ba7]" />
                  Location Details
                </h3>
                <div className="flex space-x-2 mb-3">
                  <button
                    type="button"
                    onClick={() => setActiveTab("Manual")}
                    className={`px-3 py-1 rounded-md text-sm font-medium ${
                      activeTab === "Manual"
                        ? "bg-[#2B9348] text-white dark:bg-[#2DA44E]"
                        : "bg-gray-200 text-gray-700 dark:bg-gray-600 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-500"
                    } transition-all duration-300`}
                  >
                    Manual
                  </button>
                  <button
                    type="button"
                    onClick={() => setActiveTab("Map")}
                    className={`px-3 py-1 rounded-md text-sm font-medium ${
                      activeTab === "Map"
                        ? "bg-[#2B9348] text-white dark:bg-[#2DA44E]"
                        : "bg-gray-200 text-gray-700 dark:bg-gray-600 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-500"
                    } transition-all duration-300`}
                  >
                    Map
                  </button>
                </div>
                {activeTab === "Manual" && (
                  <>
                    <div>
                      <label
                        htmlFor="details.address"
                        className="block text-sm font-medium text-gray-700 dark:text-[#C9D1D9] mb-1"
                      >
                        Address / Cross Streets
                      </label>
                      <input
                        type="text"
                        name="details.address"
                        id="details.address"
                        required
                        className={inputClasses}
                        value={formData.details.address}
                        onChange={handleChange}
                        placeholder="Corner of Lynnwood Rd and Atterbury Rd"
                      />
                    </div>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                      <div>
                        <label
                          htmlFor="details.city"
                          className="block text-sm font-medium text-gray-700 dark:text-[#C9D1D9] mb-1"
                        >
                          City
                        </label>
                        <input
                          type="text"
                          name="details.city"
                          id="details.city"
                          required
                          className={inputClasses}
                          value={formData.details.city}
                          onChange={handleChange}
                        />
                      </div>
                      <div>
                        <label
                          htmlFor="details.province"
                          className="block text-sm font-medium text-gray-700 dark:text-[#C9D1D9] mb-1"
                        >
                          Province
                        </label>
                        <input
                          type="text"
                          name="details.province"
                          id="details.province"
                          required
                          className={inputClasses}
                          value={formData.details.province}
                          onChange={handleChange}
                        />
                      </div>
                    </div>
                  </>
                )}
                {activeTab === "Map" && (
                  <div className="relative">
                    <MapContainer
                      center={[-25.7479, 28.2293]}
                      zoom={12}
                      style={{ height: "180px", width: "100%" }}
                    >
                      <TileLayer
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                        attribution=' <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                      />
                      <LocationMarker
                        setSelectedLocation={handleMapSelection}
                        setCoordinates={setCoordinates}
                        setIsSnapping={setIsSnapping}
                      />
                    </MapContainer>
                    <div className="mt-2 space-y-1">
                      {isSnapping && (
                        <div className="flex items-center text-sm text-blue-600 dark:text-blue-400">
                          <div className="animate-spin inline-block w-3 h-3 border-2 border-current border-t-transparent rounded-full mr-2"></div>
                          Snapping to nearest intersection...
                        </div>
                      )}
                      {coordinates && (
                        <p className="text-sm text-gray-700 dark:text-gray-300">
                          <span className="font-medium">Coordinates:</span>{" "}
                          {coordinates}
                        </p>
                      )}
                      <p className="text-xs text-gray-500 dark:text-gray-400">
                        Click anywhere on the map to automatically find the
                        nearest road intersection
                      </p>
                    </div>
                  </div>
                )}
              </div>
            </div>
            <div className="space-y-8">
              <div className="space-y-5">
                <h3 className="flex items-center gap-3 text-xl font-semibold text-gray-800 dark:text-[#E6EDF3] border-b border-gray-200 dark:border-[#30363D] pb-3">
                  <TrafficCone size={20} className="text-[#0f5ba7]" />
                  Simulation Parameters
                </h3>
                <div>
                  <label
                    htmlFor="default_parameters.intersection_type"
                    className="block text-sm font-medium text-gray-700 dark:text-[#C9D1D9] mb-1"
                  >
                    Intersection Type
                  </label>
                  <select
                    name="default_parameters.intersection_type"
                    id="default_parameters.intersection_type"
                    required
                    className={inputClasses}
                    value={formData.default_parameters.intersection_type}
                    onChange={handleChange}
                  >
                    <option value="INTERSECTION_TYPE_TRAFFICLIGHT">
                      Traffic Light
                    </option>
                  </select>
                </div>
                <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
                  <div>
                    <label
                      htmlFor="default_parameters.green"
                      className="block text-sm font-medium text-gray-700 dark:text-[#C9D1D9] mb-1"
                    >
                      Green (s)
                    </label>
                    <input
                      type="number"
                      name="default_parameters.green"
                      id="default_parameters.green"
                      required
                      min="1"
                      className={inputClasses}
                      value={formData.default_parameters.green}
                      onChange={handleChange}
                    />
                  </div>
                  <div>
                    <label
                      htmlFor="default_parameters.yellow"
                      className="block text-sm font-medium text-gray-700 dark:text-[#C9D1D9] mb-1"
                    >
                      Yellow (s)
                    </label>
                    <input
                      type="number"
                      name="default_parameters.yellow"
                      id="default_parameters.yellow"
                      required
                      min="1"
                      className={inputClasses}
                      value={formData.default_parameters.yellow}
                      onChange={handleChange}
                    />
                  </div>
                  <div>
                    <label
                      htmlFor="default_parameters.red"
                      className="block text-sm font-medium text-gray-700 dark:text-[#C9D1D9] mb-1"
                    >
                      Red (s)
                    </label>
                    <input
                      type="number"
                      name="default_parameters.red"
                      id="default_parameters.red"
                      required
                      min="1"
                      className={inputClasses}
                      value={formData.default_parameters.red}
                      onChange={handleChange}
                    />
                  </div>
                </div>
                <div>
                  <label
                    htmlFor="default_parameters.speed"
                    className="block text-sm font-medium text-gray-700 dark:text-[#C9D1D9] mb-1"
                  >
                    Vehicle Speed (km/h)
                  </label>
                  <input
                    type="number"
                    name="default_parameters.speed"
                    id="default_parameters.speed"
                    required
                    min="1"
                    className={inputClasses}
                    value={formData.default_parameters.speed}
                    onChange={handleChange}
                  />
                </div>
              </div>
            </div>
          </div>

          {error && (
            <p className="text-red-500 text-sm text-center mt-6">{error}</p>
          )}
          <div className="flex justify-end space-x-4 pt-0 border-t border-gray-200 dark:border-[#30363D] mt-4">
            <button
              type="button"
              onClick={onClose}
              className="px-6 py-2.5 bg-gray-100 dark:bg-[#21262D] border-2 border-gray-300 dark:border-[#30363D] text-gray-700 dark:text-[#C9D1D9] rounded-lg font-medium hover:bg-gray-200 dark:hover:bg-[#30363D] transition-colors duration-150"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className="px-6 py-2.5 bg-[#2da44e] text-white rounded-lg font-medium hover:bg-[#288c42] disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center transition-colors duration-150 shadow-sm"
            >
              {isLoading ? (
                <>
                  <svg
                    className="animate-spin -ml-1 mr-2 h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    ></circle>
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                  {isEditing ? "Updating..." : "Creating..."}
                </>
              ) : isEditing ? (
                "Update Intersection"
              ) : (
                "Create Intersection"
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// =================================================================
// MAIN PAGE COMPONENT
// =================================================================

const getAuthToken = () => {
  return localStorage.getItem("authToken");
};

const intersectionTypeDisplayMap: { [key: string]: string } = {
  INTERSECTION_TYPE_TRAFFICLIGHT: "Traffic Light",
  INTERSECTION_TYPE_UNSPECIFIED: "Traffic Light",
};

const Intersections = () => {
  const navigate = useNavigate();
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

  //  FIX: Added a robust helper function to handle various density formats from the API.
  const normalizeDensityFromAPI = (density: string): TrafficDensityUI => {
    const lowerCaseDensity = density.toLowerCase();
    if (lowerCaseDensity.includes("low")) return "low";
    if (lowerCaseDensity.includes("high")) return "high";
    return "medium"; // Defaults medium, unspecified, or other values to "medium"
  };

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
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Unexpected error");
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
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Unexpected error");
      setIntersections([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateIntersection = async (formData: IntersectionFormData) => {
    setIsCreating(true);
    setCreateError(null);
    try {
      const payload = {
        ...formData,
        //  FIX: traffic_density is now sent directly from formData
        // The 'traffic_density: uiToApiDensityMap[...]' line was removed.
        default_parameters: {
          ...formData.default_parameters,
          intersection_type: "trafficlight",
        },
      };

      const res = await fetch(`${API_BASE_URL}/intersections`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${getAuthToken()}`,
        },
        body: JSON.stringify(payload),
      });
      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.message || "Failed to create intersection");
      }
      setIsModalOpen(false);
      fetchIntersections();
    } catch (err: unknown) {
      setCreateError(err instanceof Error ? err.message : "Unexpected error");
    } finally {
      setIsCreating(false);
    }
  };

  const handleUpdateIntersection = async (formData: IntersectionFormData) => {
    if (!selectedIntersectionId) return;
    setIsCreating(true);
    setCreateError(null);
    try {
      const updatePayload = {
        name: formData.name,
        //  FIX: traffic_density is now sent directly from formData
        traffic_density: formData.traffic_density,
        details: formData.details,
        default_parameters: {
          ...formData.default_parameters,
          intersection_type: "trafficlight",
        },
      };

      const res = await fetch(
        `${API_BASE_URL}/intersections/${selectedIntersectionId}`,
        {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${getAuthToken()}`,
          },
          body: JSON.stringify(updatePayload),
        },
      );

      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.message || "Failed to update intersection");
      }

      setIsModalOpen(false);
      setIsEditing(false);
      setSelectedIntersectionId(null);
      fetchIntersections();
    } catch (err: unknown) {
      setCreateError(err instanceof Error ? err.message : "Unexpected error");
    } finally {
      setIsCreating(false);
    }
  };

  const handleDeleteClick = (id: string) => {
    const intersection = intersections.find((i) => i.id === id);
    if (!intersection) return;

    setIntersectionToDelete({ id: intersection.id, name: intersection.name });
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
          headers: { Authorization: `Bearer ${getAuthToken()}` },
        },
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
    } catch (err: unknown) {
      setError(
        err instanceof Error ? err.message : "Failed to delete intersection",
      );
    } finally {
      setIsDeleting(false);
    }
  };

  const handleSimulate = (id: string) => {
    const intersection = intersections.find((i) => i.id === id);
    if (!intersection) return;

    navigate("/simulation-results", {
      state: {
        name: `Simulation Results for ${intersection.name}`,
        description: `Viewing simulation results for ${intersection.name}`,
        intersections: [intersection.name],
        intersectionIds: [id],
        type: "simulations",
      },
    });
  };

  const handleEditClick = async (id: string) => {
    setCreateError(null);
    try {
      const res = await fetch(`${API_BASE_URL}/intersections/${id}`, {
        headers: { Authorization: `Bearer ${getAuthToken()}` },
      });
      if (!res.ok) {
        throw new Error("Failed to fetch intersection details for editing.");
      }
      const intersection: Intersection = await res.json();

      if (!intersection?.default_parameters?.simulation_parameters) {
        throw new Error("Incomplete intersection data received from server.");
      }

      setEditData({
        name: intersection.name,
        //  FIX: Use the new helper to safely handle density value from the API
        traffic_density: normalizeDensityFromAPI(intersection.traffic_density),
        details: intersection.details,
        default_parameters: {
          ...intersection.default_parameters.simulation_parameters,
          intersection_type:
            intersection.default_parameters.simulation_parameters
              .intersection_type !== "INTERSECTION_TYPE_UNSPECIFIED"
              ? intersection.default_parameters.simulation_parameters
                  .intersection_type
              : "INTERSECTION_TYPE_TRAFFICLIGHT",
        },
      });

      setSelectedIntersectionId(id);
      setIsEditing(true);
      setIsModalOpen(true);
    } catch (err: unknown) {
      alert(
        `Error: ${err instanceof Error ? err.message : "An unexpected error occurred."}`,
      );
    }
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
          intersection.name.toLowerCase().includes(searchQuery.toLowerCase()),
        )
      : intersections;

  return (
    <>
      <div className="intersectionBody flex flex-col min-h-screen bg-gray-50 dark:bg-[#0D1117]">
        <Navbar />
        <main className="main-content flex-grow w-full">
          <div className="max-w-6xl mx-auto w-full px-4 py-8 pb-24">
            <div className="topBar flex flex-col sm:flex-row justify-between items-center mb-8 gap-4">
              <div className="searchContainer relative w-full max-w-md">
                <input
                  type="text"
                  placeholder="Search by Name..."
                  className="searchBar w-full pl-10 pr-4 py-2 border-2 border-gray-200 dark:border-[#30363D] bg-white dark:bg-[#161B22] text-gray-900 dark:text-[#E6EDF3] rounded-full focus:outline-none focus:ring-2 focus:ring-red-500 dark:focus:ring-red-500 transition-colors"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
                <div className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 dark:text-[#7D8590]">
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
                className="addIntersectionBtn w-full sm:w-auto flex-shrink-0 bg-[#0F5BA7] dark:bg-[#238636] hover:bg-blue-700 dark:hover:bg-[#2DA44E] text-white font-medium py-2 px-6 rounded-lg transition-colors shadow-sm"
              >
                Add Intersection
              </button>
            </div>
            <div className="intersections space-y-6">
              {isLoading ? (
                <p className="text-center text-gray-500 dark:text-gray-400">
                  Loading intersections...
                </p>
              ) : error ? (
                <p className="text-center text-red-500">{error}</p>
              ) : filteredIntersections.length > 0 ? (
                filteredIntersections.map((intersection) => {
                  const apiType =
                    intersection.default_parameters.simulation_parameters
                      .intersection_type;
                  const displayType =
                    intersectionTypeDisplayMap[apiType] || "Traffic Light";

                  return (
                    <IntersectionCard
                      key={intersection.id}
                      id={intersection.id}
                      name={intersection.name}
                      location={`${intersection.details.address}`}
                      lanes={displayType}
                      onSimulate={handleSimulate}
                      onEdit={handleEditClick}
                      onDelete={handleDeleteClick}
                    />
                  );
                })
              ) : (
                <p className="text-center text-gray-500 dark:text-gray-400">
                  No intersections found.
                </p>
              )}
            </div>
          </div>
        </main>
        <Footer />
        <HelpMenu />
      </div>

      <CreateIntersectionModal
        isOpen={isModalOpen}
        onClose={() => {
          setIsModalOpen(false);
          setCreateError(null);
          setIsEditing(false);
          setEditData(null);
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
