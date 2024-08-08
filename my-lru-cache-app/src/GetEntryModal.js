import React from 'react';

function GetEntryModal({
  show,
  handleClose,
  handleGetEntry,
  getKey,
  setGetKey,
}) {
  return (
    <div
      className={`modal fade ${show ? 'show' : ''}`}
      style={{ display: show ? 'block' : 'none' }}
      tabIndex="-1"
      role="dialog"
      aria-labelledby="getCacheModalLabel"
      aria-hidden={!show}
    >
      <div className="modal-dialog" role="document">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title" id="getCacheModalLabel">Get Cache Entry</h5>
            <button
              type="button"
              className="close"
              onClick={handleClose}
              aria-label="Close"
            >
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div className="modal-body">
            <div className="form-group">
              <label htmlFor="getKey">Key:</label>
              <input
                type="number"
                id="getKey"
                className="form-control"
                value={getKey}
                onChange={(e) => setGetKey(e.target.value)}
              />
            </div>
            <button
              className="btn btn-primary btn-lg"
              onClick={handleGetEntry}
            >
              Get Entry
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default GetEntryModal;
