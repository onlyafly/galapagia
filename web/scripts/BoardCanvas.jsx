import React from 'react';
import api from './api';
import ReactDOM from 'react-dom';

export default React.createClass({

  componentDidMount: function () {
    ReactDOM.findDOMNode(this).appendChild(this.props.canvas);
    this.renderGrid();
  },

  componentDidUpdate: function(prevProps, prevState) {
    this.renderGrid();
  },

  renderGrid: function() {
    var grid = this.props.grid;
    var ctx = this.props.canvas.getContext("2d");

    const cellWidth = 5;
    const cellHeight = 5;
    const xCells = 100;
    const yCells = 100;
    const xOffset = 3;
    const yOffset = 3;

    ctx.fillStyle = "rgb(200,200,200)";
    ctx.fillRect(xOffset - 1, yOffset - 1, xCells*(cellWidth + 1) + 1, yCells*(cellHeight + 1) + 1);

    if (grid.length >= xCells && grid[0].length >= yCells) {
      for (var i = 0; i < xCells; i++) {
        for (var j = 0; j < yCells; j++) {
          var x = xOffset + ((cellWidth + 1) * i);
          var y = yOffset + ((cellHeight + 1) * j);
          var w = cellWidth;
          var h = cellHeight;

          switch (grid[i][j]) {
            case 0:
              ctx.fillStyle = "rgb(255,255,255)";
              break;
            case 1:
              ctx.fillStyle = "rgb(200,0,0)";
              break;
            default:
              ctx.fillStyle = "rgb(255,0,255)";
          }
          ctx.fillRect(x, y, w, h);
        }
      }
    }
  },

  render: function() {
    return <div />;
  }

});
