interface IntersectionCardProps {
  id: number;
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
    <div className="bg-white p-8 rounded-2xl shadow-lg flex justify-between items-center min-h-[220px]">
      <div className="flex items-center space-x-8">
        <div className="w-36 h-36 bg-gray-200 rounded-lg flex items-center justify-center">
          {image ? (
            <img src={image} alt={name} className="w-full h-full object-cover rounded-lg" />
          ) : (
            <svg className="w-20 h-20 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          )}
        </div>

        <div>
          <h3 className="text-3xl font-extrabold text-black mb-2">{name}</h3>
          <p className="text-xl text-gray-700">ID: {id}</p>
          <p className="text-xl text-gray-700">Location: {location}</p>
          <p className="text-xl text-gray-700">Lanes: {lanes}</p>
        </div>
      </div>

      <div className="flex flex-col space-y-3">
        <button
          onClick={() => onSimulate(id)}
          className="bg-blue-600 hover:bg-blue-700 text-white text-lg font-semibold py-2 px-8 rounded-full"
        >
          ▶ Simulate
        </button>

        <button
          onClick={() => onEdit(id)}
          className="bg-green-600 hover:bg-green-700 text-white text-lg font-semibold py-2 px-8 rounded-full"
        >
          ✏️ Edit
        </button>

        <button
          onClick={() => onDelete(id)}
          className="bg-red-600 hover:bg-red-700 text-white text-lg font-semibold py-2 px-8 rounded-full"
        >
          🗑️ Delete
        </button>
      </div>
    </div>
  );
};

export default IntersectionCard;
