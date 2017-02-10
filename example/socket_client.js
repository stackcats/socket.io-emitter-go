const socket = require('socket.io-client')('http://localhost:3000');

socket.on('connect', () => {
  console.log('connect');
});

socket.on('ping', (msg) => {
  console.log('ping: ', msg);
  socket.emit('pong', 'Hello World');
});

socket.on('disconnect', () => {
  console.log('disconnect');
});

socket.on('error', (err) => {
  console.log(err);
});
