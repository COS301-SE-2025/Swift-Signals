import unittest
from unittest.mock import patch, mock_open
import xml.etree.ElementTree as ET
import sys
import pathlib

if "tJunction" not in sys.modules:
    sys.path.insert(0, str(pathlib.Path(__file__).resolve().parent.parent / "intersections"))
    import tJunction


class TestTJunction(unittest.TestCase):

    @patch("tJunction.writeNodeFile")
    @patch("tJunction.writeEdgeFile")
    @patch("tJunction.writeConnectionFile")
    @patch("tJunction.generateTrips")
    @patch("tJunction.subprocess.run")
    @patch("builtins.open", new_callable=mock_open)
    @patch("tJunction.extractTrajectories", return_value=[])
    @patch("tJunction.parseNodes", return_value=[{"id": "n1"}])
    @patch("tJunction.parseEdges", return_value=[{"id": "e1"}])
    @patch("tJunction.parseConnections", return_value=[{"from": "e1", "to": "e2"}])
    @patch("os.makedirs")
    @patch("os.remove")
    def test_generate_full_flow(self, mock_remove, mock_makedirs, mock_parseCon, mock_parseEdg, mock_parseNod,
                                mock_extractTraj, mock_open_file, mock_subprocess_run, mock_generateTrips,
                                mock_writeCon, mock_writeEdge, mock_writeNode):
        # Setup mocks for reading logfile and tripinfo XML
        logfile_content = [
            "Vehicle 'veh1' performs emergency braking\n",
            "Vehicle 'veh2' performs emergency stop\n",
            "No issues here\n"
        ]
        tripinfo_xml = """<root>
            <tripinfo duration="100" waitingTime="10" routeLength="500"/>
            <tripinfo duration="200" waitingTime="20" routeLength="1000"/>
        </root>"""

        # Mock open() for logfile, tripinfoFile, and simulation_output.json
        def open_side_effect(file, mode='r', *args, **kwargs):
            if file == "tjunction_warnings.log" and 'r' in mode:
                mock_file = mock_open(read_data="".join(logfile_content)).return_value
                mock_file.__iter__.return_value = logfile_content
                return mock_file
            elif file == "tjunction_tripinfo.xml" and 'r' in mode:
                mock_file = mock_open(read_data=tripinfo_xml).return_value
                return mock_file
            elif file == "out/simulationOut/simulation_output.json" and 'w' in mode:
                return mock_open().return_value
            else:
                return mock_open().return_value
        mock_open_file.side_effect = open_side_effect

        params = {"Speed": 60, "Traffic Density": "medium", "seed": 42}
        results = tJunction.generate(params)

        # Check speed conversion and warnings
        self.assertIn("Total Vehicles", results)
        self.assertEqual(results["Total Vehicles"], 2)
        self.assertEqual(results["Emergency Brakes"], 1)
        self.assertEqual(results["Emergency Stops"], 1)
        self.assertEqual(results["Near collisions"], 2)

        # Check write files called with correct filenames
        mock_writeNode.assert_called_once_with("tjInt.nod.xml")
        mock_writeEdge.assert_called_once()
        args, _ = mock_writeEdge.call_args
        self.assertAlmostEqual(args[1], 60 * 1000 / 3600)  # speed conversion

        mock_writeCon.assert_called_once_with("tjInt.con.xml")

        # Check subprocess calls for netconvert and sumo run
        calls = [call_args[0][0] for call_args in mock_subprocess_run.call_args_list]
        self.assertIn([
            "netconvert",
            "--node-files=tjInt.nod.xml",
            "--edge-files=tjInt.edg.xml",
            "--connection-files=tjInt.con.xml",
            "-o", "tjunction.net.xml"
        ], calls)

        self.assertIn([
            "sumo",
            "-c", "tjunction.sumocfg",
            "--tripinfo-output", "tjunction_tripinfo.xml",
            "--fcd-output", "tjunction_fcd.xml",
            "--no-warnings", "false",
            "--message-log", "tjunction_warnings.log"
        ], calls)

        mock_generateTrips.assert_called_once_with(
            "tjunction.net.xml",
            "tjunction.rou.xml",
            "medium",
            params
        )

        # Check JSON output file creation
        mock_open_file.assert_any_call("out/simulationOut/simulation_output.json", "w")

        # Check directories created
        mock_makedirs.assert_any_call("out/simulationOut", exist_ok=True)

        # Check temp files removal
        expected_files = [
            "tjunction.net.xml", "tjunction.rou.xml", "tjunction.sumocfg", "tjunction_tripinfo.xml",
            "tjInt.nod.xml", "tjInt.edg.xml", "tjInt.con.xml", "tjunction_fcd.xml", "tjunction_warnings.log"
        ]
        for filename in expected_files:
            mock_remove.assert_any_call(filename)

    def test_writeNodeFile_writes_correct_content(self):
        with patch("builtins.open", mock_open()) as m:
            tJunction.writeNodeFile("nodes.xml")
            m().write.assert_called_once()
            written_content = m().write.call_args[0][0]
            self.assertIn('<node id="center"', written_content)

    def test_writeEdgeFile_writes_correct_content_and_speed(self):
        with patch("builtins.open", mock_open()) as m:
            tJunction.writeEdgeFile("edges.xml", speed=20)
            m().write.assert_called_once()
            content = m().write.call_args[0][0]
            self.assertIn('speed="20"', content)

    def test_writeConnectionFile_writes_correct_content(self):
        with patch("builtins.open", mock_open()) as m:
            tJunction.writeConnectionFile("conns.xml")
            m().write.assert_called_once()
            content = m().write.call_args[0][0]
            self.assertIn('<connection from="in_n1"', content)

    @patch("os.environ", {"SUMO_HOME": "/fake/sumo"})
    @patch("os.makedirs")
    @patch("subprocess.run")
    @patch("builtins.open", new_callable=mock_open)
    def test_generateTrips_runs_command(self, mock_open_file, mock_run, mock_makedirs):
        netFile = "net.xml"
        tripFile = "some_dir/trips.rou.xml"  # <-- add directory here
        params = {"seed": 123}
        tJunction.generateTrips(netFile, tripFile, "high", params)
        mock_makedirs.assert_called_with("some_dir", exist_ok=True)

    def test_extractTrajectories_parses_xml_correctly(self):
        xml_data = """<root>
            <timestep time="0.0">
                <vehicle id="veh1" x="1" y="2" speed="3"/>
                <vehicle id="veh2" x="4" y="5" speed="6"/>
            </timestep>
            <timestep time="1.0">
                <vehicle id="veh1" x="7" y="8" speed="9"/>
            </timestep>
        </root>"""
        with patch("xml.etree.ElementTree.parse") as mock_parse:
            mock_parse.return_value.getroot.return_value = ET.fromstring(xml_data)
            result = tJunction.extractTrajectories("dummy.xml")
            self.assertEqual(len(result), 2)
            veh1 = next(v for v in result if v["id"] == "veh1")
            self.assertEqual(len(veh1["positions"]), 2)
            veh2 = next(v for v in result if v["id"] == "veh2")
            self.assertEqual(len(veh2["positions"]), 1)

    def test_parseNodes_parses_correctly(self):
        xml = """<nodes>
            <node id="n1" x="1.1" y="2.2" type="priority"/>
        </nodes>"""
        with patch("xml.etree.ElementTree.parse") as mock_parse:
            mock_parse.return_value.getroot.return_value = ET.fromstring(xml)
            result = tJunction.parseNodes("file.xml")
            self.assertEqual(len(result), 1)
            self.assertEqual(result[0]["id"], "n1")
            self.assertEqual(result[0]["x"], 1.1)
            self.assertEqual(result[0]["y"], 2.2)
            self.assertEqual(result[0]["type"], "priority")

    def test_parseEdges_parses_correctly(self):
        xml = """<edges>
            <edge id="e1" from="n1" to="n2" speed="10" numLanes="3"/>
        </edges>"""
        with patch("xml.etree.ElementTree.parse") as mock_parse:
            mock_parse.return_value.getroot.return_value = ET.fromstring(xml)
            result = tJunction.parseEdges("file.xml")
            self.assertEqual(len(result), 1)
            self.assertEqual(result[0]["id"], "e1")
            self.assertEqual(result[0]["from"], "n1")
            self.assertEqual(result[0]["to"], "n2")
            self.assertEqual(result[0]["speed"], 10.0)
            self.assertEqual(result[0]["lanes"], 3)

    def test_parseConnections_parses_correctly(self):
        xml = """<connections>
            <connection from="e1" to="e2" fromLane="0" toLane="1"/>
        </connections>"""
        with patch("xml.etree.ElementTree.parse") as mock_parse:
            mock_parse.return_value.getroot.return_value = ET.fromstring(xml)
            result = tJunction.parseConnections("file.xml")
            self.assertEqual(len(result), 1)
            self.assertEqual(result[0]["from"], "e1")
            self.assertEqual(result[0]["to"], "e2")
            self.assertEqual(result[0]["fromLane"], 0)
            self.assertEqual(result[0]["toLane"], 1)

    @patch("builtins.print")
    @patch("tJunction.writeNodeFile")
    @patch("tJunction.writeEdgeFile")
    @patch("tJunction.writeConnectionFile")
    @patch("tJunction.generateTrips")
    @patch("tJunction.subprocess.run")
    @patch("tJunction.extractTrajectories", return_value=[])
    @patch("tJunction.parseNodes", return_value=[{"id": "n1"}])
    @patch("tJunction.parseEdges", return_value=[{"id": "e1"}])
    @patch("tJunction.parseConnections", return_value=[{"from": "e1", "to": "e2"}])
    @patch("os.makedirs")
    @patch("os.remove")
    @patch("xml.etree.ElementTree.parse")
    def test_generate_speed_warning_and_default(self, mock_et_parse, mock_remove, mock_makedirs, mock_parseCon,
                                                mock_parseEdg, mock_parseNod, mock_extractTraj, mock_subprocess_run,
                                                mock_generateTrips, mock_writeCon, mock_writeEdge, mock_writeNode,
                                                mock_print):

        tripinfo_xml = """<root>
            <tripinfo duration="100" waitingTime="10" routeLength="500"/>
        </root>"""
        mock_et_parse.return_value = ET.ElementTree(ET.fromstring(tripinfo_xml))

        params = {"Speed": 999, "Traffic Density": "low", "seed": 1}
        tJunction.generate(params)

        mock_print.assert_any_call("Warning: Speed 999km/h not allowed. Using default 40km/h.")


if __name__ == "__main__":
    unittest.main()
