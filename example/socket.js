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

const nsp = io.of('stackcats');

nsp.on('connection', (socket) => {
  socket.emit('helloclient', 'test');
  socket.join('test1');
  socket.join('test2');
  socket.on('helloserver', (msg) => {
    console.log(msg);
  });
});

nsp.on('error', (err) => {
  console.log(err);
});
