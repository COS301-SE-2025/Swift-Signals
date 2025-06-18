import React, { useEffect, useRef, useState } from "react";
import "../styles/HelpMenu.css";
import InteractiveTutorial from "./InteractiveTutorial"; 

// --- 1. IMPORT YOUR IMAGE ---
import DashboardTutorialPreviewImage from "../assets/Dashboard_Tutorial.png"; // Adjust the path to your image

// Icons
import { FaTimes, FaCommentDots, FaBook, FaChevronLeft, FaChevronDown } from "react-icons/fa";
import { IoSend } from "react-icons/io5";

import { intents } from "../lib/botLogic";
import type { ChatResponse } from "../lib/botLogic";

// Type definitions remain the same
type QuickReply = {
  text: string;
  payload: string;
};

type ChatMessage = {
  text: string;
  sender: "user" | "bot";
  quickReplies?: QuickReply[];
};

const tutorialsData = [
    {
        title: "How to Start a Simulation",
        content: "Navigate to the 'Simulations' tab on the main dashboard and click the 'Run New Simulation' button. Fill in the required parameters and submit the form to begin."
    },
    {
        title: "Viewing Simulation Results",
        content: "Once a simulation is complete, you can click on its entry in the dashboard table to view detailed results, including performance metrics and visualizations."
    },
    {
        title: "Using the Interactive Map",
        content: "Go to the 'Map View' page to see all your monitored intersections. Click on any intersection pin to get real-time status and performance data."
    }
];

const faqData = [
    {
        question: "What do the different status colors mean?",
        answer: "Green indicates optimal traffic flow. Yellow suggests moderate congestion. Red signals heavy congestion or an incident. Grey means the intersection is offline or data is unavailable."
    },
    {
        question: "How often is the traffic data updated?",
        answer: "Traffic data is updated in real-time, with a typical delay of less than 5 seconds, ensuring you have the most current information."
    },
    {
        question: "Can I export data from a simulation?",
        answer: "Yes, on the simulation results page, you will find an 'Export' button that allows you to download the data in various formats like CSV or PDF."
    }
];


