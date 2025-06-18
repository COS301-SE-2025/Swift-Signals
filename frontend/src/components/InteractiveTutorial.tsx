import React, { useEffect, useLayoutEffect, useState, useCallback } from 'react';
import '../styles/InteractiveTutorial.css';
import { FaTimes } from 'react-icons/fa';

type TutorialStep = {
    selector: string;
    title: string;
    text: string;
    position?: 'top' | 'bottom' | 'left' | 'right';
};

const tutorialSteps: TutorialStep[] = [
    {
        selector: '.card-grid', 
        title: 'Summary Cards',
        text: 'These cards give you a quick, at-a-glance overview of your key metrics, like total simulations and active intersections.',
        position: 'bottom',
    },
    {
        selector: '.recent-simulations-tab', 
        title: 'Simulations Table',
        text: 'Here you can see a list of all your recent simulations. Click on any row to see more details.',
        position: 'right',
    },
    {
        selector: '.quick-action-button.bg-customIndigo', 
        title: 'Add a New Intersection',
        text: 'Click this button to open the form for creating a new traffic intersection.',
        position: 'bottom',
    },
    {
        selector: '.quick-action-button.bg-customGreen',
        title: 'Run a Simulation',
        text: 'Click this button to open the form for running a traffic simulation.',
        position: 'bottom',
    },
    {
        selector: '.quick-action-button.bg-customPurple',
        title: 'View Map',
        text: 'This will take you to a full-screen map view of all your monitored intersections.',
        position: 'bottom',
    },
    {
        selector: '.graph-card',
        title: 'Traffic Volume Chart',
        text: 'This chart shows the traffic volume over time for your key intersections, helping you identify peak hours and trends.',
        position: 'left',
    },
    {
        selector: '.inter-card',
        title: 'Top Intersections',
        text: 'This card displays the top intersections based on traffic volume, helping you focus on the busiest areas.',
        position: 'left',
    },
];

type Position = {
    highlight: React.CSSProperties;
    popover: React.CSSProperties;
    isError?: boolean;
}

type Props = {
    onClose: () => void;
};

const InteractiveTutorial: React.FC<Props> = ({ onClose }) => {
    const [stepIndex, setStepIndex] = useState(0);
    const [position, setPosition] = useState<Position | null>(null);

    const currentStep = tutorialSteps[stepIndex];

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
        
        // --- THE FIX IS HERE ---
        // First, we get the element's current position on the screen.
        let rect = element.getBoundingClientRect();

        // Then, we check if it's already fully inside the visible area of the window.
        const isElementInViewport = rect.top >= 0 && rect.bottom <= window.innerHeight;

        // We ONLY scroll the page if the element is NOT already fully visible.
        if (!isElementInViewport) {
            element.scrollIntoView({ behavior: 'auto', block: 'center' });
            // After the instant scroll, we MUST re-measure the element's position
            // because its `top` and `left` values will have changed.
            rect = element.getBoundingClientRect();
        }

        // Now we can proceed with a `rect` variable that has the correct, final coordinates.
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
        if (stepIndex < tutorialSteps.length - 1) {
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
                         <span className="tutorial-step-count">{stepIndex + 1} / {tutorialSteps.length}</span>
                         <div className='nav-buttons'>
                             {stepIndex > 0 && <button onClick={handlePrev}>Previous</button>}
                             <button onClick={handleNext}>
                                 {stepIndex === tutorialSteps.length - 1 ? 'Finish' : 'Next'}
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