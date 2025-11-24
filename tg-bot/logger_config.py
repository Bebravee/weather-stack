import sys
from pathlib import Path
from loguru import logger

def setup_logging():
    """бля короче это централная хуйня"""
    
    Path("logs").mkdir(exist_ok=True)
    logger.remove()
    
    logger.add(
        sys.stdout,
        format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> | <level>{level: <8}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>",
        level="INFO",
        colorize=True,
        backtrace=True,
        diagnose=True
    )
    
    logger.add(
        "logs/bot.log",
        format="{time:YYYY-MM-DD HH:mm:ss} | {level: <8} | {name}:{function}:{line} - {message}",
        level="DEBUG",
        rotation="10 MB",  
        retention="30 days",  
        compression="zip",
        backtrace=True,
        diagnose=True
        #тут короче файловый вывод максимум 10мб и храниться могут 30 дней, если это хуйня - уберёшь
    )
    
    logger.add(
        "logs/errors.log",
        format="{time:YYYY-MM-DD HH:mm:ss} | {level: <8} | {name}:{function}:{line} - {message}",
        level="ERROR",
        rotation="10 MB",
        retention="90 days",
        compression="zip",
        backtrace=True,
        diagnose=True
    )

    logger.info("Логирование:")
    
    return logger

logger = setup_logging()
