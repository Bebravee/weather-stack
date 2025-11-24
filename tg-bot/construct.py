from loguru import logger
from datetime import datetime, timedelta
import requests

from llm_integration import LLM
import request as req
from database import session, user_pack
from config import config

db = session()
llm = LLM()

async def register(message, data) -> bool:
    user = message.from_user
    user_id = user.id
    username = user.username or user.first_name
    
    logger.info(f"регистрация {username} (ID: {user_id})")
    
    welcome_msg_error = '❌ Ошибка при регистрации'
    welcome_msg = f'''🎉 Добро пожаловать, {username}!
Вы успешно зарегистрированы в системе!

Доступные команды:
/weather - прогноз погоды
/news - последние новости
/profile - ваш профиль
/clothes - рекомендации по одежде
/settemperatures - настройка комфортных температур
/help - помощь'''

    if user_exists(user_id): 
        logger.warning(f"{user_id} уже зарегистрирован")
        return False

    try:
        data[0] = username
        logger.debug(f"отправка данных регистрации для {user_id}: {data}")
        r = req.post_register_user(user_id, data[:])
        
        if r.status_code != 200:
            logger.error(f"ошибка регистрации: статус {r.status_code}")
            raise AssertionError(f"HTTP {r.status_code}")
        else:
            logger.debug("регистрация успешна, запрос профиля")
            r2 = req.get_user_profile(user_id)
            try:
                assert r2[0] == 200
                logger.success(f"{username} успешно зарегистрирован")
                await message.reply_text(welcome_msg)
                return True
            except AssertionError as e:
                logger.error(f"ошибка получения профиля после регистрации: {e}")
                await message.reply_text(welcome_msg_error)
                return False
            
    except AssertionError as e:
        logger.warning(f"Ошибка при регистрации: {e}")
        USER_REGISTRATED = False
        try:
            if not db.user_exists(user_id):
                logger.info(f"регистрация {username} в локальной БД")
                try:
                    db.insrt_with_id(user_id, [0], username)
                    logger.success(f"{username} зарегистрирован в локальной БД")
                    USER_REGISTRATED = True
                except Exception as db_error:
                    logger.error(f"ошибка регистрации в локальной БД: {db_error}")
            else:
                USER_REGISTRATED = True
                logger.info(f"{username} уже существует в локальной БД")
                
            if USER_REGISTRATED:
                await message.reply_text(welcome_msg)
                return True
            else:
                raise AssertionError("не удалось зарегистрировать пользователя")
                
        except AssertionError as e:
            logger.error(f"ошибка регистрации в локальной БД: {e}")
            await message.reply_text(welcome_msg_error)
            return False
            
    except Exception as e:
        logger.error(f"непонятная ошибка при регистрации: {e}")
        await message.reply_text(welcome_msg_error)
        return False

def user_exists(user_id: str) -> bool:
    logger.debug(f"проверка существования {user_id}")
    
    try:
        r2 = req.get_user_profile(user_id)
        if r2[0] == 200:
            logger.debug(f"{user_id} существует")
            return True
        else:
            logger.debug(f"{user_id} не найден, статус: {r2[0]}")
            raise AssertionError(f"HTTP {r2[0]}")
    except AssertionError as e:
        logger.warning(f"ошибка при проверке пользователя: {e}")
        try:
            exists = db.user_exists(int(user_id))
            logger.debug(f"{user_id} существует в локальной БД: {exists}")
            return exists
        except Exception as db_error:
            logger.error(f"ошибка проверки пользователя в локальной БД: {db_error}")
            return False
    except Exception as e:
        logger.error(f"Непонятная ошибка при проверке пользователя: {e}")
        return False

def get_user_profile(user_id):
    logger.debug(f"запрос профиля {user_id}")
    
    try:
        r2 = req.get_user_profile(user_id)
        if r2[0] == 200:
            logger.debug(f"профиль {user_id} получен из апи")
            return r2[1]
        else:
            logger.warning(f"ошибка при получении профиля: {r2[0]}")
            raise AssertionError(f"HTTP {r2[0]}")
    except AssertionError as e:
        logger.warning(f"ошибка апи, используем локальную БД: {e}")
        try:
            profile = db.get_user_profile(int(user_id))
            if profile:
                logger.debug(f"профиль {user_id} получен из локальной БД")
                return profile
            else:
                logger.warning(f"профиль {user_id} не найден в локальной БД")
                return None
        except Exception as db_error:
            logger.error(f"ошибка получения профиля из локальной БД: {db_error}")
            return None
    except Exception as e:
        logger.error(f"Непонятная ошибка при получении профиля: {e}")
        return None

class WeatherStruct():
    def __init__(self, temp=15, descr='Ясно', hum=50, pres=1100, wind=2):
        self.temperature = temp
        self.description = descr
        self.humidity = hum
        self.pressure = pres
        self.wind_speed = wind
        
    def __str__(self):
        return f"""🌤️ Погода в городе:
🌡️ Температура: {self.temperature}°C
📝 Описание: {self.description}
💧 Влажность: {self.humidity}%
📊 Давление: {self.pressure} гПа
💨 Скорость ветра: {self.wind_speed} м/с"""

