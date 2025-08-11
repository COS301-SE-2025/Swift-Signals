import React, { useRef, useState, useEffect, useMemo } from "react";
import type { FC } from "react";
import { Canvas, useFrame } from "@react-three/fiber";
import { MapControls, OrthographicCamera } from "@react-three/drei";
import * as THREE from "three";
import { SimulationUI } from "../components/SimulationUI";

// Data Interfaces & Helpers
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
const findNodeById = (nodes: Node[], id: string): Node | null =>
  nodes.find((n) => n.id === id) || null;
const lerp = (start: number, end: number, t: number): number =>
  start * (1 - t) + end * t;
const realisticCarColors = [
  0xeaeaea, 0xb0b0b0, 0x4b8bbe, 0xd44d5c, 0x3e8e7e, 0xf2a34f, 0x6a6a6a,
  0x99b898,
];
const getRandomCarColor = () =>
  realisticCarColors[Math.floor(Math.random() * realisticCarColors.length)];

// 3D Scene Components
const GroundPlane: FC = () => (
  <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, -0.5, 0]} receiveShadow>
    <planeGeometry args={[2000, 2000]} />
    <meshStandardMaterial color={0xaaaaaa} />
  </mesh>
);

const Roads: FC<{ edges: Edge[]; nodes: Node[]; center: THREE.Vector2 }> = ({
  edges,
  nodes,
  center,
}) => {
  return (
    <group>
      {edges.map((edge) => {
        const fromNode = findNodeById(nodes, edge.from);
        const toNode = findNodeById(nodes, edge.to);
        if (!fromNode || !toNode) {
          return null;
        }

        const start = new THREE.Vector2(
          fromNode.x - center.x,
          fromNode.y - center.y
        );
        const end = new THREE.Vector2(toNode.x - center.x, toNode.y - center.y);
        const length = start.distanceTo(end);
        const angle = Math.atan2(end.y - start.y, end.x - start.x);
        const roadWidth = edge.lanes * 10;
        const laneWidth = roadWidth / edge.lanes;
        const position = new THREE.Vector3(
          (start.x + end.x) / 2,
          0,
          (start.y + end.y) / 2
        );

        return (
          <group key={edge.id} position={position} rotation={[0, -angle, 0]}>
            <mesh rotation={[-Math.PI / 2, 0, 0]}>
              <planeGeometry args={[length, roadWidth]} />
              <meshStandardMaterial color={0x282828} />
            </mesh>
            {Array.from({ length: edge.lanes === 1 ? 1 : edge.lanes - 1 }).map(
              (_, laneIndex) => {
                const offset =
                  edge.lanes === 1
                    ? 0
                    : (laneIndex + 1) * laneWidth - roadWidth / 2;
                return (
                  <group key={`divider-${edge.id}-${laneIndex}`}>
                    {Array.from({ length: Math.floor(length / 10) }).map(
                      (_, i) => (
                        <mesh
                          key={i}
                          position={[-length / 2 + i * 10 + 5, 0.1, offset]}
                          rotation={[-Math.PI / 2, 0, 0]}
                        >
                          <planeGeometry args={[4, 0.2]} />
                          <meshStandardMaterial color={0xffffff} />
                        </mesh>
                      )
                    )}
                  </group>
                );
              }
            )}
          </group>
        );
      })}
    </group>
  );
};

