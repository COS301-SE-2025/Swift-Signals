import unittest
from unittest.mock import patch, mock_open, MagicMock
import os
import pathlib
import sys
import xml.etree.ElementTree as ET


if "circle" not in sys.modules:
    sys.path.insert(
        0, str(pathlib.Path(__file__).resolve().parent.parent / "intersections")
    )
    import circle


class TestRoundabout(unittest.TestCase):

    @patch("circle.os.remove")
    @patch("circle.extractTrajectories", return_value=[])
    @patch("circle.parseNodes", return_value=[])
    @patch("circle.parseEdges", return_value=[])
    @patch("circle.parseConnections", return_value=[])
    @patch("circle.generateTrips")
    @patch("circle.subprocess.run")
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
            <tripinfo duration="60" waitingTime="5" routeLength="500"/>
            <tripinfo duration="50" waitingTime="3" routeLength="400"/>
        </root>"""

        tree_mock = MagicMock()
        root_mock = ET.ElementTree(ET.fromstring(tripinfo_xml)).getroot()
        tree_mock.getroot.return_value = root_mock
        mock_et_parse.return_value = tree_mock

        params = {"Speed": 60, "Traffic Density": "medium", "seed": 42}

        results, fullOutput = circle.generate(params)

        self.assertIn("Total Vehicles", results)
        self.assertEqual(results["Total Vehicles"], 2)
        self.assertGreater(results["Average Speed"], 0)
        self.assertAlmostEqual(results["Average Waiting Time"], 4)

        mock_subprocess.assert_any_call(
            [
                "netconvert",
                "--node-files=roundabout.nod.xml",
                "--edge-files=roundabout.edg.xml",
                "--connection-files=roundabout.con.xml",
                "-o",
                "roundabout.net.xml",
            ],
            check=True,
        )

        mock_generateTrips.assert_called_once()
        mock_extractTraj.assert_called_once()
        mock_remove.assert_any_call("roundabout.net.xml")

    @patch("builtins.open", new_callable=mock_open)
    def test_writeNodeFile(self, mock_file):
        circle.writeNodeFile("test.nod.xml")
        mock_file.assert_called_with("test.nod.xml", "w")

    @patch("builtins.open", new_callable=mock_open)
    def test_writeEdgeFile(self, mock_file):
        circle.writeEdgeFile("test.edg.xml", speed=15.0)
        mock_file.assert_called_with("test.edg.xml", "w")

    @patch("builtins.open", new_callable=mock_open)
    def test_writeConnectionFile(self, mock_file):
        circle.writeConnectionFile("test.con.xml")
        mock_file.assert_called_with("test.con.xml", "w")

    @patch("circle.subprocess.run")
    @patch("builtins.open", new_callable=mock_open)
    @patch.dict(os.environ, {"SUMO_HOME": "C:/Sumo"})
    def test_generateTrips(self, mock_file, mock_subprocess):
        params = {"seed": 7}
        circle.generateTrips("test.net.xml", "test.rou.xml", "high", params)
        self.assertTrue(mock_subprocess.called)

    def test_extractTrajectories(self):
        fcd_xml = """<root>
            <timestep time="1.0">
                <vehicle id="veh1" x="10.0" y="20.0" speed="5.5"/>
                <vehicle id="veh2" x="15.0" y="25.0" speed="6.0"/>
            </timestep>
            <timestep time="2.0">
                <vehicle id="veh1" x="12.0" y="22.0" speed="5.0"/>
            </timestep>
        </root>"""

        with patch("xml.etree.ElementTree.parse") as mock_parse:
            tree = ET.ElementTree(ET.fromstring(fcd_xml))
            mock_parse.return_value = tree

            result = circle.extractTrajectories("dummy.xml")
            self.assertEqual(len(result), 2)
            self.assertEqual(result[0]["id"], "veh1")
            self.assertEqual(len(result[0]["positions"]), 2)

    def test_parseNodes(self):
        xml = """<nodes>
            <node id="n1" x="0" y="0" type="priority"/>
            <node id="n2" x="1" y="2" type="priority"/>
        </nodes>"""

        with patch("xml.etree.ElementTree.parse") as mock_parse:
            tree = ET.ElementTree(ET.fromstring(xml))
            mock_parse.return_value = tree
            nodes = circle.parseNodes("dummy.xml")
            self.assertEqual(len(nodes), 2)
            self.assertEqual(nodes[1]["x"], 1)

    def test_parseEdges(self):
        xml = """<edges>
            <edge id="e1" from="n1" to="n2" speed="15" numLanes="2"/>
        </edges>"""

        with patch("xml.etree.ElementTree.parse") as mock_parse:
            tree = ET.ElementTree(ET.fromstring(xml))
            mock_parse.return_value = tree
            edges = circle.parseEdges("dummy.xml")
            self.assertEqual(edges[0]["speed"], 15.0)

    def test_parseConnections(self):
        xml = """<connections>
            <connection from="e1" to="e2" fromLane="0" toLane="0"/>
        </connections>"""

        with patch("xml.etree.ElementTree.parse") as mock_parse:
            tree = ET.ElementTree(ET.fromstring(xml))
            mock_parse.return_value = tree
            conns = circle.parseConnections("dummy.xml")
            self.assertEqual(conns[0]["from"], "e1")


if __name__ == "__main__":
    unittest.main()
