import tensorflow as tf
import numpy as np
import logging
import os

# Setup logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class LatencyModel:
    def __init__(self, model_path='path_to_your_model.h5'):
        self.model_path = model_path
        self.model = self.load_model()

    def load_model(self):
        try:
            model = tf.keras.models.load_model(self.model_path)
            logging.info(f"Model loaded successfully from {self.model_path}")
            return model
        except Exception as e:
            logging.error(f"Failed to load model: {e}")
            raise

    def prepare_input_data(self, connection_status):
        # Example: Normalize and reshape input data based on expected model input
        try:
            # Assuming connection_status is a dictionary with relevant features
            features = [
                connection_status.get('latency', 0),
                connection_status.get('bandwidth', 0),
                connection_status.get('signal_strength', 0),
                # Add more features as needed
            ]
            # Normalize features (example normalization)
            features = np.array(features) / np.array([1000, 1000, 100])  # Example normalization factors
            return features.reshape(1, -1)  # Reshape for model input
        except Exception as e:
            logging.error(f"Error preparing input data: {e}")
            raise

    async def predict_latency(self, connection_status):
        try:
            input_data = self.prepare_input_data(connection_status)
            predicted_latency = self.model.predict(input_data)
            return predicted_latency[0][0]  # Assuming the model outputs a single value
        except Exception as e:
            logging.error(f"Error during prediction: {e}")
            raise

    def evaluate_model(self, test_data, test_labels):
        try:
            loss, accuracy = self.model.evaluate(test_data, test_labels)
            logging.info(f"Model evaluation - Loss: {loss}, Accuracy: {accuracy}")
            return loss, accuracy
        except Exception as e:
            logging.error(f"Error during model evaluation: {e}")
            raise

# Example usage
async def main():
    model = LatencyModel('path_to_your_model.h5')
    connection_status = {
        'latency': 50,
        'bandwidth': 200,
        'signal_strength': 75,
    }
    
    try:
        predicted_latency = await model.predict_latency(connection_status)
        print(f"Predicted Latency: {predicted_latency} ms")
    except Exception as e:
        logging.error(f"Failed to predict latency: {e}")

if __name__ == "__main__":
    import asyncio
    asyncio.run(main())
