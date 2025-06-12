import pytest
from unittest import mock
import sys
import os

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
import SimLoad


@pytest.mark.parametrize(
    "user_inputs, expected_function",
    [
        (["1", "medium"], "circle"),
        (["2", "high"], "stopStreet"),
        (["3", "low"], "tJunction"),
        (["4", "medium", "15", "10"], "trafficLight"),
    ],
)
def test_main_dispatch(user_inputs, expected_function, mocker):
    """
    Test that the correct generate() function is called based on user menu input.
    """
    mocker.patch("builtins.input", side_effect=user_inputs)

    circle = mocker.patch("SimLoad.circle.generate")
    stop = mocker.patch("SimLoad.stopStreet.generate")
    tjunc = mocker.patch("SimLoad.tJunction.generate")
    light = mocker.patch("SimLoad.trafficLight.generate")

    SimLoad.main()

    if expected_function == "circle":
        circle.assert_called_once_with({"Traffic Density": user_inputs[1]})
    elif expected_function == "stopStreet":
        stop.assert_called_once_with({"Traffic Density": user_inputs[1]})
    elif expected_function == "tJunction":
        tjunc.assert_called_once_with({"Traffic Density": user_inputs[1]})
    elif expected_function == "trafficLight":
        light.assert_called_once_with(
            {
                "Traffic Density": user_inputs[1],
                "Green": int(user_inputs[2]),
                "Red": int(user_inputs[3]),
            }
        )


def test_main_invalid_choice(mocker):
    """
    Input an invalid choice first, then a valid choice '1' to exit recursion.
    """
    inputs = ["invalid", "1", "medium"]
    mocker.patch("builtins.input", side_effect=inputs)

    circle = mocker.patch("SimLoad.circle.generate")

    SimLoad.main()

    circle.assert_called_once_with({"Traffic Density": "medium"})


def test_main_invalid_choice_with_recursion(mocker):
    """
    Sequence: invalid input, then valid '1' with density.
    This test is similar to test_main_invalid_choice but kept for completeness.
    """
    inputs = ["invalid", "1", "medium"]
    mocker.patch("builtins.input", side_effect=inputs)

    circle = mocker.patch("SimLoad.circle.generate")

    SimLoad.main()

    circle.assert_called_once_with({"Traffic Density": "medium"})
