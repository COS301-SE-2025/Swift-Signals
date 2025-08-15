import subprocess
import json
import random
from unittest import mock
import GP

FAKE_RESULTS = {
    "intersection": {
        "results": {
            "Total Waiting Time": 10,
            "Total Travel Time": 5,
            "Emergency Brakes": 1,
            "Emergency Stops": 2,
            "Near collisions": 0,
        }
    }
}


def fake_json_load(f):
    return FAKE_RESULTS


@mock.patch("builtins.open", new_callable=mock.mock_open, read_data=json.dumps(FAKE_RESULTS))
@mock.patch("GP.subprocess.run")
@mock.patch("GP.datetime")
def test_full_run(mock_datetime, mock_run, mock_open):
    mock_datetime.now.return_value.strftime.return_value = "20200101-000000-000000"
    mock_run.return_value = None

    GP.ngen_waiting = 1
    GP.ngen_safety = 1
    GP.pop_size = 2

    GP.PARAMS_FOLDER = "tmp_params"
    GP.RESULTS_FOLDER = "tmp_results"
    GP.RESULT_FILE_TEMPLATE = "tmp_results/simulation_results_{}.json"
    GP.REFERENCE_RESULT = "reference.json"

    GP.generated_param_files.clear()
    GP.generated_result_files.clear()

    GP.main()

    assert mock_run.called
    assert mock_open.called


def test_custom_mutate_for_full_coverage():
    ind = GP.creator.Individual([10, 3, 10, 40, 1408])
    with mock.patch("random.random", return_value=0.0), \
         mock.patch("random.randint", side_effect=[20, 5, 30]), \
         mock.patch("random.choice", side_effect=[60]):
        mutated, = GP.custom_mutate(ind, indpb=1.0, min_speed=40)
    assert mutated[0] == 20
    assert mutated[1] == 5
    assert mutated[2] == 30
    assert mutated[3] == 60


@mock.patch("GP.run_simulation", return_value=FAKE_RESULTS["intersection"]["results"])
def test_evaluate_functions(mock_run_sim):
    ind_safe = GP.creator.Individual([10, 3, 10, 60, 1408])
    ind_unsafe = GP.creator.Individual([10, 3, 10, 40, 1408])

    wait_val = GP.evaluate_waiting_and_travel(ind_safe)
    safe_val = GP.evaluate_safety_given_waiting(ind_safe)
    unsafe_val = GP.evaluate_safety_given_waiting(ind_unsafe)

    assert isinstance(wait_val, tuple)
    assert isinstance(safe_val, tuple)
    assert unsafe_val[0] >= 1e6


@mock.patch("subprocess.run", side_effect=subprocess.CalledProcessError(1, "cmd"))
def test_run_simulation_subprocess_error(mock_run):
    ind = GP.creator.Individual([10, 3, 10, 40, 1408])
    result = GP.run_simulation(ind)
    assert result is None


@mock.patch("GP.subprocess.run")
@mock.patch("GP.json.load", side_effect=fake_json_load)
@mock.patch("GP.datetime")
def test_run_final_simulation_and_compare_normal(mock_datetime, mock_json_load, mock_run):
    mock_datetime.now.return_value.strftime.return_value = "20200101-000000-000000"
    best_params = {"Green": 10, "Yellow": 3, "Red": 10, "Speed": 60, "Seed": 1408}
    GP.run_final_simulation_and_compare(best_params)


@mock.patch("subprocess.run", side_effect=subprocess.CalledProcessError(1, "cmd"))
def test_run_final_simulation_error(mock_run):
    best_params = {"Green": 10, "Yellow": 3, "Red": 10, "Speed": 60, "Seed": 1408}
    GP.run_final_simulation_and_compare(best_params)


def test_run_final_simulation_read_error_specific_fixed():
    best_params = {"Green": 10, "Yellow": 3, "Red": 10, "Speed": 60, "Seed": 1408}

    real_open = open

    def side_effect_open(file, *args, **kwargs):
        if "final_simulation_result" in file:
            return mock.mock_open(read_data=json.dumps(FAKE_RESULTS)).return_value
        elif "simulation_results.json" in file:
            raise Exception("Read error")
        else:
            return real_open(file, *args, **kwargs)

    with mock.patch("builtins.open", side_effect=side_effect_open):
        GP.run_final_simulation_and_compare(best_params)
