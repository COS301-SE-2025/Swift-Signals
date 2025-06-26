// mock-api/server.js
const express = require('express');
const cors = require('cors');
const app = express();
const port = 9090;

console.log("Mock API Server starting...");

// =============================================================================
// CORS and JSON Parser Setup
// =============================================================================
app.options(/.*/, cors()); 
app.use(cors());
app.use(express.json());
// =============================================================================


// =============================================================================
// Mock Database
// =============================================================================
let mockIntersections = [
  {
    id: "1",
    name: "Main St & 1st Ave",
    details: {
      address: "123 Main St, Pretoria CBD",
      city: "Pretoria",
      province: "Gauteng",
    },
    default_parameters: {
      green: 15,
      yellow: 3,
      red: 5,
      speed: 60,
      seed: 12345,
      intersection_type: "4-way",
    },
    status: "optimised",
    created_at: "2025-06-24T15:04:05Z",
    last_run_at: "2025-06-26T11:30:00Z",
    run_count: 12,
  },
  {
    id: "2",
    name: "Church St & Park Rd",
    details: {
      address: "45 Park Rd, Hatfield",
      city: "Pretoria",
      province: "Gauteng",
    },
    default_parameters: {
      green: 10,
      yellow: 2,
      red: 4,
      speed: 50,
      seed: 67890,
      intersection_type: "t-junction",
    },
    status: "unoptimised",
    created_at: "2025-06-20T10:00:00Z",
    last_run_at: null,
    run_count: 0,
  },
  {
    id: "3",
    name: "University & Lynnwood",
    details: {
      address: "Corner of University and Lynnwood, Hatfield",
      city: "Pretoria",
      province: "Gauteng",
    },
    default_parameters: {
      green: 20,
      yellow: 4,
      red: 6,
      speed: 60,
      seed: 54321,
      intersection_type: "4-way",
    },
    status: "unoptimised",
    created_at: "2025-05-15T09:00:00Z",
    last_run_at: null,
    run_count: 0,
  }
];

// =============================================================================
// Middleware for Authorization
// =============================================================================
const checkAuth = (req, res, next) => {
  const authHeader = req.headers['authorization'];
  // In a real app, you'd validate the JWT. Here we just check for its presence.
  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    console.log('[AUTH] Failure: Missing or invalid token.');
    return res.status(401).json({
      code: 'UNAUTHORIZED',
      message: 'Unauthorized: Token missing or invalid'
    });
  }
  console.log('[AUTH] Success: Token present.');
  next();
};


// =============================================================================
// Authentication Endpoints (Existing)
// =============================================================================
app.post('/login', (req, res) => {
  const { email, password } = req.body;
  console.log(`[LOGIN] Attempt for email: ${email}`);

  if (email === 'user@example.com' && password === 'password123') {
    console.log('[LOGIN] Success.');
    return res.status(200).json({
      message: 'Login successful',
      // Using a static token for easy testing
      token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ik1vY2sgVXNlciIsImlhdCI6MTUxNjIzOTAyMn0.dummy-token-for-testing',
    });
  }
  
  console.log('[LOGIN] Failure: Invalid credentials.');
  return res.status(400).json({
    code: 'BAD_REQUEST',
    message: 'Invalid email or password.',
  });
});

app.post('/register', (req, res) => {
  const { username, email, password } = req.body;
  console.log(`[REGISTER] Attempt for username: ${username}, email: ${email}`);

  if (email === 'taken@example.com') {
    console.log('[REGISTER] Failure: Email already exists.');
    return res.status(400).json({
      code: 'BAD_REQUEST',
      message: 'This email address is already registered.',
    });
  }

  if (username && email && password) {
    console.log('[REGISTER] Success.');
    return res.status(201).json({ 
      user_id: 'mock-user-' + Date.now(),
    });
  }
  
  console.log('[REGISTER] Failure: Missing fields.');
  return res.status(400).json({
    code: 'BAD_REQUEST',
    message: 'Missing required fields.',
  });
});


// =============================================================================
// Intersection Endpoints (NEW)
// =============================================================================

// --- GET /intersections (Get All) ---
app.get('/intersections', checkAuth, (req, res) => {
  console.log(`[GET /intersections] Fetching all ${mockIntersections.length} intersections.`);
  // The API spec wraps the array in an object
  return res.status(200).json({
    intersections: mockIntersections
  });
});

// --- GET /intersections/:id (Get by ID) ---
app.get('/intersections/:id', checkAuth, (req, res) => {
  const { id } = req.params;
  console.log(`[GET /intersections/:id] Searching for ID: ${id}`);
  const intersection = mockIntersections.find(i => i.id === id);

  if (intersection) {
    console.log(`[GET /intersections/:id] Found: ${intersection.name}`);
    return res.status(200).json(intersection);
  }

  console.log(`[GET /intersections/:id] ID ${id} not found.`);
  return res.status(404).json({
    code: 'NOT_FOUND',
    message: 'Not Found: Intersection does not exist'
  });
});

// --- POST /intersections (Create) ---
app.post('/intersections', checkAuth, (req, res) => {
  const { name, details, default_parameters, traffic_density } = req.body;
  console.log(`[POST /intersections] Attempting to create: ${name}`);

  if (!name || !details || !default_parameters || !details.address) {
    console.log('[POST /intersections] Failure: Missing required fields.');
    return res.status(400).json({
      code: 'BAD_REQUEST',
      message: 'Invalid request payload or missing fields'
    });
  }

  const newId = String(Date.now()); // Create a unique ID
  const newIntersection = {
    id: newId,
    name,
    details,
    default_parameters,
    traffic_density: traffic_density || 'low',
    status: 'unoptimised',
    created_at: new Date().toISOString(),
    last_run_at: null,
    run_count: 0,
  };

  mockIntersections.push(newIntersection);

  console.log(`[POST /intersections] Success. New intersection '${name}' created with ID: ${newId}`);
  // The API spec returns the new ID in an object
  return res.status(201).json({
    id: newId
  });
});


// =============================================================================
// Start Server
// =============================================================================
app.listen(port, () => {
  console.log(`🚀 Mock API server is running at http://localhost:${port}`);
  console.log('Endpoints ready: /login, /register, /intersections, /intersections/:id');
});