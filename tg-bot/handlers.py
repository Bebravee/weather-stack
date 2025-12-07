from loguru import logger
from telegram.ext import (
    CommandHandler,
    MessageHandler,
    ConversationHandler,
    ContextTypes,
    filters,
)

import request as req
import construct as cnst

(
    WAITING_CITY_WEATHER,
    WAITING_TSHIRT,
    WAITING_HOODIE,
    WAITING_JACKET,
    WAITING_CITY1,
    WAITING_CITY2,
    DROP_TIME,
    WAITING_CITY_NEWS,
    WAITING_CITY_CLOTHES,
    WAITING_WEATHER,
    CANCEL,
) = range(11)


async def start(update, context):
    user = update.effective_user
    logger.info(f"Команда /start от {user.username or user.first_name} (ID: {user.id})")

    welcome_msg = (
        f"Приветствую, {user.first_name}!\n\n"
        "Доступные команды:\n"
        "/register - регистрация\n"
        "/weather - прогноз погоды\n"
        "/news - последние новости\n"
        "/profile - ваш профиль\n"
        "/clothes - рекомендации по одежде\n"
        "/settemperatures - настройка комфортных температур\n"
        "/help - помощь\n"
    )
    await update.message.reply_text(welcome_msg)


async def help_command(update, context):
    user_id = update.effective_user.id
    logger.info(f"Команда /help от {user_id}")

    help_text = (
        "Доступные команды:\n"
        "/start - регистрация и начало работы\n"
        "/weather - прогноз погоды для вашего города\n"
        "/news - последние новости по интересующему городу\n"
        "/profile - информация о вашем профиле\n"
        "/clothes - рекомендации по одежде для вашего города\n"
        "/settemperatures - настройка комфортных температур\n"
    )
    await update.message.reply_text(help_text)


async def profile_command(update, context):
    user = update.effective_user
    logger.info(f"Запрос профиля {user.id}")

    profile = cnst.get_user_profile(user.id)
    if not profile:
        await update.message.reply_text("Вы не зарегистрированы. Используйте /start для регистрации.")
        return

    if isinstance(profile, dict):
        data = {
            "id": profile.get("id"),
            "name": profile.get("name"),
            "sex": profile.get("sex"),
            "age": profile.get("age"),
            "city_n": profile.get("city_n"),
            "city_w": profile.get("city_w") or profile.get("city"),
            "t_comfort": profile.get("t_comfort"),
            "t_tol": profile.get("t_tol"),
            "t_puh": profile.get("t_puh"),
        }
    elif isinstance(profile, (list, tuple)):
        data = {
            "id": profile[0] if len(profile) > 0 else None,
            "name": profile[1] if len(profile) > 1 else None,
            "sex": profile[2] if len(profile) > 2 else None,
            "age": profile[3] if len(profile) > 3 else None,
            "city_n": profile[5] if len(profile) > 5 else None,
            "city_w": profile[6] if len(profile) > 6 else None,
            "t_comfort": profile[7] if len(profile) > 7 else None,
            "t_tol": profile[8] if len(profile) > 8 else None,
            "t_puh": profile[9] if len(profile) > 9 else None,
        }
    else:
        data = {"id": profile}

    def fmt(val):
        return val if val not in (None, "", 0) else "Не указано"

    text = (
        "Ваш профиль:\n"
        f"• ID: {fmt(data.get('id'))}\n"
        f"• Имя: {fmt(data.get('name'))}\n"
        f"• Пол: {fmt(data.get('sex'))}\n"
        f"• Возраст: {fmt(data.get('age'))}\n"
        f"• Город для новостей: {fmt(data.get('city_n'))}\n"
        f"• Город для погоды: {fmt(data.get('city_w'))}\n"
        f"• Комфорт в футболке: {fmt(data.get('t_comfort'))}°C\n"
        f"• Комфорт в толстовке: {fmt(data.get('t_tol'))}°C\n"
        f"• Комфорт в пуховике: {fmt(data.get('t_puh'))}°C"
    )
    await update.message.reply_text(text)


async def register_user(update, context):
    user = update.effective_user
    logger.info(f"Старт регистрации {user.id}")

    context.user_data["registration"] = user.id
    return await clothes_command(update, context)


def _render_clothes_message(data: dict) -> str:
    base_lines = []
    if data.get("stub"):
        base_lines.append("LLM-режим выключен, показана заглушка.")
    message = data.get("message")
    if message:
        base_lines.append(str(message))
    temps = data.get("user_temps") or {}
    weather = data.get("weather_used") or {}
    extras = []
    if temps:
        extras.append(f"Ваши настройки: футболка {temps.get('comf')}°C, толстовка {temps.get('tol')}°C, пуховик {temps.get('puh')}°C.")
    if weather:
        extras.append(
            f"Погода учтена: {weather.get('temperature')}°C, ощущается {weather.get('feels_like') or weather.get('feels')}°C, {weather.get('description')}."
        )
    if extras:
        base_lines.append("\n".join(extras))
    return "\n".join(base_lines) if base_lines else "Ответ без данных."


