import React, { useState, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import TrafficSimulation from "./TrafficSimulation";
import HelpMenu from "../components/HelpMenu";

interface ComparisonViewProps {
  originalIntersectionId?: string;
  optimizedIntersectionId?: string;
}

const ComparisonView: React.FC<ComparisonViewProps> = ({
  originalIntersectionId: propOriginalId,
  optimizedIntersectionId: propOptimizedId,
}) => {
  const location = useLocation();
  const navigate = useNavigate();
  
  // Get intersection IDs from props or location state
  const [originalIntersectionId, setOriginalIntersectionId] = useState<string>(
    propOriginalId || location.state?.originalIntersectionId || "1"
  );
  const [optimizedIntersectionId, setOptimizedIntersectionId] = useState<string>(
    propOptimizedId || location.state?.optimizedIntersectionId || "2"
  );
  const [originalIntersectionName, setOriginalIntersectionName] = useState<string>(
    location.state?.originalIntersectionName || "Original Simulation"
  );
  const [optimizedIntersectionName, setOptimizedIntersectionName] = useState<string>(
    location.state?.optimizedIntersectionName || "Optimized Simulation"
  );
  const [expanded, setExpanded] = useState<"none" | "left" | "right">("none");

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
      if (event.key === 'Escape') {
        handleExit();
      }
    };
    
    document.addEventListener('keydown', handleKeyDown);
    
    return () => {
      document.title = "Swift Signals";
      // Restore original body styles
      document.body.style.overflow = originalOverflow;
      document.body.style.margin = originalMargin;
      document.body.style.padding = originalPadding;
      // Remove keyboard event listener
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, [originalIntersectionName]);

  const containerStyle: React.CSSProperties = {
    display: "flex",
    flexDirection: "row",
    width: "100vw",
    height: "100vh",
    backgroundColor: "#1e1e1e",
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
    backgroundColor: "rgba(0,0,0,0.75)",
    color: "white",
    padding: "8px 16px",
    borderRadius: "8px",
    fontSize: "1em",
    fontWeight: "600",
    pointerEvents: "none",
  };

  const dividerStyle: React.CSSProperties = {
    flexShrink: 0,
    width: "2px",
    backgroundColor: "#333",
    transition: "width 0.5s ease-in-out",
  };

  const modernButtonStyle: React.CSSProperties = {
    position: "absolute",
    top: "24px",
    zIndex: 10,
    background:
      "linear-gradient(135deg, rgba(255,255,255,0.15) 0%, rgba(255,255,255,0.08) 100%)",
    backdropFilter: "blur(16px)",
    border: "1px solid rgba(255,255,255,0.18)",
    borderRadius: "16px",
    padding: "12px 20px",
    cursor: "pointer",
    fontWeight: "600",
    fontSize: "14px",
    color: "rgba(255,255,255,0.95)",
    transition: "all 0.3s cubic-bezier(0.4, 0, 0.2, 1)",
    display: "flex",
    alignItems: "center",
    gap: "8px",
    boxShadow: "0 8px 32px rgba(0,0,0,0.3)",
    minWidth: "140px",
    justifyContent: "center",
  };

  const buttonHoverStyle: React.CSSProperties = {
    background:
      "linear-gradient(135deg, rgba(255,255,255,0.25) 0%, rgba(255,255,255,0.15) 100%)",
    transform: "translateY(-2px)",
    boxShadow: "0 12px 40px rgba(0,0,0,0.4)",
    border: "1px solid rgba(255,255,255,0.3)",
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
    background:
      "linear-gradient(135deg, rgba(220,38,38,0.9) 0%, rgba(185,28,28,0.8) 100%)",
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
    boxShadow: "0 8px 32px rgba(220,38,38,0.3)",
    minWidth: "120px",
    justifyContent: "center",
  };

  const exitButtonHoverStyle: React.CSSProperties = {
    background:
      "linear-gradient(135deg, rgba(239,68,68,0.95) 0%, rgba(220,38,38,0.9) 100%)",
    transform: "translateX(-50%) translateY(2px)",
    boxShadow: "0 12px 40px rgba(220,38,38,0.4)",
    border: "1px solid rgba(255,255,255,0.3)",
  };

  const toggleLeft = () =>
    setExpanded((prev) => (prev === "left" ? "none" : "left"));
  const toggleRight = () => {
    if (optimizedIntersectionId && optimizedIntersectionId !== "2") {
      setExpanded((prev) => (prev === "right" ? "none" : "right"));
    } else {
      // Show info about no optimization available
      alert("No optimization available for this simulation. Run an optimization first to enable comparison.");
    }
  };

  const handleExit = () => {
    if (window.history.length > 1) {
      window.history.back();
    } else {
      window.close();
    }
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
      if (side === "right" && (!optimizedIntersectionId || optimizedIntersectionId === "2")) {
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
            background-color: #1e1e1e !important;
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
          {optimizedIntersectionId && optimizedIntersectionId !== "2" ? (
            <>
              <TrafficSimulation
                intersectionId={optimizedIntersectionId}
                scale={expanded === "right" ? 1.0 : 0.65}
                isExpanded={expanded === "right"}
              />
              <div style={labelStyle} className="comparison-label">
                {optimizedIntersectionName}
              </div>
              <div className="comparison-button">
                <ModernButton side="right" onClick={toggleRight} position="right" />
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
                  backgroundColor: "#3d3d3d",
                  color: "white",
                  textAlign: "center",
                  padding: "2rem",
                }}
              >
                <div>
                  <h3 className="text-2xl font-bold mb-4">No Optimization Available</h3>
                  <p className="text-lg text-gray-300 mb-6">
                    This simulation hasn't been optimized yet.
                  </p>
                  <p className="text-sm text-gray-400">
                    Run an optimization to compare results side-by-side.
                  </p>
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
