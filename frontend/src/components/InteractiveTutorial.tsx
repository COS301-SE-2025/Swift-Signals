import React, { useEffect, useLayoutEffect, useState, useCallback } from 'react';
import '../styles/InteractiveTutorial.css';
import { FaTimes } from 'react-icons/fa';

// The definition of a single step remains the same
export type TutorialStep = {
    selector: string;
    title: string;
    text: string;
    position?: 'top' | 'bottom' | 'left' | 'right';
};

type Position = {
    highlight: React.CSSProperties;
    popover: React.CSSProperties;
    isError?: boolean;
}

// --- PROPS ARE UPDATED ---
// It now accepts an array of steps
type Props = {
    steps: TutorialStep[];
    onClose: () => void;
};

// --- The hardcoded 'tutorialSteps' array has been REMOVED from this file ---

const InteractiveTutorial: React.FC<Props> = ({ steps, onClose }) => { // Accept 'steps' from props
    const [stepIndex, setStepIndex] = useState(0);
    const [position, setPosition] = useState<Position | null>(null);

    const currentStep = steps[stepIndex]; // Use the 'steps' prop

    const calculatePosition = useCallback(() => {
        if (!currentStep) return;

        const element = document.querySelector(currentStep.selector) as HTMLElement;

        if (!element) {
            setPosition({
                isError: true,
                highlight: {},
                popover: {
                    top: '50%',
                    left: '50%',
                    transform: 'translate(-50%, -50%)'
                }
            });
            return;
        }
        
        let rect = element.getBoundingClientRect();
        const isElementInViewport = rect.top >= 0 && rect.bottom <= window.innerHeight;

        if (!isElementInViewport) {
            element.scrollIntoView({ behavior: 'auto', block: 'center' });
            rect = element.getBoundingClientRect();
        }

        const popoverRect = { width: 320, height: 150 };

        const highlightStyles: React.CSSProperties = {
            width: `${rect.width + 16}px`,
            height: `${rect.height + 16}px`,
            top: `${rect.top - 8}px`,
            left: `${rect.left - 8}px`,
        };

        let popoverTop = 0;
        let popoverLeft = 0;

        switch (currentStep.position) {
            case 'top':
                popoverTop = rect.top - popoverRect.height - 20;
                popoverLeft = rect.left + (rect.width / 2) - (popoverRect.width / 2);
                break;
            case 'left':
                popoverTop = rect.top + (rect.height / 2) - (popoverRect.height / 2);
                popoverLeft = rect.left - popoverRect.width - 20;
                break;
            case 'right':
                popoverTop = rect.top + (rect.height / 2) - (popoverRect.height / 2);
                popoverLeft = rect.right + 20;
                break;
            default: // bottom
                popoverTop = rect.bottom + 20;
                popoverLeft = rect.left + (rect.width / 2) - (popoverRect.width / 2);
                break;
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

        window.addEventListener('resize', calculatePosition);
        return () => {
            window.removeEventListener('resize', calculatePosition);
            clearTimeout(timer);
        };
    }, [calculatePosition]);

    const handleNext = () => {
        if (stepIndex < steps.length - 1) { // Use 'steps' prop
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

    return (
        <div className="tutorial-overlay">
            {position && !position.isError && (
                 <div className="tutorial-highlight" style={position.highlight}></div>
            )}
            
            {position && (
                 <div className="tutorial-popover" style={position.popover}>
                     {position.isError ? (
                         <>
                             <h4>Element Not Found</h4>
                             <p>
                                 Could not find the element for this step.
                                 <br/>
                                 Required selector: <code>{currentStep.selector}</code>
                             </p>
                         </>
                     ) : (
                         <>
                             <h4>{currentStep.title}</h4>
                             <p>{currentStep.text}</p>
                         </>
                     )}
                     <div className="tutorial-navigation">
                         {/* Use 'steps' prop for the count */}
                         <span className="tutorial-step-count">{stepIndex + 1} / {steps.length}</span>
                         <div className='nav-buttons'>
                             {stepIndex > 0 && <button onClick={handlePrev}>Previous</button>}
                             <button onClick={handleNext}>
                                 {stepIndex === steps.length - 1 ? 'Finish' : 'Next'}
                             </button>
                         </div>
                     </div>
                     <button className="tutorial-close-button" onClick={onClose}><FaTimes/></button>
                 </div>
            )}
        </div>
    );
};

export default InteractiveTutorial;