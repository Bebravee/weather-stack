import sqlite3 as sq
from loguru import logger
from config import config

alf = ['name', 'sex', 'age', 'temp1', 'temp2', 'temp3', 'news', 'right']

class user_pack():
    
    def __init__(self, columns: list[int], *values):
        logger.debug(f"Создание user_pack с колонками: {columns}, значениями: {values}")
        m = [0] * 8
        for i, v in enumerate(columns): 
            m[v] = values[i]
    
        self.d = {
            'name': m[0], 'sex': m[1], 'age': m[2],
            'temp1': m[3], 'temp2': m[4], 'temp3': m[5],
            'news': m[6], 'right': m[7]
        }
        logger.debug(f"user_pack создан: {self.d}")
    
    def inone(self, col: str, value): 
        logger.debug(f"Обновление поля {col} на значение {value}")
        self.d[col] = value
    
    def inmany(self, columns: list[int], *values):
        logger.debug(f"Массовое обновление: колонки {columns}, значения {values}")
        keys = list(self.d.keys())
        for i, v in enumerate(columns): 
            self.d[keys[v]] = values[i]
    
    def unpack(self): 
        result = '(' + str(list(self.d.values()))[1:-1] + ')'
        logger.debug(f"Распакованные данные: {result}")
        return result
    
    def get_pack(self): 
        return list(self.d.values())

class session():
    
    def __init__(self): 
        logger.info(f"подключение к БД: {config.DATABASE_NAME}")
        try:
            self.con = sq.connect(config.DATABASE_NAME, check_same_thread=False)
            self.cur = self.con.cursor()
            self.create_table()
            logger.success("подключение к БД прошло норм")
        except Exception as error:
            logger.error(f"ошибка подключения к БД: {error}")
            raise

    def create_table(self):
        logger.info("проверка существования user_params")
        try:
            self.cur.execute('''
                CREATE TABLE IF NOT EXISTS user_params (
                    id INTEGER PRIMARY KEY,
                    name TEXT,
                    sex TEXT,
                    age INTEGER,
                    temp1 REAL,
                    temp2 REAL,
                    temp3 REAL,
                    news TEXT,
                    right TEXT
                )
            ''')
            self.con.commit()
            logger.success("таблица user_params готова")
        except Exception as error:
            logger.error(f"ошибка создания таблицы: {error}")
            raise

    def print_(self, sort: str, *num):
        if len(num) > 0:
            sym = ",".join(num)
        else: 
            sym = '*'
        data = 'SELECT ' + sym + ' FROM user_params' + sort
        result = self.cur.execute(data).fetchall()
        logger.debug(f"результат: {len(result)} записей")
        print(result)

    def getby_id(self, columns: list[int], idd: int) -> list:
        clms = ','.join([alf[i] for i in columns])
        data = 'SELECT ' + clms + f' FROM user_params WHERE id = {idd}'
        result = self.cur.execute(data).fetchone()
        logger.debug(f"результат запроса по id: {result}")
        return result
    
    def getby_name(self, columns: list[int], name: str) -> list:
        clms = ','.join([alf[i] for i in columns])
        data = 'SELECT ' + clms + f' FROM user_params WHERE name = "{name}"'
        result = self.cur.execute(data).fetchone()
        logger.debug(f"результат запроса по имени: {result}")
        return result

    def insrt_pack(self, pack: user_pack):
        clms = '(' + ','.join(alf) + ')'
        value = pack.unpack()
        data = f'INSERT INTO user_params {clms} VALUES {value}'
        self.cur.execute(data)
        self.con.commit()
        logger.success("запись добавлена через user_pack")
        
    def insrt_with_id(self, user_id: int, columns: list[int], *values):
        clms = 'id,' + ','.join([alf[i] for i in columns])
        value_str = str((user_id,) + values)
        data = f'INSERT INTO user_params ({clms}) VALUES {value_str}'
        self.cur.execute(data)
        self.con.commit()
        logger.success(f"запись с id {user_id} добавлена")
        
    def insrt_(self, columns: list[int], *values):
        clms = '(' + ','.join([alf[i] for i in columns]) + ')'
        value = str(values)
        data = f'INSERT INTO user_params {clms} VALUES {value}'
        self.cur.execute(data)
        self.con.commit()
        logger.success("запись успешно добавлена")

    def insrt_many_columns(self, columns: list[int], dat):
        q = ('?, ' * len(columns))[:-2]
        clms = '(' + ','.join([alf[i] for i in columns]) + ')'
        data = f'INSERT INTO user_params {clms} VALUES({q})'
        self.cur.executemany(data, dat)
        self.con.commit()
        logger.success(f"добавлено {len(dat)} записей")
        
    def insrt_many_packs(self, packs: list[user_pack]):
        clms = '(' + ','.join(alf) + ')'
        values = [pack.get_pack() for pack in packs]
        q = ('?, ' * len(alf))[:-2]
        data = f'INSERT INTO user_params {clms} VALUES({q})'
        self.cur.executemany(data, values)
        self.con.commit()
        logger.success(f"добавлено {len(packs)} packs")

    def updatecl(self, columns: list[int], user_id: int, values: list):
        set_clause = ','.join([f'{alf[col]} = ?' for col in columns])
        data = f'UPDATE user_params SET {set_clause} WHERE id = {user_id}'
        self.cur.execute(data, values)
        self.con.commit()
        logger.success(f"данные {user_id} обновлены")

    def delr(self, param: str):
        data = f'DELETE FROM user_params WHERE {param}'
        self.cur.execute(data)
        self.con.commit()
        logger.info("записи удалены")

    def user_exists(self, user_id: int) -> bool:
        data = f'SELECT id FROM user_params WHERE id = {user_id}'
        result = self.cur.execute(data).fetchone()
        exists = result is not None
        logger.debug(f"пользователь {user_id} существует: {exists}")
        return exists

    def get_user_profile(self, user_id: int):
        data = f'SELECT * FROM user_params WHERE id = {user_id}'
        result = self.cur.execute(data).fetchone()
        logger.debug(f"пользователь {user_id} получен: {result is not None}")
        return result

    def close(self):
        self.con.close()
        logger.debug("бд закрыта")
