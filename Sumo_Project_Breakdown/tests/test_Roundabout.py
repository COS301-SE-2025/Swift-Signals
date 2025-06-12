import os
from unittest.mock import patch, mock_open
import sys
import pathlib

if "circle" not in sys.modules:
    sys.path.insert(0, str(pathlib.Path(__file__).resolve().parent.parent / "intersections"))
    import circle


@patch("circle.subprocess.run")
@patch("builtins.open", new_callable=mock_open)
def test_generate_calls_and_files(mock_file, mock_run):
    os.environ["SUMO_HOME"] = "/fake/sumo"

    params = {"Traffic Density": "medium"}

    circle.generate(params)

    expected_files = ["tlInt.nod.xml", "tlInt.edg.xml", "tlInt.con.xml"]
    actual_files = [call_args[0][0] for call_args in mock_file.call_args_list]
    for f in expected_files:
        assert f in actual_files

    calls = [call_args[0][0] for call_args in mock_run.call_args_list]
    netconvert_call = calls[0]
    assert netconvert_call[0] == "netconvert"
    assert any(arg.startswith("--node-files=tlInt.nod.xml") for arg in netconvert_call)
    assert any(arg.startswith("--edge-files=tlInt.edg.xml") for arg in netconvert_call)
    assert any(arg.startswith("--connection-files=tlInt.con.xml") for arg in netconvert_call)
    assert "-o" in netconvert_call and "tl_intersection.net.xml" in netconvert_call

    trips_call = calls[1]
    assert trips_call[0] == "python"
    assert "randomTrips.py" in trips_call[1]
    assert "--period" in trips_call and "5" in trips_call

    route_call = calls[2]
    assert route_call[0] == "duarouter"
    assert "-n" in route_call and "tl_intersection.net.xml" in route_call
    assert "-t" in route_call and "tl_intersection.trips.xml" in route_call
    assert "-o" in route_call and "tl_intersection.rou.xml" in route_call

    gui_call = calls[3]
    assert gui_call[0] == "sumo-gui"
    assert "-c" in gui_call and "tl_intersection.sumocfg" in gui_call

    config_handle = mock_file()
    config_handle.write.assert_any_call(
        '<configuration>\n        <input>\n            <net-file value="tl_intersection.net.xml"/>\n            <route-files value="tl_intersection.rou.xml"/>\n        </input>\n        <time>\n            <begin value="0"/>\n            <end value="1000"/>\n        </time>\n    </configuration>'
    )


@patch("circle.subprocess.run")
def test_generateTrips_density_values(mock_run):
    os.environ["SUMO_HOME"] = "/fake/sumo"

    netFile = "net.xml"
    routeFile = "route.rou.xml"

    densities = {
        "low": "10",
        "medium": "5",
        "high": "2",
        "unknown": "5",
        None: "5",
    }

    for density, expected_period in densities.items():
        circle.generateTrips(netFile, routeFile, density)
        args = mock_run.call_args_list[-2][0][0]
        assert "--period" in args
        assert expected_period in args


if __name__ == "__main__":
    import pytest
    pytest.main(["-v", __file__])
