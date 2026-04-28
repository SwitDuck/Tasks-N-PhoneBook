import FuncFile as fun
from datetime import datetime
from zoneinfo import ZoneInfo

def main():
    Android = fun.Consts
    parameters = Android(adb_ip="127.0.0.1", adb_port=5037)
    parameters.connect()
    output1, output2 = fun.NmeaRepr.represent(parameters)
    user = "user1"
    current_time = datetime.now(ZoneInfo('Asia/Vladivostok'))
    current_time_str = current_time.isoformat()
    
    data = {
        "user1": user, 
        "current_time": current_time_str,
        "gpgga": output1, 
        "gprmc": output2
    }
    netw = fun.Network.Send_v(data=data, http_addr="localhost", port=8080, endpoint="/gps")

if __name__ == "__main__":
    main()