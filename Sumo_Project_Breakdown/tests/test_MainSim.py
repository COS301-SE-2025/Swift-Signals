import sys
import unittest
from unittest.mock import patch, mock_open
import uuid
import pathlib
from io import StringIO
import builtins


real_open = builtins.open

if "SimLoad" not in sys.modules:
    sys.path.insert(0, str(pathlib.Path(__file__).resolve().parent.parent))
    import SimLoad


class TestSimLoad(unittest.TestCase):

    @patch("builtins.input", side_effect=["params_test.json"])
    @patch("os.path.exists", return_value=True)
    @patch("builtins.open")
    @patch("sys.argv", new=["scriptname"])
    def test_loadParams(self, mock_open_file, mock_exists, mock_input):

        def open_side_effect(file, mode="r", *args, **kwargs):
            if file == "params_test.json":
                return StringIO(
                    '{"intersection":{"simulation_parameters":{"Intersection Type":1, "Speed":60, "seed":42}, "Traffic Density":2}}'
                )
            else:
                return real_open(file, mode, *args, **kwargs)

        mock_open_file.side_effect = open_side_effect

        params = SimLoad.loadParams()

        self.assertIn("mapped", params)
        self.assertIn("raw", params)

        self.assertEqual(params["mapped"]["Intersection Type"], "trafficlight")
        self.assertEqual(params["mapped"]["Traffic Density"], "high")
        self.assertEqual(params["mapped"]["Speed"], 60)
        self.assertEqual(params["mapped"]["seed"], 42)

        self.assertEqual(params["raw"]["Intersection Type"], 1)
        self.assertEqual(params["raw"]["Traffic Density"], 2)
        self.assertEqual(params["raw"]["Speed"], 60)
        self.assertEqual(params["raw"]["seed"], 42)

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
        self.assertEqual(
            SimLoad.getDefaultTimingsBySpeed(40), {"Green": 25, "Yellow": 3, "Red": 30}
        )
        self.assertEqual(
            SimLoad.getDefaultTimingsBySpeed(80), {"Green": 30, "Yellow": 5, "Red": 35}
        )
        self.assertEqual(
            SimLoad.getDefaultTimingsBySpeed(120), {"Green": 30, "Yellow": 5, "Red": 35}
        )

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
    @patch("time.strftime", return_value="20250623T000000Z")
    @patch("json.dump")
    def test_saveParams(self, mock_json_dump, mock_time, mock_file, mock_uuid):
        SimLoad.saveParams(
            {"Traffic Density": "medium", "Speed": 60}, "trafficlight", "testSim"
        )
        args, kwargs = mock_json_dump.call_args
        written_data = args[0]
        self.assertIn("intersection", written_data)
        self.assertEqual(
            written_data["intersection"]["parameters"]["Intersection Type"],
            1,
        )

    @patch("SimLoad.trafficLight.generate", return_value=(
        {"parameters": {"Green": 25, "Yellow": 3, "Red": 30, "Intersection Type": 1}},
        {"simulation_log": "details_here"},
    ))
    @patch("SimLoad.loadParams", return_value={
        "mapped": {
            "Intersection Type": "trafficlight",
            "Green": 25,
            "Yellow": 3,
            "Red": 30,
            "Traffic Density": "medium",
            "Speed": 60,
            "seed": 42,
        },
        "raw": {
            "Intersection Type": 1,
            "Traffic Density": 1,
            "Speed": 60,
            "seed": 42,
        },
    })
    @patch("SimLoad.saveRunCount")
    @patch("SimLoad.loadRunCount", return_value=0)
    @patch("json.dump")
    @patch("os.makedirs")
    @patch("os.remove")
    @patch("builtins.open", new_callable=mock_open)
    def test_main_traffic_light(
        self,
        mock_file,
        mock_remove,
        mock_makedirs,
        mock_json_dump,
        mock_run_count,
        mock_save_run_count,
        mock_params,
        mock_generate,
    ):
        SimLoad.main()

        mock_generate.assert_called_once()

        self.assertEqual(mock_json_dump.call_count, 2)

        first_output = mock_json_dump.call_args_list[0][0][0]
        second_output = mock_json_dump.call_args_list[1][0][0]

        self.assertIn("intersection", first_output)
        self.assertEqual(second_output, {"simulation_log": "details_here"})

        expected_dirs = {"out/results"}
        actual_dirs = {call_args[0][0] for call_args in mock_makedirs.call_args_list}

        self.assertTrue(expected_dirs.issubset(actual_dirs))

    @patch("builtins.input", return_value="nonexistent.json")
    @patch("os.path.exists", return_value=False)
    @patch("sys.argv", new=["scriptname"])
    def test_loadParams_file_not_found(self, mock_exists, mock_input):
        with self.assertRaises(SystemExit) as cm:
            SimLoad.loadParams()
        self.assertEqual(cm.exception.code, 1)

    @patch("os.path.exists", return_value=False)
    def test_loadRunCount_file_missing(self, mock_exists):
        count = SimLoad.loadRunCount()
        self.assertEqual(count, 0)

    @patch("builtins.input", side_effect=["high", "invalid_speed", "y"])
    def test_getParams_invalid_speed_defaults_to_40(self, mock_input):
        result = SimLoad.getParams(True)
        self.assertEqual(result["Speed"], 40)
        self.assertIn("Green", result)

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
        self.assertEqual(result["Green"], 25)
        self.assertEqual(result["Yellow"], 4)
        self.assertEqual(result["Red"], 30)

    @patch("SimLoad.circle.generate", return_value=({"result": "circle_ok"}, {}))
    @patch(
        "SimLoad.loadParams",
        return_value={
            "mapped": {
                "Intersection Type": "roundabout",
                "Traffic Density": "medium",
                "Speed": 60,
                "seed": 42,
            },
            "raw": {
                "Intersection Type": 2,
                "Traffic Density": 1,
                "Speed": 60,
                "seed": 42,
            },
        },
    )
    @patch("SimLoad.saveRunCount")
    @patch("SimLoad.loadRunCount", return_value=0)
    @patch("os.makedirs")
    @patch("os.remove")
    @patch("builtins.open", new_callable=mock_open)
    def test_main_roundabout(
        self,
        mock_file,
        mock_remove,
        mock_makedirs,
        mock_run_count,
        mock_save_run_count,
        mock_params,
        mock_generate,
    ):
        SimLoad.main()
        mock_generate.assert_called_once()

    @patch("SimLoad.stopStreet.generate", return_value=({"result": "stop_ok"}, {}))
    @patch(
        "SimLoad.loadParams",
        return_value={
            "mapped": {
                "Intersection Type": "fourwaystop",
                "Traffic Density": "medium",
                "Speed": 60,
                "seed": 42,
            },
            "raw": {
                "Intersection Type": 3,
                "Traffic Density": 1,
                "Speed": 60,
                "seed": 42,
            },
        },
    )
    @patch("SimLoad.saveRunCount")
    @patch("SimLoad.loadRunCount", return_value=0)
    @patch("os.makedirs")
    @patch("os.remove")
    @patch("builtins.open", new_callable=mock_open)
    def test_main_fourwaystop(
        self,
        mock_file,
        mock_remove,
        mock_makedirs,
        mock_run_count,
        mock_save_run_count,
        mock_params,
        mock_generate,
    ):
        SimLoad.main()
        mock_generate.assert_called_once()

    @patch("SimLoad.tJunction.generate", return_value=({"result": "tj_ok"}, {}))
    @patch(
        "SimLoad.loadParams",
        return_value={
            "mapped": {
                "Intersection Type": "tjunction",
                "Traffic Density": "medium",
                "Speed": 60,
                "seed": 42,
            },
            "raw": {
                "Intersection Type": 4,
                "Traffic Density": 1,
                "Speed": 60,
                "seed": 42,
            },
        },
    )
    @patch("SimLoad.saveRunCount")
    @patch("SimLoad.loadRunCount", return_value=0)
    @patch("os.makedirs")
    @patch("os.remove")
    @patch("builtins.open", new_callable=mock_open)
    def test_main_tjunction(
        self,
        mock_file,
        mock_remove,
        mock_makedirs,
        mock_run_count,
        mock_save_run_count,
        mock_params,
        mock_generate,
    ):
        SimLoad.main()
        mock_generate.assert_called_once()

    @patch(
        "SimLoad.loadParams",
        return_value={
            "mapped": {"Intersection Type": "invalidtype"},
            "raw": {"Intersection Type": -1},
        },
    )
    @patch("builtins.print")
    def test_main_invalid_intersection_type(self, mock_print, mock_params):
        result = SimLoad.main()
        mock_print.assert_called_with("Invalid intersection type in parameters.")
        self.assertIsNone(result)


if __name__ == "__main__":
    unittest.main()
