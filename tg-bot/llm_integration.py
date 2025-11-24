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
        
        if not self.API_KEY:
            logger.error("апи ключ LLM не настроен")
            return "❌ API ключ не настроен. Проверьте файл secured_data.json"
        
        model = config.LLM_MODEL
        
        user_prefs = ""
        if user_temps and len(user_temps) == 3:
            tshirt_temp, hoodie_temp, jacket_temp = user_temps
            user_prefs = f"\n\nПерсональные температурные предпочтения пользователя:\n- Комфортная температура в футболке: {tshirt_temp}°C\n- Комфортная температура в толстовке: {hoodie_temp}°C\n- Комфортная температура в пуховике: {jacket_temp}°C\n\nУчти эти персональные предпочтения при составлении рекомендаций."
            logger.debug(f"учет персональных предпочтений: {user_temps}")
        
        prompt = f"""Ты — полезный ассистент, который дает рекомендации по выбору одежды исходя из погодных условий.

Погода:
- Город: {city}
- Температура: {temperature}°C
- Погодные условия: {weather_conditions}
- Влажность: {humidity}%
- Скорость ветра: {wind_speed} м/с{user_prefs}

Проанализируй погодные условия и дай практичные рекомендации по одежде. Учитывай температуру, осадки, влажность и ветер.

Ответь в следующем формате:
- Основная одежда: [рекомендация по верхней одежде]
- Обувь: [рекомендация по обуви] 
- Аксессуары: [рекомендация по аксессуарам]
- Общий совет: [краткое итоговое замечание]

Будь кратким и практичным, отвечай на русском языке:"""
        
        try:
            payload = json.dumps({
                "model": model,
                "messages": [
                    {
                        "role": "system", 
                        "content": "Ты — помощник по подбору одежды по погоде. Отвечай кратко и по делу на русском языке."
                    },
                    {
                        "role": "user", 
                        "content": f"{prompt}"
                    },
                ],
                "max_tokens": config.LLM_MAX_TOKENS,
                "temperature": config.LLM_TEMPERATURE,
                "top_p": config.LLM_TOP_P
            })
            
            headers = {
                'Content-Type': 'application/json',
                'Authorization': f"Bearer {self.API_KEY}"
            }

            logger.debug(f"отправка запроса к апи LLM: {self.url}")
            response = requests.post(self.url, headers=headers, data=payload, timeout=30)
            response_data = response.json()
            
            if response.status_code == 200 and "choices" in response_data:
                content = response_data["choices"][0]["message"]["content"]
                logger.success("рекомендации успешно получены")
                return content
            else:
                error_msg = response_data.get('message', 'Неизвестная ошибка')
                logger.error(f"ошибка LLM апи: {error_msg}, статус: {response.status_code}")
                return f"❌ Ошибка API: {error_msg}"
                
        except requests.exceptions.Timeout:
            logger.error("таймаут при запросе к апи LLM")
            return "❌ Таймаут при получении рекомендаций. Попробуйте позже."
        except Exception as e:
            logger.error(f"ошибка при запросе к LLM: {e}")
            return f"❌ Произошла ошибка при получении рекомендаций: {str(e)}"

if __name__ == "__main__":
    llm = LLM()
    
    logger.info("тест LLM интеграции")
    result = llm.ask_for_clothes('Москва', '15', 'легкий дождь', '75', '3', [25, 18, 10])
    print(result)
