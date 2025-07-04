import React from "react";
import "../styles/IntersectionCard.css";

interface IntersectionCardProps {
  id: string;
  name: string;
  location: string;
  lanes: string;
  image?: string;
  onSimulate: (id: number) => void;
  onEdit: (id: number) => void;
  onDelete: (id: number) => void;
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
  return (
    <div className="intersectionCard bg-white p-8 rounded-2xl shadow-lg flex justify-between items-center">
      <div className="flex items-center space-x-8">
        <div className="cardImage w-36 h-36 bg-gray-200 rounded-lg flex items-center justify-center">
          {image ? (
            <img
              src={image}
              alt={name}
              className="w-full h-full object-cover rounded-lg"
            />
          ) : (
            <svg
              className="cardImageSVG w-20 h-20 text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
          )}
        </div>

        <div>
          <h3 className="intersectionName text-3xl font-extrabold text-black mb-2">
            {name}
          </h3>
          <p className="intersectionID text-xl text-gray-700">ID: {id}</p>
          <p className="intersectionLocation text-xl text-gray-700">
            Location: {location}
          </p>
          <p className="intersectionLanes text-xl text-gray-700">
            Type: {lanes} 
          </p>
        </div>
      </div>

      <div className="intBtns flex flex-col space-y-3">
        <button
          onClick={() => onSimulate(Number(id))}
          className="intersectionBtn bg-blue-600 hover:bg-blue-700 text-white text-lg font-semibold py-2 px-8 rounded-full"
        >
          ▶ Simulate
        </button>

        <button
          onClick={() => onEdit(Number(id))}
          className="intersectionBtn bg-green-600 hover:bg-green-700 text-white text-lg font-semibold py-2 px-8 rounded-full"
        >
          ✏️ Edit
        </button>

        <button
          onClick={() => onDelete(Number(id))}
          className="intersectionBtn bg-red-600 hover:bg-red-700 text-white text-lg font-semibold py-2 px-8 rounded-full"
        >
          🗑️ Delete
        </button>
      </div>
      <div className="mobileIntBtns flex flex-col space-y-3">
        <button
          onClick={() => onSimulate(Number(id))}
          className="intersectionBtn bg-blue-600 hover:bg-blue-700 text-white text-lg font-semibold py-1 px-2 rounded-full"
        >
          ▶
        </button>

        <button
          onClick={() => onEdit(Number(id))}
          className="intersectionBtn bg-green-600 hover:bg-green-700 text-white text-lg font-semibold py-1 px-2 rounded-full"
        >
          ✏️
        </button>
        <button
          onClick={() => onDelete(Number(id))}
          className="intersectionBtn bg-red-600 hover:bg-red-700 text-white text-lg font-semibold py-1 px-2 rounded-full"
        >
          🗑️
        </button>
      </div>
    </div>
  );
};

export default IntersectionCard;
