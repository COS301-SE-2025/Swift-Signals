import unittest
from unittest.mock import patch, mock_open, MagicMock
import os
import sys
import pathlib
import xml.etree.ElementTree as ET

if "stopStreet" not in sys.modules:
    sys.path.insert(
        0, str(pathlib.Path(__file__).resolve().parent.parent / "intersections")
    )
    import stopStreet


class TestStopIntersection(unittest.TestCase):

    @patch("stopStreet.os.remove")
    @patch("stopStreet.extractTrajectories", return_value=[])
    @patch("stopStreet.parseNodes", return_value=[])
    @patch("stopStreet.parseEdges", return_value=[])
    @patch("stopStreet.parseConnections", return_value=[])
    @patch("stopStreet.generateTrips")
    @patch("stopStreet.subprocess.run")
    @patch("builtins.open", new_callable=mock_open)
    @patch("xml.etree.ElementTree.parse")
    def test_generate_full_flow(
        self,
        mock_et_parse,
        mock_open_file,
        mock_subprocess,
        mock_generateTrips,
        mock_parseCon,
        mock_parseEdg,
        mock_parseNod,
        mock_extractTraj,
        mock_remove,
    ):

        tripinfo_xml = """<root>
            <tripinfo duration="100" waitingTime="20" routeLength="500"/>
            <tripinfo duration="200" waitingTime="40" routeLength="1000"/>
        </root>"""

        tree_mock = MagicMock()
        root_mock = ET.ElementTree(ET.fromstring(tripinfo_xml)).getroot()
        tree_mock.getroot.return_value = root_mock
        mock_et_parse.return_value = tree_mock

        mock_cfg_file = mock_open().return_value
        mock_cfg_file.__enter__.return_value = mock_cfg_file

        mock_log_file = mock_open().return_value
        mock_log_file.__enter__.return_value = mock_log_file

        mock_default_file = mock_open().return_value
        mock_default_file.__enter__.return_value = mock_default_file

        def open_side_effect(file, mode="r", *args, **kwargs):
            if file.endswith("stop_intersection.sumocfg"):
                return mock_cfg_file
            elif file.endswith("stop_intersection_warnings.log"):
                return mock_log_file
            else:
                return mock_default_file

        mock_open_file.side_effect = open_side_effect

        params = {"Speed": 60, "Traffic Density": "medium", "seed": 42}
        results, fullOutput = stopStreet.generate(params)

        self.assertIn("Total Vehicles", results)
        self.assertEqual(results["Total Vehicles"], 2)

        calls = [call[0][0] for call in mock_subprocess.call_args_list]

        self.assertIn(
            [
                "netconvert",
                "--node-files=stopInt.nod.xml",
                "--edge-files=stopInt.edg.xml",
                "--connection-files=stopInt.con.xml",
                "-o",
                "stop_intersection.net.xml",
            ],
            calls,
        )

        self.assertIn(
            [
                "sumo",
                "-c",
                "stop_intersection.sumocfg",
                "--tripinfo-output",
                "stop_intersection_tripinfo.xml",
                "--fcd-output",
                "stop_intersection_fcd.xml",
                "--no-warnings",
                "false",
                "--message-log",
                "stop_intersection_warnings.log",
            ],
            calls,
        )

        mock_generateTrips.assert_called_once()
        mock_extractTraj.assert_called_once()
        mock_remove.assert_any_call("stop_intersection.net.xml")

    def test_writeNodeFile(self):
        with patch("builtins.open", mock_open()) as m:
            stopStreet.writeNodeFile("test_nodes.xml")
            m.assert_called_with("test_nodes.xml", "w")

    def test_writeEdgeFile(self):
        with patch("builtins.open", mock_open()) as m:
            stopStreet.writeEdgeFile("test_edges.xml", speed=20.0)
            m.assert_called_with("test_edges.xml", "w")

    def test_writeConnectionFile(self):
        with patch("builtins.open", mock_open()) as m:
            stopStreet.writeConnectionFile("test_conns.xml")
            m.assert_called_with("test_conns.xml", "w")

    def test_parseNodes(self):
        xml = """<nodes>
            <node id="n1" x="0" y="0" type="priority"/>
        </nodes>"""
        with patch("xml.etree.ElementTree.parse") as mock_parse:
            tree = ET.ElementTree(ET.fromstring(xml))
            mock_parse.return_value = tree
            nodes = stopStreet.parseNodes("dummy.xml")
            self.assertEqual(len(nodes), 1)
            self.assertEqual(nodes[0]["id"], "n1")

    def test_parseEdges(self):
        xml = """<edges>
            <edge id="e1" from="n1" to="n2" speed="15" numLanes="2"/>
        </edges>"""
        with patch("xml.etree.ElementTree.parse") as mock_parse:
            tree = ET.ElementTree(ET.fromstring(xml))
            mock_parse.return_value = tree
            edges = stopStreet.parseEdges("dummy.xml")
            self.assertEqual(edges[0]["speed"], 15.0)

    def test_parseConnections(self):
        xml = """<connections>
            <connection from="e1" to="e2" fromLane="0" toLane="0"/>
        </connections>"""
        with patch("xml.etree.ElementTree.parse") as mock_parse:
            tree = ET.ElementTree(ET.fromstring(xml))
            mock_parse.return_value = tree
            conns = stopStreet.parseConnections("dummy.xml")
            self.assertEqual(conns[0]["from"], "e1")

    def test_extractTrajectories(self):
        xml = """<root>
            <timestep time="1">
                <vehicle id="veh1" x="10" y="20" speed="5"/>
            </timestep>
        </root>"""
        with patch("xml.etree.ElementTree.parse") as mock_parse:
            tree = ET.ElementTree(ET.fromstring(xml))
            mock_parse.return_value = tree
            traj = stopStreet.extractTrajectories("dummy.xml")
            self.assertEqual(len(traj), 1)
            self.assertEqual(traj[0]["id"], "veh1")
            self.assertEqual(len(traj[0]["positions"]), 1)

    @patch("stopStreet.subprocess.run")
    @patch.dict(os.environ, {"SUMO_HOME": "C:/Sumo"})
    def test_generateTrips(self, mock_run):
        stopStreet.generateTrips("net.xml", "routes.rou.xml", "high", {"seed": 5})
        mock_run.assert_called()
