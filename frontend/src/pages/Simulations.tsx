import React, { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { Eye, Trash2, ChevronDown } from "lucide-react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import { Bar } from "react-chartjs-2";
import { Chart, registerables } from "chart.js";
import type { ChartOptions } from "chart.js";
import { MapContainer, TileLayer, Marker, useMapEvents } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import type { LatLng } from "leaflet";
import "../styles/Simulations.css";
import "@fortawesome/fontawesome-free/css/all.min.css";
import HelpMenu from "../components/HelpMenu";

Chart.register(...registerables);

const simulationsTable1 = [
  {
    id: "SIM001",
    intersection: "Main St & 1st Ave",
    avgWaitTime: 45.2,
    vehicleThroughput: 1200,
    status: "Complete",
  },
  {
    id: "SIM002",
    intersection: "Broadway & 5th St",
    avgWaitTime: 30.8,
    vehicleThroughput: 1500,
    status: "Running",
  },
  {
    id: "SIM003",
    intersection: "Elm St & Park Rd",
    avgWaitTime: 52.1,
    vehicleThroughput: 900,
    status: "Failed",
  },
  {
    id: "SIM004",
    intersection: "Main St & 1st Ave",
    avgWaitTime: 45.2,
    vehicleThroughput: 1200,
    status: "Complete",
  },
  {
    id: "SIM005",
    intersection: "Broadway & 5th St",
    avgWaitTime: 30.8,
    vehicleThroughput: 1500,
    status: "Running",
  },
  {
    id: "SIM006",
    intersection: "Elm St & Park Rd",
    avgWaitTime: 52.1,
    vehicleThroughput: 900,
    status: "Failed",
  },
];

const simulationsTable2 = [
  {
    id: "SIM004",
    intersection: "Oak Ave & Central Blvd",
    avgWaitTime: 28.5,
    vehicleThroughput: 1800,
    status: "Complete",
  },
  {
    id: "SIM005",
    intersection: "Pine St & River Dr",
    avgWaitTime: 47.3,
    vehicleThroughput: 1100,
    status: "Running",
  },
  {
    id: "SIM006",
    intersection: "Maple Rd & 2nd Ave",
    avgWaitTime: 35.6,
    vehicleThroughput: 1300,
    status: "Failed",
  },
  {
    id: "SIM012",
    intersection: "Oak Ave & Central Blvd",
    avgWaitTime: 28.5,
    vehicleThroughput: 1800,
    status: "Complete",
  },
  {
    id: "SIM020",
    intersection: "Pine St & River Dr",
    avgWaitTime: 47.3,
    vehicleThroughput: 1100,
    status: "Running",
  },
  {
    id: "SIM007",
    intersection: "Maple Rd & 2nd Ave",
    avgWaitTime: 35.6,
    vehicleThroughput: 1300,
    status: "Failed",
  },
];

// #region API and Data Types
// Types for the street search functionality
interface Street {
  name: string;
  city?: string;
  province?: string;
  lat?: number | null;
  lon?: number | null;
  osmId?: string | number;
  osmType?: string;
}

interface Intersection {
  lat: number;
  lon: number;
  roads: string[];
  intersection: string;
}

// Type for Nominatim API results
interface NominatimResult {
  place_id: number;
  osm_type: string;
  osm_id: number;
  lat: string;
  lon: string;
  class: string;
  type: string;
  name: string | null;
  display_name: string;
  address?: {
    road?: string;
    suburb?: string;
    city?: string;
    town?: string;
    state?: string;
    province?: string;
    country_code?: string;
  };
}

// Type for Overpass API elements
interface OverpassElement {
  type: "node" | "way" | "relation";
  id: number;
  lat?: number;
  lon?: number;
  tags?: {
    [key: string]: string | undefined; // The fix is applied here
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
type IntersectionWithDistance = Intersection & { distance: number };
// #endregion

const LocationMarker: React.FC<{
  setSelectedLocation: (location: string) => void;
  setCoordinates: (coords: string) => void;
  setIsSnapping?: (snapping: boolean) => void;
}> = ({ setSelectedLocation, setCoordinates, setIsSnapping }) => {
  const [position, setPosition] = useState<LatLng | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  // Function to find nearest intersection using Overpass API
  const findNearestIntersection = async (
    lat: number,
    lon: number,
  ): Promise<Intersection | null> => {
    try {
      setIsSnapping?.(true);
      setIsLoading(true);

      const radius = 0.005;
      const bbox = `${lat - radius},${lon - radius},${lat + radius},${lon + radius}`;

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
          // Use a more explicit type guard to ensure a string[] type
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
                // This line will now work correctly
                roads: uniqueRoads.slice(0, 2),
                intersection: `${uniqueRoads[0]} & ${uniqueRoads[1]}`,
                distance,
              });
            }
          }
        }
      }

      intersections.sort((a, b) => a.distance - b.distance);

      if (intersections.length > 0) {
        return intersections[0];
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

      // First set the clicked position as a temporary marker
      setPosition(e.latlng);

      try {
        // Try to find nearest intersection
        const nearestIntersection = await findNearestIntersection(
          e.latlng.lat,
          e.latlng.lng,
        );

        if (nearestIntersection) {
          // Snap to intersection
          const snappedPosition = {
            lat: nearestIntersection.lat,
            lng: nearestIntersection.lon,
          } as LatLng;

          setPosition(snappedPosition);
          setSelectedLocation(nearestIntersection.intersection);
          setCoordinates(
            `${nearestIntersection.lat.toFixed(6)}, ${nearestIntersection.lon.toFixed(6)}`,
          );

          console.log(
            "Snapped to intersection:",
            nearestIntersection.intersection,
          );
        } else {
          // Fallback to clicked coordinates if no intersection found
          const coordinates = `${e.latlng.lat.toFixed(6)}, ${e.latlng.lng.toFixed(6)}`;
          setSelectedLocation(coordinates);
          setCoordinates(coordinates);

          console.log("No intersection found, using clicked coordinates");
        }
      } catch (error) {
        console.error("Error processing map click:", error);
        // Fallback to clicked coordinates
        const coordinates = `${e.latlng.lat.toFixed(6)}, ${e.latlng.lng.toFixed(6)}`;
        setSelectedLocation(coordinates);
        setCoordinates(coordinates);
      }
    },
  });

  return (
    <>
      {position && <Marker position={position} />}
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

const StreetSearchComponent: React.FC<{
  onIntersectionSelect: (intersection: string) => void;
}> = ({ onIntersectionSelect }) => {
  const [firstStreet, setFirstStreet] = useState("");
  const [secondStreet, setSecondStreet] = useState("");
  const [firstStreetSuggestions, setFirstStreetSuggestions] = useState<
    Street[]
  >([]);
  const [secondStreetSuggestions, setSecondStreetSuggestions] = useState<
    Street[]
  >([]);
  const [showFirstDropdown, setShowFirstDropdown] = useState(false);
  const [isLoadingFirst, setIsLoadingFirst] = useState(false);
  const [isLoadingSecond, setIsLoadingSecond] = useState(false);
  const [selectedFirstStreet, setSelectedFirstStreet] = useState<Street | null>(
    null,
  );

  const firstStreetRef = useRef<HTMLDivElement>(null);
  const secondStreetRef = useRef<HTMLDivElement>(null);
  const debounceRef = useRef<NodeJS.Timeout | null>(null);

  // API function for fetching South African streets
  const searchStreets = async (query: string): Promise<Street[]> => {
    if (!query || query.length < 3) return [];

    try {
      const cleanQuery = query.trim();
      console.log("Searching for streets with query:", cleanQuery);

      // Use multiple strategies to find streets
      const strategies = [
        // Strategy 1: Direct road search with South Africa constraint
        () => searchNominatimStreets(cleanQuery, "road"),
        // Strategy 2: Search with highway tag
        () => searchNominatimStreets(cleanQuery + " highway", "way"),
        // Strategy 3: Search with street/road suffix if not present
        () => {
          const hasStreetSuffix =
            /\b(street|road|avenue|drive|lane|way|boulevard|crescent)\b/i.test(
              cleanQuery,
            );
          if (!hasStreetSuffix) {
            return searchNominatimStreets(cleanQuery + " road", "road");
          }
          return [];
        },
      ];

      for (const strategy of strategies) {
        const results = await strategy();
        if (results.length > 0) {
          console.log("Found streets via strategy:", results);
          return results;
        }
        // Small delay between strategies
        await new Promise((resolve) => setTimeout(resolve, 300));
      }

      // Fallback to common streets
      return getCommonSouthAfricanStreets(query);
    } catch (error) {
      console.error("Error in searchStreets:", error);
      return getCommonSouthAfricanStreets(query);
    }
  };

  // Nominatim search function
  const searchNominatimStreets = async (
    query: string,
    type: string = "road",
  ): Promise<Street[]> => {
    try {
      const params = new URLSearchParams({
        format: "json",
        addressdetails: "1",
        limit: "20",
        countrycodes: "za",
        q: query,
        extratags: "1",
      });

      // Add type-specific parameters
      if (type === "road") {
        params.append("class", "highway");
      }

      const url = `https://nominatim.openstreetmap.org/search?${params.toString()}`;
      console.log("Nominatim URL:", url);

      const response = await fetch(url);
      if (!response.ok) {
        console.warn(`Nominatim request failed: ${response.status}`);
        return [];
      }

      const data: NominatimResult[] = await response.json();
      console.log("Nominatim raw results:", data.length);

      if (!Array.isArray(data) || data.length === 0) {
        return [];
      }

      return processSouthAfricanStreets(data);
    } catch (error) {
      console.error("Error in searchNominatimStreets:", error);
      return [];
    }
  };

  // Improved street data processing
  const processSouthAfricanStreets = (data: NominatimResult[]): Street[] => {
    if (!Array.isArray(data)) return [];

    const processed = data
      .filter((item: NominatimResult) => {
        // Must be in South Africa
        const isInSA =
          item.address?.country_code === "za" ||
          item.display_name?.includes("South Africa");

        // Must be a road/highway
        const isRoad =
          item.class === "highway" ||
          item.type === "road" ||
          item.address?.road ||
          /\b(street|road|avenue|drive|lane|way|boulevard|crescent|highway)\b/i.test(
            item.display_name,
          );

        // Filter out non-road types
        const validType = ![
          "hamlet",
          "village",
          "isolated_dwelling",
          "city",
          "town",
          "suburb",
        ].includes(item.type);

        return isInSA && isRoad && validType;
      })
      .map((item: NominatimResult) => {
        const streetName = extractStreetName(item);
        const coordinates = getCoordinates(item);

        return {
          name: streetName,
          city:
            item.address?.city ||
            item.address?.town ||
            item.address?.suburb ||
            extractCityFromDisplayName(item.display_name),
          province:
            item.address?.state || item.address?.province || "South Africa",
          lat: coordinates.lat,
          lon: coordinates.lon,
          osmId: item.osm_id,
          osmType: item.osm_type,
        };
      })
      .filter((street: Street) => street.name && street.name.length > 0)
      .filter(
        (street: Street, index: number, self: Street[]) =>
          // Remove duplicates based on name and city
          self.findIndex(
            (s) =>
              s.name.toLowerCase() === street.name.toLowerCase() &&
              s.city?.toLowerCase() === street.city?.toLowerCase(),
          ) === index,
      )
      .slice(0, 15);

    console.log("Processed streets:", processed);
    return processed;
  };

  // street name extraction
  const extractStreetName = (item: NominatimResult): string => {
    // Priority order for extracting street name
    const candidates = [
      item.address?.road,
      item.name,
      item.display_name?.split(",")[0]?.trim(),
    ].filter(Boolean) as string[];

    for (const candidate of candidates) {
      if (candidate && isValidStreetName(candidate)) {
        return cleanStreetName(candidate);
      }
    }

    // Extract from display name if no direct match
    const displayParts = item.display_name?.split(",") || [];
    for (const part of displayParts) {
      const cleaned = part.trim();
      if (isValidStreetName(cleaned)) {
        return cleanStreetName(cleaned);
      }
    }

    return item.display_name?.split(",")[0]?.trim() || "Unknown Street";
  };

  // Clean and standardize street names
  const cleanStreetName = (name: string): string => {
    return name
      .replace(/\s+/g, " ")
      .trim()
      .replace(/\b\w/g, (l) => l.toUpperCase());
  };

  // Get coordinates with fallback
  const getCoordinates = (
    item: NominatimResult,
  ): { lat: number | null; lon: number | null } => {
    const lat = parseFloat(item.lat);
    const lon = parseFloat(item.lon);

    return {
      lat: !isNaN(lat) ? lat : null,
      lon: !isNaN(lon) ? lon : null,
    };
  };

  // Extract city from display name
  const extractCityFromDisplayName = (displayName: string): string => {
    if (!displayName) return "South Africa";

    const parts = displayName.split(",").map((p) => p.trim());
    for (let i = 1; i < Math.min(parts.length, 4); i++) {
      const part = parts[i];
      if (part && !isValidStreetName(part) && part !== "South Africa") {
        return part;
      }
    }

    return "South Africa";
  };

  // find actual intersecting streets using Overpass API
  const findIntersectingStreets = async (street: Street): Promise<Street[]> => {
    if (!street.name) return [];

    try {
      console.log("Finding intersections for:", street.name, "in", street.city);

      // Strategy 1: Use coordinates if available
      if (street.lat && street.lon) {
        const coordResults = await findIntersectionsAtCoordinates(
          street.lat,
          street.lon,
          street.name,
        );
        if (coordResults.length > 0) {
          console.log(
            "Found intersections via coordinates:",
            coordResults.length,
          );
          return coordResults;
        }
      }

      // Strategy 2: Use Overpass API with street name and city
      const overpassResults = await findIntersectionsViaOverpass(street);
      if (overpassResults.length > 0) {
        console.log(
          "Found intersections via Overpass:",
          overpassResults.length,
        );
        return overpassResults;
      }

      // Strategy 3: Search Nominatim for intersection mentions
      const nominatimResults = await searchIntersectionMentions(
        street.name,
        street.city,
      );
      if (nominatimResults.length > 0) {
        console.log(
          "Found intersections via Nominatim search:",
          nominatimResults.length,
        );
        return nominatimResults;
      }

      // Fallback to contextual intersections
      return getContextualIntersectingStreets(street);
    } catch (error) {
      console.error("Error finding intersecting streets:", error);
      return getContextualIntersectingStreets(street);
    }
  };

  // Find intersections using precise coordinate-based search
  const findIntersectionsAtCoordinates = async (
    lat: number,
    lon: number,
    streetName: string,
  ): Promise<Street[]> => {
    try {
      console.log(
        `Finding intersections at coordinates: ${lat}, ${lon} for ${streetName}`,
      );

      // Create Overpass query to find intersections at specific coordinates
      const radius = 0.005; // About 500m radius
      const overpassQuery = `
        [out:json][timeout:30];
        (
          way["highway"~"^(primary|secondary|tertiary|residential|trunk|motorway|unclassified|road)$"]
             ["name"]
             (${lat - radius},${lon - radius},${lat + radius},${lon + radius});
        );
        (._;>;);
        out geom;
      `;

      const response = await fetch("https://overpass-api.de/api/interpreter", {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: `data=${encodeURIComponent(overpassQuery)}`,
      });

      if (!response.ok) {
        console.warn("Overpass API failed:", response.status);
        return [];
      }

      const data = await response.json();
      return extractIntersectionsFromOverpassData(data.elements, streetName);
    } catch (error) {
      console.error("Error finding intersections at coordinates:", error);
      return [];
    }
  };

  // Overpass API intersection finding
  const findIntersectionsViaOverpass = async (
    street: Street,
  ): Promise<Street[]> => {
    try {
      const searchArea =
        street.city && street.city !== "South Africa"
          ? `area["name"~"${street.city}",i]["place"~"^(city|town)$"];`
          : `area["ISO3166-1"="ZA"][admin_level=2];`;

      const overpassQuery = `
        [out:json][timeout:30];
        (
          ${searchArea}
        )->.searchArea;
        (
          way["highway"]
             ["name"~"${street.name}",i]
             (area.searchArea);
        );
        (._;>;);
        out geom;
      `;

      console.log("Overpass query for street search:", overpassQuery);

      const response = await fetch("https://overpass-api.de/api/interpreter", {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: `data=${encodeURIComponent(overpassQuery)}`,
      });

      if (!response.ok) {
        console.warn("Overpass API failed:", response.status);
        return [];
      }

      const data = await response.json();

      if (!data.elements || data.elements.length === 0) {
        console.log("No ways found for street:", street.name);
        return [];
      }

      return await findAllIntersectingWays(data.elements, street.name);
    } catch (error) {
      console.error("Error with Overpass API:", error);
      return [];
    }
  };

  // Find all ways that intersect with the target street ways
  const findAllIntersectingWays = async (
    targetWays: OverpassElement[],
    streetName: string,
  ): Promise<Street[]> => {
    try {
      const targetNodeIds = new Set<number>();
      targetWays.forEach((way) => {
        if (way.type === "way" && way.nodes) {
          way.nodes.forEach((nodeId: number) => targetNodeIds.add(nodeId));
        }
      });

      if (targetNodeIds.size === 0) {
        console.log("No nodes found in target ways");
        return [];
      }

      console.log("Target street has", targetNodeIds.size, "nodes");

      // Find all ways that share nodes with our target street
      const nodeList = Array.from(targetNodeIds).slice(0, 100);
      const nodeFilter = nodeList.map((id) => `node(${id})`).join("");

      const intersectionQuery = `
        [out:json][timeout:30];
        (
          ${nodeFilter}
        )->.targetNodes;
        (
          way["highway"]
             ["name"]
             (bn.targetNodes);
        );
        out geom;
      `;

      const response = await fetch("https://overpass-api.de/api/interpreter", {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: `data=${encodeURIComponent(intersectionQuery)}`,
      });

      if (!response.ok) {
        console.warn("Intersection query failed:", response.status);
        return [];
      }

      const data = await response.json();
      return extractIntersectionsFromOverpassData(data.elements, streetName);
    } catch (error) {
      console.error("Error finding intersecting ways:", error);
      return [];
    }
  };

  // Extract and process intersection data from Overpass results
  const extractIntersectionsFromOverpassData = (
    elements: OverpassElement[],
    originalStreet: string,
  ): Street[] => {
    if (!Array.isArray(elements)) return [];

    const intersectingStreets: Street[] = [];
    const originalStreetLower = originalStreet.toLowerCase();

    // Filter ways that are different from our original street
    const intersectingWays = elements.filter(
      (el) =>
        el.type === "way" &&
        el.tags &&
        el.tags.name &&
        el.tags.highway &&
        !el.tags.name.toLowerCase().includes(originalStreetLower) &&
        !originalStreetLower.includes(el.tags.name.toLowerCase()),
    );

    console.log(
      `Found ${intersectingWays.length} potentially intersecting ways`,
    );

    intersectingWays.forEach((way) => {
      if (way.tags?.name && isValidStreetName(way.tags.name)) {
        // Calculate center point for the way
        const bounds = way.bounds;
        const centerLat = bounds ? (bounds.minlat + bounds.maxlat) / 2 : null;
        const centerLon = bounds ? (bounds.minlon + bounds.maxlon) / 2 : null;

        intersectingStreets.push({
          name: cleanStreetName(way.tags.name),
          city: way.tags["addr:city"] || "South Africa",
          province: way.tags["addr:state"] || way.tags["addr:province"] || "",
          lat: centerLat,
          lon: centerLon,
          osmId: way.id,
          osmType: "way",
        });
      }
    });

    // Remove duplicates and sort by name
    const uniqueStreets = intersectingStreets
      .filter(
        (street, index, self) =>
          self.findIndex(
            (s) => s.name.toLowerCase() === street.name.toLowerCase(),
          ) === index,
      )
      .sort((a, b) => a.name.localeCompare(b.name));

    console.log("Extracted unique intersecting streets:", uniqueStreets.length);
    return uniqueStreets.slice(0, 20);
  };

  // Search for explicit intersection mentions in Nominatim
  const searchIntersectionMentions = async (
    streetName: string,
    city?: string,
  ): Promise<Street[]> => {
    try {
      const searchQueries = [
        `"${streetName}" intersection ${city || "South Africa"}`,
        `"${streetName}" junction ${city || "South Africa"}`,
        `"${streetName}" & ${city || "South Africa"}`,
        `intersection "${streetName}" ${city || "South Africa"}`,
      ];

      for (const query of searchQueries) {
        const url =
          `https://nominatim.openstreetmap.org/search?` +
          new URLSearchParams({
            format: "json",
            addressdetails: "1",
            limit: "30",
            countrycodes: "za",
            q: query,
          }).toString();

        console.log("Searching intersection mentions:", query);

        const response = await fetch(url);
        if (response.ok) {
          const data: NominatimResult[] = await response.json();
          const intersections = extractIntersectionsFromDisplayNames(
            data,
            streetName,
          );
          if (intersections.length > 0) {
            return intersections;
          }
        }

        // Delay between requests
        await new Promise((resolve) => setTimeout(resolve, 500));
      }

      return [];
    } catch (error) {
      console.error("Error searching intersection mentions:", error);
      return [];
    }
  };

  // Extract intersections from Nominatim display names
  const extractIntersectionsFromDisplayNames = (
    data: NominatimResult[],
    originalStreet: string,
  ): Street[] => {
    if (!Array.isArray(data)) return [];

    const intersectingStreets: Street[] = [];
    const originalStreetLower = originalStreet.toLowerCase();

    data.forEach((item) => {
      if (!item.display_name) return;

      const displayName = item.display_name.toLowerCase();

      // Look for intersection patterns
      const intersectionPatterns = [
        /(.+?)\s*&\s*(.+?),/,
        /(.+?)\s+and\s+(.+?),/,
        /(.+?)\s*\/\s*(.+?),/,
        /intersection\s+of\s+(.+?)\s+and\s+(.+?),/,
        /(.+?)\s+x\s+(.+?),/,
      ];

      for (const pattern of intersectionPatterns) {
        const match = displayName.match(pattern);
        if (match) {
          const street1 = match[1].trim();
          const street2 = match[2].trim();

          let intersectingStreet = null;

          if (street1.includes(originalStreetLower)) {
            intersectingStreet = street2;
          } else if (street2.includes(originalStreetLower)) {
            intersectingStreet = street1;
          }

          if (intersectingStreet && isValidStreetName(intersectingStreet)) {
            intersectingStreets.push({
              name: toTitleCase(intersectingStreet),
              city: item.address?.city || item.address?.town || "South Africa",
              province: item.address?.state || "",
              lat: parseFloat(item.lat),
              lon: parseFloat(item.lon),
            });
          }
        }
      }
    });

    // Remove duplicates
    const uniqueStreets = intersectingStreets.filter(
      (street, index, self) =>
        self.findIndex(
          (s) => s.name.toLowerCase() === street.name.toLowerCase(),
        ) === index,
    );

    return uniqueStreets.slice(0, 15);
  };

  // Helper to check if a string is a valid street name
  const isValidStreetName = (name: string): boolean => {
    if (!name || typeof name !== "string") return false;

    const cleaned = name.trim();
    if (cleaned.length < 3) return false;

    // Must contain street-like suffixes or be a recognized road name
    const streetSuffixes =
      /(street|road|avenue|drive|lane|way|boulevard|crescent|place|close|court|grove|gardens|park|square|circle|terrace|highway|route)$/i;
    const isNamedRoad = /^[A-Za-z][A-Za-z0-9\s-'.]*$/i.test(cleaned);

    return streetSuffixes.test(cleaned) || (isNamedRoad && cleaned.length >= 5);
  };

  // Helper to convert to title case
  const toTitleCase = (str: string): string => {
    return str.replace(
      /\w\S*/g,
      (txt) => txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase(),
    );
  };

  // common South African streets as fallback
  const getCommonSouthAfricanStreets = (query: string): Street[] => {
    const commonStreets: Street[] = [
      {
        name: "Main Road",
        city: "Various Cities",
        province: "All Provinces",
        lat: null,
        lon: null,
      },
      {
        name: "Church Street",
        city: "Various Cities",
        province: "All Provinces",
        lat: null,
        lon: null,
      },
      {
        name: "Voortrekker Road",
        city: "Cape Town",
        province: "Western Cape",
        lat: null,
        lon: null,
      },
      {
        name: "Jan Smuts Avenue",
        city: "Johannesburg",
        province: "Gauteng",
        lat: null,
        lon: null,
      },
      {
        name: "Nelson Mandela Drive",
        city: "Various Cities",
        province: "All Provinces",
        lat: null,
        lon: null,
      },
      {
        name: "Commissioner Street",
        city: "Johannesburg",
        province: "Gauteng",
        lat: null,
        lon: null,
      },
      {
        name: "Long Street",
        city: "Cape Town",
        province: "Western Cape",
        lat: null,
        lon: null,
      },
      {
        name: "Smith Street",
        city: "Durban",
        province: "KwaZulu-Natal",
        lat: null,
        lon: null,
      },
      {
        name: "Paul Kruger Street",
        city: "Pretoria",
        province: "Gauteng",
        lat: null,
        lon: null,
      },
      {
        name: "William Nicol Drive",
        city: "Johannesburg",
        province: "Gauteng",
        lat: null,
        lon: null,
      },
      {
        name: "Barry Hertzog Avenue",
        city: "Johannesburg",
        province: "Gauteng",
        lat: null,
        lon: null,
      },
      {
        name: "Louis Botha Avenue",
        city: "Johannesburg",
        province: "Gauteng",
        lat: null,
        lon: null,
      },
      {
        name: "Oxford Road",
        city: "Johannesburg",
        province: "Gauteng",
        lat: null,
        lon: null,
      },
      {
        name: "Adderley Street",
        city: "Cape Town",
        province: "Western Cape",
        lat: null,
        lon: null,
      },
      {
        name: "West Street",
        city: "Durban",
        province: "KwaZulu-Natal",
        lat: null,
        lon: null,
      },
      {
        name: "Schoeman Street",
        city: "Pretoria",
        province: "Gauteng",
        lat: null,
        lon: null,
      },
      {
        name: "Simon Vermooten Road",
        city: "Pretoria",
        province: "Gauteng",
        lat: null,
        lon: null,
      },
      {
        name: "Beyers Naude Drive",
        city: "Johannesburg",
        province: "Gauteng",
        lat: null,
        lon: null,
      },
    ];

    const filtered = commonStreets.filter((street) =>
      street.name.toLowerCase().includes(query.toLowerCase()),
    );

    return filtered.length > 0 ? filtered : commonStreets.slice(0, 10);
  };

  // Enhanced contextual intersections
  const getContextualIntersectingStreets = (street: Street): Street[] => {
    const streetName = street.name.toLowerCase();

    // Known major intersections for specific streets
    const knownIntersections: { [key: string]: Street[] } = {
      "simon vermooten road": [
        {
          name: "Jan Shoba Street",
          city: "Pretoria",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
        {
          name: "Zambezi Drive",
          city: "Pretoria",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
        {
          name: "Solomon Mahlangu Drive",
          city: "Pretoria",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
        {
          name: "Lynnwood Road",
          city: "Pretoria",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
        {
          name: "Atterbury Road",
          city: "Pretoria",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
      ],
      "jan smuts avenue": [
        {
          name: "Rosebank Road",
          city: "Johannesburg",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
        {
          name: "Oxford Road",
          city: "Johannesburg",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
        {
          name: "Barry Hertzog Avenue",
          city: "Johannesburg",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
        {
          name: "Empire Road",
          city: "Johannesburg",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
      ],
      "william nicol drive": [
        {
          name: "Republic Road",
          city: "Johannesburg",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
        {
          name: "Witkoppen Road",
          city: "Johannesburg",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
        {
          name: "Main Road",
          city: "Johannesburg",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
        {
          name: "Bryanston Drive",
          city: "Johannesburg",
          province: "Gauteng",
          lat: null,
          lon: null,
        },
      ],
      "main road": [
        {
          name: "Church Street",
          city: street.city || "Various",
          province: street.province || "",
          lat: null,
          lon: null,
        },
        {
          name: "Market Street",
          city: street.city || "Various",
          province: street.province || "",
          lat: null,
          lon: null,
        },
        {
          name: "Station Road",
          city: street.city || "Various",
          province: street.province || "",
          lat: null,
          lon: null,
        },
      ],
    };

    // Check for known intersections
    for (const [knownStreet, intersections] of Object.entries(
      knownIntersections,
    )) {
      if (
        streetName.includes(knownStreet) ||
        knownStreet.includes(streetName)
      ) {
        return intersections;
      }
    }

    // Generic fallback based on location
    const genericIntersections: Street[] = [
      {
        name: "Main Road",
        city: street.city || "South Africa",
        province: street.province || "",
        lat: null,
        lon: null,
      },
      {
        name: "Church Street",
        city: street.city || "South Africa",
        province: street.province || "",
        lat: null,
        lon: null,
      },
      {
        name: "Market Street",
        city: street.city || "South Africa",
        province: street.province || "",
        lat: null,
        lon: null,
      },
      {
        name: "Station Street",
        city: street.city || "South Africa",
        province: street.province || "",
        lat: null,
        lon: null,
      },
    ];

    return genericIntersections.filter(
      (s) => s.name.toLowerCase() !== streetName,
    );
  };

  // Debounced search for first street
  const handleFirstStreetChange = (value: string) => {
    setFirstStreet(value);
    setSelectedFirstStreet(null);
    setSecondStreet("");
    setSecondStreetSuggestions([]);

    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    if (value.length >= 3) {
      setIsLoadingFirst(true);
      setShowFirstDropdown(true);

      debounceRef.current = setTimeout(async () => {
        try {
          const suggestions = await searchStreets(value);
          setFirstStreetSuggestions(suggestions);
        } catch (error) {
          console.error("Error fetching street suggestions:", error);
          setFirstStreetSuggestions(getCommonSouthAfricanStreets(value));
        } finally {
          setIsLoadingFirst(false);
        }
      }, 1000);
    } else {
      setFirstStreetSuggestions([]);
      setShowFirstDropdown(false);
      setIsLoadingFirst(false);
    }
  };

  // Handle first street selection
  const handleFirstStreetSelect = async (street: Street) => {
    setFirstStreet(street.name);
    setSelectedFirstStreet(street);
    setShowFirstDropdown(false);
    setIsLoadingSecond(true);

    try {
      console.log("Finding intersections for selected street:", street);
      const intersecting = await findIntersectingStreets(street);
      console.log("Found intersecting streets:", intersecting);
      setSecondStreetSuggestions(intersecting);
    } catch (error) {
      console.error("Error finding intersecting streets:", error);
      setSecondStreetSuggestions(getContextualIntersectingStreets(street));
    } finally {
      setIsLoadingSecond(false);
    }
  };

  // Handle second street selection
  const handleSecondStreetSelect = (street: Street) => {
    setSecondStreet(street.name);

    const intersection = `${firstStreet} & ${street.name}`;
    onIntersectionSelect(intersection);

    // Reset form
    setFirstStreet("");
    setSecondStreet("");
    setSelectedFirstStreet(null);
    setFirstStreetSuggestions([]);
    setSecondStreetSuggestions([]);
  };

  // Click outside handlers
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        firstStreetRef.current &&
        !firstStreetRef.current.contains(event.target as Node)
      ) {
        setShowFirstDropdown(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div className="space-y-4">
      {/* First Street Input */}
      <div className="relative" ref={firstStreetRef}>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          First Street
        </label>
        <div className="relative">
          <input
            type="text"
            value={firstStreet}
            onChange={(e) => handleFirstStreetChange(e.target.value)}
            onFocus={() =>
              firstStreetSuggestions.length > 0 && setShowFirstDropdown(true)
            }
            className="w-full p-2 pr-8 rounded-md border-2 border-gray-300 dark:border-[#30363D] bg-white dark:bg-[#161B22] text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="Type a street name (min 3 characters)..."
          />
          <ChevronDown className="absolute right-2 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
        </div>

        {showFirstDropdown && (
          <div className="absolute z-10 w-full mt-1 bg-white dark:bg-[#161B22] border border-gray-300 dark:border-[#30363D] rounded-md shadow-lg max-h-60 overflow-y-auto">
            {isLoadingFirst ? (
              <div className="p-3 text-center text-gray-500 dark:text-gray-400">
                <div className="animate-spin inline-block w-4 h-4 border-2 border-current border-t-transparent rounded-full mr-2"></div>
                Searching South African streets...
              </div>
            ) : firstStreetSuggestions.length > 0 ? (
              firstStreetSuggestions.map((street, index) => (
                <button
                  key={index}
                  onClick={() => handleFirstStreetSelect(street)}
                  className="w-full text-left p-3 hover:bg-gray-100 dark:hover:bg-[#21262D] border-b border-gray-200 dark:border-[#30363D] last:border-b-0"
                >
                  <div className="font-medium text-gray-900 dark:text-gray-200">
                    {street.name}
                  </div>
                  {street.city && street.city !== "Various Cities" && (
                    <div className="text-sm text-gray-500 dark:text-gray-400">
                      {street.city}
                      {street.province &&
                        street.province !== "South Africa" &&
                        `, ${street.province}`}
                    </div>
                  )}
                </button>
              ))
            ) : (
              <div className="p-3 text-center text-gray-500 dark:text-gray-400">
                No streets found. Try typing more characters.
              </div>
            )}
          </div>
        )}
      </div>

      {/* Second Street Dropdown */}
      {selectedFirstStreet && (
        <div className="relative" ref={secondStreetRef}>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Intersecting Street
            <span className="text-xs text-gray-500 block">
              Streets that actually intersect with "{selectedFirstStreet.name}"
            </span>
          </label>
          <div className="relative">
            <select
              value={secondStreet}
              onChange={(e) => {
                const selectedStreet = secondStreetSuggestions.find(
                  (s) => s.name === e.target.value,
                );
                if (selectedStreet) {
                  handleSecondStreetSelect(selectedStreet);
                }
              }}
              className="w-full p-2 pr-8 rounded-md border-2 border-gray-300 dark:border-[#30363D] bg-white dark:bg-[#161B22] text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              disabled={isLoadingSecond}
            >
              <option value="">
                {isLoadingSecond
                  ? "Finding actual intersecting streets..."
                  : secondStreetSuggestions.length > 0
                    ? "Select an intersecting street..."
                    : "No intersecting streets found"}
              </option>
              {!isLoadingSecond &&
                secondStreetSuggestions.map((street, index) => (
                  <option key={index} value={street.name}>
                    {street.name}
                    {street.city &&
                      street.city !== "Various Cities" &&
                      street.city !== "South Africa" &&
                      ` - ${street.city}`}
                    {street.province &&
                      street.province !== "All Provinces" &&
                      street.province !== "South Africa" &&
                      street.province !== "" &&
                      `, ${street.province}`}
                  </option>
                ))}
            </select>
            <ChevronDown className="absolute right-2 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" />
          </div>

          {isLoadingSecond && (
            <div className="mt-2 text-sm text-gray-500 dark:text-gray-400 flex items-center">
              <div className="animate-spin inline-block w-4 h-4 border-2 border-current border-t-transparent rounded-full mr-2"></div>
              Searching for actual intersections using OpenStreetMap data...
            </div>
          )}

          {!isLoadingSecond &&
            secondStreetSuggestions.length === 0 &&
            selectedFirstStreet && (
              <div className="mt-2 text-sm text-yellow-600 dark:text-yellow-400">
                ⚠️ No intersections found via API. This may be due to:
                <ul className="list-disc list-inside mt-1 text-xs">
                  <li>Limited OpenStreetMap data for this street</li>
                  <li>Street name variations in the database</li>
                  <li>API rate limiting</li>
                </ul>
              </div>
            )}
        </div>
      )}

      {selectedFirstStreet && (
        <div className="mt-4 p-3 bg-gray-50 dark:bg-gray-800 rounded-md text-xs">
          <div className="font-medium text-gray-700 dark:text-gray-300 mb-1">
            Debug Info:
          </div>
          <div className="text-gray-600 dark:text-gray-400">
            Selected: {selectedFirstStreet.name}
            {selectedFirstStreet.city && ` (${selectedFirstStreet.city})`}
            {selectedFirstStreet.lat &&
              selectedFirstStreet.lon &&
              ` - Coords: ${selectedFirstStreet.lat.toFixed(4)}, ${selectedFirstStreet.lon.toFixed(4)}`}
          </div>
          <div className="text-gray-600 dark:text-gray-400">
            Found {secondStreetSuggestions.length} intersecting streets
          </div>
        </div>
      )}
    </div>
  );
};

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
  const navigate = useNavigate();
  const [simulationName, setSimulationName] = useState("");
  const [simulationDescription, setSimulationDescription] = useState("");
  const [selectedIntersections, setSelectedIntersections] = useState<string[]>(
    [],
  );
  const [activeTab, setActiveTab] = useState<"List" | "Search" | "Map">("List");
  const [coordinates, setCoordinates] = useState<string | null>(null);
  const [isSnapping, setIsSnapping] = useState(false);

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

  const handleMapSelection = (location: string) => {
    if (!selectedIntersections.includes(location)) {
      setSelectedIntersections([...selectedIntersections, location]);
    }
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
    setCoordinates(null);
    navigate("/simulation-results", { state: simulationData });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div className="simulation-modal-content bg-white dark:bg-[#161B22] rounded-lg shadow-xl w-full max-w-md p-6 relative max-h-[90vh] overflow-y-auto">
        <button
          onClick={onClose}
          className="crossBtn absolute top-4 right-4 text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-100"
        >
          ✕
        </button>
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
                  >
                    ✕
                  </button>
                </div>
              ))}
            </div>
            <div className="intersection-tabs flex space-x-2 mb-3">
              <button
                onClick={() => setActiveTab("List")}
                className={`px-3 py-1 rounded-md text-sm font-medium ${activeTab === "List" ? "bg-[#2B9348] text-white dark:bg-[#2DA44E]" : "bg-gray-200 text-gray-700 dark:bg-gray-600 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-500"} transition-all duration-300`}
              >
                List
              </button>
              <button
                onClick={() => setActiveTab("Search")}
                className={`px-3 py-1 rounded-md text-sm font-medium ${activeTab === "Search" ? "bg-[#2B9348] text-white dark:bg-[#2DA44E]" : "bg-gray-200 text-gray-700 dark:bg-gray-600 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-500"} transition-all duration-300`}
              >
                Search
              </button>
              <button
                onClick={() => setActiveTab("Map")}
                className={`px-3 py-1 rounded-md text-sm font-medium ${activeTab === "Map" ? "bg-[#2B9348] text-white dark:bg-[#2DA44E]" : "bg-gray-200 text-gray-700 dark:bg-gray-600 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-500"} transition-all duration-300`}
              >
                Map
              </button>
            </div>
            {activeTab === "List" && (
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
            )}
            {activeTab === "Search" && (
              <StreetSearchComponent
                onIntersectionSelect={handleAddIntersection}
              />
            )}
            {activeTab === "Map" && (
              <div className="relative">
                <MapContainer
                  center={[-26.2041, 28.0473]}
                  zoom={6}
                  style={{ height: "200px", width: "100%" }}
                >
                  <TileLayer
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    attribution='© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
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
                    💡 Click anywhere on the map to automatically find the
                    nearest road intersection
                  </p>
                </div>
              </div>
            )}
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
  simulations: Array<{
    id: string;
    intersection: string;
    avgWaitTime: number;
    vehicleThroughput: number;
    status: string;
  }>;
  currentPage: number;
  setCurrentPage: (page: number) => void;
}> = ({ simulations, currentPage, setCurrentPage }) => {
  const rowsPerPage = 4;
  const totalPages = Math.ceil(simulations.length / rowsPerPage);
  const startIndex = currentPage * rowsPerPage;
  const endIndex = startIndex + rowsPerPage;
  const paginatedSimulations = simulations.slice(startIndex, endIndex);

  const chartOptions: ChartOptions<"bar"> = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: {
        backgroundColor: "rgba(0, 0, 0, 0.8)",
        cornerRadius: 8,
        padding: 10,
        titleFont: { size: 12, weight: "bold" },
        bodyFont: { size: 12 },
        displayColors: false,
      },
    },
    scales: {
      x: { display: false },
      y: { beginAtZero: true, display: false },
    },
    animation: {
      duration: 1000,
      easing: "easeOutQuart",
    },
    elements: {
      bar: {
        borderRadius: 6,
        borderWidth: 0,
      },
    },
  };

  const handleViewResults = (simId: string) => {
    alert(`Viewing results for simulation ${simId}`);
    // Replace with actual logic
  };

  const handleDelete = (simId: string) => {
    alert(`Deleting simulation ${simId}`);
    // Replace with actual delete logic
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const statusClass = (status: string) => {
    switch (status) {
      case "Complete":
        return "bg-green-200 text-green-800 border-green-300";
      case "Running":
        return "bg-yellow-200 text-yellow-800 border-yellow-300";
      case "Failed":
        return "bg-red-200 text-red-800 border-red-300";
      default:
        return "bg-gray-200 text-gray-800 border-gray-300";
    }
  };

  return (
    <div className="simTable bg-white dark:bg-[#161B22] shadow-md rounded-lg overflow-hidden table-fixed-height relative">
      <table className="simulationTable min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead className="simTableHead bg-gray-50 dark:bg-[#161B22]">
          <tr>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
              Simulation ID
            </th>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
              Intersection
            </th>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
              Avg Wait Time
            </th>
            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
              Throughput
            </th>
            <th className="graphTHead px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
              Graph
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
          {paginatedSimulations.map((sim) => {
            const chartData = {
              labels: ["Wait", "Throughput"],
              datasets: [
                {
                  data: [sim.avgWaitTime, sim.vehicleThroughput / 10],
                  backgroundColor: ["#2B9348", "#0F5BA7"],
                  hoverBackgroundColor: ["#6EE7B7", "#60A5FA"],
                  borderWidth: 0,
                },
              ],
            };

            return (
              <tr key={sim.id}>
                <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-gray-200">
                  {sim.id}
                </td>
                <td className="intersectionCell px-4 py-3 whitespace-wrap text-sm text-gray-900 dark:text-gray-200">
                  {sim.intersection}
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-gray-200">
                  {sim.avgWaitTime.toFixed(1)}
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-gray-200">
                  {sim.vehicleThroughput}
                </td>
                <td className="chartCell px-4 py-3 whitespace-nowrap text-sm">
                  <div className="h-16 w-24">
                    <Bar data={chartData} options={chartOptions} />
                  </div>
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-sm">
                  <span
                    className={`sim-status inline-flex items-center px-3 py-1 rounded-full border ${statusClass(sim.status)}`}
                  >
                    {sim.status}
                  </span>
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-sm">
                  <div className="flex flex-col space-y-2">
                    <button
                      onClick={() => handleViewResults(sim.id)}
                      className="viewBtn text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300 text-sm font-medium w-full text-center"
                      title="View Results"
                    >
                      <Eye size={18} strokeWidth={2} />
                    </button>
                    <button
                      onClick={() => handleDelete(sim.id)}
                      className="deleteBtn text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300 text-sm font-medium w-full text-center"
                      title="Delete Simulation"
                    >
                      <Trash2 size={18} strokeWidth={2} />
                    </button>
                  </div>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
      {simulations.length > rowsPerPage && (
        <div className="pagination absolute bottom-0 left-0 right-0 flex justify-center items-center p-4 space-x-2 bg-white dark:bg-[#161B22] border-t border-gray-200 dark:border-gray-700">
          <button
            onClick={() => handlePageChange(currentPage - 1)}
            disabled={currentPage === 0}
            className={`px-3 py-1 rounded-full text-sm font-medium bg-[#0F5BA7] dark:bg-[#388BFD] text-white hover:from-indigo-600 hover:to-indigo-700 dark:from-indigo-400 dark:to-indigo-500 dark:hover:from-indigo-500 dark:hover:to-indigo-600 transition-all duration-300 ${currentPage === 0 ? "opacity-50 cursor-not-allowed" : ""}`}
          >
            Prev
          </button>
          {Array.from({ length: totalPages }, (_, index) => (
            <button
              key={index}
              onClick={() => handlePageChange(index)}
              className={`px-3 py-1 rounded-full text-sm font-medium ${currentPage === index ? "bg-[#0F5BA7] text-white dark:bg-[#388BFD]" : "bg-gray-200 text-gray-700 dark:bg-gray-600 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-500"} transition-all duration-300`}
            >
              {index + 1}
            </button>
          ))}
          <button
            onClick={() => handlePageChange(currentPage + 1)}
            disabled={currentPage === totalPages - 1}
            className={`px-3 py-1 rounded-full text-sm font-medium bg-[#0F5BA7] dark:bg-[#388BFD] text-white hover:from-indigo-600 hover:to-indigo-700 dark:from-indigo-400 dark:to-indigo-500 dark:hover:from-indigo-500 dark:hover:to-indigo-600 transition-all duration-300 ${currentPage === totalPages - 1 ? "opacity-50 cursor-not-allowed" : ""}`}
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
};

const Simulations: React.FC = () => {
  const [filter1, setFilter1] = useState<string>("All Intersections");
  const [filter2, setFilter2] = useState<string>("All Intersections");
  const [page1, setPage1] = useState<number>(0);
  const [page2, setPage2] = useState<number>(0);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalType, setModalType] = useState<"simulations" | "optimizations">(
    "simulations",
  );

  const filteredSimulations1 =
    filter1 === "All Intersections"
      ? simulationsTable1
      : simulationsTable1.filter((sim) => sim.intersection === filter1);
  const filteredSimulations2 =
    filter2 === "All Intersections"
      ? simulationsTable2
      : simulationsTable2.filter((sim) => sim.intersection === filter2);
  const allIntersections = Array.from(
    new Set(
      [...simulationsTable1, ...simulationsTable2].map(
        (sim) => sim.intersection,
      ),
    ),
  );

  const handleNewSimulation = (type: "simulations" | "optimizations") => {
    setModalType(type);
    setIsModalOpen(true);
  };

  const handleModalSubmit = (data: {
    name: string;
    description: string;
    intersections: string[];
  }) => {
    console.log(
      `New ${modalType === "simulations" ? "Simulation" : "Optimization"} Created:`,
      data,
    );
    setIsModalOpen(false);
  };

  return (
    <div className="simulationBody min-h-screen bg-gray-100 dark:bg-gray-900">
      <Navbar />
      <div className="sim-main-content flex-grow p-6">
        <div className="simGrid grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="simTableContainer sims">
            <div className="flex items-center justify-between mb-4">
              <h1 className="text-3xl font-bold text-gray-800 dark:text-[#E6EDF3]">
                Recent Simulations
              </h1>
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => handleNewSimulation("simulations")}
                  className="new-simulation-button px-4 py-2 rounded-md text-sm font-medium bg-[#0F5BA7] dark:bg-[#388BFD] text-white hover:from-green-600 hover:to-green-700 dark:from-green-400 dark:to-green-500 dark:hover:from-green-500 dark:hover:to-green-600 transition-all duration-300 shadow-md hover:shadow-lg"
                >
                  New Simulation
                </button>
                <select
                  value={filter1}
                  onChange={(e) => {
                    setFilter1(e.target.value);
                    setPage1(0);
                  }}
                  className="w-48 p-2 rounded-md border border-gray-300 dark:border-[#388BFD] bg-white dark:bg-[#161B22] text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                  {[
                    "All Intersections",
                    ...new Set(
                      simulationsTable1.map((sim) => sim.intersection),
                    ),
                  ].map((intersection) => (
                    <option key={intersection} value={intersection}>
                      {intersection}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <SimulationTable
              simulations={filteredSimulations1}
              currentPage={page1}
              setCurrentPage={setPage1}
            />
          </div>
          <div className="simTableContainer opts">
            <div className="flex items-center justify-between mb-4">
              <h1 className="text-3xl font-bold text-gray-800 dark:text-[#E6EDF3]">
                Recent Optimizations
              </h1>
              <div className="flex items-center space-x-2">
                <select
                  value={filter2}
                  onChange={(e) => {
                    setFilter2(e.target.value);
                    setPage2(0);
                  }}
                  className="w-48 p-2 rounded-md border border-gray-300 dark:border-[#388BFD] bg-white dark:bg-[#161B22] text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                  {[
                    "All Intersections",
                    ...new Set(
                      simulationsTable2.map((sim) => sim.intersection),
                    ),
                  ].map((intersection) => (
                    <option key={intersection} value={intersection}>
                      {intersection}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <SimulationTable
              simulations={filteredSimulations2}
              currentPage={page2}
              setCurrentPage={setPage2}
            />
          </div>
        </div>
      </div>
      <Footer />
      <NewSimulationModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleModalSubmit}
        intersections={allIntersections}
        type={modalType}
      />
      <HelpMenu />
    </div>
  );
};

export default Simulations;
