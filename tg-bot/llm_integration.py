import requests
import json
from loguru import logger
from config import config

class LLM:
    def __init__(self):
        self.API_KEY = config.LLM_API_KEY
        self.url = config.LLM_URL
        logger.info("инициализация LLM клиента")
    
    def ask_for_clothes(self, city, temperature, weather_conditions, humidity, wind_speed, user_temps=None):
        # LLM отключен, возвращаем заглушку без сетевых вызовов.
        logger.warning("LLM-запросы временно выключены, возвращаем заглушку")
        return (
            "ℹ️ Рекомендации по одежде временно без LLM.\n"
            f"Город: {city}, t={temperature}°C, условия: {weather_conditions}, "
            f"влажность: {humidity}%, ветер: {wind_speed} м/с.\n"
            "Как только включим LLM, сюда вернутся персональные подсказки."
        )

if __name__ == "__main__":
    llm = LLM()
    
    logger.info("тест LLM интеграции")
    result = llm.ask_for_clothes('Москва', '15', 'легкий дождь', '75', '3', [25, 18, 10])
    print(result)
