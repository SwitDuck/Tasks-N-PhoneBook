import FuncFile as fun
import datetime as dt
def main():
    Android = fun.Consts
    parameters = Android(adb_ip="127.0.0.1", adb_port=5037)
    parameters.connect()
    output1, output2 = fun.NmeaRepr.represent(parameters)
    user = "user1"
    current_time = dt.datetime.now("Asia/Vladivostok")
    data = {"user1": user, "current_time": current_time, "gpgga":output1, "gprmc": output2}
    netw = fun.Network.Send_v(data=data, http_addr="localhost/gps", port=8080)
    #print(output1, output2)

if __name__ == "__main__":
    main()