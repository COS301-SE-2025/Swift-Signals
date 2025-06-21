import sys
import unittest
from unittest.mock import patch, mock_open
import uuid
import pathlib

if "SimLoad" not in sys.modules:
    sys.path.insert(0, str(pathlib.Path(__file__).resolve().parent.parent))
    import SimLoad


class TestSimLoad(unittest.TestCase):

    @patch("builtins.input", side_effect=["2"])
    def test_showMenu(self, mock_input):
        choice = SimLoad.showMenu()
        self.assertEqual(choice, "2")

    @patch("builtins.input", side_effect=["params_test.json"])
    @patch("os.path.exists", return_value=True)
    @patch("builtins.open", new_callable=mock_open, read_data='{"intersection":{"simulation_parameters":{"Intersection Type":"trafficlight"}}}')
    def test_loadParams(self, mock_open_file, mock_exists, mock_input):
        params = SimLoad.loadParams()
        self.assertEqual(params["Intersection Type"], "trafficlight")

    @patch("os.path.exists", return_value=True)
    @patch("builtins.open", new_callable=mock_open, read_data="5")
    def test_loadRunCount_existing(self, mock_file, mock_exists):
        count = SimLoad.loadRunCount()
        self.assertEqual(count, 5)

    @patch("builtins.open", new_callable=mock_open)
    def test_saveRunCount(self, mock_file):
        SimLoad.saveRunCount(7)
        mock_file().write.assert_called_with("7")

    def test_getDefaultTimingsBySpeed(self):
        self.assertEqual(SimLoad.getDefaultTimingsBySpeed(40), {"Green": 25, "Yellow": 3, "Red": 30})
        self.assertEqual(SimLoad.getDefaultTimingsBySpeed(80), {"Green": 30, "Yellow": 5, "Red": 35})
        self.assertEqual(SimLoad.getDefaultTimingsBySpeed(120), {"Green": 30, "Yellow": 5, "Red": 35})

    @patch("builtins.input", side_effect=["medium", "60", "y"])
    def test_getParams_trafficlight_default(self, mock_input):
        result = SimLoad.getParams(True)
        self.assertEqual(result["Speed"], 60)
        self.assertIn("Green", result)

    @patch("builtins.input", side_effect=["high", "80"])
    def test_getParams_non_trafficlight(self, mock_input):
        result = SimLoad.getParams(False)
        self.assertEqual(result["Speed"], 80)
        self.assertNotIn("Green", result)

    @patch("uuid.uuid4", return_value=uuid.UUID("12345678-1234-5678-1234-567812345678"))
    @patch("builtins.open", new_callable=mock_open)
    def test_saveParams(self, mock_file, mock_uuid):
        SimLoad.saveParams({"Traffic Density": "medium", "Speed": 60}, "trafficlight", "testSim")
        self.assertTrue(mock_file.called)

    @patch("SimLoad.trafficLight.generate", return_value={"result": "ok"})
    @patch("SimLoad.loadParams", return_value={"Intersection Type": "trafficlight", "Speed": 60})
    @patch("SimLoad.saveRunCount")
    @patch("SimLoad.loadRunCount", return_value=0)
    @patch("SimLoad.json.dump")  # <-- patch json.dump instead of open().write
    def test_main_traffic_light(self, mock_json_dump, mock_run_count, mock_save_run_count, mock_params, mock_generate):
        SimLoad.main()
        args, kwargs = mock_json_dump.call_args
        data_written = args[0]  # this is the object passed to json.dump
        self.assertIn("intersection", data_written)
        self.assertEqual(data_written["intersection"]["parameters"]["Intersection Type"], "trafficlight")

    @patch("builtins.input", return_value="nonexistent.json")
    @patch("os.path.exists", return_value=False)
    def test_loadParams_file_not_found(self, mock_exists, mock_input):
        with self.assertRaises(SystemExit) as cm:
            SimLoad.loadParams()
        self.assertEqual(cm.exception.code, 1)

    @patch("os.path.exists", return_value=False)
    def test_loadRunCount_file_missing(self, mock_exists):
        count = SimLoad.loadRunCount()
        self.assertEqual(count, 0)

    @patch("builtins.input", side_effect=["medium", "invalid_speed", "y"])
    def test_getParams_invalid_speed_defaults_to_40(self, mock_input):
        result = SimLoad.getParams(True)
        self.assertEqual(result["Speed"], 40)
        self.assertIn("Green", result)  # Confirm light timings were added

    @patch("builtins.input", side_effect=["high", "80", "n", "20", "4", "25"])
    def test_getParams_custom_light_timings(self, mock_input):
        result = SimLoad.getParams(True)
        self.assertEqual(result["Green"], 20)
        self.assertEqual(result["Yellow"], 4)
        self.assertEqual(result["Red"], 25)
        self.assertEqual(result["Speed"], 80)

    @patch("builtins.input", side_effect=["high", "60", "n", "oops", "nope", "uhh"])
    def test_getParams_invalid_light_timings_fallback(self, mock_input):
        result = SimLoad.getParams(True)
        self.assertEqual(result["Green"], 25)   # default for 60 km/h
        self.assertEqual(result["Yellow"], 4)
        self.assertEqual(result["Red"], 30)

    @patch("builtins.open", new_callable=mock_open)
    @patch("SimLoad.circle.generate", return_value={"result": "stop_ok"})
    @patch("SimLoad.loadParams", return_value={"Intersection Type": "roundabout"})
    @patch("SimLoad.saveRunCount")
    @patch("SimLoad.loadRunCount", return_value=0)
    def test_main_roundabout(self, mock_run_count, mock_save_run_count, mock_params, mock_generate, mock_file):
        SimLoad.main()
        mock_generate.assert_called_once_with({"Intersection Type": "roundabout"})

    @patch("builtins.open", new_callable=mock_open)
    @patch("SimLoad.stopStreet.generate", return_value={"result": "stop_ok"})
    @patch("SimLoad.loadParams", return_value={"Intersection Type": "fourwaystop"})
    @patch("SimLoad.saveRunCount")
    @patch("SimLoad.loadRunCount", return_value=0)
    def test_main_fourwaystop(self, mock_run_count, mock_save_run_count, mock_params, mock_generate, mock_file):
        SimLoad.main()
        mock_generate.assert_called_once_with({"Intersection Type": "fourwaystop"})

    @patch("builtins.open", new_callable=mock_open)
    @patch("SimLoad.tJunction.generate", return_value={"result": "stop_ok"})
    @patch("SimLoad.loadParams", return_value={"Intersection Type": "tjunction"})
    @patch("SimLoad.saveRunCount")
    @patch("SimLoad.loadRunCount", return_value=0)
    def test_main_tjunction(self, mock_run_count, mock_save_run_count, mock_params, mock_generate, mock_file):
        SimLoad.main()
        mock_generate.assert_called_once_with({"Intersection Type": "tjunction"})

    @patch("SimLoad.loadParams", return_value={"Intersection Type": "invalidtype"})
    @patch("builtins.print")
    def test_main_invalid_intersection_type(self, mock_print, mock_params):
        result = SimLoad.main()
        mock_print.assert_called_with("Invalid intersection type in parameters.")
        self.assertIsNone(result)


if __name__ == "__main__":
    unittest.main()