const Vehicle: FC<{
  vehicleData: VehicleData;
  simulationTime: number;
  offset: THREE.Vector2;
  color: number;
}> = ({ vehicleData, simulationTime, offset, color }) => {
  const vehicleRef = useRef<THREE.Group>(null);
  useFrame(() => {
    if (!vehicleRef.current || vehicleData.positions.length < 2) return;
    let currentPos: Position | null = null;
    let nextPos: Position | null = null;
    for (let i = 0; i < vehicleData.positions.length - 1; i++) {
      if (
        simulationTime >= vehicleData.positions[i].time &&
        simulationTime < vehicleData.positions[i + 1].time
      ) {
        currentPos = vehicleData.positions[i];
        nextPos = vehicleData.positions[i + 1];
        break;
      }
    }
    if (currentPos && nextPos) {
      vehicleRef.current.visible = true;
      const timeDiff = nextPos.time - currentPos.time;
      const progress =
        timeDiff > 0 ? (simulationTime - currentPos.time) / timeDiff : 0;
      const rawX = lerp(currentPos.x, nextPos.x, progress);
      const rawY = lerp(currentPos.y, nextPos.y, progress);
      const alignedX = rawX - offset.x;
      const alignedY = rawY - offset.y;
      vehicleRef.current.position.set(alignedX, 0, alignedY);

      const dx = nextPos.x - currentPos.x;
      const dy = nextPos.y - currentPos.y;

      if (Math.abs(dx) > 0.01 || Math.abs(dy) > 0.01) {
        vehicleRef.current.lookAt(
          nextPos.x - offset.x,
          0,
          nextPos.y - offset.y
        );
      }
    } else {
      vehicleRef.current.visible = false;
    }
  });
  return (
    <group ref={vehicleRef} castShadow>
      <mesh castShadow position={[0, 0.6, 0]}>
        <boxGeometry args={[2.4, 1.2, 4.8]} />
        <meshStandardMaterial color={color} metalness={0.6} roughness={0.3} />
      </mesh>
      <mesh position={[0, 1.5, -0.5]}>
        <boxGeometry args={[2.2, 0.6, 2.5]} />
        <meshStandardMaterial
          color={0x111111}
          metalness={0.2}
          roughness={0.5}
        />
      </mesh>
    </group>
  );
};

const TrafficLightVisual: FC<{ state: string }> = ({ state }) => {
  const s = state ? state.toLowerCase() : "r";
  const redColor = 0xff0000;
  const yellowColor = 0xffff00;
  const greenColor = 0x00ff00;

  return (
    <group>
      <mesh position={[0, 2, -0.1]} renderOrder={0}>
        <boxGeometry args={[0.8, 2.2, 1]} />
        <meshStandardMaterial color={0x222222} />
      </mesh>
      {s === "r" && (
        <>
          <mesh position={[0, 2.7, 0.5]} renderOrder={1}>
            <sphereGeometry args={[0.3, 32, 32]} />
            <meshBasicMaterial color={redColor} toneMapped={false} />
          </mesh>
          <pointLight position={[0, 2.7, 0.5]} color={redColor} intensity={5} distance={4} decay={2} renderOrder={2} />
        </>
      )}
      {(s === "y" || s === "u") && (
        <>
          <mesh position={[0, 2.0, 0.5]} renderOrder={1}>
            <sphereGeometry args={[0.3, 32, 32]} />
            <meshBasicMaterial color={yellowColor} toneMapped={false} />
          </mesh>
          <pointLight position={[0, 2.0, 0.5]} color={yellowColor} intensity={5} distance={4} decay={2} renderOrder={2} />
        </>
      )}
      {s === "g" && (
        <>
          <mesh position={[0, 1.3, 0.5]} renderOrder={1}>
            <sphereGeometry args={[0.3, 32, 32]} />
            <meshBasicMaterial color={greenColor} toneMapped={false} />
          </mesh>
          <pointLight position={[0, 1.3, 0.5]} color={greenColor} intensity={5} distance={4} decay={2} renderOrder={2} />
        </>
      )}
    </group>
  );
};

