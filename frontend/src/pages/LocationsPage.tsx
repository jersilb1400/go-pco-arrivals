import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { motion } from 'framer-motion';
import { 
  Users, 
  Plus, 
  Edit, 
  Trash2, 
  Eye,
  MapPin,
  Calendar,
  AlertCircle
} from 'lucide-react';
import apiService from '../services/api';
import type { Location } from '../types/api';

interface AddLocationForm {
  pco_location_id: string;
  name: string;
}

const LocationsPage: React.FC = () => {
  const queryClient = useQueryClient();
  const [showAddForm, setShowAddForm] = useState(false);
  const [formData, setFormData] = useState<AddLocationForm>({
    pco_location_id: '',
    name: '',
  });

  const { data: locationsResponse, isLoading } = useQuery({
    queryKey: ['locations'],
    queryFn: apiService.getLocations,
  });

  // PCO Integration Test
  const { data: pcoLocationsResponse, isLoading: pcoLoading, refetch: refetchPCO } = useQuery({
    queryKey: ['pco-locations'],
    queryFn: apiService.getPCOLocations,
    enabled: false, // Don't fetch automatically
  });

  const addLocationMutation = useMutation({
    mutationFn: (data: AddLocationForm) => 
      apiService.addLocation(data.pco_location_id, data.name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['locations'] });
      setShowAddForm(false);
      setFormData({ pco_location_id: '', name: '' });
    },
    onError: (error) => {
      console.error('Failed to add location:', error);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.pco_location_id && formData.name) {
      addLocationMutation.mutate(formData);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const locations = locationsResponse?.locations || [];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Locations</h1>
          <p className="mt-1 text-sm text-gray-500">
            Manage your billboard locations
          </p>
        </div>
        
        <div className="flex space-x-3">
          <button
            onClick={() => refetchPCO()}
            disabled={pcoLoading}
            className="btn-secondary"
          >
            {pcoLoading ? 'Loading...' : 'Test PCO Integration'}
          </button>
          <button
            onClick={() => setShowAddForm(true)}
            className="btn-primary"
          >
            <Plus className="h-4 w-4 mr-2" />
            Add Location
          </button>
        </div>
      </div>

      {/* Add Location Form */}
      {showAddForm && (
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -20 }}
          className="card p-6"
        >
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-medium text-gray-900">Add New Location</h2>
            <button
              onClick={() => setShowAddForm(false)}
              className="text-gray-400 hover:text-gray-600"
            >
              <AlertCircle className="h-5 w-5" />
            </button>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="pco_location_id" className="block text-sm font-medium text-gray-700">
                PCO Location ID
              </label>
              <input
                type="text"
                id="pco_location_id"
                name="pco_location_id"
                value={formData.pco_location_id}
                onChange={handleInputChange}
                className="input mt-1"
                placeholder="Enter PCO location ID"
                required
              />
              <p className="mt-1 text-xs text-gray-500">
                The location ID from Planning Center Online
              </p>
            </div>

            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                Display Name
              </label>
              <input
                type="text"
                id="name"
                name="name"
                value={formData.name}
                onChange={handleInputChange}
                className="input mt-1"
                placeholder="Enter display name"
                required
              />
              <p className="mt-1 text-xs text-gray-500">
                The name that will appear on the billboard
              </p>
            </div>

            <div className="flex items-center justify-end space-x-3 pt-4">
              <button
                type="button"
                onClick={() => setShowAddForm(false)}
                className="btn-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={addLocationMutation.isPending}
                className="btn-primary"
              >
                {addLocationMutation.isPending ? 'Adding...' : 'Add Location'}
              </button>
            </div>
          </form>
        </motion.div>
      )}

      {/* PCO Integration Test Results */}
      {pcoLocationsResponse && (
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="card p-6 bg-blue-50 border-blue-200"
        >
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-medium text-blue-900">PCO Integration Test Results</h2>
            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
              Success
            </span>
          </div>
          
          <div className="space-y-3">
            <p className="text-sm text-blue-700">
              Successfully fetched {pcoLocationsResponse.locations?.length || 0} locations from PCO API
            </p>
            
            {pcoLocationsResponse.locations && pcoLocationsResponse.locations.length > 0 && (
              <div className="bg-white rounded-lg p-4 border">
                <h3 className="text-sm font-medium text-gray-900 mb-2">PCO Locations:</h3>
                <div className="space-y-2">
                  {pcoLocationsResponse.locations.slice(0, 5).map((location: any, index: number) => (
                    <div key={index} className="flex items-center justify-between text-sm">
                      <span className="text-gray-700">{location.name}</span>
                      <span className="text-gray-500">ID: {location.id}</span>
                    </div>
                  ))}
                  {pcoLocationsResponse.locations.length > 5 && (
                    <p className="text-xs text-gray-500">
                      ... and {pcoLocationsResponse.locations.length - 5} more locations
                    </p>
                  )}
                </div>
              </div>
            )}
          </div>
        </motion.div>
      )}

      {/* Locations List */}
      <div className="card">
        {isLoading ? (
          <div className="p-6 text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
            <p className="mt-2 text-sm text-gray-500">Loading locations...</p>
          </div>
        ) : locations.length === 0 ? (
          <div className="p-12 text-center">
            <MapPin className="mx-auto h-12 w-12 text-gray-400" />
            <h3 className="mt-2 text-sm font-medium text-gray-900">No locations</h3>
            <p className="mt-1 text-sm text-gray-500">
              Get started by adding your first location.
            </p>
            <div className="mt-6">
              <button
                onClick={() => setShowAddForm(true)}
                className="btn-primary"
              >
                <Plus className="h-4 w-4 mr-2" />
                Add Location
              </button>
            </div>
          </div>
        ) : (
          <div className="overflow-hidden">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Location
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    PCO ID
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Created
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {locations.map((location) => (
                  <motion.tr
                    key={location.id}
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="hover:bg-gray-50"
                  >
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <div className="flex-shrink-0 h-10 w-10">
                          <div className="h-10 w-10 rounded-full bg-primary-100 flex items-center justify-center">
                            <MapPin className="h-5 w-5 text-primary-600" />
                          </div>
                        </div>
                        <div className="ml-4">
                          <div className="text-sm font-medium text-gray-900">
                            {location.name}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        {location.pco_location_id}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        location.is_active
                          ? 'bg-green-100 text-green-800'
                          : 'bg-gray-100 text-gray-800'
                      }`}>
                        {location.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      <div className="flex items-center">
                        <Calendar className="h-4 w-4 mr-1" />
                        {new Date(location.created_at).toLocaleDateString()}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <div className="flex items-center justify-end space-x-2">
                        <Link
                          to={`/billboard/${location.pco_location_id}`}
                          className="text-primary-600 hover:text-primary-900"
                          title="View Billboard"
                        >
                          <Eye className="h-4 w-4" />
                        </Link>
                        <button
                          className="text-gray-600 hover:text-gray-900"
                          title="Edit Location"
                        >
                          <Edit className="h-4 w-4" />
                        </button>
                        <button
                          className="text-danger-600 hover:text-danger-900"
                          title="Delete Location"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </motion.tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Stats */}
      {locations.length > 0 && (
        <div className="grid grid-cols-1 gap-5 sm:grid-cols-3">
          <div className="card p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0 bg-blue-100 rounded-md p-3">
                <MapPin className="h-6 w-6 text-blue-600" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Total Locations
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {locations.length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>

          <div className="card p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0 bg-green-100 rounded-md p-3">
                <Users className="h-6 w-6 text-green-600" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Active Locations
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {locations.filter(l => l.is_active).length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>

          <div className="card p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0 bg-purple-100 rounded-md p-3">
                <Calendar className="h-6 w-6 text-purple-600" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Recently Added
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {locations.filter(l => {
                      const created = new Date(l.created_at);
                      const weekAgo = new Date();
                      weekAgo.setDate(weekAgo.getDate() - 7);
                      return created > weekAgo;
                    }).length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default LocationsPage; 