import os
import sys
import pathlib
from unittest.mock import patch, mock_open


if "trafficLight" not in sys.modules:
    sys.path.insert(0, str(pathlib.Path(__file__).resolve().parent.parent / "intersections"))
    import trafficLight as tl


def test_writeNodeFile_creates_expected_file(tmp_path):
    file = tmp_path / "nodes.xml"
    tl.writeNodeFile(str(file))
    content = file.read_text()
    assert "<node id=" in content
    assert "traffic_light" in content


def test_writeEdgeFile_creates_expected_file(tmp_path):
    file = tmp_path / "edges.xml"
    tl.writeEdgeFile(str(file))
    content = file.read_text()
    assert "<edge id=" in content
    assert "in_n2_1" in content


def test_writeConnectionFile_creates_expected_file(tmp_path):
    file = tmp_path / "connections.xml"
    tl.writeConnectionFile(str(file))
    content = file.read_text()
    assert "<connection from=" in content
    assert 'tl="0"' in content


def test_writeTrafficLightLogic_writes_correct_phases(tmp_path):
    file = tmp_path / "tl_logic.xml"
    tl.writeTrafficLightLogic(str(file), greenDuration="30", redDuration="40")
    content = file.read_text()
    assert 'duration="30"' in content
    assert 'duration="40"' in content
    assert "<phase" in content


@patch("subprocess.run")
def test_generateTrips_runs_correct_command(mock_run, tmp_path):
    os.environ["SUMO_HOME"] = "/fake/sumo"
    netFile = str(tmp_path / "net.xml")
    tripFile = str(tmp_path / "trips.xml")

    tl.generateTrips(netFile, tripFile, "high")
    args = mock_run.call_args[0][0]

    assert any("randomTrips.py" in arg for arg in args)
    assert netFile in args
    assert tripFile in args
    assert "--period" in args
    assert "2" in args


@patch("builtins.open", new_callable=mock_open)
@patch("subprocess.run")
@patch("trafficLight.writeNodeFile")
@patch("trafficLight.writeEdgeFile")
@patch("trafficLight.writeConnectionFile")
@patch("trafficLight.writeTrafficLightLogic")
def test_generate_calls_all_write_functions_and_runs_subprocesses(
    mock_writeTL, mock_writeCon, mock_writeEdge, mock_writeNode, mock_subproc_run, mock_open_file, tmp_path
):

    params = {"Green": "30", "Red": "40", "Traffic Density": "medium"}
    tl.generate(params)

    mock_writeNode.assert_called_once()
    mock_writeEdge.assert_called_once()
    mock_writeCon.assert_called_once()
    mock_writeTL.assert_called_once_with("tl_intersection.tll.xml", "30", "40")

    assert mock_subproc_run.call_count == 3

    netconvert_call = mock_subproc_run.call_args_list[0][0][0]
    assert "netconvert" in netconvert_call[0]

    randomtrips_call = mock_subproc_run.call_args_list[1][0][0]
    assert any("randomTrips" in arg for arg in randomtrips_call)

    sumo_call = mock_subproc_run.call_args_list[2][0][0]
    assert "sumo-gui" in sumo_call[0]

    mock_open_file.assert_called_with("tl_intersection.sumocfg", "w")
    handle = mock_open_file()
    written = "".join(call_arg[0][0] for call_arg in handle.write.call_args_list)
    assert "<net-file value=" in written
    assert "<route-files value=" in written
    assert "<additional-files value=" in written


@patch("subprocess.run")
def test_generateTrips_low_density(mock_run, tmp_path):
    os.environ["SUMO_HOME"] = "/fake/sumo"
    netFile = str(tmp_path / "net.xml")
    tripFile = str(tmp_path / "trips.xml")

    tl.generateTrips(netFile, tripFile, "low")
    args = mock_run.call_args[0][0]

    assert any("randomTrips" in arg for arg in args)
    assert "--period" in args
    assert "10" in args


@patch("subprocess.run")
def test_generateTrips_default_density(mock_run, tmp_path):
    os.environ["SUMO_HOME"] = "/fake/sumo"
    netFile = str(tmp_path / "net.xml")
    tripFile = str(tmp_path / "trips.xml")

    tl.generateTrips(netFile, tripFile, "invalid")
    args = mock_run.call_args[0][0]

    assert any("randomTrips" in arg for arg in args)
    assert "--period" in args
    assert "5" in args  # Default fallback
    