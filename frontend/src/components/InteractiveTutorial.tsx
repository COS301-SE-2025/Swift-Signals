import React, {
  useEffect,
  useLayoutEffect,
  useState,
  useCallback,
} from "react";
import { X } from "lucide-react";

export type TutorialStep = {
  selector: string;
  title: string;
  text: string;
  position?: "top" | "bottom" | "left" | "right" | "center";
  action?: () => void;
  autoAdvance?: boolean;
  waitFor?: string; // CSS selector to wait for before auto-advancing
};

type Position = {
  highlight: React.CSSProperties;
  popover: React.CSSProperties;
  isError?: boolean;
};

type Props = {
  steps: TutorialStep[];
  onClose: () => void;
  tutorialType: string;
};

const InteractiveTutorial: React.FC<Props> = ({ steps, onClose, tutorialType }) => {
  const [stepIndex, setStepIndex] = useState(0);
  const [position, setPosition] = useState<Position | null>(null);
  const [highlightRect, setHighlightRect] = useState<DOMRect | null>(null);
  const [isDark, setIsDark] = useState(false);

  const currentStep = steps[stepIndex];

  useEffect(() => {
    const checkDarkMode = () => {
      setIsDark(document.documentElement.classList.contains("dark"));
    };

    checkDarkMode();

    const observer = new MutationObserver(checkDarkMode);
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ["class"],
    });

    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    if (currentStep && typeof currentStep.action === "function") {
      currentStep.action();
    }
  }, [currentStep]);

  // Auto-advance logic
  useEffect(() => {
    if (!currentStep?.autoAdvance) return;

    let timeoutId: NodeJS.Timeout;

    const checkForAdvance = () => {
      if (currentStep.waitFor) {
        const element = document.querySelector(currentStep.waitFor);
        if (element) {
          timeoutId = setTimeout(() => {
            if (stepIndex < steps.length - 1) {
              setStepIndex(stepIndex + 1);
            } else {
              onClose();
            }
          }, 500); // Small delay to ensure smooth transition
        } else {
          // Keep checking every 100ms
          timeoutId = setTimeout(checkForAdvance, 100);
        }
      } else {
        // Auto-advance after action without waiting for element
        timeoutId = setTimeout(() => {
          if (stepIndex < steps.length - 1) {
            setStepIndex(stepIndex + 1);
          } else {
            onClose();
          }
        }, 1000);
      }
    };

    checkForAdvance();

    return () => {
      if (timeoutId) clearTimeout(timeoutId);
    };
  }, [currentStep, stepIndex, steps.length, onClose]);

  const calculatePosition = useCallback(() => {
    if (!currentStep) return;

    if (currentStep.position === "center") {
      setPosition({
        highlight: { display: "none" },
        popover: {
          top: "50%",
          left: "50%",
          transform: "translate(-50%, -50%)",
        },
      });
      setHighlightRect(null);
      return;
    }

    const element = document.querySelector(currentStep.selector) as HTMLElement;

    if (!element) {
      setPosition({
        isError: true,
        highlight: {},
        popover: {
          top: "50%",
          left: "50%",
          transform: "translate(-50%, -50%)",
        },
      });
      setHighlightRect(null);
      return;
    }

    let rect = element.getBoundingClientRect();
    const isElementInViewport =
      rect.top >= 0 && rect.bottom <= window.innerHeight;

    if (!isElementInViewport) {
      element.scrollIntoView({ behavior: "auto", block: "center" });
      rect = element.getBoundingClientRect();
    }

    setHighlightRect(rect);

    const popoverRect = { width: 320, height: 150 };

    const highlightStyles: React.CSSProperties = {
      width: `${rect.width + 16}px`,
      height: `${rect.height + 16}px`,
      top: `${rect.top - 8}px`,
      left: `${rect.left - 8}px`,
    };

    let popoverTop = 0;
    let popoverLeft = 0;

    if (tutorialType === "navigation" && stepIndex === 2) {
      popoverTop = rect.top - popoverRect.height - 90;
      popoverLeft = rect.left + rect.width / 2 - popoverRect.width / 2;
    } else if (tutorialType === "comparison-view" && stepIndex === 2) {
      popoverTop = rect.top + 90;
      popoverLeft = rect.left + rect.width / 2 - popoverRect.width / 2;
    } else {
      switch (currentStep.position) {
        case "top":
          popoverTop = rect.top - popoverRect.height - 20;
          popoverLeft = rect.left + rect.width / 2 - popoverRect.width / 2;
          break;
        case "left":
          popoverTop = rect.top + rect.height / 2 - popoverRect.height / 2;
          popoverLeft = rect.left - popoverRect.width - 20;
          break;
        case "right":
          popoverTop = rect.top + rect.height / 2 - popoverRect.height / 2;
          popoverLeft = rect.right + 20;
          break;
        default:
          popoverTop = rect.bottom + 20;
          popoverLeft = rect.left + rect.width / 2 - popoverRect.width / 2;
          break;
      }
    }

    const popoverStyles: React.CSSProperties = {
      top: `${Math.max(20, popoverTop)}px`,
      left: `${Math.max(20, Math.min(popoverLeft, window.innerWidth - popoverRect.width - 20))}px`,
    };

    setPosition({ highlight: highlightStyles, popover: popoverStyles });
  }, [currentStep]);

  useLayoutEffect(() => {
    setPosition(null);
    const timer = setTimeout(calculatePosition, 150);

    window.addEventListener("resize", calculatePosition);
    return () => {
      window.removeEventListener("resize", calculatePosition);
      clearTimeout(timer);
    };
  }, [calculatePosition]);

  const handleNext = () => {
    if (stepIndex < steps.length - 1) {
      setStepIndex(stepIndex + 1);
    } else {
      onClose();
    }
  };

  const handlePrev = () => {
    if (stepIndex > 0) {
      setStepIndex(stepIndex - 1);
    }
  };

  if (!currentStep) return null;

  const overlayColor = isDark ? "rgba(0, 0, 0, 0.85)" : "rgba(0, 0, 0, 0.65)";
  const glowColor = isDark
    ? "rgba(56, 139, 253, 0.9)"
    : "rgba(78, 140, 255, 0.6)";
  const glowBorderColor = isDark
    ? "rgba(147, 197, 253, 0.6)"
    : "rgba(255, 255, 255, 0.3)";

  return (
    <div className="fixed inset-0 z-[10000]">
      <svg className="absolute inset-0 w-full h-full pointer-events-none">
        <defs>
          <mask id="tutorial-mask">
            <rect x="0" y="0" width="100%" height="100%" fill="white" />
            {highlightRect && (
              <rect
                x={highlightRect.left - 8}
                y={highlightRect.top - 8}
                width={highlightRect.width + 16}
                height={highlightRect.height + 16}
                rx="8"
                fill="black"
              />
            )}
          </mask>
        </defs>
        <rect
          x="0"
          y="0"
          width="100%"
          height="100%"
          fill={overlayColor}
          mask="url(#tutorial-mask)"
          className="pointer-events-auto"
        />
      </svg>

      {position && !position.isError && highlightRect && (
        <div
          className="absolute pointer-events-none transition-all duration-400 ease-out"
          style={{
            ...position.highlight,
            boxShadow: isDark
              ? `0 0 0 6px ${glowBorderColor}, 0 0 32px 12px ${glowColor}, inset 0 0 20px 4px rgba(56, 139, 253, 0.2)`
              : `0 0 0 4px ${glowBorderColor}, 0 0 24px 8px ${glowColor}`,
            borderRadius: "8px",
          }}
        />
      )}

      {position && (
        <div
          className={`absolute rounded-lg p-6 shadow-lg transition-all duration-400 ease-out border ${
            isDark
              ? "bg-gray-800 text-gray-100 border-gray-700"
              : "bg-white text-gray-800 border-gray-200"
          }`}
          style={{
            ...position.popover,
            width: "320px",
            maxWidth: "calc(100vw - 40px)",
            zIndex: 10002,
          }}
        >
          {position.isError ? (
            <>
              <h4
                className={`text-xl font-semibold mb-3 ${
                  isDark ? "text-blue-400" : "text-blue-600"
                }`}
              >
                Element Not Found
              </h4>
              <p className="text-sm leading-relaxed mb-6">
                Could not find the element for this step.
                <br />
                Required selector:{" "}
                <code
                  className={`px-2 py-1 rounded font-mono text-xs ${
                    isDark
                      ? "bg-gray-700 text-red-400"
                      : "bg-gray-100 text-red-600"
                  }`}
                >
                  {currentStep.selector}
                </code>
              </p>
            </>
          ) : (
            <>
              <h4
                className={`text-xl font-semibold mb-3 ${
                  isDark ? "text-blue-400" : "text-blue-600"
                }`}
              >
                {currentStep.title}
              </h4>
              <p className="text-sm leading-relaxed mb-6">{currentStep.text}</p>
            </>
          )}
          <div
            className={`flex justify-between items-center border-t pt-4 mt-4 ${
              isDark ? "border-gray-700" : "border-gray-200"
            }`}
          >
            <span
              className={`text-sm ${
                isDark ? "text-gray-400" : "text-gray-600"
              }`}
            >
              {stepIndex + 1} / {steps.length}
            </span>
            <div className="flex gap-2">
              {stepIndex > 0 && (
                <button
                  onClick={handlePrev}
                  className={`px-4 py-2 rounded-md font-medium transition-colors ${
                    isDark
                      ? "bg-gray-700 text-gray-100 border border-gray-600 hover:bg-gray-600"
                      : "bg-gray-100 text-gray-800 border border-gray-300 hover:bg-gray-200"
                  }`}
                >
                  Previous
                </button>
              )}
              <button
                onClick={handleNext}
                className={`px-4 py-2 rounded-md font-medium transition-colors ${
                  isDark
                    ? "bg-green-600 text-white hover:bg-green-700"
                    : "bg-green-600 text-white hover:bg-green-700"
                }`}
              >
                {stepIndex === steps.length - 1 ? "Finish" : "Next"}
              </button>
            </div>
          </div>
          <button
            className={`absolute top-3 right-3 p-2 rounded-full transition-colors ${
              isDark
                ? "text-gray-400 hover:bg-gray-700 hover:text-gray-100"
                : "bg-white text-gray-500 hover:bg-gray-100 hover:text-gray-800 shadow-sm border border-gray-200"
            }`}
            onClick={onClose}
          >
            <X size={20} />
          </button>
        </div>
      )}
    </div>
  );
};

export default InteractiveTutorial;
