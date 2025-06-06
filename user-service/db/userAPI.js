const express = require('express');
const fs = require('fs');
const cors = require('cors');
const { json } = require('stream/consumers');

const app = express();
const PORT = 3001;

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
function loadUsers(){
    const raw = fs.readFileSync('./data/users.json');
    const data = JSON.parse(raw);
    return data.users;
}

function saveUsers(users){
    fs.writeFileSync('./data/users.json', JSON.stringify({users}, null, 2));
}

app.get('/users', (req, res) => {
    const users = loadUsers();
    res.json(users);
});

app.get('/users/:id', (req, res) => {
    const users = loadUsers();
    const user = users.find(u => u.user_id === req.params.id);

    if(!user){
        return res.status(404).json({error: 'User not found.'});
    }

    res.json(user);
})

app.get('/users/:id/simulations', (req, res) => {
    const users = loadUsers();
    const user = users.find(u => u.user_id === req.params.id);

    if(!user){
        return res.status(404).json({error: 'User not found.'});
    }

    res.json({simulations: user.simulations});
})

app.post('/users', (req, res) => {
    const users = loadUsers();
    const newUser = req.body;
    const lastUser = users[users.length - 1];
    const lastUserId = lastUser ? lastUser.user_id : "u000"
    const lastUserNumber = parseInt(lastUserId.slice(1));
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

    const userWithDefaults = {...defaultValues, ...newUser};

    if(!userWithDefaults.user_id || !userWithDefaults.userName){
        return res.status(400).json({error: 'User id and username required.'});
    } 

    users.push(userWithDefaults);
    fs.writeFileSync('./data/users.json', JSON.stringify({users}, null, 2));

    res.status(201).json(userWithDefaults);
});

app.patch('/users/:id/profileArt', (req, res) => {
    const users = loadUsers();
    const user = users.find(u => u.user_id === req.params.id);

    if(!user){
        return saveUsers.status(404).json({error: 'User not found.'});
    }

    const { profileArt } = req.body;

    if(!profileArt){
        return res.status(400).json({error: 'Profile image URL required'});
    }

    user.profileArt = profileArt;
    saveUsers(users);

    res.json({profileArt: user.profileArt});
});

app.patch('/users/:id/password', (req, res) => {
    const users = loadUsers();
    const user = users.find(u => u.user_id === req.params.id);

    if(!user){
        return res.status(404).json({error: 'User not found.'});
    }

    const { newPassword } = req.body;

    if(!newPassword){
        return res.status(400).json({error: 'New password required.'});
    }

    user.auth = newPassword;
    saveUsers(users);

    res.json({message: 'Password updated successfully.'});
});

app.delete('/users/:id', (req, res) => {
    let users = loadUsers();
    const user = users.find(u => u.user_id === req.params.id);

    if(!user){
        return res.status(404).json({error: 'User not found.'});
    }

    users = users.filter(u => u.user_id !== req.params.id);
    saveUsers(users);

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
