const express = require('express');
const fs = require('fs');
const cors = require('cors');
const User = require('./models/User.js');
const mongoose = require('mongoose');

const app = express();
const PORT = 3001;

mongoose.connect('mongodb+srv://u17104361:ThisIsInsideInsights1@swiftsignalsdb.m860tfn.mongodb.net/?retryWrites=true&w=majority&appName=SWiftSignalsDB', {
  useNewUrlParser: true,
  useUnifiedTopology: true,
})
.then(() => console.log('Connected to DB'))
.catch((err) => console.error('DB connection error:', err));

app.use(cors());
app.use(express.json());

const loadJSON = (filename) => {
  return JSON.parse(fs.readFileSync(`./data/${filename}`, 'utf8'));
};

app.get('/api/simulations', (_, res) => res.json(loadJSON('simulations.json')));
app.get('/api/intersections', (_, res) => res.json(loadJSON('intersections.json')));
app.get('/api/optimizations', (_, res) => res.json(loadJSON('optimizations.json')));
app.get('/api/traffic', (_, res) => res.json(loadJSON('traffic_data.json')));

//user routes
app.get('/users', async (req, res) => {
    const users = await User.find();
    res.json(users);
});

app.get('/users/:id', async (req, res) => {
    const users = await User.fnidOne({user_id: req.params.id});

    if(!user){
        return res.status(404).json({error: 'User not found.'});
    }

    res.json(user);
})

app.get('/users/:id/simulations', async (req, res) => {
    const users = await User.findOne({user_id: req.params.id});

    if(!user){
        return res.status(404).json({error: 'User not found.'});
    }

    res.json({simulations: user.simulations});
})

app.post('/users', async (req, res) => {
    const lastUser = await User.fnidOne().sort({user_id: -1}).exec();
    const lastUserNumber = parseInt(lastUser.user_id.slice(1));
    const newUserId = `u${(lastUserNumber + 1).toString().padStart(3, '0')}`;

    const defaultValues = {
        user_id: newUserId,
        name: "",
        lastName: "",
        userName: "",
        email: "",
        role: "analyst",
        profileArt: "https://www.google.com/imgres?q=placeholder%20image&imgurl=https%3A%2F%2Fwww.canbind.ca%2Fwp-content%2Fuploads%2F2025%2F01%2Fplaceholder-image-person-jpg.jpg&imgrefurl=https%3A%2F%2Fwww.canbind.ca%2Fabout-can-bind%2Four-team%2Fexecutive-committee%2Fplaceholder-image-person-jpg%2F&docid=-btyGbUwWLtnKM&tbnid=Cn7x48J0PlsGsM&vet=12ahUKEwjJ0JTSxaCNAxVvV0EAHY8xEkAQM3oECFYQAA..i&w=820&h=678&hcb=2&ved=2ahUKEwjJ0JTSxaCNAxVvV0EAHY8xEkAQM3oECFYQAA",
        created_at: new Date().toISOString(),
        simulations: [],
        auth: "",
    };

    const user = new User({...defaultValues, ...req.body});
    await user.save();
    res.status(201).json(user);
});

app.patch('/users/:id/profileArt', async (req, res) => {
    const {profileArt} = req.body;
    if (!profileArt){
        return res.status(400).json({error: 'ProfileArt required'});
    } 

    const user = await User.findOneAndUpdate(
        {user_id: req.params.id},
        {profileArt},
        {new: true}
    );

    if (!user){
        return res.status(404).json({error: 'User not found'});
    } 

    res.json({ profileArt: user.profileArt });
});

app.patch('/users/:id/password', async (req, res) => {
    const {newPassword} = req.body;
    if (!newPassword){
        return res.status(400).json({error: 'New password required'});
    } 

    const user = await User.findOneAndUpdate(
        {user_id: req.params.id},
        {auth: newPassword},
        {new: true}
    );

    if (!user){
        return res.status(404).json({error: 'User not found'});
    }

    res.json({message: 'Password updated successfully'});
});

app.delete('/users/:id', async (req, res) => {
    const user = await User.findOneAndDelete({user_id: req.params.id});
    if (!user){
        return res.status(404).json({error: 'User not found'});
    } 
    res.json({message: `User ${req.params.id} deleted.`});
});

app.post('/users/:id/simulations', (req, res) => {
    const users = loadUsers();
    const user = users.find(u => u.user_id === req.params.id);

    if(!user){
        return res.status(404).json({error: 'User not found.'});
    }

    const { simulation_id } = req.body;

    if(!simulation_id){
        return res.status(400).json({error: 'Simulation ID required.'});
    }

    if(!user.simulations.includes(simulation_id)){
        user.simulations.push(simulation_id);
    }

    saveUsers(users);

    res.json({simulations: user.simulations});
});

app.delete('/users/:id/simulations/:sim_id', (req, res) => {
    const users = loadUsers();
    const user = users.find(u => u.user_id === req.params.id);

    if(!user){
        return res.status(404).json({error: 'User not found.'});
    }

    user.simulations = user.simulations.filter(sim => sim !== req.params.sim_id);
    saveUsers(users);
    
    res.json({simulations: user.simulations});
});

app.listen(PORT, () => {
  console.log(`✅ Backend API running on http://localhost:${PORT}`);
});
