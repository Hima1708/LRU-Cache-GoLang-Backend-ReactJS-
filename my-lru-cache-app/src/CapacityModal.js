import React from 'react';

function CapacityModal({ show, handleClose, handleChange, handleSubmit, newCapacity }) {
  return (
    <div className={`modal fade ${show ? 'show' : ''}`} style={{ display: show ? 'block' : 'none' }} tabIndex="-1" role="dialog" aria-labelledby="capacityModalLabel" aria-hidden={!show}>
      <div className="modal-dialog" role="document">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title" id="capacityModalLabel">Set Cache Capacity</h5>
            <button type="button" className="close" onClick={handleClose} aria-label="Close">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div className="modal-body">
            <div className="form-group">
              <label htmlFor="capacityInput">Capacity:</label>
              <input
                type="number"
                id="capacityInput"
                className="form-control"
                value={newCapacity}
                onChange={handleChange}
              />
            </div>
            <button
              className="btn btn-primary btn-lg"
              onClick={handleSubmit}
            >
              Set Capacity
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default CapacityModal;
