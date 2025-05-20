const mongoose = require('mongoose');

const userSchema = new mongoose.Schema({
  user_id: { type: String, required: true, unique: true },
  name: String,
  lastName: String,
  userName: { type: String, required: true },
  email: String,
  role: { type: String, default: 'analyst' },
  profileArt: String,
  created_at: { type: Date, default: Date.now },
  simulations: [String],
  auth: String // For real apps, use hashed passwords
});

module.exports = mongoose.model('User', userSchema);
