import React, { useEffect, useLayoutEffect, useState, useCallback } from 'react';
import '../styles/InteractiveTutorial.css';
import { FaTimes } from 'react-icons/fa';

// --- TYPE IS UPDATED ---
// Added optional 'action' function and 'center' position
export type TutorialStep = {
    selector: string;
    title: string;
    text: string;
    position?: 'top' | 'bottom' | 'left' | 'right' | 'center';
    action?: () => void;
};

type Position = {
    highlight: React.CSSProperties;
    popover: React.CSSProperties;
    isError?: boolean;
}

type Props = {
    steps: TutorialStep[];
    onClose: () => void;
};

const InteractiveTutorial: React.FC<Props> = ({ steps, onClose }) => {
    const [stepIndex, setStepIndex] = useState(0);
    const [position, setPosition] = useState<Position | null>(null);

    const currentStep = steps[stepIndex];

    // --- NEW: useEffect to handle actions ---
    // This hook runs when the step changes. If the new step has an action, it executes it.
    useEffect(() => {
        if (currentStep && typeof currentStep.action === 'function') {
            currentStep.action();
        }
    }, [currentStep]);

    const calculatePosition = useCallback(() => {
        if (!currentStep) return;

        // --- UPDATED: Handle 'center' position for action steps ---
        if (currentStep.position === 'center') {
            setPosition({
                highlight: { display: 'none' }, // No highlight for centered steps
                popover: {
                    top: '50%',
                    left: '50%',
                    transform: 'translate(-50%, -50%)'
                }
            });
            return;
        }

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