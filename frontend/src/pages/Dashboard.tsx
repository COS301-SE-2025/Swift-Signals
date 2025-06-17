import React, { useEffect, useRef, useState } from "react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import "../styles/Dashboard.css";
import { Chart, registerables } from "chart.js";

//icons
import {
  FaRoad,
  FaPlay,
  FaChartLine,
  FaPlus,
  FaMap,
  FaQuestion,
  FaTimes,
} from "react-icons/fa";
import { IoSend } from "react-icons/io5";

// Register Chart.js components
Chart.register(...registerables);

const simulations = [
  {
    id: "#1234",
    intersection: "Main St & 5th Ave",
    status: "Complete",
    statusColor: "bg-statusGreen",
    textColor: "text-statusTextGreen",
  },
  {
    id: "#1233",
    intersection: "Broadway & 7th St",
    status: "Running",
    statusColor: "bg-statusYellow",
    textColor: "text-statusTextYellow",
  },
  {
    id: "#1232",
    intersection: "Park Ave & 3rd St",
    status: "Failed",
    statusColor: "bg-statusRed",
    textColor: "text-statusTextRed",
  },
  {
    id: "#1231",
    intersection: "Broadway & 7th St",
    status: "Running",
    statusColor: "bg-statusYellow",
    textColor: "text-statusTextYellow",
  },
];

const topIntersections = [
  { name: "Main St & 5th Ave", volume: "15,000 vehicles" },
  { name: "Broadway & 7th St", volume: "13,500 vehicles" },
  { name: "Park Ave & 3rd St", volume: "12,000 vehicles" },
];

type QuickReply = {
  text: string;
  payload: string;
};

type ChatMessage = {
  text: string;
  sender: "user" | "bot";
  quickReplies?: QuickReply[];
};