const TrafficLightController: FC<{
  lightData: TrafficLightData;
  nodes: Node[];
  edges: Edge[];
  connections: Connection[];
  center: THREE.Vector2;
  simulationTime: number;
  roadDirections: { [key: string]: string };
  onStateUpdate: (states: { [key: string]: string }) => void;
}> = ({ lightData, nodes, edges, connections, center, simulationTime, roadDirections, onStateUpdate }) => {
  const currentStateString = useMemo(() => {
    if (!lightData.states || !Array.isArray(lightData.states) || lightData.states.length === 0) return "";
    let activeState = lightData.states[0].state;
    for (const state of lightData.states) {
      if (simulationTime >= state.time) activeState = state.state; else break;
    }
    return activeState;
  }, [simulationTime, lightData.states]);

  const lightInfoByRoad = useMemo(() => {
    const roadMap = new Map<string, { stateChar: string; edge: Edge }>();
    if (!currentStateString || !connections || !edges) return roadMap;
    const relevantConnections = connections.filter((c) => c.from.indexOf(":") === -1);
    const connectionsByRoad = relevantConnections.reduce((acc, conn) => {
        if (!acc[conn.from]) acc[conn.from] = [];
        acc[conn.from].push(conn);
        return acc;
      }, {} as Record<string, Connection[]>);

    for (const roadId in connectionsByRoad) {
      const roadConnections = connectionsByRoad[roadId];
      const roadStates = roadConnections.map((conn) => (currentStateString[parseInt(conn.tl, 10)] || "r").toLowerCase());
      let finalStateChar = "r";
      if (roadStates.includes("g")) finalStateChar = "g";
      else if (roadStates.includes("y") || roadStates.includes("u")) finalStateChar = "y";
      const edge = edges.find((e) => e.id === roadId);
      if (edge) roadMap.set(roadId, { stateChar: finalStateChar, edge });
    }
    return roadMap;
  }, [connections, currentStateString, edges]);

  useEffect(() => {
    const newStates: { [key: string]: string } = {};
    lightInfoByRoad.forEach(({ stateChar }, roadId) => {
      const direction = roadDirections[roadId];
      if (direction) newStates[direction] = stateChar;
    });
    onStateUpdate(newStates);
  }, [lightInfoByRoad, roadDirections, onStateUpdate]);

  if (lightInfoByRoad.size === 0) return null;

  return (
    <group>
      {Array.from(lightInfoByRoad.entries()).map(([roadId, { stateChar, edge }]) => {
        const fromNode = findNodeById(nodes, edge.from);
        const toNode = findNodeById(nodes, edge.to);
        if (!fromNode || !toNode) return null;

        const start = new THREE.Vector2(fromNode.x - center.x, fromNode.y - center.y);
        const end = new THREE.Vector2(toNode.x - center.x, toNode.y - center.y);
        const dir = end.clone().sub(start).normalize();
        const angle = Math.atan2(dir.y, dir.x);
        const perp = new THREE.Vector2(-dir.y, dir.x);
        const roadLaneWidth = 3.5;
        const offsetDistance = (edge.lanes * roadLaneWidth) / 2 + 1;
        const lightPos = new THREE.Vector3(end.x - dir.x * 5, 0, end.y - dir.y * 5).add(new THREE.Vector3(perp.x * offsetDistance, 0, perp.y * offsetDistance));
        const lightRotationY = -angle + Math.PI / 2;

        return (
          <group key={roadId} position={lightPos} rotation={[0, lightRotationY, 0]}>
            <TrafficLightVisual state={stateChar} />
          </group>
        );
      })}
    </group>
  );
};

