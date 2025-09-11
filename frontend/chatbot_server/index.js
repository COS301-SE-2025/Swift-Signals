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
          if (intersections && intersections.length > 0) {
            let responseText = "Here are your intersections:\n";
            intersections.forEach(intersection => {
              responseText += `\n- ${intersection.name}`;
            });
            responseText += "\n\nTo get more details about a specific intersection, please say something like 'Tell me about [intersection name]' or 'Show details for [intersection name]'.";
            result.fulfillmentText = responseText;
          } else {
            result.fulfillmentText = "You don't have any intersections yet.";
          }
        } catch (apiError) {
          console.error("❌ API Gateway Error on /intersections:", apiError.message);
          result.fulfillmentText = "I was unable to fetch your intersections at the moment. Please try again later.";
        }
      }
    }

    // NEW LOGIC FOR Get_Intersection_Details
    if (result.intent && result.intent.displayName === "Get_Intersection_Details") {
      console.log("✅ Matched intent: Get_Intersection_Details");
      const intersectionIdentifier = result.parameters.fields.intersection_identifier.stringValue;

      if (!token) {
        result.fulfillmentText = "I can't get intersection details without knowing who you are. Please make sure you are logged in.";
      } else if (!intersectionIdentifier) {
        result.fulfillmentText = "Which intersection are you interested in? Please provide its name.";
      } else {
        try {
          console.log(`Attempting to get details for: ${intersectionIdentifier}`);
          // First, get all intersections to find the matching one by name or ID
          const allIntersectionsResponse = await axios.get("http://api-gateway:9090/intersections", {
            headers: { Authorization: `Bearer ${token}` },
          });
          const allIntersections = allIntersectionsResponse.data.intersections;

          let targetIntersection = null;
          // Try to find by ID first (assuming ID is a string)
          targetIntersection = allIntersections.find(
            (int) => int.id === intersectionIdentifier
          );

          // If not found by ID, try to find by name (case-insensitive)
          if (!targetIntersection) {
            targetIntersection = allIntersections.find(
              (int) => int.name.toLowerCase() === intersectionIdentifier.toLowerCase()
            );
          }

          if (targetIntersection) {
            let fulfillmentText = `--- Details for ${targetIntersection.name} ---
`;
            fulfillmentText += `
ID: ${targetIntersection.id}`;
            fulfillmentText += `
Status: ${targetIntersection.status}`;

            if (targetIntersection.details) {
              fulfillmentText += `

--- Location Details ---`;
              if (targetIntersection.details.address) {
                fulfillmentText += `
Address: ${targetIntersection.details.address}`;
              }
              if (targetIntersection.details.city) {
                fulfillmentText += `
City: ${targetIntersection.details.city}`;
              }
              if (targetIntersection.details.province) {
                fulfillmentText += `
Province: ${targetIntersection.details.province}`;
              }
            }

            fulfillmentText += `

--- Simulation Info ---`;
            fulfillmentText += `
Traffic Density: ${targetIntersection.traffic_density}`;
            fulfillmentText += `
Run Count: ${targetIntersection.run_count}`;
            fulfillmentText += `
Last Run: ${new Date(targetIntersection.last_run_at).toLocaleString()}`;
            fulfillmentText += `
Created: ${new Date(targetIntersection.created_at).toLocaleString()}`;
            result.fulfillmentText = fulfillmentText;
          } else {
            result.fulfillmentText = `I couldn't find an intersection named or with ID '${intersectionIdentifier}'. Please check the name and try again.`;
          }
        } catch (apiError) {
          console.error("❌ API Gateway Error on Get_Intersection_Details:", apiError.message);
          result.fulfillmentText = "I encountered an error while trying to retrieve the intersection details. Please try again later.";
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
              intersection_type: "traffic_light", // Defaulting to a standard traffic light intersection
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

app.get("/api/reverse-geocode", async (req, res) => {
  const { lat, lon } = req.query;

  if (!lat || !lon) {
    return res.status(400).send({ error: "Latitude and longitude are required" });
  }

  const url = `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lon}`;

  try {
    const response = await axios.get(url, {
        headers: {
            'User-Agent': 'Swift-Signals/1.0'
        }
    });
    res.status(200).send(response.data);
  } catch (error) {
    console.error("Nominatim Error:", error.message);
    res.status(500).send({ error: "Failed to fetch data from Nominatim" });
  }
});

app.get("/api/search-streets", async (req, res) => {
  const { q, type } = req.query;

  if (!q) {
    return res.status(400).send({ error: "Query is required" });
  }

  const params = new URLSearchParams({
    format: "json",
    addressdetails: "1",
    limit: "20",
    countrycodes: "za",
    q: q,
    extratags: "1",
  });

  if (type === "road") {
    params.append("class", "highway");
  }

  const url = `https://nominatim.openstreetmap.org/search?${params.toString()}`;

  try {
    const response = await axios.get(url, {
        headers: {
            'User-Agent': 'Swift-Signals/1.0'
        }
    });
    res.status(200).send(response.data);
  } catch (error) {
    console.error("Nominatim Error:", error.message);
    res.status(500).send({ error: "Failed to fetch data from Nominatim" });
  }
});

// --- START THE SERVER ---
app.listen(port, () => {
  console.log(`Chatbot backend server listening at http://localhost:${port}`);
});
