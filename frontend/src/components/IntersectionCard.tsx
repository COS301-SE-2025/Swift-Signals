import { PlayCircle, PencilLine, Trash2 } from "lucide-react";
import React from "react";
import "../styles/IntersectionCard.css";
import placeholderImage from "../assets/placeholder.png";

interface IntersectionCardProps {
  id: string;
  name: string;
  location: string;
  lanes: string;
  image?: string;
  onSimulate: (id: string) => void;
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
}

const IntersectionCard: React.FC<IntersectionCardProps> = ({
  id,
  name,
  location,
  lanes,
  image,
  onSimulate,
  onEdit,
  onDelete,
}) => {
  const displayName = name.split(' [')[0];
  const displayLocation = location.split(',')[0];
  return (
    <div className="intersectionCard bg-white p-8 rounded-2xl shadow-lg flex justify-between items-center">
      <div className="flex items-center space-x-8">
        <div className="cardImage w-36 h-36 rounded-lg flex items-center justify-center">
          <img
            src={image || placeholderImage}
            alt={name}
            className="w-full h-32 object-contain rounded-t-lg"
          />
        </div>

        <div>
          <h3 className="intersectionName text-3xl font-extrabold text-black dark:text-[#E6EDF3] mb-2">
            {displayName}
          </h3>
          {/* <p className="intersectionID text-xl text-gray-700 dark:text-[#8B949E]">
            ID: {id}
          </p> */}
          <p className="intersectionLocation text-xl text-gray-700 dark:text-[#8B949E]">
            Location: {displayLocation}
          </p>
          <p className="intersectionLanes text-xl text-gray-700 dark:text-[#8B949E]">
            Type: {lanes}
          </p>
        </div>
      </div>

      <div className="intBtns flex flex-col space-y-3">
        <button
          onClick={() => onSimulate(id)}
          className="simButton intersectionBtn bg-[#0F5BA7] dark:bg-[#388BFD] hover:bg-blue-700 text-white text-lg font-semibold py-2 px-8 rounded-full flex items-center justify-center gap-2"
        >
          <PlayCircle size={22} strokeWidth={2} /> Simulate
        </button>

        <button
          onClick={() => onEdit(id)}
          className="editButton intersectionBtn bg-[#2B9348] dark:bg-[#2DA44E] hover:bg-green-700 text-white text-lg font-semibold py-2 px-8 rounded-full flex items-center justify-center gap-2"
        >
          <PencilLine size={22} strokeWidth={2} /> Edit
        </button>

        <button
          onClick={() => onDelete(id)}
          className="deleteIntButton intersectionBtn bg-red-600 dark:bg-[#DA3633] hover:bg-red-700 text-white text-lg font-semibold py-2 px-8 rounded-full flex items-center justify-center gap-2"
        >
          <Trash2 size={22} strokeWidth={2} /> Delete
        </button>
      </div>
      <div className="mobileIntBtns flex flex-col space-y-3">
        <button
          onClick={() => onSimulate(id)}
          className="intersectionBtn bg-[#0F5BA7] dark:bg-[#388BFD] hover:bg-blue-700 text-white text-lg font-semibold py-1 px-2 rounded-full flex items-center justify-center"
        >
          <PlayCircle size={18} strokeWidth={2} />
        </button>

        <button
          onClick={() => onEdit(id)}
          className="intersectionBtn bg-[#2B9348] dark:bg-[#2DA44E] hover:bg-green-700 text-white text-lg font-semibold py-1 px-2 rounded-full flex items-center justify-center"
        >
          <PencilLine size={18} strokeWidth={2} />
        </button>
        <button
          onClick={() => onDelete(id)}
          className="intersectionBtn bg-red-600 dark:bg-[#DA3633] hover:bg-red-700 text-white text-lg font-semibold py-1 px-2 rounded-full flex items-center justify-center"
        >
          <Trash2 size={18} strokeWidth={2} />
        </button>
      </div>
    </div>
  );
};

export default IntersectionCard;