async def get_weather(city):
    logger.info(f"запрос погоды для города: {city}")
    
    try:
        url = f'{config.WEATHER_URL}?q={city}&units=metric&lang=ru&appid={config.WEATHER_API_KEY}'
        
        response = requests.get(url)
        weather_data = response.json()

        if response.status_code == 200:
            ws = WeatherStruct(
                weather_data['main']['temp'],
                weather_data['weather'][0]['description'],
                humidity=weather_data['main']['humidity'],
                pressure=weather_data['main']['pressure'],
                wind_speed=weather_data['wind']['speed']
            )
            
            weather_message = str(ws)
            logger.success(f"погода для {city} успешно получена")
            
            return weather_message
            
        else:
            error_msg = weather_data.get('message', 'Неизвестная ошибка')
            logger.error(f"ошибка получения погоды для {city}: {error_msg}")
            return f'❌ Ошибка: {error_msg}'
            
    except Exception as e:
        logger.error(f"ошибка при запросе погоды для {city}: {e}")
        return '❌ Произошла ошибка. Попробуйте еще раз.'

async def get_news(city):
    logger.info(f"запрос новостей для города: {city}")
    
    try:
        week_ago = (datetime.now() - timedelta(days=7)).strftime('%Y-%m-%d')
        
        url = f'{config.NEWS_URL}?q={city}&from={week_ago}&sortBy=publishedAt&language=ru&apiKey={config.NEWS_API_KEY}'
        
        response = requests.get(url)
        news_data = response.json()
        
        if response.status_code == 200 and news_data['status'] == 'ok':
            articles = news_data['articles'][:5]
            
            if not articles:
                logger.warning(f"новости для города {city} не найдены")
                return f'❌ Новости по городу {city} не найдены за последнюю неделю.'
            
            news_message = f"📰 Последние новости по городу {city}:\n\n"
            
            for i, article in enumerate(articles, 1):
                title = article['title']
                source = article['source']['name']
                published_at = datetime.strptime(article['publishedAt'], '%Y-%m-%dT%H:%M:%SZ').strftime('%d.%m.%Y %H:%M')
                url = article['url']
                
                news_message += f"{i}. **{title}**\n"
                news_message += f"   📋 Источник: {source}\n"
                news_message += f"   🕒 Дата: {published_at}\n"
                news_message += f"   🔗 [Читать полностью]({url})\n\n"
            
            logger.success(f"новости для {city} получены ({len(articles)} статей)")
            return news_message
            
        else:
            error_msg = news_data.get('message', 'Неизвестная ошибка')
            logger.error(f"ошибка получения новостей для {city}: {error_msg}")
            return f'❌ Ошибка при получении новостей: {error_msg}'
            
    except Exception as e:
        logger.error(f"ошибка при запросе новостей для {city}: {e}")
        return '❌ Произошла ошибка при получении новостей. Попробуйте еще раз.'

async def get_clothes_recommendation(city, user_temps=None):
    logger.info(f"запрос рекомендаций по одежде для города: {city}")
    
    try:
        url = f'{config.WEATHER_URL}?q={city}&units=metric&lang=ru&appid={config.WEATHER_API_KEY}'
        
        response = requests.get(url)
        weather_data = response.json()
        
        if response.status_code == 200:
            temperature = weather_data['main']['temp']
            description = weather_data['weather'][0]['description']
            humidity = weather_data['main']['humidity']
            wind_speed = weather_data['wind']['speed']
            
            logger.debug(f"погодные данные для рекомендаций: {temperature}°C, {description}")
            
            clothes_recommendation = llm.ask_for_clothes(
                city=city,
                temperature=temperature,
                weather_conditions=description,
                humidity=humidity,
                wind_speed=wind_speed,
                user_temps=user_temps
            )
            
            clothes_message = f"👕 Рекомендации по одежде для {city}:\n\n"
            clothes_message += f"🌤️ Погодные условия:\n"
            clothes_message += f"• Температура: {temperature}°C\n"
            clothes_message += f"• Описание: {description}\n"
            clothes_message += f"• Влажность: {humidity}%\n"
            clothes_message += f"• Ветер: {wind_speed} м/с\n\n"
            
            if user_temps:
                clothes_message += f"🎯 Ваши персональные настройки:\n"
                clothes_message += f"• Футболка: {user_temps[0]}°C\n"
                clothes_message += f"• Толстовка: {user_temps[1]}°C\n"
                clothes_message += f"• Пуховик: {user_temps[2]}°C\n\n"
            
            clothes_message += f"💡 Рекомендации:\n{clothes_recommendation}"
            
            logger.success(f"Рекомендации по одежде для {city} успешно сгенерированы")
            return clothes_message
            
        else:
            logger.error(f"ошибка получения погоды для рекомендаций: {response.status_code}")
            return '❌ Город не найден. Проверьте правильность написания.'
            
    except Exception as e:
        logger.error(f"ошибка при получении рекомендаций по одежде для {city}: {e}")
        return '❌ Произошла ошибка при получении рекомендаций. Попробуйте еще раз.'

async def send_weather_success(message, temp_data):
    user_id = message.from_user.id
    logger.info(f"отправка подтверждения температурных настроек user_id}")
    
    success_msg = f'''✅ Ваши температурные предпочтения сохранены!

👕 Футболка: {temp_data[2]}°C
🧥 Толстовка: {temp_data[3]}°C
🧥 Пуховик: {temp_data[4]}°C

Теперь рекомендации по одежде будут учитывать ваши персональные предпочтения!'''
                    
    await message.reply_text(success_msg)
    logger.debug(f"подтверждение температур отправлено {user_id}")
