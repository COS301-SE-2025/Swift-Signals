from ga.main import main as run_optimisation
from pprint import pprint


def test_run_optimisation():
    """
    Test the run_optimisation function.
    """
    # Call the function to test
    result = run_optimisation()

    # Check if the result is a dictionary
    assert isinstance(result, dict), "Result should be a dictionary"

    print("Test passed: run_optimisation returns a dictionary.")
    pprint(result)  # Print the result for inspection


if __name__ == "__main__":
    test_run_optimisation()
