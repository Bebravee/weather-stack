import os
import json
from loguru import logger

class Config:
    def __init__(self):
        current_dir = os.path.dirname(os.path.abspath(__file__))
        json_path = os.path.join(current_dir, "secured_data.json")
        
        try:
            with open(json_path, "r", encoding="utf-8") as file:
                data = json.load(file)
                self.LLM_API_KEY = data.get("LLM_API_KEY")
                self.TELEGRAM_TOKEN = data.get("TELEGRAM_TOKEN")
                self.NEWS_API_KEY = data.get("NEWS_API_KEY")
                self.WEATHER_API_KEY = data.get("WEATHER_API_KEY")
                
            logger.info("информация загружена из secured_data.json")
            
        except FileNotFoundError:
            logger.error(f"файл не найден: {json_path}")
            self._set_defaults()
        except Exception as error:
            logger.error(f"ошибка запуска!!!: {error}")
            self._set_defaults()
    
    def _set_defaults(self):
        logger.warning("хоть какие то значения")
        self.LLM_API_KEY = None
        self.TELEGRAM_TOKEN = None
        self.NEWS_API_KEY = None
        self.WEATHER_API_KEY = None
    
    LLM_URL = "https://api.arliai.com/v1/chat/completions"
    LLM_MODEL = "(TRIAL) Llama-3.3-70B-Instruct"
    LLM_MAX_TOKENS = 500
    LLM_TEMPERATURE = 0.7
    LLM_TOP_P = 0.9

    DATABASE_NAME = "NWBot.db"
    
    WEATHER_URL = "https://api.openweathermap.org/data/2.5/weather"
    NEWS_URL = "https://newsapi.org/v2/everything"

config = Config()
