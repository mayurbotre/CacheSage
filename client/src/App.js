import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

const App = () => {
  const [key, setKey] = useState('');
  const [value, setValue] = useState('');
  const [expiration, setExpiration] = useState(5);
  const [cacheData, setCacheData] = useState({});
  const [error, setError] = useState('');

  const apiUrl = 'http://localhost:8080/cache';
  const wsUrl = 'ws://localhost:8080/cache-updates';

  useEffect(() => {
    fetchCacheData();

    const socket = new WebSocket(wsUrl);

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setCacheData(data);
    };

    socket.onerror = () => {};

    return () => {
      socket.close();
    };
  }, [wsUrl]);

  const fetchCacheData = async () => {
    try {
      const response = await axios.get(apiUrl);
      setCacheData(response.data);
    } catch (error) {
      console.error('Error fetching cache data:', error);
      setError('Failed to fetch cache data');
    }
  };

  const handleSet = async () => {
    if (!key || !value || !expiration) {
      setError('Please provide valid inputs.');
      return;
    }

    setError('');
    try {
      await axios.post(apiUrl, {
        key,
        value,
        expires_in: expiration,
      });
      setKey('');
      setValue('');
      setExpiration(5);
      fetchCacheData();
    } catch (error) {
      console.error('Error setting cache:', error);
      setError('Failed to set cache');
    }
  };

  const handleGet = async () => {
    if (!key) {
      setError('Please provide a key.');
      return;
    }

    setError('');
    try {
      const response = await axios.get(`${apiUrl}/${key}`);
      alert(`Value: ${response.data.value}`);
    } catch (error) {
      console.error('Error fetching cache:', error);
      setError('Failed to fetch cache');
    }
  };

  const handleDelete = async () => {
    if (!key) {
      setError('Please provide a key.');
      return;
    }

    setError('');
    try {
      await axios.delete(`${apiUrl}/${key}`);
      fetchCacheData();
    } catch (error) {
      console.error('Error deleting cache:', error);
      setError('Failed to delete cache');
    }
  };

  return (
    <div className="App">
      <header className="header">CacheSage</header>
      {Object.keys(cacheData).length > 0 ? (
        <div className="container">
          <div className="left-container">
            <div className="content">
              <h1>Set Data</h1>
              <div className="input-container">
                <input
                  type="text"
                  placeholder="Key"
                  value={key}
                  onChange={(e) => setKey(e.target.value)}
                />
                <input
                  type="text"
                  placeholder="Value"
                  value={value}
                  onChange={(e) => setValue(e.target.value)}
                />
                <input
                  type="number"
                  placeholder="Expiration (seconds)"
                  value={expiration}
                  onChange={(e) => setExpiration(parseInt(e.target.value, 10))}
                />
              </div>
              <div className="button-container">
                <button onClick={handleSet}>Set</button>
                <button onClick={handleGet}>Get</button>
                <button onClick={handleDelete}>Delete</button>
              </div>
              {error && <div className="error">{error}</div>}
            </div>
          </div>
          <div className="right-container">
            <div className="content">
              <h2>Cache Data</h2>
              <table>
                <thead>
                  <tr>
                    <th>Key</th>
                    <th>Value</th>
                    <th>Expires In (seconds)</th>
                  </tr>
                </thead>
                <tbody>
                  {Object.entries(cacheData).map(([key, { value, expiration }]) => (
                    <tr key={key}>
                      <td>{key}</td>
                      <td>{value}</td>
                      <td>{expiration} seconds</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      ) : (
        <div className="full-screen-container">
          <div className="content">
            <h1>Set Data</h1>
            <div className="input-container">
              <input
                type="text"
                placeholder="Key"
                value={key}
                onChange={(e) => setKey(e.target.value)}
              />
              <input
                type="text"
                placeholder="Value"
                value={value}
                onChange={(e) => setValue(e.target.value)}
              />
              <input
                type="number"
                placeholder="Expiration (seconds)"
                value={expiration}
                onChange={(e) => setExpiration(parseInt(e.target.value, 10))}
              />
            </div>
            <div className="button-container">
              <button onClick={handleSet}>Set</button>
              <button onClick={handleGet}>Get</button>
              <button onClick={handleDelete}>Delete</button>
            </div>
            {error && <div className="error">{error}</div>}
          </div>
        </div>
      )}
    </div>
  );
};

export default App;
