import React, { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";
import TrafficSimulation from "./TrafficSimulation";
import HelpMenu from "../components/HelpMenu";

// Define SimulationData interface to match the structure from TrafficSimulation.tsx
interface Node {
  id: string;
  x: number;
  y: number;
  type: string;
}
interface Edge {
  id: string;
  from: string;
  to: string;
  speed: number;
  lanes: number;
}
interface Position {
  time: number;
  x: number;
  y: number;
  speed: number;
}
interface VehicleData {
  id: string;
  positions: Position[];
}
interface TrafficLightPhase {
  duration: number;
  state: string;
}
interface TrafficLightState {
  time: number;
  state: string;
}
interface TrafficLightData {
  id: string;
  phases: TrafficLightPhase[];
  states?: TrafficLightState[];
}
interface Connection {
  from: string;
  to: string;
  fromLane: number;
  toLane: number;
  tl: string;
}
interface SimulationData {
  intersection: {
    nodes: Node[];
    edges: Edge[];
    trafficLights?: TrafficLightData[];
    connections: Connection[];
  };
  vehicles: VehicleData[];
}

interface ComparisonViewProps {
  originalIntersectionId?: string;
  optimizedIntersectionId?: string;
}

const ComparisonView: React.FC<ComparisonViewProps> = ({
  originalIntersectionId: propOriginalId,
  optimizedIntersectionId: propOptimizedId,
}) => {
  const location = useLocation();

  // Get simulation data from location state
  const simulationData = location.state?.simulationData as
    | SimulationData
    | undefined;
  const optimizedData = location.state?.optimizedData as
    | SimulationData
    | undefined;

  // Get intersection IDs from props or location state
  const [originalIntersectionId, setOriginalIntersectionId] = useState<string>(
    propOriginalId || location.state?.originalIntersectionId || "1",
  );
  const [optimizedIntersectionId, setOptimizedIntersectionId] =
    useState<string>(
      propOptimizedId ||
        location.state?.optimizedIntersectionId ||
        originalIntersectionId,
    );
  const [originalIntersectionName] = useState<string>(
    location.state?.originalIntersectionName || "Original Simulation",
  );
  const [optimizedIntersectionName] = useState<string>("Optimized Simulation");
  const [expanded, setExpanded] = useState<"none" | "left" | "right">("none");
  const [hasOptimizedData, setHasOptimizedData] =
    useState<boolean>(!!optimizedData);
  const [isLoadingOptimized] = useState<boolean>(false);
  const [optimizedDataError, setOptimizedDataError] = useState<string | null>(
    null,
  );
  const [optimizedDataSuccess] = useState<string | null>(null);
  const [isDarkMode, setIsDarkMode] = useState(false);

  useEffect(() => {
    const checkTheme = () => {
      const savedTheme = localStorage.getItem("theme");
      setIsDarkMode(savedTheme === "dark");
    };

    checkTheme();

    window.addEventListener("storage", checkTheme);

    return () => {
      window.removeEventListener("storage", checkTheme);
    };
  }, []);

  // Check for optimized data when the component mounts or data changes
  useEffect(() => {
    setHasOptimizedData(!!optimizedData);
    if (optimizedData) {
      setOptimizedIntersectionId(originalIntersectionId);
    }
  }, [optimizedData, originalIntersectionId]);

  // Update IDs if props change
  useEffect(() => {
    if (propOriginalId) setOriginalIntersectionId(propOriginalId);
    if (propOptimizedId) setOptimizedIntersectionId(propOptimizedId);
  }, [propOriginalId, propOptimizedId]);

  // Set document title
  useEffect(() => {
    document.title = `Traffic Simulation Comparison - ${originalIntersectionName}`;

    // Store original body styles
    const originalBodyStyle = window.getComputedStyle(document.body);
    const originalOverflow = originalBodyStyle.overflow;
    const originalMargin = originalBodyStyle.margin;
    const originalPadding = originalBodyStyle.padding;

    // Add keyboard shortcut for escape key
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        handleExit();
      }
    };

    document.addEventListener("keydown", handleKeyDown);

    return () => {
      document.title = "Swift Signals";
      // Restore original body styles
      document.body.style.overflow = originalOverflow;
      document.body.style.margin = originalMargin;
      document.body.style.padding = originalPadding;
      // Remove keyboard event listener
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [originalIntersectionName]);

  const containerStyle: React.CSSProperties = {
    display: "flex",
    flexDirection: "row",
    width: "100vw",
    height: "100vh",
    backgroundColor: isDarkMode ? "#1e1e1e" : "#f0f0f0",
    paddingBottom: "0",
    margin: "0",
    position: "fixed",
    top: "0",
    left: "0",
    zIndex: "1000",
  };

  const viewStyle: React.CSSProperties = {
    position: "relative",
    height: "100%",
    overflow: "hidden",
    transition: "flex-basis 0.5s ease-in-out",
  };

  const labelStyle: React.CSSProperties = {
    position: "absolute",
    bottom: "20px",
    left: "50%",
    transform: "translateX(-50%)",
    zIndex: 1000,
    backgroundColor: isDarkMode ? "rgba(0,0,0,0.75)" : "rgba(255,255,255,0.75)",
    color: isDarkMode ? "white" : "black",
    padding: "8px 16px",
    borderRadius: "8px",
    fontSize: "1em",
    fontWeight: "600",
    pointerEvents: "none",
  };

  const dividerStyle: React.CSSProperties = {
    flexShrink: 0,
    width: "2px",
    backgroundColor: isDarkMode ? "#333" : "#ccc",
    transition: "width 0.5s ease-in-out",
  };

  const modernButtonStyle: React.CSSProperties = {
    position: "absolute",
    top: "24px",
    zIndex: 10,
    background: isDarkMode
      ? "linear-gradient(135deg, rgba(255,255,255,0.15) 0%, rgba(255,255,255,0.08) 100%)"
      : "linear-gradient(135deg, rgba(0,0,0,0.1) 0%, rgba(0,0,0,0.05) 100%)",
    backdropFilter: "blur(16px)",
    border: isDarkMode
      ? "1px solid rgba(255,255,255,0.18)"
      : "1px solid rgba(0,0,0,0.1)",
    borderRadius: "16px",
    padding: "12px 20px",
    cursor: "pointer",
    fontWeight: "600",
    fontSize: "14px",
    color: isDarkMode ? "rgba(255,255,255,0.95)" : "rgba(0,0,0,0.8)",
    transition: "all 0.3s cubic-bezier(0.4, 0, 0.2, 1)",
    display: "flex",
    alignItems: "center",
    gap: "8px",
    boxShadow: isDarkMode
      ? "0 8px 32px rgba(0,0,0,0.3)"
      : "0 8px 32px rgba(0,0,0,0.1)",
    minWidth: "140px",
    justifyContent: "center",
  };

  const buttonHoverStyle: React.CSSProperties = {
    background: isDarkMode
      ? "linear-gradient(135deg, rgba(255,255,255,0.25) 0%, rgba(255,255,255,0.15) 100%)"
      : "linear-gradient(135deg, rgba(0,0,0,0.15) 0%, rgba(0,0,0,0.1) 100%)",
    transform: "translateY(-2px)",
    boxShadow: isDarkMode
      ? "0 12px 40px rgba(0,0,0,0.4)"
      : "0 12px 40px rgba(0,0,0,0.2)",
    border: isDarkMode
      ? "1px solid rgba(255,255,255,0.3)"
      : "1px solid rgba(0,0,0,0.2)",
  };

  const iconStyle: React.CSSProperties = {
    fontSize: "16px",
    transition: "transform 0.3s ease",
  };

  const exitButtonStyle: React.CSSProperties = {
    position: "absolute",
    bottom: "70px",
    left: "50%",
    transform: "translateX(-50%)",
    zIndex: 1001,
    background: isDarkMode
      ? "linear-gradient(135deg, rgba(220,38,38,0.9) 0%, rgba(185,28,28,0.8) 100%)"
      : "linear-gradient(135deg, rgba(255,100,100,0.9) 0%, rgba(230,80,80,0.8) 100%)",
    backdropFilter: "blur(16px)",
    border: "1px solid rgba(255,255,255,0.18)",
    borderRadius: "16px",
    padding: "12px 24px",
    cursor: "pointer",
    fontWeight: "600",
    fontSize: "14px",
    color: "rgba(255,255,255,0.95)",
    transition: "all 0.3s cubic-bezier(0.4, 0, 0.2, 1)",
    display: "flex",
    alignItems: "center",
    gap: "8px",
    boxShadow: isDarkMode
      ? "0 8px 32px rgba(220,38,38,0.3)"
      : "0 8px 32px rgba(220,38,38,0.2)",
    minWidth: "120px",
    justifyContent: "center",
  };

  const exitButtonHoverStyle: React.CSSProperties = {
    background: isDarkMode
      ? "linear-gradient(135deg, rgba(239,68,68,0.95) 0%, rgba(220,38,38,0.9) 100%)"
      : "linear-gradient(135deg, rgba(255,120,120,0.95) 0%, rgba(240,100,100,0.9) 100%)",
    transform: "translateX(-50%) translateY(2px)",
    boxShadow: isDarkMode
      ? "0 12px 40px rgba(220,38,38,0.4)"
      : "0 12px 40px rgba(220,38,38,0.3)",
    border: "1px solid rgba(255,255,255,0.3)",
  };

  const toggleLeft = () =>
    setExpanded((prev) => (prev === "left" ? "none" : "left"));
  const toggleRight = () => {
    if (hasOptimizedData) {
      setExpanded((prev) => (prev === "right" ? "none" : "right"));
    } else {
      // Show info about no optimization available
      alert(
        "No optimization available for this simulation. Run an optimization first to enable comparison.",
      );
    }
  };

  const handleExit = () => {
    if (window.history.length > 1) {
      window.history.back();
    } else {
      window.close();
    }
  };

  const handleRefreshOptimizedData = () => {
    setOptimizedDataError(
      "Refresh not available. Please go back and re-run the optimization.",
    );
  };

  const getDynamicStyles = (side: "left" | "right") => {
    const isExpanded = expanded === side;
    const isCollapsed =
      (side === "left" && expanded === "right") ||
      (side === "right" && expanded === "left");

    let flexBasis = "50%";
    if (isExpanded) flexBasis = "100%";
    if (isCollapsed) flexBasis = "0%";

    return { ...viewStyle, flex: `1 1 ${flexBasis}` };
  };

  const getButtonContent = (side: "left" | "right") => {
    const isExpanded = expanded === side;

    if (isExpanded) {
      return {
        icon: "⤢",
        text: "Exit Fullscreen",
        tooltip: "Exit fullscreen to show both views",
      };
    } else {
      if (side === "right" && !hasOptimizedData) {
        return {
          icon: "ℹ️",
          text: "Info",
          tooltip: "No optimization available for this simulation",
        };
      }
      return {
        icon: "⛶",
        text: "Fullscreen",
        tooltip: `View ${side === "left" ? "original" : "optimized"} simulation in fullscreen`,
      };
    }
  };

  const ModernButton = ({
    side,
    onClick,
    position,
  }: {
    side: "left" | "right";
    onClick: () => void;
    position: "left" | "right";
  }) => {
    const [isHovered, setIsHovered] = useState(false);
    const content = getButtonContent(side);

    return (
      <button
        onClick={onClick}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        style={{
          ...modernButtonStyle,
          ...(isHovered ? buttonHoverStyle : {}),
          [position]: "24px",
        }}
        title={content.tooltip}
      >
        <span
          style={{
            ...iconStyle,
            transform: isHovered ? "scale(1.1)" : "scale(1)",
          }}
        >
          {content.icon}
        </span>
        <span>{content.text}</span>
      </button>
    );
  };

  const ExitButton = () => {
    const [isHovered, setIsHovered] = useState(false);

    return (
      <button
        onClick={handleExit}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        style={{
          ...exitButtonStyle,
          ...(isHovered ? exitButtonHoverStyle : {}),
        }}
        title="Press Escape or click to return to previous page"
      >
        <span
          style={{
            ...iconStyle,
            transform: isHovered ? "scale(1.1)" : "scale(1)",
          }}
        >
          ✕
        </span>
        <span>Exit</span>
      </button>
    );
  };

  return (
    <>
      <style>
        {`
          /* Reset all styles for this page */
          * {
            box-sizing: border-box;
          }
          
          body {
            margin: 0 !important;
            padding: 0 !important;
            overflow: hidden !important;
            background-color: ${isDarkMode ? "#1e1e1e" : "#f0f0f0"} !important;
          }
          
          /* Ensure this component takes full control */
          #root {
            margin: 0 !important;
            padding: 0 !important;
          }
          
          /* Mobile responsiveness - stack views vertically on screens smaller than 768px */
          @media (max-width: 767px) {
            body {
              overflow-y: auto !important;
              touch-action: pan-y !important;
            }
            
            .comparison-container {
              flex-direction: column !important;
              height: 200vh !important; /* Double height to accommodate both views */
              overflow-y: auto !important;
            }
          
          .comparison-view {
            height: 100vh !important; /* Each view takes full viewport height */
            flex: 1 1 50% !important; /* Equal distribution */
            width: 100% !important; /* Ensure full width */
          }
          
          .comparison-divider {
            width: 100% !important;
            height: 2px !important;
          }
          
          .comparison-label {
            bottom: 60px !important; /* Move up to avoid overlap with controls */
          }
          
          .comparison-exit-button {
            position: fixed !important;
            bottom: 20px !important;
            left: 50% !important;
            transform: translateX(-50%) !important;
          }
          
          .comparison-button {
            top: 10px !important;
            font-size: 12px !important;
            padding: 8px 12px !important;
            min-width: 100px !important;
          }
          
          /* Fix canvas container width on mobile */
          .traffic-simulation-root > div:first-child {
            width: 100% !important;
          }
          
          /* Prevent zoom/pan interference with page scrolling on mobile */
          .traffic-simulation-root canvas {
            touch-action: pan-y !important;
          }
        }
        
        /* Tablet adjustments */
        @media (max-width: 1024px) and (min-width: 768px) {
          .comparison-button {
            font-size: 12px !important;
            padding: 10px 16px !important;
            min-width: 120px !important;
          }
        }
      `}
      </style>

      <div style={containerStyle} className="comparison-container">
        <ExitButton />

        {/* Left side: Original Simulation */}
        <div style={getDynamicStyles("left")} className="comparison-view">
          <TrafficSimulation
            intersectionId={originalIntersectionId}
            scale={expanded === "left" ? 1.0 : 0.65}
            isExpanded={expanded === "left"}
            endpoint="simulate"
            simulationData={simulationData}
            isDarkMode={isDarkMode}
          />
          <div style={labelStyle} className="comparison-label">
            {originalIntersectionName}
          </div>
          <div className="comparison-button">
            <ModernButton side="left" onClick={toggleLeft} position="right" />
          </div>
        </div>

        <div
          style={{
            ...dividerStyle,
            width: expanded === "none" ? "2px" : "0px",
          }}
          className="comparison-divider"
        />

        {/* Right side: Optimized Simulation or Message */}
        <div style={getDynamicStyles("right")} className="comparison-view">
          {isLoadingOptimized ? (
            <>
              <div
                style={{
                  height: "100%",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  backgroundColor: isDarkMode ? "#3d3d3d" : "#e0e0e0",
                  color: isDarkMode ? "white" : "black",
                  textAlign: "center",
                  padding: "2rem",
                }}
              >
                <div>
                  <div className="animate-spin inline-block w-8 h-8 border-4 border-current border-t-transparent rounded-full mb-4 mx-auto"></div>
                  <h3 className="text-xl font-bold mb-2">
                    Checking Optimization Status
                  </h3>
                  <p className="text-sm text-gray-500">
                    Verifying if optimized data is available...
                  </p>
                </div>
              </div>
              <div style={labelStyle} className="comparison-label">
                Loading...
              </div>
            </>
          ) : hasOptimizedData ? (
            <>
              <TrafficSimulation
                intersectionId={optimizedIntersectionId}
                scale={expanded === "right" ? 1.0 : 0.65}
                isExpanded={expanded === "right"}
                endpoint="optimise"
                simulationData={optimizedData}
                isDarkMode={isDarkMode}
              />
              <div style={labelStyle} className="comparison-label">
                {optimizedIntersectionName}
              </div>
              <div className="comparison-button">
                <ModernButton
                  side="right"
                  onClick={toggleRight}
                  position="right"
                />
              </div>
            </>
          ) : (
            <>
              <div
                style={{
                  height: "100%",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  backgroundColor: isDarkMode ? "#3d3d3d" : "#e0e0e0",
                  color: isDarkMode ? "white" : "black",
                  textAlign: "center",
                  padding: "2rem",
                }}
              >
                <div>
                  <h3 className="text-2xl font-bold mb-4">
                    No Optimization Available
                  </h3>
                  <p
                    className={`text-sm ${isDarkMode ? "text-gray-300" : "text-gray-600"} mb-6`}
                  >
                    This simulation hasn't been optimized yet.
                  </p>
                  <p
                    className={`text-sm ${isDarkMode ? "text-gray-400" : "text-gray-500"} mb-4`}
                  >
                    Run an optimization from the Simulation Results page to
                    compare results side-by-side.
                  </p>
                  <button
                    onClick={handleRefreshOptimizedData}
                    disabled={isLoadingOptimized}
                    className="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white rounded-md transition-colors duration-200 flex items-center gap-2 mx-auto"
                  >
                    {isLoadingOptimized ? (
                      <>
                        <div className="animate-spin inline-block w-4 h-4 border-2 border-current border-t-transparent rounded-full"></div>
                        Checking...
                      </>
                    ) : (
                      <>
                        <span>🔄</span>
                        Check for Optimization
                      </>
                    )}
                  </button>
                  {optimizedDataSuccess && (
                    <p className="text-sm text-green-400 mt-2">
                      {optimizedDataSuccess}
                    </p>
                  )}
                  {optimizedDataError && (
                    <p className="text-sm text-red-400 mt-2">
                      Error: {optimizedDataError}
                    </p>
                  )}
                </div>
              </div>
              <div style={labelStyle} className="comparison-label">
                No Optimization
              </div>
            </>
          )}
        </div>

        <HelpMenu />
      </div>
    </>
  );
};

export default ComparisonView;
