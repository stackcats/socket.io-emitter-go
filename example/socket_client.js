const socket = require('socket.io-client')('http://localhost:3000/stackcats');

socket.on('connect', () => {
  console.log('connect');
});

socket.on('helloclient', (msg) => {
  console.log('helloclient: ', msg);
  socket.emit('helloserver', 'Hello World');
});

socket.on('disconnect', () => {
  console.log('disconnect');
});

socket.on('error', (err) => {
  console.log(err);
});
