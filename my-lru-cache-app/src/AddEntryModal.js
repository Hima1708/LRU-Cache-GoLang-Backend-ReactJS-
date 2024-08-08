import React from 'react';

function AddEntryModal({
  show,
  handleClose,
  handleAddEntry,
  newKey,
  newValue,
  expiration,
  setNewKey,
  setNewValue,
  setExpiration,
}) {
  return (
    <div
      className={`modal fade ${show ? 'show' : ''}`}
      style={{ display: show ? 'block' : 'none' }}
      tabIndex="-1"
      role="dialog"
      aria-labelledby="addCacheModalLabel"
      aria-hidden={!show}
    >
      <div className="modal-dialog" role="document">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title" id="addCacheModalLabel">Add Cache Entry</h5>
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
              <label htmlFor="newKey">Key:</label>
              <input
                type="number"
                id="newKey"
                className="form-control"
                value={newKey}
                onChange={(e) => setNewKey(e.target.value)}
              />
            </div>
            <div className="form-group">
              <label htmlFor="newValue">Value:</label>
              <input
                type="number"
                id="newValue"
                className="form-control"
                value={newValue}
                onChange={(e) => setNewValue(e.target.value)}
              />
            </div>
            <div className="form-group">
              <label htmlFor="expiration">Expiration Time (seconds):</label>
              <input
                type="number"
                id="expiration"
                className="form-control"
                value={expiration}
                onChange={(e) => setExpiration(e.target.value)}
              />
            </div>
            <button
              className="btn btn-primary btn-lg"
              onClick={handleAddEntry}
            >
              Add Entry
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AddEntryModal;
