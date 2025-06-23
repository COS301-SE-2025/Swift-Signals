import unittest
from unittest.mock import patch, mock_open
import xml.etree.ElementTree as ET
import sys
import pathlib

if "trafficLight" not in sys.modules:
    sys.path.insert(
        0, str(pathlib.Path(__file__).resolve().parent.parent / "intersections")
    )
    import trafficLight


class TestTLIntersection(unittest.TestCase):

    @patch("builtins.open", new_callable=mock_open)
    def test_writeNodeFile_creates_expected_content(self, m):
        trafficLight.writeNodeFile("nodes.xml")
        m().write.assert_called_once()
        content = m().write.call_args[0][0]
        self.assertIn('<node id="1" x="0" y="0" type="traffic_light"/>', content)

    @patch("builtins.open", new_callable=mock_open)
    def test_writeEdgeFile_creates_expected_content_with_speed(self, m):
        trafficLight.writeEdgeFile("edges.xml", speed=15)
        m().write.assert_called_once()
        content = m().write.call_args[0][0]
        self.assertIn('speed="15"', content)

    @patch("builtins.open", new_callable=mock_open)
    def test_writeConnectionFile_creates_expected_content(self, m):
        trafficLight.writeConnectionFile("connections.xml")
        m().write.assert_called_once()
        content = m().write.call_args[0][0]
        self.assertIn('<connection from="in_n2_1" to="out_1_n3"', content)

    @patch("builtins.open", new_callable=mock_open)
    def test_writeTrafficLightLogic_creates_expected_phases(self, m):
        trafficLight.writeTrafficLightLogic("tll.xml", 30, 5, 25)
        m().write.assert_called_once()
        content = m().write.call_args[0][0]
        self.assertIn('<phase duration="30"', content)
        self.assertIn('state="GGGrrrGGGrrr"', content)

    @patch("xml.etree.ElementTree.parse")
    def test_parseNodes_returns_expected_data(self, mock_parse):
        xml = "<nodes><node id='n1' x='1' y='2' type='priority'/></nodes>"
        mock_parse.return_value.getroot.return_value = ET.fromstring(xml)
        result = trafficLight.parseNodes("file.xml")
        self.assertEqual(result[0]["id"], "n1")
        self.assertEqual(result[0]["x"], 1.0)
        self.assertEqual(result[0]["y"], 2.0)

    @patch("xml.etree.ElementTree.parse")
    def test_parseEdges_returns_expected_data(self, mock_parse):
        xml = "<edges><edge id='e1' from='n1' to='n2' speed='20' numLanes='2'/></edges>"
        mock_parse.return_value.getroot.return_value = ET.fromstring(xml)
        result = trafficLight.parseEdges("file.xml")
        self.assertEqual(result[0]["id"], "e1")
        self.assertEqual(result[0]["speed"], 20.0)

    @patch("xml.etree.ElementTree.parse")
    def test_parseConnections_returns_expected_data(self, mock_parse):
        xml = "<connections><connection from='e1' to='e2' fromLane='0' toLane='1' tl='0'/></connections>"
        mock_parse.return_value.getroot.return_value = ET.fromstring(xml)
        result = trafficLight.parseConnections("file.xml")
        self.assertEqual(result[0]["from"], "e1")
        self.assertEqual(result[0]["tl"], "0")

    @patch("xml.etree.ElementTree.parse")
    def test_parseTrafficLights_returns_expected_data(self, mock_parse):
        xml = """
        <additional>
            <tlLogic id="1" type="static">
                <phase duration="30" state="GGGrrrGGGrrr"/>
            </tlLogic>
        </additional>
        """
        mock_parse.return_value.getroot.return_value = ET.fromstring(xml)
        result = trafficLight.parseTrafficLights("file.xml")
        self.assertEqual(result[0]["id"], "1")
        self.assertEqual(result[0]["phases"][0]["duration"], 30)

    @patch("xml.etree.ElementTree.parse")
    def test_extractTrajectories_returns_expected_data(self, mock_parse):
        xml = """
        <root>
            <timestep time="0">
                <vehicle id="v1" x="1" y="2" speed="3"/>
            </timestep>
        </root>
        """
        mock_parse.return_value.getroot.return_value = ET.fromstring(xml)
        result = trafficLight.extractTrajectories("file.xml")
        self.assertEqual(result[0]["id"], "v1")
        self.assertEqual(result[0]["positions"][0]["speed"], 3.0)

    @patch("trafficLight.os.remove")
    @patch("trafficLight.extractTrajectories", return_value=[])
    @patch("trafficLight.parseTrafficLights", return_value=[])
    @patch("trafficLight.parseConnections", return_value=[])
    @patch("trafficLight.parseEdges", return_value=[])
    @patch("trafficLight.parseNodes", return_value=[])
    @patch("builtins.open", new_callable=mock_open)
    @patch("xml.etree.ElementTree.parse")
    @patch("trafficLight.generateTrips")
    @patch("trafficLight.subprocess.run")
    def test_generate_runs_full_flow(
        self,
        mock_run,
        mock_generateTrips,
        mock_et_parse,
        mock_open_file,
        mock_parseNodes,
        mock_parseEdges,
        mock_parseCon,
        mock_parseTL,
        mock_extractTraj,
        mock_remove,
    ):
        # Mock tripinfo XML
        tripinfo_xml = """
        <root>
            <tripinfo duration="100" waitingTime="10" routeLength="1000"/>
            <tripinfo duration="200" waitingTime="20" routeLength="2000"/>
        </root>
        """
        mock_et_parse.return_value.getroot.return_value = ET.fromstring(tripinfo_xml)

        # Mock open for reading logs and writing outputs
        def open_side_effect(file, mode="r", *args, **kwargs):
            if file.endswith("_warnings.log") and "r" in mode:
                mock_file = mock_open(
                    read_data="Vehicle 'v1' performs emergency braking\n"
                ).return_value
                mock_file.__iter__.return_value = [
                    "Vehicle 'v1' performs emergency braking\n"
                ]
                return mock_file
            if file.endswith("_tripinfo.xml") and "r" in mode:
                return mock_open(read_data=tripinfo_xml).return_value
            return mock_open().return_value

        mock_open_file.side_effect = open_side_effect

        params = {
            "Speed": 60,
            "Traffic Density": "low",
            "seed": 1,
            "Green": 30,
            "Yellow": 5,
            "Red": 25,
        }
        result = trafficLight.generate(params)

        self.assertIn("Total Vehicles", result)
        self.assertEqual(result["Total Vehicles"], 2)
        mock_generateTrips.assert_called_once()
        self.assertTrue(mock_remove.called)

    @patch("builtins.print")
    @patch("trafficLight.generateTrips")
    @patch("trafficLight.subprocess.run")
    @patch("builtins.open", new_callable=mock_open)
    @patch("xml.etree.ElementTree.parse")
    @patch("trafficLight.extractTrajectories", return_value=[])
    @patch("trafficLight.parseTrafficLights", return_value=[])
    @patch("trafficLight.parseConnections", return_value=[])
    @patch("trafficLight.parseEdges", return_value=[])
    @patch("trafficLight.parseNodes", return_value=[])
    @patch("trafficLight.os.remove")
    def test_generate_speed_warning_prints_message(
        self,
        mock_remove,
        mock_parseNodes,
        mock_parseEdges,
        mock_parseCon,
        mock_parseTL,
        mock_extractTraj,
        mock_et_parse,
        mock_open_file,
        mock_run,
        mock_generateTrips,
        mock_print,
    ):
        tripinfo_xml = """
        <root>
            <tripinfo duration="100" waitingTime="10" routeLength="1000"/>
        </root>
        """
        mock_et_parse.return_value.getroot.return_value = ET.fromstring(tripinfo_xml)

        def open_side_effect(file, mode="r", *args, **kwargs):
            if file.endswith("_warnings.log") and "r" in mode:
                mock_file = mock_open(read_data="").return_value
                mock_file.__iter__.return_value = []
                return mock_file
            if file.endswith("_tripinfo.xml") and "r" in mode:
                return mock_open(read_data=tripinfo_xml).return_value
            return mock_open().return_value

        mock_open_file.side_effect = open_side_effect

        params = {
            "Speed": 999,
            "Traffic Density": "medium",
            "seed": 1,
            "Green": 30,
            "Yellow": 5,
            "Red": 25,
        }
        trafficLight.generate(params)

        mock_print.assert_any_call(
            "Warnig: Speed 999km/h not allowed. Using default 40km/h."
        )

    @patch("trafficLight.os.makedirs")
    @patch("trafficLight.subprocess.run")
    def test_generateTrips_builds_expected_command(self, mock_run, mock_makedirs):
        with patch.dict("trafficLight.os.environ", {"SUMO_HOME": "/fake/sumo"}):
            trafficLight.generateTrips(
                "net.xml", "trips/trips.rou.xml", "high", {"seed": 42}
            )
            mock_makedirs.assert_called_with("trips", exist_ok=True)

            args = mock_run.call_args[0][0]

            trips_path = args[1]

            self.assertTrue(
                trips_path.endswith("randomTrips.py"), f"Unexpected path: {trips_path}"
            )

            self.assertIn("--period", args)
            self.assertIn("3", args)


if __name__ == "__main__":
    unittest.main()
