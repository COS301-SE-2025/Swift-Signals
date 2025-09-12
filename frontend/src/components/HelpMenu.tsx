import React, { useEffect, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { v4 as uuidv4 } from "uuid";
import "../styles/HelpMenu.css";
import InteractiveTutorial, { type TutorialStep } from "./InteractiveTutorial";
import {
  FaTimes,
  FaCommentDots,
  FaBook,
  FaChevronLeft,
  FaChevronDown,
} from "react-icons/fa";
import { IoSend } from "react-icons/io5";

type QuickReply = { text: string; payload: string };
type ChatMessage = {
  text: string;
  sender: "user" | "bot";
  quickReplies?: QuickReply[];
};
type TutorialType =
  | "dashboard"
  | "navigation"
  | "intersections"
  | "simulations"
  | "users"
  | "simulation-results"
  | "comparison-view";
type DialogflowMessage = {
  payload?: {
    fields?: {
      richContent: {
        listValue: {
          values: {
            listValue: {
              values: {
                structValue: {
                  fields: {
                    options: {
                      listValue: {
                        values: DialogflowQuickReplyOption[];
                      };
                    };
                  };
                };
              }[];
            };
          }[];
        };
      };
    };
  };
};

type DialogflowQuickReplyOption = {
  structValue: {
    fields: {
      text: { stringValue: string };
      link?: { stringValue: string };
    };
  };
};

const dashboardTutorialSteps: TutorialStep[] = [
  {
    selector: ".card-grid",
    title: "Summary Cards",
    text: "These cards give you a quick, at-a-glance overview of your key metrics.",
    position: "bottom",
  },
  {
    selector: ".recent-simulations-tab",
    title: "Simulations Table",
    text: "Here you can see a list of all your recent simulations. Click on any row to see more details.",
    position: "right",
  },
  {
    selector: ".newInt",
    title: "Add a New Intersection",
    text: "Click this button to open the form for creating a new traffic intersection.",
    position: "bottom",
  },
  {
    selector: ".runSim",
    title: "Run a Simulation",
    text: "Click this button to open the form for running a traffic simulation.",
    position: "bottom",
  },
  {
    selector: ".viewMap",
    title: "View Map",
    text: "This will take you to a full-screen map view of all your monitored intersections.",
    position: "bottom",
  },
  {
    selector: ".graph-card",
    title: "Traffic Volume Chart",
    text: "This chart shows the traffic volume over time for your key intersections.",
    position: "left",
  },
  {
    selector: ".inter-card",
    title: "Top Intersections",
    text: "This card displays the top intersections based on traffic volume.",
    position: "left",
  },
];

const navigationTutorialSteps: TutorialStep[] = [
  {
    selector: ".nav-links",
    title: "Main Navigation",
    text: "Use these links to switch between the main pages of the application.",
    position: "bottom",
  },
  {
    selector: ".user-profile",
    title: "User Profile",
    text: "Access your profile, settings, or log out from this menu.",
    position: "bottom",
  },
  {
    selector: ".footer-toggle",
    title: "Appearance Toggle",
    text: "Switch between light and dark modes.",
    position: "right",
  },
];

const intersectionTutorialSteps: TutorialStep[] = [
  {
    selector: ".searchContainer",
    title: "Search Bar",
    text: "This allows you to quickly find intersections by name or ID.",
    position: "bottom",
  },
  {
    selector: ".addIntersectionBtn",
    title: "Add Intersection",
    text: "Click this button to open the form for adding a new traffic intersection. Let's see how to create a new intersection - the tutorial will now open the form for you.",
    position: "bottom",
  },
  {
    selector: "body",
    title: "Opening Form",
    text: "Please wait...",
    position: "center",
    autoAdvance: true,
    waitFor: ".fixed.inset-0.bg-black.bg-opacity-50",
    action: () => {
      const button = document.querySelector(
        ".addIntersectionBtn",
      ) as HTMLElement;
      if (button) button.click();
    },
  },
  {
    selector: ".fixed.inset-0.bg-black.bg-opacity-50 > div",
    title: "Add Intersection Modal",
    text: "This modal allows you to create a new traffic intersection. Fill in all the required details to add it to your system.",
    position: "left",
  },
  {
    selector: "#name",
    title: "Intersection Name",
    text: "Enter a descriptive name for your intersection to help identify it later.",
    position: "right",
  },
  {
    selector: "#details\\.address",
    title: "Address Field",
    text: "Provide the street address where this intersection is located.",
    position: "right",
  },
  {
    selector: "#traffic_density",
    title: "Traffic Density",
    text: "Select the expected traffic volume: Low, Medium, or High. This affects simulation parameters.",
    position: "left",
  },
  {
    selector: "#default_parameters\\.green",
    title: "Green Light Duration",
    text: "Set the default duration (in seconds) for the green light phase.",
    position: "bottom",
  },
  {
    selector: "#default_parameters\\.yellow",
    title: "Yellow Light Duration",
    text: "Set the default duration (in seconds) for the yellow/amber light phase.",
    position: "bottom",
  },
  {
    selector: "#default_parameters\\.red",
    title: "Red Light Duration",
    text: "Set the default duration (in seconds) for the red light phase.",
    position: "bottom",
  },
  {
    selector: "button[type='submit']",
    title: "Create Intersection",
    text: "Once you've filled in all the details, click this button to create your new intersection. Now let's close this modal to continue exploring the intersections page.",
    position: "left",
  },
  {
    selector: "body",
    title: "Closing Modal",
    text: "Please wait...",
    position: "center",
    autoAdvance: true,
    waitFor: ".intersectionCard",
    action: () => {
      const closeButton = document.querySelector(
        ".fixed.inset-0.bg-black.bg-opacity-50 button[onClick]",
      ) as HTMLElement;
      if (closeButton) {
        closeButton.click();
      } else {
        // Try alternative selector for the X button
        const xButton = document.querySelector(
          ".absolute.top-4.right-4",
        ) as HTMLElement;
        if (xButton) xButton.click();
      }
    },
  },
  {
    selector: ".intersectionCard",
    title: "Intersection Cards",
    text: "Each card represents a traffic intersection, displaying key information like name, location, and type.",
    position: "left",
  },
  {
    selector: ".simButton",
    title: "Simulate Button",
    text: "Click this button to run a traffic simulation for the selected intersection.",
    position: "right",
  },
  {
    selector: ".editButton",
    title: "Edit Button",
    text: "Click this button to edit the details of the selected intersection.",
    position: "right",
  },
  {
    selector: ".deleteIntButton",
    title: "Delete Button",
    text: "Click this button to delete the selected intersection.",
    position: "right",
  },
];

const simulationsTutorialSteps: TutorialStep[] = [
  {
    selector: ".sims",
    title: "Simulations",
    text: "This page shows your recent simulations.",
    position: "right",
  },
  {
    selector: ".opts",
    title: "Optimizations",
    text: "This page shows your recent optimizations.",
    position: "left",
  },
  {
    selector: ".viewBtn",
    title: "View a Simulation",
    text: "This button let's you view a simulation.",
    position: "left",
  },
  {
    selector: ".deleteBtn",
    title: "Delete a Simulation",
    text: "This button let's you delete a simulation.",
    position: "left",
  },
  {
    selector: ".pagination",
    title: "Cycle Through Pages",
    text: "Here you can navigate to view multiple pages of simulations.",
    position: "right",
  },
  {
    selector: ".new-simulation-button",
    title: "Create a New Simulation",
    text: "Let's see how to create a new simulation. The tutorial will now open the form for you.",
    position: "bottom",
  },
  {
    selector: "body",
    title: "Opening Form",
    text: "Please wait...",
    position: "center",
    autoAdvance: true,
    waitFor: ".simulation-modal-content",
    action: () => {
      const button = document.querySelector(
        ".new-simulation-button",
      ) as HTMLElement;
      if (button) button.click();
    },
  },
  {
    selector: ".simulation-modal-content",
    title: "New Simulation Form",
    text: "This modal allows you to create a new simulation. Fill in all the required details to set up your traffic simulation.",
    position: "left",
  },
  {
    selector: ".simulation-name-input",
    title: "Simulation Name",
    text: "Give your simulation a unique name so you can easily identify it later.",
    position: "right",
  },
  {
    selector: "textarea",
    title: "Simulation Description",
    text: "Add an optional description to provide more details about this simulation.",
    position: "right",
  },
  {
    selector: ".intersection-tabs",
    title: "Intersection Selection Methods",
    text: "You can add intersections to your simulation using three different methods: List, Search, or Map. Let's explore each one.",
    position: "left",
  },
  {
    selector: ".intersection-tabs button:nth-child(1)",
    title: "List Tab",
    text: "The List tab shows pre-defined intersections. This is the default active tab showing available intersections in a dropdown.",
    position: "bottom",
  },
  {
    selector: ".intersection-tabs button:nth-child(2)",
    title: "Search Tab",
    text: "The Search tab allows you to find intersections by searching for street names. Click this tab to explore street search.",
    position: "bottom",
    action: () => {
      const searchButton = document.querySelector(
        ".intersection-tabs button:nth-child(2)",
      ) as HTMLElement;
      if (searchButton) searchButton.click();
    },
  },
  {
    selector: "input[placeholder*='Type a street name']",
    title: "First Street Search",
    text: "Type the name of the first street to search for real South African streets. The system will find matching streets automatically.",
    position: "right",
  },
  {
    selector: ".intersection-tabs button:nth-child(3)",
    title: "Map Tab",
    text: "The Map tab lets you visually select intersections by clicking on a map. Click this tab to explore map selection.",
    position: "bottom",
    action: () => {
      const mapButton = document.querySelector(
        ".intersection-tabs button:nth-child(3)",
      ) as HTMLElement;
      if (mapButton) mapButton.click();
    },
  },
  {
    selector: ".leaflet-container",
    title: "Interactive Map",
    text: "Click anywhere on this map to automatically find the nearest road intersection. The system will snap your click to actual intersections.",
    position: "right",
  },
  {
    selector: ".flex.flex-wrap.gap-2",
    title: "Selected Intersections Area",
    text: "Selected intersections will appear as pills in this area. You can remove them by clicking the × button on each pill when you select intersections.",
    position: "left",
  },
  {
    selector: ".create-simulation-submit-btn",
    title: "Create Simulation",
    text: "Once you've named your simulation and selected intersections, click here to create and run your simulation. Now let's close this modal to continue exploring.",
    position: "right",
  },
  {
    selector: "body",
    title: "Closing Modal",
    text: "Please wait...",
    position: "center",
    autoAdvance: true,
    waitFor: ".sims",
    action: () => {
      const closeButton = document.querySelector(".crossBtn") as HTMLElement;
      if (closeButton) {
        closeButton.click();
      } else {
        // Try alternative close method
        const modalOverlay = document.querySelector(
          ".fixed.inset-0.z-50",
        ) as HTMLElement;
        if (modalOverlay) {
          // Click outside the modal to close
          const event = new MouseEvent("click", { bubbles: true });
          modalOverlay.dispatchEvent(event);
        }
      }
    },
  },
];

const usersTutorialSteps: TutorialStep[] = [
  {
    selector: ".usersTable",
    title: "Users Table",
    text: "This displays all the users currently signed in.",
    position: "left",
  },
  {
    selector: ".editUser",
    title: "Edit User",
    text: "This allows you to edit the user's details. This can only be done by an administrator.",
    position: "left",
  },
  {
    selector: ".deleteUser",
    title: "Delete Cards",
    text: "This allows you to delete the user's details. This can only be done by an administrator.",
    position: "left",
  },
  {
    selector: ".usersPaging",
    title: "Users Page Navigation",
    text: "Here you can navigate to view multiple pages of users.",
    position: "right",
  },
];

const simulationResultsTutorialSteps: TutorialStep[] = [
  {
    selector: ".simName",
    title: "Simulation Name",
    text: "This displays the name of the current simulation being analyzed. You can customize this when creating simulations.",
    position: "bottom",
  },
  {
    selector: ".simDesc",
    title: "Simulation Description",
    text: "This shows the detailed description of what this simulation is testing or analyzing.",
    position: "bottom",
  },
  {
    selector: ".flex.flex-wrap.gap-2.mb-2",
    title: "Selected Intersections",
    text: "These pills show which intersections are included in this simulation. Each intersection contributes to the overall traffic analysis.",
    position: "bottom",
  },
  {
    selector: ".flex.flex-col.gap-3 button:nth-child(1)",
    title: "3D Rendering View",
    text: "Click this button to see an interactive 3D visualization of the traffic simulation in action.",
    position: "left",
  },
  {
    selector: ".flex.flex-col.gap-3 button:nth-child(2)",
    title: "Optimization Toggle",
    text: "This button lets you compare original simulation results with optimized traffic light timings. Green indicates optimization data is available.",
    position: "left",
  },
  {
    selector: ".stat-cube:nth-child(1)",
    title: "Average Speed Statistics",
    text: "This card shows the average speed of all vehicles throughout the simulation. Higher speeds generally indicate better traffic flow.",
    position: "bottom",
  },
  {
    selector: ".stat-cube:nth-child(2)",
    title: "Maximum Speed Statistics",
    text: "This displays the highest speed reached by any vehicle during the simulation period.",
    position: "bottom",
  },
  {
    selector: ".stat-cube:nth-child(3)",
    title: "Minimum Speed Statistics",
    text: "This shows the lowest speed recorded, which helps identify congestion points where vehicles slow down significantly.",
    position: "bottom",
  },
  {
    selector: ".stat-cube:nth-child(4)",
    title: "Total Distance Traveled",
    text: "This represents the cumulative distance traveled by all vehicles, indicating overall traffic volume and activity.",
    position: "bottom",
  },
  {
    selector: ".stat-cube:nth-child(5)",
    title: "Vehicle Count",
    text: "This shows the total number of vehicles that participated in the simulation.",
    position: "bottom",
  },
  {
    selector: ".stat-cube:nth-child(6)",
    title: "Traffic Light Phases",
    text: "This displays the number of different signal phases (red, yellow, green combinations) configured for the intersection.",
    position: "bottom",
  },
  {
    selector: ".stat-cube:nth-child(7)",
    title: "Cycle Duration",
    text: "This shows the total time for one complete traffic light cycle, measured in seconds.",
    position: "bottom",
  },
  {
    selector: ".grid.grid-cols-1.lg\\:grid-cols-2 > div:nth-child(1)",
    title: "Average Speed Over Time Chart",
    text: "This line chart shows how average vehicle speeds changed throughout the simulation. Look for patterns that indicate rush hours or congestion periods.",
    position: "right",
  },
  {
    selector: ".grid.grid-cols-1.lg\\:grid-cols-2 > div:nth-child(2)",
    title: "Vehicle Count Over Time Chart",
    text: "This chart tracks the number of active vehicles at each time point, helping identify peak traffic periods.",
    position: "left",
  },
  {
    selector: ".grid.grid-cols-1.lg\\:grid-cols-2 > div:nth-child(3)",
    title: "Final Speed Distribution",
    text: "This histogram shows how many vehicles ended the simulation at different speed ranges, indicating overall traffic efficiency.",
    position: "right",
  },
  {
    selector: ".grid.grid-cols-1.lg\\:grid-cols-2 > div:nth-child(4)",
    title: "Total Distance Distribution",
    text: "This histogram displays the distribution of total distances traveled by individual vehicles, showing travel pattern variations.",
    position: "left",
  },
];

const comparisonViewTutorialSteps: TutorialStep[] = [
  {
    selector: ".traffic-simulation-root:first-of-type",
    title: "Original Simulation View",
    text: "This left panel shows the original traffic simulation with your initial intersection settings. It displays the baseline traffic flow before optimization.",
    position: "right",
  },
  {
    selector: ".traffic-simulation-root:last-of-type",
    title: "Optimized Simulation View",
    text: "This right panel shows the optimized traffic simulation with improved traffic light timings. Compare this with the original to see the optimization benefits.",
    position: "left",
  },
  {
    selector:
      "div[style*='position: absolute'][style*='top: 24px'][style*='left: 24px']",
    title: "Simulation Control Panel",
    text: "This panel controls both simulations simultaneously. You can play/pause, restart, adjust speed, and monitor real-time statistics for the active simulation.",
    position: "right",
  },
  {
    selector: "div[style*='progress']",
    title: "Simulation Progress",
    text: "The progress bar shows how much of the simulation has completed. Both simulations run in sync, making comparison easy.",
    position: "bottom",
  },
  {
    selector: "div[style*='border-bottom: 1px solid']:nth-of-type(2)",
    title: "Vehicle Statistics",
    text: "Monitor total, active, and completed vehicles plus average speed in real-time. These metrics help you understand traffic efficiency differences.",
    position: "bottom",
  },
  {
    selector: "div[style*='border-top: 1px solid']",
    title: "Traffic Light Status",
    text: "See the current traffic light states for all directions (North, South, East, West). Colors indicate red, yellow, or green light phases.",
    position: "bottom",
  },
  {
    selector: "button[style*='flex-grow: 1']:first-child",
    title: "Play/Pause Control",
    text: "Control both simulations simultaneously. Play to start/resume or pause to analyze specific moments in the traffic flow.",
    position: "bottom",
  },
  {
    selector: "button[style*='flex-grow: 1']:last-child",
    title: "Restart Simulations",
    text: "Reset both simulations back to the beginning. Useful for comparing different scenarios from the start.",
    position: "bottom",
  },
  {
    selector: "input[type='range']",
    title: "Speed Control",
    text: "Adjust the simulation playback speed from 1x to 20x. Higher speeds let you observe long-term traffic patterns more quickly.",
    position: "top",
  },
  {
    selector:
      "button[title*='original']:first-of-type, button[title*='View left']:first-of-type",
    title: "Left Panel Fullscreen",
    text: "Click this button to expand the original simulation to fullscreen for detailed analysis. Click again to return to side-by-side view.",
    position: "bottom",
  },
  {
    selector:
      "button[title*='optimized']:last-of-type, button[title*='View right']:last-of-type",
    title: "Right Panel Fullscreen",
    text: "Click this button to expand the optimized simulation to fullscreen for detailed analysis. Click again to return to side-by-side view.",
    position: "bottom",
  },
  {
    selector: "button[style*='position: absolute'][style*='bottom: 70px']",
    title: "Exit Comparison View",
    text: "Click this button to close the comparison view and return to the previous page. Your analysis session will end.",
    position: "top",
  },
  {
    selector: "canvas",
    title: "3D Traffic Visualization",
    text: "Each panel contains a 3D visualization of the intersection. Watch vehicles move through the intersection and observe traffic light changes in real-time.",
    position: "center",
  },
  {
    selector: "div[style*='position: absolute'][style*='bottom: 20px']",
    title: "Simulation Labels",
    text: "These labels at the bottom of each panel clearly identify which simulation you're viewing: 'Original Simulation' vs 'Optimized Simulation'.",
    position: "top",
  },
];

const faqData = [
  {
    question: "What is Swift Signals?",
    answer:
      "Swift Signals is a simulation-powered, machine-learning-based platform that helps traffic departments optimize traffic light timing at intersections to reduce congestion and improve traffic flow.",
  },
  {
    question: "Who is Swift Signals designed for?",
    answer:
      "The platform is built for municipal traffic departments and urban planners who need a scalable, data-driven tool to monitor and improve intersection performance.",
  },
  {
    question: "What problems will Swift Signals solve?",
    answer:
      "It addresses urban traffic congestion, which costs South Africa's economy an estimated R1 billion annually in productivity losses. By optimizing traffic light cycles based on real-world data, Swift Signals aims to reduce wait times, improve vehicle throughput, and enhance overall intersection efficiency.",
  },
  {
    question: "How does the platform work?",
    answer:
      "Simulates traffic flow using historical and real-time data.<br/>Applies Swarm Optimization Algorithms to test multiple traffic light timing strategies.<br/>Selects the most efficient timing plans based on metrics like average wait time and throughput.",
  },
  {
    question: "What technologies power Swift Signals?",
    answer:
      "<b>Frontend:</b> React.js + TailwindCSS<br/><b>Backend:</b> Microservices architecture (containerized)<br/><b>Database:</b> MongoDB (for scalable time-series data storage)<br/><b>Optimization:</b> Particle Swarm Optimization (PSO) and other swarm algorithms",
  },
  {
    question: "Can Swift Signals be used for multiple intersections?",
    answer:
      "Yes, the system is designed for scalability, including support for optimizing multiple intersections simultaneously and accommodating more complex models (e.g., turn-only lanes).",
  },
  {
    question: "How is the system deployed?",
    answer:
      "Swift Signals uses modern DevOps pipelines, including containerization and CI/CD, allowing seamless deployment, updates, and modular testing of services.",
  },
  {
    question: "Does it support real-time traffic?",
    answer:
      "Currently, the focus is on historical data, but the platform is designed to integrate real-time feeds in future iterations.",
  },
  {
    question: "How are intersections configured in the system?",
    answer:
      "Through a responsive web portal, users can:<br/>- Configure custom intersection layouts and signal sequences.<br/>- Monitor simulation performance.<br/>- View optimization results via interactive reports and visualizations.",
  },
  {
    question: "What kind of reports or analytics does the system generate?",
    answer:
      "Users can access:<br/>- Visual simulations of traffic flow<br/>- Performance dashboards (e.g., wait times, flow efficiency)<br/>- Alerts and improvement suggestions<br/>- Exportable optimization reports",
  },
  {
    question: "Is this product finished or in development?",
    answer:
      "Swift Signals is currently under active development using an Agile process, with updates and new features released every sprint. Stakeholder feedback is integrated regularly to align with real-world needs.",
  },
  {
    question: "How can I start using Swift Signals?",
    answer:
      "You can reach out to Southern Cross Solutions for pilot access, deployment support, and integration consultation. Full release timelines will be communicated upon request.",
  },
];

const HelpMenu: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [activeTab, setActiveTab] = useState("chat");
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [userInput, setUserInput] = useState("");
  const [isBotTyping, setIsBotTyping] = useState(false);
  const [sessionId] = useState<string>(uuidv4());
  const chatBodyRef = useRef<HTMLDivElement | null>(null);
  const [openSections, setOpenSections] = useState<Record<string, boolean>>({});
  const [activeTutorial, setActiveTutorial] = useState<TutorialType | null>(
    null,
  );
  const [confirmationDetails, setConfirmationDetails] = useState<{
    pageName: string;
    path: string;
    tutorialType: TutorialType;
  } | null>(null);

  const [openFaqIndex, setOpenFaqIndex] = useState<number | null>(null);

  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    const tutorialToStart = location.state?.startTutorial as TutorialType;
    if (tutorialToStart) {
      setTimeout(() => {
        setActiveTutorial(tutorialToStart);
      }, 150);
      window.history.replaceState({}, document.title);
    }
  }, [location]);

  useEffect(() => {
    if (chatBodyRef.current) {
      chatBodyRef.current.scrollTop = chatBodyRef.current.scrollHeight;
    }
  }, [messages, isBotTyping]);

  const sendQueryToBot = async (query: { text?: string; event?: string }) => {
    const { text, event } = query;

    if (text && text.trim() === "") return;

    if (text) {
      const newUserMessage: ChatMessage = { text, sender: "user" };
      setMessages((prev) => [...prev, newUserMessage]);
      setUserInput("");
    }

    setIsBotTyping(true);

    try {
      const response = await fetch("http://localhost:3001/api/chatbot", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          message: text,
          event: event,
          sessionId: sessionId,
          token: localStorage.getItem('authToken'),
        }),
      });

      if (!response.ok) throw new Error("Network response was not ok");

      const data = await response.json();

      // --- HANDLE NAVIGATION PAYLOAD ---
      if (data.fulfillmentMessages) {
        const navigationPayload = data.fulfillmentMessages.find(
          (msg: any) => msg.payload && msg.payload.fields && msg.payload.fields.action
        );

        if (navigationPayload) {
          const action = navigationPayload.payload.fields.action.stringValue;
          const path = navigationPayload.payload.fields.path.stringValue;

          if (action === 'NAVIGATE' && path) {
            console.log(`%c✅ ACTION HANDLER PASSED: Navigating to [${path}]`,
            "color: green; font-weight: bold;");
            setTimeout(() => {
              navigate(path);
              setIsOpen(false); // Close the help menu on navigation
            }, 1000); // Wait 1 second for the user to read the message
          }
        }
      }

      let quickReplies: QuickReply[] = [];
      if (data.fulfillmentMessages) {
        const payload = data.fulfillmentMessages.find(
          (msg: DialogflowMessage) => msg.payload,
        );
        if (payload?.payload?.fields?.richContent) {
          const options =
            payload.payload.fields.richContent.listValue.values[0].listValue
              .values[0].structValue.fields.options.listValue.values;
          quickReplies = options.map((option: DialogflowQuickReplyOption) => ({
            text: option.structValue.fields.text.stringValue,
            payload:
              option.structValue.fields.link?.stringValue ||
              option.structValue.fields.text.stringValue,
          }));
        }
      }

      const botResponse: ChatMessage = {
        text: data.fulfillmentText,
        sender: "bot",
        quickReplies: quickReplies.length > 0 ? quickReplies : undefined,
      };

      setMessages((prev) => [...prev, botResponse]);

      if (
        data.action === "start.tutorial" &&
        data.parameters?.fields?.tutorial_topic
      ) {
        const tutorialType = data.parameters.fields.tutorial_topic
          .stringValue as TutorialType;

        if (tutorialType) {
          console.log(
            `%c✅ ACTION HANDLER PASSED: Starting tutorial for [${tutorialType}]`,
            "color: green; font-weight: bold;",
          );
          setTimeout(() => {
            startTutorial(tutorialType);
          }, 500);
        }
      }
    } catch (error) {
      console.error("Error communicating with chatbot backend:", error);
      const errorResponse: ChatMessage = {
        text: "Sorry, I'm having trouble connecting to my brain right now. Please try again later.",
        sender: "bot",
      };
      setMessages((prev) => [...prev, errorResponse]);
    } finally {
      setIsBotTyping(false);
    }
  };

  useEffect(() => {
    if (isOpen && messages.length === 0) {
      sendQueryToBot({ event: "WELCOME" });
    }
  }, [isOpen]);

  const startTutorial = (tutorialType: TutorialType) => {
    const tutorialConfig = {
      dashboard: { path: "/dashboard", name: "Dashboard" },
      intersections: { path: "/intersections", name: "Intersections" },
      simulations: { path: "/simulations", name: "Simulations" },
      users: { path: "/users", name: "Users" },
      navigation: { path: "", name: "Navigation" },
      "simulation-results": {
        path: "/simulation-results",
        name: "Simulation Results",
      },
      "comparison-view": {
        path: "/comparison-rendering",
        name: "3D Comparison View",
      },
    };

    const config = tutorialConfig[tutorialType];
    if (!config) return;

    if (!config.path) {
      setIsOpen(false);
      setActiveTutorial(tutorialType);
      return;
    }

    if (location.pathname !== config.path) {
      setConfirmationDetails({
        pageName: config.name,
        path: config.path,
        tutorialType: tutorialType,
      });
    } else {
      setIsOpen(false);
      setActiveTutorial(tutorialType);
    }
  };

  const handleConfirmNavigation = () => {
    if (!confirmationDetails) return;

    navigate(confirmationDetails.path, {
      state: { startTutorial: confirmationDetails.tutorialType },
    });

    setConfirmationDetails(null);
    setIsOpen(false);
  };

  const toggleSection = (section: string) => {
    setOpenSections((prev) => ({ ...prev, [section]: !prev[section] }));
  };

  const toggleFaq = (index: number) => {
    setOpenFaqIndex((prevIndex) => (prevIndex === index ? null : index));
  };

  return (
    <>
      {activeTutorial === "dashboard" && (
        <InteractiveTutorial
          steps={dashboardTutorialSteps}
          onClose={() => setActiveTutorial(null)}
        />
      )}
      {activeTutorial === "intersections" && (
        <InteractiveTutorial
          steps={intersectionTutorialSteps}
          onClose={() => setActiveTutorial(null)}
        />
      )}
      {activeTutorial === "simulations" && (
        <InteractiveTutorial
          steps={simulationsTutorialSteps}
          onClose={() => setActiveTutorial(null)}
        />
      )}
      {activeTutorial === "users" && (
        <InteractiveTutorial
          steps={usersTutorialSteps}
          onClose={() => setActiveTutorial(null)}
        />
      )}
      {activeTutorial === "navigation" && (
        <InteractiveTutorial
          steps={navigationTutorialSteps}
          onClose={() => setActiveTutorial(null)}
        />
      )}
      {activeTutorial === "simulation-results" && (
        <InteractiveTutorial
          steps={simulationResultsTutorialSteps}
          onClose={() => setActiveTutorial(null)}
        />
      )}
      {activeTutorial === "comparison-view" && (
        <InteractiveTutorial
          steps={comparisonViewTutorialSteps}
          onClose={() => setActiveTutorial(null)}
        />
      )}

      {confirmationDetails && (
        <div className="confirmation-overlay">
          <div className="confirmation-popup">
            <h4>Switch to {confirmationDetails.pageName}?</h4>
            <p>
              The {confirmationDetails.pageName} Tutorial is best viewed on the{" "}
              {confirmationDetails.pageName} page. Would you like to go there
              now?
            </p>
            <div className="confirmation-buttons">
              <button onClick={() => setConfirmationDetails(null)}>No</button>
              <button onClick={handleConfirmNavigation}>Yes</button>
            </div>
          </div>
        </div>
      )}

      <div className={`help-container ${isOpen ? "open" : ""}`}>
        <button className="help-button" onClick={() => setIsOpen(!isOpen)}>
          {isOpen ? (
            <FaTimes />
          ) : (
            <>
              {" "}
              <FaChevronLeft className="help-button-arrow" />{" "}
              <span className="help-button-text">HELP</span>{" "}
            </>
          )}
        </button>

        <div className="help-menu">
          <div className="help-menu-header">
            <button
              className="close-help-menu-button"
              onClick={() => setIsOpen(false)}
            >
              {" "}
              <FaTimes />{" "}
            </button>
            <div className="help-menu-tabs">
              <button
                className={`help-tab-button ${activeTab === "chat" ? "active" : ""}`}
                onClick={() => setActiveTab("chat")}
              >
                {" "}
                <FaCommentDots /> Swift Chat{" "}
              </button>
              <button
                className={`help-tab-button ${activeTab === "general" ? "active" : ""}`}
                onClick={() => setActiveTab("general")}
              >
                {" "}
                <FaBook /> General Help{" "}
              </button>
            </div>
            <div className="header-spacer" />
          </div>

          {activeTab === "chat" ? (
            <div className="chatbot-container">
              <div className="chatbot-body" ref={chatBodyRef}>
                {messages.map((msg, index) => (
                  <div key={index} className={`message-wrapper ${msg.sender}`}>
                    <div className="chat-message">
                      <p
                        dangerouslySetInnerHTML={{
                          __html: msg.text.replace(/\n/g, "<br />"),
                        }}
                      />
                    </div>
                    {msg.quickReplies && (
                      <div className="quick-replies">
                        {msg.quickReplies.map((reply, i) => (
                          <button
                            key={i}
                            onClick={() =>
                              sendQueryToBot({ text: reply.payload })
                            }
                          >
                            {reply.text}
                          </button>
                        ))}
                      </div>
                    )}
                  </div>
                ))}
                {isBotTyping && (
                  <div className="message-wrapper bot">
                    <div className="chat-message">
                      <div className="typing-indicator">
                        <span></span>
                        <span></span>
                        <span></span>
                      </div>
                    </div>
                  </div>
                )}
              </div>
              <div className="chatbot-input">
                <input
                  type="text"
                  placeholder="Type your message..."
                  value={userInput}
                  onChange={(e) => setUserInput(e.target.value)}
                  onKeyPress={(e) =>
                    e.key === "Enter" && sendQueryToBot({ text: userInput })
                  }
                />
                <button onClick={() => sendQueryToBot({ text: userInput })}>
                  {" "}
                  <IoSend />{" "}
                </button>
              </div>
            </div>
          ) : (
            <div className="general-help-container">
              <div className="accordion-section">
                <button
                  className="accordion-header"
                  onClick={() => toggleSection("tutorials")}
                >
                  <span>Tutorials</span>
                  <FaChevronDown
                    className={`accordion-icon ${openSections["tutorials"] ? "open" : ""}`}
                  />
                </button>
                <div
                  className={`accordion-content ${openSections["tutorials"] ? "open" : ""}`}
                >
                  <div className="accordion-item tutorial-launcher">
                    <button onClick={() => startTutorial("navigation")}>
                      <h4>Navigation Tutorial</h4>
                      <p>Learn how to use the site's navbar and footer.</p>
                    </button>
                  </div>
                  <div className="accordion-item tutorial-launcher">
                    <button onClick={() => startTutorial("dashboard")}>
                      <h4>Dashboard Tutorial</h4>
                      <p>An interactive walkthrough of the main dashboard.</p>
                    </button>
                  </div>
                  <div className="accordion-item tutorial-launcher">
                    <button onClick={() => startTutorial("intersections")}>
                      <h4>Intersections Tutorial</h4>
                      <p>Learn how to search, add, and manage intersections.</p>
                    </button>
                  </div>
                  <div className="accordion-item tutorial-launcher">
                    <button onClick={() => startTutorial("simulations")}>
                      <h4>Simulations Tutorial</h4>
                      <p>Learn how to run simulations and optimizations.</p>
                    </button>
                  </div>
                  <div className="accordion-item tutorial-launcher">
                    <button onClick={() => startTutorial("users")}>
                      <h4>Users Tutorial</h4>
                      <p>Learn how to run view, edit, and delete users.</p>
                    </button>
                  </div>
                  <div className="accordion-item tutorial-launcher">
                    <button onClick={() => startTutorial("simulation-results")}>
                      <h4>Simulation Results Tutorial</h4>
                      <p>
                        Learn how to analyze simulation data, charts, and
                        statistics.
                      </p>
                    </button>
                  </div>
                  <div className="accordion-item tutorial-launcher">
                    <button onClick={() => startTutorial("comparison-view")}>
                      <h4>3D Comparison View Tutorial</h4>
                      <p>
                        Learn how to compare original vs optimized simulations
                        in 3D.
                      </p>
                    </button>
                  </div>
                </div>
              </div>
              <div className="accordion-section">
                <button
                  className="accordion-header"
                  onClick={() => toggleSection("faq")}
                >
                  <span>Frequently Asked Questions</span>
                  <FaChevronDown
                    className={`accordion-icon ${openSections["faq"] ? "open" : ""}`}
                  />
                </button>
                <div
                  className={`accordion-content ${openSections["faq"] ? "open" : ""}`}
                >
                  <div className="faq-list">
                    {faqData.map((item, index) => (
                      <div key={index} className="faq-item">
                        <button
                          className="faq-question"
                          onClick={() => toggleFaq(index)}
                        >
                          <span>{item.question}</span>
                          <FaChevronDown
                            className={`faq-icon ${openFaqIndex === index ? "open" : ""}`}
                          />
                        </button>
                        <div
                          className={`faq-answer ${openFaqIndex === index ? "open" : ""}`}
                        >
                          <div className="faq-answer-content">
                            <p
                              dangerouslySetInnerHTML={{ __html: item.answer }}
                            />
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </>
  );
};

export default HelpMenu;
