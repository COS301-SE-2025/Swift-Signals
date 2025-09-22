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

export interface GeocodedLocation {
  name: string;
  city?: string;
  province?: string;
  lat: number | null;
  lon: number | null;
  osmId?: string | number;
  osmType?: string;
}

// Helper to check if a string is a valid street name
const isValidStreetName = (name: string): boolean => {
  if (!name || typeof name !== "string") return false;

  const cleaned = name.trim();
  if (cleaned.length < 3) return false;

  const streetSuffixes =
    /(street|road|avenue|drive|lane|way|boulevard|crescent|place|close|court|grove|gardens|park|square|circle|terrace|highway|route)$/i;
  const isNamedRoad = /^[A-Za-z][A-Za-z0-9\s-'.]*$/i.test(cleaned);

  return streetSuffixes.test(cleaned) || (isNamedRoad && cleaned.length >= 5);
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

// Nominatim search function
const searchNominatim = async (
  query: string,
  type: string = "road",
): Promise<NominatimResult[]> => {
  try {
    const params = new URLSearchParams({
      format: "json",
      addressdetails: "1",
      limit: "1", // We only need one result for geocoding an address
      countrycodes: "za",
      q: query,
      extratags: "1",
    });

    if (type === "road") {
      params.append("class", "highway");
    }

    const url = `https://nominatim.openstreetmap.org/search?${params.toString()}`;
    const response = await fetch(url);
    if (!response.ok) {
      console.warn(`Nominatim request failed: ${response.status}`);
      return [];
    }

    const data: NominatimResult[] = await response.json();
    return data;
  } catch (error) {
    console.error("Error in searchNominatim:", error);
    return [];
  }
};

// Process Nominatim results to a more usable Street format
const processNominatimResults = (
  data: NominatimResult[],
): GeocodedLocation[] => {
  if (!Array.isArray(data)) return [];

  const processed = data
    .filter((item: NominatimResult) => {
      const isInSA =
        item.address?.country_code === "za" ||
        item.display_name?.includes("South Africa");

      const isRoad =
        item.class === "highway" ||
        item.type === "road" ||
        item.address?.road ||
        /\b(street|road|avenue|drive|lane|way|boulevard|crescent|highway)\b/i.test(
          item.display_name,
        );

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
      const streetName =
        item.address?.road ||
        item.name ||
        item.display_name.split(",")[0].trim();
      const coordinates = getCoordinates(item);

      return {
        name: cleanStreetName(streetName),
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
    .filter((loc: GeocodedLocation) => loc.name && loc.name.length > 0)
    .filter(
      (loc: GeocodedLocation, index: number, self: GeocodedLocation[]) =>
        self.findIndex(
          (s) =>
            s.name.toLowerCase() === loc.name.toLowerCase() &&
            s.city?.toLowerCase() === loc.city?.toLowerCase(),
        ) === index,
    );

  return processed;
};

/**
 * Geocodes an address string to latitude and longitude using Nominatim.
 * @param address The address string to geocode.
 * @returns A Promise that resolves to GeocodedLocation (with lat/lon) or null if not found.
 */
export const geocodeAddress = async (
  address: string,
): Promise<GeocodedLocation | null> => {
  if (!address || address.length < 3) return null;

  try {
    const results = await searchNominatim(address, "road");
    const processed = processNominatimResults(results);
    if (
      processed.length > 0 &&
      processed[0].lat !== null &&
      processed[0].lon !== null
    ) {
      return processed[0];
    }
    return null;
  } catch (error) {
    console.error("Error geocoding address:", error);
    return null;
  }
};

/**
 * Generates a static map image URL from OpenStreetMap.
 * @param lat Latitude.
 * @param lon Longitude.
 * @param zoom Zoom level (e.g., 15-18 for street level).
 * @param width Image width in pixels.
 * @param height Image height in pixels.
 * @returns The URL for the static map image.
 */
export const getStaticMapUrl = (
  lat: number,
  lon: number,
  zoom: number = 17,
  width: number = 300,
  height: number = 200,
): string => {
  // Using staticmap.openstreetmap.de for simplicity.
  // Consider self-hosting a tile server or using a commercial static map API for production.
  return `https://staticmap.openstreetmap.de/staticmap.php?center=${lat},${lon}&zoom=${zoom}&size=${width}x${height}&maptype=mapnik&markers=${lat},${lon},red-pushpin`;
};
