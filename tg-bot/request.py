from loguru import logger
import requests

url = "http://localhost:8080/"

def post_register_user(usr_id: str, data):
    
    address = "profile/register"
    payload = {
        "name": data[0],
        "sex": "male",
        "age": 69,
        "city_n": data[1],
        "city_w": data[5],
        "drop": data[6],
        "comf": data[2],
        "tol": data[3],
        "puh": data[4],
        "temp1": float(str(data[2]) + str(data[3]) + str(data[4])),
        "TelegramID": int(usr_id)
    }

    try:
        logger.debug(f"отправка POST запроса на {url + address}")
        response = requests.post(url + address, json=payload, timeout=10)
        logger.info(f"ответ сервера при регистрации: {response.status_code}")
        return response
    except Exception as e:
        logger.error(f"ошибка отправки данных регистрации: {e}")
        raise

def get_user_profile(usr_id: str):
    
    address = f"profile/by-telegram-id/{usr_id}"

    try:
        logger.debug(f"отправка GET запроса на {url + address}")
        response = requests.get(url + address, timeout=10)
        logger.info(f"ответ сервера при запросе профиля: {response.status_code}")
        
        user_data = None
        if response.status_code == 200:
            try:
                user_data = response.json()
                logger.debug(f"данные профиля получены: {user_data}")
            except Exception as e:
                logger.error(f"ошибка парсинга (профиль): {e}")
        else:
            logger.warning(f"неуспешный статус при запросе профиля: {response.status_code}")

        return response.status_code, user_data
        
    except Exception as e:
        logger.error(f"ошибка при запросе профиля пользователя: {e}")
        return None, None

def get_weather_with_profile(user_id: str):
    
    address = f"weather/{user_id}"

    try:
        response = requests.get(url + address, timeout=10)
        logger.info(f"ответ сервера при запросе погоды через профиль: {response.status_code}")
        
        user_data = None
        if response.status_code == 200:
            try:
                user_data = response.json()
                logger.debug("данные погоды через профиль получены")
            except Exception as e:
                logger.error(f"ошибка парсинга (погода/профиль): {e}")
        else:
            logger.warning(f"неуспешный статус при запросе погоды через профиль: {response.status_code}")

        return response.status_code, user_data
        
    except Exception as e:
        logger.error(f"ошибка при запросе погоды через профиль: {e}")
        return None, None

def get_clothes_with_profile(user_id: str):
    
    address = f"weather/clothes/{user_id}"

    try:
        response = requests.get(url + address, timeout=10)
        logger.info(f"ответ сервера при запросе одежды через профиль: {response.status_code}")
        
        user_data = None
        if response.status_code == 200:
            try:
                user_data = response.json()
                logger.debug("данные одежды через профиль получены")
            except Exception as e:
                logger.error(f"Ошибка парсинга (одежда/профиль): {e}")
        else:
            logger.warning(f"неуспешный статус при запросе одежды через профиль: {response.status_code}")

        return response.status_code, user_data
        
    except Exception as e:
        logger.error(f"ошибка при запросе одежды через профиль: {e}")
        return None, None

def get_weather(city: str):
    
    address = f"weather/city/{city}"

    try:
        response = requests.get(url + address, timeout=10)
        logger.info(f"ответ сервера при запросе погоды по городу: {response.status_code}")
        
        user_data = None
        if response.status_code == 200:
            try:
                user_data = response.json()
                logger.debug("данные погоды по городу получены")
            except Exception as e:
                logger.error(f"ошибка парсинга (погода): {e}")
        else:
            logger.warning(f"неуспешный статус при запросе погоды по городу: {response.status_code}")

        return response.status_code, user_data
        
    except Exception as e:
        logger.error(f"ошибка при запросе погоды по городу: {e}")
        return None, None

def get_clothes(city: str):
    
    address = f"weather/clothes/city/{city}"

    try:
        response = requests.get(url + address, timeout=10)
        logger.info(f"ответ сервера при запросе одежды по городу: {response.status_code}")
        
        user_data = None
        if response.status_code == 200:
            try:
                user_data = response.json()
                logger.debug("данные одежды по городу получены")
            except Exception as e:
                logger.error(f"ошибка парсинга (одежда): {e}")
        else:
            logger.warning(f"неуспешный статус при запросе одежды по городу: {response.status_code}")

        return response.status_code, user_data
        
    except Exception as e:
        logger.error(f"ошибка при запросе одежды по городу: {e}")
        return None, None

def get_news(city: str):
    
    address = f"news/city/{city}"

    try:
        response = requests.get(url + address, timeout=10)
        logger.info(f"ответ сервера при запросе новостей: {response.status_code}")
        
        user_data = None
        if response.status_code == 200:
            try:
                user_data = response.json()
                logger.debug("данные новостей получены")
            except Exception as e:
                logger.error(f"ошибка парсинга (новости): {e}")
        else:
            logger.warning(f"неуспешный статус при запросе новостей: {response.status_code}")

        return response.status_code, user_data
        
    except Exception as e:
        logger.error(f"ошибка при запросе новостей: {e}")
        return None, None
