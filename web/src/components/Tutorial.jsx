import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Tutorial">
      <h1>Tutorial</h1>
    </div>;
  }
});