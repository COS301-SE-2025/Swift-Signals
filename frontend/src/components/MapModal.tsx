import React from "react";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import L from "leaflet";
import { X } from "lucide-react";
import "../styles/MapModal.css";
import markerIcon2x from "leaflet/dist/images/marker-icon-2x.png";
import markerIcon from "leaflet/dist/images/marker-icon.png";
import markerShadow from "leaflet/dist/images/marker-shadow.png";

// --- FIX for default icon issue with webpack ---
// Deletes the default icon's image paths and resets them after the component is loaded.
delete (L.Icon.Default.prototype as unknown as Record<string, unknown>)
  ._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: markerIcon2x,
  iconUrl: markerIcon,
  shadowUrl: markerShadow,
});
// --- END FIX ---

interface Intersection {
  id: string;
  name: string;
  details: {
    address: string;
    city: string;
    province: string;
    latitude: number;
    longitude: number;
  };
}

interface MapModalProps {
  isOpen: boolean;
  onClose: () => void;
  intersections: Intersection[];
  onSimulate: (id: string, name: string) => void;
}

const MapModal: React.FC<MapModalProps> = ({
  isOpen,
  onClose,
  intersections,
  onSimulate,
}) => {
  if (!isOpen) return null;

  // Center map on Pretoria by default
  const defaultPosition: [number, number] = [-25.7479, 28.2293];

  return (
    <div className="map-modal-overlay">
      <div className="map-modal-content">
        <div className="map-modal-header">
          <h2 className="map-modal-title">Intersections Map</h2>
          <button onClick={onClose} className="map-modal-close-btn">
            <X size={24} />
          </button>
        </div>
        <div className="map-container-wrapper">
          <MapContainer
            center={defaultPosition}
            zoom={12}
            className="map-container"
          >
            <TileLayer
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
              attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            />
            {intersections.map((intersection) => (
              <Marker
                key={intersection.id}
                position={[
                  intersection.details.latitude,
                  intersection.details.longitude,
                ]}
              >
                <Popup>
                  <b>{intersection.name}</b>
                  <br />
                  {intersection.details.address.split(',')[0]}
                  <br />
                  <button 
                    className="mt-2 px-3 py-1 bg-[#0F5BA7] text-white text-xs rounded hover:bg-blue-700 transition-colors"
                    onClick={() => onSimulate(intersection.id, intersection.name)}>
                    Simulate
                  </button>
                </Popup>
              </Marker>
            ))}
          </MapContainer>
        </div>
      </div>
    </div>
  );
};

export default MapModal;