async def clothes_command(update, context):
    user = update.effective_user
    logger.info(f"Запрос одежды от {user.id}")

    user_id = user.id
    if context.user_data.get("registration") != user_id:
        if cnst.user_exists(user_id):
            status, data = req.get_clothes_with_profile(user_id)
            if status == 200 and isinstance(data, dict):
                await update.message.reply_text(_render_clothes_message(data))
            else:
                await update.message.reply_text("Не удалось получить рекомендации. Попробуйте позже.")
            return ConversationHandler.END
        await update.message.reply_text("Нужно зарегистрироваться: используйте /register.")
        return ConversationHandler.END

    await update.message.reply_text(
        "Введите название вашего города, чтобы сохранить настройки температур для рекомендаций по одежде."
    )
    return WAITING_CITY_WEATHER


async def waiting_city_w(update, context):
    user = update.effective_user
    city = update.message.text.strip()
    logger.info(f"{user.id} ввел город: {city}")

    if context.user_data.get("registration") == user.id:
        await update.message.reply_text(
            "Записываем предпочтения по одежде.\nВведите комфортную температуру для футболки (например, 25):"
        )
        context.user_data[user.id] = [1, city, None, None, None, None, None]
        return WAITING_TSHIRT

    await update.message.reply_text("Сначала зарегистрируйтесь с помощью /register.")
    return ConversationHandler.END


async def process_input(update, context):
    user = update.effective_user
    user_id = user.id
    temp_text = update.message.text.strip()
    logger.info(f"Ввод от {user_id}: {temp_text}")

    try:
        temp_index = context.user_data[user_id][0] + 1

        if 1 < temp_index < 5:
            temp = int(temp_text)
            if not (-5 <= temp <= 40):
                await update.message.reply_text("Введите температуру в диапазоне -5..40.")
                return (
                    WAITING_TSHIRT
                    if temp_index == 2
                    else WAITING_HOODIE
                    if temp_index == 3
                    else WAITING_JACKET
                )

        if temp_index == 6:
            temp_text = temp_text.replace(".", "/").replace(":", "/").replace("-", "/")

        context.user_data[user_id][temp_index] = temp_text

        if temp_index == 2:
            await update.message.reply_text("Сохранено. Теперь комфортная температура для толстовки (например, 18):")
            context.user_data[user_id][0] = temp_index
            return WAITING_HOODIE
        elif temp_index == 3:
            await update.message.reply_text("Сохранено. Теперь комфортная температура для пуховика (например, 10):")
            context.user_data[user_id][0] = temp_index
            return WAITING_JACKET

        if context.user_data.get("registration") != user_id:
            await cnst.send_weather_success(update.message, context.user_data[user_id])
            return ConversationHandler.END

        await update.message.reply_text("Использовать этот же город для новостей? (да/нет)")
        context.user_data[user_id][0] = temp_index
        return WAITING_CITY1

    except ValueError:
        await update.message.reply_text("Введите число.")
        return (
            WAITING_TSHIRT
            if context.user_data[user_id][0] == 1
            else WAITING_HOODIE
            if context.user_data[user_id][0] == 2
            else WAITING_JACKET
        )


def yes_no(message_text: str) -> bool:
    msg = message_text.strip().lower()
    positive_words = ["yep", "да", "ага", "yes", "ok", "ок", "конечно"]
    return any(word in msg for word in positive_words)


async def set_city_news(update, context):
    user = update.effective_user
    user_id = user.id
    message_text = update.message.text.strip()

    fl = yes_no(message_text)
    current_index = context.user_data[user_id][0]

    if current_index + 1 == 4 and fl:
        context.user_data[user_id][5] = context.user_data[user_id][1]
        await update.message.reply_text("Сохранено! Когда присылать новости/погоду? (формат 10/00)")
        context.user_data[user_id][0] = current_index + 1
        return DROP_TIME
    elif not fl and current_index == 3:
        await update.message.reply_text("Введите город для новостей:")
        context.user_data[user_id][0] = current_index + 1
        return WAITING_CITY2
    elif current_index + 1 == 5:
        context.user_data[user_id][5] = message_text
        await update.message.reply_text("Сохранено! Когда присылать новости/погоду? (формат 10/00)")
        context.user_data[user_id][0] = current_index + 1
        return DROP_TIME


