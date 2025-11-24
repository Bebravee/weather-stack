from loguru import logger
from telegram.ext import (
    CommandHandler, MessageHandler, filters, 
    ConversationHandler, ContextTypes
)
import request as req
import construct as cnst

(
    WAITING_CITY_WEATHER, WAITING_TSHIRT, WAITING_HOODIE,
    WAITING_JACKET, WAITING_CITY1, WAITING_CITY2,
    DROP_TIME, WAITING_CITY_NEWS, WAITING_CITY_CLOTHES, WAITING_WEATHER, CANCEL
) = range(11)

async def start(update, context):
    user = update.effective_user
    logger.info(f"команда /start от {user.username or user.first_name} (ID: {user.id})")
    
    welcome_msg = f'''👋 Приветствую, {user.first_name}!

Чем могу помочь?
/register - регистрация
/weather - прогноз погоды
/news - последние новости
/profile - ваш профиль
/clothes - рекомендации по одежде
/settemperatures - настройка комфортных температур
/help - помощь'''
    
    await update.message.reply_text(welcome_msg)
    logger.debug(f"сообщение отправлено {user.id}")

async def help_command(update, context):
    user_id = update.effective_user.id
    logger.info(f"команда /help от {user_id}")
    
    help_text = '''
📋 Доступные команды:

/start - регистрация и начало работы
/weather - прогноз погоды для вашего города
/news - последние новости по интересующему городу
/profile - информация о вашем профиле
/clothes - рекомендации по одежде для вашего города
/settemperatures - настройка комфортных температур

💡 Просто введите команду и следуйте инструкциям!
'''
    await update.message.reply_text(help_text)
    logger.debug(f"помощь отправлена {user_id}")

async def profile_command(update, context):
    user = update.effective_user
    logger.info(f"запрос профиля {user.id}")
    
    profile = cnst.get_user_profile(user.id)
    if not profile:
        logger.warning(f"профиль {user.id} не найден")
        await update.message.reply_text('❌ Вы не зарегистрированы. Используйте /start для регистрации.')
    else:
        profile_text = f'''
📊 Ваш профиль:

🆔 ID: {profile[0]}
👤 Имя: {profile[1] or 'Не указано'}
🚻 Пол: {profile[2] or 'Не указан'}
🎂 Возраст: {profile[3] or 'Не указан'}
👕 Комфортная температура в футболке: {profile[7] or 'Не настроена'}°C
🧥 Комфортная температура в толстовке: {profile[8] or 'Не настроена'}°C
🧥 Комфортная температура в пуховике: {profile[9] or 'Не настроена'}°C
📰 Последние новости: {profile[5] or 'Нет данных'}
'''
        await update.message.reply_text(profile_text)
        logger.info(f"профиль {user.id} отправлен")

async def register_user(update, context):
    user = update.effective_user
    logger.info(f"регистрация {user.id}")
    
    context.user_data['registration'] = user.id
    return await clothes_command(update, context)

async def clothes_command(update, context):
    user = update.effective_user
    logger.info(f"запрос по одежде от {user.id}")
    
    await update.message.reply_text('👕 Введите название вашего города для получения рекомендаций по одежде:')
    
    user_id = user.id
    if context.user_data.get('registration') != user_id:
        if cnst.user_exists(user_id):
            logger.debug(f"{user_id} существует, запрос одежды")
            data = req.get_clothes_with_profile(user_id)
    
    return WAITING_CITY_WEATHER

async def waiting_city_w(update, context):
    user = update.effective_user
    city = update.message.text
    logger.info(f"{user.id} ввел город для погоды: {city}")
    
    if context.user_data.get('registration') == user.id:
        await update.message.reply_text(
            '🌡️ Давайте настроим ваши персональные температурные предпочтения!\n\n'
            'Это поможет мне давать более точные рекомендации по одежде.\n\n'
            'Введите комфортную температуру, когда вам тепло в футболке (например, 25):'
        )
        context.user_data[user.id] = [1, city, None, None, None, None, None]
        logger.debug(f"настройка температур для {user.id}")
        return WAITING_TSHIRT
    else:
        logger.debug(f"запрос погоды для города {city} от {user.id}")
        data = cnst.get_clothes(city)
        return ConversationHandler.END

