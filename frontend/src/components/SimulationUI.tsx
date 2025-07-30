import React from 'react';
import type { FC } from 'react';

interface SimulationUIProps {
  time: number;
  vehicleCount: number;
  isPlaying: boolean;
  speed: number;
  onPlayPause: () => void;
  onRestart: () => void;
  onSpeedChange: (newSpeed: number) => void;
  trafficLightStates: { [key: string]: string };
  activeVehicles: number;
  completedVehicles: number;
  avgSpeed: number;
  progress: number;
  totalSimTime: number;
  scale?: number;
}

const styles: { [key: string]: React.CSSProperties } = {
  uiPanel: {
    position: 'absolute',
    top: '24px',
    left: '24px',
    background: 'rgba(28, 32, 38, 0.98)',
    color: '#f3f6fa',
    padding: '28px 28px 20px 28px',
    borderRadius: '10px',
    fontFamily: '"Inter", "Segoe UI", Arial, sans-serif',
    width: '320px',
    boxShadow: '0 4px 24px 0 rgba(0,0,0,0.18)',
    border: '1px solid #23272f',
    backdropFilter: 'blur(2px)',
    WebkitBackdropFilter: 'blur(2px)',
    transition: 'transform 0.2s ease-in-out',
    transformOrigin: 'top left',
  },
  title: {
    margin: '0 0 18px 0',
    fontSize: '1.25em',
    borderBottom: '1px solid #23272f',
    paddingBottom: '12px',
    letterSpacing: '0.5px',
    fontWeight: 600,
  },
  dataRow: {
    display: 'flex',
    justifyContent: 'space-between',
    marginBottom: '10px',
    fontSize: '1em',
    alignItems: 'center',
  },
  label: {
    color: '#b0b8c1',
    fontWeight: 400,
    letterSpacing: '0.1px',
  },
  controls: {
    display: 'flex',
    gap: '10px',
    marginTop: '18px',
  },
  button: {
    flexGrow: 1,
    padding: '10px 0',
    cursor: 'pointer',
    border: 'none',
    borderRadius: '5px',
    background: '#23272f',
    color: '#f3f6fa',
    fontWeight: 500,
    fontSize: '1em',
    boxShadow: 'none',
    borderBottom: '2px solid #353a45',
    transition: 'background 0.15s, color 0.15s',
  },
  buttonHover: {
    background: '#353a45',
    color: '#fff',
  },
  sliderContainer: {
    marginTop: '16px',
  },
  lightIndicator: {
    width: '12px',
    height: '12px',
    borderRadius: '50%',
    display: 'inline-block',
    marginLeft: 'auto',
    boxShadow: '0 0 4px currentColor',
    border: '1px solid #353a45',
  },
  lightsContainer: {
    borderTop: '1px solid #23272f',
    marginTop: '14px',
    paddingTop: '12px',
  },
  progressBarContainer: {
    margin: '16px 0 8px 0',
    height: '6px',
    background: '#23272f',
    borderRadius: '3px',
    overflow: 'hidden',
    boxShadow: 'none',
  },
  progressBar: {
    height: '100%',
    background: 'linear-gradient(90deg, #4e8cff 0%, #7ed6df 100%)',
    transition: 'width 0.3s',
  },
  sectionHeader: {
    color: '#f3f6fa',
    fontWeight: 500,
    fontSize: '1.05em',
    margin: '18px 0 8px 0',
    letterSpacing: '0.2px',
    borderBottom: '1px solid #23272f',
    paddingBottom: '4px',
    display: 'block',
  },
};

const lightColorMap: { [key: string]: string } = {
  g: '#28a745',
  y: '#ffc107',
  r: '#dc3545',
};

export const SimulationUI: FC<SimulationUIProps> = ({
  time,
  vehicleCount,
  isPlaying,
  speed,
  onPlayPause,
  onRestart,
  onSpeedChange,
  trafficLightStates,
  activeVehicles,
  completedVehicles,
  avgSpeed,
  progress,
  totalSimTime,
  scale = 1,
}) => {
  const [hoveredBtn, setHoveredBtn] = React.useState<string | null>(null);
  const avgSpeedKmh = avgSpeed * 3.6;

  const panelStyle = {
    ...styles.uiPanel,
    transform: `scale(${scale})`,
  };

  return (
    <div style={panelStyle}>
      <h3 style={styles.title}>Simulation</h3>
      <div style={styles.progressBarContainer}>
        <div style={{ ...styles.progressBar, width: `${Math.round(progress * 100)}%` }} />
      </div>
      <div style={{ ...styles.dataRow, marginBottom: 0 }}>
        <span style={styles.label}>Progress</span>
        <span>{Math.round(progress * 100)}%</span>
      </div>
      <div style={styles.dataRow}>
        <span style={styles.label}>Time</span>
        <span>{time.toFixed(1)} / {totalSimTime.toFixed(1)} s</span>
      </div>
      <div style={styles.sectionHeader}>Vehicles</div>
      <div style={styles.dataRow}>
        <span style={styles.label}>Total</span>
        <span>{vehicleCount}</span>
      </div>
      <div style={styles.dataRow}>
        <span style={styles.label}>Active</span>
        <span>{activeVehicles}</span>
      </div>
      <div style={styles.dataRow}>
        <span style={styles.label}>Completed</span>
        <span>{completedVehicles}</span>
      </div>
      <div style={styles.dataRow}>
        <span style={styles.label}>Avg Speed</span>
        <span>{avgSpeedKmh.toFixed(1)} km/h</span>
      </div>
      <div style={styles.sectionHeader}>Traffic Lights</div>
      <div style={styles.lightsContainer}>
        {Object.entries(trafficLightStates).sort().map(([direction, state]) => (
          <div key={direction} style={styles.dataRow}>
            <span style={styles.label}>{direction}</span>
            <span style={{ ...styles.lightIndicator, backgroundColor: lightColorMap[state] || '#bbb', color: lightColorMap[state] || '#bbb' }}></span>
          </div>
        ))}
      </div>
      <div style={styles.sectionHeader}>Controls</div>
      <div style={styles.controls}>
        <button
          style={{ ...styles.button, ...(hoveredBtn === 'playpause' ? styles.buttonHover : {}), }}
          onClick={onPlayPause}
          onMouseEnter={() => setHoveredBtn('playpause')}
          onMouseLeave={() => setHoveredBtn(null)}
        >
          {isPlaying ? 'Pause' : 'Play'}
        </button>
        <button
          style={{ ...styles.button, ...(hoveredBtn === 'restart' ? styles.buttonHover : {}), }}
          onClick={onRestart}
          onMouseEnter={() => setHoveredBtn('restart')}
          onMouseLeave={() => setHoveredBtn(null)}
        >
          Restart
        </button>
      </div>
      <div style={styles.sliderContainer}>
        <div style={styles.dataRow}>
          <span style={styles.label}>Speed</span>
          <span>{speed}x</span>
        </div>
        <input
          type="range"
          min="1"
          max="20"
          step="1"
          value={speed}
          onChange={(e) => onSpeedChange(Number(e.target.value))}
          style={{ width: '100%' }}
        />
      </div>
    </div>
  );
};