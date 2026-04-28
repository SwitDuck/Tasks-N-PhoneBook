import adbutils as adb

class Consts:
    def __init__(self, adb_ip: str, adb_port: int):
        self.adb_ip = adb_ip
        self.adb_port = adb_port
        self.lat = 0
        self.lon = 0
        self.alt = 0
        self.bear = 0
    def get_floats(self):
        return self.lat, self.lon, self.alt, self.bear

    def get_gps_info(self):
        import re
        d = adb.device()
        location_info = d.shell("dumpsys location")
        for line in location_info.split("\n"):
            if "last location" in line:
                gps = re.search(r"gps\s*(-?\d+\.\d+),\s*(-?\d+\.\d+)", line)
                if gps:
                    self.lat, self.lon = map(float, gps.groups())
                parts = re.split(r"[\s=,]+", location_info)
                for i, part in enumerate(parts):
                    if part == "alt":
                        self.alt = float(parts[i + 1])
                    elif part == "bear":
                        self.bear = float(parts[i + 1])
    def get_battery_info(self):
        import re
        d = adb.device()
        location_info = d.shell("dumpsys battery")
        print(location_info)         
        
    def connect(self):
        ad = adb.AdbClient(host=self.adb_ip, port=self.adb_port, socket_timeout=100)
        print(ad.list())
        for info in ad.list():
            print(info.serial, info.state)
        try:
            output = ad.connect(f"{self.adb_ip}:{self.adb_port}")
            print(output)
        except ad.AdbTimeout as e:
            print(e)
        self.get_gps_info()
        self.get_battery_info()

class NmeaRepr:
    def dd_to_dms(deg):
        """Преобразование широты и долготы из градусов в DMS (градусы, минуты)"""
        d = int(deg)
        m = (abs(deg) - abs(d)) * 60
        return f"{abs(d):02}{m:07.4f}"
    
    @staticmethod
    def represent(const_objs: "method") -> str:
        lat, lon, alt, bear = const_objs.get_floats()
        import time
        a = NmeaRepr.dd_to_dms
    # Определение северного/южного и восточного/западного полушарий
        lat_dir = 'N' if lat >= 0 else 'S'
        lon_dir = 'E' if lon >= 0 else 'W'

        # Текущее UTC-время
        utc_time = time.strftime("%H%M%S", time.gmtime())
        utc_date = time.strftime("%d%m%y", time.gmtime())

        # Скорость и количество спутников (заглушки)
        speed_knots = 0.5  # Примерная скорость в узлах
        nsat = 10          # Количество спутников
        hdop = 1.0         # Горизонтальная точность

        gpgga = f"$GPGGA,{utc_time},{a(lat)},{lat_dir},{a(lon)},{lon_dir},1,{nsat},{hdop:.1f},{alt:.1f},M,,M,,"
        gprmc = f"$GPRMC,{utc_time},A,{a(lat)},{lat_dir},{a(lon)},{lon_dir},{speed_knots:.1f},{bear:.1f},{utc_date},,,"
        return gpgga, gprmc

class Network:
    @staticmethod
    def Send_v(data: dict, http_addr: str, port: int, endpoint: str = "/gps"):
        import requests
        url = f'http://{http_addr}:{port}{endpoint}'
        print(f"Отправка на URL: {url}")
        try:
            response = requests.post(url, json=data, timeout=5)
            print(f"Ответ сервера: {response.status_code} - {response.text}")
            return response
        except requests.exceptions.ConnectionError as e:
            print(f"Ошибка подключения: {e}")
            return None
        
            
            
