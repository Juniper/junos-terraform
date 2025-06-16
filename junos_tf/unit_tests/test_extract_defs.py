import json
import pytest
import os
from glob import glob
import sys
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "lib")))
from generate_plugin import filter_json_using_xml

DATA_DIR = os.path.join(os.path.dirname(__file__), "tests", "data_tests")

def load_test_cases():
    """Load all test cases from the test data directory."""
    json_files = sorted(glob(os.path.join(DATA_DIR, "*_input.json")))
    test_cases = []

    for json_file in json_files:
        base_name = os.path.basename(json_file).replace("_input.json", "")
        xml_file = os.path.join(DATA_DIR, f"{base_name}_input.xml")
        expected_file = os.path.join(DATA_DIR, f"{base_name}_expected.json")

        with open(json_file, "r") as jf, open(xml_file, "r") as xf, open(expected_file, "r") as ef:
            input_json = jf.read()
            input_xml = xf.read()
            expected_json = json.load(ef)
            test_cases.append((input_json, input_xml, expected_json))

    return test_cases

@pytest.mark.parametrize("input_json, input_xml, expected_json", load_test_cases())
def test_filter_json(input_json, input_xml, expected_json, tmp_path):
    """Run schema filtering test for each input/output trio."""
    json_path = tmp_path / "input.json"
    xml_path = tmp_path / "input.xml"

    json_path.write_text(input_json)
    xml_path.write_text(input_xml)

    result = filter_json_using_xml(str(json_path), str(xml_path))
    assert result == expected_json, f"""
    ❌ Output did not match expected result.
    --- Expected ---
    {json.dumps(expected_json, indent=2)}
    --- Got ---
    {json.dumps(result, indent=2)}
    """
