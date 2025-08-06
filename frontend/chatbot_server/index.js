const express = require("express");
const cors = require("cors");
const dialogflow = require("@google-cloud/dialogflow").v2beta1;
require("dotenv").config();

const app = express();
const port = process.env.PORT || 3001;

// --- CONFIGURATION ---
app.use(cors());
app.use(express.json());

const credentials = {
  client_email: process.env.DIALOGFLOW_CLIENT_EMAIL,
  private_key: process.env.DIALOGFLOW_PRIVATE_KEY.replace(/\\n/g, "\n"),
};

const sessionClient = new dialogflow.SessionsClient({
  projectId: process.env.DIALOGFLOW_PROJECT_ID,
  credentials,
});

// --- API ENDPOINT ---
app.post("/api/chatbot", async (req, res) => {
  const { message, event, sessionId } = req.body;
  const projectId = process.env.DIALOGFLOW_PROJECT_ID;

  if (!sessionId || (!message && !event)) {
    return res.status(400).send({
      error: "SessionId and either a message or an event are required",
    });
  }

  const sessionPath = sessionClient.projectAgentSessionPath(
    projectId,
    sessionId,
  );

  const request = {
    session: sessionPath,
    queryInput: {},
    queryParams: {
      knowledgeBaseNames: [
        // Please meticulously double-check this ID one last time
        `projects/swift-signals/knowledgeBases/MTUzOTA5NDkxOTcwNzkzNzk5Njk`,
      ],
    },
  };

  if (event) {
    request.queryInput.event = {
      name: event,
      languageCode: "en-US",
    };
  } else {
    request.queryInput.text = {
      text: message,
      languageCode: "en-US",
    };
  }

  // --- THIS IS THE CRUCIAL DIAGNOSTIC STEP ---
  console.log("--- FINAL REQUEST BEING SENT TO DIALOGFLOW ---");
  console.log(JSON.stringify(request, null, 2));
  // -------------------------------------------------

  try {
    const responses = await sessionClient.detectIntent(request);
    const result = responses[0].queryResult;

    // --- ENHANCED DIAGNOSTIC LOG ---
    console.log("--- FULL DIALOGFLOW RESPONSE ---");
    console.log(JSON.stringify(result, null, 2));

    // Specifically check for the presence of knowledge answers
    if (
      result.knowledgeAnswers &&
      result.knowledgeAnswers.answers &&
      result.knowledgeAnswers.answers.length > 0
    ) {
      console.log("✅ SUCCESS: Knowledge base answer was found by Dialogflow.");
    } else {
      console.log(
        "❌ FAILURE: No knowledge base answer was returned by Dialogflow.",
      );
      if (result.intent) {
        console.log(
          `Instead, the '${result.intent.displayName}' intent was matched with confidence: ${result.intentDetectionConfidence}`,
        );
      }
    }
    // ------------------------------------

    res.status(200).send(result);
  } catch (error) {
    console.error("Dialogflow Error:", error);
    res.status(500).send({ error: "Internal Server Error" });
  }
});

// --- START THE SERVER ---
app.listen(port, () => {
  console.log(`Chatbot backend server listening at http://localhost:${port}`);
});