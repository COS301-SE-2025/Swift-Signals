import { useState } from 'react';
import Navbar from '../components/Navbar';
import UsersTable from '../components/UsersTable';

// TypeScript interface for user data
interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  lastLogin: string;
}

const Users = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages] = useState(1);

  // Sample user data
  const users: User[] = [
    { id: 1, name: 'John Doe', email: 'email@email.com', role: 'Admin', lastLogin: '2025-05-13 09:00' },
    { id: 2, name: 'Jane Smith', email: 'email@email.com', role: 'Engineer', lastLogin: '2025-05-13 09:00' },
    { id: 3, name: 'John Calvin', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
    { id: 4, name: 'Paul Washer', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
    { id: 5, name: 'Joshua Garner', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
    { id: 6, name: 'Chris Xides', email: 'email@email.com', role: 'Engineer', lastLogin: '2025-05-13 09:00' },
    { id: 7, name: 'Kgosi Segale', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
    { id: 8, name: 'John Flavel', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
    { id: 9, name: 'John Owen', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
  ];

  const handleEdit = (id: number) => console.log(`Edit user ${id}`);
  const handleDelete = (id: number) => console.log(`Delete user ${id}`);

  const goToNextPage = () => {
    if (currentPage < totalPages) setCurrentPage(currentPage + 1);
  };

  const goToPreviousPage = () => {
    if (currentPage > 1) setCurrentPage(currentPage - 1);
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <Navbar />
      <div className="max-w-6xl mx-auto px-4 py-8">
        <UsersTable 
          users={users}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
        
        <div className="flex justify-center items-center py-4 gap-4 mt-4">
          <button
            onClick={goToPreviousPage}
            disabled={currentPage === 1}
            className="px-3 py-1 border rounded-md hover:bg-gray-100 disabled:opacity-50 text-black"
          >
            &lt;
          </button>
          <span className="text-black font-medium">Page {currentPage} of {totalPages}</span>
          <button
            onClick={goToNextPage}
            disabled={currentPage === totalPages}
            className="px-3 py-1 border rounded-md hover:bg-gray-100 disabled:opacity-50 text-black"
          >
            &gt;
          </button>
        </div>
      </div>
    </div>
  );
};

export default Users;
