import json
import urllib.parse

def json_to_urlencoded(json_obj, base_key=None):
    """
    Convert a JSON object to a URL-encoded string compatible with Spring MVC.
    
    Args:
    - json_obj: The JSON object to convert.
    - base_key: The base key for nested objects, used internally for recursion.
    
    Returns:
    A URL-encoded string representation of the JSON object.
    """
    urlencoded_parts = []

    if isinstance(json_obj, dict):
        for key, value in json_obj.items():
            new_key = f"{base_key}[{key}]" if base_key else key
            urlencoded_parts.extend(json_to_urlencoded(value, new_key))
    elif isinstance(json_obj, list):
        for index, item in enumerate(json_obj):
            new_key = f"{base_key}[{index}]"
            urlencoded_parts.extend(json_to_urlencoded(item, new_key))
    else:
        encoded_value = urllib.parse.quote_plus(str(json_obj))
        urlencoded_parts.append(f"{base_key}={encoded_value}")

    return urlencoded_parts

def convert_json_payload_to_urlencoded(json_payload):
    """
    Convert a JSON payload (as a string) to a URL-encoded form string.
    
    Args:
    - json_payload: The JSON payload as a string.
    
    Returns:
    A URL-encoded string compatible with Spring MVC's data binding.
    """
    json_obj = json.loads(json_payload)
    urlencoded_parts = json_to_urlencoded(json_obj)
    return "&".join(urlencoded_parts)

# Example JSON payload
json_payload = """
{
    "name": "John Doe",
    "age": 30,
    "address": {
        "street": "123 Main St",
        "city": "Anytown"
    },
    "hobbies": ["reading", "cycling", "hiking"]
}
"""

# Convert JSON to URL-encoded form
urlencoded_payload = convert_json_payload_to_urlencoded(json_payload)
print(urlencoded_payload)
