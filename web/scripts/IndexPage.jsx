import React from 'react';
import api from './api';
import BoardCanvas from './BoardCanvas';

export default React.createClass({
  getInitialState: function() {
    return {counter: "", info: "", grid: []};
  },

  componentDidMount: function() {
    this.dataWebSocket = new WebSocket("ws://localhost:8080/api/dataWebSocket");
    this.dataWebSocket.onopen = () => {
      this.dataWebSocket.send(JSON.stringify({command: "init connection"}));
    };
    this.dataWebSocket.onmessage = (e) => {
      this.setState({grid: JSON.parse(e.data)});
    };
  },

  handleButtonClick: function() {
    api.get('/api/tick', (res) => {
      this.setState({counter: res.body.counter});
    }, (err, res) => {});
  },

  handleClickShowCurrentGrid: function() {
    this.dataWebSocket.send(JSON.stringify({command: "show current grid"}));
  },

  handleClickReset: function() {
    this.dataWebSocket.send(JSON.stringify({command: "reset"}));
  },

  handleClickTick30Times: function() {
    this.dataWebSocket.send(JSON.stringify({command: "tick 30 times"}));
  },

  handleClickLogData: function() {
    this.dataWebSocket.send(JSON.stringify({command: "log data"}));
  },

  render: function() {
    var canvas = document.getElementById('boardCanvasElement');

    return (
      <div>
        <div>{this.state.counter}</div>
        <div>Info:
          <span>{this.state.info}</span>
        </div>
        <button onClick={this.handleButtonClick}>Counter: Tick</button>
        <button onClick={this.handleClickShowCurrentGrid}>Show Current Grid</button>
        <button onClick={this.handleClickReset}>Reset</button>
        <button onClick={this.handleClickTick30Times}>Tick 30 Times</button>
        <button onClick={this.handleClickLogData}>Log Data</button>
        <div>
          <BoardCanvas canvas={canvas} grid={this.state.grid} />
        </div>
      </div>
    );
  }
});
