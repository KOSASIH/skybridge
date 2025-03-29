import requests
import logging
import asyncio
import aiohttp
from aiohttp import ClientTimeout

# Setup logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class StarlinkAPI:
    def __init__(self, api_key, base_url="https://api.starlink.com/v1/"):
        self.api_key = api_key
        self.base_url = base_url
        self.session = aiohttp.ClientSession()

    async def get_connection_status(self):
        endpoint = f"{self.base_url}status"
        return await self._make_request(endpoint)

    async def _make_request(self, url):
        retries = 3
        for attempt in range(retries):
            try:
                async with self.session.get(url, headers={"Authorization": f"Bearer {self.api_key}"}, timeout=ClientTimeout(total=10)) as response:
                    response.raise_for_status()  # Raise an error for bad responses
                    data = await response.json()
                    self._validate_response(data)
                    return data
            except (aiohttp.ClientError, asyncio.TimeoutError) as e:
                logging.error(f"Request failed: {e}. Attempt {attempt + 1} of {retries}.")
                await asyncio.sleep(2 ** attempt)  # Exponential backoff
            except Exception as e:
                logging.error(f"An unexpected error occurred: {e}")
                break
        raise Exception("Max retries exceeded for request.")

    def _validate_response(self, data):
        # Validate the response structure
        if 'status' not in data:
            raise ValueError("Invalid response structure: 'status' key not found.")
        # Add more validation as needed based on expected response structure

    async def close(self):
        await self.session.close()

# Example usage
async def main():
    api_key = "your_starlink_api_key"  # Replace with your actual API key
    starlink_api = StarlinkAPI(api_key)

    try:
        connection_status = await starlink_api.get_connection_status()
        print("Connection Status:", connection_status)
    except Exception as e:
        logging.error(f"Failed to get connection status: {e}")
    finally:
        await starlink_api.close()

if __name__ == "__main__":
    asyncio.run(main())
