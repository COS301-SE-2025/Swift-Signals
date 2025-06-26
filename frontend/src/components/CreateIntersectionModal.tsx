// src/components/CreateIntersectionModal.tsx

import React, { useState } from 'react';
import { X } from 'lucide-react';

const API_BASE_URL = "http://localhost:9090";

interface CreateIntersectionModalProps {
    isOpen: boolean;
    onClose: () => void;
    onIntersectionCreated: () => void; // Callback to refresh the list
}

const CreateIntersectionModal: React.FC<CreateIntersectionModalProps> = ({ isOpen, onClose, onIntersectionCreated }) => {
    // State for the form fields
    const [name, setName] = useState('');
    const [address, setAddress] = useState('');
    const [city, setCity] = useState('');
    const [province, setProvince] = useState('');
    const [trafficDensity, setTrafficDensity] = useState('medium');
    
    // State for API interaction
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (event: React.FormEvent) => {
        event.preventDefault();
        setIsLoading(true);
        setError(null);

        const token = localStorage.getItem('authToken');
        if (!token) {
            setError('Authentication error. Please log in again.');
            setIsLoading(false);
            return;
        }

        const requestBody = {
            name,
            details: { address, city, province },
            traffic_density: trafficDensity,
            // Including default parameters as per the API spec
            default_parameters: {
                green: 10,
                red: 6,
                yellow: 2,
                intersection_type: "t-junction",
                seed: Math.floor(Math.random() * 10000000000),
                speed: 60
            }
        };

        try {
            const response = await fetch(`${API_BASE_URL}/intersections`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(requestBody)
            });

            const data = await response.json();
            if (!response.ok) {
                throw new Error(data.message || 'Failed to create intersection.');
            }
            
            // Success!
            onIntersectionCreated(); // Trigger refresh on the parent component
            onClose(); // Close the modal

        } catch (err: any) {
            setError(err.message);
        } finally {
            setIsLoading(false);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50 p-4">
            <div className="bg-white p-8 rounded-lg shadow-2xl w-full max-w-lg relative">
                <button onClick={onClose} className="absolute top-4 right-4 text-gray-500 hover:text-gray-800">
                    <X size={24} />
                </button>
                <h2 className="text-3xl font-bold mb-6">Create New Intersection</h2>
                
                {error && <div className="bg-red-100 text-red-700 p-3 rounded-md mb-4">{error}</div>}

                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label htmlFor="name" className="block text-sm font-medium text-gray-700">Intersection Name</label>
                        <input type="text" id="name" value={name} onChange={e => setName(e.target.value)} className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500" required />
                    </div>
                     <div>
                        <label htmlFor="address" className="block text-sm font-medium text-gray-700">Address</label>
                        <input type="text" id="address" value={address} onChange={e => setAddress(e.target.value)} className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500" required />
                    </div>
                     <div>
                        <label htmlFor="city" className="block text-sm font-medium text-gray-700">City</label>
                        <input type="text" id="city" value={city} onChange={e => setCity(e.target.value)} className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500" required />
                    </div>
                    <div>
                        <label htmlFor="province" className="block text-sm font-medium text-gray-700">Province</label>
                        <input type="text" id="province" value={province} onChange={e => setProvince(e.target.value)} className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500" required />
                    </div>
                     <div>
                        <label htmlFor="trafficDensity" className="block text-sm font-medium text-gray-700">Traffic Density</label>
                        <select id="trafficDensity" value={trafficDensity} onChange={e => setTrafficDensity(e.target.value)} className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500">
                            <option value="low">Low</option>
                            <option value="medium">Medium</option>
                            <option value="high">High</option>
                        </select>
                    </div>
                    <div className="flex justify-end pt-4">
                        <button type="button" onClick={onClose} className="bg-gray-200 text-gray-800 font-medium py-2 px-4 rounded-md mr-2 hover:bg-gray-300">Cancel</button>
                        <button type="submit" disabled={isLoading} className="bg-red-700 hover:bg-red-800 text-white font-medium py-2 px-4 rounded-md disabled:bg-red-400">
                            {isLoading ? 'Creating...' : 'Create Intersection'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default CreateIntersectionModal;