import logging
import os
from logging.handlers import RotatingFileHandler

class Logger:
    def __init__(self, log_file='app.log', max_bytes=5 * 1024 * 1024, backup_count=5, log_level=logging.INFO):
        self.logger = logging.getLogger(__name__)
        self.logger.setLevel(log_level)

        # Create console handler
        console_handler = logging.StreamHandler()
        console_handler.setLevel(log_level)

        # Create file handler with rotation
        file_handler = RotatingFileHandler(log_file, maxBytes=max_bytes, backupCount=backup_count)
        file_handler.setLevel(log_level)

        # Create formatter and add it to the handlers
        formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
        console_handler.setFormatter(formatter)
        file_handler.setFormatter(formatter)

        # Add the handlers to the logger
        self.logger.addHandler(console_handler)
        self.logger.addHandler(file_handler)

    def log_connection_change(self, old_connection, new_connection):
        self.logger.info(f"Connection changed from {old_connection} to {new_connection}")

    def log_error(self, message):
        self.logger.error(message)

    def log_warning(self, message):
        self.logger.warning(message)

    def log_info(self, message):
        self.logger.info(message)

    def log_debug(self, message):
        self.logger.debug(message)

# Example usage
if __name__ == "__main__":
    # Set up the logger
    log_level = os.getenv("LOG_LEVEL", "INFO").upper()
    logger = Logger(log_level=log_level)

    # Log some messages
    logger.log_info("Application started.")
    logger.log_connection_change("satellite", "terrestrial")
    logger.log_error("An error occurred.")
    logger.log_warning("This is a warning message.")
    logger.log_debug("Debugging information.")