const SimulationController: FC<{
  simulationData: SimulationData;
  isPlaying: boolean;
  speed: number;
  offset: THREE.Vector2;
  roadCenter: THREE.Vector2;
  onTimeUpdate: (time: number) => void;
  roadDirections: { [key: string]: string };
  onTrafficLightStateUpdate: (states: { [key: string]: string }) => void;
}> = ({ simulationData, isPlaying, speed, offset, roadCenter, onTimeUpdate, roadDirections, onTrafficLightStateUpdate }) => {
  const [simulationTime, setSimulationTime] = useState(0);
  const vehicleColors = useMemo(() => simulationData.vehicles.map(() => getRandomCarColor()), [simulationData.vehicles]);
  const totalTime = useMemo(() => simulationData.vehicles.reduce((max, v) => Math.max(max, v.positions[v.positions.length - 1]?.time ?? 0), 0), [simulationData]);

  useFrame((_, delta) => {
    if (isPlaying && totalTime > 0) {
      const newTime = simulationTime + delta * speed;
      setSimulationTime(newTime);
      onTimeUpdate(newTime);
    }
  });

  return (
    <>
      <GroundPlane />
      <Roads edges={simulationData.intersection.edges} nodes={simulationData.intersection.nodes} center={roadCenter} />
      {simulationData.vehicles.map((vehicle, index) => (
        <Vehicle key={vehicle.id} vehicleData={vehicle} simulationTime={simulationTime} offset={offset} color={vehicleColors[index]} />
      ))}
      {simulationData.intersection.trafficLights?.map((lightData) => {
        const totalLoopTime = lightData.states && lightData.states.length > 1 ? lightData.states[lightData.states.length - 1].time : 1;
        return (
          <TrafficLightController
            key={lightData.id}
            lightData={lightData}
            nodes={simulationData.intersection.nodes}
            edges={simulationData.intersection.edges}
            connections={simulationData.intersection.connections}
            center={roadCenter}
            simulationTime={simulationTime % totalLoopTime}
            roadDirections={roadDirections}
            onStateUpdate={onTrafficLightStateUpdate}
          />
        );
      })}
    </>
  );
};

interface TrafficSimulationProps {
  dataUrl: string;
  scale?: number;
  isExpanded: boolean;
}

// Hook to detect mobile screen size
const useIsMobile = () => {
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const checkIsMobile = () => {
      setIsMobile(window.innerWidth <= 767);
    };

    checkIsMobile();
    window.addEventListener('resize', checkIsMobile);

    return () => window.removeEventListener('resize', checkIsMobile);
  }, []);

  return isMobile;
};