const HelpMenu: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [activeTab, setActiveTab] = useState("chat");
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [userInput, setUserInput] = useState("");
  const [isBotTyping, setIsBotTyping] = useState(false);
  const [context, setContext] = useState<string | null>(null);
  const chatBodyRef = useRef<HTMLDivElement | null>(null);
  const [openSections, setOpenSections] = useState<Record<string, boolean>>({});
  const [isTutorialActive, setIsTutorialActive] = useState(false);

  useEffect(() => {
    if (chatBodyRef.current) {
      chatBodyRef.current.scrollTop = chatBodyRef.current.scrollHeight;
    }
  }, [messages, isBotTyping]);

  useEffect(() => {
    if (isOpen && messages.length === 0) {
      const welcomeMessage: ChatMessage = {
        text: "Hello! I'm here to help. What can I assist you with today?",
        sender: "bot",
        quickReplies: [
          { text: "Tell me about simulations", payload: "simulations" },
          { text: "How do I see the map?", payload: "map" },
          { text: "What do the statuses mean?", payload: "status_colors" },
        ],
      };
      setMessages([welcomeMessage]);
    }
  }, [isOpen, messages.length]);

  const getBotResponse = (
    input: string,
    currentContext: string | null
  ): { response: ChatResponse; newContext: string | null } => {
    const text = input.toLowerCase();
    let intent = intents.find(i => i.name.toLowerCase() === text);
    if (!intent) {
        intent = intents.find(i => i.name === currentContext);
    }
    if (!intent) {
        intent = intents.find(i => i.keywords.some(k => text.includes(k)));
    }

    if (intent) {
      return {
        response: intent.getResponse(),
        newContext: intent.nextContext !== undefined ? intent.nextContext : currentContext,
      };
    }

    return {
      response: {
        text: "I'm sorry, I don't understand that. Could you try rephrasing? You can ask me about 'simulations', 'maps', or the 'chart'.",
        sender: "bot",
      },
      newContext: null,
    };
  };

  const handleSendMessage = (text: string) => {
    if (text.trim() === "") return;

    const newUserMessage: ChatMessage = { text, sender: "user" };
    setMessages(prev => [...prev, newUserMessage]);
    setUserInput("");
    setIsBotTyping(true);

    setTimeout(() => {
      const { response, newContext } = getBotResponse(text, context);
      setMessages(prev => [...prev, response]);
      setContext(newContext);
      setIsBotTyping(false);
    }, 1200);
  };
    
    const toggleSection = (section: string) => {
        setOpenSections(prev => ({
            ...prev,
            [section]: !prev[section]
        }));
    };

    const startTutorial = () => {
        setIsOpen(false); 
        setIsTutorialActive(true);
    };

  return (
    <>
      {isTutorialActive && <InteractiveTutorial onClose={() => setIsTutorialActive(false)} />}
      
      <div className={`help-container ${isOpen ? "open" : ""}`}>
        <button className="help-button" onClick={() => setIsOpen(!isOpen)}>
          {isOpen ? (
            <FaTimes />
          ) : (
            <>
              <FaChevronLeft className="help-button-arrow" />
              <span className="help-button-text">HELP</span>
            </>
          )}
        </button>

        <div className="help-menu">
          <div className="help-menu-header">
            <button
              className="close-help-menu-button"
              onClick={() => setIsOpen(false)}
            >
              <FaTimes />
            </button>
            <div className="help-menu-tabs">
              <button
                className={`help-tab-button ${activeTab === "chat" ? "active" : ""}`}
                onClick={() => setActiveTab("chat")}
              >
                <FaCommentDots />
                Swift Chat
              </button>
              <button
                className={`help-tab-button ${activeTab === "general" ? "active" : ""}`}
                onClick={() => setActiveTab("general")}
              >
                <FaBook />
                General Help
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
                              <p dangerouslySetInnerHTML={{ __html: msg.text.replace(/\n/g, "<br />") }} />
                          </div>
                          {msg.quickReplies && (
                              <div className="quick-replies">
                                  {msg.quickReplies.map((reply, i) => (
                                      <button key={i} onClick={() => handleSendMessage(reply.payload)}>
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
                      onKeyPress={(e) => e.key === "Enter" && handleSendMessage(userInput)}
                  />
                  <button onClick={() => handleSendMessage(userInput)}>
                      <IoSend />
                  </button>
              </div>
            </div>
          ) : (
            <div className="general-help-container">
              <div className="accordion-section">
                  <button className="accordion-header" onClick={() => toggleSection('tutorials')}>
                      <span>Tutorials</span>
                      <FaChevronDown className={`accordion-icon ${openSections['tutorials'] ? 'open' : ''}`} />
                  </button>
                  <div className={`accordion-content ${openSections['tutorials'] ? 'open' : ''}`}>
                      <div className="accordion-item tutorial-launcher">
                          <button onClick={startTutorial}>
                              <h4>Dashboard Tutorial</h4>

                              <div className="tutorial-launcher-image-container">
                                  <img 
                                      src={DashboardTutorialPreviewImage} 
                                      alt="Dashboard Tutorial Preview" 
                                  />
                              </div>
                              <p>Start an interactive walkthrough of the main dashboard.</p>
                          </button>
                      </div>

                      {tutorialsData.map((item, index) => (
                          <div key={index} className="accordion-item">
                              <h4>{item.title}</h4>
                              <p>{item.content}</p>
                          </div>
                      ))}
                  </div>
              </div>

              <div className="accordion-section">
                  <button className="accordion-header" onClick={() => toggleSection('faq')}>
                      <span>FAQ</span>
                      <FaChevronDown className={`accordion-icon ${openSections['faq'] ? 'open' : ''}`} />
                  </button>
                  <div className={`accordion-content ${openSections['faq'] ? 'open' : ''}`}>
                      {faqData.map((item, index) => (
                          <div key={index} className="accordion-item">
                              <h4>{item.question}</h4>
                              <p>{item.answer}</p>
                          </div>
                      ))}
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