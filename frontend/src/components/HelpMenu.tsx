import React, { useEffect, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import "../styles/HelpMenu.css";
import InteractiveTutorial, { type TutorialStep } from "./InteractiveTutorial";

// Icons
import { FaTimes, FaCommentDots, FaBook, FaChevronLeft, FaChevronDown } from "react-icons/fa";
import { IoSend } from "react-icons/io5";

import { intents } from "../lib/botLogic";
import type { ChatResponse } from "../lib/botLogic";

// Import all tutorial preview images
import DashboardTutPreview from "../assets/Dashboard_Tutorial.png";
import NavigationTutPreview from "../assets/Navigation_Tutorial.png";
import IntersectionsTutPreview from "../assets/Intersections_Tutorial.png";
// --- FIXED: Corrected image path for simulations tutorial ---
import SimulationsTutPreview from "../assets/Intersections_Tutorial.png"; 

// Other types
type QuickReply = { text: string; payload: string; };
type ChatMessage = { text: string; sender: "user" | "bot"; quickReplies?: QuickReply[]; };
type TutorialType = 'dashboard' | 'navigation' | 'intersections' | 'simulations' | 'users';

// --- TUTORIAL STEP DEFINITIONS ---
const dashboardTutorialSteps: TutorialStep[] = [
    { selector: '.card-grid', title: 'Summary Cards', text: 'These cards give you a quick, at-a-glance overview of your key metrics.', position: 'bottom' },
    { selector: '.recent-simulations-tab', title: 'Simulations Table', text: 'Here you can see a list of all your recent simulations. Click on any row to see more details.', position: 'right' },
    { selector: '.quick-action-button.bg-customIndigo', title: 'Add a New Intersection', text: 'Click this button to open the form for creating a new traffic intersection.', position: 'bottom' },
    { selector: '.quick-action-button.bg-customGreen', title: 'Run a Simulation', text: 'Click this button to open the form for running a traffic simulation.', position: 'bottom' },
    { selector: '.quick-action-button.bg-customPurple', title: 'View Map', text: 'This will take you to a full-screen map view of all your monitored intersections.', position: 'bottom' },
    { selector: '.graph-card', title: 'Traffic Volume Chart', text: 'This chart shows the traffic volume over time for your key intersections.', position: 'left' },
    { selector: '.inter-card', title: 'Top Intersections', text: 'This card displays the top intersections based on traffic volume.', position: 'left' },
];

const navigationTutorialSteps: TutorialStep[] = [
    { selector: '.nav-links', title: 'Main Navigation', text: 'Use these links to switch between the main pages of the application.', position: 'bottom' },
    { selector: '.user-profile', title: 'User Profile', text: 'Access your profile, settings, or log out from this menu.', position: 'bottom' },
    { selector: '.dark-mode-toggle', title: 'Appearance Toggle', text: 'Switch between light and dark modes.', position: 'top' }
];

const intersectionTutorialSteps: TutorialStep[] = [
    { selector: '.searchContainer', title: 'Search Bar', text: 'This allows you to quickly find intersections by name or ID.', position: 'bottom' },
    { selector: '.addIntersectionBtn', title: 'Add Intersection', text: 'Click this button to open the form for adding a new traffic intersection.', position: 'bottom' },
    { selector: '.intersectionCard', title: 'Intersection Cards', text: 'Each card represents a traffic intersection, displaying key information.', position: 'left' },
    { selector: '.intersectionBtn.bg-blue-600', title: 'Simulate Button', text: 'Click this button to run a traffic simulation for the selected intersection.', position: 'right' },
    { selector: '.intersectionBtn.bg-green-600', title: 'Edit Button', text: 'Click this button to edit the details of the selected intersection.', position: 'right' },
    { selector: '.intersectionBtn.bg-red-600', title: 'Delete Button', text: 'Click this button to delete the selected intersection.', position: 'right' }
];

const simulationsTutorialSteps: TutorialStep[] = [
    { selector: '.sims', title: 'Simulations', text: 'This page shows your recent simulations.', position: 'right' },
    { selector: '.opts', title: 'Optimizations', text: 'This page shows your recent optimizations.', position: 'left' },
    { selector: '.viewBtn', title: 'View a Simulation', text: 'This button let\'s you view a simulation.', position: 'left' },
    { selector: '.deleteBtn', title: 'Delete a Simulation', text: 'This button let\'s you delete a simulation.', position: 'left' },
    { selector: '.pagination', title: 'Cycle Through Pages', text: 'Here you can navigate to view multiple pages of simulations.', position: 'right' },
    { selector: '.new-simulation-button', title: 'Create a New Simulation', text: 'Let\'s see how to create a new simulation. The tutorial will now open the form for you.', position: 'bottom' },
    {
        selector: 'body',
        title: 'Opening Form',
        text: 'Please wait...',
        position: 'center',
        action: () => {
            const button = document.querySelector('.new-simulation-button') as HTMLElement;
            if (button) button.click();
        },
    },
    { selector: '.simulation-modal-content', title: 'New Simulation Form', text: 'In this form, you can define all the parameters for your new simulation.', position: 'left' },
    { selector: '.simulation-name-input', title: 'Name and Description', text: 'Give your simulation a unique name and an optional description so you can easily identify it later.', position: 'right' },
    { selector: '.intersection-tabs', title: 'Add Intersections', text: 'You can add intersections to your simulation from a pre-defined list, by searching, or by selecting them on a map.', position: 'left' },
    { selector: '.create-simulation-submit-btn', title: 'Create Simulation', text: 'Once you have filled out the form, click here to create and run your simulation.', position: 'right' }
];

const usersTutorialSteps: TutorialStep[] = [
    { selector: '.usersTable', title: 'Users Table', text: 'This displays all the users currently signed in.', position: 'left' },
    { selector: '.p-2.bg-green-500', title: 'Edit User', text: 'This allows you to edit the user\'s details. This can only be done by an administrator.', position: 'left' },
    { selector: '.p-2.bg-red-500', title: 'Delete Cards', text: 'This allows you to delete the user\'s details. This can only be done by an administrator.', position: 'left' },
    { selector: '.usersPaging', title: 'Users Page Navigation', text: 'Here you can navigate to view multiple pages of users.', position: 'right' },
];

const faqData = [ { question: "What do the different status colors mean?", answer: "Green indicates optimal traffic flow. Yellow suggests moderate congestion. Red signals heavy congestion or an incident. Grey means the intersection is offline or data is unavailable." }, { question: "How often is the traffic data updated?", answer: "Traffic data is updated in real-time, with a typical delay of less than 5 seconds." }, { question: "Can I export data from a simulation?", answer: "Yes, on the simulation results page, you will find an 'Export' button that allows you to download the data in various formats like CSV or PDF." } ];

const HelpMenu: React.FC = () => {
    const [isOpen, setIsOpen] = useState(false);
    const [activeTab, setActiveTab] = useState("chat");
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [userInput, setUserInput] = useState("");
    const [isBotTyping, setIsBotTyping] = useState(false);
    const [context, setContext] = useState<string | null>(null);
    const chatBodyRef = useRef<HTMLDivElement | null>(null);
    const [openSections, setOpenSections] = useState<Record<string, boolean>>({});

    const [activeTutorial, setActiveTutorial] = useState<TutorialType | null>(null);

    const [confirmationDetails, setConfirmationDetails] = useState<{
        pageName: string;
        path: string;
        tutorialType: TutorialType;
    } | null>(null);

    const location = useLocation();
    const navigate = useNavigate();

    useEffect(() => {
        const tutorialToStart = location.state?.startTutorial as TutorialType;
        if (tutorialToStart === 'dashboard' || tutorialToStart === 'intersections' || tutorialToStart === 'simulations' || tutorialToStart === 'users') {
            setTimeout(() => {
                setActiveTutorial(tutorialToStart);
            }, 150);
            window.history.replaceState({}, document.title);
        }
    }, [location]);

    useEffect(() => { if (chatBodyRef.current) { chatBodyRef.current.scrollTop = chatBodyRef.current.scrollHeight; } }, [messages, isBotTyping]);
    useEffect(() => { if (isOpen && messages.length === 0) { const welcomeMessage: ChatMessage = { text: "Hello! I'm here to help. What can I assist you with today?", sender: "bot", quickReplies: [ { text: "Tell me about simulations", payload: "simulations" }, { text: "How do I see the map?", payload: "map" }, { text: "What do the statuses mean?", payload: "status_colors" }, ], }; setMessages([welcomeMessage]); } }, [isOpen, messages.length]);
    const getBotResponse = ( input: string, currentContext: string | null ): { response: ChatResponse; newContext: string | null } => { const text = input.toLowerCase(); let intent = intents.find(i => i.name.toLowerCase() === text); if (!intent) { intent = intents.find(i => i.name === currentContext); } if (!intent) { intent = intents.find(i => i.keywords.some(k => text.includes(k))); } if (intent) { return { response: intent.getResponse(), newContext: intent.nextContext !== undefined ? intent.nextContext : currentContext, }; } return { response: { text: "I'm sorry, I don't understand that. Could you try rephrasing? You can ask me about 'simulations', 'maps', or the 'chart'.", sender: "bot", }, newContext: null, }; };
    const handleSendMessage = (text: string) => { if (text.trim() === "") return; const newUserMessage: ChatMessage = { text, sender: "user" }; setMessages(prev => [...prev, newUserMessage]); setUserInput(""); setIsBotTyping(true); setTimeout(() => { const { response, newContext } = getBotResponse(text, context); setMessages(prev => [...prev, response]); setContext(newContext); setIsBotTyping(false); }, 1200); };
    const toggleSection = (section: string) => { setOpenSections(prev => ({ ...prev, [section]: !prev[section] })); };

    const startTutorial = (tutorialType: TutorialType) => {
        const tutorialConfig = {
            dashboard: { path: '/dashboard', name: 'Dashboard' },
            intersections: { path: '/intersections', name: 'Intersections' },
            simulations: { path: '/simulations', name: 'Simulations' },
            users: { path: '/users', name: 'Users' },
            navigation: { path: '', name: 'Navigation' }
        };

        const config = tutorialConfig[tutorialType];

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
            state: { startTutorial: confirmationDetails.tutorialType } 
        });

        setConfirmationDetails(null);
        setIsOpen(false);
    };

    return (
        <>
            {activeTutorial === 'dashboard' && <InteractiveTutorial steps={dashboardTutorialSteps} onClose={() => setActiveTutorial(null)} />}
            {activeTutorial === 'intersections' && <InteractiveTutorial steps={intersectionTutorialSteps} onClose={() => setActiveTutorial(null)} />}
            {activeTutorial === 'simulations' && <InteractiveTutorial steps={simulationsTutorialSteps} onClose={() => setActiveTutorial(null)} />}
            {activeTutorial === 'users' && <InteractiveTutorial steps={usersTutorialSteps} onClose={() => setActiveTutorial(null)} />}
            {activeTutorial === 'navigation' && <InteractiveTutorial steps={navigationTutorialSteps} onClose={() => setActiveTutorial(null)} />}

            {confirmationDetails && (
                <div className="confirmation-overlay">
                    <div className="confirmation-popup">
                        <h4>Switch to {confirmationDetails.pageName}?</h4>
                        <p>The {confirmationDetails.pageName} Tutorial is best viewed on the {confirmationDetails.pageName} page. Would you like to go there now?</p>
                        <div className="confirmation-buttons">
                            <button onClick={() => setConfirmationDetails(null)}>No</button>
                            <button onClick={handleConfirmNavigation}>Yes</button>
                        </div>
                    </div>
                </div>
            )}

            <div className={`help-container ${isOpen ? "open" : ""}`}>
                <button className="help-button" onClick={() => setIsOpen(!isOpen)}>
                    {isOpen ? <FaTimes /> : ( <> <FaChevronLeft className="help-button-arrow" /> <span className="help-button-text">HELP</span> </> )}
                </button>

                <div className="help-menu">
                    <div className="help-menu-header">
                        <button className="close-help-menu-button" onClick={() => setIsOpen(false)} > <FaTimes /> </button>
                        <div className="help-menu-tabs">
                            <button className={`help-tab-button ${activeTab === "chat" ? "active" : ""}`} onClick={() => setActiveTab("chat")} > <FaCommentDots /> Swift Chat </button>
                            <button className={`help-tab-button ${activeTab === "general" ? "active" : ""}`} onClick={() => setActiveTab("general")} > <FaBook /> General Help </button>
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
                                                 <span></span><span></span><span></span>
                                             </div>
                                         </div>
                                     </div>
                                 )}
                             </div>
                             <div className="chatbot-input">
                                 <input type="text" placeholder="Type your message..." value={userInput} onChange={(e) => setUserInput(e.target.value)} onKeyPress={(e) => e.key === "Enter" && handleSendMessage(userInput)} />
                                 <button onClick={() => handleSendMessage(userInput)}> <IoSend /> </button>
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
                                        <button onClick={() => startTutorial('navigation')}>
                                            <h4>Navigation Tutorial</h4>
                                            <p>Learn how to use the site's navbar and footer.</p>
                                            {/* <div className="tutorial-launcher-image-container">
                                                <img src={NavigationTutPreview} alt="Navigation preview" />
                                            </div> */}
                                        </button>
                                    </div>
                                    <div className="accordion-item tutorial-launcher">
                                        <button onClick={() => startTutorial('dashboard')}>
                                            <h4>Dashboard Tutorial</h4>
                                            <p>An interactive walkthrough of the main dashboard.</p>
                                            {/* <div className="tutorial-launcher-image-container">
                                                <img src={DashboardTutPreview} alt="Dashboard preview" />
                                            </div> */}
                                        </button>
                                    </div>
                                    <div className="accordion-item tutorial-launcher">
                                        <button onClick={() => startTutorial('intersections')}>
                                            <h4>Intersections Tutorial</h4>
                                            <p>Learn how to search, add, and manage intersections.</p>
                                            {/* <div className="tutorial-launcher-image-container">
                                                <img src={IntersectionsTutPreview} alt="Intersections preview" />
                                            </div> */}
                                        </button>
                                    </div>
                                    <div className="accordion-item tutorial-launcher">
                                        <button onClick={() => startTutorial('simulations')}>
                                            <h4>Simulations Tutorial</h4>
                                            <p>Learn how to run simulations and optimizations.</p>
                                            {/* <div className="tutorial-launcher-image-container">
                                                <img src={SimulationsTutPreview} alt="Simulations preview" />
                                            </div> */}
                                        </button>
                                    </div>
                                    <div className="accordion-item tutorial-launcher">
                                        <button onClick={() => startTutorial('users')}>
                                            <h4>Users Tutorial</h4>
                                            <p>Learn how to run view, edit, and delete users.</p>
                                            {/* <div className="tutorial-launcher-image-container">
                                                <img src={SimulationsTutPreview} alt="Simulations preview" />
                                            </div> */}
                                        </button>
                                    </div>
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
