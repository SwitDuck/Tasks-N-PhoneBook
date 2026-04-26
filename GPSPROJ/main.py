import FuncFile as fun
import datetime as dt
def main():
    Android = fun.Consts
    parameters = Android(adb_ip="127.0.0.1", adb_port=5037)
    parameters.connect()
    output1, output2 = fun.NmeaRepr.represent(parameters)
    user = "user1"
    current_time = dt.datetime.now("Asia/Vladivostok")
    netw = fun.Network.Send_v({"user1": user, "current_time": current_time, "gpgga":output1, "gprmc": output2})
    #print(output1, output2)

if __name__ == "__main__":
    main()