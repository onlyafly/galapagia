import React from 'react';
import api from './api';

export default React.createClass({
  getInitialState: function() {
    return {counter: ""};
  },

  componentDidMount: function() {
  },

  handleButtonClick: function() {
    api.get('/api/tick', (res) => {
      this.setState({
        counter: res.body.counter
      });
    },
    (err, res) => {});
  },

  render: function() {
    return (
      <div>
        <div>{this.state.counter}</div>
        <button onClick={this.handleButtonClick}>Tick</button>
      </div>
    );
  }
});