async def finish_registration(update, context):
    user = update.effective_user
    user_id = user.id

    try:
        await update.message.reply_text("Сохраняю данные, завершаем регистрацию...")
        context.user_data[user_id][6] = update.message.text.strip()
        data = context.user_data[user_id]

        await cnst.register(update.message, data[:])
        context.user_data.clear()
        return ConversationHandler.END

    except Exception as e:
        logger.error(f"Ошибка при завершении регистрации {user_id}: {e}")
        await update.message.reply_text("Произошла ошибка при регистрации. Попробуйте снова.")
        return ConversationHandler.END


async def weather_command(update, context):
    user = update.effective_user
    user_id = user.id
    logger.info(f"Запрос погоды от {user_id}")

    if cnst.user_exists(user_id):
        profile = cnst.get_user_profile(user_id)
        city = None
        if isinstance(profile, dict):
            city = profile.get("city_w") or profile.get("city")
        elif isinstance(profile, (list, tuple)) and len(profile) > 6:
            city = profile[6]
        try:
            status, data = req.get_weather(city)
            if status == 200 and isinstance(data, dict):
                text = (
                    f"Погода в {city}:\n"
                    f"Температура: {data.get('temperature')}°C (ощущается {data.get('feels')}°C)\n"
                    f"Описание: {data.get('description')}\n"
                    f"Влажность: {data.get('humidity')}%\n"
                    f"Давление: {data.get('pressure')} hPa\n"
                    f"Ветер: {data.get('wind_speed')} м/с"
                )
                await update.message.reply_text(text)
            else:
                raise AssertionError(f"HTTP {status}")
        except AssertionError:
            await update.message.reply_text("Не удалось получить погоду. Попробуйте позже.")
        return ConversationHandler.END

    await update.message.reply_text("Нужно завершить регистрацию, чтобы получить погоду: /register")
    return ConversationHandler.END


async def news_command(update, context):
    user = update.effective_user
    user_id = user.id
    logger.info(f"Запрос новостей от {user_id}")

    if cnst.user_exists(user_id):
        profile = cnst.get_user_profile(user_id)
        city = None
        if isinstance(profile, dict):
            city = profile.get("city_n") or profile.get("city")
        elif isinstance(profile, (list, tuple)) and len(profile) > 5:
            city = profile[5]
        try:
            status, data = req.get_news(city)
            if status == 200 and isinstance(data, dict):
                articles = data.get("articles") or []
                if not articles:
                    await update.message.reply_text(
                        f"Новостей по городу {city} не найдено за последнее время."
                    )
                else:
                    lines = [f"Свежие новости по {city}:"]
                    for i, art in enumerate(articles[:5], 1):
                        title = art.get("title") or "Без заголовка"
                        source = (art.get("source") or {}).get("name") or "Не указан"
                        url = art.get("url") or ""
                        lines.append(f"{i}. {title}\n   Источник: {source}\n   Ссылка: {url}")
                    await update.message.reply_text("\n".join(lines))
            else:
                raise AssertionError(f"HTTP {status}")
        except AssertionError:
            await update.message.reply_text("Не удалось получить новости. Попробуйте позже.")
        return ConversationHandler.END

    await update.message.reply_text("Нужно завершить регистрацию, чтобы получить новости: /register")
    return ConversationHandler.END


async def send_weather_success(message, temp_data):
    user_id = message.from_user.id
    logger.info(f"Отправка подтверждения температурных настроек {user_id}")
    success_msg = (
        "Настройки комфортных температур сохранены!\n\n"
        f"Футболка: {temp_data[2]}°C\n"
        f"Толстовка: {temp_data[3]}°C\n"
        f"Пуховик: {temp_data[4]}°C\n"
    )
    await message.reply_text(success_msg)


async def cancel(update, context):
    user = update.effective_user
    logger.info(f"Отмена операции {user.id}")

    context.user_data.clear()
    await update.message.reply_text("Операция отменена.")
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
    entry_points=[CommandHandler("weather", weather_command)],
    states=weather_state,
    fallbacks=[CommandHandler("cancel", cancel)],
)

conv_handler_clothes = ConversationHandler(
    entry_points=[CommandHandler("clothes", clothes_command)],
    states=clothes_states,
    fallbacks=[CommandHandler("cancel", cancel)],
)

conv_handler_news = ConversationHandler(
    entry_points=[CommandHandler("news", news_command)],
    states=news_state,
    fallbacks=[CommandHandler("cancel", cancel)],
)

conv_handler_register = ConversationHandler(
    entry_points=[CommandHandler("register", register_user)],
    states=reg_states,
    fallbacks=[CommandHandler("cancel", cancel)],
)


def get_handlers() -> list:
    logger.debug("Создание списка обработчиков")

    handlers = [
        CommandHandler("start", start),
        CommandHandler("help", help_command),
        CommandHandler("profile", profile_command),
        conv_handler_register,
        conv_handler_weather,
        conv_handler_clothes,
        conv_handler_news,
    ]

    logger.info(f"Создано {len(handlers)} обработчиков")
    return handlers
