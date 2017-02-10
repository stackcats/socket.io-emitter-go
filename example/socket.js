const adapter = require('socket.io-redis');
const redis = require('redis');
const io = require('socket.io')(3000);

const pub = redis.createClient();
const sub = redis.createClient();

io.adapter(adapter({
  pubClient: pub,
  subClient: sub,
  key: 'socket.io'
}));

io.on('connection', (socket) => {
  socket.emit('ping', 'test');
  socket.on('pong', (msg) => {
    console.log(msg);
  });
});

io.on('error', (err) => {
  console.log(err);
});
