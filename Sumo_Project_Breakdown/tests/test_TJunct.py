import builtins
import os
import subprocess
import pytest
from unittest import mock
import sys
import pathlib
sys.path.insert(0, str(pathlib.Path(__file__).resolve().parent.parent / "intersections"))
import tJunction

@pytest.fixture(autouse=True)
def no_file_write(monkeypatch):
    mock_open = mock.mock_open()
    monkeypatch.setattr("builtins.open", mock_open)
    yield mock_open

@pytest.fixture(autouse=True)
def no_subprocess(monkeypatch):
    mock_run = mock.Mock()
    monkeypatch.setattr(subprocess, "run", mock_run)
    return mock_run

@pytest.fixture(autouse=True)
def no_makedirs(monkeypatch):
    mock_makedirs = mock.Mock()
    monkeypatch.setattr(os, "makedirs", mock_makedirs)
    return mock_makedirs

def test_generate_calls(monkeypatch, no_file_write, no_subprocess):
    params = {"Traffic Density": "medium"}
    
    called_with = {}
    def fake_generateTrips(netFile, tripFile, density):
        called_with["netFile"] = netFile
        called_with["tripFile"] = tripFile
        called_with["density"] = density
    monkeypatch.setattr(tJunction, "generateTrips", fake_generateTrips)

    tJunction.generate(params)

    open_calls = [call.args[0] for call in no_file_write.call_args_list]
    assert "tjInt.nod.xml" in open_calls
    assert "tjInt.edg.xml" in open_calls
    assert "tjInt.con.xml" in open_calls
    assert "t_junction.sumocfg" in open_calls

    assert called_with["netFile"] == "t_junction.net.xml"
    assert called_with["tripFile"] == "t_junction.rou.xml"
    assert called_with["density"] == "medium"

    cmds = [call.args[0] for call in no_subprocess.call_args_list]
    assert any("netconvert" in cmd[0] for cmd in cmds)
    assert any("sumo-gui" in cmd[0] for cmd in cmds)

def test_write_node_file(no_file_write):
    filename = "test_nodes.xml"
    tJunction.writeNodeFile(filename)
    no_file_write().write.assert_called_with("""<nodes>
    <node id="J1" x="0" y="0" type="priority"/> <!-- T-junction -->
    <node id="nNorth" x="0" y="100" type="priority"/>
    <node id="nEast" x="100" y="0" type="priority"/>
    <node id="nWest" x="-100" y="0" type="priority"/>
</nodes>""")

def test_write_edge_file(no_file_write):
    filename = "test_edges.xml"
    tJunction.writeEdgeFile(filename)
    no_file_write().write.assert_called_with("""<edges>
    <edge id="in_nNorth_J1" from="nNorth" to="J1" priority="1" numLanes="1" speed="13.9"/>
    <edge id="in_nEast_J1" from="nEast" to="J1" priority="3" numLanes="1" speed="13.9"/>
    <edge id="in_nWest_J1" from="nWest" to="J1" priority="3" numLanes="1" speed="13.9"/>

    <edge id="out_J1_nNorth" from="J1" to="nNorth" priority="1" numLanes="1" speed="13.9"/>
    <edge id="out_J1_nEast" from="J1" to="nEast" priority="3" numLanes="1" speed="13.9"/>
    <edge id="out_J1_nWest" from="J1" to="nWest" priority="3" numLanes="1" speed="13.9"/>
</edges>""")

def test_write_connection_file(no_file_write):
    filename = "test_connections.xml"
    tJunction.writeConnectionFile(filename)
    no_file_write().write.assert_called_with("""<connections>
    <connection from="in_nNorth_J1" to="out_J1_nEast" fromLane="0" toLane="0"/>
    <connection from="in_nNorth_J1" to="out_J1_nWest" fromLane="0" toLane="0"/>

    <connection from="in_nEast_J1" to="out_J1_nNorth" fromLane="0" toLane="0"/>
    <connection from="in_nEast_J1" to="out_J1_nWest" fromLane="0" toLane="0"/>

    <connection from="in_nWest_J1" to="out_J1_nNorth" fromLane="0" toLane="0"/>
    <connection from="in_nWest_J1" to="out_J1_nEast" fromLane="0" toLane="0"/>
</connections>""")

@pytest.mark.parametrize("density, expected_period", [
    ("low", "10"),
    ("medium", "5"),
    ("high", "2"),
    ("unknown", "5")  
])
def test_generateTrips_runs_commands(monkeypatch, no_subprocess, no_makedirs, density, expected_period):
    monkeypatch.setenv("SUMO_HOME", "/fake/sumo")

    called_cmds = []
    def fake_run(cmd, check):
        called_cmds.append(cmd)

    no_subprocess.side_effect = fake_run

    netFile = "net.xml"
    tripFile = "out.rou.xml"
    tJunction.generateTrips(netFile, tripFile, density)

    trip_dir = os.path.dirname(tripFile)
    if trip_dir:
        no_makedirs.assert_called_with(trip_dir, exist_ok=True)

    found_period = any(f"--period={expected_period}" in arg or expected_period in arg for cmd in called_cmds for arg in cmd)
    assert found_period

    found_netfile = any(f"-n" in cmd and netFile in cmd for cmd in called_cmds)
    found_tripfile = any(f"-o" in cmd and tripFile in cmd for cmd in called_cmds)
    assert any(netFile in arg for cmd in called_cmds for arg in cmd)
    assert any(tripFile in arg for cmd in called_cmds for arg in cmd)

def test_generateTrips_creates_dir_when_needed(mocker):
    # Setup a tripFile with a directory in the path
    tripFile = "some_dir/trips.xml"
    netFile = "net.xml"
    density = "medium"

    makedirs_mock = mocker.patch("os.makedirs")
    run_mock = mocker.patch("subprocess.run")

    tJunction.generateTrips(netFile, tripFile, density)

    # Should call os.makedirs since dirname is 'some_dir'
    makedirs_mock.assert_called_once_with("some_dir", exist_ok=True)

def test_generateTrips_no_dir_creation_for_file_in_cwd(mocker):
    # Setup a tripFile with just a filename (no directory)
    tripFile = "trips.xml"
    netFile = "net.xml"
    density = "medium"

    makedirs_mock = mocker.patch("os.makedirs")
    run_mock = mocker.patch("subprocess.run")

    tJunction.generateTrips(netFile, tripFile, density)

    # Should NOT call os.makedirs since dirname is empty string
    makedirs_mock.assert_not_called()
