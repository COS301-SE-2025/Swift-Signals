import { useState } from 'react';
import Navbar from '../components/Navbar';
import { Search } from 'lucide-react';
import IntersectionCard from '../components/IntersectionCard';
import '../styles/Intersections.css';

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
    <div className="intersectionBody min-h-screen bg-gray-100 overflow-y-auto">
      <Navbar />
    <div className="main-content flex-grow">
      <div className="max-w-6xl mx-auto px-4 py-8">
        <div className="topBar flex justify-between items-center mb-6">
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
            onClick={() => console.log('Add Intersection')}
            className="addIntersectionBtn bg-red-700 hover:bg-red-800 text-white font-medium py-2 px-4 rounded-md"
          >
            Add Intersection
          </button>
        </div>

        <div className="intersections space-y-6 overflow-y-auto max-h-[calc(100vh-0px)] pr-2">
          {filteredIntersections.map((intersection) => (
            <IntersectionCard
              key={intersection.id}
              {...intersection}
              onSimulate={(id) => console.log(`Simulate ${id}`)}
              onEdit={(id) => console.log(`Edit ${id}`)}
              onDelete={(id) => console.log(`Delete ${id}`)}
            />
          ))}
        </div>
      </div>
      </div>
    </div>
    
  );
};

export default Intersections;
