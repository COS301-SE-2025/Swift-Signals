import os
import sys
import pathlib


if "stopStreet" not in sys.modules:
    sys.path.insert(0, str(pathlib.Path(__file__).resolve().parent.parent / "intersections"))
    import stopStreet

from unittest.mock import patch, mock_open


def test_write_node_file(tmp_path):
    file = tmp_path / "test_nodes.xml"
    stopStreet.writeNodeFile(str(file))
    content = file.read_text()
    assert "<node id=" in content
    assert "n1" in content


def test_write_edge_file(tmp_path):
    file = tmp_path / "test_edges.xml"
    stopStreet.writeEdgeFile(str(file))
    content = file.read_text()
    assert "<edge id=" in content
    assert "in_n2_1" in content


def test_write_connection_file(tmp_path):
    file = tmp_path / "test_connections.xml"
    stopStreet.writeConnectionFile(str(file))
    content = file.read_text()
    assert "<connection from=" in content
    assert "in_n2_1" in content


def test_write_stop_logic(tmp_path):
    file = tmp_path / "test_stop_logic.xml"
    stopStreet.writeStopLogic(str(file))
    content = file.read_text()
    assert "<additional>" in content
    assert "priority0" in content


@patch("stopStreet.subprocess.run")
def test_generate_trips_low_density(mock_run, tmp_path):
    os.environ["SUMO_HOME"] = "/fake/sumo"

    netFile = str(tmp_path / "net.xml")
    tripFile = str(tmp_path / "trips.xml")
    stopStreet.generateTrips(netFile, tripFile, "low")

    args = mock_run.call_args[0][0]
    assert any("randomTrips.py" in arg for arg in args)
    assert "-n" in args
    assert netFile in args
    assert "-o" in args
    assert tripFile in args
    assert "--period" in args
    assert "10" in args


@patch("stopStreet.subprocess.run")
def test_generate_trips_default_density(mock_run, tmp_path):
    os.environ["SUMO_HOME"] = "/fake/sumo"
    stopStreet.generateTrips("net.xml", "trips.xml", "unknown")

    args = mock_run.call_args[0][0]
    assert "--period" in args
    assert "5" in args


@patch("stopStreet.subprocess.run")
@patch("stopStreet.writeNodeFile")
@patch("stopStreet.writeEdgeFile")
@patch("stopStreet.writeConnectionFile")
@patch("stopStreet.writeStopLogic")
def test_generate_function_calls(writeStopLogic, writeConnectionFile, writeEdgeFile, writeNodeFile, mock_run, tmp_path):
    params = {"Traffic Density": "medium"}

    with patch("builtins.open", mock_open()):
        stopStreet.generate(params)

    writeNodeFile.assert_called_once()
    writeEdgeFile.assert_called_once()
    writeConnectionFile.assert_called_once()
    writeStopLogic.assert_called_once()

    assert mock_run.call_count == 3

    first_call_args = mock_run.call_args_list[0][0][0]
    assert "netconvert" in first_call_args

    second_call_args = mock_run.call_args_list[1][0][0]
    assert any("randomTrips.py" in arg for arg in second_call_args)

    third_call_args = mock_run.call_args_list[2][0][0]
    assert "sumo-gui" in third_call_args


@patch("stopStreet.subprocess.run")
def test_generate_trips_high_density(mock_run, tmp_path):
    os.environ["SUMO_HOME"] = "/fake/sumo"

    netFile = str(tmp_path / "net.xml")
    tripFile = str(tmp_path / "trips.xml")
    stopStreet.generateTrips(netFile, tripFile, "high")

    args = mock_run.call_args[0][0]
    assert any("randomTrips.py" in arg for arg in args)
    assert "-n" in args
    assert netFile in args
    assert "-o" in args
    assert tripFile in args
    assert "--period" in args
    assert "2" in args
