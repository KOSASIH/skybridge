import os

# Load API key from environment variable or use a default value
API_KEY = os.getenv("STARLINK_API_KEY", "your_default_starlink_api_key")

# Load threshold latency from environment variable or use a default value
THRESHOLD_LATENCY = int(os.getenv("THRESHOLD_LATENCY", 100))  # Default to 100 ms

# Validate configuration values
if not API_KEY or API_KEY == "your_default_starlink_api_key":
    raise ValueError("API key is not set. Please set the STARLINK_API_KEY environment variable.")

if THRESHOLD_LATENCY <= 0:
    raise ValueError("Threshold latency must be a positive integer.")

# Optional: Print configuration for debugging purposes
print(f"Using API Key: {API_KEY}")
print(f"Threshold Latency: {THRESHOLD_LATENCY} ms")
