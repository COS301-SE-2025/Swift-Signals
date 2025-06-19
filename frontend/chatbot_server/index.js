// server/index.js

const express = require('express');
const cors = require('cors');
const dialogflow = require('@google-cloud/dialogflow');
require('dotenv').config(); // Load environment variables from .env file

const app = express();
const port = process.env.PORT || 3001; // Use port 3001 or one from environment

// --- CONFIGURATION ---

// Use CORS to allow communication from your React app's origin
// (e.g., http://localhost:3000)
app.use(cors()); 

// Middleware to parse JSON bodies from incoming requests
app.use(express.json());

// Configure the Dialogflow credentials
const credentials = {
  client_email: process.env.DIALOGFLOW_CLIENT_EMAIL,
  private_key: process.env.DIALOGFLOW_PRIVATE_KEY.replace(/\\n/g, '\n'),
};

const sessionClient = new dialogflow.SessionsClient({
  projectId: process.env.DIALOGFLOW_PROJECT_ID,
  credentials,
});

// --- API ENDPOINT ---

// Define the single endpoint for your chatbot
app.post('/api/chatbot', async (req, res) => {
  // It can receive a 'message' for text, or an 'event' for things like the welcome signal
  const { message, event, sessionId } = req.body;
  const projectId = process.env.DIALOGFLOW_PROJECT_ID;

  if (!sessionId || (!message && !event)) {
    return res.status(400).send({ error: 'SessionId and either a message or an event are required' });
  }

  const sessionPath = sessionClient.projectAgentSessionPath(projectId, sessionId);

  const request = {
    session: sessionPath,
    queryInput: {},
  };

  // This logic is CRUCIAL. It builds the request differently for events vs. text.
  if (event) {
    request.queryInput.event = {
      name: event,
      languageCode: 'en-US',
    };
  } else {
    request.queryInput.text = {
      text: message,
      languageCode: 'en-US',
    };
  }

  try {
    const responses = await sessionClient.detectIntent(request);
    const result = responses[0].queryResult;
    res.status(200).send(result);
  } catch (error) {
    console.error('Dialogflow Error:', error);
    res.status(500).send({ error: 'Internal Server Error' });
  }
});



// --- START THE SERVER ---

app.listen(port, () => {
  console.log(`Chatbot backend server listening at http://localhost:${port}`);
});