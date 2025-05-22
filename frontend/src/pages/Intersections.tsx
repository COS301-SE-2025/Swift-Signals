import { useState } from 'react';
import Navbar from '../components/Navbar';
import { Search } from 'lucide-react';
import IntersectionCard from '../components/IntersectionCard';

// TypeScript interface for intersection data
interface Intersection {
  id: number;
  name: string;
  location: string;
  lanes: string;
  image?: string;
}

const Intersections = () => {
  const [searchQuery, setSearchQuery] = useState('');
  
  // Sample data for intersections
  const intersections: Intersection[] = [
    {
      id: 1,
      name: 'Main St & 1st Ave',
      location: 'Pretoria CBD',
      lanes: '4-way, 2 lanes each',
    },
    {
      id: 2,
      name: 'Church St & Park Rd',
      location: 'Hatfield',
      lanes: '3-way, 1 lane each',
    },
    {
      id: 3,
      name: 'University Rd & Lynnwood',
      location: 'Hatfield',
      lanes: '4-way, 2 lanes each',
    },
  ];

  const filteredIntersections = intersections.filter(intersection => 
    intersection.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    intersection.id.toString().includes(searchQuery)
  );

  return (
    <div className="min-h-screen bg-gray-100">
      <Navbar />
      
      <div className="max-w-6xl mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-6">
          <div className="relative w-full max-w-md">
            <input
              type="text"
              placeholder="Search by Name or ID..."
              className="w-full pl-4 pr-10 py-2 border border-gray-300 rounded-full focus:outline-none focus:ring-2 focus:ring-red-500"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
            <div className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500">
              <Search size={20} />
            </div>
          </div>
          
          <button
            onClick={handleAddIntersection}
            className="bg-red-700 hover:bg-red-800 text-white font-medium py-2 px-4 rounded-md"
          >
            Add Intersection
          </button>
        </div>
        
        <div className="space-y-4">
          {filteredIntersections.map((intersection) => (
            <div 
              key={intersection.id}
              className="bg-white p-4 rounded-lg shadow-sm flex justify-between items-center"
            >
              <div className="flex items-center space-x-4">
                <div className="w-24 h-24 bg-gray-200 rounded-md flex items-center justify-center">
                  {intersection.image ? (
                    <img 
                      src={intersection.image} 
                      alt={intersection.name} 
                      className="w-full h-full object-cover rounded-md"
                    />
                  ) : (
                    <div className="text-gray-400">
                      <svg className="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                      </svg>
                    </div>
                  )}
                </div>
                
                <div>
                  <h3 className="text-xl font-bold">{intersection.name}</h3>
                  <p className="text-gray-600">ID: {intersection.id}</p>
                  <p className="text-gray-600">Location: {intersection.location}</p>
                  <p className="text-gray-600">Lanes: {intersection.lanes}</p>
                </div>
              </div>
              
              <div className="flex flex-col space-y-2">
                <button
                  onClick={() => handleSimulate(intersection.id)}
                  className="bg-blue-500 hover:bg-blue-600 text-white font-medium py-1 px-6 rounded-full flex items-center justify-center"
                >
                  <span className="mr-1">▶</span> Simulate
                </button>
                
                <button
                  onClick={() => handleEdit(intersection.id)}
                  className="bg-green-500 hover:bg-green-600 text-white font-medium py-1 px-6 rounded-full flex items-center justify-center"
                >
                  <span className="mr-1">✏️</span> Edit
                </button>
                
                <button
                  onClick={() => handleDelete(intersection.id)}
                  className="bg-red-500 hover:bg-red-600 text-white font-medium py-1 px-6 rounded-full flex items-center justify-center"
                >
                  <span className="mr-1">🗑️</span> Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Intersections;
