import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Profile">
      <h1>Edit your profile</h1>
    </div>;
  }
});