const Dashboard: React.FC = () => {
  const chartRef = useRef<HTMLCanvasElement | null>(null);
  const chartInstanceRef = useRef<Chart | null>(null);
  const [isHelpMenuOpen, setIsHelpMenuOpen] = useState(false);
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
  const [userInput, setUserInput] = useState("");
  const [isBotTyping, setIsBotTyping] = useState(false);
  const [conversationContext, setConversationContext] = useState<string | null>(
    null
  );
  const chatBodyRef = useRef<HTMLDivElement | null>(null);

  // Auto-scroll to the latest message
  useEffect(() => {
    if (chatBodyRef.current) {
      chatBodyRef.current.scrollTop = chatBodyRef.current.scrollHeight;
    }
  }, [chatMessages, isBotTyping]);

  // Greet the user when the chat opens
  useEffect(() => {
    if (isHelpMenuOpen && chatMessages.length === 0) {
      const welcomeMessage: ChatMessage = {
        text: "Hello! I'm here to help. What can I assist you with today?",
        sender: "bot",
        quickReplies: [
          { text: "Tell me about simulations", payload: "simulations" },
          { text: "How do I see the map?", payload: "map" },
          { text: "What do the statuses mean?", payload: "status_colors" },
        ],
      };
      setChatMessages([welcomeMessage]);
    }
  }, [isHelpMenuOpen]);

  const handleSendMessage = (text: string) => {
    if (text.trim() === "") return;

    const newUserMessage: ChatMessage = { text, sender: "user" };
    const newMessages = [...chatMessages, newUserMessage];
    setChatMessages(newMessages);

    if (text === userInput) {
      setUserInput("");
    }

    setIsBotTyping(true);

    // Simulate bot thinking time and generate a response
    setTimeout(() => {
      const { response, newContext } = getBotResponse(
        text,
        conversationContext
      );
      setChatMessages([...newMessages, response]);
      if (newContext) {
        setConversationContext(newContext);
      }
      setIsBotTyping(false);
    }, 1200);
  };

  const getBotResponse = (
    input: string,
    context: string | null
  ): { response: ChatMessage; newContext: string | null } => {
    const text = input.toLowerCase();
    let response: ChatMessage;
    let newContext: string | null = null;

    // --- Main keywords and context handling ---
    if (
      text.includes("simulation") ||
      text.includes("run") ||
      context === "simulations"
    ) {
      newContext = "simulations";
      response = {
        text: "Simulations allow you to model traffic flow. You can see recent simulations on the main dashboard or start a new one with the 'Run Simulation' button.",
        sender: "bot",
        quickReplies: [
          { text: "How do I view details?", payload: "simulation_details" },
          {
            text: "What is an optimization run?",
            payload: "optimization_runs",
          },
          { text: "Thanks!", payload: "end_conversation" },
        ],
      };
    } else if (text.includes("map") || context === "map") {
      newContext = "map";
      response = {
        text: "You can view a map of all intersections by clicking the 'View Map' button on the dashboard. This gives you a geographical overview of your traffic network.",
        sender: "bot",
        quickReplies: [
          { text: "Tell me about intersections", payload: "intersections" },
          { text: "Thanks!", payload: "end_conversation" },
        ],
      };
    } else if (text.includes("intersection") || context === "intersections") {
      newContext = "intersections";
      response = {
        text: "Intersections are the points you are monitoring. You can add new ones via the 'New Intersection' button. The 'Total Intersections' card shows you a count of all active locations.",
        sender: "bot",
        quickReplies: [
          { text: "What about the chart?", payload: "chart" },
          { text: "Thanks!", payload: "end_conversation" },
        ],
      };
    } else if (text.includes("status") || text.includes("color")) {
      newContext = "statuses";
      response = {
        text: "The status colors indicate the state of a simulation: \n- Green (Complete): The simulation finished successfully. \n- Yellow (Running): The simulation is currently in progress. \n- Red (Failed): The simulation encountered an error.",
        sender: "bot",
        quickReplies: [
          { text: "Tell me about simulations", payload: "simulations" },
          { text: "Thanks!", payload: "end_conversation" },
        ],
      };
    } else if (text.includes("chart") || text.includes("graph")) {
      newContext = "chart";
      response = {
        text: "The Traffic Volume chart shows vehicle counts over time for key intersections, helping you visualize peak hours and traffic patterns.",
        sender: "bot",
        quickReplies: [
          { text: "What are top intersections?", payload: "top_intersections" },
          { text: "Thanks!", payload: "end_conversation" },
        ],
      };
    } else if (
      text.includes("thank") ||
      text.includes("bye") ||
      text === "end_conversation"
    ) {
      response = {
        text: "You're welcome! Let me know if you need anything else.",
        sender: "bot",
      };
      newContext = null;
    } else {
      response = {
        text: "I'm not sure I understand. Can you rephrase? You can ask me about simulations, maps, intersections, or the chart.",
        sender: "bot",
        quickReplies: [
          { text: "What are simulations?", payload: "simulations" },
          { text: "Show me map info", payload: "map" },
        ],
      };
    }
    return { response, newContext };
  };

  useEffect(() => {
    if (chartRef.current) {
      if (chartInstanceRef.current) {
        chartInstanceRef.current.destroy();
      }

      const ctx = chartRef.current.getContext("2d");
      if (!ctx) return;

      // Create gradient fill
      const gradient = ctx.createLinearGradient(0, 0, 0, 180);
      gradient.addColorStop(0, "rgba(153, 25, 21, 0.3)");
      gradient.addColorStop(1, "rgba(153, 25, 21, 0)");

      chartInstanceRef.current = new Chart(ctx, {
        type: "line",
        data: {
          labels: ["6 AM", "7 AM", "8 AM", "9 AM", "10 AM"],
          datasets: [
            {
              label: "Traffic Volume",
              data: [5000, 10000, 8000, 12000, 9000],
              fill: true,
              backgroundColor: gradient,
              borderColor: "#991915",
              borderWidth: 3,
              pointBackgroundColor: "#991915",
              pointBorderColor: "#fff",
              pointHoverRadius: 6,
              pointRadius: 4,
              pointHoverBackgroundColor: "#fff",
              pointHoverBorderColor: "#991915",
              tension: 0.4,
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          layout: {
            padding: {
              top: 10,
              bottom: 10,
              left: 0,
              right: 0,
            },
          },
          scales: {
            x: {
              grid: {
                display: false,
              },
              ticks: {
                color: "#6B7280",
                font: {
                  size: 14,
                  weight: 500,
                },
              },
              border: {
                display: false,
              },
            },
            y: {
              grid: {
                color: "#E5E7EB",
                drawTicks: false,
              },
              ticks: {
                color: "#6B7280",
                stepSize: 2000,
                font: {
                  size: 14,
                  weight: 500,
                },
              },
              border: {
                display: false,
              },
            },
          },
          plugins: {
            legend: {
              display: false,
            },
            tooltip: {
              backgroundColor: "#111827", // Tailwind's gray-900
              titleColor: "#F9FAFB", // Tailwind's gray-50
              bodyColor: "#E5E7EB", // Tailwind's gray-200
              cornerRadius: 4,
              padding: 10,
              titleFont: {
                weight: "bold",
                size: 14,
              },
              bodyFont: {
                size: 13,
              },
            },
          },
        },
      });
    }

    return () => {
      if (chartInstanceRef.current) {
        chartInstanceRef.current.destroy();
        chartInstanceRef.current = null;
      }
    };
  }, []);

  return (
    <div className="dashboard-screen min-h-screen bg-gray-100 dark:bg-gray-900">
      <Navbar />
      <div className="main-content flex-grow">
        <h1 className="Dashboard-h1">Dashboard Overview</h1>

        {/* Summary Cards */}
        <div className="card-grid">
          <div className="card">
            <div className="card-icon-1">
              <span className="text-blue-600">
                <FaRoad />
              </span>
            </div>
            <div>
              <h3 className="card-h3">Total Intersections</h3>
              <p className="card-p">24</p>
            </div>
          </div>
          <div className="card">
            <div className="card-icon-2">
              <span className="text-green-600">
                <FaPlay />
              </span>
            </div>
            <div>
              <h3 className="card-h3">Active Simulations</h3>
              <p className="card-p">8</p>
            </div>
          </div>
          <div className="card">
            <div className="card-icon-3">
              <span className="text-purple-600">
                <FaChartLine />
              </span>
            </div>
            <div>
              <h3 className="card-h3">Optimization Runs</h3>
              <p className="card-p">156</p>
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="quick-actions">
          <button className="quick-action-button bg-customIndigo text-white px-4 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2">
            <FaPlus />
            New Intersection
          </button>
          <button className="quick-action-button bg-customGreen text-white px-4 py-2 rounded-lg hover:bg-green-700 flex items-center gap-2">
            <FaPlay />
            Run Simulation
          </button>
          <button className="quick-action-button bg-customPurple text-white px-4 py-2 rounded-lg hover:bg-gray-700 flex items-center gap-2">
            <FaMap />
            View Map
          </button>
        </div>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {/* Recent Simulations */}
          <div className="recent-simulations-tab bg-white p-4 rounded-lg shadow-md">
            <h2 className="text-lg font-semibold text-gray-800 mb-4">
              Recent Simulations
            </h2>
            <table className="table-auto w-full text-left">
              <thead>
                <tr className="text-gray-600">
                  <th className="p-2">ID</th>
                  <th className="p-2">Intersection</th>
                  <th className="p-2">Status</th>
                  <th className="p-2">Actions</th>
                </tr>
              </thead>
              <tbody>
                {simulations.map((sim) => (
                  <tr key={sim.id} className="border-t">
                    <td className="p-2">{sim.id}</td>
                    <td className="p-2">{sim.intersection}</td>
                    <td className="p-2">
                      <span
                        className={`status px-2 py-1 rounded-full text-xs ${sim.statusColor} ${sim.textColor}`}
                      >
                        {sim.status}
                      </span>
                    </td>
                    <td className="p-2">
                      <button className="view-details-button text-blue-600 hover:underline">
                        View Details
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
        {/* Traffic Volume Chart and Top Intersections */}
        <div className="stats bg-white p-4 rounded-lg shadow-md">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">
            Traffic Volume
          </h2>
          <div className="traffic-graph mb-4">
            <div className="traffic-chart">
              <canvas ref={chartRef}></canvas>
            </div>
          </div>
          <div className="top-intersections">
            <h3 className="text-md font-semibold text-gray-700 mb-2">
              Top Intersections
            </h3>
            {topIntersections.map((intersection, index) => (
              <div key={index} className="flex justify-between py-2 border-t">
                <span className="text-gray-600 dark:text-gray-200">
                  {intersection.name}
                </span>
                <span className="text-gray-800 font-semibold dark:text-gray-200">
                  {intersection.volume}
                </span>
              </div>
            ))}
            <div className="total flex justify-between py-2 border-t">
              <span className="text-gray-600 font-bold dark:text-gray-100">
                Avg Daily Volume:
              </span>
              <span className="text-gray-800 font-bold dark:text-gray-100">
                12,000 vehicles
              </span>
            </div>
          </div>
        </div>
      </div>
      <Footer />

      {/* Help Button and Menu */}
      <div className={`help-container ${isHelpMenuOpen ? "open" : ""}`}>
        <button
          className="help-button"
          onClick={() => setIsHelpMenuOpen(!isHelpMenuOpen)}
        >
          {isHelpMenuOpen ? <FaTimes /> : <FaQuestion />}
        </button>
        <div className="help-menu">
          <div className="help-menu-header">
            <h2>Support Chat</h2>
          </div>
          <div className="chatbot-container">
            <div className="chatbot-body" ref={chatBodyRef}>
              {chatMessages.map((msg, index) => (
                /* We add a new class 'message-wrapper' and move the sender class here */
                <div key={index} className={`message-wrapper ${msg.sender}`}>
                  <div className="chat-message">
                    {" "}
                    {/* The sender class is removed from here */}
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
                          onClick={() => handleSendMessage(reply.payload)}
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
                  {" "}
                  {/* Also apply wrapper here for consistency */}
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
                  e.key === "Enter" && handleSendMessage(userInput)
                }
              />
              <button onClick={() => handleSendMessage(userInput)}>
                <IoSend />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