async def process_input(update, context):
    user = update.effective_user
    user_id = user.id
    temp_text = update.message.text.strip()
    
    logger.info(f"ввод от {user_id}: {temp_text}")
    
    try:
        temp_index = context.user_data[user_id][0] + 1
        
        if 1 < temp_index < 5:
            temp = int(temp_text)
            if not (-5 <= temp <= 40):
                logger.warning(f"{user_id} ввел некорректную температуру: {temp}")
                await update.message.reply_text('❌ Пожалуйста, введите температуру от -5 до 40 градусов:')
                return WAITING_TSHIRT if temp_index == 2 else WAITING_HOODIE if temp_index == 3 else WAITING_JACKET

        if temp_index == 6:
            try:
                temp_text = temp_text.replace('.', '/').replace(':', '/').replace('-', '/')
                logger.debug(f"время рассылки: {temp_text}")
            except Exception as e:
                logger.warning(f"ошибка времени, установлено значение по умолчанию: {e}")
                temp_text = '10/00'

        context.user_data[user_id][temp_index] = temp_text
        logger.debug(f"данные {user_id} обновлены: {context.user_data[user_id]}")

        if temp_index == 2:
            await update.message.reply_text('✅ Сохранено! Теперь введите комфортную температуру, когда вам тепло в толстовке (например, 18):')
            context.user_data[user_id][0] = temp_index
            return WAITING_HOODIE
        elif temp_index == 3:
            await update.message.reply_text('✅ Сохранено! Теперь введите комфортную температуру, когда вам тепло в пуховике (например, 10):')
            context.user_data[user_id][0] = temp_index
            return WAITING_JACKET
                        
        if context.user_data.get('registration') != user_id:
            await cnst.send_weather_success(update.message, context.user_data[user_id])
            logger.info(f"настройка температур завершена для {user_id}")
            return ConversationHandler.END
        else:
            await update.message.reply_text('✅ Сохранено! Вы хотите получать новости по тому же городу?')
            context.user_data[user_id][0] = temp_index
            return WAITING_CITY1

    except ValueError as e:
        logger.error(f"ошибка преобразования температуры {user_id}: {e}")
        await update.message.reply_text('❌ Пожалуйста, введите число (например, 25):')
        return WAITING_TSHIRT if context.user_data[user_id][0] == 1 else WAITING_HOODIE if context.user_data[user_id][0] == 2 else WAITING_JACKET

def yes_no(message_text: str) -> bool:
    msg = message_text.strip().lower()
    positive_words = ['yep', 'да', 'конечно', 'ага', 'yes', 'ok', 'ок', 'хорошо']
    result = any(word in msg for word in positive_words)
    logger.debug(f"анализ ответа '{message_text}': положительный = {result}")
    return result

async def set_city_news(update, context):
    user = update.effective_user
    user_id = user.id
    message_text = update.message.text
    
    logger.info(f"обработка выбора города для новостей {user_id}: {message_text}")
    
    fl = yes_no(message_text)
    current_index = context.user_data[user_id][0]
    
    if current_index + 1 == 4 and fl:
        context.user_data[user_id][5] = context.user_data[user_id][1]
        await update.message.reply_text('✅ Сохранено! Когда вы хотите получать рассылку новостей и погоды?')
        context.user_data[user_id][0] = current_index + 1
        logger.debug(f"{user_id} использует тот же город для новостей")
        return DROP_TIME
    elif not fl and current_index == 3:
        await update.message.reply_text('✅ Введите по какому городу получать новости')
        context.user_data[user_id][0] = current_index + 1
        return WAITING_CITY2
    elif current_index + 1 == 5:
        context.user_data[user_id][5] = message_text.strip()
        await update.message.reply_text('✅ Сохранено! Когда вы хотите получать рассылку новостей и погоды?')
        context.user_data[user_id][0] = current_index + 1
        logger.debug(f"{user_id} установил город для новостей: {message_text}")
        return DROP_TIME

async def finish_registration(update, context):
    user = update.effective_user
    user_id = user.id
    
    logger.info(f"завершение регистрации {user_id}")
    
    try:
        await update.message.reply_text('✅ Сохранено! Завершаем регистрацию...')
        context.user_data[user_id][6] = update.message.text.strip()
        data = context.user_data[user_id]
        
        await cnst.register(update.message, data[:])
        context.user_data.clear()
        logger.success(f"{user_id} зарегистрирован")
        return ConversationHandler.END
        
    except Exception as e:
        logger.error(f"ошибка при завершении регистрации {user_id}: {e}")
        await update.message.reply_text('❌ Произошла ошибка при регистрации. Попробуйте снова.')
        return ConversationHandler.END

async def weather_command(update, context):
    user = update.effective_user
    user_id = user.id
    
    logger.info(f"запрос погоды от {user_id}")
    
    if cnst.user_exists(user_id):
        city = cnst.get_user_profile(user_id)[6]
        logger.debug(f"{user_id} существует, город: {city}")
        try:
            r = req.get_weather(city)
            assert r.status_code == 200
            logger.info(f"погода для {user_id} отправлена")
        except AssertionError:    
            await cnst.get_weather(city)
        return ConversationHandler.END
    else:
        logger.debug(f"{user_id} не зарегистрирован, запрос города")
        await update.message.reply_text('🌤️ Введите название вашего города для получения прогноза погоды:')
        return WAITING_WEATHER

