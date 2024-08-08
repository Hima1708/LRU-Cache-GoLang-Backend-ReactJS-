import React from 'react';

function CacheTable({ cacheData }) {
  const renderTable = () => {
    if (cacheData === null) {
      return <p>Loading...</p>;
    }

    if (Object.keys(cacheData).length === 0) {
      return <p>No data available</p>;
    }

    return (
      <div className="table-responsive">
        <table className="table table-bordered table-striped">
          <thead>
            <tr>
              <th>Key</th>
              <th>Value</th>
              <th>Expiration Time (seconds)</th>
            </tr>
          </thead>
          <tbody>
            {Object.entries(cacheData).map(([key, { value, expiration }]) => {
              const remainingTime = Math.max(0, expiration); // Adjust this based on backend data

              return (
                <tr key={key}>
                  <td>{key}</td>
                  <td>{value}</td>
                  <td>{remainingTime}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    );
  };

  return renderTable();
}

export default CacheTable;