const TrafficSimulation: FC<TrafficSimulationProps> = ({ dataUrl, scale, isExpanded }) => {
  const [simulationData, setSimulationData] = useState<SimulationData | null>(null);
  const [isPlaying, setIsPlaying] = useState(true);
  const [speed, setSpeed] = useState(5);
  const [simulationTime, setSimulationTime] = useState(0);
  const [restartKey, setRestartKey] = useState(0);
  const [trafficLightStates, setTrafficLightStates] = useState<{ [key: string]: string }>({});
  const isMobile = useIsMobile();

  const roadDirections: { [key: string]: string } = useMemo(() => ({
    in_n2_1: "North", in_n3_1: "South", in_n4_1: "West", in_n5_1: "East",
  }), []);

  const metrics = useMemo(() => {
    if (!simulationData) return { activeVehicles: 0, completedVehicles: 0, avgSpeed: 0, progress: 0, totalVehicles: 0, totalSimTime: 0 };
    const totalVehicles = simulationData.vehicles.length;
    let activeVehicles = 0, completedVehicles = 0, speedSum = 0, speedCount = 0, maxTime = 0;
    simulationData.vehicles.forEach((vehicle) => {
      const positions = vehicle.positions;
      if (positions.length === 0) return;
      const firstTime = positions[0].time;
      const lastTime = positions[positions.length - 1].time;
      if (simulationTime >= firstTime && simulationTime <= lastTime) {
        activeVehicles++;
        let idx = positions.findIndex((p) => p.time > simulationTime);
        if (idx === -1) idx = positions.length - 1; else if (idx > 0) idx = idx - 1;
        speedSum += positions[idx].speed;
        speedCount++;
      } else if (simulationTime > lastTime) {
        completedVehicles++;
      }
      if (lastTime > maxTime) maxTime = lastTime;
    });
    const avgSpeed = speedCount > 0 ? speedSum / speedCount : 0;
    const progress = maxTime > 0 ? Math.min(simulationTime / maxTime, 1) : 0;
    return { activeVehicles, completedVehicles, avgSpeed, progress, totalVehicles, totalSimTime: maxTime };
  }, [simulationData, simulationTime]);

  useEffect(() => {
    fetch(dataUrl)
      .then((res) => (res.ok ? res.json() : Promise.reject(res)))
      .then((data: SimulationData) => {
        if (data.intersection.trafficLights) {
          const directionToSignalIndices: { [key: string]: number[] } = { North: [], South: [], East: [], West: [] };
          const allConnectionIndices = new Set<number>();
          data.intersection.connections.forEach((conn) => {
            if (conn.from.indexOf(":") === -1) {
              const direction = roadDirections[conn.from];
              const signalIndex = parseInt(conn.tl, 10);
              if (direction && !directionToSignalIndices[direction].includes(signalIndex)) {
                directionToSignalIndices[direction].push(signalIndex);
              }
              allConnectionIndices.add(signalIndex);
            }
          });
          const maxSignalIndex = Math.max(...Array.from(allConnectionIndices));
          const stateArrayLength = maxSignalIndex >= 0 ? maxSignalIndex + 1 : 12;
          const newPhases: TrafficLightPhase[] = [];
          const nsGreenDuration = 30;
          const nsGreenState = Array(stateArrayLength).fill("r");
          directionToSignalIndices["North"].forEach((index) => { nsGreenState[index] = "G"; });
          directionToSignalIndices["South"].forEach((index) => { nsGreenState[index] = "G"; });
          newPhases.push({ duration: nsGreenDuration, state: nsGreenState.join("") });
          const nsYellowDuration = 5;
          const nsYellowState = Array(stateArrayLength).fill("r");
          directionToSignalIndices["North"].forEach((index) => { nsYellowState[index] = "y"; });
          directionToSignalIndices["South"].forEach((index) => { nsYellowState[index] = "y"; });
          newPhases.push({ duration: nsYellowDuration, state: nsYellowState.join("") });
          const ewGreenDuration = 30;
          const ewGreenState = Array(stateArrayLength).fill("r");
          directionToSignalIndices["East"].forEach((index) => { ewGreenState[index] = "G"; });
          directionToSignalIndices["West"].forEach((index) => { ewGreenState[index] = "G"; });
          newPhases.push({ duration: ewGreenDuration, state: ewGreenState.join("") });
          const ewYellowDuration = 5;
          const ewYellowState = Array(stateArrayLength).fill("r");
          directionToSignalIndices["East"].forEach((index) => { ewYellowState[index] = "y"; });
          directionToSignalIndices["West"].forEach((index) => { ewYellowState[index] = "y"; });
          newPhases.push({ duration: ewYellowDuration, state: ewYellowState.join("") });
          const processedTrafficLights = data.intersection.trafficLights.map(
            (light) => {
              let time = 0;
              const newStates = newPhases.map((phase) => {
                const state = { time: time, state: phase.state };
                time += phase.duration;
                return state;
              });
              newStates.push({ time: time, state: newPhases[0].state });
              return { ...light, phases: newPhases, states: newStates };
            }
          );
          const newSimData = { ...data, intersection: { ...data.intersection, trafficLights: processedTrafficLights } };
          setSimulationData(newSimData);
        } else {
          setSimulationData(data);
        }
      })
      .catch((error) => console.error(`Error loading simulation data from ${dataUrl}:`, error));
  }, [dataUrl, roadDirections]);

  const { roadCenter, offset } = useMemo(() => {
    if (!simulationData) {
      return {
        roadCenter: new THREE.Vector2(0, 0),
        offset: new THREE.Vector2(0, 0),
      };
    }
    const bounds = {
      road: { minX: Infinity, maxX: -Infinity, minY: Infinity, maxY: -Infinity },
      vehicle: { minX: Infinity, maxX: -Infinity, minY: Infinity, maxY: -Infinity },
    };
    simulationData.intersection.nodes.forEach((n) => {
      bounds.road.minX = Math.min(bounds.road.minX, n.x);
      bounds.road.maxX = Math.max(bounds.road.maxX, n.x);
      bounds.road.minY = Math.min(bounds.road.minY, n.y);
      bounds.road.maxY = Math.max(bounds.road.maxY, n.y);
    });
    simulationData.vehicles.forEach((v) =>
      v.positions.forEach((p) => {
        bounds.vehicle.minX = Math.min(bounds.vehicle.minX, p.x);
        bounds.vehicle.maxX = Math.max(bounds.vehicle.maxX, p.x);
        bounds.vehicle.minY = Math.min(bounds.vehicle.minY, p.y);
        bounds.vehicle.maxY = Math.max(bounds.vehicle.maxY, p.y);
      })
    );
    const roadCenter = new THREE.Vector2(
      (bounds.road.minX + bounds.road.maxX) / 2,
      (bounds.road.minY + bounds.road.maxY) / 2
    );
    const vehicleCenter = new THREE.Vector2(
      (bounds.vehicle.minX + bounds.vehicle.maxX) / 2,
      (bounds.vehicle.minY + bounds.vehicle.maxY) / 2
    );
    return { roadCenter, offset: vehicleCenter };
  }, [simulationData]);

  const handleRestart = () => {
    setSimulationTime(0);
    setRestartKey((prevKey) => prevKey + 1);
  };

  if (!simulationData) {
    return (
      <div style={{ height: "100vh", display: "grid", placeContent: "center", backgroundColor: "#3d3d3d", color: "white", }}>
        Loading Simulation from {dataUrl}...
      </div>
    );
  }

  const canvasContainerWidth = isExpanded ? '100%' : (isMobile ? '100%' : '50%');
  const uiScale = isMobile ? 0.4 : (scale || 1);

  return (
    <div className="traffic-simulation-root" style={{ position: "relative", height: "100vh", backgroundColor: "#3d3d3d", border: "none", boxShadow: "none", }}>
      <div style={{ position: "absolute", top: 0, left: 0, width: canvasContainerWidth, height: "100%", zIndex: 0, transition: 'width 0.5s ease-in-out',}}>
        <Canvas shadows>
          <ambientLight intensity={0.6} />
          <directionalLight castShadow position={[100, 100, 50]} intensity={1.5} shadow-mapSize-width={2048} shadow-mapSize-height={2048} />
          <MapControls 
            enablePan={false} 
            enableRotate={false}
            enableZoom={!isMobile}
            enableDamping={false}
            touches={{ 
              ONE: isMobile ? undefined : THREE.TOUCH.ROTATE,
              TWO: isMobile ? undefined : THREE.TOUCH.DOLLY_PAN 
            }}
          />
          
          <OrthographicCamera
            makeDefault
            position={[roadCenter.x, 100, roadCenter.y]}
            zoom={4}
          />
          
          <SimulationController
            key={restartKey}
            simulationData={simulationData}
            isPlaying={isPlaying}
            speed={speed}
            roadCenter={roadCenter}
            offset={offset}
            onTimeUpdate={setSimulationTime}
            roadDirections={roadDirections}
            onTrafficLightStateUpdate={setTrafficLightStates}
          />
        </Canvas>
      </div>
      <SimulationUI
        time={simulationTime}
        vehicleCount={simulationData.vehicles.length}
        isPlaying={isPlaying}
        speed={speed}
        onPlayPause={() => setIsPlaying(!isPlaying)}
        onRestart={handleRestart}
        onSpeedChange={setSpeed}
        trafficLightStates={trafficLightStates}
        activeVehicles={metrics.activeVehicles}
        completedVehicles={metrics.completedVehicles}
        avgSpeed={metrics.avgSpeed}
        progress={metrics.progress}
        totalSimTime={metrics.totalSimTime}
        scale={uiScale}
      />
    </div>
  );
};

export default TrafficSimulation;