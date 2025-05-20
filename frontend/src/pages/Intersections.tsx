import { useState } from 'react';
import Navbar from '../components/Navbar';
import { Search } from 'lucide-react';

// TypeScript interface for user data
interface User {
  id: number;
  name: string;
}

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

  const handleAddIntersection = () => {
    // Implementation for adding a new intersection
    console.log('Add intersection clicked');
  };

  const handleSimulate = (id: number) => {
    // Implementation for simulating the intersection
    console.log(`Simulate intersection ${id}`);
  };

  const handleEdit = (id: number) => {
    // Implementation for editing the intersection
    console.log(`Edit intersection ${id}`);
  };

  const handleDelete = (id: number) => {
    // Implementation for deleting the intersection
    console.log(`Delete intersection ${id}`);
  };

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

// const Users = () => {
//   const [currentPage, setCurrentPage] = useState(1);
//   const [totalPages] = useState(1);

//   // Sample user data
//   const users: User[] = [
//     { id: 1, name: 'John Doe', email: 'email@email.com', role: 'Admin', lastLogin: '2025-05-13 09:00' },
//     { id: 2, name: 'Jane Smith', email: 'email@email.com', role: 'Engineer', lastLogin: '2025-05-13 09:00' },
//     { id: 3, name: 'John Calvin', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
//     { id: 4, name: 'Paul Washer', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
//     { id: 5, name: 'Joshua Garner', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
//     { id: 6, name: 'Chris Xides', email: 'email@email.com', role: 'Engineer', lastLogin: '2025-05-13 09:00' },
//     { id: 7, name: 'Kgosi Segale', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
//     { id: 8, name: 'John Flavel', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
//     { id: 9, name: 'John Owen', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
//   ];

//   const handleEdit = (id: number) => console.log(`Edit user ${id}`);
//   const handleDelete = (id: number) => console.log(`Delete user ${id}`);

//   const goToNextPage = () => {
//     if (currentPage < totalPages) setCurrentPage(currentPage + 1);
//   };

//   const goToPreviousPage = () => {
//     if (currentPage > 1) setCurrentPage(currentPage - 1);
//   };

//   return (
//     <div className="min-h-screen bg-gray-100">
//       <Navbar />
//       <div className="max-w-6xl mx-auto px-4 py-8">
//         <div className="bg-white rounded-lg shadow-sm overflow-hidden">
//           <table className="w-full border-collapse">
//             <thead className="text-left border-b">
//               <tr>
//                 <th className="px-4 py-3 font-bold">ID</th>
//                 <th className="px-4 py-3 font-bold">Name</th>
//                 <th className="px-4 py-3 font-bold">Email</th>
//                 <th className="px-4 py-3 font-bold">Role</th>
//                 <th className="px-4 py-3 font-bold">Last Login</th>
//                 <th className="px-4 py-3 font-bold">Actions</th>
//               </tr>
//             </thead>
//             <tbody>
//               {users.map((user) => (
//                 <tr key={user.id} className="border-b hover:bg-gray-50">
//                   <td className="px-4 py-3">{user.id}</td>
//                   <td className="px-4 py-3">{user.name}</td>
//                   <td className="px-4 py-3">{user.email}</td>
//                   <td className="px-4 py-3">{user.role}</td>
//                   <td className="px-4 py-3">{user.lastLogin}</td>
//                   <td className="px-4 py-3">
//                     <div className="flex gap-2 justify-center">
//                       <button
//                         onClick={() => handleEdit(user.id)}
//                         className="p-2 bg-green-500 text-white rounded-full flex items-center justify-center"
//                         aria-label="Edit user"
//                       >
//                         <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
//                           strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
//                           <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
//                           <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
//                         </svg>
//                       </button>
//                       <button
//                         onClick={() => handleDelete(user.id)}
//                         className="p-2 bg-red-500 text-white rounded-full flex items-center justify-center"
//                         aria-label="Delete user"
//                       >
//                         <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
//                           strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
//                           <polyline points="3 6 5 6 21 6" />
//                           <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
//                         </svg>
//                       </button>
//                     </div>
//                   </td>
//                 </tr>
//               ))}
//             </tbody>
//           </table>

//           <div className="flex justify-center items-center py-4 gap-4">
//             <button
//               onClick={goToPreviousPage}
//               disabled={currentPage === 1}
//               className="px-3 py-1 border rounded-md hover:bg-gray-100 disabled:opacity-50"
//             >
//               &lt;
//             </button>
//             <span className="text-gray-700">Page {currentPage} of {totalPages}</span>
//             <button
//               onClick={goToNextPage}
//               disabled={currentPage === totalPages}
//               className="px-3 py-1 border rounded-md hover:bg-gray-100 disabled:opacity-50"
//             >
//               &gt;
//             </button>
//           </div>
//         </div>
//       </div>
//     </div>
//   );
// };

// export default Users;
