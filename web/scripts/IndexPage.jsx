import React from 'react';
import api from './api';

export default React.createClass({
  getInitialState: function() {
    return {counter: "", info: ""};
  },

  componentDidMount: function() {
    this.dataWebSocket = new WebSocket("ws://localhost:8080/api/dataWebSocket");
    this.dataWebSocket.onopen = () => {
      this.dataWebSocket.send(JSON.stringify({info: "Connection init"}));
    };
    this.dataWebSocket.onmessage = (e) => {
      this.setState({info: e.data});
    };
  },

  handleButtonClick: function() {
    api.get('/api/tick', (res) => {
      this.setState({counter: res.body.counter});
    }, (err, res) => {});
  },

  handleSendSomethingClick: function() {
    this.dataWebSocket.send(JSON.stringify({"test": 34}));
  },

  render: function() {
    return (
      <div>
        <div>{this.state.counter}</div>
        <div>Info:
          <span>{this.state.info}</span>
        </div>
        <button onClick={this.handleButtonClick}>Tick</button>
        <button onClick={this.handleSendSomethingClick}>SendSomething</button>
      </div>
    );
  }
});
