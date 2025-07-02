import React, { useEffect, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { v4 as uuidv4 } from "uuid";
import "../styles/HelpMenu.css";
import InteractiveTutorial, { type TutorialStep } from "./InteractiveTutorial";

// Icons
import {
  FaTimes,
  FaCommentDots,
  FaBook,
  FaChevronLeft,
  FaChevronDown,
} from "react-icons/fa";
import { IoSend } from "react-icons/io5";

// Other types
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
  | "users";
// Add these new types
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

// --- TUTORIAL STEP DEFINITIONS ---
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
    selector: ".quick-action-button.bg-customIndigo",
    title: "Add a New Intersection",
    text: "Click this button to open the form for creating a new traffic intersection.",
    position: "bottom",
  },
  {
    selector: ".quick-action-button.bg-customGreen",
    title: "Run a Simulation",
    text: "Click this button to open the form for running a traffic simulation.",
    position: "bottom",
  },
  {
    selector: ".quick-action-button.bg-customPurple",
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
    selector: ".dark-mode-toggle",
    title: "Appearance Toggle",
    text: "Switch between light and dark modes.",
    position: "top",
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
    text: "Click this button to open the form for adding a new traffic intersection.",
    position: "bottom",
  },
  {
    selector: ".intersectionCard",
    title: "Intersection Cards",
    text: "Each card represents a traffic intersection, displaying key information.",
    position: "left",
  },
  {
    selector: ".intersectionBtn.bg-blue-600",
    title: "Simulate Button",
    text: "Click this button to run a traffic simulation for the selected intersection.",
    position: "right",
  },
  {
    selector: ".intersectionBtn.bg-green-600",
    title: "Edit Button",
    text: "Click this button to edit the details of the selected intersection.",
    position: "right",
  },
  {
    selector: ".intersectionBtn.bg-red-600",
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
    text: "In this form, you can define all the parameters for your new simulation.",
    position: "left",
  },
  {
    selector: ".simulation-name-input",
    title: "Name and Description",
    text: "Give your simulation a unique name and an optional description so you can easily identify it later.",
    position: "right",
  },
  {
    selector: ".intersection-tabs",
    title: "Add Intersections",
    text: "You can add intersections to your simulation from a pre-defined list, by searching, or by selecting them on a map.",
    position: "left",
  },
  {
    selector: ".create-simulation-submit-btn",
    title: "Create Simulation",
    text: "Once you have filled out the form, click here to create and run your simulation.",
    position: "right",
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
    selector: ".p-2.bg-green-500",
    title: "Edit User",
    text: "This allows you to edit the user's details. This can only be done by an administrator.",
    position: "left",
  },
  {
    selector: ".p-2.bg-red-500",
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

// --- UPDATED: New, more detailed FAQ data ---
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

  // --- ADDED: State to manage the open FAQ item ---
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

  // This single, unified function handles all communication with the bot
  const sendQueryToBot = async (query: { text?: string; event?: string }) => {
    const { text, event } = query;

    // Exit if the user tries to send an empty text message
    if (text && text.trim() === "") return;

    // Add the user's message to the chat window UI only if it's a text message
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
        }),
      });

      if (!response.ok) throw new Error("Network response was not ok");

      const data = await response.json();

      let quickReplies: QuickReply[] = [];
      if (data.fulfillmentMessages) {
        // FIXED: Replaced 'any' with 'DialogflowMessage'
        const payload = data.fulfillmentMessages.find(
          (msg: DialogflowMessage) => msg.payload,
        );
        if (payload?.payload?.fields?.richContent) {
          const options =
            payload.payload.fields.richContent.listValue.values[0].listValue
              .values[0].structValue.fields.options.listValue.values;
          // FIXED: Replaced 'any' with 'DialogflowQuickReplyOption'
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

      // --- THE CORRECTED AND FINAL ACTION HANDLER ---
      // It now looks inside the 'fields' object for the lowercase 'tutorialtopic'
      if (
        data.action === "start.tutorial" &&
        data.parameters?.fields?.tutorialtopic
      ) {
        // Get the actual value from inside the object structure
        const tutorialType = data.parameters.fields.tutorialtopic
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

  // This useEffect hook correctly calls our unified function for the welcome message
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

  // --- ADDED: Handler to toggle individual FAQ items ---
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
                  {/* --- UPDATED: Renders the new nested FAQ accordion --- */}
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
