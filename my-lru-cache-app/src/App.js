import React, { useEffect, useState } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import CacheTable from './CacheTable';
import CapacityModal from './CapacityModal';
import AddEntryModal from './AddEntryModal';
import GetEntryModal from './GetEntryModal';

function App() {
  const [cacheData, setCacheData] = useState(null);
  const [capacity, setCapacity] = useState(5);
  const [newCapacity, setNewCapacity] = useState(5);
  const [showAddModal, setShowAddModal] = useState(false);
  const [showGetModal, setShowGetModal] = useState(false);
  const [showCapacityModal, setShowCapacityModal] = useState(false);
  const [newKey, setNewKey] = useState('');
  const [newValue, setNewValue] = useState('');
  const [expiration, setExpiration] = useState('');
  const [getKey, setGetKey] = useState('');

  useEffect(() => {
    const socket = new WebSocket('ws://localhost:8080/ws');

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        setCacheData(data);
      } catch (error) {
        console.error('Error parsing WebSocket data:', error);
      }
    };

    return () => {
      socket.close();
    };
  }, []);

  const handleCapacityChange = (event) => {
    setNewCapacity(event.target.value);
  };

  const handleCapacitySubmit = async () => {
    try {
      const response = await fetch('http://localhost:8080/cache/capacity', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ capacity: parseInt(newCapacity, 10) }),
      });
      const result = await response.json();
      if (response.ok) {
        setCapacity(result.capacity);
        alert(`Capacity set to ${result.capacity}`);
        setShowCapacityModal(false);
      } else {
        alert(`Error: ${result.message}`);
      }
    } catch (error) {
      console.error('Error setting cache capacity:', error);
    }
  };

  const handleAddEntry = async () => {
    try {
      if (newKey.trim() === '' || newValue.trim() === '' || expiration.trim() === '') {
        alert('Please fill in all fields.');
        return;
      }

      const requestBody = {
        key: parseInt(newKey, 10),
        value: parseInt(newValue, 10),
        expiration: parseInt(expiration, 10) || 0,
      };

      const response = await fetch('http://localhost:8080/cache', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody),
      });

      const result = await response.json();

      if (response.ok) {
        setShowAddModal(false);
        setNewKey('');
        setNewValue('');
        setExpiration('');
        alert('Cache entry added');
      } else {
        alert(`Error: ${result.message}`);
      }
    } catch (error) {
      console.error('Error adding cache entry:', error);
    }
  };

  const handleGetEntry = async () => {
    try {
      const response = await fetch(`http://localhost:8080/cache/${parseInt(getKey, 10)}`);
      const result = await response.json();
      if (response.ok) {
        alert(`Value: ${result.value}`);
        setShowGetModal(false);
      } else {
        alert(`Error: ${result.message}`);
      }
    } catch (error) {
      console.error('Error getting cache entry:', error);
    }
  };

  return (
    <div className="container mt-5">
      <h1>Cache Data</h1>
      <div className="d-flex justify-content-start mb-4">
        <button
          className="btn btn-primary btn-lg mx-2"
          onClick={() => setShowCapacityModal(true)}
        >
          Set Capacity
        </button>
        <button
          className="btn btn-success btn-lg mx-2"
          onClick={() => setShowAddModal(true)}
        >
          Add Cache Entry
        </button>
        <button
          className="btn btn-info btn-lg mx-2"
          onClick={() => setShowGetModal(true)}
        >
          Get Cache Entry
        </button>
      </div>
      <CacheTable cacheData={cacheData} />
      <CapacityModal
        show={showCapacityModal}
        handleClose={() => setShowCapacityModal(false)}
        handleChange={handleCapacityChange}
        handleSubmit={handleCapacitySubmit}
        newCapacity={newCapacity}
      />
      <AddEntryModal
        show={showAddModal}
        handleClose={() => setShowAddModal(false)}
        handleAddEntry={handleAddEntry}
        newKey={newKey}
        newValue={newValue}
        expiration={expiration}
        setNewKey={setNewKey}
        setNewValue={setNewValue}
        setExpiration={setExpiration}
      />
      <GetEntryModal
        show={showGetModal}
        handleClose={() => setShowGetModal(false)}
        handleGetEntry={handleGetEntry}
        getKey={getKey}
        setGetKey={setGetKey}
      />
    </div>
  );
}

export default App;
