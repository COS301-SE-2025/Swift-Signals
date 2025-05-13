const express = require('express');
const fs = require('fs');
const cors = require('cors');

const app = express();
const PORT = 3001;

app.use(cors());
app.use(express.json());

const loadJSON = (filename) => {
  return JSON.parse(fs.readFileSync(`./data/${filename}`, 'utf8'));
};

app.get('/api/users', (_, res) => res.json(loadJSON('users.json')));
app.get('/api/simulations', (_, res) => res.json(loadJSON('simulations.json')));
app.get('/api/intersections', (_, res) => res.json(loadJSON('intersections.json')));
app.get('/api/optimizations', (_, res) => res.json(loadJSON('optimizations.json')));
app.get('/api/traffic', (_, res) => res.json(loadJSON('traffic_data.json')));

app.listen(PORT, () => {
  console.log(`✅ Backend API running on http://localhost:${PORT}`);
});