async def news_command(update, context):
    user = update.effective_user
    user_id = user.id
    
    logger.info(f"запрос новостей от {user_id}")
    
    if cnst.user_exists(user_id):
        city = cnst.get_user_profile(user_id)[5]
        logger.debug(f"{user_id} существует, город для новостей: {city}")
        try:
            r = req.get_news(city)
            assert r.status_code == 200
            logger.info(f"новости для {user_id} отправлены")
        except AssertionError:    
            await cnst.get_news(city)
        return ConversationHandler.END
    
    logger.debug(f"{user_id} не зарегистрирован, запрос города для новостей")
    await update.message.reply_text('📰 Введите интересующий вас город для поиска новостей:')
    return WAITING_CITY_NEWS

async def send_weather_success(message, temp_data):
    user_id = message.from_user.id
    logger.info(f"отправка подтверждения температурных настроек {user_id}")
    
    success_msg = f'''✅ Ваши температурные предпочтения сохранены!

👕 Футболка: {temp_data[2]}°C
🧥 Толстовка: {temp_data[3]}°C
🧥 Пуховик: {temp_data[4]}°C

Теперь рекомендации по одежде будут учитывать ваши персональные предпочтения!'''
                    
    await message.reply_text(success_msg)
    logger.debug(f"подтверждение температур отправлено {user_id}")

async def cancel(update, context):
    user = update.effective_user
    logger.info(f"отмена операции {user.id}")
    
    context.user_data.clear()
    await update.message.reply_text('❌ Операция отменена.')
    return ConversationHandler.END

reg_states = {
    WAITING_CITY_WEATHER: [MessageHandler(filters.TEXT & ~filters.COMMAND, waiting_city_w)],
    WAITING_TSHIRT: [MessageHandler(filters.TEXT & ~filters.COMMAND, process_input)],
    WAITING_HOODIE: [MessageHandler(filters.TEXT & ~filters.COMMAND, process_input)],
    WAITING_JACKET: [MessageHandler(filters.TEXT & ~filters.COMMAND, process_input)],
    WAITING_CITY1: [MessageHandler(filters.TEXT & ~filters.COMMAND, set_city_news)],
    WAITING_CITY2: [MessageHandler(filters.TEXT & ~filters.COMMAND, set_city_news)],
    DROP_TIME: [MessageHandler(filters.TEXT & ~filters.COMMAND, finish_registration)],
    CANCEL: [MessageHandler(filters.TEXT & ~filters.COMMAND, cancel)],
}

clothes_states = {
    WAITING_TSHIRT: [MessageHandler(filters.TEXT & ~filters.COMMAND, process_input)],
    WAITING_HOODIE: [MessageHandler(filters.TEXT & ~filters.COMMAND, process_input)],
    WAITING_JACKET: [MessageHandler(filters.TEXT & ~filters.COMMAND, process_input)],
}

weather_state = {
    WAITING_CITY_WEATHER: [MessageHandler(filters.TEXT & ~filters.COMMAND, waiting_city_w)],
}

news_state = {
    WAITING_CITY_NEWS: [MessageHandler(filters.TEXT & ~filters.COMMAND, set_city_news)],
}

conv_handler_weather = ConversationHandler(
    entry_points=[CommandHandler('weather', weather_command)],
    states=weather_state,
    fallbacks=[CommandHandler('cancel', cancel)],
)

conv_handler_clothes = ConversationHandler(
    entry_points=[CommandHandler("clothes", clothes_command)],
    states=clothes_states,
    fallbacks=[CommandHandler('cancel', cancel)],
)

conv_handler_news = ConversationHandler(
    entry_points=[CommandHandler('news', news_command)],
    states=news_state,
    fallbacks=[CommandHandler('cancel', cancel)],
)

conv_handler_register = ConversationHandler(
    entry_points=[CommandHandler("register", register_user)],
    states=reg_states,
    fallbacks=[CommandHandler('cancel', cancel)],
)

def get_handlers() -> list:
    logger.debug("создание списка обработчиков")
    
    handlers = [
        CommandHandler("start", start),
        CommandHandler("help", help_command),
        CommandHandler("profile", profile_command),
        conv_handler_register,
        conv_handler_weather,
        conv_handler_clothes,
        conv_handler_news,
    ]
    
    logger.info(f"создано {len(handlers)} обработчиков")
    return handlers
