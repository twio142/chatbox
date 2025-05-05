window.API_URL = process.env.NODE_ENV === 'development' ? 'localhost:8080' : window.location.host;
var socket = new WebSocket(`ws://${window.API_URL}/ws`);

let connect = cb => {
  console.log('Attempting Connection...');

  socket.onopen = () => {
    console.log('Successfully Connected');
  };

  socket.onmessage = msg => {
    console.log(msg);
    cb(msg);
  };

  socket.onclose = event => {
    console.log('Socket Closed Connection: ', event);
  };

  socket.onerror = error => {
    console.log('Socket Error: ', error);
  };
};

let sendMsg = ({text, fileName, fileURL}) => {
  let msg = JSON.stringify({
    text, fileName, fileURL
  });
  socket.send(msg);
};

export { connect, sendMsg };
