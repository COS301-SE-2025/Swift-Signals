const express = require("express");
const cors = require("cors");
const dialogflow = require("@google-cloud/dialogflow").v2beta1;
const axios = require("axios");
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

    // Extract the token from the request body
    const { token } = req.body;
    console.log("--- SERVER RECEIVED ---", { token: token ? `A token was provided (length: ${token.length})` : "No token provided" });

    // Check for the Get_Intersections intent
    if (result.intent && result.intent.displayName === "Get_Intersections") {
      console.log("✅ Matched intent: Get_Intersections");
      if (!token) {
        console.log("❌ Condition failed: No token provided.");
        result.fulfillmentText = "I can't get your intersections without knowing who you are. Please make sure you are logged in.";
      } else {
        try {
          console.log("Attempting to call API: GET /intersections");
          const apiResponse = await axios.get("http://api-gateway:9090/intersections", {
            headers: { Authorization: `Bearer ${token}` },
          });
          console.log("✅ API call successful.", { data: apiResponse.data });

          const intersections = apiResponse.data.intersections;
          let responseText = "Here are your intersections:\n";
          intersections.forEach(intersection => {
            responseText += `\n- ${intersection.name} (Status: ${intersection.status})`;
          });
          result.fulfillmentText = responseText;
        } catch (apiError) {
          console.error("❌ API Gateway Error on /intersections:", apiError.message);
          result.fulfillmentText = "I was unable to fetch your intersections at the moment. Please try again later.";
        }
      }
    }

    // Check for the Default Welcome Intent to add the user's name
    if (result.intent && result.intent.displayName === "Default Welcome Intent") {
      console.log("✅ Matched intent: Default Welcome Intent");
      if (!token) {
        console.log("❌ Condition failed: No token provided. Using default greeting.");
      } else {
        try {
          console.log("Attempting to call API: GET /me");
          const apiResponse = await axios.get("http://api-gateway:9090/me", {
            headers: { Authorization: `Bearer ${token}` },
          });
          console.log("✅ API call successful.", { data: apiResponse.data });
          const userName = apiResponse.data.username || 'there';
          result.fulfillmentText = `Hello, ${userName}! I'm here to help. What can I assist you with today?`;
        } catch (apiError) {
          console.error("❌ API Gateway Error on /me:", apiError.message);
          // If the API call fails, we just fall back to the default non-personalized greeting.
        }
      }
    }

    // --- NEW LOGIC FOR Create.Intersection ---
    if (result.intent && result.intent.displayName === "Create.Intersection") {
      console.log("✅ Matched intent: Create.Intersection");
      console.log("Parameters collected so far:", result.parameters.fields);
      console.log("All required params present:", result.allRequiredParamsPresent);

      if (!token) {
        result.fulfillmentText = "I can't create an intersection without knowing who you are. Please make sure you are logged in.";
      } else if (!result.allRequiredParamsPresent) {
        // If not all parameters are present, let Dialogflow handle the prompts
        // The fulfillmentText will contain Dialogflow's prompt for the next parameter
        console.log("Not all parameters present. Returning Dialogflow's prompt.");
        // No need to modify result.fulfillmentText here, Dialogflow already set it.
      } else {
        // All parameters are present, proceed with API call
        try {
          const params = result.parameters.fields;
          const requestBody = {
            name: params['intersection-name'].stringValue,
            details: {
              address: params.address.stringValue,
              city: params.city ? params.city.stringValue : "",
              province: params.province ? params.province.stringValue : ""
            },
            traffic_density: params['traffic-density'].stringValue.toLowerCase(),
            default_parameters: {
              green: params['green-light'].numberValue,
              yellow: params['yellow-light'].numberValue,
              red: params['red-light'].numberValue,
              speed: params['vehicle-speed'] ? params['vehicle-speed'].numberValue : 0,
              intersection_type: "default", // Placeholder, needs to be added to Dialogflow if dynamic
              seed: Math.floor(Math.random() * 1000000000) // Add a random seed, as it's required
            }
          };

          console.log("Attempting to call API: POST /intersections with body:", requestBody);
          const apiResponse = await axios.post("http://api-gateway:9090/intersections", requestBody, {
            headers: { Authorization: `Bearer ${token}` },
          });
          console.log("✅ API call successful.", { data: apiResponse.data });

          result.fulfillmentText = `Intersection '${requestBody.name}' created successfully!`;

        } catch (apiError) {
          console.error("❌ API Gateway Error on /intersections (POST):", apiError.message);
          result.fulfillmentText = "Sorry, I couldn't create the intersection. Please check the details and try again.";
        }
      }
    }

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
