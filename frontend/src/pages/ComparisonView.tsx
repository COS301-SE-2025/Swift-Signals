import React, { useState } from 'react';
import TrafficSimulation from './TrafficSimulation';

const ComparisonView: React.FC = () => {
  const originalDataUrl = "/simulation_output (1).json";
  const optimizedDataUrl = "/optimized_output.json";
  const [expanded, setExpanded] = useState<'none' | 'left' | 'right'>('none');

  const containerStyle: React.CSSProperties = {
    display: 'flex',
    flexDirection: 'row',
    width: '100vw',
    height: '100vh',
    backgroundColor: '#1e1e1e'
  };

  const viewStyle: React.CSSProperties = {
    position: 'relative',
    height: '100%',
    overflow: 'hidden',
    transition: 'flex-basis 0.5s ease-in-out',
  };

  const labelStyle: React.CSSProperties = {
    position: 'absolute',
    bottom: '20px',
    left: '50%',
    transform: 'translateX(-50%)',
    zIndex: 1000,
    backgroundColor: 'rgba(0,0,0,0.75)',
    color: 'white',
    padding: '8px 16px',
    borderRadius: '8px',
    fontSize: '1em',
    fontWeight: '600',
    pointerEvents: 'none',
  };
  
  const dividerStyle: React.CSSProperties = {
    flexShrink: 0,
    width: '2px',
    backgroundColor: '#333',
    transition: 'width 0.5s ease-in-out',
  };

  const modernButtonStyle: React.CSSProperties = {
    position: 'absolute',
    top: '24px',
    zIndex: 10,
    background: 'linear-gradient(135deg, rgba(255,255,255,0.15) 0%, rgba(255,255,255,0.08) 100%)',
    backdropFilter: 'blur(16px)',
    border: '1px solid rgba(255,255,255,0.18)',
    borderRadius: '16px',
    padding: '12px 20px',
    cursor: 'pointer',
    fontWeight: '600',
    fontSize: '14px',
    color: 'rgba(255,255,255,0.95)',
    transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    boxShadow: '0 8px 32px rgba(0,0,0,0.3)',
    minWidth: '140px',
    justifyContent: 'center',
  };

  const buttonHoverStyle: React.CSSProperties = {
    background: 'linear-gradient(135deg, rgba(255,255,255,0.25) 0%, rgba(255,255,255,0.15) 100%)',
    transform: 'translateY(-2px)',
    boxShadow: '0 12px 40px rgba(0,0,0,0.4)',
    border: '1px solid rgba(255,255,255,0.3)',
  };

  const iconStyle: React.CSSProperties = {
    fontSize: '16px',
    transition: 'transform 0.3s ease',
  };

  const toggleLeft = () => setExpanded(prev => prev === 'left' ? 'none' : 'left');
  const toggleRight = () => setExpanded(prev => prev === 'right' ? 'none' : 'right');

  const getDynamicStyles = (side: 'left' | 'right') => {
    const isExpanded = expanded === side;
    const isCollapsed = (side === 'left' && expanded === 'right') || (side === 'right' && expanded === 'left');

    let flexBasis = '50%';
    if (isExpanded) flexBasis = '100%';
    if (isCollapsed) flexBasis = '0%';

    return { ...viewStyle, flex: `1 1 ${flexBasis}` };
  };

  const getButtonContent = (side: 'left' | 'right') => {
    const isExpanded = expanded === side;
    const isOtherExpanded = expanded !== 'none' && expanded !== side;
    
    if (isExpanded) {
      return {
        icon: side === 'left' ? '→' : '←',
        text: 'Expand',
        tooltip: 'Expand to show both views'
      };
    } else if (isOtherExpanded) {
      return {
        icon: side === 'left' ? '←' : '→',
        text: 'Collapse',
        tooltip: `Collapse ${side === 'left' ? 'original' : 'optimized'} view (will collapse other side)`
      };
    } else {
      return {
        icon: side === 'left' ? '←' : '→',
        text: 'Collapse',
        tooltip: `Collapse ${side === 'left' ? 'original' : 'optimized'} view (will collapse other side)`
      };
    }
  };

  const ModernButton = ({ side, onClick, position }: { 
    side: 'left' | 'right', 
    onClick: () => void,
    position: 'left' | 'right'
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
          [position]: '24px',
        }}
        title={content.tooltip}
      >
        <span style={{
          ...iconStyle,
          transform: isHovered ? 'scale(1.1)' : 'scale(1)'
        }}>
          {content.icon}
        </span>
        <span>{content.text}</span>
      </button>
    );
  };

  return (
    <div style={containerStyle}>
      {/* Left side: Original Simulation */}
      <div style={getDynamicStyles('left')}>
        <TrafficSimulation
          dataUrl={originalDataUrl}
          scale={0.65}
          isExpanded={expanded === 'left'}
        />
        <div style={labelStyle}>Original Simulation</div>
        <ModernButton 
          side="left" 
          onClick={toggleLeft}
          position="right"
        />
      </div>

      <div style={{ ...dividerStyle, width: expanded === 'none' ? '2px' : '0px' }} />

      {/* Right side: Optimized Simulation */}
      <div style={getDynamicStyles('right')}>
        <TrafficSimulation
          dataUrl={optimizedDataUrl}
          scale={0.65}
          isExpanded={expanded === 'right'}
        />
        <div style={labelStyle}>Optimized Simulation</div>
        <ModernButton 
          side="right" 
          onClick={toggleRight}
          position="left"
        />
      </div>
    </div>
  );
};

export default ComparisonView;